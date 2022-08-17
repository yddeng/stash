package cs

import (
	fnet "github.com/sniperHW/flyfish/pkg/net"
	codecs "initialthree/codec/cs"
	"initialthree/common"
	"initialthree/network"
	_ "initialthree/protocol/cs" //触发pb注册
	cs_proto "initialthree/protocol/cs/message"
	"time"
)

var (
	sessions map[*fnet.Socket]time.Time
	queue    chan func()
)

func DialTcp(peerAddr string, timeout time.Duration, dispatcher ClientDispatcher) {
	go func() {
		conn, err := network.Dial("tcp", peerAddr, timeout)
		if err != nil {
			dispatcher.OnConnectFailed(peerAddr, err)
		} else {
			session := network.CreateSession(conn)
			queue <- func() {
				sessions[session] = time.Now().Add(common.HeartBeat_Timeout_Client / 2)
				//session.SetRecvTimeout(common.HeartBeat_Timeout_Client)
				session.SetInBoundProcessor(codecs.NewReceiver("sc"))
				session.SetEncoder(codecs.NewEncoder("cs"))
				session.SetCloseCallBack(func(sess *fnet.Socket, reason error) {
					queue <- func() {
						delete(sessions, sess)
					}
					dispatcher.OnClose(sess, reason)
				})
				dispatcher.OnEstablish(session)
				session.BeginRecv(func(s *fnet.Socket, m interface{}) {
					msg := m.(*codecs.Message)
					switch msg.GetData().(type) {
					case *cs_proto.HeartbeatToC:
						//fmt.Printf("on HeartbeatToC\n")
						break
					default:
						dispatcher.Dispatch(session, msg)
						break
					}
				})
			}
		}
	}()
}

func init() {
	sessions = make(map[*fnet.Socket]time.Time)
	queue = make(chan func(), 10000)

	go func() {
		for {
			f := <-queue
			f()
		}
	}()

	go func() {
		for {
			queue <- func() {
				now := time.Now()
				for k, v := range sessions {
					if !now.Before(v) {
						sessions[k] = now.Add(common.HeartBeat_Timeout_Client / 2)
						//发送心跳
						Heartbeat := &cs_proto.HeartbeatToS{}
						k.Send(codecs.NewMessage(0, Heartbeat))
					}
				}
			}
			time.Sleep(time.Millisecond * 1000)
		}
	}()
}
