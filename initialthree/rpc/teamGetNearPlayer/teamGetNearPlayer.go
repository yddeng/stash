
package teamGetNearPlayer

import (
	"initialthree/cluster"
	"initialthree/cluster/addr"
	"initialthree/pkg/rpc"
	ss_rpc "initialthree/protocol/ss/rpc"
	"time"
)

type TeamGetNearPlayerReplyer struct {
	replyer_ *rpc.RPCReplyer
}

func (this *TeamGetNearPlayerReplyer) Reply(result *ss_rpc.TeamGetNearPlayerResp) {
	this.replyer_.Reply(result,nil)
}

/*
func (this *TeamGetNearPlayerReplyer) DropResponse() {
	this.replyer_.DropResponse()
}
*/

/*
func (this *TeamGetNearPlayerReplyer) Error(err error) {
	this.replyer_.Reply(nil,err)
}
*/

func (this *TeamGetNearPlayerReplyer) GetChannel() rpc.RPCChannel {
	return this.replyer_.GetChannel()
}


type TeamGetNearPlayer interface {
	OnCall(*TeamGetNearPlayerReplyer,*ss_rpc.TeamGetNearPlayerReq)
}

func Register(methodObj TeamGetNearPlayer) {
	f := func(r *rpc.RPCReplyer, arg interface{}) {
		replyer_ := &TeamGetNearPlayerReplyer{replyer_:r}
		methodObj.OnCall(replyer_,arg.(*ss_rpc.TeamGetNearPlayerReq))
	}

	cluster.RegisterMethod(&ss_rpc.TeamGetNearPlayerReq{},f)
}

func AsynCall(peer addr.LogicAddr,arg *ss_rpc.TeamGetNearPlayerReq,timeout time.Duration,cb func(*ss_rpc.TeamGetNearPlayerResp,error)) {
	callback := func(r interface{},e error) {
		if nil != r {
			cb(r.(*ss_rpc.TeamGetNearPlayerResp),e)
		} else {
			cb(nil,e)
		}
	}
	cluster.AsynCall(peer,arg,timeout,callback)
}

/*
func SyncCall(peer addr.LogicAddr,arg *ss_rpc.TeamGetNearPlayerReq,timeout time.Duration) (ret *ss_rpc.TeamGetNearPlayerResp, err error) {
	waitC := make(chan struct{})
	f := func(ret_ *ss_rpc.TeamGetNearPlayerResp, err_ error) {
		ret = ret_
		err = err_
		close(waitC)
	}
	AsynCall(peer,arg,timeout,f)
	<-waitC
	return
}
*/
