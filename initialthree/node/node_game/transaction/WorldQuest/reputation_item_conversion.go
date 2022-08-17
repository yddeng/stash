package WorldQuest

import (
	"initialthree/codec/cs"
	"initialthree/node/common/enumType"
	"initialthree/node/common/inoutput"
	"initialthree/node/common/transaction"
	"initialthree/node/node_game/module"
	"initialthree/node/node_game/module/worldQuest"
	"initialthree/node/node_game/user"
	"initialthree/node/table/excel/ConstTable/PlayerCamp"
	"initialthree/node/table/excel/DataTable/PlayerCampReputationShopItem"
	"initialthree/protocol/cmdEnum"
	cs_msg "initialthree/protocol/cs/message"
	"initialthree/zaplogger"
)

type transactionReputationItemConversion struct {
	transaction.TransactionBase
	user    *user.User
	req     *cs.Message
	resp    *cs_msg.ReputationItemConversionToC
	errCode cs_msg.ErrCode
}

func (t *transactionReputationItemConversion) Begin() {
	defer func() { t.EndTrans(t.resp, t.errCode) }()
	t.errCode = cs_msg.ErrCode_OK
	t.resp = &cs_msg.ReputationItemConversionToC{}
	reqMsg := t.req.GetData().(*cs_msg.ReputationItemConversionToS)
	zaplogger.GetSugar().Infof("%s Call ReputationItemConversionToS %v", t.user.GetUserLogName(), reqMsg)

	if reqMsg.GetCount() <= 0 {
		zaplogger.GetSugar().Debugf("%s ReputationItemConversionToS argument failed. ", t.user.GetUserLogName())
		t.errCode = cs_msg.ErrCode_Request_Argument_Err
		return
	}

	def := PlayerCampReputationShopItem.GetID(reqMsg.GetId())
	if def == nil {
		zaplogger.GetSugar().Debugf("%s ReputationItemConversionToS %d def is nil. ", t.user.GetUserLogName(), reqMsg.GetId())
		t.errCode = cs_msg.ErrCode_Config_NotExist
		return
	}

	userModule := t.user.GetSubModule(module.WorldQuest).(*worldQuest.WorldQuest)

	limit := def.ItemCapacity != -1
	if limit && userModule.ShopItemTimes(reqMsg.GetId())+reqMsg.GetCount() > def.ItemCapacity {
		zaplogger.GetSugar().Debugf("%s ReputationItemConversionToS %d limit capacity. ", t.user.GetUserLogName(), reqMsg.GetId())
		t.errCode = cs_msg.ErrCode_Item_Not_Enough
		return
	}
	constCamp := PlayerCamp.GetID(1)
	itemId := constCamp.ReputationItem[reqMsg.GetId()-1].ID
	in := []inoutput.ResDesc{{ID: itemId, Type: enumType.IOType_Item, Count: def.CostItemCount * reqMsg.GetCount()}}
	out := []inoutput.ResDesc{{ID: def.ItemID, Type: enumType.IOType_Item, Count: def.ItemCount * reqMsg.GetCount()}}

	if t.errCode = t.user.DoInputOutput(in, out, false); t.errCode != cs_msg.ErrCode_OK {
		zaplogger.GetSugar().Debugf("%s ReputationItemConversionToS error: inouput %s  ", t.user.GetUserLogName(), t.errCode.String())
		return
	}

	if limit {
		userModule.ShopItem(reqMsg.GetId(), reqMsg.GetCount())
	}

	zaplogger.GetSugar().Infof("%v ReputationItemConversionToS OK", t.user.GetUserLogName())
}

func (t *transactionReputationItemConversion) GetModuleName() string {
	return "WorldQuest"
}

func init() {
	user.RegisterTransFunc(cmdEnum.CS_ReputationItemConversion, func(user *user.User, msg *cs.Message) transaction.Transaction {
		return &transactionReputationItemConversion{
			user: user,
			req:  msg,
		}
	})
}
