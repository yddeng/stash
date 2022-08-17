package ScarsIngrain

import (
	codecs "initialthree/codec/cs"
	"initialthree/node/common/transaction"
	"initialthree/node/node_game/module/character"
	"initialthree/zaplogger"

	"initialthree/node/node_game/module"
	"initialthree/node/node_game/module/scarsIngrain"
	"initialthree/node/node_game/user"
	"initialthree/protocol/cmdEnum"
	cs_message "initialthree/protocol/cs/message"
)

type transactionScarsIngrainTeamPrefabSet struct {
	transaction.TransactionBase
	user    *user.User
	req     *codecs.Message
	errcode cs_message.ErrCode
	resp    *cs_message.ScarsIngrainTeamPrefabSetToC
}

func (this *transactionScarsIngrainTeamPrefabSet) GetModuleName() string {
	return "ScarsIngrain"
}

func (this *transactionScarsIngrainTeamPrefabSet) Begin() {
	defer func() { this.EndTrans(this.resp, this.errcode) }()
	this.resp = &cs_message.ScarsIngrainTeamPrefabSetToC{}
	msg := this.req.GetData().(*cs_message.ScarsIngrainTeamPrefabSetToS)
	zaplogger.GetSugar().Infof("%s ScarsIngrainTeamPrefabSetToS %v ", this.user.GetUserLogName(), msg)

	siModule := this.user.GetSubModule(module.ScarsIngrain).(*scarsIngrain.ScarsIngrain)
	charaModule := this.user.GetSubModule(module.Character).(*character.UserCharacter)

	for _, id := range msg.GetCharacterList() {
		if id != 0 {
			if nil == charaModule.GetCharacter(id) {
				this.errcode = cs_message.ErrCode_Config_NotExist
				return
			}
		}
	}

	if !siModule.BossTeamPrefabSet(msg.GetBossID(), msg.GetCharacterList()) {
		this.errcode = cs_message.ErrCode_Request_Argument_Err
		return
	}

	this.errcode = cs_message.ErrCode_OK
}

func init() {
	user.RegisterTransFunc(cmdEnum.CS_ScarsIngrainTeamPrefabSet, func(user *user.User, msg *codecs.Message) transaction.Transaction {
		return &transactionScarsIngrainTeamPrefabSet{
			user: user,
			req:  msg,
		}
	})
}
