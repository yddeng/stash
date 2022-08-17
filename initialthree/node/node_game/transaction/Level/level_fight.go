package Level

import (
	"github.com/golang/protobuf/proto"
	codecs "initialthree/codec/cs"
	"initialthree/node/common/attr"
	"initialthree/node/common/battleAttr"
	"initialthree/node/common/enumType"
	"initialthree/node/common/timeDisposal"
	"initialthree/node/common/transaction"
	"initialthree/node/node_game/battleAtt"
	scarsIngrain2 "initialthree/node/node_game/global/scarsIngrain"
	"initialthree/node/node_game/module/bigSecret"
	"initialthree/node/node_game/module/trialDungeon"
	"initialthree/node/node_game/module/worldQuest"
	"initialthree/node/table/excel/DataTable/BigSecretDungeon"
	"initialthree/node/table/excel/DataTable/SecretDungeon"
	"initialthree/node/table/excel/DataTable/TrialDungeon"
	"initialthree/node/table/excel/DataTable/WorldQuest"
	"initialthree/zaplogger"

	"initialthree/node/node_game/module"
	"initialthree/node/node_game/module/character"
	"initialthree/node/node_game/module/maindungeons"
	"initialthree/node/node_game/module/materialDungeon"
	"initialthree/node/node_game/module/scarsIngrain"
	"initialthree/node/node_game/temporary"
	"initialthree/node/node_game/user"
	MainChapter2 "initialthree/node/table/excel/ConstTable/MainChapter"
	DungeonTable "initialthree/node/table/excel/DataTable/Dungeon"
	"initialthree/node/table/excel/DataTable/MainChapter"
	"initialthree/node/table/excel/DataTable/MainDungeon"
	"initialthree/node/table/excel/DataTable/MaterialDungeon"
	"initialthree/node/table/excel/DataTable/Resurrect"
	"initialthree/node/table/excel/DataTable/ScarsIngrainBossInstance"
	"initialthree/protocol/cmdEnum"
	cs_message "initialthree/protocol/cs/message"
	"time"
)

type transactionLevelFight struct {
	transaction.TransactionBase
	user    *user.User
	req     *codecs.Message
	errcode cs_message.ErrCode
	resp    *cs_message.LevelFightToC
}

func (t *transactionLevelFight) Begin() {

	defer func() { t.EndTrans(t.resp, t.errcode) }()

	msg := t.req.GetData().(*cs_message.LevelFightToS)

	zaplogger.GetSugar().Infof("%s LevelFight %v", t.user.GetUserLogName(), msg)
	// 验证是否有未进行完的战斗
	tmpFightInfo := t.user.GetTemporary(temporary.TempLevelFight)
	// todo  单机暂不检查
	//if tmpFightInfo != nil {
	//	log.GetLogger().Debugf("%s LevelFight failed, tmpFightInfo %v != nil", t.user.GetUserLogName(), tmpFightInfo)
	//	t.errcode = cs_message.ErrCode_LevelFight_FightExist
	//	t.EndTrans()
	//	return
	//}

	isMainChapterActivityOpen := false
	dungeonCfg := DungeonTable.GetID(msg.GetDungeonID())
	if dungeonCfg == nil {
		zaplogger.GetSugar().Infof("%s LevelFight failed,  dungeonCfg %d is nil", t.user.GetUserLogName(), msg.GetDungeonID())
		t.errcode = cs_message.ErrCode_Config_Error
		return
	}

	isFirstFight := false
	switch dungeonCfg.SystemTypeEnum {
	case enumType.DungeonSystemType_MainChapter:
		isMainChapterActivityOpen, isFirstFight, t.errcode = t.checkMainDungeon(msg)
	case enumType.DungeonSystemType_Material:
		isFirstFight, t.errcode = t.checkMaterial(msg)
	case enumType.DungeonSystemType_ScarsIngrain:
		t.errcode = t.checkScarsIngrain(msg)
	case enumType.DungeonSystemType_WorldQuest:
		t.errcode = t.checkWorldQuest(msg)
	case enumType.DungeonSystemType_BigSecret:
		t.errcode = t.checkBigSecret(msg)
	case enumType.DungeonSystemType_Trial:
		isFirstFight, t.errcode = t.checkTrial(msg)
	case enumType.DungeonSystemType_Secret:
		t.errcode = t.checkSecret(msg)
	default:
		zaplogger.GetSugar().Infof("%s LevelFight failed,  SystemType %s invalid", t.user.GetUserLogName(), dungeonCfg.SystemType)
		t.errcode = cs_message.ErrCode_Request_Argument_Err
	}

	if t.errcode == cs_message.ErrCode_OK {

		if t.errcode = t.checkDungeon(msg, isMainChapterActivityOpen, isFirstFight); t.errcode != cs_message.ErrCode_OK {
			return
		}

		// 打包战斗属性
		teamList := msg.GetCharacterTeam().GetCharacterList()
		battleAttrSlice := make([]*cs_message.BattleAttrSlice, 0, len(teamList))
		for _, id := range teamList {
			sliceData := &cs_message.BattleAttrSlice{
				BattleAttrs: make([]*cs_message.BattleAttr, 0, battleAttr.AttrMax),
			}
			if id != 0 {
				bAttrs := battleAtt.RecalculateBattleAttr(t.user.PackBattleCharacter(id))
				for id := int32(1); id <= battleAttr.AttrMax; id++ {
					val := float64(0)
					if v, ok := bAttrs[id]; ok {
						val = v
					}
					sliceData.BattleAttrs = append(sliceData.BattleAttrs, &cs_message.BattleAttr{
						Id:  proto.Int32(id),
						Val: proto.Float64(val),
					})
				}
			}
			battleAttrSlice = append(battleAttrSlice, sliceData)
		}
		t.resp = &cs_message.LevelFightToC{
			CharacterTeam:   msg.GetCharacterTeam(),
			BattleAttrSlice: battleAttrSlice,
			StartTime:       proto.Int64(time.Now().Unix()),
		}

		tmpFightInfo = temporary.NewLevelFightInfo(t.user, msg, t.resp)
		fightInfo := tmpFightInfo.(*temporary.LevelFightInfo)
		t.resp.FightID = proto.Int64(fightInfo.FightID)
		t.user.SetTemporary(temporary.TempLevelFight, tmpFightInfo)

		zaplogger.GetSugar().Infof("%s LevelFight ok", t.user.GetUserLogName())
	}
}

