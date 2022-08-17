package BigSecret

import (
	codecs "initialthree/codec/cs"
	"initialthree/node/common/transaction"
	"initialthree/node/node_game/module"
	"initialthree/node/node_game/module/bigSecret"
	"initialthree/node/node_game/user"
	"initialthree/node/table/excel/DataTable/BigSecretBlessing"
	"initialthree/protocol/cmdEnum"
	"initialthree/protocol/cs/message"
	"initialthree/zaplogger"
)

type transactionBigSecretBlessingLvUp struct {
	transaction.TransactionBase
	user    *user.User
	req     *codecs.Message
	errcode message.ErrCode
}

func (this *transactionBigSecretBlessingLvUp) GetModuleName() string {
	return "Level"
}

func (this *transactionBigSecretBlessingLvUp) Begin() {
	defer func() { this.EndTrans(&message.BigSecretBlessingLvUpToC{}, this.errcode) }()
	this.errcode = message.ErrCode_OK
	zaplogger.GetSugar().Infof("%s BigSecretBlessingLvUpToS", this.user.GetUserLogName())

	m := this.user.GetSubModule(module.BigSecret).(*bigSecret.BigSecretDungeon)
	level := m.GetBlessingLevel()
	//count := m.GetBlessingCount()
	def := BigSecretBlessing.GetID(level)
	if def == nil {
		zaplogger.GetSugar().Infof("%s BigSecretWeaknessRefreshToS failed, level %d config is not exist", this.user.GetUserLogName(), level)
		this.errcode = message.ErrCode_Config_NotExist
		return
	}

}

func init() {
	user.RegisterTransFunc(cmdEnum.CS_BigSecretBlessingLvUp, func(user *user.User, msg *codecs.Message) transaction.Transaction {
		return &transactionBigSecretBlessingLvUp{
			user: user,
			req:  msg,
		}
	})
}
