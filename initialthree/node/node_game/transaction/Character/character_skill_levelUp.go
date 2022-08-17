package Character

import (
	codecs "initialthree/codec/cs"
	"initialthree/node/common/attr"
	"initialthree/node/common/enumType"
	"initialthree/node/common/inoutput"
	"initialthree/node/common/transaction"
	"initialthree/node/table/excel/DataTable/PlayerSkill"
	"initialthree/zaplogger"

	"initialthree/node/node_game/module"
	"initialthree/node/node_game/module/character"
	"initialthree/node/node_game/user"
	"initialthree/protocol/cmdEnum"
	cs_message "initialthree/protocol/cs/message"
)

type transactionCharacterSkillLevelUp struct {
	transaction.TransactionBase
	user    *user.User
	req     *codecs.Message
	resp    *cs_message.CharacterSkillLevelUpToC
	errcode cs_message.ErrCode
}

func (this *transactionCharacterSkillLevelUp) GetModuleName() string {
	return "Character"
}

func (this *transactionCharacterSkillLevelUp) Begin() {
	defer func() { this.EndTrans(this.resp, this.errcode) }()
	this.resp = &cs_message.CharacterSkillLevelUpToC{}
	msg := this.req.GetData().(*cs_message.CharacterSkillLevelUpToS)

	zaplogger.GetSugar().Infof("%s %d Call CharacterSkillLevelUp %v", this.user.GetUserID(), this.user.GetID(), msg)

	userCharacter := this.user.GetSubModule(module.Character).(*character.UserCharacter)

	chara := userCharacter.GetCharacter(msg.GetCharacterID())
	skillDef := PlayerSkill.GetID(msg.GetSkillID())
	if chara == nil || skillDef == nil {
		zaplogger.GetSugar().Infof("%v %v %s", this.user.GetUserID(), this.user.GetID(), "CharacterSkillLevelUp fail:chara or skillDef is nil")
		this.errcode = cs_message.ErrCode_Config_NotExist
		return
	}

	var skill *cs_message.Skill
	for _, v := range chara.Skills {
		if v.GetSkillID() == msg.GetSkillID() {
			skill = v
			break
		}
	}

	if skill == nil {
		zaplogger.GetSugar().Infof("%v %v %s", this.user.GetUserID(), this.user.GetID(), "CharacterSkillLevelUp fail:skill id is failed")
		this.errcode = cs_message.ErrCode_Request_Argument_Err
		return
	}

	if skill.GetLevel() >= skillDef.GetMaxLevel() {
		zaplogger.GetSugar().Infof("%v %v %s", this.user.GetUserID(), this.user.GetID(), "CharacterSkillLevelUp fail:skill level max")
		this.errcode = cs_message.ErrCode_Level_High
		return
	}

	nextLevel := skillDef.Skill[skill.GetLevel()]
	if chara.Level < nextLevel.LimitLevel {
		zaplogger.GetSugar().Infof("%s CharacterSkillLevelUp fail:character level %d low %d", this.user.GetUserLogName(), chara.Level, nextLevel.LimitLevel)
		this.errcode = cs_message.ErrCode_Level_Low
		return
	}

	useRes := make([]inoutput.ResDesc, 0, len(nextLevel.CostItems()))
	if nextLevel.Gold > 0 {
		useRes = append(useRes, inoutput.ResDesc{ID: attr.Gold, Type: enumType.IOType_UsualAttribute, Count: nextLevel.Gold})
	}
	for _, v := range nextLevel.CostItems() {
		useRes = append(useRes, inoutput.ResDesc{ID: v.ItemID, Type: enumType.IOType_Item, Count: v.Count})
	}

	if this.errcode = this.user.DoInputOutput(useRes, nil, false); this.errcode != cs_message.ErrCode_OK {
		zaplogger.GetSugar().Debugf("%s CharacterSkillLevelUp error: inouput %s  ", this.user.GetUserLogName(), this.errcode.String())
		return
	}

	userCharacter.SkillLevelUp(chara, skill)

	zaplogger.GetSugar().Infof("%v %v %s", this.user.GetUserID(), this.user.GetID(), "CharacterSkillLevelUp OK")
	this.errcode = cs_message.ErrCode_OK
}

func init() {
	user.RegisterTransFunc(cmdEnum.CS_CharacterSkillLevelUp, func(user *user.User, msg *codecs.Message) transaction.Transaction {
		return &transactionCharacterSkillLevelUp{
			user: user,
			req:  msg,
		}
	})
}
