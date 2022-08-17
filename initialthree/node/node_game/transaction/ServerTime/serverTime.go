package ServerTime

import (
	"github.com/golang/protobuf/proto"
	codecs "initialthree/codec/cs"
	"initialthree/node/common/transaction"
	"initialthree/node/node_game/user"
	"initialthree/protocol/cmdEnum"
	cs_message "initialthree/protocol/cs/message"
	"time"
)

type transactionServerTime struct {
	transaction.TransactionBase
	user *user.User
	req  *codecs.Message
}

func (t *transactionServerTime) Begin() {
	resp := &cs_message.ServerTimeToC{
		PhysTime: proto.Int64(time.Now().Unix()),
	}
	t.EndTrans(resp, cs_message.ErrCode_OK)
}

func (t *transactionServerTime) GetModuleName() string {
	return "ServerTime"
}

func init() {
	//  ServerTime 特殊的 trans，没有前置检查
	user.RegisterTransStep(cmdEnum.CS_ServerTime, func(user *user.User, msg *codecs.Message) transaction.Transaction {
		return &transactionServerTime{
			user: user,
			req:  msg,
		}
	})
}
