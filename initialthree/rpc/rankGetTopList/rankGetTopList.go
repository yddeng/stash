
package rankGetTopList

import (
	"initialthree/cluster"
	"initialthree/cluster/addr"
	"initialthree/pkg/rpc"
	ss_rpc "initialthree/protocol/ss/rpc"
	"time"
)

type RankGetTopListReplyer struct {
	replyer_ *rpc.RPCReplyer
}

func (this *RankGetTopListReplyer) Reply(result *ss_rpc.RankGetTopListResp) {
	this.replyer_.Reply(result,nil)
}

/*
func (this *RankGetTopListReplyer) DropResponse() {
	this.replyer_.DropResponse()
}
*/

/*
func (this *RankGetTopListReplyer) Error(err error) {
	this.replyer_.Reply(nil,err)
}
*/

func (this *RankGetTopListReplyer) GetChannel() rpc.RPCChannel {
	return this.replyer_.GetChannel()
}


type RankGetTopList interface {
	OnCall(*RankGetTopListReplyer,*ss_rpc.RankGetTopListReq)
}

func Register(methodObj RankGetTopList) {
	f := func(r *rpc.RPCReplyer, arg interface{}) {
		replyer_ := &RankGetTopListReplyer{replyer_:r}
		methodObj.OnCall(replyer_,arg.(*ss_rpc.RankGetTopListReq))
	}

	cluster.RegisterMethod(&ss_rpc.RankGetTopListReq{},f)
}

func AsynCall(peer addr.LogicAddr,arg *ss_rpc.RankGetTopListReq,timeout time.Duration,cb func(*ss_rpc.RankGetTopListResp,error)) {
	callback := func(r interface{},e error) {
		if nil != r {
			cb(r.(*ss_rpc.RankGetTopListResp),e)
		} else {
			cb(nil,e)
		}
	}
	cluster.AsynCall(peer,arg,timeout,callback)
}

/*
func SyncCall(peer addr.LogicAddr,arg *ss_rpc.RankGetTopListReq,timeout time.Duration) (ret *ss_rpc.RankGetTopListResp, err error) {
	waitC := make(chan struct{})
	f := func(ret_ *ss_rpc.RankGetTopListResp, err_ error) {
		ret = ret_
		err = err_
		close(waitC)
	}
	AsynCall(peer,arg,timeout,f)
	<-waitC
	return
}
*/
