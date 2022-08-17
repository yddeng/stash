package Team

import (
	"github.com/golang/protobuf/proto"
	codecs "initialthree/codec/cs"
	"initialthree/node/common/transaction"
	"initialthree/zaplogger"
	"time"

	"initialthree/node/node_game/temporary"
	"initialthree/node/node_game/user"
	"initialthree/protocol/cmdEnum"
	cs_message "initialthree/protocol/cs/message"
	"initialthree/protocol/ss/rpc"
	"initialthree/rpc/teamDismiss"
)

type transactionTeamDismiss struct {
	transaction.TransactionBase
	user    *user.User
	req     *codecs.Message
	resp    *cs_message.TeamDismissToC
	errcode cs_message.ErrCode
}

func (this *transactionTeamDismiss) GetModuleName() string {
	return "Team"
}

func (this *transactionTeamDismiss) Begin() {

	req := this.req.GetData().(*cs_message.TeamDismissToS)
	this.resp = &cs_message.TeamDismissToC{}
	zaplogger.GetSugar().Debugf("%s %s %v", this.user.GetUserLogName(), "teamDismiss", req)

	tempTeam := this.user.GetTemporary(temporary.TempTeamInfo)
	if tempTeam == nil {
		zaplogger.GetSugar().Debugf("%s %s", this.user.GetUserLogName(), "teamDismiss, but no team")
		this.errcode = cs_message.ErrCode_Team_NotTeamMember
		this.EndTrans()
		return
	}

	teamInfo := tempTeam.(*temporary.TeamInfo)
	if !teamInfo.IsHeader {
		zaplogger.GetSugar().Debugf("%s %s ", this.user.GetUserLogName(), "teamDismiss, but is not header")
		this.errcode = cs_message.ErrCode_Team_NotHeader
		this.EndTrans()
		return
	}

	msg := &rpc.TeamDismissReq{
		TeamID:   proto.Uint32(teamInfo.TeamID),
		PlayerID: proto.Uint64(this.user.GetID()),
	}

	this.AsynWrap(teamDismiss.AsynCall)(teamInfo.TeamAddr, msg, 8*time.Second, func(resp *rpc.TeamDismissResp, e error) {
		if e != nil {
			this.errcode = cs_message.ErrCode_ERROR
			this.EndTrans()
			return
		}

		// 队伍 解散 ，清理
		if resp.GetErrCode() == cs_message.ErrCode_OK {
			this.user.ClearTemporary(temporary.TempTeamInfo)
		}

		this.errcode = resp.GetErrCode()
		this.EndTrans()
	})

}

func (this *transactionTeamDismiss) End() {
	if this.errcode == cs_message.ErrCode_OK {
		this.user.Reply(this.req.GetSeriNo(), this.resp)
	} else {
		this.user.ReplyErr(this.req.GetSeriNo(), this.req.GetCmd(), this.errcode)
	}

}

func (this *transactionTeamDismiss) Timeout() {
	this.errcode = cs_message.ErrCode_RETRY
}

func init() {
	user.RegisterTransFunc(cmdEnum.CS_TeamDismiss, func(user *user.User, msg *codecs.Message) transaction.Transaction {
		return &transactionTeamDismiss{
			user: user,
			req:  msg,
		}
	})
}
