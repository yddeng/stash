/*
 *   非本集群内直连的外部服务，请求需要通过harbor转发
 */

package cluster

import (
	"initialthree/cluster/addr"
	cluster_proto "initialthree/cluster/proto"
	"initialthree/pkg/rpc"
	"math/rand"
	"sort"
)

type typeForignServiceMap struct {
	tt       uint32
	services []addr.LogicAddr
}

func (this *typeForignServiceMap) sort() {
	sort.Slice(this.services, func(i, j int) bool {
		return uint32(this.services[i]) < uint32(this.services[j])
	})
}

func (this *typeForignServiceMap) remove(addr_ addr.LogicAddr) {
	for i, v := range this.services {
		if addr_ == v {
			this.services[i] = this.services[len(this.services)-1]
			this.services = this.services[:len(this.services)-1]
			this.sort()
			break
		}
	}
}

func (this *typeForignServiceMap) add(addr_ addr.LogicAddr) {

	find := false
	for _, v := range this.services {
		if addr_ == v {
			find = true
			break
		}
	}

	if !find {
		this.services = append(this.services, addr_)
		this.sort()
	}
}

func (this *typeForignServiceMap) mod(num int) (addr.LogicAddr, error) {
	size := len(this.services)
	if size > 0 {
		i := num % size
		return this.services[i], nil
	} else {
		return addr.LogicAddr(0), ERR_NO_AVAILABLE_SERVICE
	}
}

func (this *typeForignServiceMap) random() (addr.LogicAddr, error) {
	size := len(this.services)
	if size > 0 {
		i := rand.Int() % size
		return this.services[i], nil
	} else {
		return addr.LogicAddr(0), ERR_NO_AVAILABLE_SERVICE
	}
}

func (this *serviceManager) addForginService(addr_ addr.LogicAddr) {
	this.Lock()
	defer this.Unlock()

	m, ok := this.ttForignServiceMap[addr_.Type()]
	if !ok {
		m = &typeForignServiceMap{
			tt:       addr_.Type(),
			services: []addr.LogicAddr{},
		}
		this.ttForignServiceMap[addr_.Type()] = m
	}

	m.add(addr_)

	if this.isSelfHarbor() {
		logger.Sugar().Infof("harbor %s addForginService %s\n", this.cluster.serverState.selfAddr.Logic.String(), addr_.String())
	} else {
		logger.Sugar().Infof("%s addForginService %s\n", this.cluster.serverState.selfAddr.Logic.String(), addr_.String())
	}
}

func (this *serviceManager) removeForginService(addr_ addr.LogicAddr) {
	this.Lock()
	defer this.Unlock()

	m, ok := this.ttForignServiceMap[addr_.Type()]
	if ok {
		m.remove(addr_)
	}
}

func (this *serviceManager) getAllForginService() []uint32 {
	this.RLock()
	defer this.RUnlock()
	current := []uint32{}
	for _, v1 := range this.ttForignServiceMap {
		for _, v2 := range v1.services {
			current = append(current, uint32(v2))
		}
	}
	return current
}

func (this *serviceManager) init() {

	this.cluster.RegisterMethod(&cluster_proto.NotifyForginServicesH2SReq{}, func(replyer *rpc.RPCReplyer, req interface{}) {
		if testRPCTimeout && rand.Int()%2 == 0 {
			replyer.DropResponse()
		} else {
			if !this.isSelfHarbor() {

				msg := req.(*cluster_proto.NotifyForginServicesH2SReq)

				logger.Sugar().Infof("NotifyForginServicesH2SReq %v", msg.GetNodes())

				current := this.getAllForginService()

				add, remove := diff2(msg.GetNodes(), current)

				for _, v := range add {
					this.addForginService(addr.LogicAddr(v))
				}

				for _, v := range remove {
					this.removeForginService(addr.LogicAddr(v))
				}

			}
			replyer.Reply(&cluster_proto.NotifyForginServicesH2SResp{}, nil)
		}
	})
	this.initHarbor()
}
