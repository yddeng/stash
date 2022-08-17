package main

import (
    "fmt"
    "initialthree/node/table/excel"
    "initialthree/zaplogger"
    "os"
    "runtime"
    _ "initialthree/node/table/excel/ConstTable/BigSecret"
    _ "initialthree/node/table/excel/ConstTable/Equip"
    _ "initialthree/node/table/excel/ConstTable/Global"
    _ "initialthree/node/table/excel/ConstTable/MainChapter"
    _ "initialthree/node/table/excel/ConstTable/Player"
    _ "initialthree/node/table/excel/ConstTable/PlayerCamp"
    _ "initialthree/node/table/excel/ConstTable/Quest"
    _ "initialthree/node/table/excel/ConstTable/ScarsIngrain"
    _ "initialthree/node/table/excel/ConstTable/Talent"
    _ "initialthree/node/table/excel/ConstTable/Weapon"
    _ "initialthree/node/table/excel/DataTable/Map"
    _ "initialthree/node/table/excel/DataTable/Skill"
    _ "initialthree/node/table/excel/DataTable/RewardQuest"
    _ "initialthree/node/table/excel/DataTable/RewardQuestPosition"
    _ "initialthree/node/table/excel/DataTable/Quest"
    _ "initialthree/node/table/excel/DataTable/NewbieGift"
    _ "initialthree/node/table/excel/DataTable/Dungeon"
    _ "initialthree/node/table/excel/DataTable/WorldQuest"
    _ "initialthree/node/table/excel/DataTable/WorldQuestPool"
    _ "initialthree/node/table/excel/DataTable/MainChapter"
    _ "initialthree/node/table/excel/DataTable/MainDungeon"
    _ "initialthree/node/table/excel/DataTable/BigSecretCompetition"
    _ "initialthree/node/table/excel/DataTable/BigSecretDungeon"
    _ "initialthree/node/table/excel/DataTable/BigSecretBlessing"
    _ "initialthree/node/table/excel/DataTable/ScarsIngrainArea"
    _ "initialthree/node/table/excel/DataTable/ScarsIngrainBoss"
    _ "initialthree/node/table/excel/DataTable/ScarsIngrainBossChallenge"
    _ "initialthree/node/table/excel/DataTable/ScarsIngrainBossInstance"
    _ "initialthree/node/table/excel/DataTable/ScarsIngrainRankReward"
    _ "initialthree/node/table/excel/DataTable/ScarsIngrainScoreReward"
    _ "initialthree/node/table/excel/DataTable/MaterialDungeon"
    _ "initialthree/node/table/excel/DataTable/BossKilledReward"
    _ "initialthree/node/table/excel/DataTable/SecretDungeon"
    _ "initialthree/node/table/excel/DataTable/SecretDungeonPool"
    _ "initialthree/node/table/excel/DataTable/TrialDungeon"
    _ "initialthree/node/table/excel/DataTable/TrialStageReward"
    _ "initialthree/node/table/excel/DataTable/Resurrect"
    _ "initialthree/node/table/excel/DataTable/Function"
    _ "initialthree/node/table/excel/DataTable/Pay"
    _ "initialthree/node/table/excel/DataTable/Product"
    _ "initialthree/node/table/excel/DataTable/ProductLibrary"
    _ "initialthree/node/table/excel/DataTable/Talent"
    _ "initialthree/node/table/excel/DataTable/TalentGroup"
    _ "initialthree/node/table/excel/DataTable/TalentLevel"
    _ "initialthree/node/table/excel/DataTable/TalentType"
    _ "initialthree/node/table/excel/DataTable/DrawCardsLib"
    _ "initialthree/node/table/excel/DataTable/DrawCardsPool"
    _ "initialthree/node/table/excel/DataTable/DropPool"
    _ "initialthree/node/table/excel/DataTable/Weapon"
    _ "initialthree/node/table/excel/DataTable/WeaponAttribute"
    _ "initialthree/node/table/excel/DataTable/WeaponBreakThough"
    _ "initialthree/node/table/excel/DataTable/WeaponDecompose"
    _ "initialthree/node/table/excel/DataTable/WeaponLevelMaxExp"
    _ "initialthree/node/table/excel/DataTable/WeaponRarity"
    _ "initialthree/node/table/excel/DataTable/FatigueSupply"
    _ "initialthree/node/table/excel/DataTable/GoldSupply"
    _ "initialthree/node/table/excel/DataTable/PlayerLevel"
    _ "initialthree/node/table/excel/DataTable/PlayerCard"
    _ "initialthree/node/table/excel/DataTable/PlayerPortrait"
    _ "initialthree/node/table/excel/DataTable/PlayerPortraitFrame"
    _ "initialthree/node/table/excel/DataTable/PlayerCampReputationLevel"
    _ "initialthree/node/table/excel/DataTable/PlayerCampReputationShopItem"
    _ "initialthree/node/table/excel/DataTable/Sign"
    _ "initialthree/node/table/excel/DataTable/SignAward"
    _ "initialthree/node/table/excel/DataTable/Equip"
    _ "initialthree/node/table/excel/DataTable/EquipAttribute"
    _ "initialthree/node/table/excel/DataTable/EquipDecompose"
    _ "initialthree/node/table/excel/DataTable/EquipLevelMaxExp"
    _ "initialthree/node/table/excel/DataTable/EquipQuality"
    _ "initialthree/node/table/excel/DataTable/EquipRandomAttributePool"
    _ "initialthree/node/table/excel/DataTable/EquipSkill"
    _ "initialthree/node/table/excel/DataTable/CharacterBreakThrough"
    _ "initialthree/node/table/excel/DataTable/CharacterLevelUpAttribute"
    _ "initialthree/node/table/excel/DataTable/CharacterLevelUpExp"
    _ "initialthree/node/table/excel/DataTable/CharacterResource"
    _ "initialthree/node/table/excel/DataTable/DemoCharacter"
    _ "initialthree/node/table/excel/DataTable/PlayerCharacter"
    _ "initialthree/node/table/excel/DataTable/PlayerGene"
    _ "initialthree/node/table/excel/DataTable/PlayerSkill"
    _ "initialthree/node/table/excel/DataTable/AccountInitializeAssets"
    _ "initialthree/node/table/excel/DataTable/FragmentChange"
    _ "initialthree/node/table/excel/DataTable/Gift"
    _ "initialthree/node/table/excel/DataTable/InputOutput"
    _ "initialthree/node/table/excel/DataTable/Item"
    _ "initialthree/node/table/excel/DataTable/Mail"

)

func main() {
    if len(os.Args) < 2 {
        panic("need argument: excel_path")
    }

    logger := zaplogger.NewZapLogger("check.log","log","debug", 100, 14, 10, true)
    zaplogger.InitLogger(logger)

    excel_path := os.Args[1]
    if excel_path == "" {
        panic("need argument: excel_path")
    }
	defer func() {
		if r := recover(); r != nil {
			buf := make([]byte, 65535)
			l := runtime.Stack(buf, false)
			fmt.Println("check failed :", r)
			fmt.Println(string(buf[:l]))
		} else {
			fmt.Println("check table ok")
		}
	}()
    
	excel.Load(excel_path)
}
