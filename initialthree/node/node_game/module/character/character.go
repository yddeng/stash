package character

import (
	"fmt"
	"initialthree/node/node_game/event"
	"initialthree/node/node_game/module"
	"initialthree/node/table/excel/ConstTable/Global"
	"initialthree/node/table/excel/DataTable/CharacterBreakThrough"
	"initialthree/node/table/excel/DataTable/CharacterLevelUpExp"
	table_Character "initialthree/node/table/excel/DataTable/PlayerCharacter"
	"initialthree/node/table/excel/DataTable/PlayerSkill"
	"initialthree/pkg/json"
	cs_message "initialthree/protocol/cs/message"
	"math"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/sniperHW/flyfish/client"
)

const (
	dbGroupDefault = "group_default"

	/*
	   编队预设
	   0 1人队
	   1 2人队
	   2 3人队
	   3 第一个预设队伍
	*/
	groupDefLen = 3 + 6 // 编队预设容量大小

	equipCount = 5 // 5件装备 0号位为主位置

	slotCount = 10
)

func slotName(slotIdx int) string {
	return fmt.Sprintf("slot%d", slotIdx)
}

func calcSlotIdx(id int) int {
	return id % slotCount
}

// 角色
type Character struct {
	CharacterID int32               `json:"cid"`      // ID
	Level       int32               `json:"l"`        // 等级
	CurrentExp  int32               `json:"ce"`       // 当前经验
	BreakLevel  int32               `json:"bl"`       // 突破等级
	Weapon      uint32              `json:"weapon"`   // 武器
	EquipIDs    []uint32            `json:"e_id_s"`   // 装备
	GetTime     int64               `json:"get_time"` // 获取时间
	FavorLevel  int32               `json:"fl"`       // 好感度等级
	FavorExp    int32               `json:"fe"`       // 好感度经验
	IsLike      bool                `json:"is_like"`
	Skills      []*cs_message.Skill `json:"skills"`
	GeneLevel   int32               `json:"gene_level"`
	HitTimes    int32               `json:"hit_times"` // 角色的获取次数
}

func (this *Character) GetEquipIDs() []uint32 {
	return this.EquipIDs
}

func (this *Character) Pack() *cs_message.Character {
	return &cs_message.Character{
		CharacterID: proto.Int32(this.CharacterID),
		Level:       proto.Int32(this.Level),
		CurrentExp:  proto.Int32(this.CurrentExp),
		WeaponID:    proto.Uint32(this.Weapon),
		BreakLevel:  proto.Int32(this.BreakLevel),
		EquipIDs:    this.EquipIDs,
		GetTime:     proto.Int64(this.GetTime),
		FavorLevel:  proto.Int32(this.FavorLevel),
		FavorExp:    proto.Int32(this.FavorExp),
		GeneLevel:   proto.Int32(this.GeneLevel),
		IsLike:      proto.Bool(this.IsLike),
		Skills:      this.Skills,
	}
}

type UserCharacter struct {
	userI module.UserI
	// data
	data     []map[int32]*Character // 10 个分片
	groupDef []*cs_message.CharacterTeamPrefab

	charaDirty    map[int32]struct{}
	groupDefDirty bool
	*module.ModuleSaveBase
}

func (this *UserCharacter) NewCharacter(id int32, def *table_Character.PlayerCharacter) *Character {
	c := &Character{
		CharacterID: id,
		Level:       1,
		EquipIDs:    make([]uint32, equipCount),
		GetTime:     time.Now().Unix(),
		FavorLevel:  1,
		Skills:      make([]*cs_message.Skill, 0, len(def.PlayeSkillsArray)),
	}

	for _, v := range def.PlayeSkillsArray {
		skillDef := PlayerSkill.GetID(v.ID)
		if skillDef != nil {
			skill := &cs_message.Skill{
				SkillID: proto.Int32(v.ID),
				Level:   proto.Int32(0),
			}
			def, ok := skillDef.GetDefaultUnlockSkillCond()
			if ok && c.Level >= def.LimitLevel && c.GeneLevel >= skillDef.RequiredGeneLevel {
				skill.Level = proto.Int32(1)
			}
			c.Skills = append(c.Skills, skill)
		}
	}

	return c
}

