
package teamCreate

import (
	"initialthree/cluster"
	"initialthree/cluster/addr"
	"initialthree/pkg/rpc"
	ss_rpc "initialthree/protocol/ss/rpc"
	"time"
)

type TeamCreateReplyer struct {
	replyer_ *rpc.RPCReplyer
}

func (this *TeamCreateReplyer) Reply(result *ss_rpc.TeamCreateResp) {
	this.replyer_.Reply(result,nil)
}

/*
func (this *TeamCreateReplyer) DropResponse() {
	this.replyer_.DropResponse()
}
*/

/*
func (this *TeamCreateReplyer) Error(err error) {
	this.replyer_.Reply(nil,err)
}
*/

func (this *TeamCreateReplyer) GetChannel() rpc.RPCChannel {
	return this.replyer_.GetChannel()
}


type TeamCreate interface {
	OnCall(*TeamCreateReplyer,*ss_rpc.TeamCreateReq)
}

func Register(methodObj TeamCreate) {
	f := func(r *rpc.RPCReplyer, arg interface{}) {
		replyer_ := &TeamCreateReplyer{replyer_:r}
		methodObj.OnCall(replyer_,arg.(*ss_rpc.TeamCreateReq))
	}

	cluster.RegisterMethod(&ss_rpc.TeamCreateReq{},f)
}

func AsynCall(peer addr.LogicAddr,arg *ss_rpc.TeamCreateReq,timeout time.Duration,cb func(*ss_rpc.TeamCreateResp,error)) {
	callback := func(r interface{},e error) {
		if nil != r {
			cb(r.(*ss_rpc.TeamCreateResp),e)
		} else {
			cb(nil,e)
		}
	}
	cluster.AsynCall(peer,arg,timeout,callback)
}

/*
func SyncCall(peer addr.LogicAddr,arg *ss_rpc.TeamCreateReq,timeout time.Duration) (ret *ss_rpc.TeamCreateResp, err error) {
	waitC := make(chan struct{})
	f := func(ret_ *ss_rpc.TeamCreateResp, err_ error) {
		ret = ret_
		err = err_
		close(waitC)
	}
	AsynCall(peer,arg,timeout,f)
	<-waitC
	return
}
*/