// 前置检查, 仅检查表配置是否存在，上阵队伍是否符合要求
func (t *transactionLevelFight) checkDungeon(msg *cs_message.LevelFightToS, isMainChapterActivityOpen bool, isFirst bool) cs_message.ErrCode {
	levelID := msg.GetDungeonID()
	team := msg.GetCharacterTeam().GetCharacterList()
	if levelID == 0 || len(team) != 3 {
		zaplogger.GetSugar().Debugf("%s LevelFight failed, request argument err", t.user.GetUserLogName())
		return cs_message.ErrCode_Request_Argument_Err
	}

	dungeonCfg := DungeonTable.GetID(levelID)
	if dungeonCfg == nil {
		zaplogger.GetSugar().Debugf("%s LevelFight failed, user(%s, %d) dungeon %d not exist", t.user.GetUserID(), t.user.GetID(), levelID)
		return cs_message.ErrCode_Level_DungeonCfgNotFound
	}

	// 复活配置
	if dungeonCfg.ResurrectID != 0 {
		resurrectDef := Resurrect.GetID(dungeonCfg.ResurrectID)
		if resurrectDef == nil {
			zaplogger.GetSugar().Infof("%s LevelFight failed,  resurrectDef %d is nil ", t.user.GetUserLogName(), dungeonCfg.ResurrectID)
			return cs_message.ErrCode_Config_NotExist
		}
	}

	// 队伍检查
	userCharacter := t.user.GetSubModule(module.Character).(*character.UserCharacter)
	ids := map[int32]struct{}{}
	for _, id := range team {
		if id != 0 {
			if _, ok := ids[id]; ok {
				zaplogger.GetSugar().Debugf("%s LevelFight failed,character %d repeated set", t.user.GetUserLogName(), id)
				return cs_message.ErrCode_Character_TeamRepeated
			}
			ids[id] = struct{}{}

			c := userCharacter.GetCharacter(id)
			if c == nil {
				zaplogger.GetSugar().Debugf("%s LevelFight failed, character%d not exist", t.user.GetUserLogName(), id)
				return cs_message.ErrCode_Character_NotExist
			}
		}
	}

	// 试玩关
	isTrans := len(dungeonCfg.DemoCharactersArray) != 0
	if !isTrans {
		if len(ids) == 0 {
			zaplogger.GetSugar().Debugf("%s LevelFight failed, characterCount %d ", t.user.GetUserLogName(), len(ids))
			return cs_message.ErrCode_Level_CharacterCountError
		}
		switch dungeonCfg.TeamLimitTypeEnum {
		case enumType.DungeonTeamLimitType_CharacterCountLimit:
			if int32(len(ids)) > dungeonCfg.TeamLimitValue {
				zaplogger.GetSugar().Debugf("%s LevelFight failed, character count limit with dungeon %d", t.user.GetUserLogName(), levelID)
				return cs_message.ErrCode_Level_CharacterCountError
			}

		case enumType.DungeonTeamLimitType_AppointCharacter:
			find := false
			for id := range ids {
				if id == dungeonCfg.TeamLimitValue {
					find = true
					break
				}
			}
			if !find {
				zaplogger.GetSugar().Debugf("%s LevelFight failed, character not match with dungeon %d", t.user.GetUserLogName(), levelID)
				return cs_message.ErrCode_Level_NoSpecifiedCharacter
			}
		}
	}

	// 解锁条件
	needCheckOtherThanLevel := true
	if isMainChapterActivityOpen {
		chapterCfg := MainChapter.GetID(dungeonCfg.SystemConfig)
		if chapterCfg == nil {
			zaplogger.GetSugar().Errorf("MainDungeons Fight:  dungeon %d, chapter config not found.", t.user.GetUserLogName(), dungeonCfg.SystemConfig)
			return cs_message.ErrCode_Config_NotExist
		}
		needCheckOtherThanLevel = chapterCfg.DungeonsArray[0].ID != dungeonCfg.ID
	}
	unlocks := DungeonTable.GetUnlock(msg.GetDungeonID())
	for _, unlock := range unlocks {
		if unlock.Type == enumType.DungeonUnlockType_PlayerLevel {
			levelLimit := unlock.Args[0]
			if t.user.GetLevel() < levelLimit {
				zaplogger.GetSugar().Debugf("%s LevelFight failed, %d player level too low ", t.user.GetUserLogName(), levelLimit)
				return cs_message.ErrCode_Level_DungeonLock
			}
		} else if needCheckOtherThanLevel {
			switch unlock.Type {
			case enumType.DungeonUnlockType_MainChapter:
				userMainDungeon := t.user.GetSubModule(module.MainDungeons).(*maindungeons.MainDungeons)
				for _, id := range unlock.Args {
					if !userMainDungeon.IsDungeonPass(id) {
						zaplogger.GetSugar().Debugf("%s LevelFight failed,  mainDungeon %d is unPass ", t.user.GetUserLogName(), id)
						return cs_message.ErrCode_Level_DungeonLock
					}
				}
			case enumType.DungeonUnlockType_MaterialDungeon:
				userMaterialDungeon := t.user.GetSubModule(module.MaterialDungeon).(*materialDungeon.MaterialDungeon)
				for _, id := range unlock.Args {
					if userMaterialDungeon.GetMaterialDungeon(id) == nil {
						zaplogger.GetSugar().Debugf("%s LevelFight failed, MaterialDungeon %d is unPass ", t.user.GetUserLogName(), id)
						return cs_message.ErrCode_Level_DungeonLock
					}
				}
			}
		}
	}

	// 消耗

	costTypeEnum := int32(0)
	costValue := int32(0)
	if isFirst {
		costTypeEnum = dungeonCfg.FirstCostTypeEnum
		costValue = dungeonCfg.FirstCostArgs
	} else {
		costTypeEnum = dungeonCfg.CommonCostTypeEnum
		costValue = dungeonCfg.CommonCostArgs
	}

	if costTypeEnum != 0 && costValue != 0 {
		if dungeonCfg.SystemTypeEnum == enumType.DungeonSystemType_Material {
			if msg.GetMultiple() <= 0 {
				zaplogger.GetSugar().Debugf("%s LevelFight failed, material multiple %d is error ", t.user.GetUserLogName(), msg.GetMultiple())
				return cs_message.ErrCode_Request_Argument_Err
			}
			costValue = costValue * msg.GetMultiple()
		}

		switch costTypeEnum {
		case enumType.DungeonCostType_FatigueCost:
			if t.user.GetAttr(attr.CurrentFatigue) < int64(costValue) {
				zaplogger.GetSugar().Debugf("%s LevelFight failed, CurrentFatigue is not enough", t.user.GetUserLogName())
				return cs_message.ErrCode_Attr_Low
			}
		}
	}

	return cs_message.ErrCode_OK
}

