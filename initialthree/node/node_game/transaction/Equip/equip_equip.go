package Equip

import (
	codecs "initialthree/codec/cs"
	"initialthree/node/common/transaction"
	"initialthree/node/node_game/event"
	"initialthree/zaplogger"

	"initialthree/node/node_game/module"
	"initialthree/node/node_game/module/character"
	"initialthree/node/node_game/module/equip"
	"initialthree/node/node_game/user"
	"initialthree/node/table/excel/DataTable/Equip"
	"initialthree/protocol/cmdEnum"
	cs_message "initialthree/protocol/cs/message"
)

const totalWeight = 10

type transactionEquipEquip struct {
	transaction.TransactionBase
	user    *user.User
	req     *codecs.Message
	resp    *cs_message.EquipEquipToC
	errcode cs_message.ErrCode
}

func (this *transactionEquipEquip) GetModuleName() string {
	return "Equip"
}

func (this *transactionEquipEquip) Begin() {
	defer func() { this.EndTrans(this.resp, this.errcode) }()
	this.resp = &cs_message.EquipEquipToC{}
	msg := this.req.GetData().(*cs_message.EquipEquipToS)

	zaplogger.GetSugar().Infof("%s Call EquipEquipToS %v", this.user.GetUserLogName(), msg)

	userCharacter := this.user.GetSubModule(module.Character).(*character.UserCharacter)
	userEquip := this.user.GetSubModule(module.Equip).(*equip.UserEquip)

	msgEquipID := msg.GetEquipID()
	chara := userCharacter.GetCharacter(msg.GetCharacterID())
	doEquip := userEquip.GetEquip(msgEquipID)
	if chara == nil || doEquip == nil {
		zaplogger.GetSugar().Debugf("%s EquipEquipToS error: chara or mod is nil", this.user.GetUserLogName())
		this.errcode = cs_message.ErrCode_Equip_NotExist
		return
	}

	def := Equip.GetID(doEquip.ConfigID)
	if def == nil {
		zaplogger.GetSugar().Debugf("%s EquipEquipToS error: %d modDef is nil", this.user.GetUserLogName(), doEquip.ConfigID)
		this.errcode = cs_message.ErrCode_Config_NotExist
		return
	}

	modIds := chara.GetEquipIDs()
	if msg.GetPosition() >= int32(len(modIds)) {
		zaplogger.GetSugar().Debugf("%s EquipEquipToS error: pos %d is failed", this.user.GetUserLogName(), msg.GetPosition())
		this.errcode = cs_message.ErrCode_Request_Argument_Err
		return
	}

	wight := float64(0)
	for _, id := range modIds {
		if id != 0 && id != msgEquipID {
			// 待操作的装备已经在自己身上， 更换穿戴位置
			mm := userEquip.GetEquip(id)
			conf := Equip.GetID(mm.ConfigID)
			wight += conf.Cost
		}
	}

	equipId := modIds[msg.GetPosition()]
	if equipId == msgEquipID {
		this.errcode = cs_message.ErrCode_OK
		return
	}

	// 待穿戴的位置是否已经有
	if equipId == 0 {
		if wight+def.Cost > totalWeight {
			zaplogger.GetSugar().Debugf("%s EquipEquipToS error: wight is not enough", this.user.GetUserLogName())
			this.errcode = cs_message.ErrCode_Equip_Weight
			return
		}

		// 待穿戴的装备是否已经在其他角色身上
		if doEquip.EquipCharacterId != 0 {
			oldChara := userCharacter.GetCharacter(doEquip.EquipCharacterId)
			userCharacter.EquipDemount(oldChara, msgEquipID)
			userEquip.Demount(doEquip)
		}

		userCharacter.EquipEquip(chara, msg.GetPosition(), msgEquipID)
		userEquip.Equip(doEquip, msg.GetCharacterID())
	} else {
		oldEquip := userEquip.GetEquip(equipId)
		oldDef := Equip.GetID(oldEquip.ConfigID)

		// 判断重量
		if wight-oldDef.Cost+def.Cost > totalWeight {
			zaplogger.GetSugar().Debugf("%s EquipEquipToS error: weight is not enough", this.user.GetUserLogName())
			this.errcode = cs_message.ErrCode_Equip_Weight
			return
		}

		// 卸下当前位置mod
		userEquip.Demount(oldEquip)

		// 待穿戴的装备是否已经在其他角色身上
		if doEquip.EquipCharacterId != 0 {
			oldChara := userCharacter.GetCharacter(doEquip.EquipCharacterId)
			userCharacter.EquipDemount(oldChara, msgEquipID)
			userEquip.Demount(doEquip)
		}
		// 装备
		userEquip.Equip(doEquip, msg.GetCharacterID())
		userCharacter.EquipEquip(chara, msg.GetPosition(), msgEquipID)
	}

	this.user.EmitEvent(event.EventEquipEquipped)

	this.errcode = cs_message.ErrCode_OK
	zaplogger.GetSugar().Infof("%s Call EquipEquipToS ok", this.user.GetUserLogName())
}

func init() {
	user.RegisterTransFunc(cmdEnum.CS_EquipEquip, func(user *user.User, msg *codecs.Message) transaction.Transaction {
		return &transactionEquipEquip{
			user: user,
			req:  msg,
		}
	})
}
