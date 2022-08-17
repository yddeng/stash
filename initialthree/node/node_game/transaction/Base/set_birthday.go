package Base

import (
	codecs "initialthree/codec/cs"
	"initialthree/node/common/transaction"
	"initialthree/node/node_game/module"
	"initialthree/node/node_game/module/base"
	"initialthree/node/node_game/user"
	"initialthree/protocol/cmdEnum"
	cs_message "initialthree/protocol/cs/message"
	"initialthree/zaplogger"
)

type transactionBaseSetBirthday struct {
	transaction.TransactionBase
	user    *user.User
	req     *codecs.Message
	resp    *cs_message.BaseSetBirthdayToC
	errcode cs_message.ErrCode
}

func (this *transactionBaseSetBirthday) GetModuleName() string {
	return "Base"
}

func (this *transactionBaseSetBirthday) Begin() {
	defer func() { this.EndTrans(this.resp, this.errcode) }()
	this.errcode = cs_message.ErrCode_OK
	this.resp = &cs_message.BaseSetBirthdayToC{}
	msg := this.req.GetData().(*cs_message.BaseSetBirthdayToS)

	zaplogger.GetSugar().Infof("%s Call BaseSetBirthdayToS %v", this.user.GetUserLogName(), msg)

	baseModule := this.user.GetSubModule(module.Base).(*base.UserBase)
	if baseModule.GetBirthday() != "" {
		zaplogger.GetSugar().Infof("%s Call BaseSetSignatureToS failed, already set", this.user.GetUserLogName())
		this.errcode = cs_message.ErrCode_ERROR
		return
	}

	baseModule.SetBirthday(msg.GetBirthday())
	zaplogger.GetSugar().Infof("%s Call BaseSetSignatureToS ok", this.user.GetUserLogName())
}

func init() {
	user.RegisterTransFunc(cmdEnum.CS_BaseSetBirthday, func(user *user.User, msg *codecs.Message) transaction.Transaction {
		return &transactionBaseSetBirthday{
			user: user,
			req:  msg,
		}
	})
}
