package Rank

import (
	"github.com/gogo/protobuf/proto"
	"initialthree/cluster/addr"
	codecs "initialthree/codec/cs"
	"initialthree/node/common/transaction"
	"initialthree/node/node_game/global"
	"initialthree/node/node_game/module"
	"initialthree/node/node_game/module/rankData"
	"initialthree/node/node_game/user"
	"initialthree/protocol/cmdEnum"
	cs_message "initialthree/protocol/cs/message"
	"initialthree/protocol/ss/rpc"
	"initialthree/zaplogger"
)

type transactionRankGetRank struct {
	transaction.TransactionBase
	user    *user.User
	req     *codecs.Message
	resp    *cs_message.RankGetRankToC
	errcode cs_message.ErrCode
}

func (this *transactionRankGetRank) GetModuleName() string {
	return "Rank"
}

func (this *transactionRankGetRank) Begin() {
	this.errcode = cs_message.ErrCode_OK

	msg := this.req.GetData().(*cs_message.RankGetRankToS)
	zaplogger.GetSugar().Debugf("%s RankGetRankToS %v", this.user.GetUserLogName(), msg)

	rankID := msg.GetRankID()
	this.resp = &cs_message.RankGetRankToC{
		RankID:  proto.Int32(rankID),
		Rank:    proto.Int32(0),
		Total:   proto.Int32(0),
		Percent: proto.Int32(0),
	}

	if rankID == 0 {
		this.errcode = cs_message.ErrCode_Request_Argument_Err
		this.EndTrans(this.resp, this.errcode)
		return
	}

	rankModule := this.user.GetSubModule(module.RankData).(*rankData.RankData)
	rankInfo := rankModule.GetRankInfo(rankID)
	if rankInfo == nil {
		// 无数据
		this.EndTrans(this.resp, this.errcode)
		return
	}

	logicAddr, err := addr.MakeLogicAddr(rankInfo.RankLogic)
	if err != nil {
		rankModule.DelRank(rankID)
		this.EndTrans(this.resp, this.errcode)
		return
	}

	rankData.CallGetRank(logicAddr, rankID, []uint64{this.user.GetID()}, func(results []*rpc.RankGetRankResult, err error) {
		if err != nil {
			zaplogger.GetSugar().Errorf("%s RankGetRankToS rankData callGetRank err %s", this.user.GetUserLogName(), err.Error())
			this.EndTrans(this.resp, this.errcode)
			return
		}

		result := results[0]
		// 0 没有玩家数据， 1 排行榜正在进行， 2 排行榜完结，且有排名
		switch result.GetCode() {
		case 0:
			zaplogger.GetSugar().Debugf("%s RankGetRankToS rankData rank %d is not exit", this.user.GetUserLogName(), rankID)
			rankModule.DelRank(rankID)
		case 1, 2:
			rank, total := result.GetRank(), result.GetTotal()
			this.resp.Rank = proto.Int32(rank)
			this.resp.Total = proto.Int32(rank)
			this.resp.Score = proto.Int32(result.GetScore())
			if rank > 0 && total > 0 {
				this.resp.Percent = proto.Int32(int32(global.RankPercent(rank, total) * 100))
			}
		default:
		}
		this.EndTrans(this.resp, this.errcode)
	})
}

func init() {
	user.RegisterTransFunc(cmdEnum.CS_RankGetRank, func(user *user.User, msg *codecs.Message) transaction.Transaction {
		return &transactionRankGetRank{
			user: user,
			req:  msg,
		}
	})
}
