package center

import (
	"github.com/golang/protobuf/proto"
	fnet "github.com/sniperHW/flyfish/pkg/net"
	"initialthree/pkg/rpc"
	"reflect"
	"time"
)

var timeout = time.Second * 8 // 超时

type RPCChannel struct {
	session *fnet.Socket
}

func (this *RPCChannel) SendRequest(message interface{}) error {
	return this.session.Send(message)
}

func (this *RPCChannel) SendResponse(message interface{}) error {
	return this.session.Send(message)
}

func (this *RPCChannel) Name() string {
	return this.session.RemoteAddr().String() + "<->" + this.session.LocalAddr().String()
}

func (this *RPCChannel) UID() uint64 {
	return uint64(uintptr(reflect.ValueOf(this.session.GetUnderConn()).Pointer()))
}

func (this *RPCChannel) GetSession() *fnet.Socket {
	return this.session
}

type encoder struct {
}

func (this *encoder) Encode(message rpc.RPCMessage) (interface{}, error) {
	return message, nil
}

type decoder struct {
}

func (this *decoder) Decode(o interface{}) (rpc.RPCMessage, error) {
	return o.(rpc.RPCMessage), nil
}

/*
 *  注册RPC服务,无锁保护，务必在初始化时完成
 */
func RegisterMethod(rpcServer *rpc.RPCServer, arg proto.Message, handler rpc.RPCMethodHandler) {
	rpcServer.RegisterMethod(reflect.TypeOf(arg).String(), handler)
}

func AsynCall(rpcClient *rpc.RPCClient, ses *fnet.Socket, arg proto.Message, cb rpc.RPCResponseHandler) error {
	return rpcClient.AsynCall(&RPCChannel{session: ses}, reflect.TypeOf(arg).String(), arg, timeout, cb)
}

func OnRPCRequest(rpcServer *rpc.RPCServer, ses *fnet.Socket, msg interface{}) {
	rpcServer.OnRPCMessage(&RPCChannel{session: ses}, msg)
}

func OnRPCResponse(rpcClient *rpc.RPCClient, msg interface{}) {
	rpcClient.OnRPCMessage(msg)
}

func NewClient() *rpc.RPCClient {
	return rpc.NewClient(&decoder{}, &encoder{})
}

func NewServer() *rpc.RPCServer {
	return rpc.NewRPCServer(&decoder{}, &encoder{})
}
