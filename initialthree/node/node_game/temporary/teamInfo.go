package temporary

import (
	"github.com/golang/protobuf/proto"
	"initialthree/cluster/addr"
	"initialthree/zaplogger"

	"initialthree/protocol/cs/message"
	"initialthree/protocol/ss/rpc"
	"initialthree/rpc/teamPlayerLeave"
	"time"
)

/*
	team 临时数据
	new    : create(header), teamPlayerGet
	delete : leave, kick, dismiss
	update : headerChanged,
*/

type TeamInfo struct {
	UserI    UserI
	TeamID   uint32
	TeamAddr addr.LogicAddr
	IsHeader bool
	pInfo    *message.TeamPlayer
	isDirty  bool // 需要向team同步的，外观
}

func NewTeamInfo(userI UserI, logicAddr addr.LogicAddr, teamID uint32) *TeamInfo {
	return &TeamInfo{
		UserI:    userI,
		TeamID:   teamID,
		TeamAddr: logicAddr,
	}
}

func (this *TeamInfo) UpdatePlayer(pInfo *message.TeamPlayer) {
	this.pInfo = pInfo
	this.isDirty = true
}

func (this *TeamInfo) LeaveTeam(cb func(errCode message.ErrCode)) {
	teamPlayerLeave.AsynCall(this.TeamAddr, &rpc.TeamPlayerLeaveReq{
		PlayerID: proto.Uint64(this.UserI.GetID()),
		TeamID:   proto.Uint32(this.TeamID),
	}, time.Second*5, func(resp *rpc.TeamPlayerLeaveResp, e error) {
		if e == nil {
			if resp.GetErrCode() == message.ErrCode_OK {
				zaplogger.GetSugar().Debugf("%v %v leaveTeam ok", this.UserI.GetUserID(), this.UserI.GetID())
				this.UserI.ClearTemporary(TempTeamInfo)
			}
			if cb != nil {
				cb(resp.GetErrCode())
			}
		} else {
			zaplogger.GetSugar().Debugf("%v %v leaveTeam fail:%v", this.UserI.GetUserID(), this.UserI.GetID(), e)
			if cb != nil {
				cb(message.ErrCode_ERROR)
			}
		}
	})
}

func (this *TeamInfo) UserDisconnect() {
	this.LeaveTeam(nil)
}

func (this *TeamInfo) UserLogout() {
	this.LeaveTeam(nil)
}

func (this *TeamInfo) Tick(now time.Time) {
	if this.isDirty {
		this.isDirty = false
	}
}
