
return {
	[1] = {name='heartbeat',desc='心跳'},
	[2] = {name='echo',desc='测试用回射协议'},
	[3] = {name='announcement',desc='公告'},
	[4] = {name='login',desc='登陆'},
	[5] = {name='gameLogin',desc='登陆游戏'},
	[6] = {name='reconnect',desc='重连'},
	[7] = {name='kick',desc='踢人'},
	[8] = {name='serverList',desc='获取服务器列表'},
	[9] = {name='enterMap',desc='进入地图'},
	[10] = {name='startAoi',desc='开启Aoi'},
	[11] = {name='enterSee',desc='角色进入视野'},
	[12] = {name='leaveSee',desc='角色离开视野'},
	[13] = {name='move',desc='移动'},
	[14] = {name='updatePos',desc='位置更新'},
	[15] = {name='leaveMap',desc='离开地图回到主界面'},
	[17] = {name='serverTime',desc='获取服务器时间'},
	[18] = {name='gameMaster',desc='gm'},
	[101] = {name='createRole',desc='创角'},
	[201] = {name='attrSync',desc='玩家属性同步'},
	[202] = {name='battleAttrSync',desc='玩家战斗属性同步'},
	[203] = {name='baseSync',desc='玩家基础信息同步'},
	[204] = {name='buyFatigue',desc='购买体力'},
	[205] = {name='attrSet',desc='设置属性值，特定属性能设置'},
	[206] = {name='buyGold',desc='购买金币'},
	[301] = {name='characterSync',desc='角色同步'},
	[302] = {name='characterBreak',desc='角色突破'},
	[303] = {name='characterTeamPrefabSync',desc='角色队伍预设同步'},
	[305] = {name='characterLevelUp',desc='角色升级'},
	[322] = {name='characterTeamPrefabSet',desc='角色编队预设修改'},
	[327] = {name='characterGift',desc='角色赠送礼物'},
	[328] = {name='characterLike',desc='角色喜欢'},
	[329] = {name='characterSkillLevelUp',desc='角色技能升级'},
	[330] = {name='characterGeneLevelUp',desc='命座等级提升'},
	[403] = {name='assetSync',desc='资源同步'},
	[402] = {name='drawCardSync',desc='抽卡数据同步'},
	[405] = {name='drawCardDraw',desc='抽取卡库'},
	[404] = {name='drawCardPoolChoose',desc='抽卡卡库选择'},
	[406] = {name='drawCardHistory',desc='抽卡历史记录'},
	[501] = {name='backpackSync',desc='背包同步'},
	[502] = {name='backpackUse',desc='背包使用物品'},
	[503] = {name='backpackSell',desc='背包出售物品'},
	[504] = {name='backpackRem',desc='背包移除物品'},
	[505] = {name='inoutput',desc='物品转换'},
	[601] = {name='questSync',desc='任务同步'},
	[603] = {name='questComplete',desc='完成任务'},
	[701] = {name='equipSync',desc='装备同步'},
	[703] = {name='equipEquip',desc='装备装备'},
	[704] = {name='equipDemount',desc='装备卸下'},
	[708] = {name='equipStrengthen',desc='装备强化'},
	[713] = {name='equipDecompose',desc='装备分解'},
	[714] = {name='equipLock',desc='装备锁定'},
	[715] = {name='equipRefine',desc='装备精炼'},
	[801] = {name='teamSync',desc='队伍同步'},
	[802] = {name='teamPosSync',desc='队伍位置同步'},
	[803] = {name='teamStatusSync',desc='队伍状态同步'},
	[804] = {name='teamGetNearTeam',desc='获取周围队伍'},
	[805] = {name='teamGetNearPlayer',desc='获取附近玩家'},
	[811] = {name='teamCreate',desc='队伍创建'},
	[812] = {name='teamDismiss',desc='队伍解散'},
	[813] = {name='teamDismissNotify',desc='队伍解散 通知到各队员'},
	[814] = {name='teamChangeInfo',desc='更改队伍信息'},
	[815] = {name='teamChangeInfoNotify',desc='队伍信息变更 通知给所有人'},
	[816] = {name='teamHeaderChanged',desc='队长变更'},
	[817] = {name='teamHeaderChangedNotify',desc='队长变更 通知给新队长'},
	[821] = {name='teamInvitedApply',desc='邀请玩家'},
	[822] = {name='teamInvitedApplyNotify',desc='邀请玩家 通知到被邀请方'},
	[823] = {name='teamInvitedReply',desc='邀请玩家 被邀请方处理'},
	[824] = {name='teamInvitedReplyNotify',desc='邀请玩家 被邀请方处理结果通知到邀请方'},
	[831] = {name='teamJoinApply',desc='申请加入队伍'},
	[832] = {name='teamJoinApplyNotify',desc='申请加入 通知到队长'},
	[833] = {name='teamJoinReply',desc='队长处理的入队请求'},
	[834] = {name='teamJoinReplyNotify',desc='队长处理的入队请求 通知到申请人'},
	[835] = {name='teamPlayerJoinNotify',desc='加入队伍通知 通知给除自己以外的所有人'},
	[841] = {name='teamPlayerLeave',desc='离开队伍'},
	[842] = {name='teamPlayerLeaveNotify',desc='离开队伍通知  通知给除自己以外的所有人'},
	[851] = {name='teamKickPlayer',desc='队伍踢人'},
	[852] = {name='teamKickPlayerNotify',desc='踢人 通告被踢方'},
	[861] = {name='teamBattle',desc='发起战斗'},
	[862] = {name='teamPlayerVote',desc='队员投票'},
	[863] = {name='teamPlayerVoteNotify',desc='队员投票通知'},
	[901] = {name='mainDungeonsSync',desc='主线剧情副本同步'},
	[1001] = {name='rankGetTopList',desc='获取上榜玩家列表'},
	[1002] = {name='rankGetRank',desc='获取排行榜自己名次'},
	[1101] = {name='materialDungeonSync',desc='材料副本同步'},
	[1111] = {name='secretRandomLevel',desc='秘境关卡随机'},
	[1150] = {name='scarsIngrainSync',desc='战痕印刻信息同步'},
	[1151] = {name='scarsIngrainGetScoreAward',desc='战痕印刻获取分数奖励'},
	[1152] = {name='scarsIngrainTeamPrefabSet',desc='战痕印刻Boss预设队伍'},
	[1153] = {name='scarsIngrainSaveScore',desc='战痕印刻保存分数记录'},
	[1210] = {name='trialDungeonSync',desc='试炼塔同步'},
	[1211] = {name='trialStageReward',desc='试炼塔阶段奖励'},
	[1212] = {name='trialDungeonReward',desc='试炼塔关卡奖励'},
	[1301] = {name='levelFight',desc='战斗开始'},
	[1302] = {name='levelFightEnd',desc='战斗结束'},
	[1401] = {name='rewardQuestSync',desc='悬赏任务同步'},
	[1402] = {name='rewardQuestAccept',desc='悬赏任务接取'},
	[1403] = {name='rewardQuestComplete',desc='悬赏任务完成'},
	[1404] = {name='rewardQuestRefresh',desc='悬赏任务刷新'},
	[1501] = {name='battleResurrect',desc='战斗复活'},
	[1601] = {name='fragmentChange',desc='碎片转化'},
	[1701] = {name='weaponSync',desc='武器同步'},
	[1703] = {name='weaponEquip',desc='武器装备'},
	[1708] = {name='weaponStrengthen',desc='武器强化'},
	[1709] = {name='weaponBreak',desc='武器突破'},
	[1710] = {name='weaponRefine',desc='武器精炼'},
	[1711] = {name='weaponLock',desc='武器锁定'},
	[1713] = {name='weaponDecompose',desc='武器分解'},
	[1801] = {name='baseSetSignature',desc='玩家签名'},
	[1802] = {name='baseSetBirthday',desc='玩家生日'},
	[1803] = {name='baseSetName',desc='玩家名称'},
	[1804] = {name='baseSetCharacterDisplay',desc='玩家角色展示'},
	[1805] = {name='baseSetPortraitFrameCard',desc='玩家的头像，头像框，名片'},
	[1806] = {name='baseSetSex',desc='玩家性别'},
	[1901] = {name='teachNone',desc='教学空包'},
	[2001] = {name='shopSync',desc='商店同步'},
	[2002] = {name='shopBuy',desc='商店购买'},
	[2003] = {name='shopRefresh',desc='商店刷新'},
	[2004] = {name='shopPay',desc='商城充值'},
	[2101] = {name='worldQuestSync',desc='世界任务同步'},
	[2102] = {name='worldQuestRefresh',desc='世界任务刷新'},
	[2104] = {name='reputationSync',desc='声望同步'},
	[2105] = {name='reputationItemConversion',desc='声望道具兑换'},
	[2106] = {name='reputationRefresh',desc='声望商店刷新'},
	[2201] = {name='chatMessageSend',desc='发送聊天消息'},
	[2202] = {name='chatMessageSync',desc='聊天消息同步'},
	[2301] = {name='queryRoleInfo',desc='查询玩家数据'},
	[2401] = {name='mailSync',desc='玩家邮件同步'},
	[2402] = {name='mailRead',desc='邮件设置已读'},
	[2403] = {name='mailDelete',desc='已读邮件删除'},
	[2501] = {name='signSync',desc='签到同步'},
	[2502] = {name='signIn',desc='签到'},
	[2601] = {name='bigSecretSync',desc='大秘境同步'},
	[2602] = {name='bigSecretCompetitionSync',desc='大秘境赛季同步'},
	[2603] = {name='bigSecretWeakness',desc='大秘境层级弱点'},
	[2604] = {name='bigSecretWeaknessRefresh',desc='大秘境层级弱点刷新'},
	[2605] = {name='bigSecretBlessingLvUp',desc='大秘境祝福升级'},
	[2606] = {name='bigSecretRandomLevel',desc='大秘境关卡随机'},
	[2701] = {name='talentSync',desc='天赋同步'},
	[2702] = {name='talentLevelUp',desc='天赋升级'},
	[2703] = {name='talentReset',desc='天赋重置'},
}
