
package rankGetRank

import (
	"initialthree/cluster"
	"initialthree/cluster/addr"
	"initialthree/pkg/rpc"
	ss_rpc "initialthree/protocol/ss/rpc"
	"time"
)

type RankGetRankReplyer struct {
	replyer_ *rpc.RPCReplyer
}

func (this *RankGetRankReplyer) Reply(result *ss_rpc.RankGetRankResp) {
	this.replyer_.Reply(result,nil)
}

/*
func (this *RankGetRankReplyer) DropResponse() {
	this.replyer_.DropResponse()
}
*/

/*
func (this *RankGetRankReplyer) Error(err error) {
	this.replyer_.Reply(nil,err)
}
*/

func (this *RankGetRankReplyer) GetChannel() rpc.RPCChannel {
	return this.replyer_.GetChannel()
}


type RankGetRank interface {
	OnCall(*RankGetRankReplyer,*ss_rpc.RankGetRankReq)
}

func Register(methodObj RankGetRank) {
	f := func(r *rpc.RPCReplyer, arg interface{}) {
		replyer_ := &RankGetRankReplyer{replyer_:r}
		methodObj.OnCall(replyer_,arg.(*ss_rpc.RankGetRankReq))
	}

	cluster.RegisterMethod(&ss_rpc.RankGetRankReq{},f)
}

func AsynCall(peer addr.LogicAddr,arg *ss_rpc.RankGetRankReq,timeout time.Duration,cb func(*ss_rpc.RankGetRankResp,error)) {
	callback := func(r interface{},e error) {
		if nil != r {
			cb(r.(*ss_rpc.RankGetRankResp),e)
		} else {
			cb(nil,e)
		}
	}
	cluster.AsynCall(peer,arg,timeout,callback)
}

/*
func SyncCall(peer addr.LogicAddr,arg *ss_rpc.RankGetRankReq,timeout time.Duration) (ret *ss_rpc.RankGetRankResp, err error) {
	waitC := make(chan struct{})
	f := func(ret_ *ss_rpc.RankGetRankResp, err_ error) {
		ret = ret_
		err = err_
		close(waitC)
	}
	AsynCall(peer,arg,timeout,f)
	<-waitC
	return
}
*/
