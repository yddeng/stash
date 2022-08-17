package proto_def

var SS_rpc = []St{
	{"echo", 1},
	{"synctoken", 2},

	{"enterWorld", 11},
	{"enterMap", 12},
	{"leaveMap", 13},
	{"move", 14},

	{"worldObjPush", 21}, // 由世界服推送给各地图服
	{"worldObjDo", 22},   // 由地图服请求世界服

	{"getLogicTime", 31}, // 获取逻辑时间

	{"teamCreate", 51},
	{"teamGetNearTeam", 52},
	{"teamGetNearPlayer", 53},
	{"teamPlayerLeave", 54},
	{"teamKickPlayer", 55},
	{"teamJoinApply", 56},         // 入队申请
	{"teamJoinReply", 57},         // 入队申请回复
	{"teamPlayerGetFromGame", 58}, // 从 game 拉取玩家形象
	{"teamDismiss", 59},           // 队伍解散

	{"rankGetTopList", 71},
	{"rankSetScore", 72},
	{"rankCreate", 74},
	{"rankDeleteScore", 75}, // 玩家排行榜信息移除
	{"rankGetRank", 76},
}
