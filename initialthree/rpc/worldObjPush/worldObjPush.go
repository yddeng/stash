
package worldObjPush

import (
	"initialthree/cluster"
	"initialthree/cluster/addr"
	"initialthree/pkg/rpc"
	ss_rpc "initialthree/protocol/ss/rpc"
	"time"
)

type WorldObjPushReplyer struct {
	replyer_ *rpc.RPCReplyer
}

func (this *WorldObjPushReplyer) Reply(result *ss_rpc.WorldObjPushResp) {
	this.replyer_.Reply(result,nil)
}

/*
func (this *WorldObjPushReplyer) DropResponse() {
	this.replyer_.DropResponse()
}
*/

/*
func (this *WorldObjPushReplyer) Error(err error) {
	this.replyer_.Reply(nil,err)
}
*/

func (this *WorldObjPushReplyer) GetChannel() rpc.RPCChannel {
	return this.replyer_.GetChannel()
}


type WorldObjPush interface {
	OnCall(*WorldObjPushReplyer,*ss_rpc.WorldObjPushReq)
}

func Register(methodObj WorldObjPush) {
	f := func(r *rpc.RPCReplyer, arg interface{}) {
		replyer_ := &WorldObjPushReplyer{replyer_:r}
		methodObj.OnCall(replyer_,arg.(*ss_rpc.WorldObjPushReq))
	}

	cluster.RegisterMethod(&ss_rpc.WorldObjPushReq{},f)
}

func AsynCall(peer addr.LogicAddr,arg *ss_rpc.WorldObjPushReq,timeout time.Duration,cb func(*ss_rpc.WorldObjPushResp,error)) {
	callback := func(r interface{},e error) {
		if nil != r {
			cb(r.(*ss_rpc.WorldObjPushResp),e)
		} else {
			cb(nil,e)
		}
	}
	cluster.AsynCall(peer,arg,timeout,callback)
}

/*
func SyncCall(peer addr.LogicAddr,arg *ss_rpc.WorldObjPushReq,timeout time.Duration) (ret *ss_rpc.WorldObjPushResp, err error) {
	waitC := make(chan struct{})
	f := func(ret_ *ss_rpc.WorldObjPushResp, err_ error) {
		ret = ret_
		err = err_
		close(waitC)
	}
	AsynCall(peer,arg,timeout,f)
	<-waitC
	return
}
*/
