package Shop

import (
	"initialthree/codec/cs"
	"initialthree/node/common/attr"
	"initialthree/node/common/enumType"
	"initialthree/node/common/inoutput"
	"initialthree/node/common/transaction"
	"initialthree/node/node_game/module"
	"initialthree/node/node_game/module/shop"
	"initialthree/node/node_game/user"
	"initialthree/node/table/excel/DataTable/ProductLibrary"
	"initialthree/protocol/cmdEnum"
	cs_msg "initialthree/protocol/cs/message"
	"initialthree/zaplogger"
)

type transactionShopRefresh struct {
	transaction.TransactionBase
	user    *user.User
	req     *cs.Message
	errCode cs_msg.ErrCode
}

func (t *transactionShopRefresh) Begin() {
	defer func() { t.EndTrans(&cs_msg.ShopRefreshToC{}, t.errCode) }()
	t.errCode = cs_msg.ErrCode_OK

	reqMsg := t.req.GetData().(*cs_msg.ShopRefreshToS)
	zaplogger.GetSugar().Infof("%s Call ShopRefresh %v", t.user.GetUserLogName(), reqMsg)

	def := ProductLibrary.GetID(reqMsg.GetShopID())
	if def == nil {
		zaplogger.GetSugar().Infof("%s ShopRefresh fail:ProductLibrary def is nil", t.user.GetUserLogName())
		t.errCode = cs_msg.ErrCode_Config_NotExist
		return
	}

	userShop := t.user.GetSubModule(module.Shop).(*shop.ShopData)
	times := userShop.GetShopRefreshTimes(reqMsg.GetShopID())
	if int(times) >= len(def.RedreshPrice) {
		zaplogger.GetSugar().Infof("%s ShopRefresh fail:RefreshTimes is not enough", t.user.GetUserLogName())
		t.errCode = cs_msg.ErrCode_Shop_RefreshTimesNotEnough
		return
	}

	attrID := attr.GetIdByName(def.RefreshPriceType)
	if attrID == 0 {
		zaplogger.GetSugar().Infof("%s ShopRefresh fail:RefreshPriceType %s is error", t.user.GetUserLogName(), def.RefreshPriceType)
		t.errCode = cs_msg.ErrCode_Config_Error
		return
	}

	cost := def.RedreshPrice[times].Price
	input := []inoutput.ResDesc{{ID: attrID, Type: enumType.IOType_UsualAttribute, Count: cost}}

	if t.errCode = t.user.DoInputOutput(input, nil, false); t.errCode != cs_msg.ErrCode_OK {
		zaplogger.GetSugar().Debugf("%s ShopRefresh error: inouput %s  ", t.user.GetUserLogName(), t.errCode.String())
		return
	}

	products := make([]int32, 0, len(def.ProductsArray))
	for _, v := range def.ProductsArray {
		products = append(products, v.ID)
	}

	userShop.ShopRefresh(reqMsg.GetShopID(), products)
	zaplogger.GetSugar().Infof("%v ShopRefresh OK", t.user.GetUserLogName())
}

func (t *transactionShopRefresh) GetModuleName() string {
	return "Shop"
}

func init() {
	user.RegisterTransFunc(cmdEnum.CS_ShopRefresh, func(user *user.User, msg *cs.Message) transaction.Transaction {
		return &transactionShopRefresh{
			user: user,
			req:  msg,
		}
	})
}
