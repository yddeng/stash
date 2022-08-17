package Quest

import (
	"errors"
	"fmt"
	"initialthree/node/common/enumType"
	"initialthree/node/table/excel"
	"strconv"
	"strings"
	"sync/atomic"
)

type Condition struct {
	Type      int32
	DoneTimes int // 完成次数 默认1次

	// AccumulateLogin	累计登陆	登陆次数	1
	// WeeklyLogin	本周累积登录count	登录次数	5	查历史进度
	LoginTimes int64

	// OnlineTime	在线时长	时间（秒）600
	OnlineTime int64

	// ConsumeResource	资源消耗	资源类型#资源ID#数量  Item#1001#1	UsualAttribute#5#100000
	ResourceType  int32
	ResourceID    int32
	ResourceCount int32

	// DailyActiveness	日常活跃度达到某值	活跃度 20
	// WeeklyActiveness	本周累积每日活跃度	本周每日活跃度之和	500	查历史进度
	ActivenessCount int64

	// ExchangeFatigue	通过道具兑换体力	次数 1

	// PlayerLevelUp	玩家等级达到某级	玩家等级 100
	PlayerLevel int32

	// AnyInstanceSucceed	通关任意副本	通关次数 1
	// AnySecretInstanceSucceed	通关任意秘境副本	通关次数	10	需要执行

	// InstanceSucceed	通关指定副本	副本DungeonID 10101
	// AnyMaterialInstanceSucceed	通关任意材料副本	通关次数 1
	// AnyScarsInstanceSucceed	通关任意战痕副本	通关次数 1
	// AnyWorldQuestInstanceSucceed	任意悬赏任务副本 通关次数 1
	DungeonID int32

	// AnyRoleLevelUp	任意角色升级	无参数

	// AnyRoleBreak	任意角色突破count次	角色突破次数 5
	// RoleBreak	指定角色突破	角色ID#角色突破等级 14#3
	// RoleLevelUp	指定角色升级	角色ID#角色等级 13#80
	// RoleLevel	任意count角色等级达到指定level	角色数量#角色等级	3#10	查历史进度
	CharacterID         int32
	CharacterBreakLevel int32
	CharacterLevel      int32
	CharacterCount      int32

	// AnyWeaponLevelUp	任意武器升级count次	武器升级次数 6
	// AnyWeaponBreak	任意武器突破count次	武器突破次数 6
	// AnyWeaponRefine	任意武器精炼count次	武器精炼次数 6

	// WeaponOwned	拥有等级高于level的武器count件	武器等级#武器数量 50#5
	// WeaponEquipped	装备一件稀有度为rarity，等级为level的武器	稀有度#武器等级 Star5#40
	// WeaponRarityOwned	拥有稀有度rarity以上的武器count件	稀有度#数量	Star5#4	查历史进度
	WeaponLevel  int32
	WeaponCount  int32
	WeaponRarity int32

	// AnyEquipLevelUp	任意装备升级count次	装备升级次数	5
	// AnyEquipRefine	任意装备精炼count次	装备精炼次数	5

	// EquipOwned	拥有等级高于level的束流器count件	束流器等级#束流器数量	20#3
	// EquipEquipped	装备等级高于level的束流器count件	束流器等级#束流器数量	40#5
	// RoleEquipCount	任意角色已装备束流器件数	装备数量	5	查历史进度
	// EquipRarityOwned	拥有稀有度rarity以上的装备count件	稀有度#数量	Star5#4	查历史进度
	EquipLevel  int32
	EquipCount  int32
	EquipRarity int32

	// QuestComplete	完成指定任务	任务ID	1001
	QuestID int32

	// FatigueRegen	回复体力	次数	1
	// AnyShopPurchase	任意商店购买道具	次数	1

	// BranchInstancePass	通过任意支线章节	通关支线章节数目	1	查历史进度
	// CharacterInstancePass	通过任意角色传记	通关传记章节数目	1	查历史进度

	// RoleSkillLevelUp	任意角色技能升级count(所有角色累计)	角色技能升级次数	5	查历史进度

	// RoleSkillLevel	任意count技能达到指定level	技能数量#技能等级	5#5	查历史进度
	RoleSkillLevel int32
	RoleSkillCount int32

	// TrailCount	黑月探秘累计进度	累计进度	300	查历史进度
	TrailCount int32
	// AnyTrailInstanceSucceed	通过黑月探秘涨潮期任意关卡	通关次数	20	需要执行(可重复)

	// ScarsIngrainDifficultInstanceSucceed	通关战痕印刻副本指定难度指定次数	难度枚举#通关次数	5#10	需要执行(不可重复)
	ScarsDifficult int32
	ScarsPassTimes int32
	// ScarsIngrainScore	战痕印刻积分达到指定分数	积分	10000	需要执行
	ScarsScore int32
	// AnyDrawcardTimes	任意卡池抽卡count	抽卡次数	10	需要执行

	// BigSecretPassTimes	大秘境成功挑战次数	次数	5	需要执行
	// BigSecretBestScore	大秘境通关最高层级	最高层级	100	查历史进度
}

