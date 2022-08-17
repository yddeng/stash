package dispatcher

import (
	"fmt"
	fnet "github.com/sniperHW/flyfish/pkg/net"
	codecs "initialthree/codec/cs"
	"initialthree/pkg/event"
	"initialthree/protocol/cs/message"
	"reflect"
)

type handler func()

type Dispatcher struct {
	evHandler    *event.EventHandler
	processQueue *event.EventQueue
}

func (this *Dispatcher) Emit(e interface{}, args ...interface{}) {
	this.evHandler.EmitToEventQueue(event.EventQueueParam{
		Q: this.processQueue,
	}, e, args...)
}

func (this *Dispatcher) RegisterOnce(ev interface{}, callback interface{}) {
	switch ev.(type) {
	case string:
		this.evHandler.RegisterOnce(ev, callback)
		break
	case uint16:
		this.evHandler.RegisterOnce(ev, callback)
		break
	default:
		break
	}
}

func (this *Dispatcher) Register(ev interface{}, callback interface{}) {
	switch ev.(type) {
	case string:
		this.evHandler.Register(ev, callback)
		break
	case uint16:
		this.evHandler.Register(ev, callback)
		break
	default:
		fmt.Println(reflect.TypeOf(ev).String())
		break
	}
}

func (this *Dispatcher) Dispatch(session *fnet.Socket, msg *codecs.Message) {
	if nil != msg {
		cmd := msg.GetCmd()
		if message.ErrCode(msg.GetErrCode()) != message.ErrCode_OK {
			fmt.Printf("(seqNo:%d cmd:%d) errCode %s\n", msg.GetSeriNo(), msg.GetCmd(), message.ErrCode(msg.GetErrCode()).String())
			return
		}
		this.Emit(cmd, session, msg)
	}
}

func (this *Dispatcher) OnClose(session *fnet.Socket, reason error) {
	fmt.Println("OnClose", reason)
}

func (this *Dispatcher) OnEstablish(session *fnet.Socket) {
	fmt.Println("OnEstablish")
	this.Emit("Establish", session)
}

func (this *Dispatcher) OnConnectFailed(peerAddr string, err error) {
	fmt.Println("OnConnectFailed", err)
	this.Emit("ConnectFailed", peerAddr, err)
}

func New(processQueue *event.EventQueue) *Dispatcher {
	return &Dispatcher{
		processQueue: processQueue,
		evHandler:    event.NewEventHandler(),
	}
}
