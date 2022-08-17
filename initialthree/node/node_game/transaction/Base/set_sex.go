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

type transactionBaseSetSex struct {
	transaction.TransactionBase
	user    *user.User
	req     *codecs.Message
	resp    *cs_message.BaseSetSexToC
	errcode cs_message.ErrCode
}

func (this *transactionBaseSetSex) GetModuleName() string {
	return "Base"
}

func (this *transactionBaseSetSex) Begin() {
	defer func() { this.EndTrans(this.resp, this.errcode) }()
	this.errcode = cs_message.ErrCode_OK
	this.resp = &cs_message.BaseSetSexToC{}
	msg := this.req.GetData().(*cs_message.BaseSetSexToS)

	zaplogger.GetSugar().Infof("%s Call BaseSetSexToS %v", this.user.GetUserLogName(), msg)

	if !(msg.GetSex() == 0 || msg.GetSex() == 1) {
		zaplogger.GetSugar().Infof("%s Call BaseSetSexToS failed, sex %d", this.user.GetUserLogName(), msg.GetSex())
		this.errcode = cs_message.ErrCode_Request_Argument_Err
		return
	}

	baseModule := this.user.GetSubModule(module.Base).(*base.UserBase)
	baseModule.SetSex(msg.GetSex())

	zaplogger.GetSugar().Infof("%s Call BaseSetSexToS ok", this.user.GetUserLogName())
}

func init() {
	user.RegisterTransFunc(cmdEnum.CS_BaseSetSex, func(user *user.User, msg *codecs.Message) transaction.Transaction {
		return &transactionBaseSetSex{
			user: user,
			req:  msg,
		}
	})
}
