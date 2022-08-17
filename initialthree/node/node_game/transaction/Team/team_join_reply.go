package Team

import (
	"github.com/gogo/protobuf/proto"
	codecs "initialthree/codec/cs"
	"initialthree/node/common/transaction"
	"initialthree/zaplogger"
	"time"

	"initialthree/node/node_game/temporary"
	"initialthree/node/node_game/user"
	"initialthree/protocol/cmdEnum"
	cs_message "initialthree/protocol/cs/message"
	"initialthree/protocol/ss/rpc"
	"initialthree/rpc/teamJoinReply"
)

type transactionTeamJoinReply struct {
	transaction.TransactionBase
	user    *user.User
	req     *codecs.Message
	resp    *cs_message.TeamJoinReplyToC
	errcode cs_message.ErrCode
}

func (this *transactionTeamJoinReply) GetModuleName() string {
	return "Team"
}

func (this *transactionTeamJoinReply) Begin() {

	req := this.req.GetData().(*cs_message.TeamJoinReplyToS)
	this.resp = &cs_message.TeamJoinReplyToC{}
	zaplogger.GetSugar().Debugf("%s %s %v", this.user.GetUserLogName(), "teamJoinReply", req)

	tempTeam := this.user.GetTemporary(temporary.TempTeamInfo)
	if tempTeam == nil {
		zaplogger.GetSugar().Debugf("%s %s", this.user.GetUserLogName(), "teamJoinReply, but no team")
		this.errcode = cs_message.ErrCode_Team_NotTeamMember
		this.EndTrans()
		return
	}

	teamInfo := tempTeam.(*temporary.TeamInfo)
	if !teamInfo.IsHeader {
		zaplogger.GetSugar().Debugf("%s %s ", this.user.GetUserLogName(), "teamJoinReply, but is not header")
		this.errcode = cs_message.ErrCode_Team_NotHeader
		this.EndTrans()
		return
	}

	msg := &rpc.TeamJoinReplyReq{
		Agree:   proto.Bool(req.GetAgree()),
		AgreeID: proto.Uint64(req.GetAgreeID()),
		TeamID:  proto.Uint32(teamInfo.TeamID),
		MineID:  proto.Uint64(this.user.GetID()),
	}
	this.AsynWrap(teamJoinReply.AsynCall)(teamInfo.TeamAddr, msg, 8*time.Second, func(resp *rpc.TeamJoinReplyResp, e error) {
		if e != nil {
			this.errcode = cs_message.ErrCode_ERROR
			this.EndTrans()
			return
		}
		this.errcode = resp.GetErrCode()
		this.EndTrans()
	})

}

func (this *transactionTeamJoinReply) End() {
	if this.errcode == cs_message.ErrCode_OK {
		this.user.Reply(this.req.GetSeriNo(), this.resp)
	} else {
		this.user.ReplyErr(this.req.GetSeriNo(), this.req.GetCmd(), this.errcode)
	}

}

func (this *transactionTeamJoinReply) Timeout() {
	this.errcode = cs_message.ErrCode_RETRY
}

func init() {
	user.RegisterTransFunc(cmdEnum.CS_TeamJoinReply, func(user *user.User, msg *codecs.Message) transaction.Transaction {
		return &transactionTeamJoinReply{
			user: user,
			req:  msg,
		}
	})
}
