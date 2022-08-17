package Shop

import (
	"initialthree/codec/cs"
	"initialthree/node/common/attr"
	"initialthree/node/common/enumType"
	"initialthree/node/common/inoutput"
	"initialthree/node/common/transaction"
	"initialthree/node/node_game/module"
	"initialthree/node/node_game/module/assets"
	"initialthree/node/node_game/user"
	"initialthree/node/table/excel/DataTable/Pay"
	"initialthree/protocol/cmdEnum"
	cs_msg "initialthree/protocol/cs/message"
	"initialthree/zaplogger"
)

type transactionShopPay struct {
	transaction.TransactionBase
	user    *user.User
	req     *cs.Message
	errCode cs_msg.ErrCode
}

func (t *transactionShopPay) Begin() {
	defer func() { t.EndTrans(&cs_msg.ShopPayToC{}, t.errCode) }()
	t.errCode = cs_msg.ErrCode_OK

	reqMsg := t.req.GetData().(*cs_msg.ShopPayToS)
	zaplogger.GetSugar().Infof("%s Call ShopPay %v", t.user.GetUserLogName(), reqMsg)

	def := Pay.GetID(reqMsg.GetPayID())
	if def == nil {
		zaplogger.GetSugar().Infof("%s ShopPay fail:Pay def is nil", t.user.GetUserLogName())
		t.errCode = cs_msg.ErrCode_Config_NotExist
		return
	}

	resID := attr.GetIdByName(def.Type)
	if resID == 0 {
		zaplogger.GetSugar().Infof("%s ShopPay fail:Pay Type %s is not found ", t.user.GetUserLogName(), def.Type)
		t.errCode = cs_msg.ErrCode_Config_NotExist
		return
	}

	count := int32(0)
	needSetAsset := false
	userAsset := t.user.GetSubModule(module.Assets).(*assets.UserAssets)
	if _, ok := userAsset.GetAssetCount(int32(cs_msg.AssetType_ShopFirstPay), reqMsg.GetPayID()); !ok {
		count = def.FirstPayCount + def.FirstPresentedCount
		needSetAsset = true
	} else {
		count = def.CommonPayCount + def.CommonPresentedCount
	}

	output := []inoutput.ResDesc{{ID: resID, Type: enumType.IOType_UsualAttribute, Count: count}}

	if t.errCode = t.user.DoInputOutput(nil, output, true); t.errCode != cs_msg.ErrCode_OK {
		zaplogger.GetSugar().Debugf("%s ShopPay error: inouput %s  ", t.user.GetUserLogName(), t.errCode.String())
		return
	}

	if needSetAsset {
		userAsset.SetAsset(int32(cs_msg.AssetType_ShopFirstPay), reqMsg.GetPayID(), 1)
	}
}

func (t *transactionShopPay) GetModuleName() string {
	return "Shop"
}

func init() {
	user.RegisterTransFunc(cmdEnum.CS_ShopPay, func(user *user.User, msg *cs.Message) transaction.Transaction {
		return &transactionShopPay{
			user: user,
			req:  msg,
		}
	})
}