// 主线副本
func (t *transactionLevelFight) checkMainDungeon(msg *cs_message.LevelFightToS) (bool, bool, cs_message.ErrCode) {
	dungeonCfg := DungeonTable.GetID(msg.GetDungeonID())
	def := MainDungeon.GetID(dungeonCfg.SystemConfig)
	if def == nil {
		zaplogger.GetSugar().Infof("%s LevelFight failed,  levelID %d MainDungeon is nil ", t.user.GetUserLogName(), dungeonCfg.SystemConfig)
		return false, false, cs_message.ErrCode_Config_NotExist
	}

	chapterCfg := MainChapter.GetID(def.ChapterID)
	if chapterCfg == nil {
		zaplogger.GetSugar().Errorf("MainDungeons Fight: user(%s, %d) dungeon %d, chapter config not found.", t.user.GetUserID(), t.user.GetID(), dungeonCfg.SystemConfig)
		return false, false, cs_message.ErrCode_Config_NotExist
	}

	if chapterCfg.PlayerLevelLimit > t.user.GetLevel() {
		zaplogger.GetSugar().Errorf("MainDungeons Fight: user(%s, %d) dungeon %d, level low.", t.user.GetUserID(), t.user.GetID(), dungeonCfg.SystemConfig)
		return false, false, cs_message.ErrCode_MainDungeons_DungeonNotOpen
	}

	userMainDungeons := t.user.GetSubModule(module.MainDungeons).(*maindungeons.MainDungeons)
	chapterConst := MainChapter2.Get()
	isActivityChapterOpen := chapterConst.GetActivityChapterID() == def.ChapterID && chapterConst.IsActivityChapterOpen(timeDisposal.Now())

	// 不是活动章节或活动章节未开放
	if !isActivityChapterOpen {
		for _, v := range chapterCfg.UnlockByMainDungeonsArray {
			if !userMainDungeons.IsDungeonPass(v.DungeonID) {
				zaplogger.GetSugar().Errorf("MainDungeons Fight: user(%s, %d) dungeon %d, chapter not open.", t.user.GetUserID(), t.user.GetID(), dungeonCfg.SystemConfig)
				return false, false, cs_message.ErrCode_MainDungeons_DungeonNotOpen
			}
		}
	}

	return isActivityChapterOpen, userMainDungeons.IsDungeonPass(dungeonCfg.SystemConfig) == false, cs_message.ErrCode_OK
}

