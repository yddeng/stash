
package teamPlayerGetFromGame

import (
	"initialthree/cluster"
	"initialthree/cluster/addr"
	"initialthree/pkg/rpc"
	ss_rpc "initialthree/protocol/ss/rpc"
	"time"
)

type TeamPlayerGetFromGameReplyer struct {
	replyer_ *rpc.RPCReplyer
}

func (this *TeamPlayerGetFromGameReplyer) Reply(result *ss_rpc.TeamPlayerGetFromGameResp) {
	this.replyer_.Reply(result,nil)
}

/*
func (this *TeamPlayerGetFromGameReplyer) DropResponse() {
	this.replyer_.DropResponse()
}
*/

/*
func (this *TeamPlayerGetFromGameReplyer) Error(err error) {
	this.replyer_.Reply(nil,err)
}
*/

func (this *TeamPlayerGetFromGameReplyer) GetChannel() rpc.RPCChannel {
	return this.replyer_.GetChannel()
}


type TeamPlayerGetFromGame interface {
	OnCall(*TeamPlayerGetFromGameReplyer,*ss_rpc.TeamPlayerGetFromGameReq)
}

func Register(methodObj TeamPlayerGetFromGame) {
	f := func(r *rpc.RPCReplyer, arg interface{}) {
		replyer_ := &TeamPlayerGetFromGameReplyer{replyer_:r}
		methodObj.OnCall(replyer_,arg.(*ss_rpc.TeamPlayerGetFromGameReq))
	}

	cluster.RegisterMethod(&ss_rpc.TeamPlayerGetFromGameReq{},f)
}

func AsynCall(peer addr.LogicAddr,arg *ss_rpc.TeamPlayerGetFromGameReq,timeout time.Duration,cb func(*ss_rpc.TeamPlayerGetFromGameResp,error)) {
	callback := func(r interface{},e error) {
		if nil != r {
			cb(r.(*ss_rpc.TeamPlayerGetFromGameResp),e)
		} else {
			cb(nil,e)
		}
	}
	cluster.AsynCall(peer,arg,timeout,callback)
}

/*
func SyncCall(peer addr.LogicAddr,arg *ss_rpc.TeamPlayerGetFromGameReq,timeout time.Duration) (ret *ss_rpc.TeamPlayerGetFromGameResp, err error) {
	waitC := make(chan struct{})
	f := func(ret_ *ss_rpc.TeamPlayerGetFromGameResp, err_ error) {
		ret = ret_
		err = err_
		close(waitC)
	}
	AsynCall(peer,arg,timeout,f)
	<-waitC
	return
}
*/
