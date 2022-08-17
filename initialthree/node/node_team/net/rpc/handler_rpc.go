package rpc

import (
	"github.com/golang/protobuf/proto"
	"initialthree/cluster"
	"initialthree/node/node_team/team"
	"initialthree/protocol/cs/message"
	"initialthree/protocol/ss/rpc"
	"initialthree/rpc/teamCreate"
	"initialthree/rpc/teamDismiss"
	"initialthree/rpc/teamJoinApply"
	"initialthree/rpc/teamJoinReply"
	"initialthree/rpc/teamKickPlayer"
	"initialthree/rpc/teamPlayerGetFromGame"
	"initialthree/rpc/teamPlayerLeave"
	"initialthree/util"
)

/***********************************************************************************************************************
	队伍创建接口
***********************************************************************************************************************/
type teamCreateHandler struct{}

func (h teamCreateHandler) OnCall(replier *teamCreate.TeamCreateReplyer, arg *rpc.TeamCreateReq) {
	resp := &rpc.TeamCreateResp{
		ErrCode: message.ErrCode_OK.Enum(),
	}
	teamManager := team.Mgr()

	pId := arg.GetPlayer().GetPlayerID()
	gameAddr := replier.GetChannel().(*cluster.RPCChannel).PeerAddr()

	// 检查用户是否已有队伍
	player := teamManager.GetPlayer(pId)
	if player != nil && player.Team() != nil {
		util.Logger().Debugf("%s request create team, but already have team:%d\n", player.LogStr(), player.Team().TeamID())
		//resp.ErrCode = message.ErrCode_Team_AlreadyInTeam.Enum()
		//replier.Reply(resp)

		// 数据修正，通告玩家已经在队伍
		player.Team().SyncToClients()
		resp.TeamID = proto.Uint32(player.Team().TeamID())
		replier.Reply(resp)
		return
	}

	// 玩家
	if player == nil {
		player = team.NewPlayer(gameAddr, arg.GetPlayer())
		teamManager.AddPlayer(player)
	}

	// 队伍
	insTeam := team.NewTeam()
	teamManager.AddTeam(insTeam)

	insTeam.SetTarget(arg.GetTarget())
	insTeam.AddPlayer(player)

	util.Logger().Debugf("%s request create team %d ok\n", player.LogStr(), insTeam.TeamID())
	resp.TeamID = proto.Uint32(insTeam.TeamID())
	replier.Reply(resp)
}

/***********************************************************************************************************************
	队伍解散接口
***********************************************************************************************************************/
type teamDismissHandler struct{}

func (h teamDismissHandler) OnCall(replier *teamDismiss.TeamDismissReplyer, arg *rpc.TeamDismissReq) {
	resp := &rpc.TeamDismissResp{
		ErrCode: message.ErrCode_OK.Enum(),
	}
	teamManager := team.Mgr()

	pId := arg.GetPlayerID()
	// 检查用户是否有队伍 ,todo 还需要什么处理
	player := teamManager.GetPlayer(pId)
	if player == nil || player.Team() == nil {
		util.Logger().Debugf("%d request teamDismiss, but no team\n", pId)
		resp.ErrCode = message.ErrCode_Team_NotTeamMember.Enum()
		replier.Reply(resp)
		return
	}

	insTeam := player.Team()
	// 是否是队长
	if pId != insTeam.Header() {
		util.Logger().Debugf("%s request teamDismiss, but is not header\n", player.LogStr())
		resp.ErrCode = message.ErrCode_Team_NotHeader.Enum()
		replier.Reply(resp)
		return
	}

	// 通知给除队长以外的所有
	insTeam.NotifyAllExceptMe(&message.TeamDismissNotifyToC{}, player)
	// 删除玩家实例
	insTeam.RangePlayer(func(p *team.Player) bool {
		teamManager.RemovePlayer(p.PlayerID())
		return true
	})
	// 删除队伍实例
	teamManager.RemoveTeam(insTeam.TeamID())

	util.Logger().Debugf("%s request teamDismiss %d ok\n", player.LogStr(), insTeam.TeamID())
	replier.Reply(resp)
}

/***********************************************************************************************************************
	离开队伍接口
***********************************************************************************************************************/
type teamPlayerLeaveHandler struct{}