func (this *UserCharacter) setDirty(id int32) {
	slotIdx := calcSlotIdx(int(id))

	this.charaDirty[id] = struct{}{}
	this.SetDirty(slotIdx)
}

func (this *UserCharacter) Range(cb func(c *Character) bool) {
	for _, slot := range this.data {
		for _, c := range slot {
			if !cb(c) {
				return
			}
		}
	}
}

func (this *UserCharacter) AddCharacter(c *Character) {
	slotIdx := calcSlotIdx(int(c.CharacterID))
	slot := this.data[slotIdx]
	if _, ok := slot[c.CharacterID]; ok {
		return
	} else {
		slot[c.CharacterID] = c
		this.setDirty(c.CharacterID)
	}
}

func (this *UserCharacter) RemoveCharacter(cId int32) {
	slotIdx := calcSlotIdx(int(cId))
	slot := this.data[slotIdx]
	if _, ok := slot[cId]; !ok {
		return
	}

	delete(slot, cId)
	this.setDirty(cId)

}

func (this *UserCharacter) GetCharacter(id int32) *Character {
	slotIdx := calcSlotIdx(int(id))
	return this.data[slotIdx][id]
}

func (this *UserCharacter) AddHitTimes(c *Character, dt int32) {
	c.HitTimes += dt
	this.setDirty(c.CharacterID)
}

func (this *UserCharacter) IsLike(c *Character, b bool) {
	if b != c.IsLike {
		c.IsLike = b
		this.setDirty(c.CharacterID)
	}
}

func (this *UserCharacter) WeaponReplace(c *Character, cInsID uint32) {
	c.Weapon = cInsID
	this.setDirty(c.CharacterID)
}

func (this *UserCharacter) EquipEquip(c *Character, i int32, cInsID uint32) {
	c.EquipIDs[i] = cInsID
	this.setDirty(c.CharacterID)
}

func (this *UserCharacter) EquipDemount(c *Character, cInsID uint32) {
	for i, id := range c.EquipIDs {
		if id == cInsID {
			c.EquipIDs[i] = 0
			this.setDirty(c.CharacterID)
			break
		}
	}

}

func (this *UserCharacter) SkillLevelUp(c *Character, skill *cs_message.Skill) {
	skill.Level = proto.Int32(skill.GetLevel() + 1)
	this.setDirty(c.CharacterID)
	this.userI.EmitEvent(event.EventCharacterSkillLevelUp, c.CharacterID, skill.GetSkillID(), skill.GetLevel()-1, skill.GetLevel())
}

func (this *UserCharacter) GeneLevelUp(c *Character) {
	c.GeneLevel++
	this.userI.EmitEvent(event.EventCharacterGeneLevelUp, c.CharacterID, c.GeneLevel-1, c.GeneLevel)
	this.setDirty(c.CharacterID)
}

func (this *UserCharacter) GetMaxLevel(character *Character) int32 {
	charaDef := table_Character.GetID(character.CharacterID)
	maxLevel := int32(math.MaxInt32)
	breakDef := CharacterBreakThrough.GetID(charaDef.GetBreakID(character.BreakLevel + 1))
	if breakDef != nil && breakDef.LevelRequirement < maxLevel {
		maxLevel = breakDef.LevelRequirement
	}
	return maxLevel
}

// 角色升级实际消耗的经验
func (this *UserCharacter) CalcUseExp(character *Character, exp, maxLevel int32) (usedExp, level, currentExp int32) {
	rarity := table_Character.GetID(character.CharacterID).RarityEnum
	currentExp = character.CurrentExp
	level = character.Level
	usedExp = exp
	currentExp += exp

	needExp, exist := CharacterLevelUpExp.GetMaxExp(rarity, level+1)
	for exist && level < maxLevel {
		if needExp != 0 && currentExp >= needExp {
			currentExp -= needExp
			level += 1
			needExp, exist = CharacterLevelUpExp.GetMaxExp(rarity, level+1)
		} else {
			break
		}
	}

	// 超过下一等级所需的经验，但是玩家等级限制，舍弃多余经验保证小于下级所需经验
	if !exist || level >= maxLevel {
		usedExp -= currentExp
		currentExp = 0
	}
	return
}

