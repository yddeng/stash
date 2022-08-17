package weapon

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	flyfish "github.com/sniperHW/flyfish/client"
	"initialthree/node/node_game/event"
	"initialthree/node/node_game/module"
	tableWeapon "initialthree/node/table/excel/DataTable/Weapon"
	"initialthree/node/table/excel/DataTable/WeaponLevelMaxExp"
	"initialthree/pkg/json"
	"initialthree/protocol/cs/message"
	"initialthree/zaplogger"
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

type Weapon struct {
	InsID            uint32 `json:"iid"`
	ConfigID         int32  `json:"cid"`
	Level            int32  `json:"lev"`
	Exp              int32  `json:"exp"` // 当前经验
	Refine           int32  `json:"ref"`
	BreakTimes       int32  `json:"bts"`  // 突破次数
	EquipCharacterID int32  `json:"ecid"` //装备后的角色ID
	IsLock           bool   `json:"il"`
	GetTime          int64  `json:"gt"` // 获取时间
}

func (v *Weapon) Pack() *message.Weapon {
	return &message.Weapon{
		ID:               proto.Uint32(v.InsID),
		ConfigID:         proto.Int32(v.ConfigID),
		Level:            proto.Int32(v.Level),
		Exp:              proto.Int32(v.Exp),
		RefineLevel:      proto.Int32(v.Refine),
		BreakLevel:       proto.Int32(v.BreakTimes),
		EquipCharacterID: proto.Int32(v.EquipCharacterID),
		IsLock:           proto.Bool(v.IsLock),
		GetTime:          proto.Int64(v.GetTime),
	}
}

type UserWeapon struct {
	userI module.UserI

	// 数据
	weapons []map[uint32]*Weapon

	weaponUseCap int // 装备占用

	syncDirty      map[uint32]struct{}
	syncGroupDirty bool
	*module.ModuleSaveBase
}

func (this *UserWeapon) setDirty(insID uint32) {
	slotIdx := calcSlotIdx(int(insID))
	this.syncDirty[insID] = struct{}{}
	this.SetDirty(slotIdx)
}

func (this *UserWeapon) NewWeapon(cfgID int32, insID uint32) *Weapon {
	return &Weapon{
		InsID:    insID,
		ConfigID: cfgID,
		GetTime:  time.Now().Unix(),
		Level:    1,
	}
}

func (this *UserWeapon) Range(cb func(w *Weapon) bool) {
	for _, slot := range this.weapons {
		for _, w := range slot {
			if !cb(w) {
				return
			}
		}
	}
}

// 判断同配置的ID是否已经存在
func (this *UserWeapon) ConfigIDIsExist(cfgID int32) bool {
	for _, s := range this.weapons {
		for _, w := range s {
			if w.ConfigID == cfgID {
				return true
			}
		}
	}
	return false
}

func (this *UserWeapon) AddWeapon(w *Weapon) {
	slotIdx := calcSlotIdx(int(w.InsID))
	this.weapons[slotIdx][w.InsID] = w
	this.weaponUseCap += 1
	this.setDirty(w.InsID)
}

func (this *UserWeapon) Remove(insId uint32) {
	slotIdx := calcSlotIdx(int(insId))
	this.weaponUseCap -= 1
	delete(this.weapons[slotIdx], insId)
	this.setDirty(insId)
}

func (this *UserWeapon) GetWeapon(insID uint32) *Weapon {
	slotIdx := calcSlotIdx(int(insID))
	return this.weapons[slotIdx][insID]
}

func (this *UserWeapon) GetUseCap() int {
	return this.weaponUseCap
}

func (this *UserWeapon) Equip(w *Weapon, equipCId int32) {
	w.EquipCharacterID = equipCId
	this.setDirty(w.InsID)
}

func (this *UserWeapon) Demount(w *Weapon) {
	w.EquipCharacterID = 0
	this.setDirty(w.InsID)
}

func (this *UserWeapon) AddBreakTimes(w *Weapon) {
	w.BreakTimes++
	this.userI.EmitEvent(event.EventWeaponBreak, w.InsID, w.BreakTimes-1, w.BreakTimes)

	this.setDirty(w.InsID)
}

