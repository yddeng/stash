package handler

import (
	"initialthree/cluster/addr"
	"initialthree/codec/cs"
	"initialthree/node/node_game/temporary"
	"initialthree/node/node_game/user"
	"initialthree/protocol/cs/message"
	"initialthree/zaplogger"
)

// 自己被踢， 清理team数据
func teamKickPlayerNotifyHook(from addr.LogicAddr, u *user.User, msg *cs.Message) bool {

	zaplogger.GetSugar().Debugf("%s teamKickPlayerNotify ", u.GetUserLogName())

	tempTeam := u.GetTemporary(temporary.TempTeamInfo)
	if tempTeam == nil {
		zaplogger.GetSugar().Errorf("%s temp team is nil", u.GetUserLogName())
		return true
	}

	u.ClearTemporary(temporary.TempTeamInfo)

	return true
}

// 队长变更给自己
func teamHeaderChangedNotifyHook(from addr.LogicAddr, u *user.User, msg *cs.Message) bool {

	//req := msg.GetData().(*message.TeamHeaderChangedNotifyToC)
	zaplogger.GetSugar().Debugf("%s teamHeaderChangedNotify ", u.GetUserLogName())

	tempTeam := u.GetTemporary(temporary.TempTeamInfo)
	if tempTeam == nil {
		zaplogger.GetSugar().Errorf("%s temp team is nil", u.GetUserLogName())
		return false
	}

	teamInfo := tempTeam.(*temporary.TeamInfo)
	teamInfo.IsHeader = true

	return true
}

// 队伍解散 通知给自己
func teamDismissNotifyHook(from addr.LogicAddr, u *user.User, msg *cs.Message) bool {

	//req := msg.GetData().(*message.TeamHeaderChangedNotifyToC)
	zaplogger.GetSugar().Debugf("%s teamDismissNotify ", u.GetUserLogName())

	tempTeam := u.GetTemporary(temporary.TempTeamInfo)
	if tempTeam == nil {
		zaplogger.GetSugar().Errorf("%s temp team is nil", u.GetUserLogName())
		return true
	}

	u.ClearTemporary(temporary.TempTeamInfo)
	return true
}

func init() {
	user.RegisterHooker(&message.TeamHeaderChangedNotifyToC{}, teamHeaderChangedNotifyHook)
	user.RegisterHooker(&message.TeamKickPlayerNotifyToC{}, teamKickPlayerNotifyHook)
	user.RegisterHooker(&message.TeamDismissNotifyToC{}, teamDismissNotifyHook)
}
