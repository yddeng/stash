package Character

import (
	codecs "initialthree/codec/cs"
	"initialthree/node/common/transaction"
	"initialthree/node/node_game/module/character"
	"initialthree/zaplogger"

	"initialthree/node/node_game/module"
	"initialthree/node/node_game/user"
	"initialthree/protocol/cmdEnum"
	cs_message "initialthree/protocol/cs/message"
)

type transactionCharacterLike struct {
	transaction.TransactionBase
	user    *user.User
	req     *codecs.Message
	resp    *cs_message.CharacterLikeToC
	errcode cs_message.ErrCode
}

func (this *transactionCharacterLike) GetModuleName() string {
	return "Character"
}

func (this *transactionCharacterLike) Begin() {
	defer func() { this.EndTrans(this.resp, this.errcode) }()
	this.resp = &cs_message.CharacterLikeToC{}
	msg := this.req.GetData().(*cs_message.CharacterLikeToS)

	zaplogger.GetSugar().Infof("%s Call CharacterLikeToS %v", this.user.GetUserLogName(), msg)

	userCharacter := this.user.GetSubModule(module.Character).(*character.UserCharacter)
	doCharacter := userCharacter.GetCharacter(msg.GetCharacterID())
	if doCharacter == nil {
		zaplogger.GetSugar().Debugf("%s CharacterDemountToS error:  equip %d is nil", this.user.GetUserLogName(), msg.GetCharacterID())
		this.errcode = cs_message.ErrCode_Character_NotExist
		return
	}

	userCharacter.IsLike(doCharacter, msg.GetIsLike())
	this.errcode = cs_message.ErrCode_OK
}

func init() {
	user.RegisterTransFunc(cmdEnum.CS_CharacterLike, func(user *user.User, msg *codecs.Message) transaction.Transaction {
		return &transactionCharacterLike{
			user: user,
			req:  msg,
		}
	})
}
