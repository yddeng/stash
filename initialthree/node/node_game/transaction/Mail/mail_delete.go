package Mail

import (
	codecs "initialthree/codec/cs"
	"initialthree/node/common/transaction"
	"initialthree/node/node_game/module"
	"initialthree/node/node_game/module/mail"
	"initialthree/node/node_game/user"
	"initialthree/protocol/cmdEnum"
	cs_message "initialthree/protocol/cs/message"
	"initialthree/zaplogger"
)

type transactionMailDelete struct {
	transaction.TransactionBase
	user    *user.User
	req     *codecs.Message
	errcode cs_message.ErrCode
}

func (this *transactionMailDelete) GetModuleName() string {
	return "Mail"
}

func (this *transactionMailDelete) Begin() {
	defer func() { this.EndTrans(&cs_message.MailDeleteToC{}, this.errcode) }()

	this.errcode = cs_message.ErrCode_OK
	msg := this.req.GetData().(*cs_message.MailDeleteToS)

	zaplogger.GetSugar().Infof("%s Call MailDelete %v", this.user.GetUserLogName(), msg)

	mailModule := this.user.GetSubModule(module.Mail).(*mail.Mail)
	ids := make([]uint32, 0, len(msg.GetMailIDs()))
	for _, id := range msg.GetMailIDs() {
		m := mailModule.GetMail(id)
		if m == nil || (!mailModule.IsExpire(m) && !m.GetRead()) {
			continue
		}
		ids = append(ids, id)
	}

	mailModule.MailDelete(ids)

}

func init() {
	user.RegisterTransFunc(cmdEnum.CS_MailDelete, func(user *user.User, msg *codecs.Message) transaction.Transaction {
		return &transactionMailDelete{
			user: user,
			req:  msg,
		}
	})
}
