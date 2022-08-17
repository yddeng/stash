package cluster

import (
	"container/list"
	"fmt"
	"github.com/golang/protobuf/proto"
	"initialthree/cluster/addr"
	"initialthree/cluster/priority"
	cluster_proto "initialthree/cluster/proto"
	"initialthree/codec/pb"
	"initialthree/codec/ss"
	"initialthree/network"
	"initialthree/pkg/event"
	"initialthree/pkg/rpc"
	"initialthree/pkg/timer"
	_ "initialthree/protocol/ss" //触发pb注册
	"net"
	"reflect"
	"sync"
	"sync/atomic"
	"time"
)

type clusterState struct {
	selfAddr addr.Addr
	started  int32
	stoped   int32
}

type rpcCall struct {
	arg         interface{}
	dialTimer   *timer.Timer
	deadline    time.Time
	to          addr.LogicAddr
	cb          rpc.RPCResponseHandler
	listElement *list.Element
}

func waitCondition(eventq *event.EventQueue, fn func() bool) {
	wg := sync.WaitGroup{}
	wg.Add(1)

	donefire := int32(0)

	if nil == eventq {
		go func() {
			for {
				time.Sleep(time.Millisecond * 100)
				if fn() {
					if atomic.LoadInt32(&donefire) == 0 {
						atomic.StoreInt32(&donefire, 1)
						wg.Done()
					}
					break
				}
			}
		}()
	} else {
		go func() {
			stoped := int32(0)
			ch := make(chan struct{})
			for atomic.LoadInt32(&stoped) == 0 {
				time.Sleep(time.Millisecond * 100)
				eventq.PostNoWait(priority.HIGH, func() {
					if fn() {
						if atomic.LoadInt32(&donefire) == 0 {
							atomic.StoreInt32(&donefire, 1)
							wg.Done()
						}
						atomic.StoreInt32(&stoped, 1)
					}
					ch <- struct{}{}
				})
				_ = <-ch
			}
		}()
	}

	wg.Wait()
}

func (this *Cluster) waitCondition(fn func() bool) {
	waitCondition(this.GetEventQueue(), fn)
}

/*
 *   如果服务要下线才将用Stop(true)调用
 *   如果只是重启更新不需要带参数
 */
func (this *Cluster) Stop(stopFunc func(), sendRemoveNode ...bool) {
	if atomic.CompareAndSwapInt32(&this.serverState.stoped, 0, 1) {
		if nil != stopFunc {
			stopFunc()
		}

		//等待接收以及发送的rpc都处理完成
		waitCondition(nil, func() bool {
			return this.rpcMgr.server.PendingCount() == 0 && atomic.LoadInt32(&this.pendingRPCRequestCount) == 0
		})

		this.queue.Close()

		this.l.Close()
		if len(sendRemoveNode) > 0 && sendRemoveNode[0] {
			this.centerClient.Close(true)
			if nil != this.uniLocker {
				this.uniLocker.Unlock()
			}
		} else {
			this.centerClient.Close(false)
		}
	}
}

func (this *Cluster) IsStoped() bool {
	return atomic.LoadInt32(&this.serverState.stoped) == 1
}

func (this *Cluster) postToEndPoint(end *endPoint, msg interface{}) {
	end.Lock()
	defer end.Unlock()

	if nil != end.session {
		err := end.session.Send(msg)
		if nil != err {
			//记录日志
			logger.Sugar().Errorf("Send error: %s %s\n", err.Error(), reflect.TypeOf(msg).String())
		} else {
			end.lastActive = time.Now()
		}
	} else {
		end.pendingMsg = append(end.pendingMsg, msg)
		//尝试与对端建立连接
		this.dial(end, 0)
	}
}

//向类型为tt的本cluster节点广播
func (this *Cluster) Brocast(tt uint32, msg proto.Message, exceptSelf ...bool) {
	this.serviceMgr.RLock()
	defer this.serviceMgr.RUnlock()
	if ttmap, ok := this.serviceMgr.ttEndPointMap[tt]; ok {
		for _, v := range ttmap.endPoints {
			if len(exceptSelf) == 0 || v.addr.Logic != this.serverState.selfAddr.Logic {
				this.postToEndPoint(v, msg)
			}
		}
	}
}

//向本cluster内所有节点广播
func (this *Cluster) BrocastToAll(msg proto.Message, exceptTT ...uint32) {
	this.serviceMgr.RLock()
	defer this.serviceMgr.RUnlock()
	exceptType := uint32(0)
	if len(exceptTT) > 0 {
		exceptType = exceptTT[0]
	}

	for tt, v := range this.serviceMgr.ttEndPointMap {
		if tt != exceptType {
			for _, vv := range v.endPoints {
				this.postToEndPoint(vv, msg)
			}
		}
	}
}

/*
 *  异步投递
 */
