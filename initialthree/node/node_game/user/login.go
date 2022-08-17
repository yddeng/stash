package user

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	flyfish "github.com/sniperHW/flyfish/client"
	"github.com/sniperHW/flyfish/errcode"
	"initialthree/cluster"
	"initialthree/cluster/addr"
	"initialthree/cluster/loadb"
	cluster_proto "initialthree/cluster/proto"
	"initialthree/network/smux"
	"initialthree/node/common/db"
	"initialthree/node/common/serverType"
	"initialthree/node/common/taskPool"
	"initialthree/node/node_game/module"
	"initialthree/node/node_game/temporary"
	"initialthree/pkg/timer"
	"initialthree/protocol/cs/message"
	"initialthree/zaplogger"
	"os"
	"runtime"
	"time"
)

const userMaxConnection = 10000

var lbCollector *loadb.LBCollector

func StartLBCollector() {
	lbCollector = loadb.New(userMaxConnection)
	lbCollector.SetPidCollector(os.Getpid(), loadb.ModeAvg)
	cluster.RegisterTimerOnce(time.Millisecond, lbReportLoop, nil)
}

func lbReportLoop(_ *timer.Timer, ctx interface{}) {
	gates, _ := cluster.Select(serverType.Gate)
	if len(gates) != 0 {
		priority, count := lbCollector.Get(len(userMap))
		count /= len(gates)
		msg := &cluster_proto.LbReport{
			Logic:    proto.Uint32(uint32(cluster.SelfAddr().Logic)),
			Priority: proto.Uint32(uint32(priority)),
			Capacity: proto.Uint32(uint32(count)),
		}
		cluster.Brocast(serverType.Gate, msg)
	}
	cluster.RegisterTimerOnce(cluster.LBReportDur, lbReportLoop, nil)
}

func loginErrorLogout(user *User) {
	if !user.checkStatus(status_logout) {
		user.setStatus(status_logout)
		set := db.GetFlyfishClient("game").CompareAndSetNx("user_game_login", user.userID, "gameaddr", cluster.SelfAddr().Logic.String(), "")
		set.AsyncExec(func(ret *flyfish.ValueResult) {
			zaplogger.GetSugar().Infof("%s clear db logout code:%s ", user.userID, errcode.GetErrorDesc(ret.ErrCode))
			// 清理map
			deleteUser(user)
		})
	}
}

func onLogin(stream *smux.MuxStream, req *message.GameLoginToS, callback func(u *User, code message.ErrCode)) {
	zaplogger.GetSugar().Infof("stream %d login  %v", stream.ID(), req)

	userID := req.GetUserID()
	user := userMap[userID]
	if user != nil {
		if user.ServerID != req.GetServerID() {
			if user.checkStatus(status_login) {
				user.signalKick = true
			} else if user.checkStatus(status_playing, status_wait_connect) {
				user.kick(false)
			}
			zaplogger.GetSugar().Infof("%s %s", userID, "login failed 1 serverID")
			callback(nil, message.ErrCode_RETRY)

		} else if user.checkStatus(status_wait_connect) {
			stream.SetUserData(user)
			user.stream = stream
			callback(user, message.ErrCode_OK)
			user.loginOK()

		} else if user.checkStatus(status_playing) {
			// 踢掉gate，game直接设置
			user.sendKickToC()
			user.stream.SetUserData(nil)

			user.stream = stream
			stream.SetUserData(user)
			callback(user, message.ErrCode_OK)
			user.loginOK()

		} else {
			zaplogger.GetSugar().Infof("%s %s %d ", userID, "login failed 2 status:", user.status)
			callback(nil, message.ErrCode_RETRY)
		}
	} else {
		// check
		if len(userMap) > userMaxConnection {
			zaplogger.GetSugar().Infof("%s login failed, userMaxConnection", userID)
			callback(nil, message.ErrCode_RETRY)
			return
		}

		user = New(stream, userID, req.GetServerID())
		stream.SetUserData(user)
		user.setStatus(status_login)
		addUser(user)

		if err := loginSyncThreadPool.AddTask(func() {
			code := message.ErrCode_OK
			defer func() {
				cluster.PostTask(func() {
					switch code {
					case message.ErrCode_OK:
						if user.signalKick {
							//登录流程中踢人
							callback(nil, message.ErrCode_RETRY)
							user.logout(false)
						} else {
							if user.stream == nil {
								// 登陆流程中，收到gate断开链接
								//user.setStatus(status_playing)
								callback(nil, message.ErrCode_RETRY)
								user.setWaitReconnect()
							} else {
								callback(user, message.ErrCode_OK)
								user.loginOK()
							}
						}
					default:
						callback(nil, code)
						loginErrorLogout(user)
					}
				})
			}()

			//向数据库添加登录标记，如果成功继续后续流程，否则通知客户端重试
			ret := db.GetFlyfishClient("game").CompareAndSetNx("user_game_login", userID, "gameaddr", "", cluster.SelfAddr().Logic.String()).Exec()
			switch errcode.GetCode(ret.ErrCode) {
			case errcode.Errcode_ok, errcode.Errcode_cas_not_equal:
				if errcode.GetCode(ret.ErrCode) == errcode.Errcode_cas_not_equal {
					oldAddrStr := ret.Value.GetString()
					if cluster.SelfAddr().Logic.String() != oldAddrStr {
						zaplogger.GetSugar().Infof("%s %s %s %s", userID, "login failed 3 gameaddr", errcode.GetErrorDesc(ret.ErrCode), oldAddrStr)
						code = message.ErrCode_RETRY
						gameAddr, err := addr.MakeLogicAddr(oldAddrStr)
						if err == nil {
							sendKickGameUser(gameAddr, userID)
						}
						return
					}
				}

				// 拉取玩家基础数据
				if err := user.loadGameID(); err != nil {
					zaplogger.GetSugar().Infof("%s loadGameID failed %s", userID, err)
					code = message.ErrCode_ERROR
					return
				}
				// 模块数据
				if err := user.loadModuleData(); err != nil {
					zaplogger.GetSugar().Infof("%s loadModuleData failed %s", userID, err)
					code = message.ErrCode_ERROR
					return
				}
				// 加载临时数据
				temporary.Load(user)

			default:
				zaplogger.GetSugar().Infof("%s %s %s", userID, "login failed 1", errcode.GetErrorDesc(ret.ErrCode))
				code = message.ErrCode_ERROR
			}
		}); err != nil {
			zaplogger.GetSugar().Infof("%s %s", userID, err)
			callback(nil, message.ErrCode_RETRY)
			deleteUser(user)
		}
	}
}

