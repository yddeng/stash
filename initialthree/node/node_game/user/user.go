package user

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"initialthree/cluster"
	codecs "initialthree/codec/cs"
	"initialthree/common"
	"initialthree/network/smux"
	"initialthree/node/common/transaction"
	"initialthree/node/node_game/module"
	"initialthree/node/node_game/monitor"
	"initialthree/node/node_game/temporary"
	vendor_event "initialthree/pkg/event"
	"initialthree/pkg/timer"
	"initialthree/protocol/cmdEnum"
	"initialthree/protocol/cs/message"
	"initialthree/zaplogger"
	"strconv"
	"sync/atomic"
	"time"
)

const (
	status_login        = 1 //正在登录
	status_playing      = 2 //正常游戏中
	status_wait_connect = 3 //连接已经中断，等待重连
	status_logout       = 4 //正在登出
	status_wait_remove  = 5 //所有收尾处理已经完毕，待数据回写完毕后移除
)

var userMap = map[string]*User{}
var userIDMap = map[uint64]*User{}

var holdTime = common.ReconnectTime
var durationSaveTime = 10 * time.Second

var receiver = codecs.NewReceiver("cs")
var encoder = codecs.NewEncoder("sc")

const forwardCount = 64

func addUser(u *User) {
	_, ok := userMap[u.userID]
	if !ok {
		userMap[u.userID] = u
	}
}

func addIDUser(u *User) {
	_, ok := userIDMap[u.GetID()]
	if !ok {
		userIDMap[u.GetID()] = u
	}
}

func deleteUser(u *User) {
	if _, ok := userMap[u.userID]; ok {
		delete(userMap, u.userID)
	}

	if _, ok := userIDMap[u.GetID()]; ok {
		delete(userIDMap, u.GetID())
	}

}

func GetUser(uId string) *User {
	return userMap[uId]
}

func GetIDUser(gID uint64) *User {
	return userIDMap[gID]
}

type User struct {
	userID    string
	gameID    uint64
	gameIDStr string
	ServerID  int32 // 玩家选服的ID

	stream    *smux.MuxStream
	sendCount int32

	status        uint32
	saveP         saveProcessor // c
	modules       map[module.ModuleType]module.ModuleI
	evHandler     *vendor_event.EventHandler                       // c
	transMgr      *transaction.TransactionMgr                      // c
	temporaryData map[temporary.TemporaryType]temporary.TemporaryI // 临时数据 // c

	lastTimer    *timer.Timer //清理玩家时间
	nextSave     time.Time    // c
	signalKick   bool
	loadPipeline *userLoadPipeline
}

func New(stream *smux.MuxStream, userID string, serverId int32) *User {
	u := &User{
		userID:        userID,
		stream:        stream,
		ServerID:      serverId,
		transMgr:      transaction.New(cluster.GetEventQueue()),
		modules:       map[module.ModuleType]module.ModuleI{},
		nextSave:      time.Now().Add(durationSaveTime),
		temporaryData: map[temporary.TemporaryType]temporary.TemporaryI{},
		evHandler:     vendor_event.NewEventHandler(),
	}

	u.saveP = saveProcessor{
		u:       u,
		pending: make([]bool, 0, 4),
	}

	// 临时数据
	//u.temporaryData[temporary.TempSeqCache] = temporary.NewSeqCache(u)

	u.transMgr.SetCallTransEnd(u.callTransEnd)

	return u
}

func (this *User) GetUserID() string {
	return this.userID
}

func (this *User) SetID(id uint64) {
	this.gameID = id
	addIDUser(this)
}

func (this *User) GetID() uint64 {
	return this.gameID
}

func (this *User) GetIDStr() string {
	if this.gameIDStr == "" && this.GetID() != 0 {
		this.gameIDStr = strconv.FormatUint(this.GetID(), 10)
	}
	return this.gameIDStr
}

func (this *User) GetUserLogName() string {
	return fmt.Sprintf("(%s, %d)", this.GetUserID(), this.GetID())
}

func (this *User) GetSubModule(moduleType module.ModuleType) module.ModuleI {
	return this.modules[moduleType]
}

func (this *User) GetTemporary(tt temporary.TemporaryType) temporary.TemporaryI {
	return this.temporaryData[tt]
}

func (this *User) SetTemporary(tt temporary.TemporaryType, cache temporary.TemporaryI) {
	this.temporaryData[tt] = cache
}

func (this *User) ClearTemporary(tt temporary.TemporaryType) {
	delete(this.temporaryData, tt)
}

func (this *User) RegisterEvent(event interface{}, fn interface{}) vendor_event.Handle {
	return this.evHandler.Register(event, fn)
}

func (this *User) RegisterEventOnce(event interface{}, fn interface{}) vendor_event.Handle {
	return this.evHandler.RegisterOnce(event, fn)
}