func (h teamPlayerLeaveHandler) OnCall(replier *teamPlayerLeave.TeamPlayerLeaveReplyer, arg *rpc.TeamPlayerLeaveReq) {
	resp := &rpc.TeamPlayerLeaveResp{
		ErrCode: message.ErrCode_OK.Enum(),
	}
	teamManager := team.Mgr()

	pId := arg.GetPlayerID()
	tId := arg.GetTeamID()

	// 检查用户是否有队伍
	player := teamManager.GetPlayer(pId)
	if player == nil || player.Team() == nil { // 不存在 直接返回成功
		util.Logger().Debugf("%d request teamPlayerLeave, but no team\n", pId)
		replier.Reply(resp)
		return
	}

	insTeam := player.Team()

	// id 不一致，仍然离队
	if insTeam.TeamID() != tId {
		util.Logger().Debugf("%s request teamPlayerLeave, tId failed, ins:%d req:%d \n", player.LogStr(), insTeam.TeamID(), tId)
	}

	insTeam.RemovePlayer(player)
	teamManager.RemovePlayer(pId)

	util.Logger().Debugf("%s teamPlayerLeave %d team ok\n", player.LogStr(), insTeam.TeamID())
	replier.Reply(resp)
}

/***********************************************************************************************************************
	踢人接口
***********************************************************************************************************************/
type teamKickPlayerHandler struct{}

func (h teamKickPlayerHandler) OnCall(replier *teamKickPlayer.TeamKickPlayerReplyer, arg *rpc.TeamKickPlayerReq) {
	resp := &rpc.TeamKickPlayerResp{
		ErrCode: message.ErrCode_OK.Enum(),
	}
	teamManager := team.Mgr()

	// 检查用户是否有队伍
	header := teamManager.GetPlayer(arg.GetHeaderID())
	if header == nil || header.Team() == nil { // 不存在 直接返回成功
		util.Logger().Debugf("%d request teamKickPlayer, but no team\n", arg.GetHeaderID())
		resp.ErrCode = message.ErrCode_Team_NotTeamMember.Enum()
		replier.Reply(resp)
		return
	}

	insTeam := header.Team()
	// 是否是队长
	if arg.GetHeaderID() != insTeam.Header() {
		util.Logger().Debugf("%s request teamKickPlayer, but is not header\n", header.LogStr())
		resp.ErrCode = message.ErrCode_Team_NotHeader.Enum()
		replier.Reply(resp)
		return
	}

	kickPlayer := teamManager.GetPlayer(arg.GetKickID())
	if kickPlayer != nil && kickPlayer.Team() == insTeam {
		// todo 踢人的时候要不要通知所有人呢
		insTeam.RemovePlayer(kickPlayer)
		teamManager.RemovePlayer(kickPlayer.PlayerID())
		kickPlayer.SendMsg(&message.TeamKickPlayerNotifyToC{})
	}

	util.Logger().Debugf("%s teamKickPlayer %d team ok\n", header.LogStr(), arg.GetKickID())
	replier.Reply(resp)
}

/***********************************************************************************************************************
	申请入队接口
***********************************************************************************************************************/
type teamJoinApplyHandler struct{}

func (h teamJoinApplyHandler) OnCall(replier *teamJoinApply.TeamJoinApplyReplyer, arg *rpc.TeamJoinApplyReq) {
	teamManager := team.Mgr()
	gameAddr := replier.GetChannel().(*cluster.RPCChannel).PeerAddr()

	// 队伍是否存在
	insTeam := teamManager.GetTeam(arg.GetTeamID())
	if insTeam == nil {
		util.Logger().Debugf("(%s %d) request teamJoinApply, team %d in nil \n",
			arg.GetPlayer().GetUserID(), arg.GetPlayer().GetPlayerID(), arg.GetTeamID())
		replier.Reply(&rpc.TeamJoinApplyResp{ErrCode: message.ErrCode_Team_TeamNotExist.Enum()})
		return
	}

	// 玩家是否已有队伍
	player := teamManager.GetPlayer(arg.GetPlayer().GetPlayerID())
	if player != nil && player.Team() != nil {
		util.Logger().Debugf("%s request teamJoinApply, but in team\n", player.LogStr())
		replier.Reply(&rpc.TeamJoinApplyResp{ErrCode: message.ErrCode_Team_AlreadyInTeam.Enum()})
		return
	}

	if player == nil {
		player = team.NewPlayer(gameAddr, arg.GetPlayer())
		teamManager.AddPlayer(player)
	}

	// 队伍只有处于待命状态才能加入
	if insTeam.Status() != message.TeamStatus_Standby {
		util.Logger().Debugf("%s request teamJoinApply:%d, team can't join now status %s\n",
			player.LogStr(), arg.GetTeamID(), insTeam.Status().String())
		replier.Reply(&rpc.TeamJoinApplyResp{ErrCode: message.ErrCode_Team_TeamStatusErr.Enum()})
		return
	}

	// 队伍人数已满
	if insTeam.PlayerFull() {
		util.Logger().Debugf("%s request teamJoinApply:%d, team player full\n", player.LogStr(), arg.GetTeamID())
		replier.Reply(&rpc.TeamJoinApplyResp{ErrCode: message.ErrCode_Team_PlayerFull.Enum()})
		return
	}

	// todo 入队条件

	replier.Reply(&rpc.TeamJoinApplyResp{ErrCode: message.ErrCode_OK.Enum()})
	util.Logger().Debugf("%s request teamJoinApply %d team ok,\n", player.LogStr(), insTeam.TeamID())

	// 加入team 申请
	insTeam.PlayerJoinApply(player)
}

