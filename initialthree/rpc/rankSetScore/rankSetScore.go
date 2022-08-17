
package rankSetScore

import (
	"initialthree/cluster"
	"initialthree/cluster/addr"
	"initialthree/pkg/rpc"
	ss_rpc "initialthree/protocol/ss/rpc"
	"time"
)

type RankSetScoreReplyer struct {
	replyer_ *rpc.RPCReplyer
}

func (this *RankSetScoreReplyer) Reply(result *ss_rpc.RankSetScoreResp) {
	this.replyer_.Reply(result,nil)
}

/*
func (this *RankSetScoreReplyer) DropResponse() {
	this.replyer_.DropResponse()
}
*/

/*
func (this *RankSetScoreReplyer) Error(err error) {
	this.replyer_.Reply(nil,err)
}
*/

func (this *RankSetScoreReplyer) GetChannel() rpc.RPCChannel {
	return this.replyer_.GetChannel()
}


type RankSetScore interface {
	OnCall(*RankSetScoreReplyer,*ss_rpc.RankSetScoreReq)
}

func Register(methodObj RankSetScore) {
	f := func(r *rpc.RPCReplyer, arg interface{}) {
		replyer_ := &RankSetScoreReplyer{replyer_:r}
		methodObj.OnCall(replyer_,arg.(*ss_rpc.RankSetScoreReq))
	}

	cluster.RegisterMethod(&ss_rpc.RankSetScoreReq{},f)
}

func AsynCall(peer addr.LogicAddr,arg *ss_rpc.RankSetScoreReq,timeout time.Duration,cb func(*ss_rpc.RankSetScoreResp,error)) {
	callback := func(r interface{},e error) {
		if nil != r {
			cb(r.(*ss_rpc.RankSetScoreResp),e)
		} else {
			cb(nil,e)
		}
	}
	cluster.AsynCall(peer,arg,timeout,callback)
}

/*
func SyncCall(peer addr.LogicAddr,arg *ss_rpc.RankSetScoreReq,timeout time.Duration) (ret *ss_rpc.RankSetScoreResp, err error) {
	waitC := make(chan struct{})
	f := func(ret_ *ss_rpc.RankSetScoreResp, err_ error) {
		ret = ret_
		err = err_
		close(waitC)
	}
	AsynCall(peer,arg,timeout,f)
	<-waitC
	return
}
*/
