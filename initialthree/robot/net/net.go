package net

import (
	codecs "initialthree/codec/cs"
	"initialthree/common"
	"initialthree/network"
	"initialthree/pkg/event"
	csmsg "initialthree/protocol/cs/message"
	"sync/atomic"
	"time"

	"github.com/GodYY/gutils/assert"
	fnet "github.com/sniperHW/flyfish/pkg/net"
)

const (
	new        = 0
	connecting = 1
	docallback = 2
	done       = 3
	stopped    = 4
)

type Connector struct {
	state      int32
	service    string
	cb         func(*fnet.Socket, error)
	eventQueue *event.EventQueue
}

func NewConnecter(service string, cb func(*fnet.Socket, error), eventQueue *event.EventQueue) *Connector {
	assert.Assert(cb != nil, "cb nil")
	assert.Assert(eventQueue != nil, "eventQueue nil")

	return &Connector{
		state:      new,
		service:    service,
		cb:         cb,
		eventQueue: eventQueue,
	}
}

func (c *Connector) Connect() {
	if atomic.CompareAndSwapInt32(&c.state, new, connecting) {
		go func() {
			conn, err := network.Dial("tcp", c.service, 10*time.Second)

			if atomic.CompareAndSwapInt32(&c.state, connecting, docallback) {
				if err != nil {
					c.eventQueue.PostNoWait(1, c.onConnect, nil, err)
					return
				}

				session := network.CreateSession(conn)
				session.SetRecvTimeout(common.HeartBeat_Timeout_Client)
				session.SetInBoundProcessor(codecs.NewReceiver("sc"))
				session.SetEncoder(codecs.NewEncoder("cs"))
				c.eventQueue.PostNoWait(1, c.onConnect, session, nil)
			} else if conn != nil {
				conn.Close()
			}
		}()
	}
}

func (c *Connector) onConnect(session *fnet.Socket, err error) {

	if atomic.CompareAndSwapInt32(&c.state, docallback, done) {
		c.cb(session, err)
	} else if session != nil {
		session.Close(nil, 0)
	}
}

func (c *Connector) Stop() {
	state := atomic.LoadInt32(&c.state)
	if state != new && state != done {
		atomic.StoreInt32(&c.state, stopped)
	}
}

func ConnectService(service string, cb func(*fnet.Socket, error), eventQueue *event.EventQueue) *Connector {
	connect := NewConnecter(service, cb, eventQueue)
	connect.Connect()
	return connect

	// go func() {
	// 	conn, err := network.Dial("tcp", service, 10*time.Second)
	// 	if err != nil {
	// 		eventQueue.PostNoWait(1, cb, nil, err)
	// 		return
	// 	}

	// 	session := socket.NewStreamSocket(conn)
	// 	session.SetRecvTimeout(common.HeartBeat_Timeout_Client)
	// 	session.SetInBoundProcessor(codecs.NewReceiver("sc"))
	// 	session.SetEncoder(codecs.NewEncoder("cs"))
	// 	eventQueue.PostNoWait(1, cb, session, nil)
	// }()
}

func IsMessageOK(msg *codecs.Message) bool {
	return msg.GetErrCode() == uint16(csmsg.ErrCode_OK)
}

func GetErrCodeStr(errCode uint16) string {
	return csmsg.ErrCode(errCode).String()
}

func GetMsgErrcodeStr(msg *codecs.Message) string {
	return csmsg.ErrCode(msg.GetErrCode()).String()
}
