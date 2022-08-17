package Attr

import (
	attr2 "initialthree/node/common/attr"
	"initialthree/node/common/enumType"
	"initialthree/node/common/inoutput"
	"initialthree/node/common/transaction"
	"initialthree/node/node_game/event"
	"initialthree/zaplogger"

	"initialthree/node/node_game/module"
	"initialthree/node/node_game/module/attr"
	"initialthree/node/node_game/user"
	"initialthree/node/table/excel/ConstTable/Player"
	"initialthree/node/table/excel/DataTable/FatigueSupply"
	"initialthree/protocol/cmdEnum"

	codecs "initialthree/codec/cs"
	cs_message "initialthree/protocol/cs/message"
)

type transactionBuyFatigue struct {
	transaction.TransactionBase
	user    *user.User
	req     *codecs.Message
	errcode cs_message.ErrCode
}

func (t *transactionBuyFatigue) Begin() {
	t.errcode = cs_message.ErrCode_OK

	userAttr := t.user.GetSubModule(module.Attr).(*attr.UserAttr)
	playerConst := Player.Get()
	maxBuyCount := playerConst.FatigueSupplyCountEveryDay
	buyCount, _ := userAttr.GetAttr(attr2.FatigueBuyCount)

	if int32(buyCount) >= maxBuyCount {
		zaplogger.GetSugar().Infof("user(%s, %d) buy fatigue count: no count", t.user.GetUserID(), t.user.GetID())
		t.errcode = cs_message.ErrCode_Attr_High
		goto over
	}

	{
		cfg := FatigueSupply.GetID(int32(buyCount) + 1)
		if cfg == nil {
			zaplogger.GetSugar().Infof("user(%s, %d) buy fatigue count: %d time config not found",
				t.user.GetUserID(), t.user.GetID(), buyCount+1)
			t.errcode = cs_message.ErrCode_Config_NotExist
			goto over
		}

		in := []inoutput.ResDesc{{Type: enumType.IOType_UsualAttribute, ID: attr2.Diamond, Count: cfg.DiamondCost}}
		out := []inoutput.ResDesc{
			{Type: enumType.IOType_UsualAttribute, ID: attr2.CurrentFatigue, Count: cfg.Supply},
			{Type: enumType.IOType_UsualAttribute, ID: attr2.FatigueBuyCount, Count: 1},
		}

		if t.errcode = t.user.DoInputOutput(in, out, true); t.errcode != cs_message.ErrCode_OK {
			zaplogger.GetSugar().Errorf("%s buy fatigue count: %d err %s", t.user.GetUserLogName(), buyCount+1, t.errcode.String())
			goto over
		}

		// 事件触发
		t.user.EmitEvent(event.EventExchangeFatigue, int32(1))
	}

	zaplogger.GetSugar().Infof("user(%s, %d) buy fatigue count %d time successfully", t.user.GetUserID(), t.user.GetID(), buyCount+1)

over:
	t.EndTrans(&cs_message.BuyFatigueToC{}, t.errcode)
}

func (t *transactionBuyFatigue) GetModuleName() string {
	return "Attr"
}

func init() {
	user.RegisterTransFunc(cmdEnum.CS_BuyFatigue, func(user *user.User, msg *codecs.Message) transaction.Transaction {
		return &transactionBuyFatigue{user: user, req: msg}
	})
}
