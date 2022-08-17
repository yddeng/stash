package Resurrect

import (
	"github.com/gogo/protobuf/proto"
	"initialthree/node/common/inoutput"
	"initialthree/zaplogger"

	codecs "initialthree/codec/cs"
	"initialthree/node/common/enumType"
	"initialthree/node/common/transaction"

	"initialthree/node/node_game/temporary"
	"initialthree/node/node_game/user"
	"initialthree/node/table/excel/DataTable/Dungeon"
	"initialthree/node/table/excel/DataTable/Resurrect"
	"initialthree/protocol/cmdEnum"
	cs_message "initialthree/protocol/cs/message"
)

type transactionResurrect struct {
	transaction.TransactionBase
	user    *user.User
	req     *codecs.Message
	resp    *cs_message.BattleResurrectToC
	errCode cs_message.ErrCode
}

func (t *transactionResurrect) GetModuleName() string {
	return "Resurrect"
}

func (t *transactionResurrect) Begin() {
	defer func() { t.EndTrans(t.resp, t.errCode) }()
	t.resp = &cs_message.BattleResurrectToC{}
	t.errCode = cs_message.ErrCode_BattleResurrect_CanNotResurrect

	zaplogger.GetSugar().Infof("%s %s", t.user.GetUserID(), "Resurrect ==>")
	if temp := t.user.GetTemporary(temporary.TempLevelFight); temp != nil {
		if v, ok := temp.(*temporary.LevelFightInfo); ok {
			dungeonCfg := Dungeon.GetID(v.Tos.GetDungeonID())
			// 是否可以复活
			if dungeonCfg.AllowResurrect {
				if resurrect := Resurrect.GetID(dungeonCfg.ResurrectID); resurrect != nil {
					if v.ResurrectCount < len(resurrect.Resource) {
						if t.willResurrect(resurrect, v.ResurrectCount) {
							t.errCode = cs_message.ErrCode_OK
							t.resp.EffectID = proto.Int32(resurrect.Effect)
							v.ResurrectCount++
						}
					} else {
						t.errCode = cs_message.ErrCode_BattleResurrect_NoResurrectCount
					}
				}
			}
		}
	}

}

func (t *transactionResurrect) willResurrect(resurrect *Resurrect.Resurrect, resurrectCount int) bool {
	item := resurrect.Resource[resurrectCount]
	resType, _ := enumType.GetEnumType(item.ResourceType)
	if resType == enumType.ResourceConsumeType_Item && t.user.GetItemCountByTID(item.ResourceID) >= item.ResourceCount {
		itemDesc := inoutput.ResDesc{Type: enumType.IOType_Item, ID: item.ResourceID, Count: item.ResourceCount}
		return t.user.DoInputOutput([]inoutput.ResDesc{itemDesc}, nil, false) == cs_message.ErrCode_OK
	} else if resType == enumType.ResourceConsumeType_UsualAttribute && t.user.GetAttr(item.ResourceID) >= int64(item.ResourceCount) {
		t.user.DoInputOutput([]inoutput.ResDesc{{Type: enumType.IOType_UsualAttribute, ID: item.ResourceID, Count: item.ResourceCount}}, nil, false)
		return true
	}
	return false
}

func init() {
	user.RegisterTransFunc(cmdEnum.CS_BattleResurrect, func(user *user.User, msg *codecs.Message) transaction.Transaction {
		return &transactionResurrect{
			user: user,
			req:  msg,
		}
	})
}
