package BigSecret

import (
	codecs "initialthree/codec/cs"
	"initialthree/node/common/transaction"
	"initialthree/node/node_game/module"
	"initialthree/node/node_game/module/bigSecret"
	"initialthree/node/node_game/user"
	"initialthree/protocol/cmdEnum"
	"initialthree/protocol/cs/message"
	"initialthree/zaplogger"
)

type transactionBigSecretWeakness struct {
	transaction.TransactionBase
	user    *user.User
	req     *codecs.Message
	errcode message.ErrCode
	resp    *message.BigSecretWeaknessToC
}

func (this *transactionBigSecretWeakness) GetModuleName() string {
	return "Level"
}

func (this *transactionBigSecretWeakness) Begin() {
	defer func() { this.EndTrans(this.resp, this.errcode) }()
	this.errcode = message.ErrCode_OK
	this.resp = &message.BigSecretWeaknessToC{}
	msg := this.req.GetData().(*message.BigSecretWeaknessToS)
	zaplogger.GetSugar().Infof("%s BigSecretWeaknessToS %v ", this.user.GetUserLogName(), msg)

	m := this.user.GetSubModule(module.BigSecret).(*bigSecret.BigSecretDungeon)
	unlock := m.Unlocked(msg.GetLevel())
	if !unlock {
		zaplogger.GetSugar().Infof("%s BigSecretWeaknessToS failed, level %d is locked", this.user.GetUserLogName(), msg.GetLevel())
		this.errcode = message.ErrCode_BigSecret_Locked
		return
	}

	wk := m.GetWeakness(msg.GetLevel())
	this.resp.Weakness = wk
}

func init() {
	user.RegisterTransFunc(cmdEnum.CS_BigSecretWeakness, func(user *user.User, msg *codecs.Message) transaction.Transaction {
		return &transactionBigSecretWeakness{
			user: user,
			req:  msg,
		}
	})
}
