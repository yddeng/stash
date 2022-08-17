package BigSecret

import (
	"github.com/golang/protobuf/proto"
	codecs "initialthree/codec/cs"
	"initialthree/node/common/transaction"
	"initialthree/node/table/excel/DataTable/BigSecretDungeon"
	"initialthree/zaplogger"
	"math/rand"

	"initialthree/node/node_game/user"
	"initialthree/protocol/cmdEnum"
	cs_message "initialthree/protocol/cs/message"
)

type transactionBigSecretRandomLevel struct {
	transaction.TransactionBase
	user    *user.User
	req     *codecs.Message
	errcode cs_message.ErrCode
	resp    *cs_message.BigSecretRandomLevelToC
}

func (this *transactionBigSecretRandomLevel) GetModuleName() string {
	return "Level"
}

func (this *transactionBigSecretRandomLevel) Begin() {
	defer func() { this.EndTrans(this.resp, this.errcode) }()
	this.resp = &cs_message.BigSecretRandomLevelToC{}
	this.errcode = cs_message.ErrCode_OK
	msg := this.req.GetData().(*cs_message.BigSecretRandomLevelToS)
	zaplogger.GetSugar().Infof("%s BigSecretRandomLevelToS %v ", this.user.GetUserLogName(), msg)

	def := BigSecretDungeon.GetID(msg.GetLevel())
	if def == nil || len(def.DungeonIDPoolArray) == 0 {
		zaplogger.GetSugar().Infof("%s BigSecretRandomLevelToS failed, config %d is nil", this.user.GetUserLogName(), msg.GetLevel())
		this.errcode = cs_message.ErrCode_Config_NotExist
		return
	}

	idx := rand.Int() % len(def.DungeonIDPoolArray)
	levelID := def.DungeonIDPoolArray[idx].DungeonID

	this.resp.LevelID = proto.Int32(levelID)
}

func init() {
	user.RegisterTransFunc(cmdEnum.CS_BigSecretRandomLevel, func(user *user.User, msg *codecs.Message) transaction.Transaction {
		return &transactionBigSecretRandomLevel{
			user: user,
			req:  msg,
		}
	})
}
