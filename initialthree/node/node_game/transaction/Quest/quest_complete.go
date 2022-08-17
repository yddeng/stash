package Quest

import (
	codecs "initialthree/codec/cs"
	"initialthree/node/common/attr"
	"initialthree/node/common/droppool"
	"initialthree/node/common/enumType"
	"initialthree/node/common/transaction"
	"initialthree/node/node_game/event"
	attr2 "initialthree/node/node_game/module/attr"
	"initialthree/node/table/excel/ConstTable/Global"
	"initialthree/node/table/excel/DataTable/NewbieGift"
	"initialthree/node/table/excel/DataTable/Quest"
	"initialthree/zaplogger"
	"time"

	"initialthree/node/node_game/module"
	"initialthree/node/node_game/module/quest"
	"initialthree/node/node_game/user"
	"initialthree/protocol/cmdEnum"
	cs_message "initialthree/protocol/cs/message"
)

type transactionQuestComplete struct {
	transaction.TransactionBase
	user       *user.User
	req        *codecs.Message
	resp       *cs_message.QuestCompleteToC
	errcode    cs_message.ErrCode
	days       map[int]bool
	newbieGift bool
}

func (this *transactionQuestComplete) GetModuleName() string {
	return "Quest"
}

func (this *transactionQuestComplete) Begin() {
	defer func() { this.EndTrans(this.resp, this.errcode) }()
	this.errcode = cs_message.ErrCode_OK
	this.resp = &cs_message.QuestCompleteToC{}
	msg := this.req.GetData().(*cs_message.QuestCompleteToS)

	zaplogger.GetSugar().Infof("%s  Call QuestCompleteToS %v", this.user.GetUserLogName(), msg)

	userQuest := this.user.GetSubModule(module.Quest).(*quest.UserQuest)

	for _, id := range msg.GetQuestID() {
		q := userQuest.GetQuest(id)

		if q == nil || q.State != cs_message.QuestState_Finished {
			this.errcode = cs_message.ErrCode_Quest_StateErr
			return
		}

		if q.GetConfig().TypeEnum == enumType.QuestType_NewbieGift {
			if !this.newbieGift {
				// 7日任务能完成的时间
				nowUnix := time.Now().Unix()
				attrModule := this.user.GetSubModule(module.Attr).(*attr2.UserAttr)
				startTime, _ := attrModule.GetAttr(attr.NewbieGiftStartTime)
				endTime, _ := attrModule.GetAttr(attr.NewbieGiftEndTime)
				this.days = make(map[int]bool, 7)
				if startTime != 0 && nowUnix > startTime && nowUnix < endTime {
					startT := time.Unix(startTime, 0)
					rt := Global.Get().GetDailyRefreshTime()
					refreshTime := time.Date(startT.Year(), startT.Month(), startT.Day(), int(rt.Hour), int(rt.Minute), 0, 0, time.Local) // 第一次更新时间
					if int32(startT.Hour()) > rt.Hour || (int32(startT.Hour()) == rt.Hour && int32(startT.Minute()) > rt.Minute) {
						refreshTime = refreshTime.AddDate(0, 0, 1)
					}
					this.days[1] = true
					for day := 2; refreshTime.Unix() < nowUnix && day <= 7; day, refreshTime = day+1, refreshTime.AddDate(0, 0, 1) {
						this.days[day] = true
					}
				}
				this.newbieGift = true
			}

			if len(this.days) == 0 || !this.checkNewbieQuest(q.GetConfig()) {
				this.errcode = cs_message.ErrCode_Quest_StateErr
				return
			}

		}

		userQuest.End(q)
		this.user.EmitEvent(event.EventQuestComplete, id)

		cfg := q.GetConfig()
		if cfg.Reward != 0 {
			award := droppool.DropWithID(cfg.Reward)
			_ = this.user.ApplyDropAward(award)
		}
	}

}

func (this *transactionQuestComplete) checkNewbieQuest(qcfg *Quest.Quest) bool {
	for day := range this.days {
		def := NewbieGift.GetID(int32(day))
		for _, v := range def.QuestIDListArray {
			if v.QuestID != qcfg.ID {
				return true
			}
		}
		if def.GroupRewardQuest == qcfg.ID {
			return true
		}
	}
	return false
}

func init() {
	user.RegisterTransFunc(cmdEnum.CS_QuestComplete, func(user *user.User, msg *codecs.Message) transaction.Transaction {
		return &transactionQuestComplete{
			user: user,
			req:  msg,
		}
	})
}
