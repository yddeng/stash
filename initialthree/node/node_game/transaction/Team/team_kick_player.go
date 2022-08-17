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
	"initialthree/rpc/teamKickPlayer"
)

type transactionTeamKickPlayer struct {
	transaction.TransactionBase
	user    *user.User
	req     *codecs.Message
	resp    *cs_message.TeamKickPlayerToC
	errcode cs_message.ErrCode
}

func (this *transactionTeamKickPlayer) GetModuleName() string {
	return "Team"
}

func (this *transactionTeamKickPlayer) Begin() {

	req := this.req.GetData().(*cs_message.TeamKickPlayerToS)
	this.resp = &cs_message.TeamKickPlayerToC{}
	zaplogger.GetSugar().Debugf("%s %s %v", this.user.GetUserLogName(), "teamKickPlayer", req)

	if this.user.GetID() == req.GetKickPlayerID() {
		zaplogger.GetSugar().Debugf("%s %s", this.user.GetUserLogName(), "teamKickPlayer, can't kick self")
		this.errcode = cs_message.ErrCode_ERROR
		this.EndTrans()
		return
	}

	tempTeam := this.user.GetTemporary(temporary.TempTeamInfo)
	if tempTeam == nil {
		zaplogger.GetSugar().Debugf("%s %s", this.user.GetUserLogName(), "teamKickPlayer, but no team")
		this.errcode = cs_message.ErrCode_Team_NotTeamMember
		this.EndTrans()
		return
	}

	teamInfo := tempTeam.(*temporary.TeamInfo)
	if !teamInfo.IsHeader {
		zaplogger.GetSugar().Debugf("%s %s", this.user.GetUserLogName(), "teamKickPlayer, but is not header")
		this.errcode = cs_message.ErrCode_Team_NotHeader
		this.EndTrans()
		return
	}

	msg := &rpc.TeamKickPlayerReq{
		KickID:   proto.Uint64(req.GetKickPlayerID()),
		HeaderID: proto.Uint64(this.user.GetID()),
	}

	this.AsynWrap(teamKickPlayer.AsynCall)(teamInfo.TeamAddr, msg, 8*time.Second, func(resp *rpc.TeamKickPlayerResp, e error) {
		if e != nil {
			this.errcode = cs_message.ErrCode_ERROR
			this.EndTrans()
			return
		}
		this.errcode = resp.GetErrCode()
		this.EndTrans()
	})

}

func (this *transactionTeamKickPlayer) End() {
	if this.errcode == cs_message.ErrCode_OK {
		this.user.Reply(this.req.GetSeriNo(), this.resp)
	} else {
		this.user.ReplyErr(this.req.GetSeriNo(), this.req.GetCmd(), this.errcode)
	}

}

func (this *transactionTeamKickPlayer) Timeout() {
	this.errcode = cs_message.ErrCode_RETRY
}

func init() {
	user.RegisterTransFunc(cmdEnum.CS_TeamKickPlayer, func(user *user.User, msg *codecs.Message) transaction.Transaction {
		return &transactionTeamKickPlayer{
			user: user,
			req:  msg,
		}
	})
}