/***********************************************************************************************************************
	申请入队处理接口
***********************************************************************************************************************/
type teamJoinReplyHandler struct{}

func (h teamJoinReplyHandler) OnCall(replier *teamJoinReply.TeamJoinReplyReplyer, arg *rpc.TeamJoinReplyReq) {
	teamManager := team.Mgr()

	// 检查用户是否有队伍
	header := teamManager.GetPlayer(arg.GetMineID())
	if header == nil || header.Team() == nil { // 不存在 直接返回成功
		util.Logger().Debugf("%d request teamJoinReply, but no team\n", arg.GetMineID())
		replier.Reply(&rpc.TeamJoinReplyResp{ErrCode: message.ErrCode_Team_NotTeamMember.Enum()})
		return
	}

	insTeam := header.Team()
	// 是否是队长
	if arg.GetMineID() != insTeam.Header() {
		util.Logger().Debugf("%s request teamJoinReply, but is not header\n", header.LogStr())
		replier.Reply(&rpc.TeamJoinReplyResp{ErrCode: message.ErrCode_Team_NotHeader.Enum()})
		return
	}

	player := teamManager.GetPlayer(arg.GetAgreeID())
	// 玩家不存在
	if player == nil {
		util.Logger().Debugf("%d request teamJoinReply, but player %d is nil\n", arg.GetMineID(), arg.GetAgreeID())
		replier.Reply(&rpc.TeamJoinReplyResp{ErrCode: message.ErrCode_Team_PlayerNotExist.Enum()})
		return
	}

	if arg.GetAgree() {
		// 玩家已经加入其他队伍
		if player.Team() != nil {
			util.Logger().Debugf("%d request teamJoinReply, but player %d in another %d team\n", arg.GetMineID(), arg.GetAgreeID(), player.Team().TeamID())
			replier.Reply(&rpc.TeamJoinReplyResp{ErrCode: message.ErrCode_Team_AlreadyInTeam.Enum()})
			return
		}
		// 队伍只有处于待命状态才能加入
		if insTeam.Status() != message.TeamStatus_Standby {
			util.Logger().Debugf("%s request teamJoinReply:%d, team can't join now status %s\n", player.LogStr(), insTeam.TeamID(), insTeam.Status().String())
			replier.Reply(&rpc.TeamJoinReplyResp{ErrCode: message.ErrCode_Team_TeamStatusErr.Enum()})
			return
		}
		// 队伍人数已满
		if insTeam.PlayerFull() {
			util.Logger().Debugf("%s request teamJoinReply:%d, team player full\n", player.LogStr(), arg.GetTeamID())
			replier.Reply(&rpc.TeamJoinReplyResp{ErrCode: message.ErrCode_Team_PlayerFull.Enum()})
			return
		}

		// game 上拉取最新 teamPlayer
		get := &rpc.TeamPlayerGetFromGameReq{
			UID:    proto.String(player.UserID()),
			TeamID: proto.Uint32(insTeam.TeamID()),
		}
		teamPlayerGetFromGame.AsynCall(player.Game(), get, time.Second*5, func(resp *rpc.TeamPlayerGetFromGameResp, e error) {
			if e != nil || resp.GetErrCode() != message.ErrCode_OK {
				util.Logger().Debugf("%s teamPlayerGetFromGame failed %s\n", player.LogStr(), resp.ErrCode.String())
				replier.Reply(&rpc.TeamJoinReplyResp{ErrCode: message.ErrCode_User_OffLine.Enum()})
			} else {
				util.Logger().Debugf("%s request teamJoinReply %d team ok,\n", header.LogStr(), insTeam.TeamID())
				replier.Reply(&rpc.TeamJoinReplyResp{ErrCode: message.ErrCode_OK.Enum()})
				player.UpdateInfo(resp.GetPlayer()) // 更新外观
				insTeam.PlayerJoinReply(player, true)
			}
		})
	} else {
		replier.Reply(&rpc.TeamJoinReplyResp{ErrCode: message.ErrCode_OK.Enum()})
		util.Logger().Debugf("%s request teamJoinReply %d team ok,\n", header.LogStr(), insTeam.TeamID())
		insTeam.PlayerJoinReply(player, false)
	}
}

func init() {
	teamCreate.Register(&teamCreateHandler{})
	teamDismiss.Register(&teamDismissHandler{})
	teamPlayerLeave.Register(&teamPlayerLeaveHandler{})
	teamKickPlayer.Register(&teamKickPlayerHandler{})
	teamJoinApply.Register(&teamJoinApplyHandler{})
	teamJoinReply.Register(&teamJoinReplyHandler{})
}
