
package teamGetNearTeam

import (
	"initialthree/cluster"
	"initialthree/cluster/addr"
	"initialthree/pkg/rpc"
	ss_rpc "initialthree/protocol/ss/rpc"
	"time"
)

type TeamGetNearTeamReplyer struct {
	replyer_ *rpc.RPCReplyer
}

func (this *TeamGetNearTeamReplyer) Reply(result *ss_rpc.TeamGetNearTeamResp) {
	this.replyer_.Reply(result,nil)
}

/*
func (this *TeamGetNearTeamReplyer) DropResponse() {
	this.replyer_.DropResponse()
}
*/

/*
func (this *TeamGetNearTeamReplyer) Error(err error) {
	this.replyer_.Reply(nil,err)
}
*/

func (this *TeamGetNearTeamReplyer) GetChannel() rpc.RPCChannel {
	return this.replyer_.GetChannel()
}


type TeamGetNearTeam interface {
	OnCall(*TeamGetNearTeamReplyer,*ss_rpc.TeamGetNearTeamReq)
}

func Register(methodObj TeamGetNearTeam) {
	f := func(r *rpc.RPCReplyer, arg interface{}) {
		replyer_ := &TeamGetNearTeamReplyer{replyer_:r}
		methodObj.OnCall(replyer_,arg.(*ss_rpc.TeamGetNearTeamReq))
	}

	cluster.RegisterMethod(&ss_rpc.TeamGetNearTeamReq{},f)
}

func AsynCall(peer addr.LogicAddr,arg *ss_rpc.TeamGetNearTeamReq,timeout time.Duration,cb func(*ss_rpc.TeamGetNearTeamResp,error)) {
	callback := func(r interface{},e error) {
		if nil != r {
			cb(r.(*ss_rpc.TeamGetNearTeamResp),e)
		} else {
			cb(nil,e)
		}
	}
	cluster.AsynCall(peer,arg,timeout,callback)
}

/*
func SyncCall(peer addr.LogicAddr,arg *ss_rpc.TeamGetNearTeamReq,timeout time.Duration) (ret *ss_rpc.TeamGetNearTeamResp, err error) {
	waitC := make(chan struct{})
	f := func(ret_ *ss_rpc.TeamGetNearTeamResp, err_ error) {
		ret = ret_
		err = err_
		close(waitC)
	}
	AsynCall(peer,arg,timeout,f)
	<-waitC
	return
}
*/
