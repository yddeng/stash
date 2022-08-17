
package enterMap

import (
	"initialthree/cluster"
	"initialthree/cluster/addr"
	"initialthree/pkg/rpc"
	ss_rpc "initialthree/protocol/ss/rpc"
	"time"
)

type EnterMapReplyer struct {
	replyer_ *rpc.RPCReplyer
}

func (this *EnterMapReplyer) Reply(result *ss_rpc.EnterMapResp) {
	this.replyer_.Reply(result,nil)
}

/*
func (this *EnterMapReplyer) DropResponse() {
	this.replyer_.DropResponse()
}
*/

/*
func (this *EnterMapReplyer) Error(err error) {
	this.replyer_.Reply(nil,err)
}
*/

func (this *EnterMapReplyer) GetChannel() rpc.RPCChannel {
	return this.replyer_.GetChannel()
}


type EnterMap interface {
	OnCall(*EnterMapReplyer,*ss_rpc.EnterMapReq)
}

func Register(methodObj EnterMap) {
	f := func(r *rpc.RPCReplyer, arg interface{}) {
		replyer_ := &EnterMapReplyer{replyer_:r}
		methodObj.OnCall(replyer_,arg.(*ss_rpc.EnterMapReq))
	}

	cluster.RegisterMethod(&ss_rpc.EnterMapReq{},f)
}

func AsynCall(peer addr.LogicAddr,arg *ss_rpc.EnterMapReq,timeout time.Duration,cb func(*ss_rpc.EnterMapResp,error)) {
	callback := func(r interface{},e error) {
		if nil != r {
			cb(r.(*ss_rpc.EnterMapResp),e)
		} else {
			cb(nil,e)
		}
	}
	cluster.AsynCall(peer,arg,timeout,callback)
}

/*
func SyncCall(peer addr.LogicAddr,arg *ss_rpc.EnterMapReq,timeout time.Duration) (ret *ss_rpc.EnterMapResp, err error) {
	waitC := make(chan struct{})
	f := func(ret_ *ss_rpc.EnterMapResp, err_ error) {
		ret = ret_
		err = err_
		close(waitC)
	}
	AsynCall(peer,arg,timeout,f)
	<-waitC
	return
}
*/
