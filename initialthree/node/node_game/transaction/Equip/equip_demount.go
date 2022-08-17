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
	"initialthree/protocol/cmdEnum"
	cs_message "initialthree/protocol/cs/message"
)

type transactionEquipDemount struct {
	transaction.TransactionBase
	user    *user.User
	req     *codecs.Message
	resp    *cs_message.EquipDemountToC
	errcode cs_message.ErrCode
}

func (this *transactionEquipDemount) GetModuleName() string {
	return "Equip"
}

func (this *transactionEquipDemount) Begin() {
	defer func() { this.EndTrans(this.resp, this.errcode) }()
	this.resp = &cs_message.EquipDemountToC{}
	msg := this.req.GetData().(*cs_message.EquipDemountToS)

	zaplogger.GetSugar().Infof("%s Call EquipDemountToS %v", this.user.GetUserLogName(), msg)

	userCharacter := this.user.GetSubModule(module.Character).(*character.UserCharacter)
	userEquip := this.user.GetSubModule(module.Equip).(*equip.UserEquip)

	for _, id := range msg.GetEquipID() {
		m := userEquip.GetEquip(id)
		if m == nil {
			zaplogger.GetSugar().Debugf("%s EquipDemountToS error: mod %d is not exist", this.user.GetUserLogName(), id)
			this.errcode = cs_message.ErrCode_Equip_NotExist
			return
		}
	}

	doChara := map[int32]*character.Character{}
	for _, id := range msg.GetEquipID() {
		m := userEquip.GetEquip(id)
		c := userCharacter.GetCharacter(m.EquipCharacterId)
		if c != nil {
			userCharacter.EquipDemount(c, id)
			if _, ok := doChara[m.EquipCharacterId]; !ok {
				doChara[m.EquipCharacterId] = c
			}
		}
		userEquip.Demount(m)
		this.user.EmitEvent(event.EventEquipEquipped)
	}

	this.errcode = cs_message.ErrCode_OK
	zaplogger.GetSugar().Infof("%s Call EquipDemountToS ok", this.user.GetUserLogName())
}

func init() {
	user.RegisterTransFunc(cmdEnum.CS_EquipDemount, func(user *user.User, msg *codecs.Message) transaction.Transaction {
		return &transactionEquipDemount{
			user: user,
			req:  msg,
		}
	})
}
