
package ss
import (
	"initialthree/codec/pb"
	"initialthree/protocol/ss/ssmessage"
	"initialthree/protocol/ss/rpc"
)

func init() {
	//普通消息
	pb.Register("ss",&ssmessage.Echo{},1001)
	pb.Register("ss",&ssmessage.KickGameUser{},1002)
	pb.Register("ss",&ssmessage.KickIPGate{},1003)
	pb.Register("ss",&ssmessage.ReportGate{},1012)
	pb.Register("ss",&ssmessage.MapHeartbeat{},1014)
	pb.Register("ss",&ssmessage.UserLoginToDir{},1021)
	pb.Register("ss",&ssmessage.Relay{},1022)
	pb.Register("ss",&ssmessage.ReportStatus{},1023)
	pb.Register("ss",&ssmessage.StartAoi{},1024)
	pb.Register("ss",&ssmessage.MapToWorld{},1031)
	pb.Register("ss",&ssmessage.WorldBroadcastToMap{},1032)
	pb.Register("ss",&ssmessage.FunctionSwitchReload{},1033)
	pb.Register("ss",&ssmessage.FirewallUpdate{},1041)
	pb.Register("ss",&ssmessage.RankEnd{},1043)
	pb.Register("ss",&ssmessage.ScarsIngrainUpdate{},1044)
	pb.Register("ss",&ssmessage.MailUpdate{},1046)

	//rpc请求
	pb.Register("rpc_req",&rpc.EchoReq{},1001)
	pb.Register("rpc_req",&rpc.SynctokenReq{},1002)
	pb.Register("rpc_req",&rpc.EnterWorldReq{},1011)
	pb.Register("rpc_req",&rpc.EnterMapReq{},1012)
	pb.Register("rpc_req",&rpc.LeaveMapReq{},1013)
	pb.Register("rpc_req",&rpc.MoveReq{},1014)
	pb.Register("rpc_req",&rpc.WorldObjPushReq{},1021)
	pb.Register("rpc_req",&rpc.WorldObjDoReq{},1022)
	pb.Register("rpc_req",&rpc.GetLogicTimeReq{},1031)
	pb.Register("rpc_req",&rpc.TeamCreateReq{},1051)
	pb.Register("rpc_req",&rpc.TeamGetNearTeamReq{},1052)
	pb.Register("rpc_req",&rpc.TeamGetNearPlayerReq{},1053)
	pb.Register("rpc_req",&rpc.TeamPlayerLeaveReq{},1054)
	pb.Register("rpc_req",&rpc.TeamKickPlayerReq{},1055)
	pb.Register("rpc_req",&rpc.TeamJoinApplyReq{},1056)
	pb.Register("rpc_req",&rpc.TeamJoinReplyReq{},1057)
	pb.Register("rpc_req",&rpc.TeamPlayerGetFromGameReq{},1058)
	pb.Register("rpc_req",&rpc.TeamDismissReq{},1059)
	pb.Register("rpc_req",&rpc.RankGetTopListReq{},1071)
	pb.Register("rpc_req",&rpc.RankSetScoreReq{},1072)
	pb.Register("rpc_req",&rpc.RankCreateReq{},1074)
	pb.Register("rpc_req",&rpc.RankDeleteScoreReq{},1075)
	pb.Register("rpc_req",&rpc.RankGetRankReq{},1076)

	//rpc响应
	pb.Register("rpc_resp",&rpc.EchoResp{},1001)
	pb.Register("rpc_resp",&rpc.SynctokenResp{},1002)
	pb.Register("rpc_resp",&rpc.EnterWorldResp{},1011)
	pb.Register("rpc_resp",&rpc.EnterMapResp{},1012)
	pb.Register("rpc_resp",&rpc.LeaveMapResp{},1013)
	pb.Register("rpc_resp",&rpc.MoveResp{},1014)
	pb.Register("rpc_resp",&rpc.WorldObjPushResp{},1021)
	pb.Register("rpc_resp",&rpc.WorldObjDoResp{},1022)
	pb.Register("rpc_resp",&rpc.GetLogicTimeResp{},1031)
	pb.Register("rpc_resp",&rpc.TeamCreateResp{},1051)
	pb.Register("rpc_resp",&rpc.TeamGetNearTeamResp{},1052)
	pb.Register("rpc_resp",&rpc.TeamGetNearPlayerResp{},1053)
	pb.Register("rpc_resp",&rpc.TeamPlayerLeaveResp{},1054)
	pb.Register("rpc_resp",&rpc.TeamKickPlayerResp{},1055)
	pb.Register("rpc_resp",&rpc.TeamJoinApplyResp{},1056)
	pb.Register("rpc_resp",&rpc.TeamJoinReplyResp{},1057)
	pb.Register("rpc_resp",&rpc.TeamPlayerGetFromGameResp{},1058)
	pb.Register("rpc_resp",&rpc.TeamDismissResp{},1059)
	pb.Register("rpc_resp",&rpc.RankGetTopListResp{},1071)
	pb.Register("rpc_resp",&rpc.RankSetScoreResp{},1072)
	pb.Register("rpc_resp",&rpc.RankCreateResp{},1074)
	pb.Register("rpc_resp",&rpc.RankDeleteScoreResp{},1075)
	pb.Register("rpc_resp",&rpc.RankGetRankResp{},1076)

}