func (this *User) UnRegisterEvent(h vendor_event.Handle) {
	this.evHandler.Remove(h)
}

func (this *User) ClearEvent(event interface{}) {
	this.evHandler.Clear(event)
}

func (this *User) EmitEvent(event interface{}, args ...interface{}) {
	this.evHandler.Emit(event, args...)
}

func (this *User) removeLastTimer() {
	if this.lastTimer != nil {
		this.lastTimer.Cancel()
		this.lastTimer = nil
	}
}

func (this *User) setStatus(status uint32) {
	this.status = status
}

func (this *User) StatusOk() bool {
	return this.checkStatus(status_playing)
}

func (this *User) checkStatus(status ...uint32) bool {
	for _, v := range status {
		if v == this.status {
			return true
		}
	}
	return false
}

func (this *User) FlushAllDirtyToClient() {
	for _, v := range this.modules {
		v.FlushDirtyToClient()
	}
}

func (this *User) FlushAllToClient() {
	for _, v := range this.modules {
		v.FlushAllToClient()
	}
}

func (this *User) increaseSendCount(max int32) bool {
	for {
		sendCount := atomic.LoadInt32(&this.sendCount)
		if sendCount >= max {
			return false
		}

		if atomic.CompareAndSwapInt32(&this.sendCount, sendCount, sendCount+1) {
			return true
		}
	}
}

func (this *User) sendToClient(forwordMsg *codecs.Message) {
	if this.checkStatus(status_playing) && this.stream != nil {
		if !this.increaseSendCount(int32(forwardCount)) {
			zaplogger.GetSugar().Errorf("uer.User: sendToClient sendCount full %s %d  ", this.GetUserID(), this.sendCount)
		} else {
			if err := this.stream.AsyncSend(forwordMsg, func(stream *smux.MuxStream, e error) {
				atomic.AddInt32(&this.sendCount, -1)
				if e != nil {
					zaplogger.GetSugar().Error("uer.User: sendToClient callback error ", e)
				}
			}); err != nil {
				atomic.AddInt32(&this.sendCount, -1)
				zaplogger.GetSugar().Error("uer.User: sendToClient error ", err)
			}
		}
	}
}

func (this *User) Reply(seqNo uint32, msg proto.Message) {
	forwordMsg := codecs.NewMessage(seqNo, msg)
	this.sendToClient(forwordMsg)
}

func (this *User) ReplyErr(seqNo uint32, cmd uint16, errCode message.ErrCode) {
	forwordMsg := codecs.ErrMessage(seqNo, cmd, uint16(errCode))
	this.sendToClient(forwordMsg)
}

func (this *User) Post(msg proto.Message) {
	forwordMsg := codecs.NewMessage(0, msg)
	this.sendToClient(forwordMsg)
}

func (this *User) setWaitReconnect() {
	// gate上实体必然销毁
	//this.stream = nil

	// login ： 不改变状态，由login逻辑触发
	// 设置定时器，如果玩家在定时器到期前没有重连上来销毁对象
	this.lastTimer = cluster.RegisterTimerOnce(holdTime, func(t *timer.Timer, _ interface{}) {
		if t == this.lastTimer {
			zaplogger.GetSugar().Infof("%s %s", this.userID, "lastTimer timeout")
			this.kick(false)
		}
	}, nil)

	// 短线后，临时数据处理
	for _, v := range this.temporaryData {
		v.UserDisconnect()
	}

	this.setStatus(status_wait_connect)
}

func (this *User) Tick(t *timer.Timer, _ interface{}) {

	if !this.checkStatus(status_wait_remove, status_login) {
		now := time.Now()

		//if now.After(this.nextSave) {
		this.durationSave()
		//	this.nextSave = now.Add(durationSaveTime + (time.Duration(rand.Int()%6) * time.Second))
		//}

		// 事物
		this.transMgr.Tick(now)

		// 模块
		for _, v := range this.modules {
			v.Tick(now)
			v.FlushDirtyToClient()
		}

		// 临时数据
		for _, v := range this.temporaryData {
			v.Tick(now)
		}

		cluster.RegisterTimerOnce(time.Second, this.Tick, nil)
	}
}

func init() {
	cluster.Register(cmdEnum.SS_KickGameUser, onKickGameUser)
	monitor.RegisterOnlinePlayerFunc(func() float64 {
		ch := make(chan float64, 1)
		cluster.PostTask(func() {
			ch <- float64(len(userMap))
		})
		return <-ch
	})

	/*cluster.SetPeerDisconnected(func(peer addr.LogicAddr, err error) {
		if peer.Type() == serverType.Gate {
			//for _, u := range userMap {
			//	if u.gateInfo != nil && u.gateInfo.GateAddr == peer {
			//		u.setWaitReconnect()
			//	}
			//}
		}
	})*/
}
