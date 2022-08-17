package Chat

import (
	"github.com/golang/protobuf/proto"
	codecs "initialthree/codec/cs"
	"initialthree/node/common/transaction"
	"initialthree/node/node_game/module"
	"initialthree/node/node_game/module/base"
	"initialthree/zaplogger"

	"initialthree/node/node_game/user"
	"initialthree/protocol/cmdEnum"
	cs_message "initialthree/protocol/cs/message"
)

type transactionChatMessageSend struct {
	transaction.TransactionBase
	user    *user.User
	req     *codecs.Message
	resp    *cs_message.ChatMessageSendToC
	errcode cs_message.ErrCode
}

func (this *transactionChatMessageSend) GetModuleName() string {
	return "Attr"
}

func (this *transactionChatMessageSend) Begin() {
	msg := this.req.GetData().(*cs_message.ChatMessageSendToS)
	zaplogger.GetSugar().Infof("%s %d Call ChatMessageSendToS %v", this.user.GetUserID(), this.user.GetID(), msg)

	this.EndTrans(&cs_message.ChatMessageSendToC{}, cs_message.ErrCode_OK)
	baseModule := this.user.GetSubModule(module.Base).(*base.UserBase)
	resp := &cs_message.ChatMessageSyncToC{
		UserID:        proto.String(this.user.GetUserID()),
		GameID:        proto.Uint64(this.user.GetID()),
		Name:          proto.String(baseModule.GetName()),
		Portrait:      proto.Int32(baseModule.GetPortrait()),
		PortraitFrame: proto.Int32(baseModule.GetPortraitFrame()),
		Message:       proto.String(msg.GetMessage()),
	}
	this.user.Reply(0, resp)
}

func init() {
	user.RegisterTransFunc(cmdEnum.CS_ChatMessageSend, func(user *user.User, msg *codecs.Message) transaction.Transaction {
		return &transactionChatMessageSend{
			user: user,
			req:  msg,
		}
	})
}
