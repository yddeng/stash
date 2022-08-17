package Shop

/*
import (
	"initialthree/codec/cs"
	"initialthree/node/common/transaction"
	"initialthree/node/node_game/module"
	"initialthree/node/node_game/module/shop"
	"initialthree/node/node_game/user"
	"initialthree/protocol/cmdEnum"
	cs_msg "initialthree/protocol/cs/message"
)

type transactionShopSync struct {
	transaction.TransactionBase
	user    *user.User
	req     *cs.Message
	errCode cs_msg.ErrCode
}

func (t *transactionShopSync) Begin() {
	userShop := t.user.GetSubModule(module.Shop).(*shop.ShopData)
	userShop.FlushAllToClient(t.req.GetSeriNo())
	t.EndTrans(userShop.ShopSyncToC(), cs_msg.ErrCode_OK)
}

func (t *transactionShopSync) GetModuleName() string {
	return "Shop"
}

func init() {
	user.RegisterTransFunc(cmdEnum.CS_ShopSync, func(user *user.User, msg *cs.Message) transaction.Transaction {
		return &transactionShopSync{
			user: user,
			req:  msg,
		}
	})
}
*/
