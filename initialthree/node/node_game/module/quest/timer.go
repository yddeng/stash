package quest

import (
	"initialthree/node/common/droppool"
	"initialthree/node/common/enumType"
	"initialthree/node/node_game/module"
	TableQuest "initialthree/node/table/excel/ConstTable/Quest"
	"initialthree/node/table/excel/DataTable/BigSecretCompetition"
	"initialthree/node/table/excel/DataTable/Mail"
	"initialthree/protocol/cs/message"
	"initialthree/zaplogger"
	"time"
)

func (this *UserQuest) tryClock(name string, nowUnix int64, fn func(nowUnix int64)) {
	if timestamp, ok := this.timedata[name]; ok {
		if nowUnix >= timestamp {
			fn(nowUnix)
		}
	} else {
		fn(nowUnix)
	}
}

func (this *UserQuest) clockTimer() {
	now := time.Now().Unix()

	this.tryClock(module.WeeklyTimeName, now, this.weeklyClock)
	this.tryClock(module.DailyTimeName, now, this.dailyClock)
}

// 日更新
func (this *UserQuest) dailyClock(now int64) {
	this.clear(enumType.QuestType_Daily)
	for _, id := range TableQuest.GetDailyQuest(int(time.Now().Weekday())) {
		this.AddQuest(id)
	}
	this.clear(enumType.QuestType_DailyReward)
	for _, id := range TableQuest.GetDailyRewardQuest() {
		this.AddQuest(id)
	}

	this.FlushAllToClient()

	this.timedata[module.DailyTimeName] = module.CalDailyTime().Unix()
	this.SetDirty(timeField)
}

// 周更新
func (this *UserQuest) weeklyClock(now int64) {
	this.clear(enumType.QuestType_Weekly)
	for _, id := range TableQuest.GetWeekdayQuest() {
		this.AddQuest(id)
	}
	this.FlushAllToClient()

	this.timedata[module.WeeklyTimeName] = module.CalWeeklyTime().Unix()
	this.SetDirty(timeField)

}

// 清空某类型的任务
func (this *UserQuest) clear(tt int32) {

	// 若奖励未领取，通过邮件发送
	if tt == enumType.QuestType_NewbieGift ||
		tt == enumType.QuestType_DailyReward ||
		tt == enumType.QuestType_Weekly {

		awardPoolIDs := make([]int32, 0, len(this.quests[tt]))
		for _, q := range this.quests[tt] {
			cfg := q.GetConfig()
			if q.State == message.QuestState_Finished && cfg.UnreceivedSendMail && cfg.Reward != 0 {
				awardPoolIDs = append(awardPoolIDs, cfg.Reward)
			}
		}

		if len(awardPoolIDs) > 0 {
			award := droppool.DropAward(awardPoolIDs...)
			if !award.IsZero() {
				switch tt {
				case enumType.QuestType_DailyReward:
					m := Mail.DailyQuestMail(time.Now(), award.ToMessageAward())
					this.userI.SendMail([]*message.Mail{m})
				case enumType.QuestType_Weekly:
					m := Mail.WeeklyQuestMail(time.Now(), award.ToMessageAward())
					this.userI.SendMail([]*message.Mail{m})
				case enumType.QuestType_NewbieGift:
					m := Mail.NewbieGiftMail(time.Now(), award.ToMessageAward())
					this.userI.SendMail([]*message.Mail{m})
				default:

				}
			}
		}
	}

	// 清理事件
	for _, q := range this.quests[tt] {
		if q.State == message.QuestState_Running {
			for _, c := range q.Conditions {
				if c.condEvent != nil {
					c.condEvent.UnTrigger()
				}
			}
		}
	}

	this.quests[tt] = map[int32]*Quest{}
	this.SetDirty(tt)
}

func (this *UserQuest) EventResetBigSecret(id int32) {
	zaplogger.GetSugar().Infof("event reset bigSecret %d ", id)
	this.clear(enumType.QuestType_BigSecret)

	def := BigSecretCompetition.GetID(id)
	if def != nil {
		for _, v := range def.QuestArray {
			this.AddQuest(v.QuestID)
		}
	}
}
