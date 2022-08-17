package cluster

import (
	"errors"
	"initialthree/cluster/addr"
	cluster_proto "initialthree/cluster/proto"
	"initialthree/cluster/rpcerr"
	"initialthree/codec/ss"
	"initialthree/pkg/rpc"
	"math/rand"
	"sort"
	"strings"
	"time"
)

func (this *serviceManager) isSelfHarbor() bool {
	return this.cluster.serverState.selfAddr.Logic.Type() == harbarType
}

func (this *serviceManager) getHarbor(m addr.LogicAddr) *endPoint {
	return this.getHarborByGroup(this.cluster.serverState.selfAddr.Logic.Group(), m)
}

func (this *serviceManager) getHarborByGroup(group uint32, m addr.LogicAddr) *endPoint {
	this.RLock()
	defer this.RUnlock()
	harborGroup, ok := this.harborsByGroup[group]
	if ok && len(harborGroup) > 0 {
		return harborGroup[int(m)%len(harborGroup)]
	} else {
		return nil
	}
}

func (this *serviceManager) addHarbor(harbor *endPoint) {

	group := harbor.addr.Logic.Group()

	harborGroup, ok := this.harborsByGroup[group]
	if !ok {
		this.harborsByGroup[group] = []*endPoint{harbor}
	} else {
		for _, v := range this.harborsByGroup[group] {
			if v.addr.Logic == harbor.addr.Logic {
				return
			}
		}
		this.harborsByGroup[group] = append(harborGroup, harbor)
	}

	logger.Sugar().Infof("addHarbor %s\n", harbor.addr.Logic.String())

}

func (this *serviceManager) removeHarbor(harbor *endPoint) {
	if nil != harbor.session {
		harbor.session.Close(errors.New("remove harbor"), 0)
	}
	group := harbor.addr.Logic.Group()
	harborGroup, ok := this.harborsByGroup[group]
	if ok && len(harborGroup) > 0 {
		i := 0
		for ; i < len(harborGroup); i++ {
			if harborGroup[i].addr.Logic == harbor.addr.Logic {
				break
			}
		}
		if i != len(harborGroup) {
			harborGroup[i], harborGroup[len(harborGroup)-1] = harborGroup[len(harborGroup)-1], harborGroup[i]
			this.harborsByGroup[group] = harborGroup[:len(harborGroup)-1]
		}
	}
}

func (this *serviceManager) onEndPointJoin(end *endPoint) {
	if this.isSelfHarbor() {
		if end.addr.Logic.Type() == harbarType {
			if end.addr.Logic.Group() != this.cluster.serverState.selfAddr.Logic.Group() {
				this.NotifyForginServicesH2H(end)
			}
		} else {
			this.NotifyForginServicesH2S(end)

			if end.exportService == 1 {
				logger.Sugar().Infof("exportService %s join \n", end.addr.Logic.String())
				//通告非本group的其它harbor
				for g, v := range this.harborsByGroup {
					if g != this.cluster.serverState.selfAddr.Logic.Group() {
						for _, vv := range v {
							this.NotifyForginServicesH2H(vv)
						}
					}
				}
			}
		}
	}
}

func (this *serviceManager) onEndPointLeave(end *endPoint) {
	if this.isSelfHarbor() && end.addr.Logic.Type() != harbarType {
		//通告非本group的其它harbor,有节点离开,需要移除forginService
		for g, v := range this.harborsByGroup {
			if g != this.cluster.serverState.selfAddr.Logic.Group() {
				for _, vv := range v {
					this.NotifyForginServicesH2H(vv)
				}
			}
		}
	}
}

func nodes2Str(nodes []uint32) string {
	s := []string{}
	for _, v := range nodes {
		t := addr.LogicAddr(v)
		s = append(s, t.String())
	}

	return strings.Join(s, ",")
}

func (this *serviceManager) NotifyForginServicesH2H(end *endPoint) {

	forgins := []uint32{}
	for k, v := range this.idEndPointMap {
		if v.exportService == 1 {
			forgins = append(forgins, uint32(k))
		}
	}

	req := &cluster_proto.NotifyForginServicesH2HReq{Nodes: forgins}
	this.cluster.asynCall(end.addr.Logic, end, req, time.Second, func(_ interface{}, err error) {
		if nil == err {
			logger.Sugar().Infof("%s NotifyForginServicesH2H to %s ok %s\n", this.cluster.serverState.selfAddr.Logic.String(), end.addr.Logic.String(), nodes2Str(forgins))
		} else if nil != err {
			go func() {
				time.Sleep(time.Second)
				if end == this.getEndPoint(end.addr.Logic) {
					logger.Sugar().Infof("%s NotifyForginServicesH2H to %s error:%v, tryagain %s\n", this.cluster.serverState.selfAddr.Logic.String(), end.addr.Logic, err, nodes2Str(forgins))
					this.RLock()
					this.NotifyForginServicesH2H(end)
					this.RUnlock()
				}
			}()
		}
	})
}

func (this *serviceManager) NotifyForginServicesH2S(end *endPoint) {
	forgins := []uint32{}
	for tt, v := range this.ttForignServiceMap {
		if tt != harbarType {
			for _, vv := range v.services {
				forgins = append(forgins, uint32(vv))
			}
		}
	}
	req := &cluster_proto.NotifyForginServicesH2SReq{Nodes: forgins}
	this.cluster.asynCall(end.addr.Logic, end, req, time.Second, func(_ interface{}, err error) {
		if nil == err {
			logger.Sugar().Infof("%s NotifyForginServicesH2S to %s ok %v\n", this.cluster.serverState.selfAddr.Logic.String(), end.addr.Logic.String(), nodes2Str(forgins))
		} else if nil != err {
			go func() {
				time.Sleep(time.Second)
				if end == this.getEndPoint(end.addr.Logic) {
					logger.Sugar().Infof("%s NotifyForginServicesH2S to %s error:%v,tryagain %v\n", this.cluster.serverState.selfAddr.Logic.String(), end.addr.Logic.String(), err, nodes2Str(forgins))
					this.RLock()
					this.NotifyForginServicesH2S(end)
					this.RUnlock()
				}
			}()
		}
	})
}

