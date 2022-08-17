package worldQuest

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	flyfish "github.com/sniperHW/flyfish/client"
	event2 "initialthree/node/node_game/event"
	"initialthree/node/node_game/module"
	"initialthree/node/node_game/module/maindungeons"
	"initialthree/node/table/excel/ConstTable/PlayerCamp"
	"initialthree/node/table/excel/ConstTable/Quest"
	"initialthree/node/table/excel/DataTable/PlayerCampReputationLevel"
	WorldQuest2 "initialthree/node/table/excel/DataTable/WorldQuest"
	"initialthree/pkg/event"
	"initialthree/pkg/json"
	"initialthree/protocol/cs/message"
	"initialthree/zaplogger"
	"time"
)

const (
	table           = "worldquest"
	questField      = "quest"
	reputationField = "reputation"
	timeField       = "timedata"
	shopField       = "shop"
)

var tableFields = []string{questField, reputationField, timeField, shopField}

type CampLevel struct {
	Camp    int32 `json:"camp"`
	Level   int32 `json:"level"`
	Current int32 `json:"current"`
}

type ShopData struct {
	BuyTimes         map[int32]int32 `json:"buy_times"`
	CampRefreshTimes map[int32]int32 `json:"camp_refresh_times"`
}

type QuestData struct {
	DoneQuests    map[int32]struct{}    `json:"done_quests"`    // 每日完成的任务
	CurrentQuests []*message.WorldQuest `json:"current_quests"` // 当前任务完成状态
	RefreshTimes  int32                 `json:"refresh_times"`
	DoneTimes     int32                 `json:"done_times"`
}

type WorldQuest struct {
	userI module.UserI

	questData QuestData
	campLevel map[int32]*CampLevel
	shopData  ShopData
	timedata  map[string]int64

	dirtyQuestData   bool
	dirtyCampLev     map[int32]struct{}
	dirtyBuyTimes    map[int32]struct{}
	dirtyCampRefresh map[int32]struct{}

	unlocked    bool
	activeEvent event.Handle // 注册解锁事件
	*module.ModuleSaveBase
}

func (this *WorldQuest) clearCampRefreshTimes() {
	for id := range this.shopData.CampRefreshTimes {
		this.dirtyCampRefresh[id] = struct{}{}
	}
	this.shopData.CampRefreshTimes = map[int32]int32{}
}

func (this *WorldQuest) clearBuyTimes(ids []int32) {
	for _, id := range ids {
		delete(this.shopData.BuyTimes, id)
		this.dirtyBuyTimes[id] = struct{}{}
	}
}

func (this *WorldQuest) GetWQRefreshTimes() int32 {
	return this.questData.RefreshTimes
}

func (this *WorldQuest) WQRefresh() {
	this.questData.RefreshTimes += 1
	this.questRefresh(true)
}

func (this *WorldQuest) ShopItem(id, count int32) {
	this.shopData.BuyTimes[id] += count
	this.dirtyBuyTimes[id] = struct{}{}
	this.SetDirty(shopField)
}

func (this *WorldQuest) ShopItemTimes(id int32) int32 {
	return this.shopData.BuyTimes[id]
}

func (this *WorldQuest) CampRefresh(id int32) {
	this.shopData.CampRefreshTimes[id] += 1
	this.dirtyCampRefresh[id] = struct{}{}
	this.SetDirty(shopField)
}

func (this *WorldQuest) CampRefreshTimes(id int32) int32 {
	return this.shopData.CampRefreshTimes[id]
}

func (this *WorldQuest) addCampReputation(campType, count int32) {
	rep, ok := this.campLevel[campType]
	if !ok {
		rep = &CampLevel{
			Camp:    campType,
			Level:   1,
			Current: 0,
		}
		this.campLevel[campType] = rep
	}
	rep.Current += count

	def := PlayerCampReputationLevel.GetID(rep.Level)
	for def != nil && def.ReputationValue > 0 && rep.Current >= def.ReputationValue {
		rep.Level++
		rep.Current -= def.ReputationValue

		def = PlayerCampReputationLevel.GetID(rep.Level)
	}
	if def == nil || def.ReputationValue == 0 {
		rep.Current = 0
	}

	this.dirtyCampLev[campType] = struct{}{}
	this.SetDirty(reputationField)
}

