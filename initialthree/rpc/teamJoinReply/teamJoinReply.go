
package teamJoinReply

import (
	"initialthree/cluster"
	"initialthree/cluster/addr"
	"initialthree/pkg/rpc"
	ss_rpc "initialthree/protocol/ss/rpc"
	"time"
)

type TeamJoinReplyReplyer struct {
	replyer_ *rpc.RPCReplyer
}

func (this *TeamJoinReplyReplyer) Reply(result *ss_rpc.TeamJoinReplyResp) {
	this.replyer_.Reply(result,nil)
}

/*
func (this *TeamJoinReplyReplyer) DropResponse() {
	this.replyer_.DropResponse()
}
*/

/*
func (this *TeamJoinReplyReplyer) Error(err error) {
	this.replyer_.Reply(nil,err)
}
*/

func (this *TeamJoinReplyReplyer) GetChannel() rpc.RPCChannel {
	return this.replyer_.GetChannel()
}


type TeamJoinReply interface {
	OnCall(*TeamJoinReplyReplyer,*ss_rpc.TeamJoinReplyReq)
}

func Register(methodObj TeamJoinReply) {
	f := func(r *rpc.RPCReplyer, arg interface{}) {
		replyer_ := &TeamJoinReplyReplyer{replyer_:r}
		methodObj.OnCall(replyer_,arg.(*ss_rpc.TeamJoinReplyReq))
	}

	cluster.RegisterMethod(&ss_rpc.TeamJoinReplyReq{},f)
}

func AsynCall(peer addr.LogicAddr,arg *ss_rpc.TeamJoinReplyReq,timeout time.Duration,cb func(*ss_rpc.TeamJoinReplyResp,error)) {
	callback := func(r interface{},e error) {
		if nil != r {
			cb(r.(*ss_rpc.TeamJoinReplyResp),e)
		} else {
			cb(nil,e)
		}
	}
	cluster.AsynCall(peer,arg,timeout,callback)
}

/*
func SyncCall(peer addr.LogicAddr,arg *ss_rpc.TeamJoinReplyReq,timeout time.Duration) (ret *ss_rpc.TeamJoinReplyResp, err error) {
	waitC := make(chan struct{})
	f := func(ret_ *ss_rpc.TeamJoinReplyResp, err_ error) {
		ret = ret_
		err = err_
		close(waitC)
	}
	AsynCall(peer,arg,timeout,f)
	<-waitC
	return
}
*/
