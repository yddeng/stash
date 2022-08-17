package cluster

import (
	"errors"
	"github.com/golang/protobuf/proto"
	fnet "github.com/sniperHW/flyfish/pkg/net"
	"initialthree/cluster/addr"
	"initialthree/cluster/priority"
	"initialthree/codec/ss"
	"initialthree/common"
	event2 "initialthree/pkg/event"
	"initialthree/pkg/rpc"
	"reflect"
	"time"
)

func (this *Cluster) onEstablishClient(end *endPoint, session *fnet.Socket) {

	logger.Sugar().Infof("%v onEstablishClient\n", this.serverState.selfAddr.Logic)
	end.dialing = false

	if nil != end.session {
		/*
		 * 如果end.conn != nil 表示两端同时请求建立连接，本端已经作为服务端成功接受了对端的连接
		 */
		logger.Sugar().Infof("endPoint:%s already have connection\n", end.addr.Logic.String())
		session.Close(errors.New("duplicate endPoint connection"), 0)
	} else {
		//不主动触发心跳，超时回收连接
		/*
			RegisterTimer(time.Second*(common.HeartBeat_Timeout/2), func(t *timer.Timer, _ interface{}) {
				heartbeat := &Heartbeat{}
				heartbeat.Timestamp1 = proto.Int64(time.Now().UnixNano())
				heartbeat.OriginSender = proto.Uint32(uint32(selfAddr.Logic))
				if kendynet.ErrSocketClose == session.Send(heartbeat) {
					t.Cancel()
				}
			}, nil)
		*/

		this.onEstablish(end, session)
	}
}

func (this *Cluster) onSessionEvent(end *endPoint, session *fnet.Socket, event interface{}, updateLastActive bool) {
	if updateLastActive {
		end.lastActive = time.Now()
	}

	switch event.(type) {
	case *ss.Message:

		var err error

		msg := event.(*ss.Message)

		from := msg.From()
		if addr.LogicAddr(0) == from {
			from = end.addr.Logic
		}

		data := msg.GetData()
		switch data.(type) {
		case *rpc.RPCRequest:
			err = this.queue.PostFullReturn(priority.LOW, this.rpcMgr.onRPCRequest, end, from, data.(*rpc.RPCRequest))
		case *rpc.RPCResponse:
			//response回调已经被hook必然在queue中调用
			this.rpcMgr.onRPCResponse(data.(*rpc.RPCResponse))
		case proto.Message:
			err = this.queue.PostFullReturn(priority.LOW, this.dispatch, from, session, msg.GetCmd(), data.(proto.Message))
		default:
			logger.Sugar().Errorf("invaild message type:%s \n", reflect.TypeOf(data).String())
		}

		if err == event2.ErrQueueFull {
			logger.Sugar().Errorf("event queue full discard message\n")
		}

		break
	case *ss.RelayMessage:
		if this.serverState.selfAddr.Logic.Type() == harbarType {
			this.onRelayMessage(event.(*ss.RelayMessage))
		}
		break
	default:
		logger.Sugar().Errorf("invaild message type\n")
		break
	}
}

func (this *Cluster) onEstablishServer(end *endPoint, session *fnet.Socket) {
	logger.Sugar().Infof("%v onEstablishServer\n", this.serverState.selfAddr.Logic)
	if end.addr.Logic != this.serverState.selfAddr.Logic {
		this.onEstablish(end, session)
	} else {
		//自连接server
		session.SetInBoundProcessor(ss.NewReceiver("ss", "rpc_req", "rpc_resp", this.serverState.selfAddr.Logic))
		session.SetEncoder(ss.NewEncoder("ss", "rpc_req", "rpc_resp"))
		session.BeginRecv(func(s *fnet.Socket, e interface{}) {
			this.onSessionEvent(end, s, e, false)
		})
	}
}

func (this *Cluster) onEstablish(end *endPoint, session *fnet.Socket) {

	logger.Sugar().Infof("(self:%v) %v onEstablish %v <---> %v %v %v\n", this.serverState.selfAddr.Logic, end, this.serverState.selfAddr.Logic.String(), end.addr.Logic.String(), session.LocalAddr(), session.RemoteAddr())

	end.session = session
	session.SetInBoundProcessor(ss.NewReceiver("ss", "rpc_req", "rpc_resp", this.serverState.selfAddr.Logic))
	session.SetEncoder(ss.NewEncoder("ss", "rpc_req", "rpc_resp"))
	session.SetCloseCallBack(func(sess *fnet.Socket, reason error) {
		logger.Sugar().Infof("%s disconnected error:%s %v %v\n", end.addr.Logic.String(), reason, sess.LocalAddr(), sess.RemoteAddr())
		this.queue.PostNoWait(priority.MID, this.onPeerDisconnected, end.addr.Logic, reason)
		this.rpcMgr.onEndDisconnected(end)
	}).SetErrorCallBack(func(s *fnet.Socket, reason error) {
		end.closeSession(reason)
	}).BeginRecv(func(s *fnet.Socket, e interface{}) {
		this.onSessionEvent(end, s, e, true)
	})

	now := time.Now()

	end.lastActive = now

	//自连接不需要建立timer
	if end.addr.Logic != this.serverState.selfAddr.Logic {
		end.timer = this.RegisterTimer(common.HeartBeat_Timeout/2, end.onTimerTimeout, nil)
	}

	pendingMsg := end.pendingMsg
	end.pendingMsg = []interface{}{}
	for _, v := range pendingMsg {
		session.Send(v)
	}

	for v := end.pendingCall.Front(); v != nil; v = end.pendingCall.Front() {
		c := end.pendingCall.Remove(v).(*rpcCall)
		c.listElement = nil
		if c.dialTimer.Cancel() {
			remain := c.deadline.Sub(now)
			err := this.rpcMgr.client.AsynCall(&RPCChannel{
				to:      c.to,
				peer:    end,
				cluster: this,
			}, "rpc", c.arg, remain, c.cb)

			if nil != err {
				c.cb(nil, err)
			}
		}
	}
}