func (this *WorldQuest) CanDo() bool {
	eachCount := Quest.GetID(1).WorldQuestEachCount
	if this.questData.DoneTimes >= eachCount {
		return false
	}
	return true
}

func (this *WorldQuest) Pass(questID int32, cfg *WorldQuest2.WorldQuest) bool {
	eachCount := Quest.GetID(1).WorldQuestEachCount
	if this.questData.DoneTimes >= eachCount {
		return false
	}

	for _, v := range this.questData.CurrentQuests {
		if v.GetQuestID() == questID && !v.GetDone() {
			v.Done = proto.Bool(true)
			this.questData.DoneTimes++
			this.questData.DoneQuests[questID] = struct{}{}
			this.dirtyQuestData = true
			this.SetDirty(questField)

			// 声望奖励
			this.addCampReputation(cfg.CampTypeEnum, cfg.ReputationReward)
			return true
		}
	}

	return false
}

func (this *WorldQuest) ModuleType() module.ModuleType {
	return module.WorldQuest
}

func (this *WorldQuest) Init(fields map[string]*flyfish.Field) error {
	for _, name := range tableFields {
		field, ok := fields[name]
		if ok && len(field.GetBlob()) != 0 {
			var err error
			switch name {
			case questField:
				err = json.Unmarshal(field.GetBlob(), &this.questData)
			case reputationField:
				err = json.Unmarshal(field.GetBlob(), &this.campLevel)
			case timeField:
				err = json.Unmarshal(field.GetBlob(), &this.timedata)
			case shopField:
				err = json.Unmarshal(field.GetBlob(), &this.shopData)
			}
			if err != nil {
				zaplogger.GetSugar().Errorf("%s init worldQuest name %s err %s ", this.userI.GetUserID(), name, err)
				return fmt.Errorf("unmarshal: %s", err)
			}
		}
	}
	return nil
}

func (this *WorldQuest) ReadOut() *module.ReadOutCommand {
	return &module.ReadOutCommand{
		Table:  table,
		Key:    this.userI.GetIDStr(),
		Fields: tableFields,
		Module: this,
	}
}

func (this *WorldQuest) WriteCommand(fields map[interface{}]struct{}) *module.WriteBackCommand {
	cmd := &module.WriteBackCommand{
		Table:  table,
		Key:    this.userI.GetIDStr(),
		Module: this,
		Fields: make([]*module.WriteBackFiled, 0, len(fields)),
	}

	for field := range fields {
		name := field.(string)
		var data []byte
		switch name {
		case questField:
			data, _ = json.Marshal(this.questData)
		case reputationField:
			data, _ = json.Marshal(this.campLevel)
		case timeField:
			data, _ = json.Marshal(this.timedata)
		case shopField:
			data, _ = json.Marshal(this.shopData)
		default:
			continue
		}
		cmd.Fields = append(cmd.Fields, &module.WriteBackFiled{
			Name:  name,
			Value: data,
		})
	}

	return cmd
}

func (this *WorldQuest) Tick(now time.Time) {
	if this.unlocked {
		this.clockTimer()
	}
}

func (this *WorldQuest) FlushDirtyToClient() {
	if this.dirtyQuestData {
		msg := &message.WorldQuestSyncToC{
			IsAll:        proto.Bool(true),
			DoneTimes:    proto.Int32(this.questData.DoneTimes),
			RefreshTimes: proto.Int32(this.questData.RefreshTimes),
			WorldQuests:  this.questData.CurrentQuests,
		}

		this.userI.Post(msg)
		this.dirtyQuestData = false
	}

	if len(this.dirtyCampLev) > 0 || len(this.dirtyBuyTimes) > 0 || len(this.dirtyCampRefresh) > 0 {
		msg := &message.ReputationSyncToC{
			IsAll:           proto.Bool(false),
			Reputations:     make([]*message.Reputation, 0, len(this.dirtyCampLev)),
			ReputationItems: make([]*message.ReputationItem, 0, len(this.dirtyBuyTimes)),
			Refresh:         make([]*message.ReputationRefresh, 0, len(this.dirtyCampRefresh)),
		}

		for i := range this.dirtyCampLev {
			v := this.campLevel[i]
			msg.Reputations = append(msg.Reputations, &message.Reputation{
				CampType:          proto.Int32(v.Camp),
				ReputationLevel:   proto.Int32(v.Level),
				CurrentReputation: proto.Int32(v.Current),
			})
		}

		for id := range this.dirtyBuyTimes {
			msg.ReputationItems = append(msg.ReputationItems, &message.ReputationItem{
				Id:    proto.Int32(id),
				Count: proto.Int32(this.shopData.BuyTimes[id]),
			})
		}

		for id := range this.dirtyCampRefresh {
			msg.Refresh = append(msg.Refresh, &message.ReputationRefresh{
				Id:    proto.Int32(id),
				Times: proto.Int32(this.shopData.CampRefreshTimes[id]),
			})
		}

		this.userI.Post(msg)
		this.dirtyCampLev = map[int32]struct{}{}
		this.dirtyCampRefresh = map[int32]struct{}{}
		this.dirtyBuyTimes = map[int32]struct{}{}
	}

}

