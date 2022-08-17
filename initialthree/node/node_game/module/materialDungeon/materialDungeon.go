package materialDungeon

import (
	"fmt"
	"github.com/gogo/protobuf/proto"
	flyfish "github.com/sniperHW/flyfish/client"
	"initialthree/pkg/json"
	"initialthree/zaplogger"

	"initialthree/node/node_game/module"
	"initialthree/protocol/cs/message"
	"time"
)

const (
	tableName        = "materialdungeon"
	dungeonDataField = "dungeon_data"
)

var tableFields = []string{dungeonDataField}

type MaterialDungeon struct {
	userI       module.UserI
	dungeonData map[int32]*message.MaterialDungeon
	dirtyLevel  map[int32]struct{}
	*module.ModuleSaveBase
}

func (this *MaterialDungeon) GetMaterialDungeon(levelId int32) *message.MaterialDungeon {
	return this.dungeonData[levelId]
}

func (this *MaterialDungeon) Pass(levelID int32) bool {
	data := this.GetMaterialDungeon(levelID)
	if data == nil {
		data = &message.MaterialDungeon{
			DungeonID: proto.Int32(levelID),
		}
		this.dungeonData[levelID] = data
		this.dirtyLevel[levelID] = struct{}{}
		this.SetDirty(dungeonDataField)
		return true
	}
	return false
}

func (this *MaterialDungeon) ModuleType() module.ModuleType {
	return module.MaterialDungeon
}

func (this *MaterialDungeon) Init(fields map[string]*flyfish.Field) error {
	for _, name := range tableFields {
		field, ok := fields[name]
		if ok && len(field.GetBlob()) != 0 {
			var err error
			switch name {
			case dungeonDataField:
				err = json.Unmarshal(field.GetBlob(), &this.dungeonData)
			}
			if err != nil {
				zaplogger.GetSugar().Errorf("%s initMaterialDungeon name %s err %s ", this.userI.GetUserID(), name, err)
				return fmt.Errorf("unmarshal: %s", err)
			}
		}
	}
	return nil
}

func (this *MaterialDungeon) ReadOut() *module.ReadOutCommand {
	cmd := &module.ReadOutCommand{
		Table:  tableName,
		Key:    this.userI.GetIDStr(),
		Fields: tableFields,
		Module: this,
	}
	return cmd
}

func (this *MaterialDungeon) WriteCommand(fields map[interface{}]struct{}) *module.WriteBackCommand {
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

func (this *MaterialDungeon) Tick(now time.Time) {}

func (this *MaterialDungeon) FlushDirtyToClient() {
	if len(this.dirtyLevel) > 0 {
		msg := new(message.MaterialDungeonSyncToC)
		msg.All = proto.Bool(false)
		msg.MaterialDungeons = make([]*message.MaterialDungeon, 0, len(this.dirtyLevel))

		for id := range this.dirtyLevel {
			level := this.GetMaterialDungeon(id)
			msg.MaterialDungeons = append(msg.MaterialDungeons, level)
		}

		this.userI.Post(msg)
		this.dirtyLevel = map[int32]struct{}{}
	}

}
func (this *MaterialDungeon) FlushAllToClient(seqNo ...uint32) {
	msg := new(message.MaterialDungeonSyncToC)
	msg.All = proto.Bool(true)
	msg.MaterialDungeons = make([]*message.MaterialDungeon, 0, len(this.dungeonData))

	for _, levelData := range this.dungeonData {
		msg.MaterialDungeons = append(msg.MaterialDungeons, levelData)
	}

	this.userI.Post(msg)
}

func init() {
	module.RegisterModule(module.MaterialDungeon, func(userI module.UserI) module.ModuleI {
		m := &MaterialDungeon{
			userI:       userI,
			dungeonData: map[int32]*message.MaterialDungeon{},
			dirtyLevel:  map[int32]struct{}{},
		}

		m.ModuleSaveBase = module.NewModuleSaveBase(m)
		return m
	})
}
