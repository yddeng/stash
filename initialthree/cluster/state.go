package cluster

import (
	"go.uber.org/zap"
	center_client "initialthree/center/client"
	"initialthree/cluster/addr"
	"initialthree/cluster/priority"
	"initialthree/cluster/rpcerr"
	"initialthree/pkg/event"
	"initialthree/pkg/rpc"
	"net"
	"sync"
)

const EXPORT bool = true
const dialTerminateCount int = 5
const harbarType uint32 = 255

var logger *zap.Logger //kendynet.LoggerI

type serviceManager struct {
	sync.RWMutex
	cluster            *Cluster
	idEndPointMap      map[addr.LogicAddr]*endPoint
	ttEndPointMap      map[uint32]*typeEndPointMap
	ttForignServiceMap map[uint32]*typeForignServiceMap
	harborsByGroup     map[uint32][]*endPoint
}

//确保同一逻辑地址只能被唯一进程使用
type UniLocker interface {
	Lock(addr.Addr) bool
	Unlock()
}

type Cluster struct {
	serverState            clusterState
	queue                  *event.EventQueue
	rpcMgr                 rpcManager
	serviceMgr             serviceManager
	msgMgr                 msgManager
	centerClient           *center_client.CenterClient
	l                      net.Listener //*listener.Listener
	uniLocker              UniLocker
	pendingRPCRequestCount int32

	onNewMuxConn func(logic addr.LogicAddr, conn net.Conn)
}

func (this *Cluster) SetNewMuxConn(onNewMuxConn func(logic addr.LogicAddr, conn net.Conn)) {

	if this.onNewMuxConn != nil {
		panic("cluster.state:SetNewMuxSocket onNewMuxSocket is exist.")
	}
	this.onNewMuxConn = onNewMuxConn
}

func NewCluster() *Cluster {
	c := &Cluster{
		queue: event.NewEventQueueWithPriority(priority.HIGH+1, 10000),
		rpcMgr: rpcManager{
			server: rpc.NewRPCServer(&decoder{}, &encoder{}),
			client: rpc.NewClient(&decoder{}, &encoder{}),
		},
		msgMgr: msgManager{
			msgHandlers: map[uint16]MsgHandler{},
		},
		serviceMgr: serviceManager{
			idEndPointMap:      map[addr.LogicAddr]*endPoint{},
			ttEndPointMap:      map[uint32]*typeEndPointMap{},
			ttForignServiceMap: map[uint32]*typeForignServiceMap{},
			harborsByGroup:     map[uint32][]*endPoint{},
		},
	}
	c.serviceMgr.cluster = c
	c.rpcMgr.cluster = c
	c.rpcMgr.server.SetErrorCodeOnMissingMethod(rpcerr.Err_RPC_InvaildMethod)
	return c
}

var defaultCluster *Cluster = NewCluster()
