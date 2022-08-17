package Equip

import (
	codecs "initialthree/codec/cs"
	"initialthree/node/common/attr"
	"initialthree/node/common/enumType"
	"initialthree/node/common/inoutput"
	"initialthree/node/common/transaction"
	"initialthree/node/node_game/module"
	"initialthree/node/node_game/module/equip"
	"initialthree/node/node_game/user"
	"initialthree/node/table/excel/DataTable/Equip"
	"initialthree/node/table/excel/DataTable/EquipAttribute"
	"initialthree/node/table/excel/DataTable/EquipQuality"
	"initialthree/node/table/excel/DataTable/Skill"
	"initialthree/protocol/cmdEnum"
	cs_message "initialthree/protocol/cs/message"
	"initialthree/zaplogger"
)

type transactionEquipRefine struct {
	transaction.TransactionBase
	user    *user.User
	req     *codecs.Message
	resp    *cs_message.EquipRefineToC
	errcode cs_message.ErrCode
}

func (this *transactionEquipRefine) GetModuleName() string {
	return "Equip"
}

func (this *transactionEquipRefine) Begin() {
	defer func() { this.EndTrans(this.resp, this.errcode) }()
	this.resp = &cs_message.EquipRefineToC{}
	msg := this.req.GetData().(*cs_message.EquipRefineToS)

	zaplogger.GetSugar().Infof("%s Call EquipRefineToS %v", this.user.GetUserLogName(), msg)

	userEquip := this.user.GetSubModule(module.Equip).(*equip.UserEquip)

	doEquip := userEquip.GetEquip(msg.GetEquipID())
	if doEquip == nil {
		zaplogger.GetSugar().Debugf("%s EquipRefineToS error: equip is nil", this.user.GetUserLogName())
		this.errcode = cs_message.ErrCode_Equip_NotExist
		return
	}

	def := Equip.GetID(doEquip.ConfigID)
	if def == nil {
		zaplogger.GetSugar().Debugf("%s EquipRefineToS error: %d equipDef is nil", this.user.GetUserLogName(), doEquip.ConfigID)
		this.errcode = cs_message.ErrCode_Config_NotExist
		return
	}

	if msg.GetPosition() < 0 || msg.GetPosition() > int32(len(doEquip.Refine)) {
		zaplogger.GetSugar().Debugf("%s EquipRefineToS error: position %d is failed", this.user.GetUserLogName(), msg.GetPosition())
		this.errcode = cs_message.ErrCode_Request_Argument_Err
		return
	}

	rarityDef := EquipQuality.GetID(def.QualityEnum)
	if rarityDef == nil {
		zaplogger.GetSugar().Debugf("%s EquipRefineToS error: rarity def is nil", this.user.GetUserLogName())
		this.errcode = cs_message.ErrCode_Config_NotExist
		return
	}

	maxLevel := int32(0)
	if msg.GetPosition() == 0 {
		skillDef := Skill.GetID(def.SkillID)
		if skillDef == nil {
			zaplogger.GetSugar().Debugf("%s EquipRefineToS error: %d skill def is nil", this.user.GetUserLogName(), def.SkillID)
			this.errcode = cs_message.ErrCode_Config_NotExist
			return
		}

		maxLevel = int32(len(skillDef.Damage)) - 1
	} else {
		attrbDef := EquipAttribute.GetID(doEquip.RandomAttribId)
		if attrbDef == nil {
			zaplogger.GetSugar().Debugf("%s EquipRefineToS error: %d attribute def is nil", this.user.GetUserLogName(), doEquip.RandomAttribId)
			this.errcode = cs_message.ErrCode_Config_NotExist
			return
		}

		maxLevel = int32(len(attrbDef.Attr)) - 1
	}

	if doEquip.Refine[msg.GetPosition()] >= maxLevel {
		zaplogger.GetSugar().Debugf("%s EquipRefineToS error: %d skill level is max", this.user.GetUserLogName(), def.SkillID)
		this.errcode = cs_message.ErrCode_Equip_LevelMax
		return
	}

	useRes := make([]inoutput.ResDesc, 0, len(msg.GetCostEquips())+1)
	totalLev := int32(0)
	for _, id := range msg.GetCostEquips() {
		mm := userEquip.GetEquip(id)
		if mm == nil {
			zaplogger.GetSugar().Debugf("%s EquipRefineToS error: cost equip %d is nil", this.user.GetUserLogName(), id)
			this.errcode = cs_message.ErrCode_Equip_NotExist
			return
		}
		ddef := Equip.GetID(mm.ConfigID)
		if ddef == nil {
			zaplogger.GetSugar().Debugf("%s EquipRefineToS error:cost equip %d def is nil", this.user.GetUserLogName(), mm.ConfigID)
			this.errcode = cs_message.ErrCode_Config_NotExist
			return
		}
		if ddef.QualityEnum < def.QualityEnum {
			zaplogger.GetSugar().Debugf("%s EquipRefineToS error:cost equip rarity is too low", this.user.GetUserLogName())
			this.errcode = cs_message.ErrCode_Equip_RarityLow
			return
		}
		// 本身消耗增加1级
		totalLev += 1

		// 已经精炼过加
		for _, lev := range mm.Refine {
			totalLev += lev
		}
		useRes = append(useRes, inoutput.ResDesc{Type: enumType.IOType_Equip, ID: int32(id), Count: 1})
	}

	useRes = append(useRes, inoutput.ResDesc{Type: enumType.IOType_UsualAttribute, ID: attr.Gold, Count: rarityDef.RefineCostGold})

	// 消耗
	if this.errcode = this.user.DoInputOutput(useRes, nil, false); this.errcode != cs_message.ErrCode_OK {
		zaplogger.GetSugar().Debugf("%s EquipRefineToS error: inouput %s  ", this.user.GetUserLogName(), this.errcode.String())
		return
	}

	userEquip.Refine(doEquip, msg.GetPosition(), totalLev, maxLevel)

	this.errcode = cs_message.ErrCode_OK
	zaplogger.GetSugar().Infof("%s Call EquipRefineToS ok", this.user.GetUserLogName())
}

func init() {
	user.RegisterTransFunc(cmdEnum.CS_EquipRefine, func(user *user.User, msg *codecs.Message) transaction.Transaction {
		return &transactionEquipRefine{
			user: user,
			req:  msg,
		}
	})
}
