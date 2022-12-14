package rpc

import (
	"errors"
	"fmt"
	"initialthree/pkg/util"
	"initialthree/zaplogger"
	"sync"
	"sync/atomic"
)

type RPCReplyer struct {
	encoder RPCMessageEncoder
	channel RPCChannel
	req     *RPCRequest
	fired   int32 //防止重复Reply
	s       *RPCServer
}

func (this *RPCReplyer) Reply(ret interface{}, err error) {
	if atomic.CompareAndSwapInt32(&this.fired, 0, 1) {
		if this.req.NeedResp {
			response := &RPCResponse{Seq: this.req.Seq, Ret: ret, Err: err}
			this.reply(response)
		}
		atomic.AddInt32(&this.s.pendingCount, -1)
	}
}

func (this *RPCReplyer) DropResponse() {
	if atomic.CompareAndSwapInt32(&this.fired, 0, 1) {
		if nil != this.s {
			atomic.AddInt32(&this.s.pendingCount, -1)
		}
	}
}

func (this *RPCReplyer) reply(response RPCMessage) {
	msg, err := this.encoder.Encode(response)
	if nil != err {
		zaplogger.GetSugar().Errorf(util.FormatFileLine("Encode rpc response error:%s\n", err.Error()))
		return
	}
	err = this.channel.SendResponse(msg)
	if nil != err {
		zaplogger.GetSugar().Errorf(util.FormatFileLine("send rpc response to (%s) error:%s\n", this.channel.Name(), err.Error()))
	}
}

func (this *RPCReplyer) GetChannel() RPCChannel {
	return this.channel
}

type RPCMethodHandler func(*RPCReplyer, interface{})

type RPCServer struct {
	sync.RWMutex
	encoder            RPCMessageEncoder
	decoder            RPCMessageDecoder
	methods            sync.Map
	lastSeq            uint64
	pendingCount       int32
	errOnMissingMethod atomic.Value
}

func (this *RPCServer) PendingCount() int32 {
	return atomic.LoadInt32(&this.pendingCount)
}

func (this *RPCServer) SetErrorCodeOnMissingMethod(errOnMissingMethod error) {
	if nil != errOnMissingMethod {
		this.errOnMissingMethod.Store(errOnMissingMethod)
	}
}

/*
 *  在服务停止的情况下直接向对端返回错误响应
 */
func (this *RPCServer) OnServiceStop(channel RPCChannel, message interface{}, errorCode error) {

	msg, err := this.decoder.Decode(message)
	if nil != err {
		zaplogger.GetSugar().Error(util.FormatFileLine("RPCServer rpc message from(%s) decode err:%s\n", channel.Name(), err.Error()))
		return
	}

	switch msg.(type) {
	case *RPCRequest:
		req := msg.(*RPCRequest)

		response := &RPCResponse{Seq: req.Seq, Err: errorCode}

		msg, err := this.encoder.Encode(response)
		if nil != err {
			zaplogger.GetSugar().Errorf(util.FormatFileLine("Encode rpc response error:%s\n", err.Error()))
		} else {
			err = channel.SendResponse(msg)

			if nil != err {
				zaplogger.GetSugar().Errorf(util.FormatFileLine("send rpc response to (%s) error:%s\n", channel.Name(), err.Error()))
			}
		}
	}
}

func (this *RPCServer) RegisterMethod(name string, method RPCMethodHandler) error {
	if name == "" {
		return errors.New("nams is nil")
	} else if nil == method {
		return errors.New("method == nil")
	} else {
		if _, ok := this.methods.LoadOrStore(name, method); ok {
			return fmt.Errorf("duplicate method:%s", name)
		}
		return nil
	}
}

func (this *RPCServer) UnRegisterMethod(name string) {
	this.methods.Delete(name)
}

func (this *RPCServer) callMethod(method RPCMethodHandler, replyer *RPCReplyer, arg interface{}) {
	defer func() {
		if r := recover(); r != nil {
			replyer.Reply(nil, errors.New(fmt.Sprintf("%v", r)))
		}
	}()
	method(replyer, arg)
}

/*
*   如果需要单线程处理,可以提供eventQueue
 */
func (this *RPCServer) OnRPCMessage(channel RPCChannel, message interface{}) {
	msg, err := this.decoder.Decode(message)
	if nil != err {
		zaplogger.GetSugar().Error(util.FormatFileLine("RPCServer rpc message from(%s) decode err:%s\n", channel.Name(), err.Error()))
		return
	}

	switch msg.(type) {
	case *RPCRequest:
		{
			req := msg.(*RPCRequest)
			m, ok := this.methods.Load(req.Method)
			replyer := &RPCReplyer{encoder: this.encoder, channel: channel, req: req, s: this}
			atomic.AddInt32(&this.pendingCount, 1)

			if !ok {
				zaplogger.GetSugar().Errorf(util.FormatFileLine("rpc request from(%s) invaild method %s\n", channel.Name(), req.Method))
				errOnMissingMethod := this.errOnMissingMethod.Load()
				if nil != errOnMissingMethod {
					replyer.Reply(nil, errOnMissingMethod.(error))
				} else {
					replyer.Reply(nil, fmt.Errorf("invaild method:%s", req.Method))
				}
			} else {
				this.callMethod(m.(RPCMethodHandler), replyer, req.Arg)
			}
		}
	}

}

func NewRPCServer(decoder RPCMessageDecoder, encoder RPCMessageEncoder) *RPCServer {
	if nil == decoder || nil == encoder {
		return nil
	} else {
		return &RPCServer{
			decoder: decoder,
			encoder: encoder,
		}
	}
}
