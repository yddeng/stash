package Equip

import (
	codecs "initialthree/codec/cs"
	"initialthree/node/common/transaction"
	"initialthree/zaplogger"

	"initialthree/node/node_game/module"
	"initialthree/node/node_game/module/equip"
	"initialthree/node/node_game/user"
	"initialthree/protocol/cmdEnum"
	cs_message "initialthree/protocol/cs/message"
)

type transactionEquipLock struct {
	transaction.TransactionBase
	user    *user.User
	req     *codecs.Message
	resp    *cs_message.EquipLockToC
	errcode cs_message.ErrCode
}

func (this *transactionEquipLock) GetModuleName() string {
	return "Equip"
}

func (this *transactionEquipLock) Begin() {
	defer func() { this.EndTrans(this.resp, this.errcode) }()
	this.resp = &cs_message.EquipLockToC{}
	msg := this.req.GetData().(*cs_message.EquipLockToS)

	zaplogger.GetSugar().Infof("%s Call EquipLockToS %v", this.user.GetUserLogName(), msg)

	userEquip := this.user.GetSubModule(module.Equip).(*equip.UserEquip)
	doEquip := userEquip.GetEquip(msg.GetEquipID())
	if doEquip == nil {
		zaplogger.GetSugar().Debugf("%s EquipLockToS error:  equip %d is nil", this.user.GetUserLogName(), msg.GetEquipID())
		this.errcode = cs_message.ErrCode_Equip_NotExist
		return
	}

	userEquip.Lock(doEquip, msg.GetIsLock())

	this.errcode = cs_message.ErrCode_OK
	zaplogger.GetSugar().Infof("%s Call EquipLockToS ok", this.user.GetUserLogName())

}

func init() {
	user.RegisterTransFunc(cmdEnum.CS_EquipLock, func(user *user.User, msg *codecs.Message) transaction.Transaction {
		return &transactionEquipLock{
			user: user,
			req:  msg,
		}
	})
}
