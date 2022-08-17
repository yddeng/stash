package Weapon

import (
	codecs "initialthree/codec/cs"
	"initialthree/node/common/attr"
	"initialthree/node/common/enumType"
	"initialthree/node/common/inoutput"
	"initialthree/node/common/transaction"
	"initialthree/node/node_game/module"
	"initialthree/node/node_game/module/weapon"
	"initialthree/node/node_game/user"
	"initialthree/node/table/excel/DataTable/Skill"
	"initialthree/node/table/excel/DataTable/Weapon"
	"initialthree/node/table/excel/DataTable/WeaponRarity"
	"initialthree/protocol/cmdEnum"
	cs_message "initialthree/protocol/cs/message"
	"initialthree/zaplogger"
)

type transactionWeaponRefine struct {
	transaction.TransactionBase
	user    *user.User
	req     *codecs.Message
	resp    *cs_message.WeaponRefineToC
	errcode cs_message.ErrCode
}

func (this *transactionWeaponRefine) GetModuleName() string {
	return "Weapon"
}

func (this *transactionWeaponRefine) Begin() {
	defer func() { this.EndTrans(this.resp, this.errcode) }()
	this.resp = &cs_message.WeaponRefineToC{}
	msg := this.req.GetData().(*cs_message.WeaponRefineToS)

	zaplogger.GetSugar().Infof("%s Call WeaponRefineToS %v", this.user.GetUserLogName(), msg)

	userWeapon := this.user.GetSubModule(module.Weapon).(*weapon.UserWeapon)

	useWeapon := userWeapon.GetWeapon(msg.GetCostWeapons())
	doWeapon := userWeapon.GetWeapon(msg.GetWeaponID())
	if doWeapon == nil || useWeapon == nil {
		zaplogger.GetSugar().Debugf("%s WeaponRefineToS error: weapon is nil", this.user.GetUserLogName())
		this.errcode = cs_message.ErrCode_Weapon_NotExist
		return
	}

	if doWeapon.ConfigID != useWeapon.ConfigID {
		zaplogger.GetSugar().Debugf("%s WeaponRefineToS error: %d config id is failed", this.user.GetUserLogName(), doWeapon.ConfigID)
		this.errcode = cs_message.ErrCode_Weapon_TypeErr
		return
	}

	def := Weapon.GetID(doWeapon.ConfigID)
	if def == nil {
		zaplogger.GetSugar().Debugf("%s WeaponRefineToS error: %d weaponDef is nil", this.user.GetUserLogName(), doWeapon.ConfigID)
		this.errcode = cs_message.ErrCode_Config_NotExist
		return
	}

	skillDef := Skill.GetID(def.WeaponSkillConfig)
	rarityDef := WeaponRarity.GetID(def.RarityTypeEnum)
	if skillDef == nil || rarityDef == nil {
		zaplogger.GetSugar().Debugf("%s WeaponRefineToS error: %d skill or rarity def is nil", this.user.GetUserLogName(), def.WeaponSkillConfig)
		this.errcode = cs_message.ErrCode_Config_NotExist
		return
	}

	if doWeapon.Refine >= skillDef.GetMaxLevel() {
		zaplogger.GetSugar().Debugf("%s WeaponRefineToS error: %d skill level is max", this.user.GetUserLogName(), doWeapon.InsID)
		this.errcode = cs_message.ErrCode_Weapon_LevelMax
		return
	}

	useRes := []inoutput.ResDesc{
		{Type: enumType.IOType_Weapon, ID: int32(msg.GetCostWeapons()), Count: 1},
		{Type: enumType.IOType_UsualAttribute, ID: attr.Gold, Count: rarityDef.RefineCostGold},
	}

	// 消耗
	if this.errcode = this.user.DoInputOutput(useRes, nil, false); this.errcode != cs_message.ErrCode_OK {
		zaplogger.GetSugar().Debugf("%s WeaponRefineToS error: inouput %s  ", this.user.GetUserLogName(), this.errcode.String())
		return
	}

	userWeapon.Refine(doWeapon)

	this.errcode = cs_message.ErrCode_OK
	zaplogger.GetSugar().Infof("%s Call WeaponRefineToS ok", this.user.GetUserLogName())
}

func init() {
	user.RegisterTransFunc(cmdEnum.CS_WeaponRefine, func(user *user.User, msg *codecs.Message) transaction.Transaction {
		return &transactionWeaponRefine{
			user: user,
			req:  msg,
		}
	})
}
