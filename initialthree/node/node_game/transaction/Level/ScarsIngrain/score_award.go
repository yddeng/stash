package ScarsIngrain

import (
	codecs "initialthree/codec/cs"
	"initialthree/node/common/droppool"
	"initialthree/node/common/transaction"
	"initialthree/zaplogger"

	"initialthree/node/node_game/module"
	"initialthree/node/node_game/module/scarsIngrain"
	"initialthree/node/node_game/user"
	"initialthree/node/table/excel/DataTable/ScarsIngrainScoreReward"
	"initialthree/protocol/cmdEnum"
	cs_message "initialthree/protocol/cs/message"
)

type transactionScarsIngrainGetScoreAward struct {
	transaction.TransactionBase
	user    *user.User
	req     *codecs.Message
	errcode cs_message.ErrCode
	resp    *cs_message.ScarsIngrainGetScoreAwardToC
}

func (this *transactionScarsIngrainGetScoreAward) GetModuleName() string {
	return "ScarsIngrain"
}

func (this *transactionScarsIngrainGetScoreAward) Begin() {
	defer func() { this.EndTrans(this.resp, this.errcode) }()
	this.resp = &cs_message.ScarsIngrainGetScoreAwardToC{}
	msg := this.req.GetData().(*cs_message.ScarsIngrainGetScoreAwardToS)
	zaplogger.GetSugar().Infof("%s ScarsIngrainGetScoreAwardToS %v ", this.user.GetUserLogName(), msg)

	siModule := this.user.GetSubModule(module.ScarsIngrain).(*scarsIngrain.ScarsIngrain)
	siData := siModule.GetData()

	def := ScarsIngrainScoreReward.GetID(siData.SIID)
	if def == nil {
		zaplogger.GetSugar().Infof("%s ScarsIngrainGetScoreAwardToS failed, %d is nil", this.user.GetUserLogName(), siData.SIID)
		this.errcode = cs_message.ErrCode_Config_NotExist
		return
	}

	score := msg.GetScore()
	droppoolId, ok := def.GetScoreAward(score)
	if !ok {
		zaplogger.GetSugar().Infof("%s ScarsIngrainGetScoreAwardToS failed, %d -> score %d is nil", this.user.GetUserLogName(), siData.SIID, score)
		this.errcode = cs_message.ErrCode_Config_NotExist
		return
	}

	// 是否可以领取分数
	if score > siData.TotalScore {
		zaplogger.GetSugar().Infof("%s ScarsIngrainGetScoreAwardToS failed,score %d  < totalScore %d", this.user.GetUserLogName(), score, siData.TotalScore)
		this.errcode = cs_message.ErrCode_ScarsIngrain_GetAwardScoreErr
		return
	}

	// 是否已经领取过
	_, ok = siData.ScoreAward[score]
	if ok {
		zaplogger.GetSugar().Infof("%s ScarsIngrainGetScoreAwardToS failed,score %d is already get", this.user.GetUserLogName(), score)
		this.errcode = cs_message.ErrCode_ScarsIngrain_GetAwardScoreErr
		return
	}

	award := droppool.DropAward(droppoolId)
	_ = this.user.ApplyDropAward(award)

	this.resp.Award = award.ToMessageAward()

	siModule.SetScoreAwardGet(score)

	this.errcode = cs_message.ErrCode_OK
}

func init() {
	user.RegisterTransFunc(cmdEnum.CS_ScarsIngrainGetScoreAward, func(user *user.User, msg *codecs.Message) transaction.Transaction {
		return &transactionScarsIngrainGetScoreAward{
			user: user,
			req:  msg,
		}
	})
}
