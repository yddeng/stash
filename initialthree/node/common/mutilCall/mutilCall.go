package mutilCall

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"initialthree/cluster"
	"initialthree/cluster/addr"
	"time"
)

type context struct {
	id      int
	peer    addr.LogicAddr
	arg     proto.Message
	timeout time.Duration
}

type Result struct {
	Req interface{}
	Ret interface{}
	Err error
}

type MutilCall struct {
	contexts []*context
}

func NewMutilCall() *MutilCall {
	return &MutilCall{
		contexts: []*context{},
	}
}

func (this *MutilCall) Add(peer addr.LogicAddr, arg proto.Message, timeout time.Duration) {

	c := &context{
		peer:    peer,
		arg:     arg,
		timeout: timeout,
	}
	this.contexts = append(this.contexts, c)
}

/*
*  所有调用均返回(无论成功失败)才执行回调
 */
func (this *MutilCall) All(callback func([]Result)) {

	if nil == callback {
		return
	}

	results := make([]Result, len(this.contexts))
	remain := len(this.contexts)
	for i, v := range this.contexts {
		idx := i
		val := v
		results[idx].Req = val.arg
		cluster.AsynCall(val.peer, val.arg, val.timeout, func(r interface{}, err error) {
			results[idx].Ret = r
			results[idx].Err = err
			remain--
			if remain == 0 {
				callback(results)
			}
		})
	}
}

/*
*  任意一个调用成功/所有均失败执行回调
 */

func (this *MutilCall) AnySuccessOrAllFail(callback func(interface{}, error)) {

	if nil == callback {
		return
	}

	remain := len(this.contexts)
	for _, v := range this.contexts {
		//idx := i
		cluster.AsynCall(v.peer, v.arg, v.timeout, func(r interface{}, err error) {
			if callback == nil {
				return
			}
			if err == nil {
				callback_back := callback
				callback = nil
				callback_back(r, nil)
			} else {
				remain--
				if remain == 0 {
					callback(nil, fmt.Errorf("all failed"))
				}
			}
		})
	}
}

func (this *MutilCall) AnyFailOrAllSuccess(callback func([]Result, error)) {

	if nil == callback {
		return
	}

	results := make([]Result, len(this.contexts))
	remain := len(this.contexts)
	for i, v := range this.contexts {
		idx := i
		val := v
		results[idx].Req = val.arg
		cluster.AsynCall(val.peer, val.arg, val.timeout, func(r interface{}, err error) {

			if callback == nil {
				return
			}
			if err != nil {
				callback_back := callback
				callback = nil
				callback_back(nil, fmt.Errorf("AnyFail"))
			} else {
				results[idx].Ret = r
				results[idx].Err = err
				remain--
				if remain == 0 {
					callback(results, nil)
				}
			}

		})
	}
}
