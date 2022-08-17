package robot

import (
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	fnet "github.com/sniperHW/flyfish/pkg/net"
	codecs "initialthree/codec/cs"
	"initialthree/robot/net"
	"initialthree/robot/robot/module"
	"initialthree/robot/statistics"
	"initialthree/zaplogger"
	"time"
)

func onSessionClose(session *fnet.Socket, err error) {
	// zaplogger.GetSugar().Debugf("session(%s) close: %v", session.RemoteAddr().String(), err)
	r, ok := session.GetUserData().(*Robot)
	if ok {
		r.PostNoWait(onSessionClosed, session, err)
	}
}

func onSessionClosed(session *fnet.Socket, err error) {
	robot := session.GetUserData().(*Robot)
	if session != robot.session {
		return
	}

	robot.onDisconnect(err)
}

func onSessionRecv(session *fnet.Socket, o interface{}) {
	r, ok := session.GetUserData().(*Robot)
	if ok {
		r.PostNoWait(dispatchMsg, session, o, time.Now())
	} else {
		// 会话没有关联机器人
		zaplogger.GetSugar().Warn("session no associated robot")
	}
}

func dispatchMsg(session *fnet.Socket, o interface{}, recvTime time.Time) {
	robot := session.GetUserData().(*Robot)
	if session != robot.session {
		// 网络会话匹配，消息才有效
		return
	}

	msg := o.(*codecs.Message)
	dispatchRobotMsg(robot, msg)
	robot.onMessage(msg, recvTime)
}

func (r *Robot) SetSession(session *fnet.Socket) {
	r.session = session
	if session != nil {
		session.SetUserData(r)
		session.SetCloseCallBack(onSessionClose)
		session.BeginRecv(onSessionRecv)
	}
	r.msgSeriNo = 0
	r.msgCBList.Init()
}

func (r *Robot) CloseSession() {
	if r.session != nil {
		r.session.Close(nil, 0)
		r.SetSession(nil)
	}
}

var errStop = errors.New("stop")

func (r *Robot) stopSession() {
	if r.session != nil {
		r.session.Close(errStop, 0)
		r.SetSession(nil)
	}
}

func (r *Robot) onDisconnect(err error) {
	r.Infof("disconnected: %v", err)
	r.SetSession(nil)
	r.status.Reset()
	r.bevEntity.Stop()
}

var errActiveDisconnect = errors.New("active disconnect")

func (r *Robot) ActiveDisconnect() {
	if r.session != nil {
		r.Info("active disconnect")
		r.session.Close(errActiveDisconnect, 0)
	}
}

type msgCallback = func(RobotI, *codecs.Message) bool

type msgCallbackNode struct {
	cb   msgCallback
	next *msgCallbackNode
}

type msgCallbackChain struct {
	msg      *codecs.Message
	sendTime time.Time
	head     *msgCallbackNode
}

func newMsgCallbackChain(msg *codecs.Message, sendTime time.Time, cb ...msgCallback) *msgCallbackChain {
	if len(cb) == 0 {
		panic("empty callbacks")
	}

	c := &msgCallbackChain{msg: msg, sendTime: sendTime}
	c.head = &msgCallbackNode{cb: nil}
	tail := c.head
	for _, v := range cb {
		if v == nil {
			panic("nil callback")
		}
		n := &msgCallbackNode{
			cb:   v,
			next: nil,
		}
		tail.next = n
		tail = n
	}

	return c
}

func (c *msgCallbackChain) call(r RobotI, msg *codecs.Message) {
	for p, b := c.head.next, true; p != nil && b; b, p = p.cb(r, msg), p.next {
	}
}

func (r *Robot) SendMessage(msg proto.Message, cb ...msgCallback) {
	r.msgSeriNo += 1
	msg_ := codecs.NewMessage(r.msgSeriNo, msg)
	r.session.Send(msg_)
	if len(cb) > 0 {
		r.msgCBList.PushBack(newMsgCallbackChain(msg_, time.Now(), cb...))
	}
	r.lastSendMsgTime = time.Now()
}

func (r *Robot) onMessage(msg *codecs.Message, recvTime time.Time) {
	// if r.curState != nil {
	// 	r.curState.onMessage(msg)
	// }

	// 服务器主动下发的消息序号为0
	if msg.GetSeriNo() != 0 {
		p := r.msgCBList.Front()
		if p != nil {
			cc := p.Value.(*msgCallbackChain)
			if msg.GetSeriNo() == cc.msg.GetSeriNo() && msg.GetCmd() == cc.msg.GetCmd() {
				cc.call(r, msg)
				r.msgCBList.Remove(p)

				// 网络统计
				statistics.Message().Received(recvTime.Sub(cc.sendTime))

			} else {
				r.Panicf("onMessage seri:%d cmd:%d : do not in sequence, it should be seri:%d cmd:%d", msg.GetSeriNo(), msg.GetCmd(), cc.msg.GetSeriNo(), cc.msg.GetCmd())
			}
		} else {
			r.Panicf("onMessage seri:%d cmd:%d : no callback-chain", msg.GetSeriNo(), msg.GetCmd())
		}
	}
}

type msgHandler func(*Robot, *codecs.Message)

var msgHandlers = map[uint16]msgHandler{}

func regMsgHandler(msgId uint16, h msgHandler) {
	if _, ok := msgHandlers[msgId]; ok {
		panic(fmt.Errorf("duplcate robot msg:%d handler registration", msgId))
	}

	msgHandlers[msgId] = h
}

func dispatchRobotMsg(r *Robot, msg *codecs.Message) {
	r.Debugf("recv seri:%d cmd:%d errcode:%s", msg.GetSeriNo(), msg.GetCmd(), net.GetErrCodeStr(msg.GetErrCode()))
	h, ok := msgHandlers[msg.GetCmd()]
	if ok {
		h(r, msg)
	}
}

func onSync(r *Robot, msg *codecs.Message) {
	r.onSync(msg)
}

func regSyncHandler(msgId uint16) {
	regMsgHandler(msgId, onSync)
}

var msgID2ModuleID = map[uint16]int{}

func init() {
	module.TraverseDefines(func(id int, md *module.ModuleDefine) bool {
		if md.Name == "" || md.AssociatedMsg == nil || md.Create == nil {
			panic(fmt.Errorf("module id:%d name:%s definition invalid", id, md.Name))
		}

		for _, msgId := range md.AssociatedMsg {
			msgID2ModuleID[msgId] = id
			regSyncHandler(msgId)
		}

		return true
	})
}

func (r *Robot) onSync(msg *codecs.Message) {
	moduleID := msgID2ModuleID[msg.GetCmd()]
	if moduleID == module.Module_Unknown {
		r.Debugf("onSync msg:%d no associated moduleID", msg.GetCmd())
		return
	}

	m := r.GetModule(moduleID)
	if m == nil {
		r.Panicf("module:%s not exist", module.GetModuleName(moduleID))
	}

	r.Debugf("onSync msg:%d module:%s", msg.GetCmd(), module.GetModuleName(moduleID))
	m.OnModuleSync(r, msg.GetData())
}