func (this *WorldQuest) FlushAllToClient(seqNo ...uint32) {
	msg := &message.WorldQuestSyncToC{
		IsAll:        proto.Bool(true),
		DoneTimes:    proto.Int32(this.questData.DoneTimes),
		RefreshTimes: proto.Int32(this.questData.RefreshTimes),
		WorldQuests:  this.questData.CurrentQuests,
	}

	this.userI.Post(msg)
	this.dirtyQuestData = false

	msg1 := &message.ReputationSyncToC{
		IsAll:           proto.Bool(true),
		Reputations:     make([]*message.Reputation, 0, len(this.campLevel)),
		ReputationItems: make([]*message.ReputationItem, 0, len(this.shopData.BuyTimes)),
		Refresh:         make([]*message.ReputationRefresh, 0, len(this.shopData.CampRefreshTimes)),
	}

	for id, count := range this.shopData.BuyTimes {
		msg1.ReputationItems = append(msg1.ReputationItems, &message.ReputationItem{
			Id:    proto.Int32(id),
			Count: proto.Int32(count),
		})
	}

	for id, times := range this.shopData.CampRefreshTimes {
		msg1.Refresh = append(msg1.Refresh, &message.ReputationRefresh{
			Id:    proto.Int32(id),
			Times: proto.Int32(times),
		})
	}

	for _, v := range this.campLevel {
		msg1.Reputations = append(msg1.Reputations, &message.Reputation{
			CampType:          proto.Int32(v.Camp),
			ReputationLevel:   proto.Int32(v.Level),
			CurrentReputation: proto.Int32(v.Current),
		})
	}
	this.userI.Post(msg1)
	this.dirtyCampLev = map[int32]struct{}{}
	this.dirtyCampRefresh = map[int32]struct{}{}
	this.dirtyBuyTimes = map[int32]struct{}{}
}

func (this *WorldQuest) AfterInitAll() error {
	unlock := map[int]struct{}{}
	userDungeon := this.userI.GetSubModule(module.MainDungeons).(*maindungeons.MainDungeons)
	camp := PlayerCamp.GetID(1)
	for i, v := range camp.UnlockDungeon {
		if v.ID == 0 || userDungeon.IsDungeonPass(v.ID) {
			unlock[i] = struct{}{}
		}
	}

	// 模块解锁
	if len(unlock) == 0 {
		this.activeEvent = this.userI.RegisterEvent(event2.EventInstanceSucceed, this.EventInstanceSucceed)
	} else {
		this.unlocked = true
	}
	return nil
}

func init() {
	module.RegisterModule(module.WorldQuest, func(userI module.UserI) module.ModuleI {
		m := &WorldQuest{
			userI: userI,
			questData: QuestData{
				DoneQuests: map[int32]struct{}{},
			},
			campLevel: map[int32]*CampLevel{},
			shopData: ShopData{
				BuyTimes:         map[int32]int32{},
				CampRefreshTimes: map[int32]int32{},
			},
			timedata:         map[string]int64{},
			dirtyBuyTimes:    map[int32]struct{}{},
			dirtyCampRefresh: map[int32]struct{}{},
			dirtyCampLev:     map[int32]struct{}{},
		}

		m.ModuleSaveBase = module.NewModuleSaveBase(m)
		return m
	})
}