func (this *UserCharacter) SetLevel(character *Character, level, currentExp int32) {
	oldLevel := character.Level
	character.Level = level
	character.CurrentExp = currentExp

	this.setDirty(character.CharacterID)
	// 事件触发
	if oldLevel != character.Level {
		this.userI.EmitEvent(event.EventCharacterLevelUp, character.CharacterID, oldLevel, character.Level)
	}
}

func (this *UserCharacter) FavorLevelUp(character *Character, exp int32) {
	character.FavorExp += exp
	favorExp := Global.Get().FavorExpArray

	for character.FavorLevel-1 < int32(len(favorExp)) {
		needExp := favorExp[character.FavorLevel-1].Exp
		if needExp != 0 && character.FavorExp >= needExp {
			character.FavorExp -= needExp
			character.FavorLevel += 1
		} else {
			break
		}
	}

	if character.FavorLevel-1 == int32(len(favorExp)) {
		character.FavorExp = 0
	}

	this.setDirty(character.CharacterID)

}

func (this *UserCharacter) BreakLevel(c *Character) {
	c.BreakLevel++
	this.userI.EmitEvent(event.EventCharacterBreak, c.CharacterID, c.BreakLevel-1, c.BreakLevel)
	this.setDirty(c.CharacterID)
}

func (this *UserCharacter) GroupDefSet(idx int, prefab *cs_message.CharacterTeamPrefab) {
	this.groupDef[idx] = prefab
	this.groupDefDirty = true
	this.SetDirty(dbGroupDefault)

}

func (this *UserCharacter) ModuleType() module.ModuleType {
	return module.Character
}

func (this *UserCharacter) Init(fields map[string]*client.Field) error {

	for i := 0; i < slotCount; i++ {
		dbname := slotName(i)
		field, ok := fields[dbname]
		this.data[i] = map[int32]*Character{}
		if !ok || len(field.GetBlob()) == 0 {
			this.SetDirty(i)
		} else {
			if err := json.Unmarshal(field.GetBlob(), &this.data[i]); err != nil {
				return fmt.Errorf("unmarshal: %s", err)
			}
		}
	}

	field, ok := fields[dbGroupDefault]
	if ok && len(field.GetBlob()) != 0 {
		if err := json.Unmarshal(field.GetBlob(), &this.groupDef); err != nil {
			return fmt.Errorf("unmarshal:%s %s", string(field.GetBlob()), err)
		}
	} else {
		for i := 3; i < len(this.groupDef); i++ {
			this.groupDef[i] = &cs_message.CharacterTeamPrefab{
				Name: proto.String(fmt.Sprintf("编队%d", i-2)),
			}
		}
		this.SetDirty(dbGroupDefault)
	}

	return nil
}

func (this *UserCharacter) Tick(now time.Time) {
}

func (this *UserCharacter) ReadOut() *module.ReadOutCommand {
	out := &module.ReadOutCommand{
		Table:  this.ModuleType().String(),
		Key:    this.userI.GetIDStr(),
		Module: this,
		Fields: make([]string, 0, slotCount+1),
	}
	for i := 0; i < slotCount; i++ {
		out.Fields = append(out.Fields, slotName(i))
	}
	out.Fields = append(out.Fields, dbGroupDefault)

	return out
}

func (this *UserCharacter) WriteCommand(fields map[interface{}]struct{}) *module.WriteBackCommand {
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
			data, _ := json.Marshal(this.data[idx])
			cmd.Fields = append(cmd.Fields, &module.WriteBackFiled{
				Name:  slotName(idx),
				Value: data,
			})
		case string:
			name := field.(string)
			if name == dbGroupDefault {
				data, _ := json.Marshal(this.groupDef)
				cmd.Fields = append(cmd.Fields, &module.WriteBackFiled{
					Name:  dbGroupDefault,
					Value: data,
				})
			}
		}
	}

	return cmd
}

