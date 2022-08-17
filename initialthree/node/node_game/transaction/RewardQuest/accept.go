package RewardQuest

import (
	codecs "initialthree/codec/cs"
	"initialthree/node/common/timeDisposal"
	"initialthree/node/common/transaction"
	"initialthree/zaplogger"

	"initialthree/node/node_game/module"
	"initialthree/node/node_game/module/character"
	"initialthree/node/node_game/module/rewardQuest"
	"initialthree/node/node_game/user"
	"initialthree/node/table/excel/DataTable/PlayerCharacter"
	"initialthree/node/table/excel/DataTable/RewardQuest"
	"initialthree/protocol/cmdEnum"
	cs_message "initialthree/protocol/cs/message"
)

type transactionRewardQuestAccept struct {
	transaction.TransactionBase
	user    *user.User
	req     *codecs.Message
	resp    *cs_message.RewardQuestAcceptToC
	errcode cs_message.ErrCode
}

func (this *transactionRewardQuestAccept) GetModuleName() string {
	return "RewardQuest"
}

func (this *transactionRewardQuestAccept) Begin() {

	defer func() { this.EndTrans(this.resp, this.errcode) }()

	msg := this.req.GetData().(*cs_message.RewardQuestAcceptToS)

	zaplogger.GetSugar().Infof("%s %d Call RewardQuestAcceptToS %v", this.user.GetUserID(), this.user.GetID(), msg)

	userQuest := this.user.GetSubModule(module.RewardQuest).(*rewardQuest.RewardQuest)

	q := userQuest.GetRewardQuest(msg.GetQuestID())
	if q == nil {
		zaplogger.GetSugar().Debugf("%s RewardQuestAcceptToS failed, rewardQuest %d is not exist", this.user.GetUserLogName(), msg.GetQuestID())
		this.errcode = cs_message.ErrCode_Quest_NotExist
		return
	}

	if q.GetState() != cs_message.QuestState_Acceptable {
		zaplogger.GetSugar().Debugf("%s RewardQuestAcceptToS failed, rewardQuest %d state failed", this.user.GetUserLogName(), msg.GetQuestID())
		this.errcode = cs_message.ErrCode_Quest_StateErr
		return
	}

	def := RewardQuest.GetID(msg.GetQuestID())
	if def == nil {
		zaplogger.GetSugar().Debugf("%s RewardQuestAcceptToS failed, rewardQuest config %d is not exist", this.user.GetUserLogName(), msg.GetQuestID())
		this.errcode = cs_message.ErrCode_Config_NotExist
		return
	}

	ids := msg.GetCharacters()
	if len(ids) == 0 || len(ids) > 3 {
		zaplogger.GetSugar().Debugf("%s RewardQuestAcceptToS failed, characters length is error", this.user.GetUserLogName(), len(ids))
		this.errcode = cs_message.ErrCode_Request_Argument_Err
		return
	}

	usedRoles := userQuest.UsedRoles()

	charaModule := this.user.GetSubModule(module.Character).(*character.UserCharacter)
	var attack, quality []int32
	for _, id := range ids {
		chara := charaModule.GetCharacter(id)
		if chara == nil {
			zaplogger.GetSugar().Debugf("%s RewardQuestAcceptToS failed, character %d is not exist", this.user.GetUserLogName(), id)
			this.errcode = cs_message.ErrCode_Character_NotExist
			return
		}

		if _, ok := usedRoles[id]; ok {
			zaplogger.GetSugar().Debugf("%s RewardQuestAcceptToS failed, character %d is already used ", this.user.GetUserLogName(), id)
			this.errcode = cs_message.ErrCode_Request_Argument_Err
			return
		}

		charaDef := PlayerCharacter.GetID(id)
		if charaDef == nil {
			zaplogger.GetSugar().Debugf("%s RewardQuestAcceptToS failed, character config %d is not exist", this.user.GetUserLogName(), id)
			this.errcode = cs_message.ErrCode_Config_NotExist
			return
		}

		attack = append(attack, charaDef.GetDamageElementType())
		quality = append(quality, charaDef.RarityEnum)
	}

	ok := def.IsAccept(attack, quality)
	if !ok {
		zaplogger.GetSugar().Debugf("%s RewardQuestAcceptToS failed, condition is'not enough", this.user.GetUserLogName())
		this.errcode = cs_message.ErrCode_Attr_Low
		return
	}

	userQuest.Accept(q, ids, timeDisposal.NowUnix())
	zaplogger.GetSugar().Infof("%s RewardQuestAcceptToS ok", this.user.GetUserLogName())
}

func init() {
	user.RegisterTransFunc(cmdEnum.CS_RewardQuestAccept, func(user *user.User, msg *codecs.Message) transaction.Transaction {
		return &transactionRewardQuestAccept{
			user:    user,
			req:     msg,
			resp:    &cs_message.RewardQuestAcceptToC{},
			errcode: cs_message.ErrCode_OK,
		}
	})
}