func (this *serviceManager) brocastH2S() {
	this.RLock()
	defer this.RUnlock()
	for tt, v := range this.ttEndPointMap {
		if tt != harbarType {
			for _, vv := range v.endPoints {
				this.NotifyForginServicesH2S(vv)
			}
		}
	}
}

func diff2(a, b []uint32) ([]uint32, []uint32) {

	if len(a) == 0 {
		return nil, b
	}

	if len(b) == 0 {
		return a, nil
	}

	sort.Slice(a, func(i, j int) bool {
		return a[i] < a[j]
	})

	sort.Slice(b, func(i, j int) bool {
		return b[i] < b[j]
	})

	add := []uint32{}
	remove := []uint32{}

	i := 0
	j := 0

	for i < len(a) && j < len(b) {
		if a[i] == b[j] {
			add = append(add, a[i])
			i++
			j++
		} else if a[i] > b[j] {
			remove = append(remove, b[j])
			j++
		} else {
			add = append(add, a[i])
			i++
		}
	}

	if len(a[i:]) > 0 {
		add = append(add, a[i:]...)
	}

	if len(b[j:]) > 0 {
		remove = append(remove, b[j:]...)
	}

	return add, remove
}

func (this *serviceManager) initHarbor() {

	this.cluster.RegisterMethod(&cluster_proto.NotifyForginServicesH2HReq{}, func(replyer *rpc.RPCReplyer, req interface{}) {
		if testRPCTimeout && rand.Int()%2 == 0 {
			replyer.DropResponse()
		} else {

			if this.isSelfHarbor() {

				msg := req.(*cluster_proto.NotifyForginServicesH2HReq)

				current := this.getAllForginService()

				logger.Sugar().Infof("NotifyForginServicesH2HReq %v", msg.GetNodes())

				add, remove := diff2(msg.GetNodes(), current)

				for _, v := range add {
					this.addForginService(addr.LogicAddr(v))
				}

				for _, v := range remove {
					this.removeForginService(addr.LogicAddr(v))
				}

				if len(add) > 0 || len(remove) > 0 {
					this.brocastH2S()
				}

			}

			replyer.Reply(&cluster_proto.NotifyForginServicesH2HResp{}, nil)
		}
	})
}

func (this *Cluster) postRelayError(peer addr.LogicAddr, msg *ss.RPCRelayErrorMessage) {

	endPoint := this.serviceMgr.getEndPoint(peer)
	if nil == endPoint {
		endPoint = this.serviceMgr.getHarborByGroup(peer.Group(), peer)
		if nil != endPoint && endPoint.addr.Logic == this.serverState.selfAddr.Logic {
			logger.Sugar().Errorf("postRelayError ring!!!\n")
			return
		}
	}

	if nil != endPoint {
		endPoint.Lock()
		defer endPoint.Unlock()
		if nil != endPoint.session {
			endPoint.session.Send(msg)
		} else {
			endPoint.pendingMsg = append(endPoint.pendingMsg, msg)
			//尝试与对端建立连接
			this.dial(endPoint, 0)
		}
	} else {
		logger.Sugar().Errorf("postRelayError %s not found", peer.String())
	}
}

func (this *Cluster) onRelayError(message *ss.RelayMessage, err error) {
	if message.IsRPCReq() {
		//通告请求端消息无法送达到目的地
		msg := &ss.RPCRelayErrorMessage{
			To:    message.From,
			From:  this.serverState.selfAddr.Logic,
			Seqno: message.GetSeqno(),
			Err:   err,
		}

		logger.Sugar().Errorf("onRelayError %v\n", err)

		this.postRelayError(message.From, msg)
	}
}

func (this *Cluster) onRelayMessage(message *ss.RelayMessage) {

	logger.Sugar().Debugf("%s onRelayMessage target:%s from:%s\n", this.serverState.selfAddr.Logic.String(), message.To.String(), message.From.String())

	endPoint := this.serviceMgr.getEndPoint(message.To)
	if nil == endPoint {
		if message.To.Group() != this.serverState.selfAddr.Logic.Group() {
			//不同group要求harbor转发
			endPoint = this.serviceMgr.getHarborByGroup(message.To.Group(), message.To)
			if nil != endPoint && endPoint.addr.Logic == this.serverState.selfAddr.Logic {
				logger.Sugar().Errorf("onRelayMessage ring!!!\n")
				return
			}
		} else {
			//同group,server为0,则从本地随机选择一个符合type的server
			if message.To.Server() == 0 {
				if addr, err := this.Random(message.To.Type()); nil == err {
					endPoint = this.serviceMgr.getEndPoint(addr)
					//需要将to设置为正确的地址，否则无法转发
					message.ResetTo(addr)
				}
			}
		}
	}

	if nil != endPoint {
		endPoint.Lock()
		defer endPoint.Unlock()

		if nil != endPoint.session {
			err := endPoint.session.Send(message)
			if nil != err {
				logger.Sugar().Debug(err)
			}

		} else {
			endPoint.pendingMsg = append(endPoint.pendingMsg, message)
			//尝试与对端建立连接
			this.dial(endPoint, 0)
		}
	} else {
		logger.Sugar().Infof("unable route to target %s\n", message.To.String())
		this.onRelayError(message, rpcerr.Err_RPC_RelayError)
	}
}
