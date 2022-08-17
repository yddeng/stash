package node_login

import (
	fnet "github.com/sniperHW/flyfish/pkg/net"
	"initialthree/cluster"
	codecs "initialthree/codec/cs"
	"initialthree/zaplogger"
)

type handler func(*fnet.Socket, *codecs.Message)

var gDispatcher dispatcher

type dispatcher struct {
	handlers map[uint16]handler
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

func (this *dispatcher) Dispatch(session *fnet.Socket, msg *codecs.Message) {
	if nil != msg {
		cmd := msg.GetCmd()
		handler, ok := this.handlers[cmd]
		if ok {
			//交给cluster的任务队列单线程执行
			cluster.PostTask(func() {
				handler(session, msg)
			})
		}
	}
}

func (this *dispatcher) OnClose(session *fnet.Socket, reason error) {
	zaplogger.GetSugar().Infof("client close:%v", reason)
}

func (this *dispatcher) OnNewClient(session *fnet.Socket) {
	zaplogger.GetSugar().Infof("new client")
}

func (this *dispatcher) OnAuthenticate(session *fnet.Socket) bool {
	//remoteAddr := session.RemoteAddr()
	//if err := firewall.AuthIP(remoteAddr.(*net.TCPAddr).IP.String()); err != nil {
	//	zaplogger.GetSugar().Infof("authenticate remote addr %s: %s", remoteAddr.String(), err)
	//	return false
	//}

	return true
}
