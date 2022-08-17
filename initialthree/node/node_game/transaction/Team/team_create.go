package Team

import (
	"initialthree/cluster"
	codecs "initialthree/codec/cs"
	"initialthree/node/common/serverType"
	"initialthree/node/common/transaction"
	"initialthree/zaplogger"
	"time"

	"initialthree/node/node_game/temporary"
	"initialthree/node/node_game/user"
	"initialthree/protocol/cmdEnum"
	cs_message "initialthree/protocol/cs/message"
	"initialthree/protocol/ss/rpc"
	"initialthree/rpc/teamCreate"
)

type transactionTeamCreate struct {
	transaction.TransactionBase
	user    *user.User
	req     *codecs.Message
	resp    *cs_message.TeamCreateToC
	errcode cs_message.ErrCode
}

func (this *transactionTeamCreate) GetModuleName() string {
	return "Team"
}

func (this *transactionTeamCreate) Begin() {

	req := this.req.GetData().(*cs_message.TeamCreateToS)
	this.resp = &cs_message.TeamCreateToC{}
	zaplogger.GetSugar().Debugf("%s %s %v", this.user.GetUserLogName(), "teamCreate", req)

	teamInfo := this.user.GetTemporary(temporary.TempTeamInfo)
	if teamInfo != nil {
		zaplogger.GetSugar().Debugf("%s %s %v", this.user.GetUserLogName(), "create team, but in team", teamInfo.(*temporary.TeamInfo))
		this.errcode = cs_message.ErrCode_Team_AlreadyInTeam
		this.EndTrans()
		return
	}

	teamAddr, err := cluster.Random(serverType.Team)
	if err != nil {
		zaplogger.GetSugar().Errorf("---- no team server:%v", err)
		this.errcode = cs_message.ErrCode_ERROR
		this.EndTrans()
		return
	}

	teamCreateReq := &rpc.TeamCreateReq{
		Player: this.user.PackTeamPlayer(),
		Target: req.GetTarget(),
	}

	this.AsynWrap(teamCreate.AsynCall)(teamAddr, teamCreateReq, 8*time.Second, func(resp *rpc.TeamCreateResp, e error) {
		if resp.GetErrCode() == cs_message.ErrCode_OK {
			teamInfo := temporary.NewTeamInfo(this.user, teamAddr, resp.GetTeamID())
			teamInfo.IsHeader = true
			this.user.SetTemporary(temporary.TempTeamInfo, teamInfo)
		}

		this.errcode = resp.GetErrCode()
		this.EndTrans()
	})
}

func (this *transactionTeamCreate) End() {
	if this.errcode == cs_message.ErrCode_OK {
		this.user.Reply(this.req.GetSeriNo(), this.resp)
	} else {
		this.user.ReplyErr(this.req.GetSeriNo(), this.req.GetCmd(), this.errcode)
	}

}

func (this *transactionTeamCreate) Timeout() {
	this.errcode = cs_message.ErrCode_RETRY
}

func init() {
	user.RegisterTransFunc(cmdEnum.CS_TeamCreate, func(user *user.User, msg *codecs.Message) transaction.Transaction {
		return &transactionTeamCreate{
			user: user,
			req:  msg,
		}
	})
}
