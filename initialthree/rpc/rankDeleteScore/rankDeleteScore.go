
package rankDeleteScore

import (
	"initialthree/cluster"
	"initialthree/cluster/addr"
	"initialthree/pkg/rpc"
	ss_rpc "initialthree/protocol/ss/rpc"
	"time"
)

type RankDeleteScoreReplyer struct {
	replyer_ *rpc.RPCReplyer
}

func (this *RankDeleteScoreReplyer) Reply(result *ss_rpc.RankDeleteScoreResp) {
	this.replyer_.Reply(result,nil)
}

/*
func (this *RankDeleteScoreReplyer) DropResponse() {
	this.replyer_.DropResponse()
}
*/

/*
func (this *RankDeleteScoreReplyer) Error(err error) {
	this.replyer_.Reply(nil,err)
}
*/

func (this *RankDeleteScoreReplyer) GetChannel() rpc.RPCChannel {
	return this.replyer_.GetChannel()
}


type RankDeleteScore interface {
	OnCall(*RankDeleteScoreReplyer,*ss_rpc.RankDeleteScoreReq)
}

func Register(methodObj RankDeleteScore) {
	f := func(r *rpc.RPCReplyer, arg interface{}) {
		replyer_ := &RankDeleteScoreReplyer{replyer_:r}
		methodObj.OnCall(replyer_,arg.(*ss_rpc.RankDeleteScoreReq))
	}

	cluster.RegisterMethod(&ss_rpc.RankDeleteScoreReq{},f)
}

func AsynCall(peer addr.LogicAddr,arg *ss_rpc.RankDeleteScoreReq,timeout time.Duration,cb func(*ss_rpc.RankDeleteScoreResp,error)) {
	callback := func(r interface{},e error) {
		if nil != r {
			cb(r.(*ss_rpc.RankDeleteScoreResp),e)
		} else {
			cb(nil,e)
		}
	}
	cluster.AsynCall(peer,arg,timeout,callback)
}

/*
func SyncCall(peer addr.LogicAddr,arg *ss_rpc.RankDeleteScoreReq,timeout time.Duration) (ret *ss_rpc.RankDeleteScoreResp, err error) {
	waitC := make(chan struct{})
	f := func(ret_ *ss_rpc.RankDeleteScoreResp, err_ error) {
		ret = ret_
		err = err_
		close(waitC)
	}
	AsynCall(peer,arg,timeout,f)
	<-waitC
	return
}
*/
