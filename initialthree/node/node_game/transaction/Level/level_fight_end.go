package Level

import (
	"github.com/golang/protobuf/proto"
	codecs "initialthree/codec/cs"
	"initialthree/node/common/attr"
	"initialthree/node/common/droppool"
	"initialthree/node/common/enumType"
	"initialthree/node/common/inoutput"
	"initialthree/node/common/transaction"
	"initialthree/node/node_game/event"
	bigSecret2 "initialthree/node/node_game/global/bigSecret"
	"initialthree/node/node_game/module"
	"initialthree/node/node_game/module/base"
	"initialthree/node/node_game/module/bigSecret"
	"initialthree/node/node_game/module/character"
	"initialthree/node/node_game/module/maindungeons"
	"initialthree/node/node_game/module/materialDungeon"
	"initialthree/node/node_game/module/rankData"
	"initialthree/node/node_game/module/scarsIngrain"
	"initialthree/node/node_game/module/trialDungeon"
	"initialthree/node/node_game/module/worldQuest"
	"initialthree/node/node_game/temporary"
	"initialthree/node/node_game/user"
	"initialthree/node/table/excel/ConstTable/ScarsIngrain"
	"initialthree/node/table/excel/DataTable/BossKilledReward"
	DungeonTable "initialthree/node/table/excel/DataTable/Dungeon"
	"initialthree/node/table/excel/DataTable/ScarsIngrainBossChallenge"
	"initialthree/node/table/excel/DataTable/ScarsIngrainBossInstance"
	"initialthree/node/table/excel/DataTable/SecretDungeon"
	"initialthree/node/table/excel/DataTable/WorldQuest"
	"initialthree/protocol/cmdEnum"
	cs_message "initialthree/protocol/cs/message"
	"initialthree/zaplogger"
	"time"
)

type transactionLevelFightEnd struct {
	transaction.TransactionBase
	user    *user.User
	req     *codecs.Message
	errcode cs_message.ErrCode
	resp    *cs_message.LevelFightEndToC
}

func (t *transactionLevelFightEnd) Begin() {

	defer func() { t.EndTrans(t.resp, t.errcode) }()

	t.errcode = cs_message.ErrCode_OK
	msg := t.req.GetData().(*cs_message.LevelFightEndToS)
	zaplogger.GetSugar().Debugf("%s LevelFightEnd %v ", t.user.GetUserLogName(), msg)

	// 验证是否有未进行完的战斗
	fightInfo, ok := t.user.GetTemporary(temporary.TempLevelFight).(*temporary.LevelFightInfo)
	if !ok || fightInfo == nil {
		zaplogger.GetSugar().Debugf("%s LevelFightEnd failed, tmpFightInfo == nil", t.user.GetUserLogName())
		t.errcode = cs_message.ErrCode_Level_FightNotExist
		return
	}

	if fightInfo.FightID != msg.GetFightID() {
		zaplogger.GetSugar().Debugf("%s LevelFightEnd failed, fightId %d != %d", t.user.GetUserLogName(), fightInfo.FightID, msg.GetFightID())
		t.errcode = cs_message.ErrCode_Level_FightIDNotMatch
		return
	}

	t.resp = &cs_message.LevelFightEndToC{}
	t.resp.DungeonID = proto.Int32(fightInfo.Tos.GetDungeonID())
	t.resp.Pass = proto.Bool(msg.GetPass())

	isPass := msg.GetPass()
	t.user.ClearTemporary(temporary.TempLevelFight)
	if !isPass {
		return
	}

	t.dungeonPass(fightInfo, msg)
}

