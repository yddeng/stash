package Attr

import (
	codecs "initialthree/codec/cs"
	"initialthree/node/common/attr"
	"initialthree/node/common/enumType"
	"initialthree/node/common/inoutput"
	"initialthree/node/common/transaction"
	"initialthree/node/node_game/user"
	"initialthree/node/table/excel/DataTable/GoldSupply"
	"initialthree/protocol/cmdEnum"
	cs_message "initialthree/protocol/cs/message"
	"initialthree/zaplogger"
)

type transactionBuyGold struct {
	transaction.TransactionBase
	user *user.User
	req  *codecs.Message
}

func (t *transactionBuyGold) Begin() {
	goldBuyCount := t.user.GetAttr(attr.GoldBuyCount)
	maxGoldBuyCount := attr.GetAttrInfo(attr.GoldBuyCount).Max

	if goldBuyCount >= maxGoldBuyCount {
		zaplogger.GetSugar().Infof("user(%s, %d) no gold buy count", t.user.GetUserID(), t.user.GetID())
		t.EndTrans(&cs_message.BuyGoldToC{}, cs_message.ErrCode_Attr_High)
		return
	}

	goldBuyCount += 1
	goldSupply := GoldSupply.GetID(int32(goldBuyCount))
	if goldSupply == nil {
		zaplogger.GetSugar().Errorf("user(%s, %d) buy No.%d gold, config not found", t.user.GetUserID(), t.user.GetID(), goldBuyCount)
		t.EndTrans(&cs_message.BuyGoldToC{}, cs_message.ErrCode_Config_NotExist)
		return
	}

	gold := t.user.GetAttr(attr.Gold)
	if int64(goldSupply.Supply)+gold > attr.GetAttrInfo(attr.Gold).Max {
		zaplogger.GetSugar().Infof("user(%s, %d) buy No.%d gold, gold overflow", t.user.GetUserID(), t.user.GetID(), goldBuyCount)
		t.EndTrans(&cs_message.BuyGoldToC{}, cs_message.ErrCode_Attr_High)
		return
	}

	input := []inoutput.ResDesc{{Type: enumType.IOType_UsualAttribute, ID: attr.Diamond, Count: goldSupply.DiamondCost}}
	output := []inoutput.ResDesc{
		{Type: enumType.IOType_UsualAttribute, ID: attr.GoldBuyCount, Count: 1},
		{Type: enumType.IOType_UsualAttribute, ID: attr.Gold, Count: goldSupply.Supply},
	}
	if errCode := t.user.DoInputOutput(input, output, false); errCode != cs_message.ErrCode_OK {
		zaplogger.GetSugar().Debugf("%s buy gold error: inouput %s  ", t.user.GetUserLogName(), errCode.String())
		t.EndTrans(&cs_message.BuyGoldToC{}, errCode)
		return
	}

	zaplogger.GetSugar().Debugf("user(%s, %d) buy No.%d gold.", t.user.GetUserID(), t.user.GetID(), goldBuyCount)
	t.EndTrans(&cs_message.BuyGoldToC{}, cs_message.ErrCode_OK)
}

func (t *transactionBuyGold) GetModuleName() string {
	return "Attr"
}

func init() {
	user.RegisterTransFunc(cmdEnum.CS_BuyGold, func(user *user.User, msg *codecs.Message) transaction.Transaction {
		return &transactionBuyGold{
			user: user,
			req:  msg,
		}
	})
}
