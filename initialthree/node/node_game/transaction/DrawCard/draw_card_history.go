package DrawCard

import (
	codecs "initialthree/codec/cs"
	"initialthree/node/common/transaction"
	"initialthree/node/node_game/module/drawCard"
	"initialthree/zaplogger"

	"initialthree/node/node_game/module"
	"initialthree/node/node_game/user"
	"initialthree/protocol/cmdEnum"
	cs_message "initialthree/protocol/cs/message"
)

type transactionDrawCardHistory struct {
	transaction.TransactionBase
	user    *user.User
	req     *codecs.Message
	resp    *cs_message.DrawCardHistoryToC
	errCode cs_message.ErrCode
}

func (t *transactionDrawCardHistory) GetModuleName() string {
	return "DrawCard"
}

func (t *transactionDrawCardHistory) Begin() {
	defer func() { t.EndTrans(t.resp, t.errCode) }()
	msg := t.req.GetData().(*cs_message.DrawCardHistoryToS)
	zaplogger.GetSugar().Infof("%s %s", t.user.GetUserID(), "DrawCardHistory ==>")
	t.resp = &cs_message.DrawCardHistoryToC{}

	history := t.user.GetSubModule(module.DrawCard).(*drawCard.DrawCard)
	if items := history.GetHistory(msg.GetLibID()); items != nil {
		t.resp.History = items
	} else {
		t.resp.History = make([]*cs_message.DrawCardAward, 0)
	}
	t.errCode = cs_message.ErrCode_OK

}

func init() {
	user.RegisterTransFunc(cmdEnum.CS_DrawCardHistory, func(user *user.User, msg *codecs.Message) transaction.Transaction {
		return &transactionDrawCardHistory{
			user: user,
			req:  msg,
		}
	})
}
