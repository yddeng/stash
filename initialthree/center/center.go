package center

import (
	"errors"
	//"fmt"
	"github.com/golang/protobuf/proto"
	fnet "github.com/sniperHW/flyfish/pkg/net"
	"go.uber.org/zap"
	"initialthree/center/constant"
	"initialthree/center/protocol"
	center_rpc "initialthree/center/rpc"
	"initialthree/cluster/addr"
	"initialthree/codec/ss"
	"initialthree/common"
	"initialthree/network"
	"initialthree/pkg/event"
	"initialthree/pkg/rpc"
	"net"
	"reflect"
	"sync/atomic"
	"time"
)

type MsgHandler func(*fnet.Socket, proto.Message)

type node struct {
	addr          addr.Addr
	session       *fnet.Socket
	exportService uint32
	heartBeatTime time.Time
}

type Center struct {
	queue     *event.EventQueue
	handlers  map[uint16]MsgHandler
	nodes     map[addr.LogicAddr]*node
	rpcServer *rpc.RPCServer
	logger    *zap.Logger
	l         net.Listener
	stoped    int32
	started   int32
	selfAddr  string
	die       chan struct{}
}

func (this *Center) tick() {
	ticker := time.NewTicker(time.Second)
	for {
		now := <-ticker.C
		err := this.queue.PostNoWait(0, func() {
			for logicAddr, node := range this.nodes {
				now = time.Now()
				if now.Unix()-node.heartBeatTime.Unix() > constant.HeartBeatTimeout {
					this.removeNode(logicAddr)
				}
			}
		})
		if nil != err {
			return
		}
	}
}

func (this *Center) removeNode(logicAddr addr.LogicAddr) {
	msg := &protocol.NodeLeave{
		Nodes: []uint32{uint32(logicAddr)},
	}
	this.logger.Sugar().Infof("removeNode %s\n", logicAddr.String())
	delete(this.nodes, logicAddr)
	this.brocast(msg, nil)
}

func (this *Center) brocast(msg proto.Message, exclude *map[*fnet.Socket]bool) {
	for _, v := range this.nodes {
		if v.session != nil {
			if nil == exclude {
				this.logger.Sugar().Infof("brocast notify to %s %v\n", v.addr.Logic.String(), msg)
				v.session.Send(msg)
			} else {
				if _, ok := (*exclude)[v.session]; !ok {
					this.logger.Sugar().Infof("%s brocast notify to %s %v\n", this.selfAddr, v.addr.Logic, msg)
					v.session.Send(msg)
				}
			}
		}
	}
}

func (this *Center) onSessionClose(session *fnet.Socket, reason error) {
	ud := session.GetUserData()
	if nil != ud {
		n := ud.(*node)
		//removeNode(n.addr.Logic)
		this.logger.Sugar().Infof("node lose connection %s reason:%s\n", n.addr.Logic.String(), reason.Error())
		n.session = nil
	}
}

func (this *Center) onLogin(replyer *rpc.RPCReplyer, req interface{}) {

	session := replyer.GetChannel().(*center_rpc.RPCChannel).GetSession()

	ud := session.GetUserData()
	if nil != ud {
		return
	}

	login := req.(*protocol.Login)

	var n *node

	logicAddr := addr.LogicAddr(login.GetLogicAddr())

	netAddr, err := net.ResolveTCPAddr("tcp", login.GetNetAddr())

	if nil != err {
		loginRet := &protocol.LoginRet{
			ErrCode: proto.Int32(constant.LoginFailed),
			Msg:     proto.String("invaild netAddr:" + login.GetNetAddr()),
		}
		replyer.Reply(loginRet, nil)
		this.logger.Sugar().Info(loginRet.GetMsg())
		session.Close(errors.New("invaild netAddr"), time.Second)
		return
	}

	n, ok := this.nodes[logicAddr]
	if !ok {
		n = &node{
			addr: addr.Addr{
				Logic: logicAddr,
				Net:   netAddr,
			},
			exportService: login.GetExportService(),
		}
		this.nodes[logicAddr] = n
		this.logger.Sugar().Infof("add new node %s exportService=%v\n", logicAddr.String(), 1 == n.exportService)
	}

	if n.session != nil {
		//重复登录
		loginRet := &protocol.LoginRet{
			ErrCode: proto.Int32(constant.LoginFailed),
			Msg:     proto.String("duplicate node:" + logicAddr.String()),
		}
		replyer.Reply(loginRet, nil)
		this.logger.Sugar().Info(loginRet.GetMsg())
		session.Close(errors.New("duplicate node"), time.Second)
		return
	}

	n.session = session
	n.addr.Net = netAddr //更新网络地址
	n.heartBeatTime = time.Now()
	session.SetUserData(n)

	this.logger.Sugar().Infof("%s onLogin:%s exportService=%v\n", this.selfAddr, logicAddr.String(), 1 == n.exportService)

	loginRet := &protocol.LoginRet{
		ErrCode: proto.Int32(constant.LoginOK),
	}

	replyer.Reply(loginRet, nil)

	//记录日志
	nodeAdd := &protocol.NodeAdd{}

	nodeAdd.Nodes = append(nodeAdd.Nodes, &protocol.NodeInfo{
		LogicAddr:     proto.Uint32(uint32(n.addr.Logic)),
		NetAddr:       proto.String(n.addr.Net.String()),
		ExportService: proto.Uint32(n.exportService),
	})

	//将新节点的信息通告给除自己以外的其它节点
	exclude := map[*fnet.Socket]bool{}
	exclude[session] = true
	this.brocast(nodeAdd, &exclude)

	//将所有节点信息(包括自己)发给新到节点

	notify := &protocol.NotifyNodeInfo{}

	for _, v := range this.nodes {
		notify.Nodes = append(notify.Nodes, &protocol.NodeInfo{
			LogicAddr:     proto.Uint32(uint32(v.addr.Logic)),
			NetAddr:       proto.String(v.addr.Net.String()),
			ExportService: proto.Uint32(v.exportService),
		})
	}

	err = session.Send(notify)
	if nil != err {
		this.logger.Sugar().Error(err)
	} else {
		this.logger.Sugar().Infof("%s send notify to %s %v\n", this.selfAddr, logicAddr.String(), notify)
	}

}

