package Base

import (
	codecs "initialthree/codec/cs"
	"initialthree/node/common/transaction"
	"initialthree/node/node_game/module"
	"initialthree/node/node_game/module/base"
	"initialthree/zaplogger"

	"initialthree/node/node_game/user"
	"initialthree/protocol/cmdEnum"
	cs_message "initialthree/protocol/cs/message"
)

type transactionBaseSetSignature struct {
	transaction.TransactionBase
	user    *user.User
	req     *codecs.Message
	resp    *cs_message.BaseSetSignatureToC
	errcode cs_message.ErrCode
}

func (this *transactionBaseSetSignature) GetModuleName() string {
	return "Base"
}

func (this *transactionBaseSetSignature) Begin() {
	defer func() { this.EndTrans(this.resp, this.errcode) }()
	this.errcode = cs_message.ErrCode_OK
	this.resp = &cs_message.BaseSetSignatureToC{}
	msg := this.req.GetData().(*cs_message.BaseSetSignatureToS)

	zaplogger.GetSugar().Infof("%s Call BaseSetSignatureToS %v", this.user.GetUserLogName(), msg)

	baseModule := this.user.GetSubModule(module.Base).(*base.UserBase)
	baseModule.SetSignature(msg.GetSignature())

	zaplogger.GetSugar().Infof("%s Call BaseSetSignatureToS ok", this.user.GetUserLogName())
}

func init() {
	user.RegisterTransFunc(cmdEnum.CS_BaseSetSignature, func(user *user.User, msg *codecs.Message) transaction.Transaction {
		return &transactionBaseSetSignature{
			user: user,
			req:  msg,
		}
	})
}
