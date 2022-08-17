package Character

import (
	codecs "initialthree/codec/cs"
	"initialthree/node/common/enumType"
	"initialthree/node/common/inoutput"
	"initialthree/node/common/transaction"
	"initialthree/node/table/excel/DataTable/PlayerGene"
	"initialthree/zaplogger"

	"initialthree/node/node_game/module"
	"initialthree/node/node_game/module/character"
	"initialthree/node/node_game/user"
	"initialthree/protocol/cmdEnum"
	cs_message "initialthree/protocol/cs/message"
)

type transactionCharacterGeneLevelUp struct {
	transaction.TransactionBase
	user    *user.User
	req     *codecs.Message
	resp    *cs_message.CharacterGeneLevelUpToC
	errcode cs_message.ErrCode
}

func (this *transactionCharacterGeneLevelUp) GetModuleName() string {
	return "Character"
}

func (this *transactionCharacterGeneLevelUp) Begin() {
	defer func() { this.EndTrans(this.resp, this.errcode) }()
	this.resp = &cs_message.CharacterGeneLevelUpToC{}
	msg := this.req.GetData().(*cs_message.CharacterGeneLevelUpToS)

	zaplogger.GetSugar().Infof("%s Call CharacterGeneLevelUp %v", this.user.GetUserLogName(), msg)

	userCharacter := this.user.GetSubModule(module.Character).(*character.UserCharacter)

	chara := userCharacter.GetCharacter(msg.GetCharacterID())
	if chara == nil {
		zaplogger.GetSugar().Infof("%v %v CharacterGeneLevelUp fail:chara is nil", this.user.GetUserID(), this.user.GetID())
		this.errcode = cs_message.ErrCode_Character_NotExist
		return
	}

	geneDef := PlayerGene.GetGene(msg.GetCharacterID(), chara.GeneLevel+1)
	if geneDef == nil {
		zaplogger.GetSugar().Infof("%v %v CharacterGeneLevelUp fail:geneDef is nil", this.user.GetUserID(), this.user.GetID())
		this.errcode = cs_message.ErrCode_Config_NotExist
		return
	}

	useRes := []inoutput.ResDesc{{ID: geneDef.ItemID, Type: enumType.IOType_Item, Count: geneDef.ItemCount}}
	if this.errcode = this.user.DoInputOutput(useRes, nil, false); this.errcode != cs_message.ErrCode_OK {
		zaplogger.GetSugar().Debugf("%s CharacterGeneLevelUp error: inouput %s  ", this.user.GetUserLogName(), this.errcode.String())
		return
	}

	userCharacter.GeneLevelUp(chara)

	zaplogger.GetSugar().Infof("%v %v CharacterGeneLevelUp OK", this.user.GetUserID(), this.user.GetID())
	this.errcode = cs_message.ErrCode_OK
}

func init() {
	user.RegisterTransFunc(cmdEnum.CS_CharacterGeneLevelUp, func(user *user.User, msg *codecs.Message) transaction.Transaction {
		return &transactionCharacterGeneLevelUp{
			user: user,
			req:  msg,
		}
	})
}
