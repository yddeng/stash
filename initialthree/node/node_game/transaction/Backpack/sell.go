package Backpack

import (
	"initialthree/node/common/attr"
	"initialthree/node/common/enumType"
	"initialthree/node/common/transaction"
	"initialthree/node/table/excel/DataTable/Item"
	"initialthree/zaplogger"

	"initialthree/node/node_game/module"
	attr2 "initialthree/node/node_game/module/attr"
	backpack2 "initialthree/node/node_game/module/backpack"
	"initialthree/node/node_game/user"
	"initialthree/protocol/cmdEnum"
	"time"

	codecs "initialthree/codec/cs"
	cs_message "initialthree/protocol/cs/message"
)

type transactionBackpackSell struct {
	transaction.TransactionBase
	user    *user.User
	req     *codecs.Message
	errcode cs_message.ErrCode
}

func (t *transactionBackpackSell) Begin() {
	defer func() { t.EndTrans(&cs_message.BackpackSellToC{}, t.errcode) }()

	t.errcode = cs_message.ErrCode_OK
	reqMsg := t.req.GetData().(*cs_message.BackpackSellToS)
	backpack := t.user.GetSubModule(module.Backpack).(*backpack2.Backpack)
	userAttr := t.user.GetSubModule(module.Attr).(*attr2.UserAttr)

	var sellEntities map[uint32]int32
	var sellPrices map[int32]int64

	if len(reqMsg.GetSellEntities()) == 0 {
		t.errcode = cs_message.ErrCode_Request_Argument_Err
		return
	}

	sellEntities = make(map[uint32]int32, len(reqMsg.GetSellEntities()))
	for _, v := range reqMsg.GetSellEntities() {
		if v.GetCount() <= 0 {
			zaplogger.GetSugar().Infof("user(%s, %d) backpack sell entity(%d, %d): sell count error",
				t.user.GetUserID(), t.user.GetID(), v.GetId(), v.GetCount())
			t.errcode = cs_message.ErrCode_Backpack_EntityError
			return
		}

		sellEntities[v.GetId()] += v.GetCount()
	}

	sellPrices = make(map[int32]int64)
	for id, count := range sellEntities {
		it := backpack.GetItem(id)
		if it == nil {
			zaplogger.GetSugar().Infof("user(%s, %d) backpack sell entity(%d, %d): entity not found",
				t.user.GetUserID(), t.user.GetID(), id, count)
			t.errcode = cs_message.ErrCode_Backpack_EntityError
			return
		}

		if it.IsExpired(time.Now()) {
			zaplogger.GetSugar().Infof("user(%s, %d) backpack sell entity(%d, %d): entity expired",
				t.user.GetUserID(), t.user.GetID(), id, count)
			t.errcode = cs_message.ErrCode_Backpack_EntityExpired
			return
		}

		if count > it.Count {
			zaplogger.GetSugar().Infof("user(%s, %d) backpack sell entity(%d, %d): count not enough",
				t.user.GetUserID(), t.user.GetID(), id, count)
			t.errcode = cs_message.ErrCode_Backpack_EntityError
			return
		}

		itCfg := Item.GetID(it.TID)

		if !itCfg.AllowSell {
			zaplogger.GetSugar().Infof("user(%s, %d) backpack sell entity(%d, %d): not allowed to sell",
				t.user.GetUserID(), t.user.GetID(), id, count)
			t.errcode = cs_message.ErrCode_Backpack_NotAllowedToSell
			return
		}

		if attr.GetAttrInfo(itCfg.SellCurrencyTypeEnum) == nil {
			zaplogger.GetSugar().Infof("user(%s, %d) backpack sell entity(%d, %d): currency type error",
				t.user.GetUserID(), t.user.GetID(), id, it.TID)
			t.errcode = cs_message.ErrCode_Config_NotExist
			return
		}

		if itCfg.SellPrice <= 0 {
			zaplogger.GetSugar().Infof("user(%s, %d) backpack sell entity(%d, %d): unit price error",
				t.user.GetUserID(), t.user.GetID(), id, it.TID)
			t.errcode = cs_message.ErrCode_Config_Error
			return
		}

		switch itCfg.SellCurrencyTypeEnum {
		case enumType.CurrencyType_Gold:
			sellPrices[attr.Gold] += int64(itCfg.SellPrice * count)
		case enumType.CurrencyType_Diamond:
			sellPrices[attr.Diamond] += int64(itCfg.SellPrice * count)
		}
	}

	for id, count := range sellEntities {
		backpack.AddItemCount(id, -count)
	}

	for id, value := range sellPrices {
		_, _ = userAttr.AddAttr(id, value)
	}

}

func (t *transactionBackpackSell) GetModuleName() string {
	return "Backpack"
}

func init() {
	user.RegisterTransFunc(cmdEnum.CS_BackpackSell, func(user *user.User, msg *codecs.Message) transaction.Transaction {
		return &transactionBackpackSell{
			user: user,
			req:  msg,
		}
	})
}
