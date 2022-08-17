package Mail

import (
	codecs "initialthree/codec/cs"
	"initialthree/node/common/droppool"
	"initialthree/node/common/inoutput"
	"initialthree/node/common/transaction"
	"initialthree/node/node_game/module"
	"initialthree/node/node_game/module/mail"
	"initialthree/node/node_game/user"
	"initialthree/protocol/cmdEnum"
	cs_message "initialthree/protocol/cs/message"
	"initialthree/zaplogger"
)

type transactionMailRead struct {
	transaction.TransactionBase
	user    *user.User
	req     *codecs.Message
	errcode cs_message.ErrCode
}

func (this *transactionMailRead) GetModuleName() string {
	return "Mail"
}

func (this *transactionMailRead) Begin() {
	defer func() { this.EndTrans(&cs_message.MailReadToC{}, this.errcode) }()

	this.errcode = cs_message.ErrCode_OK
	msg := this.req.GetData().(*cs_message.MailReadToS)

	zaplogger.GetSugar().Infof("%s Call MailRead %v", this.user.GetUserLogName(), msg)

	mailModule := this.user.GetSubModule(module.Mail).(*mail.Mail)
	ids := make([]uint32, 0, len(msg.GetMailIDs()))
	for _, id := range msg.GetMailIDs() {
		m := mailModule.GetMail(id)
		if m == nil || m.GetRead() {
			zaplogger.GetSugar().Debug(m)
			continue
		}
		if !mailModule.IsExpire(m) && len(m.GetAwards().GetAwardInfos()) > 0 {
			award := m.GetAwards()
			out := make([]inoutput.ResDesc, 0, len(award.GetAwardInfos()))
			for _, info := range award.GetAwardInfos() {
				out = append(out, inoutput.ResDesc{
					Type:  droppool.DropType2IOType(int(info.GetType())),
					ID:    info.GetID(),
					Count: info.GetCount(),
				})
			}

			if this.errcode = this.user.DoInputOutput(nil, out, false); this.errcode != cs_message.ErrCode_OK {
				zaplogger.GetSugar().Debugf("%s MailRead error: inouput %s  ", this.user.GetUserLogName(), this.errcode.String())
				break
			}
		}
		ids = append(ids, id)
	}

	mailModule.MailRead(ids)

}

func init() {
	user.RegisterTransFunc(cmdEnum.CS_MailRead, func(user *user.User, msg *codecs.Message) transaction.Transaction {
		return &transactionMailRead{
			user: user,
			req:  msg,
		}
	})
}
