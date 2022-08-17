
package rankCreate

import (
	"initialthree/cluster"
	"initialthree/cluster/addr"
	"initialthree/pkg/rpc"
	ss_rpc "initialthree/protocol/ss/rpc"
	"time"
)

type RankCreateReplyer struct {
	replyer_ *rpc.RPCReplyer
}

func (this *RankCreateReplyer) Reply(result *ss_rpc.RankCreateResp) {
	this.replyer_.Reply(result,nil)
}

/*
func (this *RankCreateReplyer) DropResponse() {
	this.replyer_.DropResponse()
}
*/

/*
func (this *RankCreateReplyer) Error(err error) {
	this.replyer_.Reply(nil,err)
}
*/

func (this *RankCreateReplyer) GetChannel() rpc.RPCChannel {
	return this.replyer_.GetChannel()
}


type RankCreate interface {
	OnCall(*RankCreateReplyer,*ss_rpc.RankCreateReq)
}

func Register(methodObj RankCreate) {
	f := func(r *rpc.RPCReplyer, arg interface{}) {
		replyer_ := &RankCreateReplyer{replyer_:r}
		methodObj.OnCall(replyer_,arg.(*ss_rpc.RankCreateReq))
	}

	cluster.RegisterMethod(&ss_rpc.RankCreateReq{},f)
}

func AsynCall(peer addr.LogicAddr,arg *ss_rpc.RankCreateReq,timeout time.Duration,cb func(*ss_rpc.RankCreateResp,error)) {
	callback := func(r interface{},e error) {
		if nil != r {
			cb(r.(*ss_rpc.RankCreateResp),e)
		} else {
			cb(nil,e)
		}
	}
	cluster.AsynCall(peer,arg,timeout,callback)
}

/*
func SyncCall(peer addr.LogicAddr,arg *ss_rpc.RankCreateReq,timeout time.Duration) (ret *ss_rpc.RankCreateResp, err error) {
	waitC := make(chan struct{})
	f := func(ret_ *ss_rpc.RankCreateResp, err_ error) {
		ret = ret_
		err = err_
		close(waitC)
	}
	AsynCall(peer,arg,timeout,f)
	<-waitC
	return
}
*/
