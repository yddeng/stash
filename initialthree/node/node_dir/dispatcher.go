package node_dir

import (
	"errors"
	fnet "github.com/sniperHW/flyfish/pkg/net"
	"initialthree/cluster"
	codecs "initialthree/codec/cs"
	"initialthree/zaplogger"
	"reflect"
	"runtime"
)

//连接到服务器的session
var dirSession *fnet.Socket

type handler func(*fnet.Socket, *codecs.Message)

var gDispatcher dispatcher

type dispatcher struct {
	handlers       map[uint16]handler
	defaultHandler handler
}

func (this *dispatcher) SetDefaultHandler(h handler) {
	this.defaultHandler = h
}

func (this *dispatcher) Register(cmd uint16, h handler) {
	if nil == h {
		return
	}
	_, ok := this.handlers[cmd]
	if ok {
		return
	}

	this.handlers[cmd] = h
}

func pcall(handler_ handler, cmd uint16, session *fnet.Socket, msg *codecs.Message) {
	defer func() {
		if r := recover(); r != nil {
			buf := make([]byte, 65535)
			l := runtime.Stack(buf, false)
			zaplogger.GetSugar().Errorf("error on Dispatch:%d\nstack:%v,%s", cmd, r, buf[:l])
		}
	}()
	handler_(session, msg)
}

func (this *dispatcher) Dispatch(session *fnet.Socket, msg *codecs.Message) {
	if nil != msg {
		cmd := msg.GetCmd()
		handler, ok := this.handlers[cmd]
		if ok {
			//交给cluster的任务队列单线程执行
			cluster.PostTask(func() {
				pcall(handler, cmd, session, msg)
			})
		} else {
			cluster.PostTask(func() {
				pcall(this.defaultHandler, cmd, session, msg)
			})
		}
	}
}

func (this *dispatcher) OnClose(session *fnet.Socket, reason error) {
	u := session.GetUserData()
	if u != nil {
		u.(interface {
			OnDisConnect()
		}).OnDisConnect()
	}
	zaplogger.GetSugar().Debugf("client close:%v\n", reason)
}

func (this *dispatcher) OnNewClient(session *fnet.Socket) {
	zaplogger.GetSugar().Debug("new client\n")
}

func (this *dispatcher) OnAuthenticate(session *fnet.Socket) bool {
	//remoteAddr := session.RemoteAddr()
	//if err := firewall.AuthIP(remoteAddr.(*net.TCPAddr).IP.String()); err != nil {
	//	zaplogger.GetSugar().Infof("authenticate remote addr %s: %s", remoteAddr.String(), err)
	//	return false
	//}
	return true
}

func RegisterHandler(cmd uint16, h handler) {
	gDispatcher.Register(cmd, h)
}

func init() {
	gDispatcher.handlers = map[uint16]handler{}
	gDispatcher.defaultHandler = func(session *fnet.Socket, message *codecs.Message) {
		zaplogger.GetSugar().Errorf("defaultHandler dispatcher %v:%s", message, reflect.ValueOf(message.GetData()).String())
		session.Close(errors.New("invalid message"), 0)
	}
}
