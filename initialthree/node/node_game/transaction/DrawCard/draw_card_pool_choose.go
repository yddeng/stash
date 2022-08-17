package DrawCard

import (
	codecs "initialthree/codec/cs"
	"initialthree/node/common/transaction"
	"initialthree/node/node_game/module"
	"initialthree/node/node_game/module/drawCard"
	"initialthree/node/table/excel/DataTable/DrawCardsLib"
	"initialthree/zaplogger"

	"initialthree/node/node_game/user"
	"initialthree/protocol/cmdEnum"
	cs_message "initialthree/protocol/cs/message"
)

type transactionDrawCardPoolChoose struct {
	transaction.TransactionBase
	user    *user.User
	req     *codecs.Message
	resp    *cs_message.DrawCardPoolChooseToC
	errCode cs_message.ErrCode
}

func (t *transactionDrawCardPoolChoose) GetModuleName() string {
	return "DrawCard"
}

func (t *transactionDrawCardPoolChoose) Begin() {
	defer func() { t.EndTrans(t.resp, t.errCode) }()
	t.resp = &cs_message.DrawCardPoolChooseToC{}
	t.errCode = cs_message.ErrCode_OK

	msg := t.req.GetData().(*cs_message.DrawCardPoolChooseToS)
	zaplogger.GetSugar().Infof("%s DrawCardPoolChooseToS %v", t.user.GetUserLogName(), msg)

	if drawCardLib := DrawCardsLib.GetID(msg.GetLibID()); drawCardLib == nil {
		t.errCode = cs_message.ErrCode_Config_NotExist // 卡池不存在
	} else if !t.getDrawCardsPoolID(drawCardLib, msg.GetPoolIndex()) {
		t.errCode = cs_message.ErrCode_DrawCard_GuaranteeLibID
	}

	dc := t.user.GetSubModule(module.DrawCard).(*drawCard.DrawCard)
	dc.SetPoolIndex(msg.GetLibID(), msg.GetPoolIndex())
	zaplogger.GetSugar().Infof("%s %s", t.user.GetUserID(), "DrawCardPoolChooseToS ok")
}

func (t *transactionDrawCardPoolChoose) getDrawCardsPoolID(drawCardLib *DrawCardsLib.DrawCardsLib, poolIdx int32) bool {
	if int32(len(drawCardLib.DrawCardsPoolArray)) <= poolIdx {
		return false
	}
	v := drawCardLib.DrawCardsPoolArray[poolIdx]
	if v.PoolID == 0 {
		return false
	}
	return true
}

func init() {
	user.RegisterTransFunc(cmdEnum.CS_DrawCardPoolChoose, func(user *user.User, msg *codecs.Message) transaction.Transaction {
		return &transactionDrawCardPoolChoose{
			user: user,
			req:  msg,
		}
	})
}
