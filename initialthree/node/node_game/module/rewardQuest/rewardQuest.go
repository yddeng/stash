package rewardQuest

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	flyfish "github.com/sniperHW/flyfish/client"
	"initialthree/node/common/enumType"
	"initialthree/pkg/json"
	"initialthree/zaplogger"

	"initialthree/node/node_game/module"
	RewardQuest2 "initialthree/node/table/excel/DataTable/RewardQuest"
	"initialthree/protocol/cs/message"
	"time"
)

const (
	tableName = "rewardquest"
	dataField = "data"
	baseField = "base"
	timeField = "timedata"
)

var tableFields = []string{baseField, dataField, timeField}

type Base struct {
	//所处地图及坐标
	SSCount      int32 `json:"ssc"` // 每周重置时为0
	SCount       int32 `json:"sc"`  // 每周重置时为0
	RefreshTimes int32 `json:"rt"`  // 每日重置时为0
}

type RewardQuest struct {
	userI    module.UserI
	base     Base
	data     map[int32]*message.RewardQuest
	timedata map[string]int64
	dirty    map[int32]struct{}
	*module.ModuleSaveBase
}

func (this *RewardQuest) GetData() map[int32]*message.RewardQuest {
	return this.data
}

func (this *RewardQuest) GetRewardQuest(id int32) *message.RewardQuest {
	return this.data[id]
}

// 返回正在任务中的角色，完成任务的角色可再次使用
func (this *RewardQuest) UsedRoles() map[int32]struct{} {
	roles := map[int32]struct{}{}
	for _, v := range this.data {
		if v.GetState() == message.QuestState_Running {
			for _, id := range v.GetCharacters() {
				roles[id] = struct{}{}
			}
		}
	}

	return roles
}

func (this *RewardQuest) GetBase() Base {
	return this.base
}

func (this *RewardQuest) Accept(q *message.RewardQuest, roles []int32, acceptTime int64) {
	q.State = message.QuestState_Running.Enum()
	q.Characters = roles
	q.AcceptTimestamp = proto.Int64(acceptTime)

	this.dirty[q.GetQuestID()] = struct{}{}
	this.SetDirty(dataField)
}

func (this *RewardQuest) Complete(q *message.RewardQuest) {
	q.State = message.QuestState_End.Enum()
	this.dirty[q.GetQuestID()] = struct{}{}
	this.SetDirty(dataField)
}

// removed 移除的ID，add 添加的pos和quality key->pos,value->quality
func (this *RewardQuest) Replace(removed map[int32]struct{}, add map[int32]int32, isRefresh bool) {
	if isRefresh {
		this.base.RefreshTimes++
		this.SetDirty(baseField)
	}

	for id := range removed {
		delete(this.data, id)
		this.dirty[id] = struct{}{}
	}

	for pos, quality := range add {
		def := RewardQuest2.GetPosQuality(pos, quality)
		if def != nil {
			switch quality {
			case enumType.RarityType_Star5:
				this.base.SSCount++
				this.SetDirty(baseField)
			case enumType.RarityType_Star4:
				this.base.SCount++
				this.SetDirty(baseField)
			}

			this.data[def.ID] = &message.RewardQuest{
				QuestID: proto.Int32(def.ID),
				State:   message.QuestState_Acceptable.Enum(),
			}
			this.dirty[def.ID] = struct{}{}
		}
	}
	this.SetDirty(dataField)
}

func (this *RewardQuest) ModuleType() module.ModuleType {
	return module.RewardQuest
}

func (this *RewardQuest) Init(fields map[string]*flyfish.Field) error {
	for _, name := range tableFields {
		field, ok := fields[name]
		if ok && len(field.GetBlob()) != 0 {
			var err error
			switch name {
			case baseField:
				err = json.Unmarshal(field.GetBlob(), &this.base)
			case dataField:
				err = json.Unmarshal(field.GetBlob(), &this.data)
			case timeField:
				err = json.Unmarshal(field.GetBlob(), &this.timedata)
			}
			if err != nil {
				zaplogger.GetSugar().Errorf("%s initRewardQuest name %s err %s ", this.userI.GetUserID(), name, err)
				return fmt.Errorf("unmarshal: %s", err)
			}
		}
	}

	//log.GetLogger().Infoln(this.userI.GetUserLogName(), "init reward quest ok ")
	return nil
}

func (this *RewardQuest) AfterInitAll() error {
	this.initTime()
	return nil
}

func (this *RewardQuest) ReadOut() *module.ReadOutCommand {
	cmd := &module.ReadOutCommand{
		Table:  tableName,
		Key:    this.userI.GetIDStr(),
		Fields: tableFields,
		Module: this,
	}
	return cmd
}

func (this *RewardQuest) WriteCommand(fields map[interface{}]struct{}) *module.WriteBackCommand {
	cmd := &module.WriteBackCommand{
		Table:  tableName,
		Key:    this.userI.GetIDStr(),
		Fields: make([]*module.WriteBackFiled, 0, len(fields)),
		Module: this,
	}

	for field := range fields {
		name := field.(string)
		var data []byte
		switch name {
		case baseField:
			data, _ = json.Marshal(this.base)
		case dataField:
			data, _ = json.Marshal(this.data)
		case timeField:
			data, _ = json.Marshal(this.timedata)
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

func (this *RewardQuest) Tick(now time.Time) {}

func (this *RewardQuest) FlushDirtyToClient() {
	if len(this.dirty) > 0 {
		msg := new(message.RewardQuestSyncToC)
		msg.IsAll = proto.Bool(false)
		msg.RefreshTimes = proto.Int32(this.base.RefreshTimes)
		msg.Quests = make([]*message.RewardQuest, 0, len(this.dirty))

		for id := range this.dirty {
			v, ok := this.data[id]
			if !ok {
				v = &message.RewardQuest{
					QuestID:   proto.Int32(id),
					IsRemoved: proto.Bool(true),
				}
			}
			msg.Quests = append(msg.Quests, v)
		}

		this.userI.Post(msg)
		this.dirty = map[int32]struct{}{}
	}

}
func (this *RewardQuest) FlushAllToClient(seqNo ...uint32) {
	msg := new(message.RewardQuestSyncToC)
	msg.IsAll = proto.Bool(true)
	msg.RefreshTimes = proto.Int32(this.base.RefreshTimes)
	msg.Quests = make([]*message.RewardQuest, 0, len(this.data))
	for _, v := range this.data {
		msg.Quests = append(msg.Quests, v)
	}
	this.dirty = map[int32]struct{}{}
	this.userI.Post(msg)
}
func init() {
	module.RegisterModule(module.RewardQuest, func(userI module.UserI) module.ModuleI {
		m := &RewardQuest{
			userI:    userI,
			data:     map[int32]*message.RewardQuest{},
			timedata: map[string]int64{},
			dirty:    map[int32]struct{}{},
		}

		m.ModuleSaveBase = module.NewModuleSaveBase(m)
		return m
	})
}
