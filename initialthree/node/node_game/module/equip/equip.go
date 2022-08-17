package equip

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	flyfish "github.com/sniperHW/flyfish/client"
	"initialthree/node/node_game/event"
	"initialthree/node/table/excel/DataTable/EquipLevelMaxExp"
	"initialthree/node/table/excel/DataTable/EquipRandomAttributePool"
	"initialthree/pkg/json"
	"initialthree/zaplogger"

	"initialthree/node/node_game/module"
	tableEquip "initialthree/node/table/excel/DataTable/Equip"
	"initialthree/protocol/cs/message"
	"time"
)

const (
	slotCount = 10
)

func slotName(slotIdx int) string {
	return fmt.Sprintf("slot%d", slotIdx)
}

func calcSlotIdx(id int) int {
	return id % slotCount
}

type Equip struct {
	InsID            uint32  `json:"iid"`
	ConfigID         int32   `json:"cid"`
	Level            int32   `json:"lev"`
	Exp              int32   `json:"exp"`  // 当前经验
	RandomAttribId   int32   `json:"raid"` // 随机的技能ID，位置1
	Refine           []int32 `json:"ref"`  // 长度为2，第一个是配置表的固定技能ID，第二个是随机ID
	EquipCharacterId int32   `json:"ecid"` //装备后的角色ID
	IsLock           bool    `json:"il"`
	GetTime          int64   `json:"gt"` // 获取时间
}

func (v *Equip) Pack() *message.Equip {
	return &message.Equip{
		ID:               proto.Uint32(v.InsID),
		ConfigID:         proto.Int32(v.ConfigID),
		Level:            proto.Int32(v.Level),
		Exp:              proto.Int32(v.Exp),
		RefineLevel:      v.Refine,
		RandomAttribId:   proto.Int32(v.RandomAttribId),
		EquipCharacterID: proto.Int32(v.EquipCharacterId),
		IsLock:           proto.Bool(v.IsLock),
		GetTime:          proto.Int64(v.GetTime),
	}
}

type UserEquip struct {
	userI module.UserI

	// 数据
	equips []map[uint32]*Equip

	equipUseCap int // 装备占用

	syncDirty      map[uint32]struct{}
	syncGroupDirty bool
	*module.ModuleSaveBase
}

func (this *UserEquip) setDirty(insID uint32) {
	slotIdx := calcSlotIdx(int(insID))
	this.syncDirty[insID] = struct{}{}
	this.SetDirty(slotIdx)
}

func (this *UserEquip) NewEquip(cfgID int32, insID uint32) *Equip {
	def := tableEquip.GetID(cfgID)
	id := int32(0)
	randPool := EquipRandomAttributePool.GetID(def.RandomAttribPool)
	if randPool != nil {
		id = randPool.RandomID()
	}

	return &Equip{
		InsID:          insID,
		ConfigID:       cfgID,
		Level:          1,
		RandomAttribId: id,
		Refine:         make([]int32, 2),
		GetTime:        time.Now().Unix(),
	}

}

func (this *UserEquip) Range(cb func(e *Equip) bool) {
	for _, slot := range this.equips {
		for _, e := range slot {
			if !cb(e) {
				return
			}
		}
	}
}

func (this *UserEquip) AddEquip(e *Equip) {
	slotIdx := calcSlotIdx(int(e.InsID))
	this.equips[slotIdx][e.InsID] = e
	this.equipUseCap += 1
	this.setDirty(e.InsID)
}

func (this *UserEquip) Remove(insId uint32) {
	slotIdx := calcSlotIdx(int(insId))
	this.equipUseCap -= 1
	delete(this.equips[slotIdx], insId)
	this.setDirty(insId)
}

func (this *UserEquip) GetEquip(insID uint32) *Equip {
	slotIdx := calcSlotIdx(int(insID))
	return this.equips[slotIdx][insID]
}

func (this *UserEquip) GetUseCap() int {
	return this.equipUseCap
}

func (this *UserEquip) Equip(e *Equip, equipCId int32) {
	e.EquipCharacterId = equipCId
	this.setDirty(e.InsID)
}

func (this *UserEquip) Demount(e *Equip) {
	e.EquipCharacterId = 0
	this.setDirty(e.InsID)
}

func (this *UserEquip) Strengthen(levelDef *EquipLevelMaxExp.EquipLevelMaxExp, e *Equip, maxLev, exp int32) {
	e.Exp += exp
	needExp := levelDef.GetMaxExp(e.Level)
	oldLevel := e.Level
	for e.Level < maxLev && e.Exp >= needExp {
		e.Level++
		e.Exp -= needExp
		needExp = levelDef.GetMaxExp(e.Level)
	}

	if e.Level == maxLev {
		e.Exp = 0
	}
	newLevel := e.Level

	if newLevel != oldLevel {
		this.userI.EmitEvent(event.EventEquipLevelUp, e.InsID, oldLevel, newLevel)
		if e.EquipCharacterId != 0 {
			this.userI.EmitEvent(event.EventEquipEquipped)
		}
	}
	this.setDirty(e.InsID)
}

