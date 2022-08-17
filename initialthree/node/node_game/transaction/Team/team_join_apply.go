package Team

import (
	"github.com/gogo/protobuf/proto"
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
	"initialthree/rpc/teamJoinApply"
)

type transactionTeamJoinApply struct {
	transaction.TransactionBase
	user    *user.User
	req     *codecs.Message
	resp    *cs_message.TeamJoinApplyToC
	errcode cs_message.ErrCode
}

func (this *transactionTeamJoinApply) GetModuleName() string {
	return "Team"
}

func (this *transactionTeamJoinApply) Begin() {

	req := this.req.GetData().(*cs_message.TeamJoinApplyToS)
	this.resp = &cs_message.TeamJoinApplyToC{}
	zaplogger.GetSugar().Debugf("%s %s %v", this.user.GetUserLogName(), "teamJoinApply", req)

	tempTeam := this.user.GetTemporary(temporary.TempTeamInfo)
	if tempTeam != nil {
		zaplogger.GetSugar().Debugf("%s %s %v", this.user.GetUserLogName(), "teamJoinApply, but in team", tempTeam.(*temporary.TeamInfo))
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

	msg := &rpc.TeamJoinApplyReq{
		TeamID: proto.Uint32(req.GetTeamID()),
		Player: this.user.PackTeamPlayer(),
	}

	this.AsynWrap(teamJoinApply.AsynCall)(teamAddr, msg, 8*time.Second, func(resp *rpc.TeamJoinApplyResp, e error) {
		if e != nil {
			this.errcode = cs_message.ErrCode_ERROR
			this.EndTrans()
			return
		}
		this.errcode = resp.GetErrCode()
		this.EndTrans()
	})

}

func (this *transactionTeamJoinApply) End() {
	if this.errcode == cs_message.ErrCode_OK {
		this.user.Reply(this.req.GetSeriNo(), this.resp)
	} else {
		this.user.ReplyErr(this.req.GetSeriNo(), this.req.GetCmd(), this.errcode)
	}

}

func (this *transactionTeamJoinApply) Timeout() {
	this.errcode = cs_message.ErrCode_RETRY
}

func init() {
	user.RegisterTransFunc(cmdEnum.CS_TeamJoinApply, func(user *user.User, msg *codecs.Message) transaction.Transaction {
		return &transactionTeamJoinApply{
			user: user,
			req:  msg,
		}
	})
}
