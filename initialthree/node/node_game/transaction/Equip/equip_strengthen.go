package Equip

import (
	codecs "initialthree/codec/cs"
	"initialthree/node/common/attr"
	"initialthree/node/common/enumType"
	"initialthree/node/common/inoutput"
	"initialthree/node/common/transaction"
	"initialthree/node/node_game/module"
	"initialthree/node/node_game/module/backpack"
	"initialthree/node/node_game/module/equip"
	"initialthree/node/node_game/user"
	constEquip "initialthree/node/table/excel/ConstTable/Equip"
	"initialthree/node/table/excel/DataTable/Equip"
	"initialthree/node/table/excel/DataTable/EquipLevelMaxExp"
	"initialthree/node/table/excel/DataTable/EquipQuality"
	"initialthree/node/table/excel/DataTable/Item"
	"initialthree/protocol/cmdEnum"
	cs_message "initialthree/protocol/cs/message"
	"initialthree/zaplogger"
)

type transactionEquipStrengthen struct {
	transaction.TransactionBase
	user    *user.User
	req     *codecs.Message
	resp    *cs_message.EquipStrengthenToC
	errcode cs_message.ErrCode
}

func (this *transactionEquipStrengthen) GetModuleName() string {
	return "Equip"
}

func (this *transactionEquipStrengthen) Begin() {
	defer func() { this.EndTrans(this.resp, this.errcode) }()
	this.resp = &cs_message.EquipStrengthenToC{}
	msg := this.req.GetData().(*cs_message.EquipStrengthenToS)

	zaplogger.GetSugar().Infof("%s Call EquipStrengthenToS %v", this.user.GetUserLogName(), msg)

	userEquip := this.user.GetSubModule(module.Equip).(*equip.UserEquip)
	doEquip := userEquip.GetEquip(msg.GetEquipID())
	if doEquip == nil {
		zaplogger.GetSugar().Debugf("%s EquipStrengthenToS error: mod is nil", this.user.GetUserLogName())
		this.errcode = cs_message.ErrCode_Equip_NotExist
		return
	}

	def := Equip.GetID(doEquip.ConfigID)
	if def == nil {
		zaplogger.GetSugar().Debugf("%s EquipStrengthenToS error: %d modDef is nil", this.user.GetUserLogName(), doEquip.ConfigID)
		this.errcode = cs_message.ErrCode_Config_NotExist
		return
	}

	// 请求参数验证
	if len(msg.GetCostItems()) == 0 && len(msg.GetCostEquipIDs()) == 0 {
		zaplogger.GetSugar().Debugf("%s EquipStrengthenToS error: costEquipIDs is nil", this.user.GetUserLogName())
		this.errcode = cs_message.ErrCode_Request_Argument_Err
		return
	}

	levelDef := EquipLevelMaxExp.GetID(def.LevelMaxExp)
	qualityDef := EquipQuality.GetID(def.QualityEnum)
	if levelDef == nil || qualityDef == nil {
		zaplogger.GetSugar().Debugf("%s EquipStrengthenToS error: %d levelDef or qualityDef is nil", this.user.GetUserLogName(), msg.GetEquipID())
		this.errcode = cs_message.ErrCode_Config_NotExist
		return
	}

	maxLevel := qualityDef.MaxLevel
	if levelDef.MaxLevel() < maxLevel {
		maxLevel = levelDef.MaxLevel()
	}

	// 等级到达上限
	if doEquip.Level >= maxLevel {
		zaplogger.GetSugar().Debugf("%s EquipStrengthenToS error: %d level %d max ", this.user.GetUserLogName(), msg.GetEquipID(), maxLevel)
		this.errcode = cs_message.ErrCode_Equip_LevelMax
		return
	}

	userBackpack := this.user.GetSubModule(module.Backpack).(*backpack.Backpack)
	totalExp := int32(0)
	var useRes = make([]inoutput.ResDesc, 0, len(msg.GetCostEquipIDs())+len(msg.GetCostItems())+1)

	for _, id := range msg.GetCostEquipIDs() {
		if id == msg.GetEquipID() {
			zaplogger.GetSugar().Debugf("%s EquipStrengthenToS error: the consumption is the same as the operation %d ", this.user.GetUserLogName(), msg.GetEquipID())
			this.errcode = cs_message.ErrCode_Request_Argument_Err
			return
		}
		e := userEquip.GetEquip(id)
		if e == nil {
			zaplogger.GetSugar().Debugf("%s EquipStrengthenToS error: cost mod %d is nil", this.user.GetUserLogName(), id)
			this.errcode = cs_message.ErrCode_Equip_NotExist
			return
		}
		// 装备已被装备
		if e.IsLock || e.EquipCharacterId != 0 {
			zaplogger.GetSugar().Debugf("%s EquipStrengthenToS error: mod %d is equipped or locked", this.user.GetUserLogName(), id)
			this.errcode = cs_message.ErrCode_Request_Argument_Err
			return
		}

		def2 := Equip.GetID(e.ConfigID)
		if def2 == nil {
			zaplogger.GetSugar().Debugf("%s EquipStrengthenToS error: cost mod %d def is nil", this.user.GetUserLogName(), id)
			this.errcode = cs_message.ErrCode_Equip_NotExist
			return
		}

		levelDef2 := EquipLevelMaxExp.GetID(def.LevelMaxExp)
		if levelDef2 == nil {
			zaplogger.GetSugar().Debugf("%s EquipStrengthenToS error: cost %d levelDef is nil", this.user.GetUserLogName(), def.LevelMaxExp)
			this.errcode = cs_message.ErrCode_Config_NotExist
			return
		}
		totalExp += def2.SupplyExp + levelDef2.GetRangeTotalExp(1, e.Level) + e.Exp
		useRes = append(useRes, inoutput.ResDesc{Type: enumType.IOType_Equip, ID: int32(id), Count: 1})
	}

	for _, v := range msg.GetCostItems() {
		if v.GetCount() <= 0 || v.GetCount() > userBackpack.GetItemCountByTID(v.GetItemID()) {
			zaplogger.GetSugar().Debugf("%s EquipStrengthenToS error: costItem %d count %d is not enough nil", this.user.GetUserLogName(), v.GetItemID(), v.GetCount())
			this.errcode = cs_message.ErrCode_Item_Not_Enough
			return
		}
		itemDef := Item.GetID(v.GetItemID())
		if itemDef == nil {
			zaplogger.GetSugar().Debugf("%s EquipStrengthenToS error: item %d itemDef is nil", this.user.GetUserLogName(), v.GetItemID())
			this.errcode = cs_message.ErrCode_Equip_NotExist
			return
		}
		supplyExp := constEquip.GetSupplyExp(v.GetItemID())
		if supplyExp == 0 {
			zaplogger.GetSugar().Debugf("%s EquipStrengthenToS error: item %d supplyExp is 0", this.user.GetUserLogName(), v.GetItemID())
			this.errcode = cs_message.ErrCode_Equip_NotExist
			return
		}
		totalExp += supplyExp * v.GetCount()
		useRes = append(useRes, inoutput.ResDesc{Type: enumType.IOType_Item, ID: v.GetItemID(), Count: v.GetCount()})
	}

	if totalExp <= 0 {
		zaplogger.GetSugar().Debugf("%s EquipStrengthenToS error: totalExp = 0", this.user.GetUserLogName())
		this.errcode = cs_message.ErrCode_Resource_NotEnough
		return
	}

	useGold := int32(constEquip.Get().ExpToGoldRate * float64(totalExp))
	useRes = append(useRes, inoutput.ResDesc{Type: enumType.IOType_UsualAttribute, ID: attr.Gold, Count: useGold})
	// 消耗
	if this.errcode = this.user.DoInputOutput(useRes, nil, false); this.errcode != cs_message.ErrCode_OK {
		zaplogger.GetSugar().Debugf("%s EquipStrengthenToS error: inouput %s  ", this.user.GetUserLogName(), this.errcode.String())
		return
	}

	userEquip.Strengthen(levelDef, doEquip, maxLevel, totalExp)

	this.errcode = cs_message.ErrCode_OK
	zaplogger.GetSugar().Infof("%s Call EquipStrengthenToS ok", this.user.GetUserLogName())
}

func init() {
	user.RegisterTransFunc(cmdEnum.CS_EquipStrengthen, func(user *user.User, msg *codecs.Message) transaction.Transaction {
		return &transactionEquipStrengthen{
			user: user,
			req:  msg,
		}
	})
}
