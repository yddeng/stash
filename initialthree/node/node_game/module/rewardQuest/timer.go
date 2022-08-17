package rewardQuest

import (
	"initialthree/node/common/timeDisposal"
	"initialthree/node/node_game/module"
	"initialthree/node/table/excel/ConstTable/Global"
	"initialthree/node/table/excel/DataTable/RewardQuestPosition"
	"initialthree/protocol/cs/message"
	"time"
)

func (this *RewardQuest) tryClock(name string, nowUnix int64, fn func(nowUnix int64)) {
	if timestamp, ok := this.timedata[name]; ok {
		if nowUnix >= timestamp {
			fn(nowUnix)
		}
	} else {
		fn(nowUnix)
	}
}

func (this *RewardQuest) clockTimer() {
	now := time.Now().Unix()

	this.tryClock(module.WeeklyTimeName, now, this.weeklyClock)
	this.tryClock(module.DailyTimeName, now, this.dailyClock)
}

// 日更新
func (this *RewardQuest) dailyClock(now int64) {
	this.base.RefreshTimes = 0
	this.SetDirty(baseField)

	// 只替换 未开始，已经完成
	// 任务刷新
	exist := map[int32]struct{}{}
	removed := map[int32]struct{}{}
	for id, q := range this.data {
		if q.GetState() == message.QuestState_Running {
			pos := id / 10
			exist[pos] = struct{}{}
		} else {
			removed[id] = struct{}{}
		}
	}

	constDef := RewardQuest2.GetID(1)
	pos := RewardQuestPosition.RandPositions(int(constDef.DailyCount)-len(exist), exist)
	ss, s := constDef.SSCount-this.base.SSCount, constDef.SCount-this.base.SCount
	posQuality := RewardQuest2.RandQuality(ss, s, pos)
	this.Replace(removed, posQuality, false)

	this.FlushAllToClient()

	this.timedata[module.DailyTimeName] = module.CalDailyTime().Unix()
	this.SetDirty(timeField)

}

// 周更新
func (this *RewardQuest) weeklyClock(now int64) {
	this.base.SCount = 0
	this.base.SSCount = 0
	this.SetDirty(baseField)

	this.FlushAllToClient()

	this.timedata[module.WeeklyTimeName] = module.CalWeeklyTime().Unix()
	this.SetDirty(timeField)
}
