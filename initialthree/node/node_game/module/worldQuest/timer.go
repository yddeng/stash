package worldQuest

import (
	"github.com/golang/protobuf/proto"
	"initialthree/node/common/enumType"
	"initialthree/node/common/timeDisposal"
	"initialthree/node/node_game/module"
	"initialthree/node/node_game/module/maindungeons"
	"initialthree/node/table/excel/ConstTable/PlayerCamp"
	"initialthree/node/table/excel/ConstTable/Quest"
	"initialthree/node/table/excel/DataTable/PlayerCampReputationShopItem"
	"initialthree/node/table/excel/DataTable/WorldQuestPool"
	"initialthree/protocol/cs/message"
)

func (this *WorldQuest) tryClock(name string, nowUnix int64, fn func(nowUnix int64)) {
	if timestamp, ok := this.timedata[name]; ok {
		if nowUnix >= timestamp {
			fn(nowUnix)
		}
	} else {
		fn(nowUnix)
	}
}

func (this *WorldQuest) clockTimer() {
	now := timeDisposal.NowUnix()

	this.tryClock(module.MonthlyTimeName, now, this.monthlyClock)
	this.tryClock(module.WeeklyTimeName, now, this.weeklyClock)
	this.tryClock(module.DailyTimeName, now, this.dailyClock)
}

// 日更新
func (this *WorldQuest) dailyClock(now int64) {
	// 任务
	this.questRefresh(false)
	// 商品次数
	this.refreshBuyTimes(enumType.ProductLimitType_Daily)
	// 商店刷新次数
	this.clearCampRefreshTimes()

	this.SetDirty(questField, shopField)
	this.FlushAllToClient()

	this.timedata[module.DailyTimeName] = module.CalDailyTime().Unix()
	this.SetDirty(timeField)
}

func (this *WorldQuest) questRefresh(user bool) {
	if !user {
		this.questData.DoneTimes = 0
		this.questData.RefreshTimes = 0
		this.questData.DoneQuests = map[int32]struct{}{}

	}
	this.questData.CurrentQuests = []*message.WorldQuest{}

	unlock := make(map[int]struct{}, 6)

	userDungeon := this.userI.GetSubModule(module.MainDungeons).(*maindungeons.MainDungeons)
	camp := PlayerCamp.GetID(1)
	for i, v := range camp.UnlockDungeon {
		if v.ID == 0 || userDungeon.IsDungeonPass(v.ID) {
			unlock[i] = struct{}{}
		}
	}

	if len(unlock) == 0 {
		this.dirtyQuestData = true
		this.SetDirty(questField)
		return
	}

	// 刷新任务
	level := this.userI.GetLevel()
	pool := WorldQuestPool.GetByLevel(level)
	questDef := Quest.GetID(1)
	if pool != nil {
		ids := pool.RandomWorldQuest(questDef.WorldQuestEachCount, unlock, this.questData.DoneQuests)
		this.questData.CurrentQuests = make([]*message.WorldQuest, 0, len(ids))
		for id := range ids {
			this.questData.CurrentQuests = append(this.questData.CurrentQuests, &message.WorldQuest{
				QuestID: proto.Int32(id),
				Done:    proto.Bool(false),
			})
		}
	}

	// 阵营声望
	for idx := range unlock {
		// 索引转camp
		camp := int32(idx) + 1
		if _, exist := this.campLevel[camp]; !exist {
			this.addCampReputation(camp, 0)
		}
	}

	this.dirtyQuestData = true
	this.SetDirty(questField)
}

// 周更新
func (this *WorldQuest) weeklyClock(now int64) {
	// 商品次数
	this.refreshBuyTimes(enumType.ProductLimitType_Weekly)

	this.timedata[module.WeeklyTimeName] = module.CalWeeklyTime().Unix()
	this.SetDirty(timeField)
}

// 月更新
func (this *WorldQuest) monthlyClock(now int64) {
	// 商品次数
	this.refreshBuyTimes(enumType.ProductLimitType_Monthly)

	this.timedata[module.MonthlyTimeName] = module.CalMonthlyTime().Unix()
	this.SetDirty(timeField)
}

func (this *WorldQuest) refreshBuyTimes(limit int32) {
	products := make([]int32, 0, len(this.shopData.BuyTimes))
	for id := range this.shopData.BuyTimes {
		def := PlayerCampReputationShopItem.GetID(id)
		if def == nil {
			products = append(products, id)
		} else {
			if def.ProductLimitTypeEnum == limit {
				products = append(products, id)
			}
		}
	}
	this.clearBuyTimes(products)
	this.SetDirty(shopField)
}

func (this *WorldQuest) EventInstanceSucceed(instanceID, instanceStar int32, monster map[int32]int32) {
	//zaplogger.GetSugar().Info("EventInstanceComplete", this.userI.GetID(), instanceID)
	camp := PlayerCamp.GetID(1)
	for _, v := range camp.UnlockDungeon {
		if v.ID != 0 && v.ID == instanceID {
			this.unlocked = true
			this.userI.UnRegisterEvent(this.activeEvent)
			return
		}
	}
}
