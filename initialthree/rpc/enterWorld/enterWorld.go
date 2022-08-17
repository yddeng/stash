
package enterWorld

import (
	"initialthree/cluster"
	"initialthree/cluster/addr"
	"initialthree/pkg/rpc"
	ss_rpc "initialthree/protocol/ss/rpc"
	"time"
)

type EnterWorldReplyer struct {
	replyer_ *rpc.RPCReplyer
}

func (this *EnterWorldReplyer) Reply(result *ss_rpc.EnterWorldResp) {
	this.replyer_.Reply(result,nil)
}

/*
func (this *EnterWorldReplyer) DropResponse() {
	this.replyer_.DropResponse()
}
*/

/*
func (this *EnterWorldReplyer) Error(err error) {
	this.replyer_.Reply(nil,err)
}
*/

func (this *EnterWorldReplyer) GetChannel() rpc.RPCChannel {
	return this.replyer_.GetChannel()
}


type EnterWorld interface {
	OnCall(*EnterWorldReplyer,*ss_rpc.EnterWorldReq)
}

func Register(methodObj EnterWorld) {
	f := func(r *rpc.RPCReplyer, arg interface{}) {
		replyer_ := &EnterWorldReplyer{replyer_:r}
		methodObj.OnCall(replyer_,arg.(*ss_rpc.EnterWorldReq))
	}

	cluster.RegisterMethod(&ss_rpc.EnterWorldReq{},f)
}

func AsynCall(peer addr.LogicAddr,arg *ss_rpc.EnterWorldReq,timeout time.Duration,cb func(*ss_rpc.EnterWorldResp,error)) {
	callback := func(r interface{},e error) {
		if nil != r {
			cb(r.(*ss_rpc.EnterWorldResp),e)
		} else {
			cb(nil,e)
		}
	}
	cluster.AsynCall(peer,arg,timeout,callback)
}

/*
func SyncCall(peer addr.LogicAddr,arg *ss_rpc.EnterWorldReq,timeout time.Duration) (ret *ss_rpc.EnterWorldResp, err error) {
	waitC := make(chan struct{})
	f := func(ret_ *ss_rpc.EnterWorldResp, err_ error) {
		ret = ret_
		err = err_
		close(waitC)
	}
	AsynCall(peer,arg,timeout,f)
	<-waitC
	return
}
*/
