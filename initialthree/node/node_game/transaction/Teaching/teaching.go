package Teaching

import (
	codecs "initialthree/codec/cs"
	"initialthree/node/common/transaction"
	"initialthree/node/node_game/user"
	"initialthree/protocol/cmdEnum"
	cs_message "initialthree/protocol/cs/message"
)

type transactionTeaching struct {
	transaction.TransactionBase
	user *user.User
	req  *codecs.Message
}

func (t *transactionTeaching) Begin() {
	t.EndTrans(&cs_message.TeachNoneToC{}, cs_message.ErrCode_OK)
}

func (t *transactionTeaching) GetModuleName() string {
	return "Teaching"
}

func init() {
	//  ServerTime 特殊的 trans，没有前置检查
	user.RegisterTransStep(cmdEnum.CS_TeachNone, func(user *user.User, msg *codecs.Message) transaction.Transaction {
		return &transactionTeaching{
			user: user,
			req:  msg,
		}
	})
}