func (this *UserCharacter) FlushAllToClient(seqNo ...uint32) {
	msg1 := &cs_message.CharacterSyncToC{
		IsAll:      proto.Bool(true),
		Characters: make([]*cs_message.Character, 0, 16),
	}
	for _, sliceMap := range this.data {
		for _, v := range sliceMap {
			msg1.Characters = append(msg1.Characters, v.Pack())
		}
	}
	this.charaDirty = map[int32]struct{}{}
	this.userI.Post(msg1)

	msg3 := &cs_message.CharacterTeamPrefabSyncToC{
		CharacterTeamPrefabs: this.groupDef,
	}
	this.groupDefDirty = false
	this.userI.Post(msg3)

}

func (this *UserCharacter) FlushDirtyToClient() {
	if len(this.charaDirty) > 0 {
		msg := &cs_message.CharacterSyncToC{
			IsAll:      proto.Bool(false),
			Characters: make([]*cs_message.Character, 0, len(this.charaDirty)),
		}
		for id := range this.charaDirty {
			v := this.GetCharacter(id)
			if v != nil {
				msg.Characters = append(msg.Characters, v.Pack())
			}
		}
		this.charaDirty = map[int32]struct{}{}
		this.userI.Post(msg)
	}

	if this.groupDefDirty {
		msg := &cs_message.CharacterTeamPrefabSyncToC{
			CharacterTeamPrefabs: this.groupDef,
		}
		this.groupDefDirty = false
		this.userI.Post(msg)
	}

}

// 默认解锁技能注册
func (this *UserCharacter) AfterInitAll() error {
	for _, s := range this.data {
		for _, c := range s {
			for _, skill := range c.Skills {
				skillDef := PlayerSkill.GetID(skill.GetSkillID())
				def, defaultUnlock := skillDef.GetDefaultUnlockSkillCond()
				if skill.GetLevel() == 0 && defaultUnlock {
					if c.Level >= def.LimitLevel && c.GeneLevel >= skillDef.RequiredGeneLevel {
						this.SkillLevelUp(c, skill)
					} else {
						if c.Level < def.LimitLevel {
							cond := &defaultUnlockSkill{
								m:     this,
								c:     c,
								skill: skill,
							}
							cond.h = this.userI.RegisterEvent(event.EventCharacterLevelUp, cond.EventCharacterLevelUp)
						}
						if c.GeneLevel >= skillDef.RequiredGeneLevel {
							cond := &geneLevelUnlockSkill{
								m:     this,
								c:     c,
								skill: skill,
							}
							cond.h = this.userI.RegisterEvent(event.EventCharacterGeneLevelUp, cond.EventCharacterGeneLevelUp)
						}
					}
				}
			}
		}
	}
	return nil
}

func (this *UserCharacter) Query(arg *cs_message.QueryRoleInfoArg, ret *cs_message.QueryRoleInfoResult) error {
	ids := arg.GetCharacterIDs()
	ret.Characters = make([]*cs_message.Character, 0, len(ids))
	for _, id := range ids {
		v := this.GetCharacter(id)
		if v != nil {
			ret.Characters = append(ret.Characters, v.Pack())

			// 武器
			existWeapon := false
			for _, wID := range arg.GetWeaponIDs() {
				if wID == v.Weapon {
					existWeapon = true
					break
				}
			}
			if !existWeapon {
				arg.WeaponIDs = append(arg.WeaponIDs, v.Weapon)
			}

			// 装备
			for _, id := range v.EquipIDs {
				existEquip := false
				for _, eID := range arg.GetEquipIDs() {
					if eID == id {
						existEquip = true
						break
					}
				}
				if !existEquip {
					arg.EquipIDs = append(arg.EquipIDs, id)
				}
			}
		}
	}
	return nil
}

func init() {
	module.RegisterModule(module.Character, func(userI module.UserI) module.ModuleI {
		characters := &UserCharacter{
			userI:      userI,
			data:       make([]map[int32]*Character, slotCount),
			groupDef:   make([]*cs_message.CharacterTeamPrefab, groupDefLen),
			charaDirty: map[int32]struct{}{},
		}

		characters.ModuleSaveBase = module.NewModuleSaveBase(characters)

		return characters
	})

}
