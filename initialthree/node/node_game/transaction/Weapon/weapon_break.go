package Weapon

import (
	codecs "initialthree/codec/cs"
	attr2 "initialthree/node/common/attr"
	"initialthree/node/common/enumType"
	"initialthree/node/common/inoutput"
	"initialthree/node/common/transaction"
	"initialthree/node/node_game/module/weapon"
	"initialthree/node/table/excel/DataTable/WeaponBreakThough"
	"initialthree/node/table/excel/DataTable/WeaponRarity"
	"initialthree/zaplogger"

	"initialthree/node/node_game/module"
	"initialthree/node/node_game/user"
	"initialthree/node/table/excel/DataTable/Weapon"
	"initialthree/protocol/cmdEnum"
	cs_message "initialthree/protocol/cs/message"
)

type transactionWeaponBreak struct {
	transaction.TransactionBase
	user    *user.User
	req     *codecs.Message
	resp    *cs_message.WeaponBreakToC
	errcode cs_message.ErrCode
}

func (this *transactionWeaponBreak) GetModuleName() string {
	return "Weapon"
}

func (this *transactionWeaponBreak) Begin() {
	defer func() { this.EndTrans(this.resp, this.errcode) }()

	this.resp = &cs_message.WeaponBreakToC{}
	msg := this.req.GetData().(*cs_message.WeaponBreakToS)

	zaplogger.GetSugar().Infof("%s Call WeaponBreakToS %v", this.user.GetUserLogName(), msg)

	userWeapon := this.user.GetSubModule(module.Weapon).(*weapon.UserWeapon)
	doWeapon := userWeapon.GetWeapon(msg.GetWeaponID())
	if doWeapon == nil {
		zaplogger.GetSugar().Debugf("%s WeaponBreakToS error: weapon %d is nil", this.user.GetUserLogName(), msg.GetWeaponID())
		this.errcode = cs_message.ErrCode_Weapon_NotExist
		return
	}

	def := Weapon.GetID(doWeapon.ConfigID)
	if def == nil {
		zaplogger.GetSugar().Debugf("%s WeaponBreakToS error: %d def is nil", this.user.GetUserLogName(), msg.GetWeaponID())
		this.errcode = cs_message.ErrCode_Config_NotExist
		return
	}

	rarityDef := WeaponRarity.GetID(def.RarityTypeEnum)
	breakDef := WeaponBreakThough.GetID(def.BreakThroughConfig)
	if breakDef == nil || rarityDef == nil {
		zaplogger.GetSugar().Debugf("%s WeaponBreakToS error: %d breakDef or rarityDef is nil", this.user.GetUserLogName(), msg.GetWeaponID())
		this.errcode = cs_message.ErrCode_Config_NotExist
		return
	}

	if doWeapon.BreakTimes >= rarityDef.MaxBreakThrough || doWeapon.BreakTimes >= int32(len(breakDef.BreakLevel)) {
		zaplogger.GetSugar().Debugf("%s WeaponBreakToS error: %d break max", this.user.GetUserLogName(), msg.GetWeaponID())
		this.errcode = cs_message.ErrCode_Weapon_BreakMax
		return
	}

	nextBreak := breakDef.GetBreakTimesLevel(doWeapon.BreakTimes + 1)
	// 等级限制
	if doWeapon.Level < nextBreak.Level {
		zaplogger.GetSugar().Debugf("%s WeaponBreakToS error: %d user level low", this.user.GetUserLogName(), msg.GetWeaponID())
		this.errcode = cs_message.ErrCode_Level_Low
		return
	}

	costItems := nextBreak.CostItems()
	//扣除消耗
	useRes := make([]inoutput.ResDesc, 0, len(costItems)+1)
	useRes = append(useRes, inoutput.ResDesc{Type: enumType.IOType_UsualAttribute, ID: attr2.Gold, Count: nextBreak.Gold})
	for _, v := range costItems {
		useRes = append(useRes, inoutput.ResDesc{Type: enumType.IOType_Item, ID: v.ItemID, Count: v.Count})
	}

	if this.errcode = this.user.DoInputOutput(useRes, nil, false); this.errcode != cs_message.ErrCode_OK {
		zaplogger.GetSugar().Debugf("%s WeaponBreakToS error: inouput %s  ", this.user.GetUserLogName(), this.errcode.String())
		return
	}

	userWeapon.AddBreakTimes(doWeapon)

	this.errcode = cs_message.ErrCode_OK
	zaplogger.GetSugar().Infof("%s Call WeaponBreakToS ok", this.user.GetUserLogName())
}

func init() {
	user.RegisterTransFunc(cmdEnum.CS_WeaponBreak, func(user *user.User, msg *codecs.Message) transaction.Transaction {
		return &transactionWeaponBreak{
			user: user,
			req:  msg,
		}
	})
}
