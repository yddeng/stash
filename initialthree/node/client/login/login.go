package login

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	fnet "github.com/sniperHW/flyfish/pkg/net"
	codecs "initialthree/codec/cs"
	"initialthree/cs"
	"initialthree/node/client/dispatcher"
	"initialthree/pkg/event"
	"initialthree/protocol/cmdEnum"
	cs_msg "initialthree/protocol/cs/message"
	"sync/atomic"
)

var processQueue *event.EventQueue

type loginCallback func(*Session, *dispatcher.Dispatcher, *codecs.Message, error)

type Session struct {
	session *fnet.Socket
	SeqNo   uint32
}

func (this *Session) Send(message proto.Message) error {
	seqno := atomic.AddUint32(&this.SeqNo, 1)
	//fmt.Println("------", seqno, reflect.TypeOf(message), message)
	return this.session.Send(codecs.NewMessage(seqno, message))
}

func (this *Session) SeqNoSend(seqNo uint32, message proto.Message) error {
	//fmt.Println("------", seqNo, reflect.TypeOf(message), message)
	return this.session.Send(codecs.NewMessage(seqNo, message))
}

func onGameEstablish(session *fnet.Socket, userID string, token string, gameDispatcher *dispatcher.Dispatcher, lcallback loginCallback) {
	fmt.Println("onGameEstablish")

	gameLoginToS := &cs_msg.GameLoginToS{
		UserID:   proto.String(userID),
		Token:    proto.String(token),
		ServerID: proto.Int32(1),
	}
	err := session.Send(codecs.NewMessage(uint32(0), gameLoginToS))
	fmt.Println("gameLoginToS", gameLoginToS, err)
	if nil != err {
		fmt.Printf("send error:%s\n", err.Error())
		lcallback(nil, nil, nil, err)
		return
	}

	sess := &Session{
		session: session,
		SeqNo:   0,
	}
	stop := make(chan bool)

	gameDispatcher.RegisterOnce(cmdEnum.CS_GameLogin, func(session *fnet.Socket, msg *codecs.Message) {

		code := msg.GetErrCode()
		if code != 0 {
			panic(fmt.Sprintf("gameLogin err %s", cs_msg.ErrCode(code).String()))
		} else {
			fmt.Println("------ GameLogin ok")
			lcallback(sess, gameDispatcher, msg, nil)
		}

	})
	gameDispatcher.RegisterOnce(cmdEnum.CS_Kick, func(session *fnet.Socket, msg *codecs.Message) {
		fmt.Printf("--- --- --- --- --- user on Kick\n")
		close(stop)
	})

}

func onLoginEstablish(session *fnet.Socket, userID string, loginDispatcher *dispatcher.Dispatcher, lcallback loginCallback) {
	fmt.Println("onLoginEstablish")

	loginToS := &cs_msg.LoginToS{
		UserID: proto.String(userID),
	}
	err := session.Send(codecs.NewMessage(uint32(0), loginToS))
	if nil != err {
		fmt.Printf("send error:%s\n", err.Error())
		return
	}

	loginDispatcher.RegisterOnce(cmdEnum.CS_Login, func(session *fnet.Socket, msg *codecs.Message) {

		loginToC := msg.GetData().(*cs_msg.LoginToC)
		fmt.Println("LoginToC", loginToC)

		game := loginToC.GetGame()
		token := loginToC.GetToken()

		fmt.Println(msg.GetErrCode(), game, token)

		if msg.GetErrCode() == 0 {
			gameDispatcher := dispatcher.New(processQueue)
			gameDispatcher.RegisterOnce("Establish", func(session *fnet.Socket) {
				onGameEstablish(session, userID, token, gameDispatcher, lcallback)
			})
			cs.DialTcp(game, 0, gameDispatcher)
		} else {
			fmt.Printf("login failed: %s", cs_msg.ErrCode(msg.GetErrCode()).String())
		}
	})
}

func init() {

	processQueue = event.NewEventQueue()
	go func() {
		processQueue.Run()
		fmt.Println("queue break-------------")
	}()
}

func Login(userID, addr string, callback loginCallback) {

	fmt.Println(userID, addr)

	loginDispatcher := dispatcher.New(processQueue)
	loginDispatcher.RegisterOnce("Establish", func(session *fnet.Socket) {
		onLoginEstablish(session, userID, loginDispatcher, callback)
	})
	cs.DialTcp(addr, 0, loginDispatcher)
}
