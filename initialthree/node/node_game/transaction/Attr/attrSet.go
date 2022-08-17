package Attr

import (
	codecs "initialthree/codec/cs"
	attr2 "initialthree/node/common/attr"
	"initialthree/node/common/transaction"
	"initialthree/zaplogger"

	"initialthree/node/node_game/module"
	"initialthree/node/node_game/module/attr"
	"initialthree/node/node_game/user"
	"initialthree/protocol/cmdEnum"
	cs_message "initialthree/protocol/cs/message"
)

type transactionAttrSet struct {
	transaction.TransactionBase
	user    *user.User
	req     *codecs.Message
	resp    *cs_message.AttrSetToC
	errcode cs_message.ErrCode
}

func (this *transactionAttrSet) GetModuleName() string {
	return "Attr"
}

func (this *transactionAttrSet) Begin() {

	defer func() { this.EndTrans(this.resp, this.errcode) }()
	this.errcode = cs_message.ErrCode_OK

	this.resp = &cs_message.AttrSetToC{}
	msg := this.req.GetData().(*cs_message.AttrSetToS)

	zaplogger.GetSugar().Infof("%s %d Call AttrSetToS %v", this.user.GetUserID(), this.user.GetID(), msg)

	userAttr := this.user.GetSubModule(module.Attr).(*attr.UserAttr)

	id := msg.GetID()
	switch id {
	case attr2.YuruCharacterID, attr2.PassedPrologue, attr2.CurrentTitle:
		if _, err := userAttr.SetAttr(id, msg.GetVal(), false); err != nil {
			zaplogger.GetSugar().Debugf("%s AttrSetToS error: %s ", this.user.GetUserLogName(), err)
			this.errcode = cs_message.ErrCode_ERROR
		}
	default:
		zaplogger.GetSugar().Debugf("%s AttrSetToS error: id %d can't set", this.user.GetUserLogName(), id)
		this.errcode = cs_message.ErrCode_Request_Argument_Err
	}
}

func init() {
	user.RegisterTransFunc(cmdEnum.CS_AttrSet, func(user *user.User, msg *codecs.Message) transaction.Transaction {
		return &transactionAttrSet{
			user: user,
			req:  msg,
		}
	})
}
