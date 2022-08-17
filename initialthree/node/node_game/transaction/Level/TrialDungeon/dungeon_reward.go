package SecertDungeon

import (
	codecs "initialthree/codec/cs"
	"initialthree/node/common/droppool"
	"initialthree/node/common/transaction"
	"initialthree/node/node_game/module"
	"initialthree/node/node_game/module/trialDungeon"
	"initialthree/node/node_game/user"
	"initialthree/node/table/excel/DataTable/TrialDungeon"
	"initialthree/protocol/cmdEnum"
	cs_message "initialthree/protocol/cs/message"
	"initialthree/zaplogger"
)

type transactionTrialDungeonReward struct {
	transaction.TransactionBase
	user    *user.User
	req     *codecs.Message
	errcode cs_message.ErrCode
}

func (this *transactionTrialDungeonReward) GetModuleName() string {
	return "Level"
}

func (this *transactionTrialDungeonReward) Begin() {
	defer func() { this.EndTrans(&cs_message.TrialDungeonRewardToC{}, this.errcode) }()
	this.errcode = cs_message.ErrCode_OK
	msg := this.req.GetData().(*cs_message.TrialDungeonRewardToS)
	zaplogger.GetSugar().Infof("%s TrialDungeonRewardToS %v ", this.user.GetUserLogName(), msg)

	userTrial := this.user.GetSubModule(module.TrialDungeon).(*trialDungeon.TrialDungeon)
	for _, id := range msg.GetDungeonIDs() {
		def := TrialDungeon.GetID(id)
		if def == nil {
			this.errcode = cs_message.ErrCode_Config_NotExist
			return
		}

		trialData := userTrial.GetTrialDungeon(id)
		if trialData == nil || trialData.GetGetReward() {
			this.errcode = cs_message.ErrCode_Request_Argument_Err
			return
		}

		userTrial.DungeonReward(id)
		this.user.ApplyDropAward(droppool.DropAward(def.DropPoolID))

	}

}

func init() {
	user.RegisterTransFunc(cmdEnum.CS_TrialDungeonReward, func(user *user.User, msg *codecs.Message) transaction.Transaction {
		return &transactionTrialDungeonReward{
			user: user,
			req:  msg,
		}
	})
}
