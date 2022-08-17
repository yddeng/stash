package Weapon

import (
	codecs "initialthree/codec/cs"
	"initialthree/node/common/attr"
	"initialthree/node/common/enumType"
	"initialthree/node/common/inoutput"
	"initialthree/node/common/transaction"
	"initialthree/node/node_game/module"
	"initialthree/node/node_game/module/backpack"
	"initialthree/node/node_game/module/weapon"
	"initialthree/node/node_game/user"
	constWeapon "initialthree/node/table/excel/ConstTable/Weapon"
	"initialthree/node/table/excel/DataTable/Item"
	"initialthree/node/table/excel/DataTable/Weapon"
	"initialthree/node/table/excel/DataTable/WeaponBreakThough"
	"initialthree/node/table/excel/DataTable/WeaponLevelMaxExp"
	"initialthree/node/table/excel/DataTable/WeaponRarity"
	"initialthree/protocol/cmdEnum"
	cs_message "initialthree/protocol/cs/message"
	"initialthree/zaplogger"
)

type transactionWeaponStrengthen struct {
	transaction.TransactionBase
	user    *user.User
	req     *codecs.Message
	resp    *cs_message.WeaponStrengthenToC
	errcode cs_message.ErrCode
}

func (this *transactionWeaponStrengthen) GetModuleName() string {
	return "Weapon"
}