var questConditions atomic.Value

func (this *Quest) Conditions() []*Condition {
	return questConditions.Load().(map[*Quest][]*Condition)[this]
}

func (this *Table) AfterLoad() {
	idMap := GetIDMap()

	questConds := map[*Quest][]*Condition{}

	for _, q := range idMap {
		if len(q.Condition) == 0 {
			panic(fmt.Sprintf("quest %d condition is nil", q.ID))
		}

		conds := make([]*Condition, 0, len(q.Condition))
		for _, v := range q.Condition {
			cond := &Condition{
				Type:      v.Type,
				DoneTimes: 1,
			}

			var err error
			var number int

			switch v.Type {
			case enumType.QuestConditionType_DailyLogin: // 每日登陆	无参数

			// case enumType.QuestConditionType_AccumulateLogin: // todo
			//	number, err = strconv.Atoi(v.Arg)
			//	cond.LoginTimes = int64(number)
			case enumType.QuestConditionType_OnlineTime: // 在线时长	时间（秒）
				number, err = strconv.Atoi(v.Arg)
				cond.OnlineTime = int64(number)
			case enumType.QuestConditionType_ConsumeResource: // 资源消耗	资源类型#资源ID#数量
				s := strings.Split(v.Arg, "#")
				if len(s) != 3 {
					err = errors.New(fmt.Sprintf("condition argument %s is failed", v.Arg))
					break
				}
				cond.ResourceType = excel.ReadEnum("ResourceConsumeType_" + s[0])
				number, err = strconv.Atoi(s[1])
				if err != nil {
					break
				}
				cond.ResourceID = int32(number)
				number, err = strconv.Atoi(s[2])
				cond.ResourceCount = int32(number)
			// case enumType.QuestConditionType_YuruInteract: // todo
			case enumType.QuestConditionType_DailyActiveness: // 日常活跃度达到某值	活跃度
				number, err = strconv.Atoi(v.Arg)
				cond.ActivenessCount = int64(number)
			// case enumType.QuestConditionType_EnterDorm: // todo
			case enumType.QuestConditionType_ExchangeFatigue: // 通过道具兑换体力	次数
				number, err = strconv.Atoi(v.Arg)
				cond.DoneTimes = number
			case enumType.QuestConditionType_PlayerLevelUp: // 玩家等级达到某级	玩家等级
				number, err = strconv.Atoi(v.Arg)
				cond.PlayerLevel = int32(number)
			case enumType.QuestConditionType_AnyInstanceSucceed: // 通关任意副本	通关次数
				number, err = strconv.Atoi(v.Arg)
				cond.DoneTimes = number
			case enumType.QuestConditionType_InstanceSucceed: // 通关指定副本	副本DungeonID
				number, err = strconv.Atoi(v.Arg)
				cond.DungeonID = int32(number)
			case enumType.QuestConditionType_AnyMaterialInstanceSucceed: // 通关任意材料副本	通关次数 1
				number, err = strconv.Atoi(v.Arg)
				cond.DoneTimes = number
			case enumType.QuestConditionType_AnyScarsInstanceSucceed: // 通关任意战痕副本 通关次数 1
				number, err = strconv.Atoi(v.Arg)
				cond.DoneTimes = number
			case enumType.QuestConditionType_AnyRoleLevelUp: // 任意角色升级	无参数

			case enumType.QuestConditionType_AnyRoleBreak: // 任意角色突破count次	角色突破次数
				number, err = strconv.Atoi(v.Arg)
				cond.DoneTimes = number
			case enumType.QuestConditionType_RoleBreak: // 指定角色突破	角色ID#角色突破等级
				s := strings.Split(v.Arg, "#")
				if len(s) != 2 {
					err = errors.New(fmt.Sprintf("condition argument %s is failed", v.Arg))
					break
				}
				number, err = strconv.Atoi(s[0])
				if err != nil {
					break
				}
				cond.CharacterID = int32(number)
				number, err = strconv.Atoi(s[1])
				cond.CharacterBreakLevel = int32(number)
			case enumType.QuestConditionType_RoleLevelUp: // 指定角色升级	角色ID#角色等级
				s := strings.Split(v.Arg, "#")
				if len(s) != 2 {
					err = errors.New(fmt.Sprintf("condition argument %s is failed", v.Arg))
					break
				}
				number, err = strconv.Atoi(s[0])
				if err != nil {
					break
				}
				cond.CharacterID = int32(number)
				number, err = strconv.Atoi(s[1])
				cond.CharacterLevel = int32(number)
			case enumType.QuestConditionType_AnyWeaponLevelUp: // 任意武器升级count次	武器升级次数
				number, err = strconv.Atoi(v.Arg)
				cond.DoneTimes = number
			case enumType.QuestConditionType_AnyWeaponBreak: // 任意武器突破count次	武器突破次数
				number, err = strconv.Atoi(v.Arg)
				cond.DoneTimes = number
			case enumType.QuestConditionType_AnyWeaponRefine: // 任意武器精炼count次	武器精炼次数
				number, err = strconv.Atoi(v.Arg)
				cond.DoneTimes = number
			case enumType.QuestConditionType_WeaponOwned: // 拥有等级高于level的武器count件	武器等级#武器数量
				s := strings.Split(v.Arg, "#")
				if len(s) != 2 {
					err = errors.New(fmt.Sprintf("condition argument %s is failed", v.Arg))
					break
				}
				number, err = strconv.Atoi(s[0])
				if err != nil {
					break
				}
				cond.WeaponLevel = int32(number)
				number, err = strconv.Atoi(s[1])
				cond.WeaponCount = int32(number)
			case enumType.QuestConditionType_WeaponEquipped: // 装备一件稀有度为rarity，等级为level的武器	稀有度#武器等级
				s := strings.Split(v.Arg, "#")
				if len(s) != 2 {
					err = errors.New(fmt.Sprintf("condition argument %s is failed", v.Arg))
					break
				}
				cond.WeaponRarity, err = enumType.GetEnumType("RarityType_" + s[0])
				if err != nil {
					break
				}
				number, err = strconv.Atoi(s[1])
				cond.WeaponLevel = int32(number)
			case enumType.QuestConditionType_AnyEquipLevelUp: // 任意装备升级count次	装备升级次数
				number, err = strconv.Atoi(v.Arg)
				cond.DoneTimes = number
			case enumType.QuestConditionType_AnyEquipRefine: // 任意装备精炼count次	装备精炼次数
				number, err = strconv.Atoi(v.Arg)
				cond.DoneTimes = number
			case enumType.QuestConditionType_EquipOwned: // 拥有等级高于level的束流器count件	束流器等级#束流器数量
				s := strings.Split(v.Arg, "#")
				if len(s) != 2 {
					err = errors.New(fmt.Sprintf("condition argument %s is failed", v.Arg))
					break
				}
				number, err = strconv.Atoi(s[0])
				if err != nil {
					break
				}
				cond.EquipLevel = int32(number)
				number, err = strconv.Atoi(s[1])
				cond.EquipCount = int32(number)
			case enumType.QuestConditionType_EquipEquipped: // 装备等级高于level的束流器count件	束流器等级#束流器数量
				s := strings.Split(v.Arg, "#")
				if len(s) != 2 {
					err = errors.New(fmt.Sprintf("condition argument %s is failed", v.Arg))
					break
				}
				number, err = strconv.Atoi(s[0])
				if err != nil {
					break
				}
				cond.EquipLevel = int32(number)
				number, err = strconv.Atoi(s[1])
				cond.EquipCount = int32(number)
			case enumType.QuestConditionType_QuestComplete: // 完成指定任务	任务ID
				number, err = strconv.Atoi(v.Arg)
				cond.QuestID = int32(number)
			case enumType.QuestConditionType_FatigueRegen: // 回复体力	次数
				number, err = strconv.Atoi(v.Arg)
				cond.DoneTimes = number
			case enumType.QuestConditionType_AnyWorldQuestInstanceSucceed: // 任意悬赏任务副本	通关次数 1
				number, err = strconv.Atoi(v.Arg)
				cond.DoneTimes = number
			case enumType.QuestConditionType_AnyShopPurchase: // 任意商店购买道具	次数
				number, err = strconv.Atoi(v.Arg)
				cond.DoneTimes = number
			case enumType.QuestConditionType_BranchInstancePass: // 通过任意支线章节	通关支线章节数目	1	查历史进度
				number, err = strconv.Atoi(v.Arg)
				cond.DoneTimes = number
			case enumType.QuestConditionType_CharacterInstancePass: // 通过任意角色传记	通关传记章节数目	1	查历史进度
				number, err = strconv.Atoi(v.Arg)
				cond.DoneTimes = number
			case enumType.QuestConditionType_RoleLevel: // 任意count角色等级达到指定level	角色数量#角色等级	3#10	查历史进度
				s := strings.Split(v.Arg, "#")
				if len(s) != 2 {
					err = errors.New(fmt.Sprintf("condition argument %s is failed", v.Arg))
					break
				}
				number, err = strconv.Atoi(s[0])
				if err != nil {
					break
				}
				cond.CharacterCount = int32(number)
				number, err = strconv.Atoi(s[1])
				cond.CharacterLevel = int32(number)
			case enumType.QuestConditionType_RoleSkillLevelUp: // 任意角色技能升级count(所有角色累计)	角色技能升级次数	5	查历史进度
				number, err = strconv.Atoi(v.Arg)
				cond.DoneTimes = number
			case enumType.QuestConditionType_RoleSkillLevel: // 任意count技能达到指定level	技能数量#技能等级	5#5	查历史进度
				s := strings.Split(v.Arg, "#")
				if len(s) != 2 {
					err = errors.New(fmt.Sprintf("condition argument %s is failed", v.Arg))
					break
				}
				number, err = strconv.Atoi(s[0])
				if err != nil {
					break
				}
				cond.RoleSkillCount = int32(number)
				number, err = strconv.Atoi(s[1])
				cond.RoleSkillLevel = int32(number)
			case enumType.QuestConditionType_RoleEquipCount: // 任意角色已装备束流器件数	装备数量	5	查历史进度
				number, err = strconv.Atoi(v.Arg)
				cond.EquipCount = int32(number)
			case enumType.QuestConditionType_AnySecretInstanceSucceed: // 通关任意秘境副本	通关次数	10	需要执行
				number, err = strconv.Atoi(v.Arg)
				cond.DoneTimes = number
			case enumType.QuestConditionType_TrailCount: // 黑月探秘累计进度	累计进度	300	查历史进度
				number, err = strconv.Atoi(v.Arg)
				cond.TrailCount = int32(number)
			case enumType.QuestConditionType_AnyTrailInstanceSucceed: // 通过黑月探秘涨潮期任意关卡	通关次数	20	需要执行(可重复)
				number, err = strconv.Atoi(v.Arg)
				cond.DoneTimes = number
			case enumType.QuestConditionType_ScarsIngrainDifficultInstanceSucceed: // 通关战痕印刻副本指定难度指定次数	难度枚举#通关次数	5#10	需要执行(不可重复)
				s := strings.Split(v.Arg, "#")
				if len(s) != 2 {
					err = errors.New(fmt.Sprintf("condition argument %s is failed", v.Arg))
					break
				}
				number, err = strconv.Atoi(s[0])
				if err != nil {
					break
				}
				cond.ScarsDifficult = int32(number)
				number, err = strconv.Atoi(s[1])
				cond.ScarsPassTimes = int32(number)
			case enumType.QuestConditionType_ScarsIngrainScore: // 战痕印刻积分达到指定分数	积分	10000	需要执行
				number, err = strconv.Atoi(v.Arg)
				cond.ScarsScore = int32(number)
			case enumType.QuestConditionType_AnyDrawcardTimes: // 任意卡池抽卡count	抽卡次数	10	需要执行
				number, err = strconv.Atoi(v.Arg)
				cond.DoneTimes = number
			case enumType.QuestConditionType_WeeklyLogin: // 本周累积登录count	登录次数	5	查历史进度
				number, err = strconv.Atoi(v.Arg)
				cond.LoginTimes = int64(number)
			case enumType.QuestConditionType_WeeklyActiveness: // 本周累积每日活跃度	本周每日活跃度之和	500	查历史进度
				number, err = strconv.Atoi(v.Arg)
				cond.ActivenessCount = int64(number)
			case enumType.QuestConditionType_EquipRarityOwned: // EquipRarityOwned	拥有稀有度rarity以上的装备count件	稀有度#数量	Star5#4	查历史进度
				s := strings.Split(v.Arg, "#")
				if len(s) != 2 {
					err = errors.New(fmt.Sprintf("condition argument %s is failed", v.Arg))
					break
				}
				cond.EquipRarity, err = enumType.GetEnumType("RarityType_" + s[0])
				if err != nil {
					break
				}
				number, err = strconv.Atoi(s[1])
				cond.EquipCount = int32(number)
			case enumType.QuestConditionType_WeaponRarityOwned: // 拥有稀有度rarity以上的武器count件	稀有度#数量	Star5#4	查历史进度
				s := strings.Split(v.Arg, "#")
				if len(s) != 2 {
					err = errors.New(fmt.Sprintf("condition argument %s is failed", v.Arg))
					break
				}
				cond.WeaponRarity, err = enumType.GetEnumType("RarityType_" + s[0])
				if err != nil {
					break
				}
				number, err = strconv.Atoi(s[1])
				cond.WeaponCount = int32(number)
			case enumType.QuestConditionType_BigSecretPassTimes: // 大秘境成功挑战次数	次数	5	需要执行
				number, err = strconv.Atoi(v.Arg)
				cond.DoneTimes = number
			case enumType.QuestConditionType_BigSecretBestScore: // 大秘境通关最高层级	最高层级	100	查历史进度
				number, err = strconv.Atoi(v.Arg)
				cond.DoneTimes = number
			default:
				err = errors.New(fmt.Sprintf("condition type %d not register", v.Type))
			}
			if err != nil {
				panic(fmt.Sprintf("quest %d  err %s", q.ID, err))
			}

			conds = append(conds, cond)
		}
		questConds[q] = conds
	}

	questConditions.Store(questConds)
}
