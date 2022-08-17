package SecertDungeon

import (
	codecs "initialthree/codec/cs"
	"initialthree/node/common/droppool"
	"initialthree/node/common/transaction"
	"initialthree/node/node_game/module"
	"initialthree/node/node_game/module/trialDungeon"
	"initialthree/node/node_game/user"
	"initialthree/node/table/excel/DataTable/TrialStageReward"
	"initialthree/protocol/cmdEnum"
	cs_message "initialthree/protocol/cs/message"
	"initialthree/zaplogger"
)

type transactionTrialStageDungeon struct {
	transaction.TransactionBase
	user    *user.User
	req     *codecs.Message
	errcode cs_message.ErrCode
}

func (this *transactionTrialStageDungeon) GetModuleName() string {
	return "Level"
}

func (this *transactionTrialStageDungeon) Begin() {
	defer func() { this.EndTrans(&cs_message.TrialStageRewardToC{}, this.errcode) }()
	this.errcode = cs_message.ErrCode_OK
	msg := this.req.GetData().(*cs_message.TrialStageRewardToS)
	zaplogger.GetSugar().Infof("%s TrialStageRewardToS %v ", this.user.GetUserLogName(), msg)

	def := TrialStageReward.GetID(msg.GetStageID())
	if def == nil {
		zaplogger.GetSugar().Infof("%s TrialStageRewardToS failed, config %d is nil", this.user.GetUserLogName(), msg.GetStageID())
		this.errcode = cs_message.ErrCode_Config_NotExist
		return
	}

	userTrial := this.user.GetSubModule(module.TrialDungeon).(*trialDungeon.TrialDungeon)

	if def.Score > userTrial.GetTrialCount() {
		zaplogger.GetSugar().Infof("%s TrialStageRewardToS failed, trial count not enough ", this.user.GetUserLogName())
		this.errcode = cs_message.ErrCode_Resource_NotEnough
		return
	}

	if userTrial.GetStageReward(msg.GetStageID()) != nil {
		zaplogger.GetSugar().Infof("%s TrialStageRewardToS failed,stage %d already get reward", this.user.GetUserLogName(), msg.GetStageID())
		this.errcode = cs_message.ErrCode_Request_Argument_Err
		return
	}

	this.user.ApplyDropAward(droppool.DropAward(def.DropPoolID))

	userTrial.StageReward(msg.GetStageID())

}

func init() {
	user.RegisterTransFunc(cmdEnum.CS_TrialStageReward, func(user *user.User, msg *codecs.Message) transaction.Transaction {
		return &transactionTrialStageDungeon{
			user: user,
			req:  msg,
		}
	})
}
