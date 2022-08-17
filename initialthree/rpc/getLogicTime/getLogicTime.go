
package getLogicTime

import (
	"initialthree/cluster"
	"initialthree/cluster/addr"
	"initialthree/pkg/rpc"
	ss_rpc "initialthree/protocol/ss/rpc"
	"time"
)

type GetLogicTimeReplyer struct {
	replyer_ *rpc.RPCReplyer
}

func (this *GetLogicTimeReplyer) Reply(result *ss_rpc.GetLogicTimeResp) {
	this.replyer_.Reply(result,nil)
}

/*
func (this *GetLogicTimeReplyer) DropResponse() {
	this.replyer_.DropResponse()
}
*/

/*
func (this *GetLogicTimeReplyer) Error(err error) {
	this.replyer_.Reply(nil,err)
}
*/

func (this *GetLogicTimeReplyer) GetChannel() rpc.RPCChannel {
	return this.replyer_.GetChannel()
}


type GetLogicTime interface {
	OnCall(*GetLogicTimeReplyer,*ss_rpc.GetLogicTimeReq)
}

func Register(methodObj GetLogicTime) {
	f := func(r *rpc.RPCReplyer, arg interface{}) {
		replyer_ := &GetLogicTimeReplyer{replyer_:r}
		methodObj.OnCall(replyer_,arg.(*ss_rpc.GetLogicTimeReq))
	}

	cluster.RegisterMethod(&ss_rpc.GetLogicTimeReq{},f)
}

func AsynCall(peer addr.LogicAddr,arg *ss_rpc.GetLogicTimeReq,timeout time.Duration,cb func(*ss_rpc.GetLogicTimeResp,error)) {
	callback := func(r interface{},e error) {
		if nil != r {
			cb(r.(*ss_rpc.GetLogicTimeResp),e)
		} else {
			cb(nil,e)
		}
	}
	cluster.AsynCall(peer,arg,timeout,callback)
}

/*
func SyncCall(peer addr.LogicAddr,arg *ss_rpc.GetLogicTimeReq,timeout time.Duration) (ret *ss_rpc.GetLogicTimeResp, err error) {
	waitC := make(chan struct{})
	f := func(ret_ *ss_rpc.GetLogicTimeResp, err_ error) {
		ret = ret_
		err = err_
		close(waitC)
	}
	AsynCall(peer,arg,timeout,f)
	<-waitC
	return
}
*/
