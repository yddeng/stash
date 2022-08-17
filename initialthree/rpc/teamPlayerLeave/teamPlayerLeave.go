
package teamPlayerLeave

import (
	"initialthree/cluster"
	"initialthree/cluster/addr"
	"initialthree/pkg/rpc"
	ss_rpc "initialthree/protocol/ss/rpc"
	"time"
)

type TeamPlayerLeaveReplyer struct {
	replyer_ *rpc.RPCReplyer
}

func (this *TeamPlayerLeaveReplyer) Reply(result *ss_rpc.TeamPlayerLeaveResp) {
	this.replyer_.Reply(result,nil)
}

/*
func (this *TeamPlayerLeaveReplyer) DropResponse() {
	this.replyer_.DropResponse()
}
*/

/*
func (this *TeamPlayerLeaveReplyer) Error(err error) {
	this.replyer_.Reply(nil,err)
}
*/

func (this *TeamPlayerLeaveReplyer) GetChannel() rpc.RPCChannel {
	return this.replyer_.GetChannel()
}


type TeamPlayerLeave interface {
	OnCall(*TeamPlayerLeaveReplyer,*ss_rpc.TeamPlayerLeaveReq)
}

func Register(methodObj TeamPlayerLeave) {
	f := func(r *rpc.RPCReplyer, arg interface{}) {
		replyer_ := &TeamPlayerLeaveReplyer{replyer_:r}
		methodObj.OnCall(replyer_,arg.(*ss_rpc.TeamPlayerLeaveReq))
	}

	cluster.RegisterMethod(&ss_rpc.TeamPlayerLeaveReq{},f)
}

func AsynCall(peer addr.LogicAddr,arg *ss_rpc.TeamPlayerLeaveReq,timeout time.Duration,cb func(*ss_rpc.TeamPlayerLeaveResp,error)) {
	callback := func(r interface{},e error) {
		if nil != r {
			cb(r.(*ss_rpc.TeamPlayerLeaveResp),e)
		} else {
			cb(nil,e)
		}
	}
	cluster.AsynCall(peer,arg,timeout,callback)
}

/*
func SyncCall(peer addr.LogicAddr,arg *ss_rpc.TeamPlayerLeaveReq,timeout time.Duration) (ret *ss_rpc.TeamPlayerLeaveResp, err error) {
	waitC := make(chan struct{})
	f := func(ret_ *ss_rpc.TeamPlayerLeaveResp, err_ error) {
		ret = ret_
		err = err_
		close(waitC)
	}
	AsynCall(peer,arg,timeout,f)
	<-waitC
	return
}
*/