func (this *UserWeapon) Strengthen(levelDef *WeaponLevelMaxExp.WeaponLevelMaxExp, w *Weapon, maxLev, exp int32) {
	w.Exp += exp
	needExp := levelDef.GetMaxExp(w.Level)
	oldLevel := w.Level
	for w.Exp >= needExp && w.Level < maxLev {
		w.Level++
		w.Exp -= needExp
		needExp = levelDef.GetMaxExp(w.Level)
	}

	if w.Level == maxLev {
		w.Exp = 0
	}

	if w.Level != oldLevel {
		this.userI.EmitEvent(event.EventWeaponLevelUp, w.InsID, oldLevel, w.Level)

		// 已经装备的武器升级后可能完成 装备一件X星X级武器 任务
		if w.EquipCharacterID != 0 {
			def := tableWeapon.GetID(w.ConfigID)
			this.userI.EmitEvent(event.EventWeaponEquipped, w.Level, def.RarityTypeEnum)
		}
	}

	this.setDirty(w.InsID)
}

func (this *UserWeapon) Refine(w *Weapon) {
	w.Refine++
	this.userI.EmitEvent(event.EventWeaponRefine, w.InsID, w.Refine-1, w.Refine)
	this.setDirty(w.InsID)
}

func (this *UserWeapon) Lock(w *Weapon, b bool) {
	if b != w.IsLock {
		w.IsLock = b
		this.setDirty(w.InsID)
	}
}

func (this *UserWeapon) ModuleType() module.ModuleType {
	return module.Weapon
}

func (this *UserWeapon) Tick(now time.Time) {}

func (this *UserWeapon) ReadOut() *module.ReadOutCommand {
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

func (this *UserWeapon) WriteCommand(fields map[interface{}]struct{}) *module.WriteBackCommand {
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
			slice := this.weapons[idx]
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

func (this *UserWeapon) FlushDirtyToClient() {
	if len(this.syncDirty) > 0 {
		msg := &message.WeaponSyncToC{
			IsAll:   proto.Bool(false),
			Weapons: make([]*message.Weapon, 0, len(this.syncDirty)),
			UseCap:  proto.Int32(int32(this.weaponUseCap)),
		}
		for id := range this.syncDirty {
			v := this.GetWeapon(id)
			if v != nil {
				msg.Weapons = append(msg.Weapons, v.Pack())
			} else { // 意识被删除
				msg.Weapons = append(msg.Weapons, &message.Weapon{
					ID:       proto.Uint32(id),
					IsRemove: proto.Bool(true),
				})
			}
		}
		this.syncDirty = map[uint32]struct{}{}
		this.userI.Post(msg)
	}
}

func (this *UserWeapon) FlushAllToClient(seqNo ...uint32) {
	msg1 := &message.WeaponSyncToC{
		IsAll:   proto.Bool(true),
		Weapons: make([]*message.Weapon, 0, 16),
		UseCap:  proto.Int32(int32(this.weaponUseCap)),
	}
	for _, sliceMap := range this.weapons {
		for _, v := range sliceMap {
			msg1.Weapons = append(msg1.Weapons, v.Pack())
		}
	}
	this.userI.Post(msg1)
	this.syncDirty = map[uint32]struct{}{}
}

func (this *UserWeapon) Init(fields map[string]*flyfish.Field) error {
	for i := 0; i < slotCount; i++ {
		fieldName := slotName(i)
		field, ok := fields[fieldName]
		this.weapons[i] = map[uint32]*Weapon{}
		if !ok || len(field.GetBlob()) == 0 {
			this.SetDirty(i)
		} else {
			if err := json.Unmarshal(field.GetBlob(), &this.weapons[i]); err != nil {
				return fmt.Errorf("unmarshal:%s %s", string(field.GetBlob()), err)
			}
			this.weaponUseCap += len(this.weapons[i])
		}
	}

	return nil
}

func (this *UserWeapon) Query(arg *message.QueryRoleInfoArg, ret *message.QueryRoleInfoResult) error {
	ids := arg.GetWeaponIDs()
	ret.Weapons = make([]*message.Weapon, 0, len(ids))
	for _, id := range ids {
		v := this.GetWeapon(id)
		if v != nil {
			ret.Weapons = append(ret.Weapons, v.Pack())
		}
	}
	return nil
}

func init() {
	module.RegisterModule(module.Weapon, func(userI module.UserI) module.ModuleI {
		m := &UserWeapon{
			userI:     userI,
			weapons:   make([]map[uint32]*Weapon, slotCount),
			syncDirty: map[uint32]struct{}{},
		}
		m.ModuleSaveBase = module.NewModuleSaveBase(m)
		return m
	})
}
