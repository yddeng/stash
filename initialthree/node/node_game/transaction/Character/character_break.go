package Character

import (
	codecs "initialthree/codec/cs"
	"initialthree/node/common/attr"
	"initialthree/node/common/enumType"
	"initialthree/node/common/inoutput"
	"initialthree/node/common/transaction"
	"initialthree/node/node_game/module"
	"initialthree/node/node_game/module/character"
	"initialthree/node/node_game/user"
	"initialthree/node/table/excel/DataTable/CharacterBreakThrough"
	"initialthree/node/table/excel/DataTable/PlayerCharacter"
	"initialthree/protocol/cmdEnum"
	cs_message "initialthree/protocol/cs/message"
	"initialthree/zaplogger"
)

type transactionCharacterBreak struct {
	transaction.TransactionBase
	user    *user.User
	req     *codecs.Message
	resp    *cs_message.CharacterBreakToC
	errcode cs_message.ErrCode
}

func (this *transactionCharacterBreak) GetModuleName() string {
	return "Character"
}

func (this *transactionCharacterBreak) Begin() {
	defer func() { this.EndTrans(this.resp, this.errcode) }()
	this.resp = &cs_message.CharacterBreakToC{}
	msg := this.req.GetData().(*cs_message.CharacterBreakToS)

	zaplogger.GetSugar().Infof("%s %d Call CharacterBreakToS id:%d", this.user.GetUserID(), this.user.GetID(), msg.GetCharacterID())

	userCharacter := this.user.GetSubModule(module.Character).(*character.UserCharacter)

	chara := userCharacter.GetCharacter(msg.GetCharacterID())
	charaDef := PlayerCharacter.GetID(msg.GetCharacterID())
	if chara == nil || charaDef == nil {
		zaplogger.GetSugar().Info(this.user.GetUserID(), this.user.GetID(), "CharacterBreak fail:chara is nil")
		this.errcode = cs_message.ErrCode_Character_NotExist
		return
	}

	nextBreakDef := CharacterBreakThrough.GetID(charaDef.GetBreakID(chara.BreakLevel + 1))
	if nextBreakDef == nil {
		zaplogger.GetSugar().Info(this.user.GetUserID(), this.user.GetID(), "CharacterBreak fail:nextBreakDef is nil")
		this.errcode = cs_message.ErrCode_Config_NotExist
		return
	}

	// 等级
	if chara.Level < nextBreakDef.LevelRequirement {
		zaplogger.GetSugar().Info(this.user.GetUserID(), this.user.GetID(), "CharacterBreak fail: level < LevelRequirement")
		this.errcode = cs_message.ErrCode_Character_Level_Low
		return
	}

	useRes := make([]inoutput.ResDesc, 0, len(nextBreakDef.SpecifiedIDListArray))
	useRes = append(useRes, inoutput.ResDesc{ID: attr.Gold, Type: enumType.IOType_UsualAttribute, Count: nextBreakDef.GoldCost})
	for _, v := range nextBreakDef.SpecifiedIDListArray {
		useRes = append(useRes, inoutput.ResDesc{ID: v.ID, Type: enumType.IOType_Item, Count: v.Count})
	}
	if this.errcode = this.user.DoInputOutput(useRes, nil, false); this.errcode != cs_message.ErrCode_OK {
		zaplogger.GetSugar().Debugf("%s CharacterBreak error: inouput %s  ", this.user.GetUserLogName(), this.errcode.String())
		return
	}

	userCharacter.BreakLevel(chara)

	zaplogger.GetSugar().Info(this.user.GetUserID(), this.user.GetID(), "CharacterBreak OK")
	this.errcode = cs_message.ErrCode_OK
}

func init() {
	user.RegisterTransFunc(cmdEnum.CS_CharacterBreak, func(user *user.User, msg *codecs.Message) transaction.Transaction {
		return &transactionCharacterBreak{
			user: user,
			req:  msg,
		}
	})
}
