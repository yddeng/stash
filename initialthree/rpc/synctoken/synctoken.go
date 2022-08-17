
package synctoken

import (
	"initialthree/cluster"
	"initialthree/cluster/addr"
	"initialthree/pkg/rpc"
	ss_rpc "initialthree/protocol/ss/rpc"
	"time"
)

type SynctokenReplyer struct {
	replyer_ *rpc.RPCReplyer
}

func (this *SynctokenReplyer) Reply(result *ss_rpc.SynctokenResp) {
	this.replyer_.Reply(result,nil)
}

/*
func (this *SynctokenReplyer) DropResponse() {
	this.replyer_.DropResponse()
}
*/

/*
func (this *SynctokenReplyer) Error(err error) {
	this.replyer_.Reply(nil,err)
}
*/

func (this *SynctokenReplyer) GetChannel() rpc.RPCChannel {
	return this.replyer_.GetChannel()
}


type Synctoken interface {
	OnCall(*SynctokenReplyer,*ss_rpc.SynctokenReq)
}

func Register(methodObj Synctoken) {
	f := func(r *rpc.RPCReplyer, arg interface{}) {
		replyer_ := &SynctokenReplyer{replyer_:r}
		methodObj.OnCall(replyer_,arg.(*ss_rpc.SynctokenReq))
	}

	cluster.RegisterMethod(&ss_rpc.SynctokenReq{},f)
}

func AsynCall(peer addr.LogicAddr,arg *ss_rpc.SynctokenReq,timeout time.Duration,cb func(*ss_rpc.SynctokenResp,error)) {
	callback := func(r interface{},e error) {
		if nil != r {
			cb(r.(*ss_rpc.SynctokenResp),e)
		} else {
			cb(nil,e)
		}
	}
	cluster.AsynCall(peer,arg,timeout,callback)
}

/*
func SyncCall(peer addr.LogicAddr,arg *ss_rpc.SynctokenReq,timeout time.Duration) (ret *ss_rpc.SynctokenResp, err error) {
	waitC := make(chan struct{})
	f := func(ret_ *ss_rpc.SynctokenResp, err_ error) {
		ret = ret_
		err = err_
		close(waitC)
	}
	AsynCall(peer,arg,timeout,f)
	<-waitC
	return
}
*/
