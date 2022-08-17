package excel

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
)

//	DefType_key       // celName@key@string //类型int32
//	DefType_mutilKey1 // name@mutilKey1@celName,celName //类型string
//	DefType_mutilKey2 // name@mutilKey2@celName:1000,celName:1 //类型int32 计算星*1000+y*1
//
//  自定义字段及类型  // celName@int64 //是表中有的字段名，类型以这里为准
//	DefType_enum   // celName@enum@[default enum string]
//	DefType_struct // celName@struct@name:type,name:type@;,
//	DefType_array1 // celName@array1@name:type,name:type@:, //对当前格子的数据分,只有一个分割符只填一个，两个符填两个
//	DefType_array2 // name@array2@name:type,name:type@celName,celName //将多个格子合并. 规定表起始为1且连续
//)

var (
	constTableDef = map[string][]string{}
	dataTableDef  = map[string][]string{}
)

func addConstTableSheet(sheetName string, defTypes ...string) {
	if _, ok := constTableDef[sheetName]; ok {
		panic(fmt.Sprintf("%s const 重复定义", sheetName))
	}
	constTableDef[sheetName] = defTypes
}

func addDataTableSheet(sheetName string, defTypes ...string) {
	if _, ok := dataTableDef[sheetName]; ok {
		panic(fmt.Sprintf("%s data 重复定义", sheetName))
	}
	dataTableDef[sheetName] = defTypes
}

func walkPath(loadPath, writePath string, def map[string][]string) {
	names := map[string]string{}
	readPath := path.Join(loadPath, writePath)
	if err := filepath.Walk(readPath, func(filePath string, f os.FileInfo, err error) error {
		if f != nil && !f.IsDir() {
			filename := f.Name()
			if strings.Contains(filename, ".xlsx") {
				if p, ok := names[filename]; ok {
					panic(fmt.Sprintf("%s %s 重复文件", filePath, p))
				}
				names[filename] = filePath
				dir := path.Dir(strings.TrimPrefix(filePath, loadPath))
				name := strings.TrimSuffix(filename, ".xlsx")

				if defs, ok := def[name]; ok {
					tab := newTable(name, dir, writePath, defs...)
					gen(tab, loadPath)
				}
			}
		}
		return nil
	}); err != nil {
		panic(err)
	}
}

func Gen(loadPath string) {
	walkPath(loadPath, "ConstTable", constTableDef)
	walkPath(loadPath, "DataTable", dataTableDef)
}

func GenCheck(loadPath, writePath string) []string {
	imports := []string{}
	tmp := `    _ "initialthree/node/table/excel/%s/%s"
`
	def := map[string][]string{}
	switch writePath {
	case "ConstTable":
		def = constTableDef
	case "DataTable":
		def = dataTableDef
	default:
		return imports
	}

	readPath := path.Join(loadPath, writePath)
	if err := filepath.Walk(readPath, func(filePath string, f os.FileInfo, err error) error {
		if f != nil && !f.IsDir() {
			filename := f.Name()
			if strings.Contains(filename, ".xlsx") {
				name := strings.TrimSuffix(filename, ".xlsx")
				if _, ok := def[name]; ok {
					imports = append(imports, fmt.Sprintf(tmp, writePath, name))
				}
			}
		}
		return nil
	}); err != nil {
		panic(err)
	}
	return imports
}