func (this *UserEquip) Refine(e *Equip, pos, lev, maxLevel int32) {
	oldLevel := e.Refine[pos]
	e.Refine[pos] += lev
	if e.Refine[pos] > maxLevel {
		e.Refine[pos] = maxLevel
	}
	newLevel := e.Refine[pos]

	if newLevel != oldLevel {
		this.userI.EmitEvent(event.EventEquipRefine, e.InsID, oldLevel, newLevel)
	}
	this.setDirty(e.InsID)
}

func (this *UserEquip) Lock(e *Equip, b bool) {
	if b != e.IsLock {
		e.IsLock = b
		this.setDirty(e.InsID)
	}
}

func (this *UserEquip) ModuleType() module.ModuleType {
	return module.Equip
}

func (this *UserEquip) Tick(now time.Time) {}

func (this *UserEquip) ReadOut() *module.ReadOutCommand {
	out := &module.ReadOutCommand{
		Table:  this.ModuleType().String(),
		Key:    this.userI.GetIDStr(),
		Module: this,
		Fields: make([]string, 0, slotCount),
	}
	for i := 0; i < slotCount; i++ {
		out.Fields = append(out.Fields, slotName(i))
	}

	return out
}

func (this *UserEquip) WriteCommand(fields map[interface{}]struct{}) *module.WriteBackCommand {
	cmd := &module.WriteBackCommand{
		Table:  this.ModuleType().String(),
		Key:    this.userI.GetIDStr(),
		Module: this,
		Fields: make([]*module.WriteBackFiled, 0, len(fields)),
	}

	for field := range fields {
		switch field.(type) {
		case int:
			idx := field.(int)
			slice := this.equips[idx]
			data, err := json.Marshal(slice)
			if nil != err {
				zaplogger.GetSugar().Error(err.Error())
				return nil
			}

			cmd.Fields = append(cmd.Fields, &module.WriteBackFiled{
				Name:  slotName(idx),
				Value: data,
			})
		}
	}

	return cmd
}

func (this *UserEquip) FlushDirtyToClient() {
	if len(this.syncDirty) > 0 {
		msg := &message.EquipSyncToC{
			IsAll:  proto.Bool(false),
			Equips: make([]*message.Equip, 0, len(this.syncDirty)),
			UseCap: proto.Int32(int32(this.equipUseCap)),
		}
		for id := range this.syncDirty {
			slotIdx := calcSlotIdx(int(id))
			v := this.equips[slotIdx][id]
			if v != nil {
				msg.Equips = append(msg.Equips, v.Pack())
			} else { // 意识被删除
				msg.Equips = append(msg.Equips, &message.Equip{
					ID:       proto.Uint32(id),
					IsRemove: proto.Bool(true),
				})
			}
		}
		this.syncDirty = map[uint32]struct{}{}
		this.userI.Post(msg)
	}
}

func (this *UserEquip) FlushAllToClient(seqNo ...uint32) {
	msg1 := &message.EquipSyncToC{
		IsAll:  proto.Bool(true),
		Equips: make([]*message.Equip, 0, 16),
		UseCap: proto.Int32(int32(this.equipUseCap)),
	}
	for _, sliceMap := range this.equips {
		for _, v := range sliceMap {
			msg1.Equips = append(msg1.Equips, v.Pack())
		}
	}
	this.userI.Post(msg1)
	this.syncDirty = map[uint32]struct{}{}
}

func (this *UserEquip) Init(fields map[string]*flyfish.Field) error {
	for i := 0; i < slotCount; i++ {
		fieldName := slotName(i)
		field, ok := fields[fieldName]
		this.equips[i] = map[uint32]*Equip{}
		if !ok || len(field.GetBlob()) == 0 {
			this.SetDirty(i)
		} else {
			if err := json.Unmarshal(field.GetBlob(), &this.equips[i]); err != nil {
				return fmt.Errorf("unmarshal:%s %s", string(field.GetBlob()), err)
			}
			this.equipUseCap += len(this.equips[i])
		}
	}

	return nil
}

func (this *UserEquip) Query(arg *message.QueryRoleInfoArg, ret *message.QueryRoleInfoResult) error {
	ids := arg.GetEquipIDs()
	ret.Equips = make([]*message.Equip, 0, len(ids))
	for _, id := range ids {
		slotIdx := calcSlotIdx(int(id))
		v := this.equips[slotIdx][id]
		if v != nil {
			ret.Equips = append(ret.Equips, v.Pack())
		}
	}
	return nil
}

func init() {
	module.RegisterModule(module.Equip, func(userI module.UserI) module.ModuleI {
		m := &UserEquip{
			userI:     userI,
			equips:    make([]map[uint32]*Equip, slotCount),
			syncDirty: map[uint32]struct{}{},
		}
		m.ModuleSaveBase = module.NewModuleSaveBase(m)
		return m
	})
}
