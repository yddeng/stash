package RewardQuest

import (
	codecs "initialthree/codec/cs"
	"initialthree/node/common/transaction"
	"initialthree/zaplogger"

	"initialthree/node/node_game/module"
	"initialthree/node/node_game/module/rewardQuest"
	"initialthree/node/node_game/user"
	RewardQuest2 "initialthree/node/table/excel/ConstTable/RewardQuest"
	"initialthree/node/table/excel/DataTable/RewardQuest"
	"initialthree/node/table/excel/DataTable/RewardQuestPosition"
	"initialthree/protocol/cmdEnum"
	cs_message "initialthree/protocol/cs/message"
)

type transactionRewardQuestRefresh struct {
	transaction.TransactionBase
	user    *user.User
	req     *codecs.Message
	resp    *cs_message.RewardQuestRefreshToC
	errcode cs_message.ErrCode
}

func (this *transactionRewardQuestRefresh) GetModuleName() string {
	return "RewardQuest"
}

func (this *transactionRewardQuestRefresh) Begin() {
	defer func() { this.EndTrans(this.resp, this.errcode) }()
	msg := this.req.GetData().(*cs_message.RewardQuestRefreshToS)

	zaplogger.GetSugar().Infof("%s %d Call RewardQuestRefreshToS %v", this.user.GetUserID(), this.user.GetID(), msg)

	userQuest := this.user.GetSubModule(module.RewardQuest).(*rewardQuest.RewardQuest)

	q := userQuest.GetRewardQuest(msg.GetQuestID())
	if q == nil {
		zaplogger.GetSugar().Debugf("%s RewardQuestRefreshToS failed, rewardQuest %d is not exist", this.user.GetUserLogName(), msg.GetQuestID())
		this.errcode = cs_message.ErrCode_Quest_NotExist
		return
	}

	if q.GetState() != cs_message.QuestState_Acceptable {
		zaplogger.GetSugar().Debugf("%s RewardQuestRefreshToS failed, rewardQuest %d state failed", this.user.GetUserLogName(), msg.GetQuestID())
		this.errcode = cs_message.ErrCode_Quest_StateErr
		return
	}

	def := RewardQuest.GetID(msg.GetQuestID())
	if def == nil {
		zaplogger.GetSugar().Debugf("%s RewardQuestRefreshToS failed, rewardQuest config %d is not exist", this.user.GetUserLogName(), msg.GetQuestID())
		this.errcode = cs_message.ErrCode_Config_NotExist
		return
	}

	constDef := RewardQuest2.GetID(1)
	base := userQuest.GetBase()
	if base.RefreshTimes >= constDef.RefreshTimes {
		zaplogger.GetSugar().Debugf("%s RewardQuestRefreshToS failed, refresh times is not enough", this.user.GetUserLogName())
		this.errcode = cs_message.ErrCode_Attr_Low
		return
	}

	exist := getExist(userQuest, q.GetQuestID())
	pos := RewardQuestPosition.RandPositions(1, exist)
	posQuality := RewardQuest2.RandQuality(constDef.SSCount-base.SSCount, constDef.SCount-base.SCount, pos)
	userQuest.Replace(map[int32]struct{}{q.GetQuestID(): {}}, posQuality, true)

	zaplogger.GetSugar().Infof("%s RewardQuestRefreshToS ok", this.user.GetUserLogName())
}

func init() {
	user.RegisterTransFunc(cmdEnum.CS_RewardQuestRefresh, func(user *user.User, msg *codecs.Message) transaction.Transaction {
		return &transactionRewardQuestRefresh{
			user:    user,
			req:     msg,
			resp:    &cs_message.RewardQuestRefreshToC{},
			errcode: cs_message.ErrCode_OK,
		}
	})
}