func init() {

	/********************** ConstTable *****************************************/
	addConstTableSheet(
		"Global",
		"DailyRefreshTime@struct@hour:int32,minute:int32@:",
		"LevepUpCostItemsAndItOfferedExp@array1@ItemID:int32,Exp:int32@,#",
		"FavorExp@array1@Exp:int32@,",
		"DrawCardTenthGuaranteeRarity@enum",
		"DrawCardGuaranteeRarity@enum",
		"DefaultCharacterTeam@array1@ID:int32@,",
	)

	addConstTableSheet("Quest",
		"InitMainStoryQuest@array1@ID:int32@,",
		"DailyRewardQuest@array1@ID:int32@,",
		"DailyQuest_1@array1@ID:int32@,",
		"DailyQuest_2@array1@ID:int32@,",
		"DailyQuest_3@array1@ID:int32@,",
		"DailyQuest_4@array1@ID:int32@,",
		"DailyQuest_5@array1@ID:int32@,",
		"DailyQuest_6@array1@ID:int32@,",
		"DailyQuest_7@array1@ID:int32@,",
		"WeeklyQuest@array1@ID:int32@,",
	)

	addConstTableSheet(
		"MainChapter",
		"SpecialUnlockStarDate@struct@year:int32,month:int32,day:int32,hour:int32,min:int32@.",
		"SpecialUnlockEndDate@struct@year:int32,month:int32,day:int32,hour:int32,min:int32@.",
	)

	addConstTableSheet("ScarsIngrain",
		"BeginTime@struct@weekly:int32,hour:int32,minute:int32@,",
		"EndTime@struct@weekly:int32,hour:int32,minute:int32@,",
		"BossOpenTime1@struct@weekly:int32,hour:int32,minute:int32@,",
		"BossOpenTime2@struct@weekly:int32,hour:int32,minute:int32@,",
		"BossOpenTime3@struct@weekly:int32,hour:int32,minute:int32@,",
	)

	addConstTableSheet(
		"Equip",
		"EquipSupplyExpItem@array1@ItemID:int32,Exp:int32@,#",
	)

	addConstTableSheet(
		"Weapon",
		"WeaponSupplyExpItem@array1@ItemID:int32,Exp:int32@,#",
	)

	addConstTableSheet(
		"Player",
	)

	addConstTableSheet(
		"PlayerCamp",
		"UnlockDungeon@array2@ID:int32@UnlockDungeon_",
		"ReputationItem@array2@ID:int32@ReputationItem_",
		"ReputationRefresh@array1@Cost:int32@,",
	)

	addConstTableSheet(
		"BigSecret",
		"PassTimeUnlock@array2@PassTime:int32,Unlock:int32@PassTimeLimit_,UnlockLv_",
		"Weakness@array1@ID:int32@,",
		"WeaknessRefreshCost@array1@Cost:int32@,",
	)

	addConstTableSheet(
		"Talent",
		"Attr@array2@Attr:int32,Value:int32@InfiniteTatlentAttr_,InfiniteTalentAttrValue_",
	)

	/********************** DataTable *****************************************/

	addDataTableSheet("PlayerLevel")

	addDataTableSheet("Map",
		"DefaultPosition@struct@x:int32,y:int32,z:int32@,")

	addDataTableSheet("PlayerCharacter",
		"Rarity@enum",
		"BreakIDList@array1@ID:int32@,",
		"GiftType@enum",
		"PlayeSkills@array1@ID:int32@,",
		"DrawCardItemIDs@array1@ID:int32@,",
		"DrawCardItemCounts@array1@Count:int32@,",
		"MaxGeneLvDrawCardItemIDs@array1@ID:int32@,",
		"MaxGeneLvDrawCardItemCounts@array1@Count:int32@,",
	)

	addDataTableSheet("PlayerSkill",
		"Skill@array2@LimitLevel:int32,Gold:int32,ItemStr:string@LevelUpCharacterLvLimit_,LevelUpCostGold_,LevelUpCostItem_",
	)

	addDataTableSheet("PlayerGene",
		"PlayerSkillID@string",
		"Attri@array1@ID:int32,Val:float64@,#",
	)

	addDataTableSheet("CharacterLevelUpAttribute",
		"AttributeList@array1@ID:int32,Val:float64@,#",
	)

	addDataTableSheet("CharacterLevelUpExp")

	addDataTableSheet("CharacterResource",
		"DamageElementType@enum@",
	)

	addDataTableSheet("CharacterBreakThrough",
		"SpecifiedIDList@array1@ID:int32,Count:int32@,#",
		"AttributeBonus@array1@ID:int32,Val:float64@,#",
	)

	addDataTableSheet("DemoCharacter")

	addDataTableSheet("Item",
		"Rarity@enum",
		"Type@enum",
		"SellCurrencyType@enum",
		"TimeLimitType@enum@",
	)

	addDataTableSheet("InputOutput",
		"Input@array2@Type:enum,ID:int32,Count:int32@InputType_,InputID_,InputCount_",
		"Output@array2@Type:enum,ID:int32,Count:int32@OutputType_,OutputID_,OutputCount_")

	addDataTableSheet("DropPool",
		"Type@enum",
		"Repeatable@bool",
		"DropList@array2@Type:enum,ID:int32,Count:int32,Wave:int32,Weight:int32@DropType_,DropID_,DropCount_,Wave_,DropWeight_",
	)

	addDataTableSheet(
		"DrawCardsLib",
		"DrawCardsPool@array1@PoolID:int32@,",
	)

	addDataTableSheet("DrawCardsPool",
		"TenTimesGuaranteePools@array1@Idx:int32@,",
		"DropList@array2@ID:int32,Weight:int32@CardsPoolID_,CardsPoolWeight_",
	)

	addDataTableSheet(
		"Dungeon",
		"TeamLimitType@enum@",
		"Unlocks@array2@Type:enum,Arg:string@UnlockType_,UnlockArgs_",
		"Rewards@array2@Type:enum,Arg:int32@RewardType_,RewardArgs_",
		"FirstCostType@enum@",
		"CommonCostType@enum@",
		"SystemType@enum",
		"DemoCharacters@array1@CharacterID:int32@,",
	)

	addDataTableSheet(
		"MainChapter",
		"ChapterType@enum",
		"Dungeons@array1@ID:int32@,",
		"UnlockByMainDungeons@array1@DungeonID:int32@,",
	)

	addDataTableSheet(
		"MainDungeon",
		"QuestId@int32",
	)

	addDataTableSheet("ScarsIngrainArea",
		"BossIDs@array1@ID:int32@,",
	)

	addDataTableSheet("ScarsIngrainBossChallenge",
		"ChallengeType@enum",
		"Challenge@array2@ArgLeft:int32,ArgRight:int32,Score:int32,Buff:int32@ArgLeft_,ArgRight_,Score_,Buff_",
	)

	addDataTableSheet("ScarsIngrainBoss",
		"ChallengeConfig@array1@ID:int32@,",
		"BossDifficulty@array2@InstanceID:int32,SkillStr:string,BuffStr:string@Instance_,skills_,buffs_",
	)

	addDataTableSheet("ScarsIngrainBossInstance")

	addDataTableSheet("ScarsIngrainRankReward",
		"RankReward@array2@Start:float64,End:float64,DropPoolID:int32@RankStartPerc_,RankEndPerc_,DropPoolID_",
	)

	addDataTableSheet("ScarsIngrainScoreReward",
		"ScoreReward@array2@Score:int32,DropPoolID:int32@Score_,DropPool_",
	)

	addDataTableSheet("FatigueSupply")

	addDataTableSheet("MaterialDungeon",
		"OpenTime@array1@Weekday:int32@,",
	)

	addDataTableSheet("Resurrect",
		"Resource@array2@ResourceType:string,ResourceID:int32,ResourceCount:int32@ResourceType_,ResourceID_,ResourceCount_",
	)

	addDataTableSheet("RewardQuest")

	addDataTableSheet("RewardQuestPosition")

	addDataTableSheet("FragmentChange",
		"Item@array2@ID:int32,Count:int32@ItemID,ItemCount",
	)

	addDataTableSheet("EquipOverclock",
		"Cost@array2@Type:enum,ID:int32,Count:int32@CostType,CostID,CostCount",
	)

	addDataTableSheet("AccountInitializeAssets",
		"AssetType@enum",
	)

	addDataTableSheet("Function",
		"Unlock@array2@Type:enum,Arg:string@UnlockType_,UnlockArg_",
	)

	addDataTableSheet("Gift",
		"GiftType@enum",
	)

	addDataTableSheet("Weapon",
		"RarityType@enum",
		"WeaponType@enum",
		"AttrConfigs@array2@ID:int32@AttribConfig_",
	)

	addDataTableSheet("WeaponAttribute",
		"LevelAttr@array2@Val:float64@LevelAttrib_",
		"BreakLevelAttr@array2@Val:float64@BreakLevelAttrib_",
	)

	addDataTableSheet("WeaponBreakThough",
		"BreakLevel@array2@Level:int32,Gold:int32,ItemStr:string,AttrStr:string@Level_,GoldCost_,ItemCost_,Attribute_",
	)

	addDataTableSheet("WeaponLevelMaxExp",
		"MaxExp@array2@Exp:int32@LevelMaxExp_",
	)

	addDataTableSheet("WeaponRarity")

	addDataTableSheet("WeaponDecompose",
		"Items@array2@Id:int32,Count:int32@ReturnItemID_,ReturnItemCount_",
	)

	addDataTableSheet("Equip",
		"Quality@enum",
		"AttrConfigs@array2@ID:int32@Attrib_",
	)

	addDataTableSheet("EquipAttribute",
		"Attr@array2@Val:float64@LevelAttrib_",
	)

	addDataTableSheet("EquipDecompose",
		"Items@array2@Id:int32,Count:int32@ReturnItemID_,ReturnItemCount_",
	)

	addDataTableSheet("EquipLevelMaxExp",
		"MaxExp@array2@Exp:int32@LevelMaxExp_",
	)

	addDataTableSheet("EquipQuality")

	addDataTableSheet("EquipSkill",
		"EquipSkillLimitType@enum@",
	)

	addDataTableSheet("EquipRandomAttributePool",
		"Random@array2@ID:int32,Weight:int32@ID_,Weight_",
	)

	addDataTableSheet("Skill",
		"Damage@array2@Desc:string@Desc_",
	)

	addDataTableSheet("PlayerCard")

	addDataTableSheet("PlayerPortrait")

	addDataTableSheet("PlayerPortraitFrame")

	addDataTableSheet("ProductLibrary",
		"RedreshPrice@array2@Price:int32@RefreshPrice_",
		"Products@array1@ID:int32@,",
	)

	addDataTableSheet("Product",
		"PType@enum",
		"ProductLimitType@enum",
	)
	addDataTableSheet("Pay")

	addDataTableSheet("WorldQuest",
		"CampType@enum",
	)

	addDataTableSheet("WorldQuestPool",
		"QuestList@array2@List:string@QuestList_",
		"MinMaxLevelID@mutilKey2@MinLevel:1000,MaxLevel:1",
	)

	addDataTableSheet("PlayerCampReputationLevel",
		"ReputationLevelType@enum",
	)

	addDataTableSheet("PlayerCampReputationShopItem",
		"CampType@enum",
		"ReputationLevelType@enum",
		"ProductLimitType@enum",
	)

	addDataTableSheet("SecretDungeon",
		"DifficultyType@enum",
	)

	addDataTableSheet("SecretDungeonPool",
		"Pool@array1@ID:int32@,",
	)

	addDataTableSheet("BossKilledReward",
		"WeekdaySUpID_1@array1@ID:int32@,",
		"WeekdaySUpID_2@array1@ID:int32@,",
		"WeekdaySUpID_3@array1@ID:int32@,",
		"WeekdaySUpID_4@array1@ID:int32@,",
		"WeekdaySUpID_5@array1@ID:int32@,",
		"WeekdaySUpID_6@array1@ID:int32@,",
		"WeekdayPoolUp@array1@ID:int32@,",
	)

	addDataTableSheet("TrialDungeon")

	addDataTableSheet("TrialStageReward")

	addDataTableSheet("GoldSupply")

	addDataTableSheet("Quest",
		"Type@enum",
		"Condition@array2@Type:enum,Arg:string@CondType_,CondArg_",
		"UnlockQuests@array1@QuestID:int32@,",
	)

	addDataTableSheet("Mail")

	addDataTableSheet("Sign",
		"Type@enum",
	)

	addDataTableSheet("SignAward")
	addDataTableSheet(
		"NewbieGift",
		"QuestIDList@array1@QuestID:int32@,",
	)

	addDataTableSheet("BigSecretDungeon",
		"DungeonIDPool@array1@DungeonID:int32@,")

	addDataTableSheet("BigSecretWeakness")

	addDataTableSheet("BigSecretBlessing")

	addDataTableSheet("BigSecretCompetition",
		"Quest@array1@QuestID:int32@,")

	addDataTableSheet("Talent")

	addDataTableSheet("TalentGroup",
		"Talents@array1@ID:int32@,")

	addDataTableSheet("TalentLevel",
		"PreTalentLevel@array1@ID:int32,Level:int32@,#")

	addDataTableSheet("TalentType")
}
