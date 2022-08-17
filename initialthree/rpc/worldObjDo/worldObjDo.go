
package worldObjDo

import (
	"initialthree/cluster"
	"initialthree/cluster/addr"
	"initialthree/pkg/rpc"
	ss_rpc "initialthree/protocol/ss/rpc"
	"time"
)

type WorldObjDoReplyer struct {
	replyer_ *rpc.RPCReplyer
}

func (this *WorldObjDoReplyer) Reply(result *ss_rpc.WorldObjDoResp) {
	this.replyer_.Reply(result,nil)
}

/*
func (this *WorldObjDoReplyer) DropResponse() {
	this.replyer_.DropResponse()
}
*/

/*
func (this *WorldObjDoReplyer) Error(err error) {
	this.replyer_.Reply(nil,err)
}
*/

func (this *WorldObjDoReplyer) GetChannel() rpc.RPCChannel {
	return this.replyer_.GetChannel()
}


type WorldObjDo interface {
	OnCall(*WorldObjDoReplyer,*ss_rpc.WorldObjDoReq)
}

func Register(methodObj WorldObjDo) {
	f := func(r *rpc.RPCReplyer, arg interface{}) {
		replyer_ := &WorldObjDoReplyer{replyer_:r}
		methodObj.OnCall(replyer_,arg.(*ss_rpc.WorldObjDoReq))
	}

	cluster.RegisterMethod(&ss_rpc.WorldObjDoReq{},f)
}

func AsynCall(peer addr.LogicAddr,arg *ss_rpc.WorldObjDoReq,timeout time.Duration,cb func(*ss_rpc.WorldObjDoResp,error)) {
	callback := func(r interface{},e error) {
		if nil != r {
			cb(r.(*ss_rpc.WorldObjDoResp),e)
		} else {
			cb(nil,e)
		}
	}
	cluster.AsynCall(peer,arg,timeout,callback)
}

/*
func SyncCall(peer addr.LogicAddr,arg *ss_rpc.WorldObjDoReq,timeout time.Duration) (ret *ss_rpc.WorldObjDoResp, err error) {
	waitC := make(chan struct{})
	f := func(ret_ *ss_rpc.WorldObjDoResp, err_ error) {
		ret = ret_
		err = err_
		close(waitC)
	}
	AsynCall(peer,arg,timeout,f)
	<-waitC
	return
}
*/
