
package teamJoinApply

import (
	"initialthree/cluster"
	"initialthree/cluster/addr"
	"initialthree/pkg/rpc"
	ss_rpc "initialthree/protocol/ss/rpc"
	"time"
)

type TeamJoinApplyReplyer struct {
	replyer_ *rpc.RPCReplyer
}

func (this *TeamJoinApplyReplyer) Reply(result *ss_rpc.TeamJoinApplyResp) {
	this.replyer_.Reply(result,nil)
}

/*
func (this *TeamJoinApplyReplyer) DropResponse() {
	this.replyer_.DropResponse()
}
*/

/*
func (this *TeamJoinApplyReplyer) Error(err error) {
	this.replyer_.Reply(nil,err)
}
*/

func (this *TeamJoinApplyReplyer) GetChannel() rpc.RPCChannel {
	return this.replyer_.GetChannel()
}


type TeamJoinApply interface {
	OnCall(*TeamJoinApplyReplyer,*ss_rpc.TeamJoinApplyReq)
}

func Register(methodObj TeamJoinApply) {
	f := func(r *rpc.RPCReplyer, arg interface{}) {
		replyer_ := &TeamJoinApplyReplyer{replyer_:r}
		methodObj.OnCall(replyer_,arg.(*ss_rpc.TeamJoinApplyReq))
	}

	cluster.RegisterMethod(&ss_rpc.TeamJoinApplyReq{},f)
}

func AsynCall(peer addr.LogicAddr,arg *ss_rpc.TeamJoinApplyReq,timeout time.Duration,cb func(*ss_rpc.TeamJoinApplyResp,error)) {
	callback := func(r interface{},e error) {
		if nil != r {
			cb(r.(*ss_rpc.TeamJoinApplyResp),e)
		} else {
			cb(nil,e)
		}
	}
	cluster.AsynCall(peer,arg,timeout,callback)
}

/*
func SyncCall(peer addr.LogicAddr,arg *ss_rpc.TeamJoinApplyReq,timeout time.Duration) (ret *ss_rpc.TeamJoinApplyResp, err error) {
	waitC := make(chan struct{})
	f := func(ret_ *ss_rpc.TeamJoinApplyResp, err_ error) {
		ret = ret_
		err = err_
		close(waitC)
	}
	AsynCall(peer,arg,timeout,f)
	<-waitC
	return
}
*/