// 扣除消耗，结算奖励
func (t *transactionLevelFightEnd) dungeonPass(fightInfo *temporary.LevelFightInfo, msg *cs_message.LevelFightEndToS) {

	t.resp.UseTime = proto.Int32(msg.GetUseTime())
	t.resp.Stars = msg.GetStars()

	isFirstPass := false
	dungeonCfg := DungeonTable.GetID(fightInfo.Tos.GetDungeonID())
	switch dungeonCfg.SystemTypeEnum {
	case enumType.DungeonSystemType_MainChapter:
		isFirstPass = t.mainDungeonPass(fightInfo, msg)
	case enumType.DungeonSystemType_Material:
		isFirstPass = t.materialPass(fightInfo, msg)
	case enumType.DungeonSystemType_Secret:
		t.secretPass(fightInfo, msg)
	case enumType.DungeonSystemType_WorldQuest:
		t.worldQuestPass(fightInfo, msg)
	case enumType.DungeonSystemType_BigSecret:
		t.bigSecretPass(fightInfo, msg)
	case enumType.DungeonSystemType_Trial:
		isFirstPass = t.trialPass(fightInfo, msg)
	case enumType.DungeonSystemType_ScarsIngrain:
		// 战痕印刻结束，需要手动选择，这里仅返回分数
		t.scarsIngrainPass(fightInfo, msg)
		return
	}

	if t.errcode != cs_message.ErrCode_OK {
		return
	}

	// 扣除消耗
	tos := fightInfo.Tos
	costTypeEnum := dungeonCfg.CommonCostTypeEnum
	costValue := dungeonCfg.CommonCostArgs
	if isFirstPass {
		costTypeEnum = dungeonCfg.FirstCostTypeEnum
		costValue = dungeonCfg.FirstCostArgs
	}

	if costTypeEnum != 0 && costValue != 0 {
		switch dungeonCfg.SystemTypeEnum {
		case enumType.DungeonSystemType_Material:
			costValue = costValue * tos.GetMultiple()
		}

		switch costTypeEnum {
		case enumType.DungeonCostType_FatigueCost:
			t.user.DoInputOutput([]inoutput.ResDesc{{Type: enumType.IOType_UsualAttribute, ID: attr.CurrentFatigue, Count: costValue}}, nil, false)
		}
	}

	//玩家经验以及角色经验首通和通常的配置是独立的嘛，首通就读首通的，通常就读通常的
	//只有奖励是首通的时候是首通+通常嘛
	// 玩家经验
	playerExp := dungeonCfg.CommonPlayerExp
	if isFirstPass {
		playerExp = dungeonCfg.FirstPlayerExp
	}
	if playerExp != 0 {
		if dungeonCfg.SystemTypeEnum == enumType.DungeonSystemType_Material {
			playerExp = playerExp * tos.GetMultiple()
		}
		_ = t.user.DoInputOutput(nil, []inoutput.ResDesc{{Type: enumType.IOType_UsualAttribute, ID: attr.CurrentExp, Count: playerExp}}, true)

	}
	// 角色经验
	characterExp := dungeonCfg.CommonCharacterExp
	if isFirstPass {
		characterExp = dungeonCfg.FirstCharacterExp
	}
	if characterExp != 0 {
		if dungeonCfg.SystemTypeEnum == enumType.DungeonSystemType_Material {
			characterExp = characterExp * tos.GetMultiple()
		}
		userCharacter := t.user.GetSubModule(module.Character).(*character.UserCharacter)
		for _, id := range tos.GetCharacterTeam().GetCharacterList() {
			if id != 0 {
				chara := userCharacter.GetCharacter(id)
				if chara != nil {
					maxLevel := userCharacter.GetMaxLevel(chara)
					_, level, currentExp := userCharacter.CalcUseExp(chara, characterExp, maxLevel)
					userCharacter.SetLevel(chara, level, currentExp)
				}
			}
		}
	}
	// 掉落
	dropIDs := make([]int32, 0, 2)
	if isFirstPass && dungeonCfg.FirstDrop != 0 {
		dropIDs = append(dropIDs, dungeonCfg.FirstDrop)
	}
	if dungeonCfg.CommonDrop != 0 {
		dropIDs = append(dropIDs, dungeonCfg.CommonDrop)
		if dungeonCfg.SystemTypeEnum == enumType.DungeonSystemType_Material {
			for i := int32(1); i < tos.GetMultiple(); i++ {
				dropIDs = append(dropIDs, dungeonCfg.CommonDrop)
			}
		}
	}
	if len(dropIDs) > 0 {
		drop := droppool.DropAward(dropIDs...)
		if !drop.IsZero() {
			t.resp.AwardList = append(t.resp.AwardList, drop.ToMessageAward())
			t.user.ApplyDropAward(drop)
		}
	}

	// 事件触发
	t.user.EmitEvent(event.EventInstanceSucceed, tos.GetDungeonID(), int32(0), map[int32]int32{})
	if dungeonCfg.SystemTypeEnum == enumType.DungeonSystemType_Material {
		// 多倍挑战 任务进度应该算多次
		for i := int32(1); i < tos.GetMultiple(); i++ {
			t.user.EmitEvent(event.EventInstanceSucceed, tos.GetDungeonID(), int32(0), map[int32]int32{})
		}
	}
}