// 材料副本战斗检测
func (t *transactionLevelFight) checkMaterial(msg *cs_message.LevelFightToS) (bool, cs_message.ErrCode) {
	dungeonCfg := DungeonTable.GetID(msg.GetDungeonID())

	def := MaterialDungeon.GetID(dungeonCfg.SystemConfig)
	if def == nil {
		zaplogger.GetSugar().Infof("%s LevelFight failed,  levelID %d MaterialDungeon is nil ", t.user.GetUserLogName(), dungeonCfg.SystemConfig)
		return false, cs_message.ErrCode_Config_NotExist
	}

	if !def.DungeonOpen() {
		zaplogger.GetSugar().Infof("%s LevelFight failed,  levelID %d MaterialDungeon is not open", t.user.GetUserLogName(), dungeonCfg.SystemConfig)
		return false, cs_message.ErrCode_Level_DungeonLock
	}

	userMaterialLevel := t.user.GetSubModule(module.MaterialDungeon).(*materialDungeon.MaterialDungeon)
	return userMaterialLevel.GetMaterialDungeon(dungeonCfg.SystemConfig) == nil, cs_message.ErrCode_OK
}

// 世界任务副本战斗检测
func (t *transactionLevelFight) checkWorldQuest(msg *cs_message.LevelFightToS) cs_message.ErrCode {
	dungeonCfg := DungeonTable.GetID(msg.GetDungeonID())

	def := WorldQuest.GetID(dungeonCfg.SystemConfig)
	if def == nil {
		zaplogger.GetSugar().Infof("%s LevelFight failed,  levelID %d WorldQuestDungeon is nil ", t.user.GetUserLogName(), dungeonCfg.SystemConfig)
		return cs_message.ErrCode_Config_NotExist
	}

	userWorldQuest := t.user.GetSubModule(module.WorldQuest).(*worldQuest.WorldQuest)
	if !userWorldQuest.CanDo() {
		return cs_message.ErrCode_Attr_Low
	}

	return cs_message.ErrCode_OK
}

// 秘境副本检测
func (t *transactionLevelFight) checkSecret(msg *cs_message.LevelFightToS) cs_message.ErrCode {
	dungeonCfg := DungeonTable.GetID(msg.GetDungeonID())

	def := SecretDungeon.GetID(dungeonCfg.SystemConfig)
	if def == nil {
		zaplogger.GetSugar().Infof("%s LevelFight failed,  levelID %d SecretDungeon is nil ", t.user.GetUserLogName(), dungeonCfg.SystemConfig)
		return cs_message.ErrCode_Config_NotExist
	}

	return cs_message.ErrCode_OK
}

