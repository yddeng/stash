package Backpack

import (
	"initialthree/node/common/transaction"
	"initialthree/zaplogger"

	"initialthree/node/node_game/module"
	backpack2 "initialthree/node/node_game/module/backpack"
	"initialthree/node/node_game/user"
	"initialthree/protocol/cmdEnum"
	"time"

	codecs "initialthree/codec/cs"
	cs_message "initialthree/protocol/cs/message"
)

type transactionBackpackRem struct {
	transaction.TransactionBase
	user    *user.User
	req     *codecs.Message
	errcode cs_message.ErrCode
}

func (t *transactionBackpackRem) Begin() {
	defer func() { t.EndTrans(&cs_message.BackpackRemToC{}, t.errcode) }()
	t.errcode = cs_message.ErrCode_OK

	reqMsg := t.req.GetData().(*cs_message.BackpackRemToS)
	backpack := t.user.GetSubModule(module.Backpack).(*backpack2.Backpack)

	for _, v := range reqMsg.GetRemEntities() {
		item := backpack.GetItem(v.GetId())

		if item == nil {
			zaplogger.GetSugar().Errorf("user(%s, %d) backpack rem entity %d: entity not found",
				t.user.GetUserID(), t.user.GetID(), v.GetId())
			t.errcode = cs_message.ErrCode_Backpack_EntityError
			return
		}

		switch v.GetType() {
		case 1: // 过期删除
			if !item.IsExpired(time.Now()) {
				zaplogger.GetSugar().Errorf("user(%s, %d) backpack rem entity %d: entity not expired",
					t.user.GetUserID(), t.user.GetID(), v.GetId())
				t.errcode = cs_message.ErrCode_Backpack_EntityError
				return
			}

		case 2: // 常规删除
			if item.Count < v.GetCount() {
				zaplogger.GetSugar().Errorf("user(%s, %d) backpack rem entity(%d, %d): entity count error",
					t.user.GetUserID(), t.user.GetID(), v.GetId(), v.GetCount())
				t.errcode = cs_message.ErrCode_Backpack_EntityError
				return
			}

		default:
			zaplogger.GetSugar().Infof("user(%s, %d) backpack rem entity %d: wrong type %d",
				t.user.GetUserID(), t.user.GetID(), v.GetId(), v.GetType())
			t.errcode = cs_message.ErrCode_Backpack_RemTypeError
			return
		}
	}

	for _, v := range reqMsg.GetRemEntities() {
		switch v.GetType() {
		case 1: // 过期删除
			backpack.RemItem(v.GetId())
		case 2:
			backpack.AddItemCount(v.GetId(), -v.GetCount())
		}
	}

}

func (t *transactionBackpackRem) GetModuleName() string {
	return "Backpack"
}

func init() {
	user.RegisterTransFunc(cmdEnum.CS_BackpackRem, func(user *user.User, msg *codecs.Message) transaction.Transaction {
		return &transactionBackpackRem{
			user: user,
			req:  msg,
		}
	})
}
