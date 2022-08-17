package proto_def

type st struct {
	Name      string
	Desc      string
	MessageID int
}

var CS_message = []st{

	{"heartbeat", "心跳", 1},
	{"echo", "测试用回射协议", 2},
	{"announcement", "公告", 3},
	{"login", "登陆", 4},
	{"gameLogin", "登陆游戏", 5},
	{"reconnect", "重连", 6},
	{"kick", "踢人", 7},
	{"serverList", "获取服务器列表", 8},
	{"enterMap", "进入地图", 9},
	{"startAoi", "开启Aoi", 10},
	{"enterSee", "角色进入视野", 11},
	{"leaveSee", "角色离开视野", 12},
	{"move", "移动", 13},
	{"updatePos", "位置更新", 14},
	{"leaveMap", "离开地图回到主界面", 15},
	{"serverTime", "获取服务器时间", 17},
	{"gameMaster", "gm", 18},

	{"createRole", "创角", 101},

	{"attrSync", "玩家属性同步", 201},
	{"battleAttrSync", "玩家战斗属性同步", 202},
	{"baseSync", "玩家基础信息同步", 203},
	{"buyFatigue", "购买体力", 204},
	{"attrSet", "设置属性值，特定属性能设置", 205},
	{"buyGold", "购买金币", 206},

	{"characterSync", "角色同步", 301},
	{"characterBreak", "角色突破", 302},
	{"characterTeamPrefabSync", "角色队伍预设同步", 303},
	{"characterLevelUp", "角色升级", 305},
	{"characterTeamPrefabSet", "角色编队预设修改", 322},
	{"characterGift", "角色赠送礼物", 327},
	{"characterLike", "角色喜欢", 328},
	{"characterSkillLevelUp", "角色技能升级", 329},
	{"characterGeneLevelUp", "命座等级提升", 330},

	{"assetSync", "资源同步", 403},
	{"drawCardSync", "抽卡数据同步", 402},
	{"drawCardDraw", "抽取卡库", 405},
	{"drawCardPoolChoose", "抽卡卡库选择", 404},
	{"drawCardHistory", "抽卡历史记录", 406},

	{"backpackSync", "背包同步", 501},
	{"backpackUse", "背包使用物品", 502},
	{"backpackSell", "背包出售物品", 503},
	{"backpackRem", "背包移除物品", 504},
	{"inoutput", "物品转换", 505},

	{"questSync", "任务同步", 601},
	{"questComplete", "完成任务", 603},

	{"equipSync", "装备同步", 701},
	{"equipEquip", "装备装备", 703},
	{"equipDemount", "装备卸下", 704},
	{"equipStrengthen", "装备强化", 708},
	{"equipDecompose", "装备分解", 713},
	{"equipLock", "装备锁定", 714},
	{"equipRefine", "装备精炼", 715},

	{"teamSync", "队伍同步", 801},
	{"teamPosSync", "队伍位置同步", 802},
	{"teamStatusSync", "队伍状态同步", 803},
	{"teamGetNearTeam", "获取周围队伍", 804},
	{"teamGetNearPlayer", "获取附近玩家", 805},
	{"teamCreate", "队伍创建", 811},
	{"teamDismiss", "队伍解散", 812},
	{"teamDismissNotify", "队伍解散 通知到各队员", 813},
	{"teamChangeInfo", "更改队伍信息", 814},
	{"teamChangeInfoNotify", "队伍信息变更 通知给所有人", 815},
	{"teamHeaderChanged", "队长变更", 816},
	{"teamHeaderChangedNotify", "队长变更 通知给新队长", 817},
	{"teamInvitedApply", "邀请玩家", 821},
	{"teamInvitedApplyNotify", "邀请玩家 通知到被邀请方", 822},
	{"teamInvitedReply", "邀请玩家 被邀请方处理", 823},
	{"teamInvitedReplyNotify", "邀请玩家 被邀请方处理结果通知到邀请方", 824},
	{"teamJoinApply", "申请加入队伍", 831},
	{"teamJoinApplyNotify", "申请加入 通知到队长", 832},
	{"teamJoinReply", "队长处理的入队请求", 833},
	{"teamJoinReplyNotify", "队长处理的入队请求 通知到申请人", 834},
	{"teamPlayerJoinNotify", "加入队伍通知 通知给除自己以外的所有人", 835},
	{"teamPlayerLeave", "离开队伍", 841},
	{"teamPlayerLeaveNotify", "离开队伍通知  通知给除自己以外的所有人", 842},
	{"teamKickPlayer", "队伍踢人", 851},
	{"teamKickPlayerNotify", "踢人 通告被踢方", 852},
	{"teamBattle", "发起战斗", 861},
	{"teamPlayerVote", "队员投票", 862},
	{"teamPlayerVoteNotify", "队员投票通知", 863},

	{"mainDungeonsSync", "主线剧情副本同步", 901},

	{"rankGetTopList", "获取上榜玩家列表", 1001},
	{"rankGetRank", "获取排行榜自己名次", 1002},

	{"materialDungeonSync", "材料副本同步", 1101},

	{"secretRandomLevel", "秘境关卡随机", 1111},

	{"scarsIngrainSync", "战痕印刻信息同步", 1150},
	{"scarsIngrainGetScoreAward", "战痕印刻获取分数奖励", 1151},
	{"scarsIngrainTeamPrefabSet", "战痕印刻Boss预设队伍", 1152},
	{"scarsIngrainSaveScore", "战痕印刻保存分数记录", 1153},

	{"trialDungeonSync", "试炼塔同步", 1210},
	{"trialStageReward", "试炼塔阶段奖励", 1211},
	{"trialDungeonReward", "试炼塔关卡奖励", 1212},

	{"levelFight", "战斗开始", 1301},
	{"levelFightEnd", "战斗结束", 1302},

	{"rewardQuestSync", "悬赏任务同步", 1401},
	{"rewardQuestAccept", "悬赏任务接取", 1402},
	{"rewardQuestComplete", "悬赏任务完成", 1403},
	{"rewardQuestRefresh", "悬赏任务刷新", 1404},

	{"battleResurrect", "战斗复活", 1501},

	{"fragmentChange", "碎片转化", 1601},

	{"weaponSync", "武器同步", 1701},
	{"weaponEquip", "武器装备", 1703},
	{"weaponStrengthen", "武器强化", 1708},
	{"weaponBreak", "武器突破", 1709},
	{"weaponRefine", "武器精炼", 1710},
	{"weaponLock", "武器锁定", 1711},
	{"weaponDecompose", "武器分解", 1713},

	{"baseSetSignature", "玩家签名", 1801},
	{"baseSetBirthday", "玩家生日", 1802},
	{"baseSetName", "玩家名称", 1803},
	{"baseSetCharacterDisplay", "玩家角色展示", 1804},
	{"baseSetPortraitFrameCard", "玩家的头像，头像框，名片", 1805},
	{"baseSetSex", "玩家性别", 1806},

	{"teachNone", "教学空包", 1901},

	{"shopSync", "商店同步", 2001},
	{"shopBuy", "商店购买", 2002},
	{"shopRefresh", "商店刷新", 2003},
	{"shopPay", "商城充值", 2004},

	{"worldQuestSync", "世界任务同步", 2101},
	{"worldQuestRefresh", "世界任务刷新", 2102},
	{"reputationSync", "声望同步", 2104},
	{"reputationItemConversion", "声望道具兑换", 2105},
	{"reputationRefresh", "声望商店刷新", 2106},

	{"chatMessageSend", "发送聊天消息", 2201},
	{"chatMessageSync", "聊天消息同步", 2202},

	{"queryRoleInfo", "查询玩家数据", 2301},

	{"mailSync", "玩家邮件同步", 2401},
	{"mailRead", "邮件设置已读", 2402},
	{"mailDelete", "已读邮件删除", 2403},

	{"signSync", "签到同步", 2501},
	{"signIn", "签到", 2502},

	{"bigSecretSync", "大秘境同步", 2601},
	{"bigSecretCompetitionSync", "大秘境赛季同步", 2602},
	{"bigSecretWeakness", "大秘境层级弱点", 2603},
	{"bigSecretWeaknessRefresh", "大秘境层级弱点刷新", 2604},
	{"bigSecretBlessingLvUp", "大秘境祝福升级", 2605},
	{"bigSecretRandomLevel", "大秘境关卡随机", 2606},

	{"talentSync", "天赋同步", 2701},
	{"talentLevelUp", "天赋升级", 2702},
	{"talentReset", "天赋重置", 2703},
}
