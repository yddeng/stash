package cluster

import (
	"fmt"
	"github.com/sniperHW/flyfish/pkg/buffer"
	fnet "github.com/sniperHW/flyfish/pkg/net"
	"initialthree/cluster/addr"
	"initialthree/common"
	"initialthree/network"
	"initialthree/network/smux"
	"initialthree/pkg/timer"
	"time"
)

func (this *Cluster) dialError(end *endPoint, session *fnet.Socket, err error, counter int) {

	if nil != session {
		session.Close(err, 0)
	}

	end.Lock()
	defer end.Unlock()

	end.dialing = false

	/*
	 * 如果end.session != nil 表示两端同时请求建立连接，本端已经作为服务端成功接受了对端的连接
	 */
	if nil == end.session {
		if (err != ERR_INVAILD_ENDPOINT || err != ERR_AUTH) &&
			(end.pendingCall.Len() > 0 || counter < dialTerminateCount) {
			end.dialing = true
			this.RegisterTimerOnce(time.Second, func(t *timer.Timer, _ interface{}) {
				this._dial(end, counter+1)
			}, nil)
		} else {
			end.pendingMsg = []interface{}{}
			for v := end.pendingCall.Front(); v != nil; v = end.pendingCall.Front() {
				c := end.pendingCall.Remove(v).(*rpcCall)
				c.listElement = nil
				if c.dialTimer.Cancel() {
					c.cb(nil, err)
				}
			}
		}
	}
}

func (this *Cluster) dialOK(end *endPoint, session *fnet.Socket) {
	if end == this.serviceMgr.getEndPoint(end.addr.Logic) {
		end.Lock()
		defer end.Unlock()
		this.onEstablishClient(end, session)
	} else {
		//不再是合法的end
		logger.Sugar().Errorf("%s dial error: ERR_INVAILD_ENDPOINT ", end.addr.Logic.String())
		this.dialError(end, session, ERR_INVAILD_ENDPOINT, dialTerminateCount)
	}
}

func (this *Cluster) _dial(end *endPoint, counter int) {

	logger.Sugar().Infof("(self:%v) dial %s %v\n", this.serverState.selfAddr.Logic, end.addr.Logic.String(), end.addr.AtomicGetNetaddr())
	go func() {
		conn, err := network.Dial("tcp", end.addr.AtomicGetNetaddr().String(), time.Second*3)
		if err != nil {
			logger.Sugar().Errorf("(self:%v) %s dial error: %v ", this.serverState.selfAddr.Logic, end.addr.Logic.String(), err)
			this.dialError(end, nil, ERR_DIAL, counter)
		} else {
			session := fnet.NewSocket(conn, fnet.OutputBufLimit{
				OutPutLimitSoft:        512 * 1024,
				OutPutLimitSoftSeconds: 10,
				OutPutLimitHard:        8 * 1024 * 1024,
			})
			if err := this.login(end, conn, false); err != nil {
				logger.Sugar().Errorf("(self:%v) %s dial error: %v ", this.serverState.selfAddr.Logic, end.addr.Logic.String(), err)
				this.dialError(end, session, err, counter)
			} else {
				this.dialOK(end, session)
			}
		}
	}()
}

func (this *Cluster) dial(end *endPoint, counter int) {
	//发起异步Dial连接
	if !end.dialing {
		end.dialing = true
		this._dial(end, counter)
	}
}

func (this *Cluster) DialMuxSocket(peer addr.LogicAddr, enc func(o interface{}, b *buffer.Buffer) error, onSocketClose func(*smux.MuxSocket)) (*smux.MuxSocket, error) {
	end := this.serviceMgr.getEndPoint(peer)
	if end == nil {
		return nil, fmt.Errorf("%s not found", peer.String())
	}

	conn, err := network.Dial("tcp", end.addr.AtomicGetNetaddr().String(), time.Second*3)
	if err != nil {
		return nil, err
	}

	if err := this.login(end, conn, true); err != nil {
		return nil, err
	}
	return smux.NewMuxSocketClient(conn, common.HeartBeat_Timeout, enc, onSocketClose), nil
}
