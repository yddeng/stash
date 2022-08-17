package SecertDungeon

import (
	"github.com/golang/protobuf/proto"
	codecs "initialthree/codec/cs"
	"initialthree/node/common/transaction"
	"initialthree/node/table/excel/DataTable/SecretDungeonPool"
	"initialthree/zaplogger"
	"math/rand"

	"initialthree/node/node_game/user"
	"initialthree/protocol/cmdEnum"
	cs_message "initialthree/protocol/cs/message"
)

type transactionSecretRandomLevel struct {
	transaction.TransactionBase
	user    *user.User
	req     *codecs.Message
	errcode cs_message.ErrCode
	resp    *cs_message.SecretRandomLevelToC
}

func (this *transactionSecretRandomLevel) GetModuleName() string {
	return "Level"
}

func (this *transactionSecretRandomLevel) Begin() {
	defer func() { this.EndTrans(this.resp, this.errcode) }()
	this.resp = &cs_message.SecretRandomLevelToC{}
	this.errcode = cs_message.ErrCode_OK
	msg := this.req.GetData().(*cs_message.SecretRandomLevelToS)
	zaplogger.GetSugar().Infof("%s SecretRandomLevelToS %v ", this.user.GetUserLogName(), msg)

	def := SecretDungeonPool.GetID(msg.GetDifficult())
	if def == nil || len(def.PoolArray) == 0 {
		zaplogger.GetSugar().Infof("%s SecretRandomLevelToS failed, config %d is nil", this.user.GetUserLogName(), msg.GetDifficult())
		this.errcode = cs_message.ErrCode_Config_NotExist
		return
	}

	idx := rand.Int() % len(def.PoolArray)
	levelID := def.PoolArray[idx].ID

	this.resp.LevelID = proto.Int32(levelID)
}

func init() {
	user.RegisterTransFunc(cmdEnum.CS_SecretRandomLevel, func(user *user.User, msg *codecs.Message) transaction.Transaction {
		return &transactionSecretRandomLevel{
			user: user,
			req:  msg,
		}
	})
}
