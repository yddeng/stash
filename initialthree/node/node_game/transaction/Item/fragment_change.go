package Item

import (
	"initialthree/codec/cs"
	"initialthree/node/common/enumType"
	"initialthree/node/common/inoutput"
	"initialthree/node/common/transaction"
	"initialthree/node/node_game/user"
	"initialthree/node/table/excel/DataTable/FragmentChange"
	"initialthree/protocol/cmdEnum"
	cs_msg "initialthree/protocol/cs/message"
	"initialthree/zaplogger"
)

type transactionFragmentChange struct {
	transaction.TransactionBase
	user    *user.User
	req     *cs.Message
	errCode cs_msg.ErrCode
}

func (t *transactionFragmentChange) Begin() {
	defer func() { t.EndTrans(&cs_msg.FragmentChangeToC{}, t.errCode) }()
	t.errCode = cs_msg.ErrCode_OK

	reqMsg := t.req.GetData().(*cs_msg.FragmentChangeToS)

	reqCount := reqMsg.GetCount()
	if reqCount <= 0 {
		zaplogger.GetSugar().Errorf("user(%s, %d) fragment change: count error.",
			t.user.GetUserID(), t.user.GetID(), reqMsg.GetId())
		t.errCode = cs_msg.ErrCode_ERROR
		return
	}
	/*
		userCharacter := t.user.GetSubModule(module.Character).(*character.UserCharacter)

		characterCfg := PlayerCharacter.GetFragmentIDint32(reqMsg.GetId())
		if characterCfg == nil {
			zaplogger.GetSugar().Errorf("user(%s, %d) fragment change: fragment %d rel character config not found.",
				t.user.GetUserID(), t.user.GetID(), reqMsg.GetId())
			t.errCode = cs_msg.ErrCode_Config_NotExist
			return
		}

		character := userCharacter.GetCharacter(characterCfg.ID)
		if character == nil {
			zaplogger.GetSugar().Errorf("user(%s, %d) fragment change: fragment %d cant change.",
				t.user.GetUserID(), t.user.GetID(), reqMsg.GetId())
			t.errCode = cs_msg.ErrCode_ERROR
			return
		}
	*/
	changeCfg := FragmentChange.GetID(reqMsg.GetId())
	if changeCfg == nil {
		zaplogger.GetSugar().Errorf("user(%s, %d) fragment change: fragment %d change config not found.",
			t.user.GetUserID(), t.user.GetID(), reqMsg.GetId())
		t.errCode = cs_msg.ErrCode_Config_NotExist
		return
	}

	input := []inoutput.ResDesc{{ID: reqMsg.GetId(), Type: enumType.IOType_Item, Count: reqCount}}
	output := []inoutput.ResDesc{}

	for _, v := range changeCfg.Item {
		output = append(output, inoutput.ResDesc{ID: v.ID, Type: enumType.IOType_Item, Count: v.Count * reqCount})
	}

	t.errCode = t.user.DoInputOutput(input, output, true)
}

func (t *transactionFragmentChange) GetModuleName() string {
	return "Fragment"
}

func init() {
	user.RegisterTransFunc(cmdEnum.CS_FragmentChange, func(user *user.User, msg *cs.Message) transaction.Transaction {
		return &transactionFragmentChange{
			user: user,
			req:  msg,
		}
	})
}
