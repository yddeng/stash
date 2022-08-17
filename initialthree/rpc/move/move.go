
package move

import (
	"initialthree/cluster"
	"initialthree/cluster/addr"
	"initialthree/pkg/rpc"
	ss_rpc "initialthree/protocol/ss/rpc"
	"time"
)

type MoveReplyer struct {
	replyer_ *rpc.RPCReplyer
}

func (this *MoveReplyer) Reply(result *ss_rpc.MoveResp) {
	this.replyer_.Reply(result,nil)
}

/*
func (this *MoveReplyer) DropResponse() {
	this.replyer_.DropResponse()
}
*/

/*
func (this *MoveReplyer) Error(err error) {
	this.replyer_.Reply(nil,err)
}
*/

func (this *MoveReplyer) GetChannel() rpc.RPCChannel {
	return this.replyer_.GetChannel()
}


type Move interface {
	OnCall(*MoveReplyer,*ss_rpc.MoveReq)
}

func Register(methodObj Move) {
	f := func(r *rpc.RPCReplyer, arg interface{}) {
		replyer_ := &MoveReplyer{replyer_:r}
		methodObj.OnCall(replyer_,arg.(*ss_rpc.MoveReq))
	}

	cluster.RegisterMethod(&ss_rpc.MoveReq{},f)
}

func AsynCall(peer addr.LogicAddr,arg *ss_rpc.MoveReq,timeout time.Duration,cb func(*ss_rpc.MoveResp,error)) {
	callback := func(r interface{},e error) {
		if nil != r {
			cb(r.(*ss_rpc.MoveResp),e)
		} else {
			cb(nil,e)
		}
	}
	cluster.AsynCall(peer,arg,timeout,callback)
}

/*
func SyncCall(peer addr.LogicAddr,arg *ss_rpc.MoveReq,timeout time.Duration) (ret *ss_rpc.MoveResp, err error) {
	waitC := make(chan struct{})
	f := func(ret_ *ss_rpc.MoveResp, err_ error) {
		ret = ret_
		err = err_
		close(waitC)
	}
	AsynCall(peer,arg,timeout,f)
	<-waitC
	return
}
*/