// 材料副本战斗检测
func (t *transactionLevelFight) checkTrial(msg *cs_message.LevelFightToS) (bool, cs_message.ErrCode) {
	dungeonCfg := DungeonTable.GetID(msg.GetDungeonID())

	def := TrialDungeon.GetID(dungeonCfg.SystemConfig)
	if def == nil {
		zaplogger.GetSugar().Infof("%s LevelFight failed,  levelID %d TrialDungeon is nil ", t.user.GetUserLogName(), dungeonCfg.SystemConfig)
		return false, cs_message.ErrCode_Config_NotExist
	}

	userTrialDungeon := t.user.GetSubModule(module.TrialDungeon).(*trialDungeon.TrialDungeon)
	return userTrialDungeon.GetTrialDungeon(dungeonCfg.SystemConfig) == nil, cs_message.ErrCode_OK

}

// 大秘境
func (t *transactionLevelFight) checkBigSecret(msg *cs_message.LevelFightToS) cs_message.ErrCode {
	dungeonCfg := DungeonTable.GetID(msg.GetDungeonID())

	def := BigSecretDungeon.GetID(dungeonCfg.SystemConfig)
	if def == nil {
		zaplogger.GetSugar().Infof("%s LevelFight failed, levelID %d BigSecret is nil ", t.user.GetUserLogName(), dungeonCfg.SystemConfig)
		return cs_message.ErrCode_Config_NotExist
	}

	m := t.user.GetSubModule(module.BigSecret).(*bigSecret.BigSecretDungeon)
	if !m.Unlocked(dungeonCfg.SystemConfig) {
		zaplogger.GetSugar().Infof("%s LevelFight failed, level %d is locked", t.user.GetUserLogName(), dungeonCfg.SystemConfig)
		return cs_message.ErrCode_BigSecret_Locked
	}

	keyCount := m.GetKeyCount()
	if keyCount == 0 {
		zaplogger.GetSugar().Infof("%s LevelFight failed, key is not enough", t.user.GetUserLogName())
		return cs_message.ErrCode_BigSecret_KeyNotEnough
	}
	return cs_message.ErrCode_OK

}

// 战痕印刻战斗检测
func (t *transactionLevelFight) checkScarsIngrain(msg *cs_message.LevelFightToS) cs_message.ErrCode {
	dungeonCfg := DungeonTable.GetID(msg.GetDungeonID())

	bIns := ScarsIngrainBossInstance.GetID(dungeonCfg.SystemConfig)
	if bIns == nil {
		zaplogger.GetSugar().Infof("%s LevelFight failed,  levelID %d bossInstance is nil ", t.user.GetUserLogName(), dungeonCfg.SystemConfig)
		return cs_message.ErrCode_Config_NotExist
	}

	siModule := t.user.GetSubModule(module.ScarsIngrain).(*scarsIngrain.ScarsIngrain)
	siData := siModule.GetData()
	if siData.SIID == 0 {
		zaplogger.GetSugar().Debugf("%s LevelFight failed, SIID is 0", t.user.GetUserLogName())
		return cs_message.ErrCode_ERROR
	}
	gSiClass := scarsIngrain2.GetIDClass(siData.SIID)
	if !gSiClass.IsRunning(0) {
		zaplogger.GetSugar().Debugf("%s LevelFight failed, SIID %d is not run", t.user.GetUserLogName(), siData.SIID)
		return cs_message.ErrCode_ERROR
	}

	bossId, difficult := bIns.BossDifficult()
	// 挑战的 boss
	if unlock := siModule.BossDifficultUnlock(bossId, difficult); !unlock {
		zaplogger.GetSugar().Debugf("%s LevelFight failed, fight (bossID %d, difficult %d ) is locked ", t.user.GetUserLogName(), bossId, difficult)
		return cs_message.ErrCode_ScarsIngrain_BossUnlock
	}

	// 驻守角色
	defend := siModule.GetDefendCharacter(bossId)
	for _, id := range msg.GetCharacterTeam().GetCharacterList() {
		if _, ok := defend[id]; ok {
			zaplogger.GetSugar().Debugf("%s LevelFight failed, character %s is defend", t.user.GetUserLogName(), id)
			return cs_message.ErrCode_Request_Argument_Err
		}
	}

	// 清理该战场的零时数据
	siModule.ClearFightData()

	return cs_message.ErrCode_OK
}

func (t *transactionLevelFight) GetModuleName() string {
	return "Level"
}

func init() {
	user.RegisterTransFunc(cmdEnum.CS_LevelFight, func(user *user.User, msg *codecs.Message) transaction.Transaction {
		return &transactionLevelFight{
			user: user,
			req:  msg,
		}
	})
}
