package main

import (
	"fmt"
	"github.com/gogo/protobuf/proto"
	fnet "github.com/sniperHW/flyfish/pkg/net"
	codecs "initialthree/codec/cs"
	"initialthree/pkg/golog"
	connector "initialthree/pkg/socket/connector/tcp"
	"initialthree/protocol/cmdEnum"
	cs_msg "initialthree/protocol/cs/message"
	"initialthree/util"
	"os"
	"reflect"
	"sync/atomic"
	"time"
)

type Handler func(*Session, *codecs.Message)

type Dispatcher struct {
	handlers map[uint16]Handler
}

func NewDispatch() *Dispatcher {
	return &Dispatcher{
		handlers: map[uint16]Handler{},
	}
}

func (this *Dispatcher) Register(cmd uint16, callback Handler) {
	_, ok := this.handlers[cmd]
	if ok {
		return
	}

	this.handlers[cmd] = callback
}

func (this *Dispatcher) Dispatch(lSession *Session, msg *codecs.Message) {
	if nil != msg {
		cmd := msg.GetCmd()
		handler, ok := this.handlers[cmd]
		if ok {
			handler(lSession, msg)
		}
	}
}

func dial(name, addr string, dispatcher *Dispatcher) *Session {
	c, _ := connector.New("tcp", addr)
	sTime := time.Now()
	logger.Infof("dial %s addr %s start conn\n", name, addr)

	var sess *Session
	lSession, err := c.Dial(0)
	if nil != err {
		logger.Infof("dial %s addr %s, err %s\n", name, addr, err)
		return nil
	} else {
		eTime := time.Now()
		logger.Infof("dial %s addr %s conn ok\n", name, addr)
		logger.Infof("-- dial %s startTime %s stopTime %s, useTime %s\n", name, sTime.String(), eTime.String(), eTime.Sub(sTime).String())
		sess = &Session{lSession: lSession, name: name}
		lSession.SetReceiver(codecs.NewReceiver("sc"))
		lSession.SetEncoder(codecs.NewEncoder("cs"))
		lSession.SetCloseCallBack(func(sess *fnet.Socket, reason string) {
			//logger.Infof("%s addr %s, onClose %s\n", name, addr, reason)
		})
		lSession.Start(func(event *kendynet.Event) {
			if event.EventType == kendynet.EventTypeError {
				event.Session.Close(event.Data.(error).Error(), 0)
			} else {
				msg := event.Data.(*codecs.Message)
				dispatcher.Dispatch(sess, msg)
			}
		})
	}

	return sess
}

type Session struct {
	lSession *fnet.Socket
	name     string
	SeqNo    uint32
}

func (this *Session) Send(message proto.Message) error {
	seqno := atomic.AddUint32(&this.SeqNo, 1)
	logger.Infof("%s send message %s %v\n", this.name, reflect.TypeOf(message).String(), message)
	return this.lSession.Send(codecs.NewMessage(seqno, message))
}

func dir(addr string) {
	dispatcher := NewDispatch()
	dispatcher.Register(cmdEnum.CS_ServerList, func(lSession *Session, message *codecs.Message) {
		msg := message.GetData().(*cs_msg.ServerListToC)
		logger.Infof("dir read message ServerListToC %v\n", msg)

		server := msg.GetServerList()[0]
		loginAddr := server.GetServerAddr()
		serverID = server.GetServerId()
		loginF(loginAddr)
	})

	sess := dial("dir", addr, dispatcher)
	sess.Send(&cs_msg.ServerListToS{
		UserID: proto.String(userID),
	})
}

func loginF(addr string) {
	dispatcher := NewDispatch()
	dispatcher.Register(cmdEnum.CS_Login, func(lSession *Session, message *codecs.Message) {
		msg := message.GetData().(*cs_msg.LoginToC)
		logger.Infof("login read msg LoginToC %v\n", msg)

		gameAddr := msg.GetGame()
		token := msg.GetToken()
		gate(gameAddr, token)
	})
	sess := dial("login", addr, dispatcher)
	sess.Send(&cs_msg.LoginToS{
		UserID: proto.String(userID),
	})
}

func gate(addr, token string) {
	dispatcher := NewDispatch()
	dispatcher.Register(cmdEnum.CS_GameLogin, func(lSession *Session, message *codecs.Message) {
		msg := message.GetData().(*cs_msg.GameLoginToC)
		logger.Infof("gate read msg GameLoginToC %v\n", msg)

		if msg.GetIsFirstLogin() {
			req := &cs_msg.CreateRoleToS{
				Name: proto.String(userID),
			}
			lSession.Send(req)
		}
	})

	dispatcher.Register(cmdEnum.CS_CreateRole, func(lSession *Session, message *codecs.Message) {
		msg := message.GetData().(*cs_msg.CreateRoleToC)
		logger.Infof("gate read msg CreateRoleToC %v\n", msg)
	})

	dispatcher.Register(cmdEnum.CS_AttrSync, func(lSession *Session, message *codecs.Message) {
		msg := message.GetData().(*cs_msg.AttrSyncToC)
		logger.Infof("gate read msg AttrSyncToC %v\n", msg)
		stopTime = time.Now()

		logger.Infof("-- all end, startTime %s stopTime %s, useTime %s", startTime.String(), stopTime.String(), stopTime.Sub(startTime).String())
	})

	sess := dial("gate", addr, dispatcher)
	sess.Send(&cs_msg.GameLoginToS{
		UserID:   proto.String(userID),
		Token:    proto.String(token),
		ServerID: proto.Int32(serverID),
	})
}

var userID string
var serverID int32
var logger golog.LoggerI
var startTime time.Time
var stopTime time.Time

func main() {
	if len(os.Args) < 3 {
		fmt.Printf("usage dirAddr userID\n")
		return
	}
	dirAddr := os.Args[1]
	userID = os.Args[2]

	logger = util.NewLogger("log", "connect", 1024)

	startTime = time.Now()
	dir(dirAddr)

	select {}

}
