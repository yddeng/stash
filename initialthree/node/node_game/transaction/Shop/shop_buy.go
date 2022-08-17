package Shop

import (
	"initialthree/codec/cs"
	"initialthree/node/common/attr"
	"initialthree/node/common/enumType"
	"initialthree/node/common/inoutput"
	"initialthree/node/common/transaction"
	"initialthree/node/node_game/event"
	"initialthree/node/node_game/module"
	"initialthree/node/node_game/module/shop"
	"initialthree/node/node_game/user"
	"initialthree/node/table/excel/DataTable/Product"
	"initialthree/protocol/cmdEnum"
	cs_msg "initialthree/protocol/cs/message"
	"initialthree/zaplogger"
	"time"
)

type transactionShopBuy struct {
	transaction.TransactionBase
	user    *user.User
	req     *cs.Message
	errCode cs_msg.ErrCode
}

func (t *transactionShopBuy) Begin() {
	defer func() { t.EndTrans(&cs_msg.ShopBuyToC{}, t.errCode) }()
	t.errCode = cs_msg.ErrCode_OK

	reqMsg := t.req.GetData().(*cs_msg.ShopBuyToS)
	zaplogger.GetSugar().Infof("%s Call ShopBuy %v", t.user.GetUserLogName(), reqMsg)

	def := Product.GetID(reqMsg.GetId())
	if def == nil {
		zaplogger.GetSugar().Infof("%s ShopBuy fail:Product def is nil", t.user.GetUserLogName())
		t.errCode = cs_msg.ErrCode_Config_NotExist
		return
	}

	userShop := t.user.GetSubModule(module.Shop).(*shop.ShopData)
	if def.LimitCount > 0 {
		times := userShop.GetProductBuyTimes(reqMsg.GetId())
		if times >= def.LimitCount {
			zaplogger.GetSugar().Infof("%s ShopBuy fail:Product times %d >= limitCount %d", t.user.GetUserLogName(), times, def.LimitCount)
			t.errCode = cs_msg.ErrCode_Shop_BuyTimesNotEnough
			return
		}
	}

	// 日期限购
	if def.ProductLimitTypeEnum == enumType.ProductLimitType_Date {
		nowUnix := time.Now().Unix()
		limitDate := def.LimitDate()
		if nowUnix < limitDate.StartTime || nowUnix > limitDate.EndTime {
			zaplogger.GetSugar().Infof("%s ShopBuy fail: limit date", t.user.GetUserLogName())
			t.errCode = cs_msg.ErrCode_Shop_BuyLimit
			return
		}
	}

	attrID := attr.GetIdByName(def.PriceType)
	if attrID == 0 {
		zaplogger.GetSugar().Infof("%s ShopBuy fail:PriceType %s is error", t.user.GetUserLogName(), def.PriceType)
		t.errCode = cs_msg.ErrCode_Config_Error
		return
	}

	price := def.Price
	if def.Discount != 0 {
		price = int32(float64(price) * def.Discount)
	}
	input := []inoutput.ResDesc{{ID: attrID, Type: enumType.IOType_UsualAttribute, Count: price * reqMsg.GetCount()}}
	output := []inoutput.ResDesc{{ID: def.PID, Type: int(def.PTypeEnum), Count: def.PCount * reqMsg.GetCount()}}

	if t.errCode = t.user.DoInputOutput(input, output, false); t.errCode != cs_msg.ErrCode_OK {
		zaplogger.GetSugar().Debugf("%s ShopBuy error: inouput %s  ", t.user.GetUserLogName(), t.errCode.String())
		return
	}

	userShop.ShopBuy(reqMsg.GetId(), reqMsg.GetCount())
	zaplogger.GetSugar().Infof("%v ShopBuy OK", t.user.GetUserLogName())
	t.user.EmitEvent(event.EventShopBuy, reqMsg.GetId(), reqMsg.GetCount())
}

func (t *transactionShopBuy) GetModuleName() string {
	return "Shop"
}

func init() {
	user.RegisterTransFunc(cmdEnum.CS_ShopBuy, func(user *user.User, msg *cs.Message) transaction.Transaction {
		return &transactionShopBuy{
			user: user,
			req:  msg,
		}
	})
}
