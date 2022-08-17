package Equip

import (
	codecs "initialthree/codec/cs"
	"initialthree/node/common/attr"
	"initialthree/node/common/enumType"
	"initialthree/node/common/inoutput"
	"initialthree/node/common/transaction"
	"initialthree/node/node_game/module/equip"
	"initialthree/node/table/excel/DataTable/EquipLevelMaxExp"
	"initialthree/zaplogger"

	"initialthree/node/node_game/module"
	"initialthree/node/node_game/user"
	constEquip "initialthree/node/table/excel/ConstTable/Equip"
	"initialthree/node/table/excel/DataTable/Equip"
	"initialthree/node/table/excel/DataTable/EquipDecompose"
	"initialthree/protocol/cmdEnum"
	cs_message "initialthree/protocol/cs/message"
)

type transactionEquipDecompose struct {
	transaction.TransactionBase
	user    *user.User
	req     *codecs.Message
	resp    *cs_message.EquipDecomposeToC
	errcode cs_message.ErrCode
}

func (this *transactionEquipDecompose) GetModuleName() string {
	return "Equip"
}

func (this *transactionEquipDecompose) Begin() {
	defer func() { this.EndTrans(this.resp, this.errcode) }()
	this.resp = &cs_message.EquipDecomposeToC{}
	msg := this.req.GetData().(*cs_message.EquipDecomposeToS)

	zaplogger.GetSugar().Infof("%s Call EquipDecomposeToS %v", this.user.GetUserLogName(), msg)

	userEquip := this.user.GetSubModule(module.Equip).(*equip.UserEquip)

	useRes := make([]inoutput.ResDesc, 0, len(msg.GetEquipID()))
	addGold := float64(0)
	toItems := map[int32]float64{}
	for _, id := range msg.GetEquipID() {
		e := userEquip.GetEquip(id)
		if e == nil {
			zaplogger.GetSugar().Debugf("%s EquipDecomposeToS error: equip %d is nil", this.user.GetUserLogName(), id)
			this.errcode = cs_message.ErrCode_Equip_NotExist
			return
		}
		// 装备已被装备
		if e.IsLock || e.EquipCharacterId != 0 {
			zaplogger.GetSugar().Debugf("%s EquipDecomposeToS error: equip %d is equipped or locked", this.user.GetUserLogName(), id)
			this.errcode = cs_message.ErrCode_Equip_Equipped
			return
		}

		def := Equip.GetID(e.ConfigID)
		if def == nil {
			zaplogger.GetSugar().Debugf("%s EquipDecomposeToS error: equip %d def is nil", this.user.GetUserLogName(), id)
			this.errcode = cs_message.ErrCode_Equip_NotExist
			return
		}

		decDef := EquipDecompose.GetID(def.Decompose)
		if decDef == nil {
			zaplogger.GetSugar().Debugf("%s EquipDecomposeToS error:  EquipDecompose %d def is nil", this.user.GetUserLogName(), id)
			this.errcode = cs_message.ErrCode_Equip_NotExist
			return
		}

		levelDef := EquipLevelMaxExp.GetID(def.LevelMaxExp)
		if levelDef == nil {
			zaplogger.GetSugar().Debugf("%s EquipDecomposeToS error:  %d levelDef is nil", this.user.GetUserLogName(), def.LevelMaxExp)
			this.errcode = cs_message.ErrCode_Config_NotExist
			return
		}
		toItemId := decDef.ExpToItemID
		supplyExp := constEquip.GetSupplyExp(toItemId)
		if supplyExp == 0 {
			zaplogger.GetSugar().Debugf("%s EquipDecomposeToS error:  %d supplyExp is nil", this.user.GetUserLogName(), toItemId)
			this.errcode = cs_message.ErrCode_Config_NotExist
			return
		}

		// 总经验
		exp := levelDef.GetRangeTotalExp(1, e.Level) + e.Exp
		// 转金币
		addGold += float64(exp)/decDef.ExpToOneGoldRate + float64(decDef.ReturnGoldCount)
		// 转道具数量
		val := float64(exp) / float64(supplyExp)
		toItems[toItemId] += val
		for _, v := range decDef.Items {
			toItems[v.Id] += float64(v.Count)
		}

		useRes = append(useRes, inoutput.ResDesc{Type: enumType.IOType_Equip, ID: int32(id), Count: 1})
	}

	// 消耗添加
	addRes := make([]inoutput.ResDesc, 0, len(toItems)+1)
	addRes = append(addRes, inoutput.ResDesc{Type: enumType.IOType_UsualAttribute, ID: attr.Gold, Count: int32(addGold)})
	for itemID, itemCount := range toItems {
		addRes = append(addRes, inoutput.ResDesc{Type: enumType.IOType_Item, ID: itemID, Count: int32(itemCount)})
	}

	if this.errcode = this.user.DoInputOutput(useRes, addRes, true); this.errcode != cs_message.ErrCode_OK {
		zaplogger.GetSugar().Debugf("%s EquipDecomposeToS error: inouput %s  ", this.user.GetUserLogName(), this.errcode.String())
		return
	}

	this.errcode = cs_message.ErrCode_OK
	zaplogger.GetSugar().Infof("%s Call EquipDecomposeToS ok", this.user.GetUserLogName())
}

func init() {
	user.RegisterTransFunc(cmdEnum.CS_EquipDecompose, func(user *user.User, msg *codecs.Message) transaction.Transaction {
		return &transactionEquipDecompose{
			user: user,
			req:  msg,
		}
	})
}
