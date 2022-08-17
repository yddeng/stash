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
	"initialthree/node/table/excel/ConstTable/Global"
	"initialthree/node/table/excel/DataTable/CharacterLevelUpExp"
	"initialthree/node/table/excel/DataTable/PlayerCharacter"
	"initialthree/protocol/cmdEnum"
	cs_message "initialthree/protocol/cs/message"
	"initialthree/zaplogger"
)

type transactionCharacterLevelUp struct {
	transaction.TransactionBase
	user    *user.User
	req     *codecs.Message
	resp    *cs_message.CharacterLevelUpToC
	errcode cs_message.ErrCode
}

func (this *transactionCharacterLevelUp) GetModuleName() string {
	return "Character"
}

func (this *transactionCharacterLevelUp) Begin() {

	defer func() { this.EndTrans(this.resp, this.errcode) }()

	this.resp = &cs_message.CharacterLevelUpToC{}
	msg := this.req.GetData().(*cs_message.CharacterLevelUpToS)

	zaplogger.GetSugar().Infof("%s %d Call CharacterLevelUp %v", this.user.GetUserID(), this.user.GetID(), msg)

	userCharacter := this.user.GetSubModule(module.Character).(*character.UserCharacter)

	chara := userCharacter.GetCharacter(msg.GetCharacterID())
	charaDef := PlayerCharacter.GetID(msg.GetCharacterID())
	if chara == nil || charaDef == nil {
		zaplogger.GetSugar().Infof("%v %v %s", this.user.GetUserID(), this.user.GetID(), "CharacterLevelUp fail:chara is nil")
		this.errcode = cs_message.ErrCode_Character_NotExist
		return
	}

	maxLevel := userCharacter.GetMaxLevel(chara)

	if chara.Level >= maxLevel {
		zaplogger.GetSugar().Infof("%v %v %s", this.user.GetUserID(), this.user.GetID(), "CharacterLevelUp fail:chara level >= maxLevel")
		this.errcode = cs_message.ErrCode_Level_Low
		return
	}

	if _, exist := CharacterLevelUpExp.GetMaxExp(charaDef.RarityEnum, chara.Level+1); !exist {
		zaplogger.GetSugar().Infof("%v %v %s", this.user.GetUserID(), this.user.GetID(), "CharacterLevelUp fail:nextLevelDef is nil")
		this.errcode = cs_message.ErrCode_Config_NotExist
		return
	}

	totalExp := int32(0)
	for _, v := range msg.GetCostItems() {
		if v.GetCount() <= 0 {
			this.errcode = cs_message.ErrCode_Request_Argument_Err
			return
		}
		exp := this.getCostItemExp(v.GetItemID())
		if exp <= 0 {
			zaplogger.GetSugar().Infof("%s CharacterLevelUp fail:itemID %d is not in LeveUpCostItemsAndItOfferedExpArray", this.user.GetUserLogName(), v.GetItemID())
			this.errcode = cs_message.ErrCode_Character_LevelUp_ItemError
			return
		}
		totalExp += exp * v.GetCount()
	}

	if totalExp <= 0 {
		zaplogger.GetSugar().Debugf("%v %v %s", this.user.GetUserID(), this.user.GetID(), "CharacterLevelUp totalExp = 0")
		this.errcode = cs_message.ErrCode_Character_LevelUp_ExpNil
		return
	}

	usedExp, level, currentExp := userCharacter.CalcUseExp(chara, totalExp, maxLevel)

	// 扣除消耗
	useRes := make([]inoutput.ResDesc, 0, len(msg.GetCostItems())+1)
	// 根据消耗经验计算消耗金币
	expCostGold := Global.Get().CharacterLvUpExpCostGold
	useRes = append(useRes, inoutput.ResDesc{ID: attr.Gold, Type: enumType.IOType_UsualAttribute, Count: expCostGold * usedExp})
	for _, v := range msg.GetCostItems() {
		useRes = append(useRes, inoutput.ResDesc{ID: v.GetItemID(), Type: enumType.IOType_Item, Count: v.GetCount()})
	}
	if this.errcode = this.user.DoInputOutput(useRes, nil, false); this.errcode != cs_message.ErrCode_OK {
		zaplogger.GetSugar().Debugf("%s CharacterLevelUp error: inouput %s  ", this.user.GetUserLogName(), this.errcode.String())
		return
	}

	userCharacter.SetLevel(chara, level, currentExp)

	// 是否需要触发重新计算战斗属性
	//userCharacter.EventRecalculate(chara)

	zaplogger.GetSugar().Infof("%v %v %s", this.user.GetUserID(), this.user.GetID(), "CharacterLevelUp OK")
	this.errcode = cs_message.ErrCode_OK
}

func (this *transactionCharacterLevelUp) getCostItemExp(itemID int32) int32 {
	def := Global.GetID(1)
	for _, v := range def.LevepUpCostItemsAndItOfferedExpArray {
		if v.ItemID == itemID {
			return v.Exp
		}
	}
	return 0
}

func init() {
	user.RegisterTransFunc(cmdEnum.CS_CharacterLevelUp, func(user *user.User, msg *codecs.Message) transaction.Transaction {
		return &transactionCharacterLevelUp{
			user: user,
			req:  msg,
		}
	})
}