func (this *Cluster) Post(peer addr.LogicAddr, msg interface{}) {
	if atomic.LoadInt32(&this.serverState.started) == 0 {
		logger.Sugar().Errorf("PostMessage cluster not started\n")
		return
	}

	if peer.Empty() {
		logger.Sugar().Errorf("PostMessage to empty peer\n")
		return
	}

	endPoint := this.serviceMgr.getEndPoint(peer)

	if nil == endPoint {
		if peer.Group() == this.serverState.selfAddr.Logic.Group() {
			//记录日志
			logger.Sugar().Errorf("PostMessage %s not found\n", peer.String())
			return
		} else {
			//不同服务器组，需要通告Harbor转发
			harbor := this.serviceMgr.getHarbor(peer)
			if nil == harbor {
				logger.Sugar().Errorf("Post cross Group message failed,no harbor:%s\n", peer.String())
				return
			} else {
				endPoint = harbor
			}
		}
	}

	this.postToEndPoint(endPoint, msg)
}

func (this *Cluster) PostMessage(peer addr.LogicAddr, msg proto.Message) {
	this.Post(peer, ss.NewMessage(msg, peer, this.serverState.selfAddr.Logic))
}

func (this *Cluster) SelfAddr() addr.Addr {
	return this.serverState.selfAddr
}

/*
*  启动服务
 */
func (this *Cluster) Start(center_addr []string, selfAddr addr.Addr, uniLocker UniLocker, export ...bool) error {
	if selfAddr.Logic.Server() == uint32(0) {
		return ERR_SERVERADDR_ZERO
	}

	if !atomic.CompareAndSwapInt32(&this.serverState.started, 0, 1) {
		return ERR_STARTED
	}

	if !uniLocker.Lock(selfAddr) {
		return fmt.Errorf("lock logic addr %s failed\n", selfAddr.Logic.String())
	}

	this.uniLocker = uniLocker

	this.serverState.selfAddr = selfAddr

	l, serve, err := network.Listen("tcp", this.serverState.selfAddr.Net.String(), func(conn net.Conn) {
		go func() {
			if err := this.auth(conn); nil != err {
				logger.Sugar().Infof("auth error %s self %s", err.Error(), this.serverState.selfAddr.Logic.String())
				conn.Close()
			}
		}()
	})

	if nil != err {
		return err
	} else {
		this.l = l
		go func() {
			this.queue.Run()
		}()

		this.serviceMgr.init()
		this.centerInit(export...)
		this.connectCenter(center_addr)

		go serve()

		return nil
	}
}

func (this *Cluster) GetEventQueue() *event.EventQueue {
	return this.queue
}

/*
*  将一个闭包投递到队列中执行，args为传递给闭包的参数
 */
func (this *Cluster) PostTask(function interface{}, args ...interface{}) {
	this.queue.PostNoWait(priority.LOW, function, args...)
}

func (this *Cluster) Mod(tt uint32, num int) (addr.LogicAddr, error) {
	this.serviceMgr.RLock()
	defer this.serviceMgr.RUnlock()

	//优先从本集群查找
	if ttmap, ok := this.serviceMgr.ttEndPointMap[tt]; ok {
		addr_, err := ttmap.mod(num)
		if nil == err {
			return addr_, err
		}
	}

	//从forginService查找
	if smap, ok := this.serviceMgr.ttForignServiceMap[tt]; ok {
		return smap.mod(num)
	}

	return addr.LogicAddr(0), ERR_NO_AVAILABLE_SERVICE
}

//随机获取一个类型为tt的节点id
func (this *Cluster) Random(tt uint32) (addr.LogicAddr, error) {
	this.serviceMgr.RLock()
	defer this.serviceMgr.RUnlock()

	//优先从本集群查找
	if ttmap, ok := this.serviceMgr.ttEndPointMap[tt]; ok {
		addr_, err := ttmap.random()
		if nil == err {
			return addr_, err
		}
	}

	//从forginService查找
	if smap, ok := this.serviceMgr.ttForignServiceMap[tt]; ok {
		return smap.random()
	}

	return addr.LogicAddr(0), ERR_NO_AVAILABLE_SERVICE
}

func (this *Cluster) Select(tt uint32) ([]addr.LogicAddr, error) {
	this.serviceMgr.RLock()
	defer this.serviceMgr.RUnlock()

	ttmap := this.serviceMgr.ttEndPointMap[tt]
	if ttmap == nil {
		return nil, ERR_NO_AVAILABLE_SERVICE
	}

	if len(ttmap.endPoints) == 0 {
		return nil, ERR_NO_AVAILABLE_SERVICE
	}

	addrs := make([]addr.LogicAddr, 0, len(ttmap.endPoints))
	for _, ep := range ttmap.endPoints {
		addrs = append(addrs, ep.addr.Logic)
	}

	return addrs, nil
}

func init() {
	pb.Register("ss", &cluster_proto.Heartbeat{}, 1)
}
