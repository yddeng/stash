package BigSecret

import (
	codecs "initialthree/codec/cs"
	"initialthree/node/common/attr"
	"initialthree/node/common/enumType"
	"initialthree/node/common/inoutput"
	"initialthree/node/common/transaction"
	"initialthree/node/node_game/module"
	"initialthree/node/node_game/module/bigSecret"
	"initialthree/node/node_game/user"
	"initialthree/node/table/excel/ConstTable/BigSecret"
	"initialthree/protocol/cmdEnum"
	"initialthree/protocol/cs/message"
	"initialthree/zaplogger"
)

type transactionBigSecretWeaknessRefresh struct {
	transaction.TransactionBase
	user    *user.User
	req     *codecs.Message
	errcode message.ErrCode
	resp    *message.BigSecretWeaknessRefreshToC
}

func (this *transactionBigSecretWeaknessRefresh) GetModuleName() string {
	return "Level"
}

func (this *transactionBigSecretWeaknessRefresh) Begin() {
	defer func() { this.EndTrans(this.resp, this.errcode) }()
	this.errcode = message.ErrCode_OK
	this.resp = &message.BigSecretWeaknessRefreshToC{}
	msg := this.req.GetData().(*message.BigSecretWeaknessRefreshToS)
	zaplogger.GetSugar().Infof("%s BigSecretWeaknessRefreshToS %v ", this.user.GetUserLogName(), msg)

	m := this.user.GetSubModule(module.BigSecret).(*bigSecret.BigSecretDungeon)
	unlock := m.Unlocked(msg.GetLevel())
	if !unlock {
		zaplogger.GetSugar().Infof("%s BigSecretWeaknessRefreshToS failed, level %d is locked", this.user.GetUserLogName(), msg.GetLevel())
		this.errcode = message.ErrCode_BigSecret_Locked
		return
	}

	def := BigSecret.GetID(1)
	times := m.GetWeaknessRefreshTimes()
	if times >= int32(len(def.WeaknessRefreshCostArray)) {
		zaplogger.GetSugar().Infof("%s BigSecretWeaknessRefreshToS failed, refresh time is not enough ", this.user.GetUserLogName())
		this.errcode = message.ErrCode_BigSecret_RefreshTimesNotEnough
		return
	}

	cost := def.WeaknessRefreshCostArray[times].Cost
	if cost > 0 {
		in := []inoutput.ResDesc{{Type: enumType.IOType_UsualAttribute, ID: attr.Diamond, Count: cost}}
		if this.errcode = this.user.DoInputOutput(in, nil, false); this.errcode != message.ErrCode_OK {
			return
		}
	}

	wk := m.WeaknessRefresh(msg.GetLevel())
	this.resp.Weakness = wk
}

func init() {
	user.RegisterTransFunc(cmdEnum.CS_BigSecretWeaknessRefresh, func(user *user.User, msg *codecs.Message) transaction.Transaction {
		return &transactionBigSecretWeaknessRefresh{
			user: user,
			req:  msg,
		}
	})
}
