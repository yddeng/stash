package Character

import (
	codecs "initialthree/codec/cs"
	"initialthree/node/common/enumType"
	"initialthree/node/common/inoutput"
	"initialthree/node/common/transaction"
	"initialthree/zaplogger"

	"initialthree/node/node_game/module"
	"initialthree/node/node_game/module/character"
	"initialthree/node/node_game/user"
	"initialthree/node/table/excel/ConstTable/Global"
	"initialthree/node/table/excel/DataTable/Gift"
	"initialthree/node/table/excel/DataTable/PlayerCharacter"
	"initialthree/protocol/cmdEnum"
	cs_message "initialthree/protocol/cs/message"
)

type transactionCharacterGift struct {
	transaction.TransactionBase
	user    *user.User
	req     *codecs.Message
	errcode cs_message.ErrCode
}

func (this *transactionCharacterGift) GetModuleName() string {
	return "Character"
}

func (this *transactionCharacterGift) Begin() {

	defer func() { this.EndTrans(&cs_message.CharacterLevelUpToC{}, this.errcode) }()

	msg := this.req.GetData().(*cs_message.CharacterGiftToS)

	zaplogger.GetSugar().Infof("%s %d Call CharacterGift %v", this.user.GetUserID(), this.user.GetID(), msg)

	userCharacter := this.user.GetSubModule(module.Character).(*character.UserCharacter)

	chara := userCharacter.GetCharacter(msg.GetCharacterID())
	CharaDef := PlayerCharacter.GetID(msg.GetCharacterID())
	if chara == nil || CharaDef == nil {
		zaplogger.GetSugar().Infof("%v %v %s", this.user.GetUserID(), this.user.GetID(), "CharacterGift fail:chara is nil")
		this.errcode = cs_message.ErrCode_Character_NotExist
		return
	}

	favorExps := Global.Get().FavorExpArray
	if chara.FavorLevel-1 >= int32(len(favorExps)) {
		zaplogger.GetSugar().Infof("%v %v %s", this.user.GetUserID(), this.user.GetID(), "CharacterGift fail:chara FavorLevel is max")
		this.errcode = cs_message.ErrCode_Level_High
		return
	}

	totalExp := int32(0)
	for _, v := range msg.GetCostItems() {
		if v.GetCount() <= 0 {
			this.errcode = cs_message.ErrCode_Request_Argument_Err
			return
		}

		giftDef := Gift.GetID(v.GetItemID())
		if giftDef == nil {
			zaplogger.GetSugar().Infof("%s CharacterGift fail:itemID %d giftDef is nil", this.user.GetUserLogName(), v.GetItemID())
			this.errcode = cs_message.ErrCode_Config_NotExist
			return
		}

		exp := giftDef.BaseExp
		if giftDef.GiftTypeEnum == CharaDef.GiftTypeEnum {
			exp += giftDef.BonusExp
		}
		totalExp += exp * v.GetCount()
	}

	if totalExp <= 0 {
		zaplogger.GetSugar().Debugf("%v %v %s", this.user.GetUserID(), this.user.GetID(), "CharacterGift totalExp = 0")
		this.errcode = cs_message.ErrCode_Request_Argument_Err
		return
	}

	// 扣除消耗
	useRes := make([]inoutput.ResDesc, 0, len(msg.GetCostItems()))
	for _, v := range msg.GetCostItems() {
		useRes = append(useRes, inoutput.ResDesc{ID: v.GetItemID(), Type: enumType.IOType_Item, Count: v.GetCount()})
	}
	if this.errcode = this.user.DoInputOutput(useRes, nil, false); this.errcode != cs_message.ErrCode_OK {
		zaplogger.GetSugar().Debugf("%s CharacterGift error: inouput %s  ", this.user.GetUserLogName(), this.errcode.String())
		return
	}

	// 总和exp
	userCharacter.FavorLevelUp(chara, totalExp)

	zaplogger.GetSugar().Infof("%v %v %s", this.user.GetUserID(), this.user.GetID(), "CharacterGift OK")
	this.errcode = cs_message.ErrCode_OK
}

func init() {
	user.RegisterTransFunc(cmdEnum.CS_CharacterGift, func(user *user.User, msg *codecs.Message) transaction.Transaction {
		return &transactionCharacterGift{
			user: user,
			req:  msg,
		}
	})
}
