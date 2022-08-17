package Backpack

import (
	codecs "initialthree/codec/cs"
	"initialthree/node/common/attr"
	"initialthree/node/common/enumType"
	"initialthree/node/common/inoutput"
	"initialthree/node/common/transaction"
	"initialthree/node/node_game/event"
	"initialthree/node/table/excel/DataTable/Item"
	"initialthree/zaplogger"

	"initialthree/node/node_game/module"
	backpack2 "initialthree/node/node_game/module/backpack"
	"initialthree/node/node_game/user"
	IOTable "initialthree/node/table/excel/DataTable/InputOutput"
	"initialthree/protocol/cmdEnum"
	cs_message "initialthree/protocol/cs/message"
	"time"
)

type transactionBackpackUse struct {
	transaction.TransactionBase
	user    *user.User
	req     *codecs.Message
	errcode cs_message.ErrCode
}

func (t *transactionBackpackUse) Begin() {
	defer func() { t.EndTrans(&cs_message.BackpackUseToC{}, t.errcode) }()

	t.errcode = cs_message.ErrCode_OK
	reqMsg := t.req.GetData().(*cs_message.BackpackUseToS)
	backpack := t.user.GetSubModule(module.Backpack).(*backpack2.Backpack)

	it := backpack.GetItem(reqMsg.GetId())
	if it == nil {
		zaplogger.GetSugar().Infof("user(%s, %d) backpack use entity(%d, %d): not found",
			t.user.GetUserID(), t.user.GetID(), reqMsg.GetId(), reqMsg.GetCount())
		t.errcode = cs_message.ErrCode_Backpack_EntityError
		return
	}

	itCfg := Item.GetID(it.TID)
	zaplogger.GetSugar().Debugf("user(%s, %d) backpack use entity(%d, %d) cfgID %d",
		t.user.GetUserID(), t.user.GetID(), reqMsg.GetId(), reqMsg.GetCount(), itCfg.ID)

	if it.IsExpired(time.Now()) {
		zaplogger.GetSugar().Infof("user(%s, %d) backpack use entity(%d, %d): entity expired",
			t.user.GetUserID(), t.user.GetID(), reqMsg.GetId(), reqMsg.GetCount())
		t.errcode = cs_message.ErrCode_Backpack_EntityExpired
		return
	}

	if reqMsg.GetCount() <= 0 || reqMsg.GetCount() > it.Count {
		zaplogger.GetSugar().Infof("user(%s, %d) backpack use entity(%d, %d): count error",
			t.user.GetUserID(), t.user.GetID(), reqMsg.GetId(), reqMsg.GetCount())
		t.errcode = cs_message.ErrCode_Request_Argument_Err
		return
	}

	if !(itCfg.AllowUse || itCfg.AllowBatchUse) {
		zaplogger.GetSugar().Infof("user(%s, %d) backpack use entity(%d, %d): not allowed to use",
			t.user.GetUserID(), t.user.GetID(), reqMsg.GetId(), reqMsg.GetCount())
		t.errcode = cs_message.ErrCode_Backpack_NotAllowedToUse
		return
	}

	cfg := IOTable.GetID(itCfg.UseRuleID)

	if nil == cfg {
		t.errcode = cs_message.ErrCode_ERROR
		return
	}

	input := []inoutput.ResDesc{}
	output := []inoutput.ResDesc{}

	//for _, v := range cfg.Input {
	//	input = append(input, inoutput.ResDesc{ID: v.ID, Type: int(v.Type), Count: v.Count * reqMsg.GetCount()})
	//}
	fatigueRegen := false
	// 设置消耗
	input = append(input, inoutput.ResDesc{ID: it.TID, Type: enumType.IOType_Item, Count: reqMsg.GetCount(), InsID: reqMsg.GetId()})

	for _, v := range cfg.Output {
		output = append(output, inoutput.ResDesc{ID: v.ID, Type: int(v.Type), Count: v.Count * reqMsg.GetCount()})

		if !fatigueRegen && v.Type == enumType.IOType_UsualAttribute && v.ID == attr.CurrentFatigue {
			fatigueRegen = true
		}
	}

	if t.errcode = t.user.DoInputOutput(input, output, false); t.errcode != cs_message.ErrCode_OK {
		zaplogger.GetSugar().Debugf("%s BackpackUse error: inouput %s  ", t.user.GetUserLogName(), t.errcode.String())
		return
	}

	// 事件触发
	t.user.EmitEvent(event.EventUseItem, map[int32]int32{it.TID: reqMsg.GetCount()})

	if fatigueRegen {
		t.user.EmitEvent(event.EventFatigueRegen, reqMsg.GetCount())
	}
}

func (t *transactionBackpackUse) GetModuleName() string {
	return "Backpack"
}

func init() {
	user.RegisterTransFunc(cmdEnum.CS_BackpackUse, func(user *user.User, msg *codecs.Message) transaction.Transaction {
		return &transactionBackpackUse{
			user: user,
			req:  msg,
		}
	})
}
