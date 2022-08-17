package WorldQuest

import (
	"initialthree/codec/cs"
	"initialthree/node/common/attr"
	"initialthree/node/common/enumType"
	"initialthree/node/common/inoutput"
	"initialthree/node/common/transaction"
	"initialthree/node/node_game/module"
	"initialthree/node/node_game/module/worldQuest"
	"initialthree/node/node_game/user"
	"initialthree/node/table/excel/ConstTable/Quest"
	"initialthree/protocol/cmdEnum"
	cs_msg "initialthree/protocol/cs/message"
	"initialthree/zaplogger"
)

type transactionWorldQuestRefresh struct {
	transaction.TransactionBase
	user    *user.User
	req     *cs.Message
	errCode cs_msg.ErrCode
}

func (t *transactionWorldQuestRefresh) Begin() {
	defer func() { t.EndTrans(&cs_msg.WorldQuestRefreshToC{}, t.errCode) }()

	t.errCode = cs_msg.ErrCode_OK
	reqMsg := t.req.GetData().(*cs_msg.WorldQuestRefreshToS)
	zaplogger.GetSugar().Infof("%s Call WorldQuestRefresh %v", t.user.GetUserLogName(), reqMsg)

	userModule := t.user.GetSubModule(module.WorldQuest).(*worldQuest.WorldQuest)
	def := Quest.GetID(1)
	times := userModule.GetWQRefreshTimes()

	if times >= def.WorldQuestRefreshFreeMaxCount {
		in := make([]inoutput.ResDesc, 0, 1)
		if def.WorldQuestRefreshItemID != 0 && t.user.GetItemCountByTID(def.WorldQuestRefreshItemID) > 0 {
			in = append(in, inoutput.ResDesc{Type: enumType.IOType_Item, ID: def.WorldQuestRefreshItemID, Count: 1})
		}
		if len(in) == 0 && def.WorldQuestRefreshDiamondCount != 0 {
			if t.user.GetAttr(attr.Diamond) < int64(def.WorldQuestRefreshDiamondCount) {
				t.errCode = cs_msg.ErrCode_Resource_NotEnough
				return
			}
			in = append(in, inoutput.ResDesc{Type: enumType.IOType_UsualAttribute, ID: attr.Diamond, Count: def.WorldQuestRefreshDiamondCount})
		}
		if t.errCode = t.user.DoInputOutput(in, nil, false); t.errCode != cs_msg.ErrCode_OK {
			zaplogger.GetSugar().Debugf("%s WorldQuestRefresh error: inouput %s  ", t.user.GetUserLogName(), t.errCode.String())
			return
		}

	}

	userModule.WQRefresh()
	zaplogger.GetSugar().Infof("%v WorldQuestRefresh OK", t.user.GetUserLogName())
}

func (t *transactionWorldQuestRefresh) GetModuleName() string {
	return "WorldQuest"
}

func init() {
	user.RegisterTransFunc(cmdEnum.CS_WorldQuestRefresh, func(user *user.User, msg *cs.Message) transaction.Transaction {
		return &transactionWorldQuestRefresh{
			user: user,
			req:  msg,
		}
	})
}
