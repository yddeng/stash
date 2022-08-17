package Inoutput

import (
	codecs "initialthree/codec/cs"
	"initialthree/node/common/transaction"
	"initialthree/node/node_game/user"
	IOTable "initialthree/node/table/excel/DataTable/InputOutput"
	"initialthree/protocol/cmdEnum"
	cs_message "initialthree/protocol/cs/message"
	"initialthree/zaplogger"
)

type transactionInoutput struct {
	transaction.TransactionBase
	user    *user.User
	req     *codecs.Message
	errcode cs_message.ErrCode
}

func (t *transactionInoutput) Begin() {
	defer func() { t.EndTrans(&cs_message.InoutputToC{}, t.errcode) }()

	t.errcode = cs_message.ErrCode_OK
	reqMsg := t.req.GetData().(*cs_message.InoutputToS)
	zaplogger.GetSugar().Infof("%s InoutputToS %v", t.user.GetUserLogName(), reqMsg)

	if reqMsg.GetCount() <= 0 || reqMsg.GetId() <= 0 {
		t.errcode = cs_message.ErrCode_Request_Argument_Err
		return
	}

	cfg := IOTable.GetID(reqMsg.GetId())
	if nil == cfg {
		t.errcode = cs_message.ErrCode_Config_NotExist
		return
	}

	in, out := cfg.InOut(reqMsg.GetCount())
	if t.errcode = t.user.DoInputOutput(in, out, false); t.errcode != cs_message.ErrCode_OK {
		zaplogger.GetSugar().Debugf("%s InoutputToS error: inouput %s  ", t.user.GetUserLogName(), t.errcode.String())
		return
	}

}

func (t *transactionInoutput) GetModuleName() string {
	return "Backpack"
}

func init() {
	user.RegisterTransFunc(cmdEnum.CS_Inoutput, func(user *user.User, msg *codecs.Message) transaction.Transaction {
		return &transactionInoutput{
			user: user,
			req:  msg,
		}
	})
}