func (this *Center) onHeartBeat(session *fnet.Socket, msg proto.Message) {
	ud := session.GetUserData()
	if nil == ud {
		this.logger.Sugar().Infof("onHeartBeat,session is not login\n")
		return
	}
	node := ud.(*node)
	node.heartBeatTime = time.Now()
	//kendynet.Infoln("onHeartBeat", node)

	heartbeat := msg.(*protocol.HeartbeatToCenter)
	resp := &protocol.HeartbeatToNode{}
	resp.TimestampBack = proto.Int64(heartbeat.GetTimestamp())
	resp.Timestamp = proto.Int64(time.Now().UnixNano())
	err := session.Send(resp)
	if nil != err {
		this.logger.Sugar().Errorf("send error:%s\n", err.Error())
	}
}

func (this *Center) onRemoveNode(session *fnet.Socket, msg proto.Message) {
	this.logger.Sugar().Info("onRemoveNode")
	removeNode := msg.(*protocol.RemoveNode)
	for _, v := range removeNode.Nodes {
		this.removeNode(addr.LogicAddr(v))
	}
}

func (this *Center) dispatchMsg(session *fnet.Socket, msg *ss.Message) {
	if nil != msg {
		data := msg.GetData()
		switch data.(type) {
		case *rpc.RPCRequest:

			center_rpc.OnRPCRequest(this.rpcServer, session, data.(*rpc.RPCRequest))

		//case *rpc.RPCResponse:
		//	center_rpc.OnRPCResponse(rpcServer.data.(*rpc.RPCResponse))
		case proto.Message:
			cmd := msg.GetCmd()
			handler, ok := this.handlers[cmd]
			if ok {
				handler(session, msg.GetData().(proto.Message))
			} else {
				//记录日志
				this.logger.Sugar().Errorf("unknow cmd:%s\n", cmd)
			}
		default:
			this.logger.Sugar().Errorf("invaild message type:%s \n", reflect.TypeOf(data).String())
		}
	}
}

func (this *Center) registerHandler(cmd uint16, handler MsgHandler) {
	if nil == handler {
		//记录日志
		this.logger.Sugar().Errorf("Register %d failed: handler is nil\n", cmd)
		return
	}
	_, ok := this.handlers[cmd]
	if ok {
		//记录日志
		this.logger.Sugar().Errorf("Register %d failed: duplicate handler\n", cmd)
		return
	}
	this.handlers[cmd] = handler
}

func (this *Center) Stop() {
	if atomic.CompareAndSwapInt32(&this.stoped, 0, 1) {
		this.l.Close()
		this.queue.PostNoWait(0, func() {
			for _, v := range this.nodes {
				if nil != v.session {
					v.session.Close(errors.New("center stop"), 0)
				}
			}
		})
		this.queue.Close()
		<-this.die
	}
}

func (this *Center) Start(service string, l *zap.Logger) error {

	if !atomic.CompareAndSwapInt32(&this.started, 0, 1) {
		return errors.New("already started")
	}

	this.selfAddr = service
	this.logger = l
	this.handlers = map[uint16]MsgHandler{}
	this.nodes = map[addr.LogicAddr]*node{}
	this.rpcServer = center_rpc.NewServer()
	this.die = make(chan struct{})
	go this.tick()

	center_rpc.RegisterMethod(this.rpcServer, &protocol.Login{}, this.onLogin)

	this.registerHandler(protocol.CENTER_HeartbeatToCenter, this.onHeartBeat)

	this.registerHandler(protocol.CENTER_RemoveNode, this.onRemoveNode)

	this.queue = event.NewEventQueue()

	go func() {
		this.queue.Run()
		this.die <- struct{}{}
	}()

	//启动本地监听

	server, serve, err := network.Listen("tcp", service, func(conn net.Conn) {
		session := fnet.NewSocket(conn, fnet.OutputBufLimit{})
		session.SetRecvTimeout(common.HeartBeat_Timeout * time.Second)
		session.SetInBoundProcessor(protocol.NewReceiver())
		session.SetEncoder(protocol.NewEncoder())
		session.SetCloseCallBack(func(sess *fnet.Socket, reason error) {

			this.queue.PostNoWait(0, this.onSessionClose, sess, reason)
		}).BeginRecv(func(s *fnet.Socket, m interface{}) {
			msg := m.(*ss.Message)
			this.queue.PostNoWait(0, this.dispatchMsg, session, msg)
		})
	})

	if nil == err {
		this.l = server
		this.logger.Sugar().Infof("server running on:%s\n", service)
		go serve()
	} else {
		this.logger.Sugar().Errorf("center failed %s\n", err.Error())
	}
	return err
}

func New() *Center {
	return &Center{}
}

var defaultCenter *Center = New()

func Stop() {
	defaultCenter.Stop()
}

func Start(service string, l *zap.Logger) error {
	return defaultCenter.Start(service, l)
}