func (this *transactionWeaponStrengthen) Begin() {
	defer func() { this.EndTrans(this.resp, this.errcode) }()
	this.resp = &cs_message.WeaponStrengthenToC{}
	msg := this.req.GetData().(*cs_message.WeaponStrengthenToS)

	zaplogger.GetSugar().Infof("%s Call WeaponStrengthenToS %v", this.user.GetUserLogName(), msg)

	userWeapon := this.user.GetSubModule(module.Weapon).(*weapon.UserWeapon)
	doWeapon := userWeapon.GetWeapon(msg.GetWeaponID())
	if doWeapon == nil {
		zaplogger.GetSugar().Debugf("%s WeaponStrengthenToS error: mod is nil", this.user.GetUserLogName())
		this.errcode = cs_message.ErrCode_Weapon_NotExist
		return
	}

	def := Weapon.GetID(doWeapon.ConfigID)
	if def == nil {
		zaplogger.GetSugar().Debugf("%s WeaponStrengthenToS error: %d modDef is nil", this.user.GetUserLogName(), doWeapon.ConfigID)
		this.errcode = cs_message.ErrCode_Config_NotExist
		return
	}

	// 请求参数验证
	if len(msg.GetCostItems()) == 0 && len(msg.GetCostWeaponIDs()) == 0 {
		zaplogger.GetSugar().Debugf("%s WeaponStrengthenToS error: costWeaponIDs is nil", this.user.GetUserLogName())
		this.errcode = cs_message.ErrCode_Request_Argument_Err
		return
	}

	levelDef := WeaponLevelMaxExp.GetID(def.LevelExpConfig)
	if levelDef == nil {
		zaplogger.GetSugar().Debugf("%s WeaponStrengthenToS error: cost %d levelDef is nil", this.user.GetUserLogName(), def.LevelExpConfig)
		this.errcode = cs_message.ErrCode_Config_NotExist
		return
	}

	breakDef := WeaponBreakThough.GetID(def.BreakThroughConfig)
	rarityDef := WeaponRarity.GetID(def.RarityTypeEnum)
	if breakDef == nil || rarityDef == nil {
		zaplogger.GetSugar().Debugf("%s WeaponBreakToS error: %d breakDef or rarityDef is nil", this.user.GetUserLogName(), msg.GetWeaponID())
		this.errcode = cs_message.ErrCode_Config_NotExist
		return
	}

	maxLevel := rarityDef.MaxLevel
	if levelDef.MaxLevel() < maxLevel {
		maxLevel = levelDef.MaxLevel()
	}
	breakLevel := breakDef.GetBreakTimesLevel(doWeapon.BreakTimes + 1)
	if breakLevel != nil && breakLevel.Level < maxLevel {
		maxLevel = breakLevel.Level
	}

	// 等级到达上限
	if doWeapon.Level >= maxLevel {
		zaplogger.GetSugar().Debugf("%s WeaponStrengthenToS error: %d level %d max ", this.user.GetUserLogName(), msg.GetWeaponID(), maxLevel)
		this.errcode = cs_message.ErrCode_Weapon_LevelMax
		return
	}

	userBackpack := this.user.GetSubModule(module.Backpack).(*backpack.Backpack)
	totalExp := int32(0)
	var useRes = make([]inoutput.ResDesc, 0, len(msg.GetCostWeaponIDs())+len(msg.GetCostItems())+1)

	for _, id := range msg.GetCostWeaponIDs() {
		if id == msg.GetWeaponID() {
			zaplogger.GetSugar().Debugf("%s WeaponStrengthenToS error: the consumption is the same as the operation %d ", this.user.GetUserLogName(), msg.GetWeaponID())
			this.errcode = cs_message.ErrCode_Request_Argument_Err
			return
		}
		useWeapon := userWeapon.GetWeapon(id)
		if useWeapon == nil {
			zaplogger.GetSugar().Debugf("%s WeaponStrengthenToS error: cost mod %d is nil", this.user.GetUserLogName(), id)
			this.errcode = cs_message.ErrCode_Weapon_NotExist
			return
		}
		// 装备已被装备
		if useWeapon.IsLock || useWeapon.EquipCharacterID != 0 {
			zaplogger.GetSugar().Debugf("%s WeaponStrengthenToS error: mod %d is equipped or locked", this.user.GetUserLogName(), id)
			this.errcode = cs_message.ErrCode_Request_Argument_Err
			return
		}

		def2 := Weapon.GetID(useWeapon.ConfigID)
		if def2 == nil {
			zaplogger.GetSugar().Debugf("%s WeaponStrengthenToS error: cost mod %d def is nil", this.user.GetUserLogName(), id)
			this.errcode = cs_message.ErrCode_Weapon_NotExist
			return
		}

		levelDef2 := WeaponLevelMaxExp.GetID(def2.LevelExpConfig)
		if levelDef2 == nil {
			zaplogger.GetSugar().Debugf("%s WeaponStrengthenToS error: cost %d levelDef is nil", this.user.GetUserLogName(), def2.LevelExpConfig)
			this.errcode = cs_message.ErrCode_Config_NotExist
			return
		}
		totalExp += def2.SupplyExp + levelDef2.GetRangeTotalExp(1, useWeapon.Level) + useWeapon.Exp
		useRes = append(useRes, inoutput.ResDesc{Type: enumType.IOType_Weapon, ID: int32(id), Count: 1})
	}

	for _, v := range msg.GetCostItems() {
		if v.GetCount() <= 0 || v.GetCount() > userBackpack.GetItemCountByTID(v.GetItemID()) {
			zaplogger.GetSugar().Debugf("%s WeaponStrengthenToS error: costItem %d count %d is not enough nil", this.user.GetUserLogName(), v.GetItemID(), v.GetCount())
			this.errcode = cs_message.ErrCode_Item_Not_Enough
			return
		}
		itemDef := Item.GetID(v.GetItemID())
		if itemDef == nil {
			zaplogger.GetSugar().Debugf("%s WeaponStrengthenToS error: item %d itemDef is nil", this.user.GetUserLogName(), v.GetItemID())
			this.errcode = cs_message.ErrCode_Weapon_NotExist
			return
		}

		totalExp += constWeapon.GetSupplyExp(v.GetItemID()) * v.GetCount()
		useRes = append(useRes, inoutput.ResDesc{Type: enumType.IOType_Item, ID: v.GetItemID(), Count: v.GetCount()})
	}

	if totalExp <= 0 {
		zaplogger.GetSugar().Debugf("%s WeaponStrengthenToS error: totalExp = 0", this.user.GetUserLogName())
		this.errcode = cs_message.ErrCode_Resource_NotEnough
		return
	}

	useGold := int32(constWeapon.Get().ExpToGoldRate * float64(totalExp))
	useRes = append(useRes, inoutput.ResDesc{Type: enumType.IOType_UsualAttribute, ID: attr.Gold, Count: useGold})
	// 消耗
	if this.errcode = this.user.DoInputOutput(useRes, nil, false); this.errcode != cs_message.ErrCode_OK {
		zaplogger.GetSugar().Debugf("%s WeaponStrengthenToS error: inouput %s  ", this.user.GetUserLogName(), this.errcode.String())
		return
	}

	userWeapon.Strengthen(levelDef, doWeapon, maxLevel, totalExp)

	this.errcode = cs_message.ErrCode_OK
	zaplogger.GetSugar().Infof("%s Call WeaponStrengthenToS ok", this.user.GetUserLogName())
}

func init() {
	user.RegisterTransFunc(cmdEnum.CS_WeaponStrengthen, func(user *user.User, msg *codecs.Message) transaction.Transaction {
		return &transactionWeaponStrengthen{
			user: user,
			req:  msg,
		}
	})
}