// 主线副本结束
func (t *transactionLevelFightEnd) mainDungeonPass(fightInfo *temporary.LevelFightInfo, msg *cs_message.LevelFightEndToS) (first bool) {
	fightTos := fightInfo.Tos
	dungeonID := fightTos.GetDungeonID()

	userMainDungeons := t.user.GetSubModule(module.MainDungeons).(*maindungeons.MainDungeons)
	dungeonCfg := DungeonTable.GetID(dungeonID)
	return userMainDungeons.DungeonPass(dungeonCfg.SystemConfig)
}

// 材料副本结束
func (t *transactionLevelFightEnd) materialPass(fightInfo *temporary.LevelFightInfo, msg *cs_message.LevelFightEndToS) (first bool) {
	tos := fightInfo.Tos
	dungeonCfg := DungeonTable.GetID(tos.GetDungeonID())

	userMaterialLevel := t.user.GetSubModule(module.MaterialDungeon).(*materialDungeon.MaterialDungeon)
	return userMaterialLevel.Pass(dungeonCfg.SystemConfig)
}

// 大秘境结束
func (t *transactionLevelFightEnd) bigSecretPass(fightInfo *temporary.LevelFightInfo, msg *cs_message.LevelFightEndToS) {
	tos := fightInfo.Tos
	dungeonCfg := DungeonTable.GetID(tos.GetDungeonID())

	useTime := msg.GetUseTime()
	userBigSecret := t.user.GetSubModule(module.BigSecret).(*bigSecret.BigSecretDungeon)
	userBigSecret.Pass(dungeonCfg.SystemConfig, useTime)

	nowUnix := time.Now().Unix()
	bsData := bigSecret2.GetData()
	if bsData != nil && nowUnix > bsData.BeginTime && nowUnix < bsData.EndTime {
		// 排行榜还在进行，设置分数
		// 分数 8位。 低5位表使用时间，支持27小时。高3位表示层级支持 999 层
		// 用时越少排名越靠前，这里时间取反。时间越少数字越大
		useTime = 99999 - useTime
		if useTime < 0 {
			useTime = 0
		}
		score := dungeonCfg.SystemConfig*100000 + useTime

		baseModule := t.user.GetSubModule(module.Base).(*base.UserBase)
		rankRoleInfo := &cs_message.RankRoleInfo{
			ID:            proto.Uint64(t.user.GetID()),
			Name:          proto.String(baseModule.GetName()),
			Level:         proto.Int32(t.user.GetLevel()),
			Score:         proto.Int32(score),
			CharacterList: tos.GetCharacterTeam().GetCharacterList(),
		}

		rankModule := t.user.GetSubModule(module.RankData).(*rankData.RankData)
		rankModule.SetRank(bsData.RankID, score, bsData.RankLogic, rankRoleInfo, func(rank, total int32, err error) {
			if err != nil {
				zaplogger.GetSugar().Debugf("%s BigSecret SetScore failed, %v ", t.user.GetUserLogName(), err)
				return
			}
			zaplogger.GetSugar().Debugf("%s BigSecret SetScore, rank (idx %d, total %d)", t.user.GetUserLogName(), rank, total)
		})
	}
}

