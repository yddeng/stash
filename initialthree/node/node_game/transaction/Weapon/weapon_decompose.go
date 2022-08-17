package Weapon

import (
	codecs "initialthree/codec/cs"
	"initialthree/node/common/attr"
	"initialthree/node/common/enumType"
	"initialthree/node/common/inoutput"
	"initialthree/node/common/transaction"
	"initialthree/node/node_game/module/weapon"
	"initialthree/node/table/excel/DataTable/WeaponLevelMaxExp"
	"initialthree/zaplogger"

	"initialthree/node/node_game/module"
	"initialthree/node/node_game/user"
	constWeapon "initialthree/node/table/excel/ConstTable/Weapon"
	"initialthree/node/table/excel/DataTable/Weapon"
	"initialthree/node/table/excel/DataTable/WeaponDecompose"
	"initialthree/protocol/cmdEnum"
	cs_message "initialthree/protocol/cs/message"
)

type transactionWeaponDecompose struct {
	transaction.TransactionBase
	user    *user.User
	req     *codecs.Message
	resp    *cs_message.WeaponDecomposeToC
	errcode cs_message.ErrCode
}

func (this *transactionWeaponDecompose) GetModuleName() string {
	return "Weapon"
}

func (this *transactionWeaponDecompose) Begin() {
	defer func() { this.EndTrans(this.resp, this.errcode) }()
	this.resp = &cs_message.WeaponDecomposeToC{}
	msg := this.req.GetData().(*cs_message.WeaponDecomposeToS)

	zaplogger.GetSugar().Infof("%s Call WeaponDecomposeToS %v", this.user.GetUserLogName(), msg)

	userWeapon := this.user.GetSubModule(module.Weapon).(*weapon.UserWeapon)

	useRes := make([]inoutput.ResDesc, 0, len(msg.GetWeaponIDs()))
	addGold := float64(0)
	toItems := map[int32]float64{}
	for _, id := range msg.GetWeaponIDs() {
		w := userWeapon.GetWeapon(id)
		if w == nil {
			zaplogger.GetSugar().Debugf("%s WeaponDecomposeToS error: weapon %d is nil", this.user.GetUserLogName(), id)
			this.errcode = cs_message.ErrCode_Weapon_NotExist
			return
		}
		// 装备已被装备
		if w.IsLock || w.EquipCharacterID != 0 {
			zaplogger.GetSugar().Debugf("%s WeaponDecomposeToS error: weapon %d is equipped or locked", this.user.GetUserLogName(), id)
			this.errcode = cs_message.ErrCode_Equip_Equipped
			return
		}

		def := Weapon.GetID(w.ConfigID)
		if def == nil {
			zaplogger.GetSugar().Debugf("%s WeaponDecomposeToS error: weapon %d def is nil", this.user.GetUserLogName(), id)
			this.errcode = cs_message.ErrCode_Weapon_NotExist
			return
		}

		decDef := WeaponDecompose.GetID(def.DecomposeConfig)
		if decDef == nil {
			zaplogger.GetSugar().Debugf("%s WeaponDecomposeToS error:  WeaponDecompose %d def is nil", this.user.GetUserLogName(), id)
			this.errcode = cs_message.ErrCode_Weapon_NotExist
			return
		}

		levelDef := WeaponLevelMaxExp.GetID(def.LevelExpConfig)
		if levelDef == nil {
			zaplogger.GetSugar().Debugf("%s WeaponDecomposeToS error:  %d levelDef is nil", this.user.GetUserLogName(), def.LevelExpConfig)
			this.errcode = cs_message.ErrCode_Config_NotExist
			return
		}
		toItemId := decDef.ExpToItemID
		supplyExp := constWeapon.GetSupplyExp(toItemId)
		if supplyExp == 0 {
			zaplogger.GetSugar().Debugf("%s WeaponDecomposeToS error:  %d supplyExp is nil", this.user.GetUserLogName(), toItemId)
			this.errcode = cs_message.ErrCode_Config_NotExist
			return
		}

		// 总经验
		exp := levelDef.GetRangeTotalExp(1, w.Level) + w.Exp
		// 转金币
		addGold += float64(exp)/decDef.ExpToOneGoldRate + float64(decDef.ReturnGoldCount)
		// 转道具数量
		val := float64(exp) / float64(supplyExp)
		toItems[toItemId] += val
		for _, v := range decDef.Items {
			toItems[v.Id] += float64(v.Count)
		}

		useRes = append(useRes, inoutput.ResDesc{Type: enumType.IOType_Weapon, ID: int32(id), Count: 1})
	}

	// 消耗添加
	addRes := make([]inoutput.ResDesc, 0, len(toItems)+1)
	addRes = append(addRes, inoutput.ResDesc{Type: enumType.IOType_UsualAttribute, ID: attr.Gold, Count: int32(addGold)})
	for itemID, itemCount := range toItems {
		addRes = append(addRes, inoutput.ResDesc{Type: enumType.IOType_Item, ID: itemID, Count: int32(itemCount)})
	}

	if this.errcode = this.user.DoInputOutput(useRes, addRes, true); this.errcode != cs_message.ErrCode_OK {
		zaplogger.GetSugar().Debugf("%s WeaponDecomposeToS error: inouput %s  ", this.user.GetUserLogName(), this.errcode.String())
		return
	}

	zaplogger.GetSugar().Infof("%s Call WeaponDecomposeToS ok", this.user.GetUserLogName())
}

func init() {
	user.RegisterTransFunc(cmdEnum.CS_WeaponDecompose, func(user *user.User, msg *codecs.Message) transaction.Transaction {
		return &transactionWeaponDecompose{
			user: user,
			req:  msg,
		}
	})
}
