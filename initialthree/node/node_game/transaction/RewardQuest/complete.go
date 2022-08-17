package RewardQuest

import (
	codecs "initialthree/codec/cs"
	"initialthree/node/common/droppool"
	"initialthree/node/common/timeDisposal"
	"initialthree/node/common/transaction"
	"initialthree/zaplogger"

	"initialthree/node/node_game/module"
	"initialthree/node/node_game/module/rewardQuest"
	"initialthree/node/node_game/user"
	"initialthree/node/table/excel/ConstTable/Global"
	RewardQuest2 "initialthree/node/table/excel/ConstTable/RewardQuest"
	"initialthree/node/table/excel/DataTable/RewardQuest"
	"initialthree/node/table/excel/DataTable/RewardQuestPosition"
	"initialthree/protocol/cmdEnum"
	cs_message "initialthree/protocol/cs/message"
	"time"
)

type transactionRewardQuestComplete struct {
	transaction.TransactionBase
	user    *user.User
	req     *codecs.Message
	resp    *cs_message.RewardQuestCompleteToC
	errcode cs_message.ErrCode
}

func (this *transactionRewardQuestComplete) GetModuleName() string {
	return "RewardQuest"
}

func (this *transactionRewardQuestComplete) Begin() {

	defer func() { this.EndTrans(this.resp, this.errcode) }()

	msg := this.req.GetData().(*cs_message.RewardQuestCompleteToS)

	zaplogger.GetSugar().Infof("%s %d Call RewardQuestCompleteToS %v", this.user.GetUserID(), this.user.GetID(), msg)

	userQuest := this.user.GetSubModule(module.RewardQuest).(*rewardQuest.RewardQuest)

	q := userQuest.GetRewardQuest(msg.GetQuestID())
	if q == nil {
		zaplogger.GetSugar().Debugf("%s RewardQuestCompleteToS failed, rewardQuest %d is not exist", this.user.GetUserLogName(), msg.GetQuestID())
		this.errcode = cs_message.ErrCode_Quest_NotExist
		return
	}

	if q.GetState() != cs_message.QuestState_Running {
		zaplogger.GetSugar().Debugf("%s RewardQuestCompleteToS failed, rewardQuest %d state failed", this.user.GetUserLogName(), msg.GetQuestID())
		this.errcode = cs_message.ErrCode_Quest_StateErr
		return
	}

	def := RewardQuest.GetID(msg.GetQuestID())
	if def == nil {
		zaplogger.GetSugar().Debugf("%s RewardQuestCompleteToS failed, rewardQuest config %d is not exist", this.user.GetUserLogName(), msg.GetQuestID())
		this.errcode = cs_message.ErrCode_Config_NotExist
		return
	}

	now := timeDisposal.NowUnix()
	endTimeUnix := q.GetAcceptTimestamp() + int64(time.Duration(def.ExecutionTime)*time.Hour)
	if now < endTimeUnix {
		zaplogger.GetSugar().Debugf("%s RewardQuestCompleteToS failed, rewardQuest %d ExecutionTime not finish", this.user.GetUserLogName(), msg.GetQuestID())
		this.errcode = cs_message.ErrCode_Quest_StateErr
		return
	}

	userQuest.Complete(q)

	// 发放奖励
	this.user.ApplyDropAward(droppool.DropAward(def.DroppoolID))

	// 如果不是今天的任务，刷新一条。 如果接取时间在今天的刷新时间之前则为以前的任务，需刷新
	rt := Global.Get().GetDailyRefreshTime()
	todayTime := timeDisposal.TodayTime(rt.Hour, rt.Minute, 0).Unix()
	if now > todayTime && q.GetAcceptTimestamp() < todayTime {
		// 现在时间大于今天的更新时间，任务的接取时间在今天更新时间之前
		exist := getExist(userQuest, q.GetQuestID())

		constDef := RewardQuest2.GetID(1)
		base := userQuest.GetBase()
		pos := RewardQuestPosition.RandPositions(1, exist)
		posQuality := RewardQuest2.RandQuality(constDef.SSCount-base.SSCount, constDef.SCount-base.SCount, pos)
		userQuest.Replace(map[int32]struct{}{q.GetQuestID(): {}}, posQuality, false)
	}

	zaplogger.GetSugar().Infof("%s RewardQuestCompleteToS ok", this.user.GetUserLogName())
}

func init() {
	user.RegisterTransFunc(cmdEnum.CS_RewardQuestComplete, func(user *user.User, msg *codecs.Message) transaction.Transaction {
		return &transactionRewardQuestComplete{
			user:    user,
			req:     msg,
			resp:    &cs_message.RewardQuestCompleteToC{},
			errcode: cs_message.ErrCode_OK,
		}
	})
}
