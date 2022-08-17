package handler

import (
	"initialthree/cluster"
	"initialthree/node/node_game/temporary"
	"initialthree/node/node_game/user"
	"initialthree/protocol/cs/message"
	"initialthree/protocol/ss/rpc"
	"initialthree/rpc/teamPlayerGetFromGame"
	"initialthree/zaplogger"
)

/***********************************************************************************************************************
	玩家入队时拉取的最新外观
***********************************************************************************************************************/
type teamPlayerGetFromGameHandler struct{}

// 通过此消息，给玩家添加 team 数据。
func (_ *teamPlayerGetFromGameHandler) OnCall(replyer *teamPlayerGetFromGame.TeamPlayerGetFromGameReplyer, req *rpc.TeamPlayerGetFromGameReq) {
	zaplogger.GetSugar().Debugf("onCall teamPlayerGetFromGame %v", req)

	u := user.GetUser(req.GetUID())
	if u != nil && u.StatusOk() {

		// team 已经验证，没有队伍。故这里直接清理
		tempTeam := u.GetTemporary(temporary.TempTeamInfo)
		if tempTeam != nil {
			u.ClearTemporary(temporary.TempTeamInfo)
		}

		teamAddr := replyer.GetChannel().(*cluster.RPCChannel).PeerAddr()
		teamInfo := temporary.NewTeamInfo(u, teamAddr, req.GetTeamID())
		u.SetTemporary(temporary.TempTeamInfo, teamInfo)

		replyer.Reply(&rpc.TeamPlayerGetFromGameResp{
			ErrCode: message.ErrCode_OK.Enum(),
			Player:  u.PackTeamPlayer(),
		})

	} else {
		replyer.Reply(&rpc.TeamPlayerGetFromGameResp{
			ErrCode: message.ErrCode_User_OffLine.Enum(),
		})
	}
}

func init() {
	teamPlayerGetFromGame.Register(&teamPlayerGetFromGameHandler{})
}
