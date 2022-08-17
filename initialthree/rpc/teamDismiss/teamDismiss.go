
package teamDismiss

import (
	"initialthree/cluster"
	"initialthree/cluster/addr"
	"initialthree/pkg/rpc"
	ss_rpc "initialthree/protocol/ss/rpc"
	"time"
)

type TeamDismissReplyer struct {
	replyer_ *rpc.RPCReplyer
}

func (this *TeamDismissReplyer) Reply(result *ss_rpc.TeamDismissResp) {
	this.replyer_.Reply(result,nil)
}

/*
func (this *TeamDismissReplyer) DropResponse() {
	this.replyer_.DropResponse()
}
*/

/*
func (this *TeamDismissReplyer) Error(err error) {
	this.replyer_.Reply(nil,err)
}
*/

func (this *TeamDismissReplyer) GetChannel() rpc.RPCChannel {
	return this.replyer_.GetChannel()
}


type TeamDismiss interface {
	OnCall(*TeamDismissReplyer,*ss_rpc.TeamDismissReq)
}

func Register(methodObj TeamDismiss) {
	f := func(r *rpc.RPCReplyer, arg interface{}) {
		replyer_ := &TeamDismissReplyer{replyer_:r}
		methodObj.OnCall(replyer_,arg.(*ss_rpc.TeamDismissReq))
	}

	cluster.RegisterMethod(&ss_rpc.TeamDismissReq{},f)
}

func AsynCall(peer addr.LogicAddr,arg *ss_rpc.TeamDismissReq,timeout time.Duration,cb func(*ss_rpc.TeamDismissResp,error)) {
	callback := func(r interface{},e error) {
		if nil != r {
			cb(r.(*ss_rpc.TeamDismissResp),e)
		} else {
			cb(nil,e)
		}
	}
	cluster.AsynCall(peer,arg,timeout,callback)
}

/*
func SyncCall(peer addr.LogicAddr,arg *ss_rpc.TeamDismissReq,timeout time.Duration) (ret *ss_rpc.TeamDismissResp, err error) {
	waitC := make(chan struct{})
	f := func(ret_ *ss_rpc.TeamDismissResp, err_ error) {
		ret = ret_
		err = err_
		close(waitC)
	}
	AsynCall(peer,arg,timeout,f)
	<-waitC
	return
}
*/
