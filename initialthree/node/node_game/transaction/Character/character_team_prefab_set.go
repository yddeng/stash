package Character

import (
	codecs "initialthree/codec/cs"
	"initialthree/node/common/transaction"
	"initialthree/zaplogger"

	"initialthree/node/node_game/module"
	"initialthree/node/node_game/module/character"
	"initialthree/node/node_game/user"
	"initialthree/protocol/cmdEnum"
	cs_message "initialthree/protocol/cs/message"
)

type transactionCharacterTeamPrefabSet struct {
	transaction.TransactionBase
	user    *user.User
	req     *codecs.Message
	resp    *cs_message.CharacterTeamPrefabSetToC
	errcode cs_message.ErrCode
}

func (this *transactionCharacterTeamPrefabSet) GetModuleName() string {
	return "Character"
}

func (this *transactionCharacterTeamPrefabSet) Begin() {
	defer func() { this.EndTrans(this.resp, this.errcode) }()
	this.resp = &cs_message.CharacterTeamPrefabSetToC{}
	msg := this.req.GetData().(*cs_message.CharacterTeamPrefabSetToS)

	zaplogger.GetSugar().Infof("%s %d Call CharacterTeamPrefabSetToS %v", this.user.GetUserID(), this.user.GetID(), msg)

	idx := msg.GetTeamPrefabIdx()
	if idx < 0 || idx >= 9 {
		zaplogger.GetSugar().Debugf("%s CharacterTeamPrefabSetToS idx:%d is failed", this.user.GetUserLogName(), idx)
		this.errcode = cs_message.ErrCode_Character_TeamIndexFailed
		return
	}

	if msg.GetPrefab() == nil || len(msg.GetPrefab().GetCharacterList()) != 3 {
		zaplogger.GetSugar().Debugf("%s CharacterTeamSetToS character len = %d ", this.user.GetUserLogName(), len(msg.GetPrefab().GetCharacterList()))
		this.errcode = cs_message.ErrCode_Character_TeamRoleCountErr
		return
	}

	userCharacter := this.user.GetSubModule(module.Character).(*character.UserCharacter)

	ids := map[int32]struct{}{}
	for _, id := range msg.GetPrefab().GetCharacterList() {
		if id != 0 {
			if _, ok := ids[id]; ok {
				zaplogger.GetSugar().Debugf("%s CharacterTeamPrefabSetToS character:%d repeated set", this.user.GetUserLogName(), id)
				this.errcode = cs_message.ErrCode_Character_TeamRepeated
				return
			}
			ids[id] = struct{}{}

			c := userCharacter.GetCharacter(id)
			if c == nil {
				zaplogger.GetSugar().Debugf("%s CharacterTeamPrefabSetToS character:%d not exist", this.user.GetUserLogName(), id)
				this.errcode = cs_message.ErrCode_Character_NotExist
				return
			}
		}
	}

	userCharacter.GroupDefSet(int(idx), msg.GetPrefab())

	zaplogger.GetSugar().Debugf("%s CharacterTeamPrefabSetToS OK", this.user.GetUserLogName())
	this.errcode = cs_message.ErrCode_OK
}

func init() {
	user.RegisterTransFunc(cmdEnum.CS_CharacterTeamPrefabSet, func(user *user.User, msg *codecs.Message) transaction.Transaction {
		return &transactionCharacterTeamPrefabSet{
			user: user,
			req:  msg,
		}
	})
}
