package Weapon

import (
	codecs "initialthree/codec/cs"
	"initialthree/node/common/transaction"
	"initialthree/node/node_game/module/weapon"
	"initialthree/zaplogger"

	"initialthree/node/node_game/module"
	"initialthree/node/node_game/user"
	"initialthree/protocol/cmdEnum"
	cs_message "initialthree/protocol/cs/message"
)

type transactionWeaponLock struct {
	transaction.TransactionBase
	user    *user.User
	req     *codecs.Message
	resp    *cs_message.WeaponLockToC
	errcode cs_message.ErrCode
}

func (this *transactionWeaponLock) GetModuleName() string {
	return "Weapon"
}

func (this *transactionWeaponLock) Begin() {
	defer func() { this.EndTrans(this.resp, this.errcode) }()
	this.resp = &cs_message.WeaponLockToC{}
	msg := this.req.GetData().(*cs_message.WeaponLockToS)

	zaplogger.GetSugar().Infof("%s Call WeaponLockToS %v", this.user.GetUserLogName(), msg)

	userWeapon := this.user.GetSubModule(module.Weapon).(*weapon.UserWeapon)
	doWeapon := userWeapon.GetWeapon(msg.GetWeaponID())
	if doWeapon == nil {
		zaplogger.GetSugar().Debugf("%s WeaponLockToS error:  equip %d is nil", this.user.GetUserLogName(), msg.GetWeaponID())
		this.errcode = cs_message.ErrCode_Weapon_NotExist
		return
	}

	userWeapon.Lock(doWeapon, msg.GetIsLock())

	this.errcode = cs_message.ErrCode_OK
	zaplogger.GetSugar().Infof("%s Call WeaponLockToS ok", this.user.GetUserLogName())

}

func init() {
	user.RegisterTransFunc(cmdEnum.CS_WeaponLock, func(user *user.User, msg *codecs.Message) transaction.Transaction {
		return &transactionWeaponLock{
			user: user,
			req:  msg,
		}
	})
}
