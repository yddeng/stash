package Rank

import (
	"github.com/gogo/protobuf/proto"
	"initialthree/cluster/addr"
	codecs "initialthree/codec/cs"
	"initialthree/node/common/transaction"
	"initialthree/node/node_game/global"
	"initialthree/zaplogger"
	"time"

	"initialthree/node/node_game/user"
	"initialthree/protocol/cmdEnum"
	cs_message "initialthree/protocol/cs/message"
	"initialthree/protocol/ss/rpc"
	"initialthree/rpc/rankGetTopList"
)

type transactionRankGetTopList struct {
	transaction.TransactionBase
	user    *user.User
	req     *codecs.Message
	resp    *cs_message.RankGetTopListToC
	errcode cs_message.ErrCode
}

func (this *transactionRankGetTopList) GetModuleName() string {
	return "Rank"
}

func (this *transactionRankGetTopList) Begin() {
	this.errcode = cs_message.ErrCode_OK
	this.resp = &cs_message.RankGetTopListToC{}

	msg := this.req.GetData().(*cs_message.RankGetTopListToS)
	zaplogger.GetSugar().Debugf("%s rankGetTopList %v", this.user.GetUserLogName(), msg)

	rankId := msg.GetRankID()
	if rankId == 0 {
		zaplogger.GetSugar().Debugf("%s rankGetTopList failed: rankId %d", this.user.GetUserLogName(), rankId)
		this.errcode = cs_message.ErrCode_Request_Argument_Err
		this.EndTrans(this.resp, this.errcode)
		return
	}

	req := &rpc.RankGetTopListReq{
		Tos:    msg,
		RoleID: proto.Uint64(this.user.GetID()),
	}

	logicAddr, err := RankLogicAddr(rankId)
	if err != nil {
		zaplogger.GetSugar().Debugf("%s rankGetTopList failed: err %s", this.user.GetUserLogName(), err)
		this.errcode = cs_message.ErrCode_ERROR
		this.EndTrans(this.resp, this.errcode)
		return
	}

	this.callGetList(logicAddr, req)
}

func (this *transactionRankGetTopList) callGetList(logicAddr addr.LogicAddr, req *rpc.RankGetTopListReq) {
	this.AsynWrap(rankGetTopList.AsynCall)(logicAddr, req, 8*time.Second, func(resp *rpc.RankGetTopListResp, err error) {
		if err != nil {
			zaplogger.GetSugar().Error(err)
			this.errcode = cs_message.ErrCode_ERROR
			this.EndTrans(this.resp, this.errcode)
			return
		}

		if resp.GetRedirectRankAddr() != "" {
			req.Tos.GetLast = proto.Bool(false)
			req.Tos.RankID = proto.Int32(resp.GetCurRankID())

			logicAddr, err = addr.MakeLogicAddr(resp.GetRedirectRankAddr())
			if err != nil {
				zaplogger.GetSugar().Error(err)
				this.errcode = cs_message.ErrCode_ERROR
				this.EndTrans(this.resp, this.errcode)
				return
			}
			this.callGetList(logicAddr, req)
		} else {
			this.resp = resp.GetToc()
			rank, total := this.resp.GetRank(), this.resp.GetTotal()
			if rank > 0 && total > 0 {
				this.resp.Percent = proto.Int32(int32(global.RankPercent(rank, total) * 100))
			}
			this.EndTrans(this.resp, this.errcode)
		}
	})

}

func (this *transactionRankGetTopList) End() {
	if this.errcode == cs_message.ErrCode_OK {
		this.user.Reply(this.req.GetSeriNo(), this.resp)
	} else {
		this.user.ReplyErr(this.req.GetSeriNo(), this.req.GetCmd(), this.errcode)
	}
}

func (this *transactionRankGetTopList) Timeout() {
	this.errcode = cs_message.ErrCode_RETRY
}

func init() {
	user.RegisterTransFunc(cmdEnum.CS_RankGetTopList, func(user *user.User, msg *codecs.Message) transaction.Transaction {
		return &transactionRankGetTopList{
			user: user,
			req:  msg,
		}
	})
}
