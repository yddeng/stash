package transaction

import (
	"container/list"
	"fmt"
	"github.com/golang/protobuf/proto"
	"initialthree/cluster/priority"
	"initialthree/codec/cs"
	"initialthree/pkg/event"
	"initialthree/pkg/pcall"
	"initialthree/protocol/cs/message"
	"initialthree/zaplogger"
	"reflect"
	"runtime"
	"time"
)

const (
	TransTimeoutSec = 10
	TransTimeoutMs  = 10 * 1000
)

type TransactionBase struct {
	mgr     *TransactionMgr
	trans   Transaction
	start   time.Time
	expired time.Time
	req     *cs.Message
}

func (this *TransactionBase) setTransBase(mgr *TransactionMgr, trans Transaction, req *cs.Message) {
	this.mgr = mgr
	this.trans = trans
	this.req = req
}

func (this *TransactionBase) setExpiredTime(start, expired time.Time) {
	this.start = start
	this.expired = expired
}

func (this *TransactionBase) getRequest() *cs.Message {
	return this.req
}

func (this *TransactionBase) getStartTime() time.Time {
	return this.start
}

func (this *TransactionBase) GetGameModule() int { return 0 }

func (this *TransactionBase) IsTimeout() bool {
	if !this.expired.IsZero() && time.Now().After(this.expired) {
		return true
	}
	return false
}

func (this *TransactionBase) EndTrans(resp proto.Message, errCode message.ErrCode) {
	this.mgr.callEnd(this.trans, resp, errCode)
}

type wrapFunc func(args ...interface{}) []interface{}

/*
 *   异步函数包装器，被包装异步函数的最后一个参数是回调函数
 *   包装后的函数保证回调只在Transaction没有被End才会被调用。例如,Transaction已经超时End,此时异步返回，则回调不会被调用。
 */
func (this *TransactionBase) AsynWrap(fn interface{}) wrapFunc {
	oriF := reflect.ValueOf(fn)

	if oriF.Kind() != reflect.Func {
		return nil
	}

	if this.mgr == nil {
		return nil
	}

	return func(args ...interface{}) []interface{} {
		oriCallBack := reflect.ValueOf(args[len(args)-1])
		args[len(args)-1] = reflect.MakeFunc(oriCallBack.Type(), func(in []reflect.Value) []reflect.Value {
			//如果trans已经end就不再调用callback
			if this.mgr.CheckTrans(this.trans) {
				return oriCallBack.Call(in)
			}
			return nil
		}).Interface()
		ret, err := pcall.Call(fn, args...)
		if nil != err {
			panic(err.Error())
		}
		return ret
	}
}

type Transaction interface {
	GetGameModule() int
	GetModuleName() string
	Begin()
	IsTimeout() bool

	setTransBase(mgr *TransactionMgr, trans Transaction, req *cs.Message)
	setExpiredTime(time.Time, time.Time)
	getRequest() *cs.Message
	getStartTime() time.Time
}

type transNode struct {
	trans   Transaction
	timeout time.Duration
}

type TransactionMgr struct {
	current       Transaction //当前正在处理的Transaction
	transQueue    *list.List
	queue         *event.EventQueue
	closeCallBack func()
	callTransEnd  func(trans Transaction, req *cs.Message, resp proto.Message, errCode message.ErrCode, usedTime time.Duration)
}

func New(queue *event.EventQueue) *TransactionMgr {
	return &TransactionMgr{
		queue:      queue,
		transQueue: list.New(),
	}
}

func protectcall(f func()) (ok bool) {

	defer func() {
		if r := recover(); r != nil {
			buf := make([]byte, 65535)
			l := runtime.Stack(buf, false)
			zaplogger.GetSugar().Errorf("%v: %s\n", r, buf[:l])
			ok = false
			return
		}
		ok = true
	}()

	f()
	return
}

func (this *TransactionMgr) SetCallTransEnd(callTransEnd func(trans Transaction, req *cs.Message, resp proto.Message, errCode message.ErrCode, usedTime time.Duration)) {
	this.callTransEnd = callTransEnd
}

func (this *TransactionMgr) Close(closeCallBack func()) {
	if nil != closeCallBack {
		this.closeCallBack = closeCallBack
		if nil == this.current {
			protectcall(closeCallBack)
		}
	}
}

func (this *TransactionMgr) callBegin(trans Transaction, timeout time.Duration) {
	trans.setExpiredTime(time.Now(), time.Now().Add(timeout))
	if !protectcall(trans.Begin) {
		this.callEnd(trans, nil, message.ErrCode_ERROR)
	}
}

func (this *TransactionMgr) callEnd(trans Transaction, resp proto.Message, errCode message.ErrCode) {
	if this.CheckTrans(trans) {

		usedTime := time.Now().Sub(trans.getStartTime())
		if nil != this.callTransEnd {
			protectcall(func() { this.callTransEnd(trans, trans.getRequest(), resp, errCode, usedTime) })
		}

		e := this.transQueue.Front()
		if nil != e {
			this.transQueue.Remove(e)
			this.current = e.Value.(transNode).trans
			this.queue.PostNoWait(priority.MID, this.callBegin, this.current, e.Value.(transNode).timeout)
		} else {
			this.current = nil
			if nil != this.closeCallBack {
				protectcall(this.closeCallBack)
			}
		}
	}
}

func (this *TransactionMgr) PushTrans(trans Transaction, req *cs.Message, timeout time.Duration) error {
	if nil != this.closeCallBack {
		return fmt.Errorf("closed")
	} else {
		/*
		 * 如果当前没有Transaction正在执行立即开始执行，否则将trans添加到待处理队列中
		 */
		trans.setTransBase(this, trans, req)
		if nil == this.current {
			this.current = trans
			this.callBegin(trans, timeout)
		} else {
			this.transQueue.PushBack(transNode{
				trans:   trans,
				timeout: timeout,
			})
		}
		return nil
	}
}

func (this *TransactionMgr) Tick(now time.Time) {
	if this.current != nil && this.current.IsTimeout() {
		this.callEnd(this.current, nil, message.ErrCode_RETRY)
	}
	if this.closeCallBack != nil && this.current == nil {
		protectcall(this.closeCallBack)
	}
}

func (this *TransactionMgr) CheckTrans(trans Transaction) bool {
	return this.current == trans
}
