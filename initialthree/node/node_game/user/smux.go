package user

import (
	"github.com/golang/protobuf/proto"
	"initialthree/cluster"
	"initialthree/cluster/addr"
	"initialthree/codec/cs"
	"initialthree/common"
	"initialthree/network/smux"
	"initialthree/protocol/cmdEnum"
	"initialthree/protocol/cs/message"
	"initialthree/zaplogger"
	"net"
	"sync"
	"time"
)

var (
	gateMap  = map[addr.LogicAddr]*smuxGate{}
	gateLock sync.Mutex
)

type smuxGate struct {
	logic      addr.LogicAddr
	smuxSocket *smux.MuxSocket
}

func onNewSmuxConn(logic addr.LogicAddr, conn net.Conn) {
	zaplogger.GetSugar().Info("new SmuxConn ", logic)
	gateLock.Lock()
	defer gateLock.Unlock()

	gate, ok := gateMap[logic]
	if !ok {
		socket := smux.NewMuxSocketServer(conn, common.HeartBeat_Timeout, encoder.EnCode, func(muxSocket *smux.MuxSocket) {
			zaplogger.GetSugar().Infof("SmuxConn %s closed", logic.String())
			gateLock.Lock()
			defer gateLock.Unlock()
			if g, ok := gateMap[logic]; ok && g.smuxSocket == muxSocket {
				delete(gateMap, logic)
			}

		})
		socket.Listen(onNewStream)

		gate = &smuxGate{logic: logic, smuxSocket: socket}
		gateMap[logic] = gate
	} else {
		conn.Close()
	}
}

func onNewStream(stream *smux.MuxStream) {
	zaplogger.GetSugar().Info("new stream ", stream.ID())
	stream.SetCloseCallback(func(stream *smux.MuxStream, e error) {
		zaplogger.GetSugar().Info("stream close ", stream.ID(), e)
		cluster.PostTask(func() {
			if u, ok := stream.GetUserData().(*User); ok && u.stream == stream {
				// 用户已经在线，再登陆一次。直接替换了链接，故这里不相等
				zaplogger.GetSugar().Debugf("user %s onClient disconnect %v", u.userID, e)
				u.stream = nil
				if u.checkStatus(status_playing) {
					u.setWaitReconnect()
				}
			}
		})
	})

	stream.SetRecvTimeout(time.Second * 5)
	stream.Recv(func(stream *smux.MuxStream, bytes []byte) {

		ret, err := receiver.DirectUnpack(bytes)
		if err != nil {
			zaplogger.GetSugar().Error(err)
			stream.Close(err)
			return
		}

		msg := ret.(*cs.Message)
		if msg.GetCmd() != cmdEnum.CS_GameLogin {
			zaplogger.GetSugar().Errorf("first command %d, not login", msg.GetCmd())
			stream.Close(nil)
			return
		}

		cluster.PostTask(func() {
			onLogin(stream, msg.GetData().(*message.GameLoginToS), func(u *User, code message.ErrCode) {
				var err error
				//zaplogger.GetSugar().Infof("onlogin callback code %s %d %s", code.String(), u.GetID(), u.GetName())
				if code == message.ErrCode_OK {
					toc := &message.GameLoginToC{
						IsFirstLogin: proto.Bool(u.GetID() == 0 || u.GetName() == ""),
					}
					err = stream.AsyncSend(cs.NewMessage(msg.GetSeriNo(), toc), func(stream *smux.MuxStream, e error) {
						if e != nil {
							stream.Close(e)
							return
						}
						stream.SetRecvTimeout(common.HeartBeat_Timeout_Client)
						u.recv()
					}, time.Second*2)
				} else {
					err = stream.AsyncSend(cs.ErrMessage(msg.GetSeriNo(), msg.GetCmd(), uint16(code)), func(stream *smux.MuxStream, e error) {
						stream.Close(err)
					}, time.Second*2)
				}
				if err != nil {
					zaplogger.GetSugar().Error(err)
					stream.Close(err)
				}
			})
		})
	})
}

func (this *User) recv() {
	if this.stream != nil {
		this.stream.Recv(func(stream *smux.MuxStream, bytes []byte) {
			ret, err := receiver.DirectUnpack(bytes)
			if err != nil {
				zaplogger.GetSugar().Error(err)
				stream.Close(err)
				return
			}
			cluster.PostTask(func() {
				if this.stream == stream {
					msg := ret.(*cs.Message)
					switch msg.GetCmd() {
					case cmdEnum.CS_Heartbeat:
						this.Reply(msg.GetSeriNo(), &message.HeartbeatToC{})
						this.recv()
					default:
						TransDispatch(this, msg)
					}
				}
			})
		})
	}
}

func init() {
	cluster.SetNewMuxConn(onNewSmuxConn)
}
