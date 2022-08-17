package character

import (
	"initialthree/node/table/excel/DataTable/PlayerSkill"
	vendor_event "initialthree/pkg/event"
	"initialthree/protocol/cs/message"
)

type defaultUnlockSkill struct {
	m     *UserCharacter
	c     *Character
	skill *message.Skill
	h     vendor_event.Handle
}

// 默认技能解锁, 角色升级
func (this *defaultUnlockSkill) EventCharacterLevelUp(charaID, oldLevel, newLevel int32) {
	skillDef := PlayerSkill.GetID(this.skill.GetSkillID())
	cfg, _ := skillDef.GetDefaultUnlockSkillCond()
	if this.skill.GetLevel() == 0 && newLevel >= cfg.LimitLevel && this.c.GeneLevel >= skillDef.RequiredGeneLevel {
		this.m.SkillLevelUp(this.c, this.skill)
		this.m.userI.UnRegisterEvent(this.h)
	}
}

type geneLevelUnlockSkill struct {
	m     *UserCharacter
	c     *Character
	skill *message.Skill
	h     vendor_event.Handle
}

// 默认技能解锁，角色命座升级
func (this *geneLevelUnlockSkill) EventCharacterGeneLevelUp(charaID, oldLevel, newLevel int32) {
	skillDef := PlayerSkill.GetID(this.skill.GetSkillID())
	cfg, _ := skillDef.GetDefaultUnlockSkillCond()
	if this.skill.GetLevel() == 0 && newLevel >= skillDef.RequiredGeneLevel && this.c.Level > cfg.LimitLevel {
		this.m.SkillLevelUp(this.c, this.skill)
		this.m.userI.UnRegisterEvent(this.h)
	}
}
