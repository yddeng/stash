package user

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	codecs "initialthree/codec/cs"
	"initialthree/node/common/pipeline"
	"initialthree/node/common/transaction"
	"initialthree/node/node_game/module"
	"initialthree/node/node_game/module/assets"
	cs_message "initialthree/protocol/cs/message"
	"time"
)

type TransCreateFunc func(user *User, msg *codecs.Message) transaction.Transaction
type CheckModuleOpenFunc func(user *User, msg *codecs.Message) bool

var checkModuleOpen = map[string]CheckModuleOpenFunc{}

func RegisterCheckModuleOpen(moduleName string, fn CheckModuleOpenFunc) {
	_, ok := checkModuleOpen[moduleName]
	if ok {
		panic(fmt.Sprintf("mudule Name %s is aleary register", moduleName))
		return
	}

	checkModuleOpen[moduleName] = fn
}

func (this *User) callTransEnd(trans transaction.Transaction, req *codecs.Message, resp proto.Message, errCode cs_message.ErrCode, usedDuration time.Duration) {
	if trans.GetModuleName() != "CreateRole" {
		this.FlushAllDirtyToClient()
	}

	// 回消息
	if errCode == cs_message.ErrCode_OK {
		if resp != nil {
			this.Reply(req.GetSeriNo(), resp)
		}
	} else {
		this.ReplyErr(req.GetSeriNo(), req.GetCmd(), errCode)
	}

	// 新手教学
	if errCode == cs_message.ErrCode_OK && req.GetTeachID() != 0 {
		userAsset := this.GetSubModule(module.Assets).(*assets.UserAssets)
		userAsset.SetAsset(int32(cs_message.AssetType_Teaching), int32(req.GetTeachID()), 1)
	}

	// 拉取请求
	this.recv()
}

type TransLine struct {
	creator  TransCreateFunc
	pipeline *pipeline.Pipeline
}

var transLine = map[uint16]*TransLine{}

// Deprecated: Use RegisterTransStep
func RegisterTransFunc(cmd uint16, creator TransCreateFunc) {
	RegisterTransStep(cmd, creator, StepIsMessageDisable, StepIsModuleDisable, StepCheckUserState, StepCheckModuleOpen)
}

// 步骤之间应相互隔离，没有前后间依赖
func RegisterTransStep(cmd uint16, creator TransCreateFunc, steps ...pipeline.StepFunc) {
	if nil == creator {
		return
	}
	_, ok := transLine[cmd]
	if ok {
		panic(fmt.Sprintf("cmd %d is aleary register", cmd))
		return
	}

	line := &TransLine{
		creator:  creator,
		pipeline: pipeline.NewPipeline(),
	}

	line.pipeline.AddStep(steps...)

	transLine[cmd] = line
}

func TransDispatch(u *User, msg *codecs.Message) {
	elem := &PipelineElem{Msg: msg, User: u}
	cmd := msg.GetCmd()

	tl, ok := transLine[cmd]
	if !ok {
		u.recv()
		return
	}

	trans := tl.creator(elem.User, elem.Msg)
	elem.Trans = trans

	// trans 前置检查
	out, err := tl.pipeline.Run(elem)
	if err != nil {
		elem = out.(*PipelineElem)
		u.ReplyErr(msg.GetSeriNo(), msg.GetCmd(), elem.ErrCode)
		u.recv()
		return
	}

	if nil != u.transMgr.PushTrans(trans, msg, time.Second*transaction.TransTimeoutSec) {
		u.recv()
	}
}

/************************ step start ******************************** */
type PipelineElem struct {
	Msg     *codecs.Message
	User    *User
	Trans   transaction.Transaction
	ErrCode cs_message.ErrCode
}

func StepIsMessageDisable(in interface{}) (interface{}, error) {
	elem := in.(*PipelineElem)
	//cmd := elem.Msg.GetCmd()
	//if functionSwitch.IsMessageDisable(cmd) {
	//	//如果消息被熔断，直接响应错误
	//	elem.ErrCode = cs_message.ErrCode_FUNCTION_DISABLE
	//	return elem, fmt.Errorf("transation message %d is disable", cmd)
	//}
	return elem, nil
}

func StepIsModuleDisable(in interface{}) (interface{}, error) {
	elem := in.(*PipelineElem)
	//gameModule := elem.Trans.GetGameModule()
	//if functionSwitch.IsModuleDisable(elem.Trans.GetGameModule()) {
	//	//如果消息所属模块被熔断，直接响应错误
	//	elem.ErrCode = cs_message.ErrCode_FUNCTION_DISABLE
	//	return elem, fmt.Errorf("transation module %d is disable", gameModule)
	//}
	return elem, nil
}

// 检查账号是否初始化
func StepCheckUserState(in interface{}) (interface{}, error) {
	elem := in.(*PipelineElem)
	moduleName := elem.Trans.GetModuleName()
	u := elem.User
	if u.GetID() == 0 {
		elem.ErrCode = cs_message.ErrCode_User_NotExist
		return elem, fmt.Errorf("syatem %s is unlock, need create role. ", moduleName)
	}
	return elem, nil
}

// 检查账号 模块是否达到开启条件
func StepCheckModuleOpen(in interface{}) (interface{}, error) {
	elem := in.(*PipelineElem)
	moduleName := elem.Trans.GetModuleName()

	moduleOpen := true // 默认开启
	fn, ok := checkModuleOpen[moduleName]
	if ok {
		moduleOpen = fn(elem.User, elem.Msg)
	}

	if !moduleOpen {
		elem.ErrCode = cs_message.ErrCode_System_Unlock
		return elem, fmt.Errorf("syatem %s is unlock", moduleName)
	}
	return elem, nil
}

/************************ step end ******************************** */