// 试炼塔副本结束
func (t *transactionLevelFightEnd) trialPass(fightInfo *temporary.LevelFightInfo, msg *cs_message.LevelFightEndToS) (first bool) {
	tos := fightInfo.Tos
	dungeonCfg := DungeonTable.GetID(tos.GetDungeonID())

	userTrialDungeon := t.user.GetSubModule(module.TrialDungeon).(*trialDungeon.TrialDungeon)
	return userTrialDungeon.Pass(dungeonCfg.SystemConfig)
}

// 战痕印刻战斗结束
func (t *transactionLevelFightEnd) scarsIngrainPass(fightInfo *temporary.LevelFightInfo, msg *cs_message.LevelFightEndToS) {
	tos := fightInfo.Tos
	t.resp.ScarsIngrainEnd = &cs_message.ScarsIngrainEnd{}

	siModule := t.user.GetSubModule(module.ScarsIngrain).(*scarsIngrain.ScarsIngrain)
	parseVal := func(val, min, max int32) int32 {
		if val < min {
			return min
		}
		if val > max {
			return max
		}
		return val
	}

	// boss 死亡判断：boss剩余血量等于0
	bossDie := false

	dungeonCfg := DungeonTable.GetID(tos.GetDungeonID())

	/*
	 boss 1~4难度单关分数由该关卡通关评价等级对应的分数决定；未通关不计算分数
	 5难度单关分数由该关通关时间和造成伤害决定，具体各部分算法不变；未通关仅计算伤害得分
	 boss_C 有额外的分数加成，即idx=2。
	*/
	def := ScarsIngrainBossInstance.GetID(dungeonCfg.SystemConfig)
	bossId, difficult := def.BossDifficult()
	idx, bossGroup, historyScore := siModule.BossDifficultScore(bossId, difficult)
	// 历史分数
	if historyScore != -1 {
		t.resp.ScarsIngrainEnd.HistoryScore = proto.Int32(historyScore)
	}
	// 分数加成
	addition := float64(1)
	if idx == 2 {
		addition += ScarsIngrain.GetID(1).ScoreAddition
	}

	if difficult <= 4 {
		if msg.GetBossCurHP() > 0 {
			return
		}

		challengeDef := ScarsIngrainBossChallenge.GetID(bossGroup.GetChallengeID())

		score := int32(0)
		switch challengeDef.ChallengeTypeEnum {
		case enumType.ScarsIngrainBossChallengeType_Time:
			score, _ = challengeDef.GetScoreBuff(msg.GetUseTime())
		case enumType.ScarsIngrainBossChallengeType_BeHit:
			score, _ = challengeDef.GetScoreBuff(msg.GetBeHit())
		}

		score = int32(float64(score) * addition)
		scoreID := siModule.AddFightData(score, true, tos)
		t.resp.ScarsIngrainEnd.ScoreID = proto.Int32(scoreID)
		t.resp.ScarsIngrainEnd.TotalScore = proto.Int32(score)

	} else {

		// 分数计算
		damageMaxScore := int32(float64(def.AllScore) * def.DamageScorePercent)
		timeMaxScore := int32(float64(def.AllScore) * def.TimeScorePercent)

		// -- 造成伤害评分计算规则：造成伤害评分=造成伤害评分上限值*（1-Boss当前HP/Boss总HP），四舍五入取整
		damageScore := int32(0)
		if msg.GetBossCurHP() > 0 {
			damageScore = int32(float64(damageMaxScore) * (1 - msg.GetBossCurHP()/msg.GetBossMaxHP()))
			damageScore = parseVal(damageScore, 0, damageMaxScore)
		} else {
			bossDie = true
			damageScore = damageMaxScore
		}

		// -- 剩余时间评分计算规则：剩余时间评分 = 剩余时间评分上限值*（剩余时间/总时间）
		// 若未击杀Boss，则不入算剩余时间得分
		timeScore := int32(0)
		if bossDie {
			lastTime := dungeonCfg.ClearTimeLimit - msg.GetUseTime() // 剩余时间
			if dungeonCfg.ClearTimeLimit > 0 && lastTime > 0 {
				timeScore = int32(float64(timeMaxScore) * (float64(lastTime) / float64(dungeonCfg.ClearTimeLimit)))
				timeScore = parseVal(timeScore, 0, timeMaxScore)
			}
		}

		totalScore := timeScore + damageScore
		totalScore = int32(float64(totalScore) * addition)
		scoreID := siModule.AddFightData(totalScore, bossDie, tos)

		t.resp.ScarsIngrainEnd.ScoreID = proto.Int32(scoreID)
		t.resp.ScarsIngrainEnd.TotalScore = proto.Int32(totalScore)
		t.resp.ScarsIngrainEnd.TimeScore = proto.Int32(timeScore)
		t.resp.ScarsIngrainEnd.DamageScore = proto.Int32(damageScore)
	}
}