func (this *User) loginOK() {
	zaplogger.GetSugar().Infof("%s %s %d", this.userID, "loginOK", this.status)
	this.setStatus(status_playing)
	this.removeLastTimer()
	this.DoLoadPipeline()
}

func (this *User) loadGameID() error {
	key := fmt.Sprintf("%s:%d", this.userID, this.ServerID)
	ret := db.GetFlyfishClient("game").Get("game_user", key, "id").Exec()
	zaplogger.GetSugar().Infof("load game_user %s", this.userID)

	switch errcode.GetCode(ret.ErrCode) {
	case errcode.Errcode_ok:
		this.SetID(uint64(ret.Fields["id"].GetInt()))
		zaplogger.GetSugar().Debugf("loadGameID %s %d ", this.userID, this.GetID())

	case errcode.Errcode_record_notexist:
		zaplogger.GetSugar().Infof("user:%s first login", this.userID)
	default:
		return fmt.Errorf("flyfis code:%s", errcode.GetErrorDesc(ret.ErrCode))
	}
	return nil
}

func (this *User) loadModuleData() error {
	var (
		readOutCluster = &readOutCluster{
			readOuts: map[string]*readOutClusterEle{},
		}
		modules = make(map[module.ModuleType]module.ModuleI, module.End)
	)

	for mt, creator := range module.Modules {
		m := creator(this)
		modules[mt] = m

		if this.GetID() != 0 {
			cmd := m.ReadOut()
			if cmd != nil {
				readOutCluster.addCommand(cmd)
			}
		} else {
			m.Init(nil)
		}
	}

	if this.GetID() == 0 {
		this.modules = modules
		return nil
	}

	ret := make(chan error, 1)

	readOutCluster.do(this, func() {

		rollback := false
		readOut := readOutCluster.readOuts["user_module_data"+":"+this.GetIDStr()]
		if readOut != nil && errcode.GetCode(readOut.result.ErrCode) == errcode.Errcode_record_notexist {
			rollback = true
		}

		if errcode.GetCode(readOut.result.ErrCode) != errcode.Errcode_ok {
			zaplogger.GetSugar().Debugf("user_module_data errcode: %s", errcode.GetErrorDesc(readOut.result.ErrCode))
		}

		for _, v := range readOutCluster.readOuts {
			if !(errcode.GetCode(v.result.ErrCode) == errcode.Errcode_ok || errcode.GetCode(v.result.ErrCode) == errcode.Errcode_record_notexist) {
				zaplogger.GetSugar().Errorf("user %s init: %s in %s", this.GetUserID(), errcode.GetErrorDesc(v.result.ErrCode), v.table)
				ret <- fmt.Errorf("user %s init: %s in %s", this.GetUserID(), errcode.GetErrorDesc(v.result.ErrCode), v.table)
				return
			}

			var fields map[string]*flyfish.Field
			if !rollback {
				fields = v.result.Fields
			}

			for _, m := range v.modules {
				if err := m.Init(fields); err != nil {
					zaplogger.GetSugar().Errorf("user %s init: module %s init %s", this.GetUserID(), m.ModuleType().String(), err)
					ret <- fmt.Errorf("user %s init: module %s init %s", this.GetUserID(), m.ModuleType().String(), err)
					return
				}
			}
		}

		this.modules = modules
		ret <- nil
	})

	return <-ret
}

var loginSyncThreadPool *taskPool.TaskPool

func init() {
	loginSyncThreadPool = taskPool.NewTaskPool(runtime.NumCPU()*10, 1024)
}
