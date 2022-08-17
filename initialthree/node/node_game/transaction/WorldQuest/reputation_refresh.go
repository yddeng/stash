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
	"initialthree/protocol/cmdEnum"
	cs_msg "initialthree/protocol/cs/message"
	"initialthree/zaplogger"
)

type transactionReputationRefresh struct {
	transaction.TransactionBase
	user    *user.User
	req     *cs.Message
	resp    *cs_msg.ReputationRefreshToC
	errCode cs_msg.ErrCode
}

func (t *transactionReputationRefresh) Begin() {
	defer func() { t.EndTrans(t.resp, t.errCode) }()
	t.errCode = cs_msg.ErrCode_OK
	t.resp = &cs_msg.ReputationRefreshToC{}
	reqMsg := t.req.GetData().(*cs_msg.ReputationRefreshToS)
	zaplogger.GetSugar().Infof("%s Call ReputationRefreshToS %v", t.user.GetUserLogName(), reqMsg)

	def := PlayerCamp.GetID(1)
	userModule := t.user.GetSubModule(module.WorldQuest).(*worldQuest.WorldQuest)

	nowTimes := userModule.CampRefreshTimes(reqMsg.GetId())
	if nowTimes >= int32(len(def.ReputationRefreshArray)) {
		zaplogger.GetSugar().Infof("%s ReputationRefreshToS failed:RefreshTimes is not enough", t.user.GetUserLogName())
		t.errCode = cs_msg.ErrCode_Shop_RefreshTimesNotEnough
		return
	}

	if cost := def.ReputationRefreshArray[nowTimes].Cost; cost > 0 {
		itemID := def.ReputationItem[reqMsg.GetId()-1].ID
		in := []inoutput.ResDesc{{ID: itemID, Type: enumType.IOType_Item, Count: cost}}

		if t.errCode = t.user.DoInputOutput(in, nil, false); t.errCode != cs_msg.ErrCode_OK {
			zaplogger.GetSugar().Debugf("%s ReputationRefreshToS error: inouput %s  ", t.user.GetUserLogName(), t.errCode.String())
			return
		}

	}

	userModule.CampRefresh(reqMsg.GetId())

	zaplogger.GetSugar().Infof("%v ReputationRefreshToS OK", t.user.GetUserLogName())
}

func (t *transactionReputationRefresh) GetModuleName() string {
	return "WorldQuest"
}

func init() {
	user.RegisterTransFunc(cmdEnum.CS_ReputationRefresh, func(user *user.User, msg *cs.Message) transaction.Transaction {
		return &transactionReputationRefresh{
			user: user,
			req:  msg,
		}
	})
}
