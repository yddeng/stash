
package teamKickPlayer

import (
	"initialthree/cluster"
	"initialthree/cluster/addr"
	"initialthree/pkg/rpc"
	ss_rpc "initialthree/protocol/ss/rpc"
	"time"
)

type TeamKickPlayerReplyer struct {
	replyer_ *rpc.RPCReplyer
}

func (this *TeamKickPlayerReplyer) Reply(result *ss_rpc.TeamKickPlayerResp) {
	this.replyer_.Reply(result,nil)
}

/*
func (this *TeamKickPlayerReplyer) DropResponse() {
	this.replyer_.DropResponse()
}
*/

/*
func (this *TeamKickPlayerReplyer) Error(err error) {
	this.replyer_.Reply(nil,err)
}
*/

func (this *TeamKickPlayerReplyer) GetChannel() rpc.RPCChannel {
	return this.replyer_.GetChannel()
}


type TeamKickPlayer interface {
	OnCall(*TeamKickPlayerReplyer,*ss_rpc.TeamKickPlayerReq)
}

func Register(methodObj TeamKickPlayer) {
	f := func(r *rpc.RPCReplyer, arg interface{}) {
		replyer_ := &TeamKickPlayerReplyer{replyer_:r}
		methodObj.OnCall(replyer_,arg.(*ss_rpc.TeamKickPlayerReq))
	}

	cluster.RegisterMethod(&ss_rpc.TeamKickPlayerReq{},f)
}

func AsynCall(peer addr.LogicAddr,arg *ss_rpc.TeamKickPlayerReq,timeout time.Duration,cb func(*ss_rpc.TeamKickPlayerResp,error)) {
	callback := func(r interface{},e error) {
		if nil != r {
			cb(r.(*ss_rpc.TeamKickPlayerResp),e)
		} else {
			cb(nil,e)
		}
	}
	cluster.AsynCall(peer,arg,timeout,callback)
}

/*
func SyncCall(peer addr.LogicAddr,arg *ss_rpc.TeamKickPlayerReq,timeout time.Duration) (ret *ss_rpc.TeamKickPlayerResp, err error) {
	waitC := make(chan struct{})
	f := func(ret_ *ss_rpc.TeamKickPlayerResp, err_ error) {
		ret = ret_
		err = err_
		close(waitC)
	}
	AsynCall(peer,arg,timeout,f)
	<-waitC
	return
}
*/
