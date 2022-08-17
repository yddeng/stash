package Map

import (
	"github.com/golang/protobuf/proto"
	"initialthree/cluster"
	"initialthree/cluster/addr"
	codecs "initialthree/codec/cs"
	"initialthree/node/common/transaction"
	"initialthree/zaplogger"

	"initialthree/node/node_game/temporary"
	"initialthree/node/node_game/user"
	"initialthree/protocol/cmdEnum"
	cs_message "initialthree/protocol/cs/message"
	ss_message "initialthree/protocol/ss/ssmessage"
)

type transactionStartAoi struct {
	transaction.TransactionBase
	user     *user.User
	req      *codecs.Message
	resp     *cs_message.StartAoiToC
	errcode  cs_message.ErrCode
	node_map addr.LogicAddr
}

func (this *transactionStartAoi) GetModuleName() string {
	return "Map"
}

func (this *transactionStartAoi) Begin() {
	this.resp = &cs_message.StartAoiToC{}

	mapInfoCache := this.user.GetTemporary(temporary.TempMapInfo)
	if mapInfoCache == nil {
		zaplogger.GetSugar().Infof("%s %d %s", this.user.GetUserID(), this.user.GetID(), "fail no map enter")
		this.errcode = cs_message.ErrCode_ERROR
		this.EndTrans()
		return
	}

	mapInfo := mapInfoCache.(*temporary.MapInfo)
	startAoi := &ss_message.StartAoi{
		UserID:   proto.String(this.user.GetUserID()),
		SceneIdx: proto.Int32(mapInfo.SceneIdx),
	}
	cluster.PostMessage(mapInfo.MapAddr, startAoi)

	zaplogger.GetSugar().Infof("%s %d %s", this.user.GetUserID(), this.user.GetID(), "Strat Aoi OK")
	this.errcode = cs_message.ErrCode_OK
	this.EndTrans()
}

func (this *transactionStartAoi) End() {
	if this.errcode == cs_message.ErrCode_OK {
		this.user.Reply(this.req.GetSeriNo(), this.resp)
	} else {
		this.user.ReplyErr(this.req.GetSeriNo(), this.req.GetCmd(), this.errcode)
	}
}

func (this *transactionStartAoi) Timeout() {
	this.errcode = cs_message.ErrCode_RETRY
}

func init() {
	user.RegisterTransFunc(cmdEnum.CS_StartAoi, func(user *user.User, msg *codecs.Message) transaction.Transaction {
		return &transactionStartAoi{
			user: user,
			req:  msg,
		}
	})
}
