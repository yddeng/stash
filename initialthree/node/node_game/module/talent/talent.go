package talent

import (
	"encoding/json"
	"fmt"
	"github.com/golang/protobuf/proto"
	flyfish "github.com/sniperHW/flyfish/client"
	"initialthree/node/node_game/module"
	"initialthree/protocol/cs/message"
	"time"
)

const (
	tableName     = "talent"
	infiniteField = "infinite"
)

var tableFields []string

func initTableFields() {
	tableFields = make([]string, 0, 1+slotCount)
	tableFields = append(tableFields, infiniteField)
	for i := 0; i < slotCount; i++ {
		tableFields = append(tableFields, slotName(i))
	}
}

const (
	slotCount = 10
)

func slotName(slotIdx int) string {
	return fmt.Sprintf("slot%d", slotIdx)
}

func calcSlotIdx(id int) int {
	return id % slotCount
}

type Talent struct {
	ID    int32 `json:"id"`
	Level int32 `json:"level"`
}

type TalentGroup struct {
	ID      int32             `json:"id"`
	Talents map[int32]*Talent `json:"talents"`
	dirty   map[int32]*Talent
}

type InfiniteTalent struct {
	Level int32 `json:"level"`
}

type UserTalent struct {
	userI    module.UserI
	groups   []map[int32]*TalentGroup
	infinite InfiniteTalent

	groupDirty    map[int32]*TalentGroup
	infiniteDirty bool
	*module.ModuleSaveBase
}

func (this *UserTalent) InfiniteTalent() *InfiniteTalent {
	return this.infinite
}

func (this *UserTalent) ModuleType() module.ModuleType {
	return module.Talent
}

func (this *UserTalent) Init(fields map[string]*flyfish.Field) error {
	for i := 0; i < slotCount; i++ {
		name := slotName(i)
		field, ok := fields[name]
		this.groups[i] = map[int32]*TalentGroup{}
		if !ok || len(field.GetBlob()) == 0 {
			this.SetDirty(i)
		} else {
			if err := json.Unmarshal(field.GetBlob(), &this.groups[i]); err != nil {
				return fmt.Errorf("unmarshal:%s %s", string(field.GetBlob()), err)
			}
		}
	}

	field, ok := fields[infiniteField]
	if ok && len(field.GetBlob()) != 0 {
		if err := json.Unmarshal(field.GetBlob(), &this.infinite); err != nil {
			return fmt.Errorf("unmarshal:%s %s", string(field.GetBlob()), err)
		}
	}

	return nil
}

func (this *UserTalent) ReadOut() *module.ReadOutCommand {
	out := &module.ReadOutCommand{
		Table:  tableName,
		Key:    this.userI.GetIDStr(),
		Module: this,
		Fields: tableFields,
	}
	return out
}

func (this *UserTalent) Tick(now time.Time) {

}

func (this *UserTalent) WriteCommand(fields map[interface{}]struct{}) *module.WriteBackCommand {
	cmd := &module.WriteBackCommand{
		Table:  tableName,
		Key:    this.userI.GetIDStr(),
		Module: this,
		Fields: make([]*module.WriteBackFiled, 0, len(fields)),
	}

	for field := range fields {
		var name string
		var data []byte
		switch field.(type) {
		case int:
			name = slotName(field.(int))
			slotIdx := field.(int)
			data, _ = json.Marshal(this.groups[slotIdx])
		case string:
			name := field.(string)
			switch name {
			case infiniteField:
				data, _ = json.Marshal(this.infinite)
			default:
				continue
			}
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

func (this *UserTalent) FlushAllToClient(seqNo ...uint32) {
	msg := &message.TalentSyncToC{
		IsAll:               proto.Bool(true),
		InfiniteTalentLevel: proto.Int32(this.infinite.Level),
		Groups:              make([]*message.TalentGroup, 0, slotCount),
	}

	for _, slot := range this.groups {
		for _, g := range slot {
			mg := &message.TalentGroup{
				GroupID: proto.Int32(g.ID),
				Talents: make([]*message.Talent, 0, len(g.Talents)),
			}

			for _, t := range g.Talents {
				mg.Talents = append(mg.Talents, &message.Talent{
					ID:    proto.Int32(t.ID),
					Level: proto.Int32(t.Level),
				})
			}
			msg.Groups = append(msg.Groups, mg)
		}
	}

	this.userI.Post(msg)
	this.groupDirty = map[int32]*TalentGroup{}
	this.infiniteDirty = false
}

func (this *UserTalent) FlushDirtyToClient() {
	if len(this.groupDirty) > 0 || this.infiniteDirty {
		msg := &message.TalentSyncToC{
			IsAll:               proto.Bool(false),
			InfiniteTalentLevel: proto.Int32(this.infinite.Level),
			Groups:              make([]*message.TalentGroup, 0, len(this.groupDirty)),
		}

		for gid, g := range this.groupDirty {
			mg := &message.TalentGroup{
				GroupID: proto.Int32(gid),
				Talents: make([]*message.Talent, 0, len(g.dirty)),
			}

			for _, t := range g.dirty {
				mg.Talents = append(mg.Talents, &message.Talent{
					ID:    proto.Int32(t.ID),
					Level: proto.Int32(t.Level),
				})
			}
			msg.Groups = append(msg.Groups, mg)
		}

		this.userI.Post(msg)
		this.groupDirty = map[int32]*TalentGroup{}
		this.infiniteDirty = false
	}

}

func init() {
	initTableFields()
	module.RegisterModule(module.Talent, func(userI module.UserI) module.ModuleI {
		m := &UserTalent{
			userI:      userI,
			groupDirty: map[int32]*TalentGroup{},
			groups:     make([]map[int32]*TalentGroup, slotCount),
		}

		m.ModuleSaveBase = module.NewModuleSaveBase(m)
		return m
	})
}
