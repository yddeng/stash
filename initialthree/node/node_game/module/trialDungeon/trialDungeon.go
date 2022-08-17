package trialDungeon

import (
	"fmt"
	"github.com/gogo/protobuf/proto"
	flyfish "github.com/sniperHW/flyfish/client"
	"initialthree/node/node_game/event"
	TrialDungeon2 "initialthree/node/table/excel/DataTable/TrialDungeon"
	"initialthree/pkg/json"
	"initialthree/zaplogger"

	"initialthree/node/node_game/module"
	"initialthree/protocol/cs/message"
	"time"
)

const (
	tableName        = "trialdungeon"
	dungeonDataField = "dungeon_data"
	trialCountField  = "trial_count"
	stageRewardField = "stage_reward"
)

var tableFields = []string{dungeonDataField, trialCountField, stageRewardField}

type TrialDungeon struct {
	userI        module.UserI
	dungeonData  map[int32]*message.TrialDungeon
	stageReward  map[int32]*message.TrialStageReward
	trialCount   int64
	dirtyDungeon map[int32]struct{}
	dirtyCount   bool
	dirtyStage   map[int32]struct{}
	*module.ModuleSaveBase
}

func (this *TrialDungeon) GetTrialCount() int32 {
	return int32(this.trialCount)
}

func (this *TrialDungeon) AddTrialCount(trialCount int32) {
	this.trialCount += int64(trialCount)
	this.dirtyCount = true
	this.SetDirty(trialCountField)

	this.userI.EmitEvent(event.EventTrailCount, int32(this.trialCount-int64(trialCount)), int32(this.trialCount))
}

func (this *TrialDungeon) GetTrialDungeon(levelId int32) *message.TrialDungeon {
	return this.dungeonData[levelId]
}

func (this *TrialDungeon) Pass(levelID int32) bool {
	data := this.GetTrialDungeon(levelID)
	if data == nil {
		data = &message.TrialDungeon{
			DungeonID: proto.Int32(levelID),
		}
		this.dungeonData[levelID] = data
		this.dirtyDungeon[levelID] = struct{}{}
		this.SetDirty(dungeonDataField)

		def := TrialDungeon2.GetID(levelID)
		this.AddTrialCount(def.TrialCount)
		return true
	}
	return false
}

func (this *TrialDungeon) DungeonReward(dungeon int32) {
	dun := this.GetTrialDungeon(dungeon)
	dun.GetReward = proto.Bool(true)
	this.dirtyDungeon[dungeon] = struct{}{}
	this.SetDirty(dungeonDataField)
}

func (this *TrialDungeon) GetStageReward(stage int32) *message.TrialStageReward {
	return this.stageReward[stage]
}

func (this *TrialDungeon) StageReward(stage int32) {
	data := this.GetStageReward(stage)
	if data == nil {
		data = &message.TrialStageReward{
			Stage: proto.Int32(stage),
		}
		this.stageReward[stage] = data
		this.dirtyStage[stage] = struct{}{}
		this.SetDirty(stageRewardField)
	}
}

func (this *TrialDungeon) ModuleType() module.ModuleType {
	return module.TrialDungeon
}

func (this *TrialDungeon) Init(fields map[string]*flyfish.Field) error {
	for _, name := range tableFields {
		field, ok := fields[name]
		if ok {
			var err error
			switch name {
			case dungeonDataField:
				if len(field.GetBlob()) != 0 {
					err = json.Unmarshal(field.GetBlob(), &this.dungeonData)
				}
			case stageRewardField:
				if len(field.GetBlob()) != 0 {
					err = json.Unmarshal(field.GetBlob(), &this.stageReward)
				}
			case trialCountField:
				this.trialCount = field.GetInt()
			}
			if err != nil {
				zaplogger.GetSugar().Errorf("%s initTrialDungeon name %s err %s ", this.userI.GetUserID(), name, err)
				return fmt.Errorf("unmarshal: %s", err)
			}
		}
	}
	return nil
}

func (this *TrialDungeon) ReadOut() *module.ReadOutCommand {
	cmd := &module.ReadOutCommand{
		Table:  tableName,
		Key:    this.userI.GetIDStr(),
		Fields: tableFields,
		Module: this,
	}
	return cmd
}

func (this *TrialDungeon) WriteCommand(fields map[interface{}]struct{}) *module.WriteBackCommand {
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
		case dungeonDataField:
			data, _ = json.Marshal(this.dungeonData)
		case stageRewardField:
			data, _ = json.Marshal(this.stageReward)
		case trialCountField:
			cmd.Fields = append(cmd.Fields, &module.WriteBackFiled{
				Name:  name,
				Value: this.trialCount,
			})
			continue
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

func (this *TrialDungeon) Tick(now time.Time) {}

func (this *TrialDungeon) FlushDirtyToClient() {
	if len(this.dirtyDungeon) > 0 || this.dirtyCount || len(this.dirtyStage) > 0 {
		msg := new(message.TrialDungeonSyncToC)
		msg.All = proto.Bool(false)
		msg.TrailCount = proto.Int32(int32(this.trialCount))
		msg.TrialDungeons = make([]*message.TrialDungeon, 0, len(this.dirtyDungeon))
		msg.StageReward = make([]*message.TrialStageReward, 0, len(this.dirtyStage))

		for id := range this.dirtyDungeon {
			level := this.GetTrialDungeon(id)
			msg.TrialDungeons = append(msg.TrialDungeons, level)
		}
		for id := range this.dirtyStage {
			level := this.GetStageReward(id)
			msg.StageReward = append(msg.StageReward, level)
		}

		this.userI.Post(msg)
		this.dirtyDungeon = map[int32]struct{}{}
		this.dirtyStage = map[int32]struct{}{}
		this.dirtyCount = false
	}

}
func (this *TrialDungeon) FlushAllToClient(seqNo ...uint32) {
	msg := new(message.TrialDungeonSyncToC)
	msg.All = proto.Bool(true)
	msg.TrailCount = proto.Int32(int32(this.trialCount))
	msg.TrialDungeons = make([]*message.TrialDungeon, 0, len(this.dungeonData))
	msg.StageReward = make([]*message.TrialStageReward, 0, len(this.stageReward))

	for _, levelData := range this.dungeonData {
		msg.TrialDungeons = append(msg.TrialDungeons, levelData)
	}
	for _, levelData := range this.stageReward {
		msg.StageReward = append(msg.StageReward, levelData)
	}

	this.userI.Post(msg)
}

func init() {
	module.RegisterModule(module.TrialDungeon, func(userI module.UserI) module.ModuleI {
		m := &TrialDungeon{
			userI:        userI,
			dungeonData:  map[int32]*message.TrialDungeon{},
			stageReward:  map[int32]*message.TrialStageReward{},
			dirtyDungeon: map[int32]struct{}{},
			dirtyStage:   map[int32]struct{}{},
		}

		m.ModuleSaveBase = module.NewModuleSaveBase(m)
		return m
	})
}
