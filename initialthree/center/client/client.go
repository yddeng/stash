package center

import (
	"errors"
	"github.com/golang/protobuf/proto"
	fnet "github.com/sniperHW/flyfish/pkg/net"
	"go.uber.org/zap"
	"initialthree/center/constant"
	center_proto "initialthree/center/protocol"
	center_rpc "initialthree/center/rpc"
	"initialthree/cluster/addr"
	"initialthree/cluster/priority"
	"initialthree/codec/ss"
	"initialthree/network"
	"initialthree/pkg/event"
	"initialthree/pkg/rpc"
	"os"
	"reflect"
	"sync"
	"sync/atomic"
	"time"
)

type centerHandler func(*fnet.Socket, proto.Message)

type CenterClient struct {
	logger *zap.Logger //kendynet.LoggerI

	clientProcessQueue *event.EventQueue

	exportService uint32 //本节点是否暴露到服务器组外面

	centerHandlers map[uint16]centerHandler

	rpcClient *rpc.RPCClient

	closed int32

	centers []*center
}

func (this *CenterClient) RegisterCenterMsgHandler(cmd uint16, handler centerHandler) {
	if nil == handler {
		//记录日志
		this.logger.Sugar().Errorf("Register %d failed: handler is nil\n", cmd)
		return
	}

	_, ok := this.centerHandlers[cmd]
	if ok {
		//记录日志
		this.logger.Sugar().Errorf("Register %d failed: duplicate handler\n", cmd)
		return
	}

	this.centerHandlers[cmd] = handler
}

func (this *CenterClient) Close(sendRemoveNode bool) {
	if atomic.CompareAndSwapInt32(&this.closed, 0, 1) {
		for _, v := range this.centers {
			v.close(sendRemoveNode)
		}
	}
}

func (this *CenterClient) dispatchCenterMsg(session *fnet.Socket, msg *ss.Message) {

	if atomic.LoadInt32(&this.closed) == 0 {
		data := msg.GetData()
		switch data.(type) {
		case *rpc.RPCResponse:
			center_rpc.OnRPCResponse(this.rpcClient, data.(*rpc.RPCResponse))
		case proto.Message:
			cmd := msg.GetCmd()
			handler, ok := this.centerHandlers[cmd]
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

func (this *CenterClient) ConnectCenter(centerAddrs []string, selfAddr addr.Addr) {
	centers := map[string]bool{}
	for _, v := range centerAddrs {
		if _, ok := centers[v]; !ok {
			centers[v] = true
			c := &center{
				addr:         v,
				selfAddr:     selfAddr,
				centerClient: this,
			}
			this.centers = append(this.centers, c)
			c.connect()
		}
	}
}

func New(queue *event.EventQueue, l *zap.Logger, export uint32) *CenterClient {
	return &CenterClient{
		clientProcessQueue: queue,
		logger:             l,
		exportService:      export,
		rpcClient:          center_rpc.NewClient(),
		centerHandlers:     map[uint16]centerHandler{},
		centers:            []*center{},
	}
}

type center struct {
	sync.Mutex
	addr         string
	selfAddr     addr.Addr
	centerClient *CenterClient
	session      *fnet.Socket
	closed       int32
}

func login(centerClient *CenterClient, session *fnet.Socket, req *center_proto.Login, onResp func(interface{}, error)) error {
	return center_rpc.AsynCall(centerClient.rpcClient, session, req, onResp)
}

func (this *center) close(sendRemoveNode bool) {
	if atomic.CompareAndSwapInt32(&this.closed, 0, 1) {
		this.Lock()
		s := this.session
		this.Unlock()
		if nil != s {
			if sendRemoveNode {
				s.Send(&center_proto.RemoveNode{
					Nodes: []uint32{uint32(this.selfAddr.Logic)},
				})
			}
			s.Close(errors.New("centerClient close"), time.Second)
		}
	}
}

func (this *center) connect() {
	go func() {
		for atomic.LoadInt32(&this.closed) == 0 {
			conn, err := network.Dial("tcp", this.addr, time.Second*3)
			if err != nil {
				time.Sleep(time.Millisecond * 1000)
			} else {
				session := fnet.NewSocket(conn, fnet.OutputBufLimit{})
				this.Lock()
				this.session = session
				this.Unlock()
				session.SetRecvTimeout(time.Second * time.Duration(constant.HeartBeatTimeout))
				session.SetInBoundProcessor(center_proto.NewReceiver())
				session.SetEncoder(center_proto.NewEncoder())

				done := make(chan struct{}, 1)

				session.SetCloseCallBack(func(sess *fnet.Socket, reason error) {
					this.Lock()
					this.session = nil
					this.Unlock()
					this.centerClient.logger.Sugar().Infof("center disconnected %s self:%s\n", reason.Error(), this.selfAddr.Logic.String())
					done <- struct{}{}
					if atomic.LoadInt32(&this.closed) == 0 {
						this.connect()
					}
				}).BeginRecv(func(s *fnet.Socket, m interface{}) {
					msg := m.(*ss.Message)
					this.centerClient.clientProcessQueue.PostNoWait(priority.HIGH, this.centerClient.dispatchCenterMsg, session, msg)
				})

				loginReq := &center_proto.Login{
					LogicAddr:     proto.Uint32(uint32(this.selfAddr.Logic)),
					NetAddr:       proto.String(this.selfAddr.Net.String()),
					ExportService: proto.Uint32(this.centerClient.exportService),
				}

				var onResp func(interface{}, error)

				onResp = func(r interface{}, err error) {
					if nil != err {
						if err == rpc.ErrCallTimeout {
							this.centerClient.logger.Sugar().Errorf("login timeout %v self:%v\n", this.addr, this.selfAddr)
							login(this.centerClient, session, loginReq, onResp)
						} else {
							//登录center中如果出现除超时以外的错误，直接退出进程
							this.centerClient.logger.Sugar().Errorf("%v %v\n", this.addr, err)
							os.Exit(0)
						}
					} else {
						resp := r.(*center_proto.LoginRet)
						if resp.GetErrCode() == constant.LoginOK {
							this.centerClient.logger.Sugar().Infof("login ok %v self:%v\n", this.addr, this.selfAddr)
							ticker := time.NewTicker(time.Second * time.Duration(constant.HeartBeatTimeout/2))
							go func() {
								for {
									select {
									case <-done:
										ticker.Stop()
										return
									case <-ticker.C:
										//发送心跳
										session.Send(&center_proto.HeartbeatToCenter{
											Timestamp: proto.Int64(time.Now().UnixNano()),
										})
									}
								}
							}()
						} else {
							panic(resp.GetMsg())
						}
					}
				}
				if err := login(this.centerClient, session, loginReq, onResp); nil != err {
					this.centerClient.logger.Sugar().Errorf("%v %v\n", this.addr, err)
					os.Exit(0)
				}
				break
			}
		}
	}()
}
