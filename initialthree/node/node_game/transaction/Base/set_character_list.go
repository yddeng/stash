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

type transactionBaseSetCharacterList struct {
	transaction.TransactionBase
	user    *user.User
	req     *codecs.Message
	resp    *cs_message.BaseSetCharacterDisplayToC
	errcode cs_message.ErrCode
}

func (this *transactionBaseSetCharacterList) GetModuleName() string {
	return "Base"
}

func (this *transactionBaseSetCharacterList) Begin() {
	defer func() { this.EndTrans(this.resp, this.errcode) }()
	this.errcode = cs_message.ErrCode_OK
	this.resp = &cs_message.BaseSetCharacterDisplayToC{}
	msg := this.req.GetData().(*cs_message.BaseSetCharacterDisplayToS)

	zaplogger.GetSugar().Infof("%s Call BaseSetCharacterDisplayToS %v", this.user.GetUserLogName(), msg)

	if len(msg.GetCharacterIDs()) > 8 {
		this.errcode = cs_message.ErrCode_Request_Argument_Err
		return
	}

	baseModule := this.user.GetSubModule(module.Base).(*base.UserBase)
	baseModule.SetCharacterList(msg.GetCharacterIDs())
	zaplogger.GetSugar().Infof("%s Call BaseSetCharacterDisplayToS ok", this.user.GetUserLogName())
}

func init() {
	user.RegisterTransFunc(cmdEnum.CS_BaseSetCharacterDisplay, func(user *user.User, msg *codecs.Message) transaction.Transaction {
		return &transactionBaseSetCharacterList{
			user: user,
			req:  msg,
		}
	})
}
