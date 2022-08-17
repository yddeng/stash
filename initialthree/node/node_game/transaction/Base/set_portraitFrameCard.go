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

type transactionBaseSetPortraitFrameCard struct {
	transaction.TransactionBase
	user    *user.User
	req     *codecs.Message
	resp    *cs_message.BaseSetPortraitFrameCardToC
	errcode cs_message.ErrCode
}

func (this *transactionBaseSetPortraitFrameCard) GetModuleName() string {
	return "Base"
}

func (this *transactionBaseSetPortraitFrameCard) Begin() {
	defer func() { this.EndTrans(this.resp, this.errcode) }()
	this.errcode = cs_message.ErrCode_OK
	this.resp = &cs_message.BaseSetPortraitFrameCardToC{}
	msg := this.req.GetData().(*cs_message.BaseSetPortraitFrameCardToS)

	zaplogger.GetSugar().Infof("%s Call BaseSetPortraitFrameCardToS %v", this.user.GetUserLogName(), msg)

	baseModule := this.user.GetSubModule(module.Base).(*base.UserBase)
	if msg.GetPortrait() != -1 {
		baseModule.SetPortrait(msg.GetPortrait())
	}
	if msg.GetPortraitFrame() != -1 {
		baseModule.SetPortraitFrame(msg.GetPortraitFrame())
	}
	if msg.GetCard() != -1 {
		baseModule.SetCard(msg.GetCard())
	}

	zaplogger.GetSugar().Infof("%s Call BaseSetPortraitFrameCardToS ok", this.user.GetUserLogName())
}

func init() {
	user.RegisterTransFunc(cmdEnum.CS_BaseSetPortraitFrameCard, func(user *user.User, msg *codecs.Message) transaction.Transaction {
		return &transactionBaseSetPortraitFrameCard{
			user: user,
			req:  msg,
		}
	})
}
