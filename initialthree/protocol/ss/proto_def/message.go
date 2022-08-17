package proto_def

type St struct {
	Name      string
	MessageID int
}

var SS_message = []St{
	{"echo", 1},
	{"kickGameUser", 2},
	{"kickIPGate", 3}, // 通过ip踢人，通知到gate

	{"reportGate", 12},
	{"mapHeartbeat", 14},

	{"userLoginToDir", 21}, //用户登陆上报给dir
	{"relay", 22},          //透传
	{"reportStatus", 23},
	{"startAoi", 24},

	{"mapToWorld", 31},
	{"worldBroadcastToMap", 32},
	{"functionSwitchReload", 33},

	{"firewallUpdate", 41},
	{"rankEnd", 43},            // 排行榜结束，通知各game玩家拉取奖励
	{"scarsIngrainUpdate", 44}, // 战痕印刻副本已更新

	{"mailUpdate", 46},
}
