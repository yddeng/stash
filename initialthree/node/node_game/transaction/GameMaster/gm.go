package GameMaster

import (
	codecs "initialthree/codec/cs"
	attr2 "initialthree/node/common/attr"
	"initialthree/node/common/enumType"
	"initialthree/node/common/inoutput"
	"initialthree/node/common/transaction"
	"initialthree/node/node_game/event"
	"initialthree/node/node_game/module/attr"
	"initialthree/node/table/excel/DataTable/MainDungeon"
	"initialthree/zaplogger"

	"initialthree/node/node_game/module"
	"initialthree/node/node_game/module/maindungeons"
	"initialthree/node/node_game/user"
	"initialthree/protocol/cmdEnum"
	cs_message "initialthree/protocol/cs/message"
)

type transactionGameMaster struct {
	transaction.TransactionBase
	user    *user.User
	req     *codecs.Message
	resp    *cs_message.GameMasterToC
	errcode cs_message.ErrCode
}

func (this *transactionGameMaster) GetModuleName() string {
	return "GameMaster"
}

func (this *transactionGameMaster) Begin() {
	defer func() { this.EndTrans(this.resp, this.errcode) }()
	this.resp = &cs_message.GameMasterToC{}
	gmCmds := this.req.GetData().(*cs_message.GameMasterToS)

	zaplogger.GetSugar().Infof("%s %d %s %v", this.user.GetUserID(), this.user.GetID(), "Call GameMaster", gmCmds)

	//userAttr := this.user.GetSubModule(module.Attr).(*attr.UserAttr)

	this.errcode = cs_message.ErrCode_OK

	var res []inoutput.ResDesc

	for _, msg := range gmCmds.GetCmds() {
		tt := uint16(msg.GetType())
		switch tt {

		case 1: // 属性
			//res = append(res, inoutput.ResDesc{Type: enumType.IOType_UsualAttribute, ID: msg.GetID(), Count: msg.GetCount()})
			attrMoudle := this.user.GetSubModule(module.Attr).(*attr.UserAttr)
			if attr2.NewbieGiftStartTime == msg.GetID() || attr2.NewbieGiftEndTime == msg.GetID() {
				attrMoudle.SetAttr(msg.GetID(), int64(msg.GetCount()), false)
			} else {
				attrMoudle.AddAttr(msg.GetID(), int64(msg.GetCount()))
			}
		case 2: // 角色
			res = append(res, inoutput.ResDesc{Type: enumType.IOType_Character, ID: msg.GetID(), Count: msg.GetCount()})
		case 3: // 物品
			res = append(res, inoutput.ResDesc{Type: enumType.IOType_Item, ID: msg.GetID(), Count: msg.GetCount()})
		case 4: // 主线关卡通关
			userMainDungeons := this.user.GetSubModule(module.MainDungeons).(*maindungeons.MainDungeons)
			dungeonCfg := MainDungeon.GetID(msg.GetID())
			if dungeonCfg != nil {
				userMainDungeons.DungeonPass(msg.GetID())
				this.user.EmitEvent(event.EventInstanceSucceed, msg.GetID(), int32(0), map[int32]int32{})
			}
		case 5: // add 装备
			res = append(res, inoutput.ResDesc{Type: enumType.IOType_Equip, ID: msg.GetID(), Count: msg.GetCount()})
		case 6: // add 武器
			res = append(res, inoutput.ResDesc{Type: enumType.IOType_Weapon, ID: msg.GetID(), Count: msg.GetCount()})
		default:
			zaplogger.GetSugar().Infof("GameMaster err type %d", tt)
			this.errcode = cs_message.ErrCode_ERROR
		}
	}

	if len(res) > 0 {
		if this.errcode = this.user.DoInputOutput(nil, res, true); this.errcode != cs_message.ErrCode_OK {
			zaplogger.GetSugar().Debugf("%s GameMaster error: inouput %s  ", this.user.GetUserLogName(), this.errcode.String())
			return
		}
	}

}

func init() {
	user.RegisterTransFunc(cmdEnum.CS_GameMaster, func(user *user.User, msg *codecs.Message) transaction.Transaction {
		return &transactionGameMaster{
			user: user,
			req:  msg,
		}
	})
}