// 世界任务
func (t *transactionLevelFightEnd) worldQuestPass(fightInfo *temporary.LevelFightInfo, msg *cs_message.LevelFightEndToS) {
	tos := fightInfo.Tos

	dungeonCfg := DungeonTable.GetID(tos.GetDungeonID())
	def := WorldQuest.GetID(dungeonCfg.SystemConfig)

	userWorldQuest := t.user.GetSubModule(module.WorldQuest).(*worldQuest.WorldQuest)
	if ok := userWorldQuest.Pass(dungeonCfg.SystemConfig, def); ok {
		drop := droppool.DropAward(def.DroppoolID)
		if !drop.IsZero() {
			t.resp.AwardList = append(t.resp.AwardList, drop.ToMessageAward())
			t.user.ApplyDropAward(drop)
		}
	}

}

// 秘境副本结束
func (t *transactionLevelFightEnd) secretPass(fightInfo *temporary.LevelFightInfo, msg *cs_message.LevelFightEndToS) {
	tos := fightInfo.Tos
	dungeonCfg := DungeonTable.GetID(tos.GetDungeonID())
	sdef := SecretDungeon.GetID(dungeonCfg.SystemConfig)

	awards := BossKilledReward.GenDropAward(sdef.DifficultyTypeEnum, msg.GetKillBossCount())
	coinAward := BossKilledReward.GetSecretCoinAward(sdef.DifficultyTypeEnum)

	out := make([]inoutput.ResDesc, 0, len(awards)+1)
	for _, award := range awards {
		if !award.IsZero() {
			out = append(out, award.ToResDesc()...)
		}
	}
	for _, v := range coinAward.Infos {
		out = append(out, v.ToResDesc())
	}

	insOut := t.user.OutputIns(out)
	if t.errcode = t.user.DoInputOutput(nil, insOut, true); t.errcode != cs_message.ErrCode_OK {
		zaplogger.GetSugar().Debugf("%s error: inouput %s  ", t.user.GetUserLogName(), t.errcode.String())
		return
	}

	award := inoutput.OutputIns2Award(insOut)
	t.resp.AwardList = append(t.resp.AwardList, award)

}

func (t *transactionLevelFightEnd) GetModuleName() string {
	return "Level"
}

func init() {
	user.RegisterTransFunc(cmdEnum.CS_LevelFightEnd, func(user *user.User, msg *codecs.Message) transaction.Transaction {
		return &transactionLevelFightEnd{
			user: user,
			req:  msg,
		}
	})
}
