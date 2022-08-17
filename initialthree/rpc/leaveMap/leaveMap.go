
package leaveMap

import (
	"initialthree/cluster"
	"initialthree/cluster/addr"
	"initialthree/pkg/rpc"
	ss_rpc "initialthree/protocol/ss/rpc"
	"time"
)

type LeaveMapReplyer struct {
	replyer_ *rpc.RPCReplyer
}

func (this *LeaveMapReplyer) Reply(result *ss_rpc.LeaveMapResp) {
	this.replyer_.Reply(result,nil)
}

/*
func (this *LeaveMapReplyer) DropResponse() {
	this.replyer_.DropResponse()
}
*/

/*
func (this *LeaveMapReplyer) Error(err error) {
	this.replyer_.Reply(nil,err)
}
*/

func (this *LeaveMapReplyer) GetChannel() rpc.RPCChannel {
	return this.replyer_.GetChannel()
}


type LeaveMap interface {
	OnCall(*LeaveMapReplyer,*ss_rpc.LeaveMapReq)
}

func Register(methodObj LeaveMap) {
	f := func(r *rpc.RPCReplyer, arg interface{}) {
		replyer_ := &LeaveMapReplyer{replyer_:r}
		methodObj.OnCall(replyer_,arg.(*ss_rpc.LeaveMapReq))
	}

	cluster.RegisterMethod(&ss_rpc.LeaveMapReq{},f)
}

func AsynCall(peer addr.LogicAddr,arg *ss_rpc.LeaveMapReq,timeout time.Duration,cb func(*ss_rpc.LeaveMapResp,error)) {
	callback := func(r interface{},e error) {
		if nil != r {
			cb(r.(*ss_rpc.LeaveMapResp),e)
		} else {
			cb(nil,e)
		}
	}
	cluster.AsynCall(peer,arg,timeout,callback)
}

/*
func SyncCall(peer addr.LogicAddr,arg *ss_rpc.LeaveMapReq,timeout time.Duration) (ret *ss_rpc.LeaveMapResp, err error) {
	waitC := make(chan struct{})
	f := func(ret_ *ss_rpc.LeaveMapResp, err_ error) {
		ret = ret_
		err = err_
		close(waitC)
	}
	AsynCall(peer,arg,timeout,f)
	<-waitC
	return
}
*/
