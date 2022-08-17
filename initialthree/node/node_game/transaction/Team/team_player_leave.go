package Team

import (
	codecs "initialthree/codec/cs"
	"initialthree/node/common/transaction"
	"initialthree/zaplogger"

	"initialthree/node/node_game/temporary"
	"initialthree/node/node_game/user"
	"initialthree/protocol/cmdEnum"
	cs_message "initialthree/protocol/cs/message"
)

type transactionTeamPlayerLeave struct {
	transaction.TransactionBase
	user    *user.User
	req     *codecs.Message
	resp    *cs_message.TeamPlayerLeaveToC
	errcode cs_message.ErrCode
}

func (this *transactionTeamPlayerLeave) GetModuleName() string {
	return "Team"
}

func (this *transactionTeamPlayerLeave) Begin() {

	req := this.req.GetData().(*cs_message.TeamPlayerLeaveToS)
	this.resp = &cs_message.TeamPlayerLeaveToC{}
	zaplogger.GetSugar().Debugf("%s %s %v", this.user.GetUserLogName(), "teamPlayerLeave", req)

	tempTeam := this.user.GetTemporary(temporary.TempTeamInfo)
	if tempTeam == nil {
		zaplogger.GetSugar().Debugf("%s %s", this.user.GetUserLogName(), "teamPlayerLeave, but no team")
		this.errcode = cs_message.ErrCode_Team_NotTeamMember
		this.EndTrans()
		return
	}

	teamInfo := tempTeam.(*temporary.TeamInfo)
	teamInfo.LeaveTeam(func(errCode cs_message.ErrCode) {
		this.errcode = errCode
		this.EndTrans()
	})

}

func (this *transactionTeamPlayerLeave) End() {
	if this.errcode == cs_message.ErrCode_OK {
		this.user.Reply(this.req.GetSeriNo(), this.resp)
	} else {
		this.user.ReplyErr(this.req.GetSeriNo(), this.req.GetCmd(), this.errcode)
	}

}

func (this *transactionTeamPlayerLeave) Timeout() {
	this.errcode = cs_message.ErrCode_RETRY
}

func init() {
	user.RegisterTransFunc(cmdEnum.CS_TeamPlayerLeave, func(user *user.User, msg *codecs.Message) transaction.Transaction {
		return &transactionTeamPlayerLeave{
			user: user,
			req:  msg,
		}
	})
}
