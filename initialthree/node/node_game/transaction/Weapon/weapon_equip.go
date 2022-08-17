package Weapon

import (
	codecs "initialthree/codec/cs"
	"initialthree/node/common/transaction"
	"initialthree/node/node_game/event"
	"initialthree/node/node_game/module/weapon"
	"initialthree/node/table/excel/DataTable/Weapon"
	"initialthree/zaplogger"

	"initialthree/node/node_game/module"
	"initialthree/node/node_game/module/character"
	"initialthree/node/node_game/user"
	"initialthree/protocol/cmdEnum"
	cs_message "initialthree/protocol/cs/message"
)

type transactionWeaponEquip struct {
	transaction.TransactionBase
	user    *user.User
	req     *codecs.Message
	resp    *cs_message.WeaponEquipToC
	errcode cs_message.ErrCode
}

func (this *transactionWeaponEquip) GetModuleName() string {
	return "Equip"
}

func (this *transactionWeaponEquip) Begin() {
	defer func() { this.EndTrans(this.resp, this.errcode) }()
	this.resp = &cs_message.WeaponEquipToC{}
	msg := this.req.GetData().(*cs_message.WeaponEquipToS)

	zaplogger.GetSugar().Infof("%s Call WeaponEquipToS %v", this.user.GetUserLogName(), msg)

	userCharacter := this.user.GetSubModule(module.Character).(*character.UserCharacter)
	userWeapon := this.user.GetSubModule(module.Weapon).(*weapon.UserWeapon)

	chara := userCharacter.GetCharacter(msg.GetCharacterID())
	doWeapon := userWeapon.GetWeapon(msg.GetWeaponID())
	if chara == nil || doWeapon == nil {
		zaplogger.GetSugar().Debugf("%s WeaponEquipToS error: chara or weapon is nil", this.user.GetUserLogName())
		this.errcode = cs_message.ErrCode_Weapon_NotExist
		return
	}

	def := Weapon.GetID(doWeapon.ConfigID)
	if def == nil {
		zaplogger.GetSugar().Debugf("%s WeaponEquipToS error: %d weapon def is nil", this.user.GetUserLogName(), doWeapon.ConfigID)
		this.errcode = cs_message.ErrCode_Config_NotExist
		return
	}

	oldWeapon := userWeapon.GetWeapon(chara.Weapon)

	// 武器类型
	if def.WeaponTypeEnum != Weapon.GetID(oldWeapon.ConfigID).WeaponTypeEnum {
		zaplogger.GetSugar().Debugf("%s WeaponEquipToS error: %d equip type failed", this.user.GetUserLogName(), doWeapon.ConfigID)
		this.errcode = cs_message.ErrCode_Weapon_TypeErr
		return
	}

	// 待穿戴的装备是否已经在其他角色身上
	if doWeapon.EquipCharacterID != 0 {
		oldChara := userCharacter.GetCharacter(doWeapon.EquipCharacterID)
		userCharacter.WeaponReplace(oldChara, oldWeapon.InsID)
		userCharacter.WeaponReplace(chara, msg.GetWeaponID())

		userWeapon.Equip(oldWeapon, oldChara.CharacterID)
		userWeapon.Equip(doWeapon, chara.CharacterID)
	} else {
		userCharacter.WeaponReplace(chara, msg.GetWeaponID())
		userWeapon.Equip(doWeapon, chara.CharacterID)
		userWeapon.Demount(oldWeapon)
	}

	this.user.EmitEvent(event.EventWeaponEquipped, doWeapon.Level, def.RarityTypeEnum)

	this.errcode = cs_message.ErrCode_OK
	zaplogger.GetSugar().Infof("%s Call WeaponEquipToS ok", this.user.GetUserLogName())
}

func init() {
	user.RegisterTransFunc(cmdEnum.CS_WeaponEquip, func(user *user.User, msg *codecs.Message) transaction.Transaction {
		return &transactionWeaponEquip{
			user: user,
			req:  msg,
		}
	})
}
