
package echo

import (
	"initialthree/cluster"
	"initialthree/cluster/addr"
	"initialthree/pkg/rpc"
	ss_rpc "initialthree/protocol/ss/rpc"
	"time"
)

type EchoReplyer struct {
	replyer_ *rpc.RPCReplyer
}

func (this *EchoReplyer) Reply(result *ss_rpc.EchoResp) {
	this.replyer_.Reply(result,nil)
}

/*
func (this *EchoReplyer) DropResponse() {
	this.replyer_.DropResponse()
}
*/

/*
func (this *EchoReplyer) Error(err error) {
	this.replyer_.Reply(nil,err)
}
*/

func (this *EchoReplyer) GetChannel() rpc.RPCChannel {
	return this.replyer_.GetChannel()
}


type Echo interface {
	OnCall(*EchoReplyer,*ss_rpc.EchoReq)
}

func Register(methodObj Echo) {
	f := func(r *rpc.RPCReplyer, arg interface{}) {
		replyer_ := &EchoReplyer{replyer_:r}
		methodObj.OnCall(replyer_,arg.(*ss_rpc.EchoReq))
	}

	cluster.RegisterMethod(&ss_rpc.EchoReq{},f)
}

func AsynCall(peer addr.LogicAddr,arg *ss_rpc.EchoReq,timeout time.Duration,cb func(*ss_rpc.EchoResp,error)) {
	callback := func(r interface{},e error) {
		if nil != r {
			cb(r.(*ss_rpc.EchoResp),e)
		} else {
			cb(nil,e)
		}
	}
	cluster.AsynCall(peer,arg,timeout,callback)
}

/*
func SyncCall(peer addr.LogicAddr,arg *ss_rpc.EchoReq,timeout time.Duration) (ret *ss_rpc.EchoResp, err error) {
	waitC := make(chan struct{})
	f := func(ret_ *ss_rpc.EchoResp, err_ error) {
		ret = ret_
		err = err_
		close(waitC)
	}
	AsynCall(peer,arg,timeout,f)
	<-waitC
	return
}
*/
