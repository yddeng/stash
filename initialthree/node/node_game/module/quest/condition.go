package quest

import (
	"fmt"
	"initialthree/node/common/attr"
	"initialthree/node/common/enumType"
	"initialthree/node/node_game/event"
	"initialthree/node/node_game/module"
	attr2 "initialthree/node/node_game/module/attr"
	"initialthree/node/node_game/module/bigSecret"
	"initialthree/node/node_game/module/character"
	"initialthree/node/node_game/module/equip"
	"initialthree/node/node_game/module/maindungeons"
	"initialthree/node/node_game/module/materialDungeon"
	"initialthree/node/node_game/module/scarsIngrain"
	"initialthree/node/node_game/module/trialDungeon"
	"initialthree/node/node_game/module/weapon"
	"initialthree/node/table/excel/DataTable/Dungeon"
	"initialthree/node/table/excel/DataTable/Equip"
	"initialthree/node/table/excel/DataTable/MainChapter"
	Quest2 "initialthree/node/table/excel/DataTable/Quest"
	"initialthree/node/table/excel/DataTable/ScarsIngrainBossInstance"
	"initialthree/node/table/excel/DataTable/Weapon"
	vendor_event "initialthree/pkg/event"
	"initialthree/protocol/cs/message"
	"initialthree/zaplogger"
)

type CondEvent interface {
	/*
		验证依赖模块状态是否可用。
		例 任务条件是属性活跃度，当任务先日刷新，新任务去读取属性活跃度，但是属性还未日刷新重置值，导致获取的属性是昨天的值。
		当前检测方法： 任务条件注册时需判断依赖模块是否也刷新了，若没有不注册事件，由tick来定时检测能否注册。
	*/
	// 验证任务条件依赖的模块状态是否满足，目前仅判断刷新
	Check() bool
	// 任务条件注册前，检测是否已经满足条件
	Before() bool
	Trigger()
	UnTrigger()
	Done()
}

type condEventBase struct {
	condTemp *Quest2.Condition
	cond     *Condition
	q        *Quest
	uq       *UserQuest
	hs       []vendor_event.Handle
}

func newCondEventBase(cond *Condition, condTemp *Quest2.Condition, q *Quest, uq *UserQuest) *condEventBase {
	return &condEventBase{condTemp: condTemp, cond: cond, q: q, uq: uq, hs: make([]vendor_event.Handle, 0, 2)}
}

func (this *condEventBase) setDirty() {
	this.uq.SetDirty(this.q.questTemp.TypeEnum)
	this.uq.questDirty[this.q.ID] = struct{}{}
}

func (this *condEventBase) trigger(e interface{}, fn interface{}) {
	h := this.uq.userI.RegisterEvent(e, fn)
	this.hs = append(this.hs, h)
}

func (this *condEventBase) Check() bool {
	// 默认依赖模块满条件。
	return true
}

func (this *condEventBase) Done() {
	zaplogger.GetSugar().Debugf("CondEvent Done questID %d condition %v ", this.q.ID, this.condTemp)

	this.cond.Complete = true
	this.uq.tryAllCondComplete(this.q)
	this.setDirty()

	this.UnTrigger()
}

func (this *condEventBase) UnTrigger() {
	if len(this.hs) != 0 {
		for _, h := range this.hs {
			this.uq.userI.UnRegisterEvent(h)
		}
		this.hs = this.hs[0:0]
	}
	this.cond.condEvent = nil
}

/***************************** 每日登陆 **************************************************************************************/
// DailyLogin	每日登陆	无参数
type dailyLogin struct {
	*condEventBase
}

func (this *dailyLogin) Before() bool {
	this.cond.DoneTimes = 1
	return true
}
func (this *dailyLogin) Trigger() {}

/***************************** 累计登陆 **************************************************************************************/
// AccumulateLogin	累计登陆	登陆次数	1
type accumulateLogin struct {
	*condEventBase
}

func (this *accumulateLogin) Before() bool {
	attrModule := this.uq.userI.GetSubModule(module.Attr).(*attr2.UserAttr)
	val, _ := attrModule.GetAttr(attr.AccumulateLogin)
	if this.cond.DoneTimes != int32(val) {
		this.cond.DoneTimes = int32(val)
		this.setDirty()
		if this.cond.DoneTimes >= int32(this.condTemp.LoginTimes) {
			this.cond.DoneTimes = int32(this.condTemp.LoginTimes)
			return true
		}
	}
	return false
}
func (this *accumulateLogin) Trigger() {
	trigger := func(id int32, oldVal, newVal int64) {
		if id == attr.AccumulateLogin && this.cond.DoneTimes != int32(newVal) {
			this.cond.DoneTimes = int32(newVal)
			this.setDirty()
			if this.cond.DoneTimes >= int32(this.condTemp.LoginTimes) {
				this.cond.DoneTimes = int32(this.condTemp.LoginTimes)
				this.Done()
			}
		}
	}
	this.trigger(event.EventAttrChange, trigger)
	this.cond.condEvent = this
}

/***************************** 在线时长 **************************************************************************************/
// OnlineTime	在线时长	时间（秒）	600
type onlineTime struct {
	*condEventBase
}

func (this *onlineTime) Check() bool {
	attrModule := this.uq.userI.GetSubModule(module.Attr).(*attr2.UserAttr)
	return attrModule.TimeEqual(module.DailyTimeName, this.uq.timedata[module.DailyTimeName])
}

func (this *onlineTime) Before() bool {
	attrModule := this.uq.userI.GetSubModule(module.Attr).(*attr2.UserAttr)
	val, _ := attrModule.GetAttr(attr.DailyOnLine)
	if this.cond.DoneTimes != int32(val) {
		this.cond.DoneTimes = int32(val)
		this.setDirty()
		if this.cond.DoneTimes >= int32(this.condTemp.OnlineTime) {
			this.cond.DoneTimes = int32(this.condTemp.OnlineTime)
			return true
		}
	}
	return false
}
func (this *onlineTime) Trigger() {
	trigger := func(id int32, oldVal, newVal int64) {
		if id == attr.DailyOnLine && this.cond.DoneTimes != int32(newVal) {
			this.cond.DoneTimes = int32(newVal)
			this.setDirty()
			if this.cond.DoneTimes >= int32(this.condTemp.OnlineTime) {
				this.cond.DoneTimes = int32(this.condTemp.OnlineTime)
				this.Done()
			}
		}
	}
	this.trigger(event.EventAttrChange, trigger)
	this.cond.condEvent = this
}

/***************************** 资源消耗 **************************************************************************************/
// ConsumeResource	资源消耗	资源类型#资源ID#数量	Item#1001#1	UsualAttribute#5#100000
type consumeResource struct {
	*condEventBase
}

func (this *consumeResource) Before() bool { return false }
func (this *consumeResource) Trigger() {

	switch this.condTemp.ResourceType {
	case enumType.ResourceConsumeType_UsualAttribute:
		trigger := func(id int32, oldVal, newVal int64) {
			// newVal > oldVal 为消耗
			used := newVal - oldVal
			if id == this.condTemp.ResourceID && used < 0 {
				this.cond.DoneTimes += -int32(used)
				this.setDirty()
				if this.cond.DoneTimes >= this.condTemp.ResourceCount {
					this.cond.DoneTimes = this.condTemp.ResourceCount
					this.Done()
				}
			}
		}
		this.trigger(event.EventAttrChange, trigger)
	case enumType.ResourceConsumeType_Item:
		trigger := func(items map[int32]int32) {
			count, ok := items[this.condTemp.ResourceID]
			if ok {
				this.cond.DoneTimes += count
				this.setDirty()
				if this.cond.DoneTimes >= this.condTemp.ResourceCount {
					this.cond.DoneTimes = this.condTemp.ResourceCount
					this.Done()
				}
			}
		}
		this.trigger(event.EventUseItem, trigger)
	}
	this.cond.condEvent = this
}

/***************************** 与看板娘交互 **************************************************************************************/
// YuruInteract	与看板娘交互	次数	1
type yuruInteract struct {
	*condEventBase
}

func (this *yuruInteract) Before() bool { return false }
func (this *yuruInteract) Trigger()     {}

/***************************** 日常活跃度达到某值 **************************************************************************************/
// DailyActiveness	日常活跃度达到某值	活跃度	20
type dailyActiveness struct {
	*condEventBase
}

func (this *dailyActiveness) Check() bool {
	attrModule := this.uq.userI.GetSubModule(module.Attr).(*attr2.UserAttr)
	return attrModule.TimeEqual(module.DailyTimeName, this.uq.timedata[module.DailyTimeName])
}

func (this *dailyActiveness) Before() bool {
	attrModule := this.uq.userI.GetSubModule(module.Attr).(*attr2.UserAttr)
	val, _ := attrModule.GetAttr(attr.DailyActiveness)
	if this.cond.DoneTimes != int32(val) {
		this.cond.DoneTimes = int32(val)
		this.setDirty()
		if this.cond.DoneTimes >= int32(this.condTemp.ActivenessCount) {
			this.cond.DoneTimes = int32(this.condTemp.ActivenessCount)
			return true
		}
	}
	return false
}
func (this *dailyActiveness) Trigger() {
	trigger := func(id int32, oldVal, newVal int64) {
		if id == attr.DailyActiveness && this.cond.DoneTimes != int32(newVal) {
			this.cond.DoneTimes = int32(newVal)
			this.setDirty()
			if this.cond.DoneTimes >= int32(this.condTemp.ActivenessCount) {
				this.cond.DoneTimes = int32(this.condTemp.ActivenessCount)
				this.Done()
			}
		}
	}
	this.trigger(event.EventAttrChange, trigger)
	this.cond.condEvent = this
}

/***************************** 进入宿舍 **************************************************************************************/
// EnterDorm	进入宿舍
type enterDorm struct {
	*condEventBase
}

func (this *enterDorm) Before() bool {
	return false
}
func (this *enterDorm) Trigger() {}

/***************************** 通过道具兑换体力 **************************************************************************************/
// ExchangeFatigue	通过道具兑换体力	次数	1
type exchangeFatigue struct {
	*condEventBase
}

func (this *exchangeFatigue) Before() bool { return false }
func (this *exchangeFatigue) Trigger() {
	trigger := func(times int32) {
		this.cond.DoneTimes += times
		this.setDirty()
		if this.cond.DoneTimes >= int32(this.condTemp.DoneTimes) {
			this.cond.DoneTimes = int32(this.condTemp.DoneTimes)
			this.Done()
		}
	}
	this.trigger(event.EventExchangeFatigue, trigger)
	this.cond.condEvent = this
}

/***************************** 玩家等级达到某级 **************************************************************************************/
// PlayerLevelUp	玩家等级达到某级	玩家等级	100
type playerLevelUp struct {
	*condEventBase
}

func (this *playerLevelUp) Before() bool {
	attrModule := this.uq.userI.GetSubModule(module.Attr).(*attr2.UserAttr)
	val, _ := attrModule.GetAttr(attr.Level)
	if this.cond.DoneTimes != int32(val) {
		this.cond.DoneTimes = int32(val)
		this.setDirty()
		if this.cond.DoneTimes >= this.condTemp.PlayerLevel {
			this.cond.DoneTimes = this.condTemp.PlayerLevel
			return true
		}
	}
	return false
}
func (this *playerLevelUp) Trigger() {
	trigger := func(id int32, oldVal, newVal int64) {
		if id == attr.Level && this.cond.DoneTimes != int32(newVal) {
			this.cond.DoneTimes = int32(newVal)
			this.setDirty()
			if this.cond.DoneTimes >= int32(this.condTemp.PlayerLevel) {
				this.cond.DoneTimes = int32(this.condTemp.PlayerLevel)
				this.Done()
			}
		}
	}
	this.trigger(event.EventAttrChange, trigger)
	this.cond.condEvent = this
}

/***************************** 通关任意副本 **************************************************************************************/
// AnyInstanceSucceed	通关任意副本	通关次数	1
type anyInstanceSucceed struct {
	*condEventBase
}

func (this *anyInstanceSucceed) Before() bool { return false }
func (this *anyInstanceSucceed) Trigger() {
	trigger := func(instanceID, instanceStar int32, monster map[int32]int32) {
		this.cond.DoneTimes++
		this.setDirty()
		if this.cond.DoneTimes >= int32(this.condTemp.DoneTimes) {
			this.cond.DoneTimes = int32(this.condTemp.DoneTimes)
			this.Done()
		}
	}
	this.trigger(event.EventInstanceSucceed, trigger)
	this.cond.condEvent = this
}

/***************************** 通关指定副本 **************************************************************************************/
// InstanceSucceed	通关指定副本	副本DungeonID	10101
type instanceSucceed struct {
	*condEventBase
}

func (this *instanceSucceed) Before() bool {
	def := Dungeon.GetID(this.condTemp.DungeonID)
	if def == nil {
		return false
	}
	switch def.SystemTypeEnum {
	case enumType.DungeonSystemType_MainChapter:
		mainDungeon := this.uq.userI.GetSubModule(module.MainDungeons).(*maindungeons.MainDungeons)
		if mainDungeon.IsDungeonPass(def.SystemConfig) {
			this.cond.DoneTimes = 1
			return true
		}
	case enumType.DungeonSystemType_Material:
		userMaterialLevel := this.uq.userI.GetSubModule(module.MaterialDungeon).(*materialDungeon.MaterialDungeon)
		if userMaterialLevel.GetMaterialDungeon(def.SystemConfig) != nil {
			this.cond.DoneTimes = 1
			return true
		}
	case enumType.DungeonSystemType_Secret:
	case enumType.DungeonSystemType_WorldQuest:
	case enumType.DungeonSystemType_Trial:
		userTrialDungeon := this.uq.userI.GetSubModule(module.TrialDungeon).(*trialDungeon.TrialDungeon)
		if userTrialDungeon.GetTrialDungeon(def.SystemConfig) != nil {
			this.cond.DoneTimes = 1
			return true
		}
	case enumType.DungeonSystemType_ScarsIngrain:
	}
	return false
}
func (this *instanceSucceed) Trigger() {
	trigger := func(instanceID, instanceStar int32, monster map[int32]int32) {
		if this.condTemp.DungeonID == instanceID {
			this.cond.DoneTimes = 1
			this.Done()
		}
	}
	this.trigger(event.EventInstanceSucceed, trigger)
	this.cond.condEvent = this
}

/***************************** 通关任意材料副本 **************************************************************************************/
// AnyMaterialInstanceSucceed	通关任意材料副本	通关次数	1
type anyMaterialInstanceSucceed struct {
	*condEventBase
}

func (this *anyMaterialInstanceSucceed) Before() bool { return false }
func (this *anyMaterialInstanceSucceed) Trigger() {
	trigger := func(instanceID, instanceStar int32, monster map[int32]int32) {
		def := Dungeon.GetID(instanceID)
		if def != nil && def.SystemTypeEnum == enumType.DungeonSystemType_Material {
			this.cond.DoneTimes++
			this.setDirty()
			if this.cond.DoneTimes >= int32(this.condTemp.DoneTimes) {
				this.cond.DoneTimes = int32(this.condTemp.DoneTimes)
				this.Done()
			}
		}
	}
	this.trigger(event.EventInstanceSucceed, trigger)
	this.cond.condEvent = this
}

/***************************** 通关任意战痕副本 **************************************************************************************/
// AnyScarsInstanceSucceed	通关任意战痕副本	通关次数	1
type anyScarsInstanceSucceed struct {
	*condEventBase
}

func (this *anyScarsInstanceSucceed) Before() bool { return false }
func (this *anyScarsInstanceSucceed) Trigger() {
	trigger := func(instanceID, instanceStar int32, monster map[int32]int32) {
		def := Dungeon.GetID(instanceID)
		if def != nil && def.SystemTypeEnum == enumType.DungeonSystemType_ScarsIngrain {
			this.cond.DoneTimes++
			this.setDirty()
			if this.cond.DoneTimes >= int32(this.condTemp.DoneTimes) {
				this.cond.DoneTimes = int32(this.condTemp.DoneTimes)
				this.Done()
			}
		}
	}
	this.trigger(event.EventInstanceSucceed, trigger)
	this.cond.condEvent = this
}

/***************************** 任意角色升级 **************************************************************************************/
// AnyRoleLevelUp	任意角色升级	无参数
type anyRoleLevelUp struct {
	*condEventBase
}

func (this *anyRoleLevelUp) Before() bool { return false }
func (this *anyRoleLevelUp) Trigger() {
	trigger := func(charaID, oldLevel, newLevel int32) {
		this.cond.DoneTimes = 1
		this.Done()
	}
	this.trigger(event.EventCharacterLevelUp, trigger)
	this.cond.condEvent = this
}

/***************************** 任意角色突破count次 **************************************************************************************/
// AnyRoleBreak	任意角色突破count次	角色突破次数	5
type anyRoleBreak struct {
	*condEventBase
}

func (this *anyRoleBreak) Before() bool { return false }
func (this *anyRoleBreak) Trigger() {
	trigger := func(charaID, oldLevel, newLevel int32) {
		this.cond.DoneTimes += newLevel - oldLevel
		this.setDirty()
		if this.cond.DoneTimes >= int32(this.condTemp.DoneTimes) {
			this.cond.DoneTimes = int32(this.condTemp.DoneTimes)
			this.Done()
		}
	}
	this.trigger(event.EventCharacterBreak, trigger)
	this.cond.condEvent = this
}

/***************************** 指定角色突破 **************************************************************************************/
// RoleBreak	指定角色突破	角色ID#角色突破等级	14#3
type roleBreak struct {
	*condEventBase
}

func (this *roleBreak) Before() bool {
	charaModule := this.uq.userI.GetSubModule(module.Character).(*character.UserCharacter)
	if c := charaModule.GetCharacter(this.condTemp.CharacterID); c != nil {
		if this.cond.DoneTimes != c.BreakLevel {
			this.cond.DoneTimes = c.BreakLevel
			this.setDirty()
			if this.cond.DoneTimes >= this.condTemp.CharacterBreakLevel {
				this.cond.DoneTimes = this.condTemp.CharacterBreakLevel
				return true
			}
		}
	}
	return false
}
func (this *roleBreak) Trigger() {
	trigger := func(charaID, oldLevel, newLevel int32) {
		if this.condTemp.CharacterID == charaID && this.cond.DoneTimes != newLevel {
			this.cond.DoneTimes = newLevel
			this.setDirty()
			if this.cond.DoneTimes >= this.condTemp.CharacterBreakLevel {
				this.cond.DoneTimes = this.condTemp.CharacterBreakLevel
				this.Done()
			}
		}
	}
	this.trigger(event.EventCharacterBreak, trigger)
	this.cond.condEvent = this
}

/***************************** 指定角色升级 **************************************************************************************/
// RoleLevelUp	指定角色升级	角色ID#角色等级	14#80
type roleLevelUp struct {
	*condEventBase
}

func (this *roleLevelUp) Before() bool {
	charaModule := this.uq.userI.GetSubModule(module.Character).(*character.UserCharacter)
	if c := charaModule.GetCharacter(this.condTemp.CharacterID); c != nil {
		if this.cond.DoneTimes != c.Level {
			this.cond.DoneTimes = c.Level
			this.setDirty()
			if this.cond.DoneTimes >= this.condTemp.CharacterLevel {
				this.cond.DoneTimes = this.condTemp.CharacterLevel
				return true
			}
		}
	}
	return false
}
func (this *roleLevelUp) Trigger() {
	trigger := func(charaID, oldLevel, newLevel int32) {
		if this.condTemp.CharacterID == charaID && this.cond.DoneTimes != newLevel {
			this.cond.DoneTimes = newLevel
			this.setDirty()
			if this.cond.DoneTimes >= this.condTemp.CharacterLevel {
				this.cond.DoneTimes = this.condTemp.CharacterLevel
				this.Done()
			}
		}
	}
	this.trigger(event.EventCharacterLevelUp, trigger)
	this.cond.condEvent = this
}

/***************************** 任意武器升级count次 **************************************************************************************/
// AnyWeaponLevelUp	任意武器升级count次	武器升级次数	5
type anyWeaponLevelUp struct {
	*condEventBase
}

func (this *anyWeaponLevelUp) Before() bool { return false }
func (this *anyWeaponLevelUp) Trigger() {
	trigger := func(weaponID uint32, oldLevel, newLevel int32) {
		this.cond.DoneTimes += newLevel - oldLevel
		this.setDirty()
		if this.cond.DoneTimes >= int32(this.condTemp.DoneTimes) {
			this.cond.DoneTimes = int32(this.condTemp.DoneTimes)
			this.Done()
		}
	}
	this.trigger(event.EventWeaponLevelUp, trigger)
	this.cond.condEvent = this
}

/***************************** 任意武器突破count次 **************************************************************************************/
// AnyWeaponBreak	任意武器突破count次	武器突破次数	5
type anyWeaponBreak struct {
	*condEventBase
}

func (this *anyWeaponBreak) Before() bool { return false }
func (this *anyWeaponBreak) Trigger() {
	trigger := func(weaponID uint32, oldLevel, newLevel int32) {
		this.cond.DoneTimes += newLevel - oldLevel
		this.setDirty()
		if this.cond.DoneTimes >= int32(this.condTemp.DoneTimes) {
			this.cond.DoneTimes = int32(this.condTemp.DoneTimes)
			this.Done()
		}
	}
	this.trigger(event.EventWeaponBreak, trigger)
	this.cond.condEvent = this
}

/***************************** 任意武器精炼count次 **************************************************************************************/
// AnyWeaponRefine	任意武器精炼count次	武器精炼次数	5
type anyWeaponRefine struct {
	*condEventBase
}

func (this *anyWeaponRefine) Before() bool { return false }
func (this *anyWeaponRefine) Trigger() {
	trigger := func(weaponID uint32, oldLevel, newLevel int32) {
		this.cond.DoneTimes += newLevel - oldLevel
		this.setDirty()
		if this.cond.DoneTimes >= int32(this.condTemp.DoneTimes) {
			this.cond.DoneTimes = int32(this.condTemp.DoneTimes)
			this.Done()
		}
	}
	this.trigger(event.EventWeaponRefine, trigger)
	this.cond.condEvent = this
}

/***************************** 拥有等级高于level的武器count件 **************************************************************************************/
// WeaponOwned	拥有等级高于level的武器count件	武器等级#武器数量	20#5
type weaponOwned struct {
	*condEventBase
}

func (this *weaponOwned) do() bool {
	doneTimes := int32(0)
	weaponModule := this.uq.userI.GetSubModule(module.Weapon).(*weapon.UserWeapon)
	weaponModule.Range(func(w *weapon.Weapon) bool {
		if w.Level >= this.condTemp.WeaponLevel {
			doneTimes++
			if doneTimes >= this.condTemp.WeaponCount {
				return false
			}
		}
		return true
	})

	if this.cond.DoneTimes != doneTimes {
		this.cond.DoneTimes = doneTimes
		this.setDirty()
		if this.cond.DoneTimes >= this.condTemp.WeaponCount {
			this.cond.DoneTimes = this.condTemp.WeaponCount
			return true
		}
	}
	return false
}

func (this *weaponOwned) Before() bool {
	return this.do()
}
func (this *weaponOwned) Trigger() {
	trigger := func(weaponID uint32, oldLevel, newLevel int32) {
		// newLevel == 0,表示被消耗掉了
		if (oldLevel >= this.condTemp.WeaponLevel && newLevel == 0) ||
			(oldLevel < this.condTemp.WeaponLevel && newLevel >= this.condTemp.WeaponLevel) {
			if this.do() {
				this.Done()
			}
		}
	}
	this.trigger(event.EventWeaponLevelUp, trigger)
	this.cond.condEvent = this
}

/***************************** 装备一件稀有度为rarity，等级为level的武器 **************************************************************************************/
// WeaponEquipped	装备一件稀有度为rarity，等级为level的武器	稀有度#武器等级	Star5#40
type weaponEquipped struct {
	*condEventBase
}

func (this *weaponEquipped) Before() bool {
	doneTimes := int32(0)
	weaponModule := this.uq.userI.GetSubModule(module.Weapon).(*weapon.UserWeapon)
	weaponModule.Range(func(w *weapon.Weapon) bool {
		def := Weapon.GetID(w.ConfigID)
		if w.EquipCharacterID != 0 && w.Level >= this.condTemp.WeaponLevel && def.RarityTypeEnum >= this.condTemp.WeaponRarity {
			doneTimes = 1
			return false
		}
		return true
	})
	if doneTimes == 1 {
		this.cond.DoneTimes = 1
		return true
	}
	return false
}
func (this *weaponEquipped) Trigger() {
	trigger := func(newLevel, rarity int32) {
		if newLevel >= this.condTemp.WeaponLevel && rarity >= this.condTemp.WeaponRarity {
			this.cond.DoneTimes = 1
			this.Done()
		}
	}
	this.trigger(event.EventWeaponEquipped, trigger)
	this.cond.condEvent = this
}

/***************************** 任意装备升级count次 **************************************************************************************/
// AnyEquipLevelUp	任意装备升级count次	装备升级次数	5
type anyEquipLevelUp struct {
	*condEventBase
}

func (this *anyEquipLevelUp) Before() bool { return false }
func (this *anyEquipLevelUp) Trigger() {
	trigger := func(equipID uint32, oldLevel, newLevel int32) {
		this.cond.DoneTimes += newLevel - oldLevel
		this.setDirty()
		if this.cond.DoneTimes >= int32(this.condTemp.DoneTimes) {
			this.cond.DoneTimes = int32(this.condTemp.DoneTimes)
			this.Done()
		}
	}
	this.trigger(event.EventEquipLevelUp, trigger)
	this.cond.condEvent = this
}

/***************************** 任意装备精炼count次 **************************************************************************************/
// AnyEquipRefine	任意装备精炼count次	装备精炼次数	5
type anyEquipRefine struct {
	*condEventBase
}

func (this *anyEquipRefine) Before() bool { return false }
func (this *anyEquipRefine) Trigger() {
	trigger := func(equipID uint32, oldLevel, newLevel int32) {
		this.cond.DoneTimes += newLevel - oldLevel
		this.setDirty()
		if this.cond.DoneTimes >= int32(this.condTemp.DoneTimes) {
			this.cond.DoneTimes = int32(this.condTemp.DoneTimes)
			this.Done()
		}
	}
	this.trigger(event.EventEquipRefine, trigger)
	this.cond.condEvent = this
}

/***************************** 拥有等级高于level的束流器count件 **************************************************************************************/
// EquipOwned	拥有等级高于level的束流器count件	束流器等级#束流器数量	20#3
type equipOwned struct {
	*condEventBase
}

func (this *equipOwned) do() bool {
	doneTimes := int32(0)
	equipModule := this.uq.userI.GetSubModule(module.Equip).(*equip.UserEquip)
	equipModule.Range(func(e *equip.Equip) bool {
		if e.Level >= this.condTemp.EquipLevel {
			doneTimes++
			if doneTimes >= this.condTemp.EquipCount {
				return false
			}
		}
		return true
	})

	if this.cond.DoneTimes != doneTimes {
		this.cond.DoneTimes = doneTimes
		this.setDirty()
		if this.cond.DoneTimes >= this.condTemp.EquipCount {
			this.cond.DoneTimes = this.condTemp.EquipCount
			return true
		}
	}
	return false
}

func (this *equipOwned) Before() bool {
	return this.do()
}
func (this *equipOwned) Trigger() {
	trigger := func(equipID uint32, oldLevel, newLevel int32) {
		// newLevel == 0,表示被消耗掉了
		if (oldLevel >= this.condTemp.EquipLevel && newLevel == 0) ||
			(oldLevel < this.condTemp.EquipLevel && newLevel >= this.condTemp.EquipLevel) {
			if this.do() {
				this.Done()
			}
		}
	}
	this.trigger(event.EventEquipLevelUp, trigger)
	this.cond.condEvent = this
}

/***************************** 装备等级高于level的束流器count件 **************************************************************************************/
// EquipEquipped	装备等级高于level的束流器count件	束流器等级#束流器数量	40#5
type equipEquipped struct {
	*condEventBase
}

func (this *equipEquipped) do() bool {
	doneTimes := int32(0)
	equipModule := this.uq.userI.GetSubModule(module.Equip).(*equip.UserEquip)
	equipModule.Range(func(e *equip.Equip) bool {
		if e.EquipCharacterId != 0 && e.Level >= this.condTemp.EquipLevel {
			doneTimes++
			if doneTimes >= this.condTemp.EquipCount {
				return false
			}
		}
		return true
	})
	if this.cond.DoneTimes != doneTimes {
		this.cond.DoneTimes = doneTimes
		this.setDirty()
		if this.cond.DoneTimes >= this.condTemp.EquipCount {
			this.cond.DoneTimes = this.condTemp.EquipCount
			return true
		}
	}
	return false
}

func (this *equipEquipped) Before() bool {
	return this.do()
}
func (this *equipEquipped) Trigger() {
	trigger := func() {
		if this.do() {
			this.Done()
		}
	}
	this.trigger(event.EventEquipEquipped, trigger)
	this.cond.condEvent = this
}

/***************************** 完成指定任务 **************************************************************************************/
// QuestComplete	完成指定任务	任务ID	1001
type questComplete struct {
	*condEventBase
}

func (this *questComplete) Before() bool {
	q := this.uq.GetQuest(this.condTemp.QuestID)
	if q != nil && q.State == message.QuestState_End {
		this.cond.DoneTimes = 1
		return true
	}
	return false
}
func (this *questComplete) Trigger() {
	trigger := func(questID int32) {
		if this.condTemp.QuestID == questID {
			this.cond.DoneTimes = 1
			this.Done()
		}
	}
	this.trigger(event.EventQuestComplete, trigger)
	this.cond.condEvent = this
}

/***************************** 回复体力 **************************************************************************************/
// FatigueRegen	回复体力	次数	1
type fatigueRegen struct {
	*condEventBase
}

func (this *fatigueRegen) Before() bool { return false }
func (this *fatigueRegen) Trigger() {
	trigger := func(times int32) {
		this.cond.DoneTimes += times
		this.setDirty()
		if this.cond.DoneTimes >= int32(this.condTemp.DoneTimes) {
			this.cond.DoneTimes = int32(this.condTemp.DoneTimes)
			this.Done()
		}
	}
	this.trigger(event.EventFatigueRegen, trigger)
	this.cond.condEvent = this
}

/***************************** 任意悬赏任务副本 **************************************************************************************/
// AnyWorldQuestInstanceSucceed	任意悬赏任务副本	通关次数	1
type anyWorldQuestInstanceSucceed struct {
	*condEventBase
}

func (this *anyWorldQuestInstanceSucceed) Before() bool { return false }
func (this *anyWorldQuestInstanceSucceed) Trigger() {
	trigger := func(instanceID, instanceStar int32, monster map[int32]int32) {
		def := Dungeon.GetID(instanceID)
		if def != nil && def.SystemTypeEnum == enumType.DungeonSystemType_WorldQuest {
			this.cond.DoneTimes++
			this.setDirty()
			if this.cond.DoneTimes >= int32(this.condTemp.DoneTimes) {
				this.cond.DoneTimes = int32(this.condTemp.DoneTimes)
				this.Done()
			}
		}
	}
	this.trigger(event.EventInstanceSucceed, trigger)
	this.cond.condEvent = this
}

/***************************** 任意商店购买道具 **************************************************************************************/
// AnyShopPurchase	任意商店购买道具	次数	1
type anyShopPurchase struct {
	*condEventBase
}

func (this *anyShopPurchase) Before() bool { return false }
func (this *anyShopPurchase) Trigger() {
	trigger := func(id, count int32) {
		this.cond.DoneTimes += count
		this.setDirty()
		if this.cond.DoneTimes >= int32(this.condTemp.DoneTimes) {
			this.cond.DoneTimes = int32(this.condTemp.DoneTimes)
			this.Done()
		}
	}
	this.trigger(event.EventShopBuy, trigger)
	this.cond.condEvent = this
}

/***************************** 通过任意支线章节 **************************************************************************************/
// BranchInstancePass	通过任意支线章节	通关支线章节数目	1	查历史进度
type branchInstancePass struct {
	*condEventBase
}

func (this *branchInstancePass) do() bool {
	mainDungeon := this.uq.userI.GetSubModule(module.MainDungeons).(*maindungeons.MainDungeons)
	doneTimes := int32(0)
	mainDungeon.ChapterRange(func(c *maindungeons.Chapter) bool {
		def := MainChapter.GetID(c.ID)
		if def != nil && def.ChapterTypeEnum == enumType.MainChapterType_Branch {
			doneTimes++
		}
		if doneTimes >= int32(this.condTemp.DoneTimes) {
			return false
		}
		return true
	})
	if this.cond.DoneTimes != doneTimes {
		this.cond.DoneTimes = doneTimes
		this.setDirty()
		if this.cond.DoneTimes >= int32(this.condTemp.DoneTimes) {
			this.cond.DoneTimes = int32(this.condTemp.DoneTimes)
			return true
		}
	}
	return false
}

func (this *branchInstancePass) Before() bool {
	return this.do()
}
func (this *branchInstancePass) Trigger() {
	trigger := func(chapterID int32) {
		if this.do() {
			this.Done()
		}
	}
	this.trigger(event.EventMainChapter, trigger)
	this.cond.condEvent = this
}

/***************************** 通过任意角色传记 **************************************************************************************/
// CharacterInstancePass	通过任意角色传记	通关传记章节数目	1	查历史进度
type characterInstancePass struct {
	*condEventBase
}

func (this *characterInstancePass) do() bool {
	mainDungeon := this.uq.userI.GetSubModule(module.MainDungeons).(*maindungeons.MainDungeons)
	doneTimes := int32(0)
	mainDungeon.ChapterRange(func(c *maindungeons.Chapter) bool {
		def := MainChapter.GetID(c.ID)
		if def != nil && def.ChapterTypeEnum == enumType.MainChapterType_Character {
			doneTimes++
		}
		if doneTimes >= int32(this.condTemp.DoneTimes) {
			return false
		}
		return true
	})
	if this.cond.DoneTimes != doneTimes {
		this.cond.DoneTimes = doneTimes
		this.setDirty()
		if this.cond.DoneTimes >= int32(this.condTemp.DoneTimes) {
			this.cond.DoneTimes = int32(this.condTemp.DoneTimes)
			return true
		}
	}
	return false
}

func (this *characterInstancePass) Before() bool {
	return this.do()
}
func (this *characterInstancePass) Trigger() {
	trigger := func(chapterID int32) {
		if this.do() {
			this.Done()
		}
	}
	this.trigger(event.EventMainChapter, trigger)
	this.cond.condEvent = this
}

/***************************** 任意count角色等级达到指定level **************************************************************************************/
// RoleLevel	任意count角色等级达到指定level	角色数量#角色等级	3#10	查历史进度
type roleLevel struct {
	*condEventBase
}

func (this *roleLevel) do() bool {
	doneTimes := int32(0)
	charaModule := this.uq.userI.GetSubModule(module.Character).(*character.UserCharacter)
	charaModule.Range(func(c *character.Character) bool {
		if c.Level >= this.condTemp.CharacterLevel {
			doneTimes++
			if doneTimes >= this.condTemp.CharacterCount {
				return false
			}
		}
		return true
	})
	if this.cond.DoneTimes != doneTimes {
		this.cond.DoneTimes = doneTimes
		this.setDirty()
		if this.cond.DoneTimes >= this.condTemp.CharacterCount {
			this.cond.DoneTimes = this.condTemp.CharacterCount
			return true
		}
	}
	return false
}

func (this *roleLevel) Before() bool {
	return this.do()
}
func (this *roleLevel) Trigger() {
	trigger := func(charaID, oldLevel, newLevel int32) {
		if oldLevel < this.condTemp.CharacterLevel && newLevel >= this.condTemp.CharacterLevel {
			if this.do() {
				this.Done()
			}
		}
	}
	this.trigger(event.EventCharacterLevelUp, trigger)
	this.cond.condEvent = this
}

/***************************** 任意角色技能升级count(所有角色累计) **************************************************************************************/
// RoleSkillLevelUp	任意角色技能升级count(所有角色累计)	角色技能升级次数	5	查历史进度
type roleSkillLevelUp struct {
	*condEventBase
}

func (this *roleSkillLevelUp) Before() bool {
	doneTimes := int32(0)
	charaModule := this.uq.userI.GetSubModule(module.Character).(*character.UserCharacter)
	charaModule.Range(func(c *character.Character) bool {
		for _, sk := range c.Skills {
			if level := sk.GetLevel(); level > 1 {
				doneTimes += level - 1
				if doneTimes >= int32(this.condTemp.DoneTimes) {
					return false
				}
			}
		}
		return true
	})
	if this.cond.DoneTimes != doneTimes {
		this.cond.DoneTimes = doneTimes
		this.setDirty()
		if this.cond.DoneTimes >= int32(this.condTemp.DoneTimes) {
			this.cond.DoneTimes = int32(this.condTemp.DoneTimes)
			return true
		}
	}
	return false
}
func (this *roleSkillLevelUp) Trigger() {
	trigger := func(charaID, skillID, oldLevel, newLevel int32) {
		this.cond.DoneTimes += newLevel - oldLevel
		this.setDirty()
		if this.cond.DoneTimes >= int32(this.condTemp.DoneTimes) {
			this.cond.DoneTimes = int32(this.condTemp.DoneTimes)
			this.Done()
		}
	}
	this.trigger(event.EventCharacterSkillLevelUp, trigger)
	this.cond.condEvent = this
}

/***************************** 任意count技能达到指定level **************************************************************************************/
// RoleSkillLevel	任意count技能达到指定level	技能数量#技能等级	5#5	查历史进度
type roleSkillLevel struct {
	*condEventBase
}

func (this *roleSkillLevel) do() bool {
	doneTimes := int32(0)
	charaModule := this.uq.userI.GetSubModule(module.Character).(*character.UserCharacter)
	charaModule.Range(func(c *character.Character) bool {
		for _, sk := range c.Skills {
			if sk.GetLevel() >= this.condTemp.RoleSkillLevel {
				doneTimes++
			}
			if doneTimes >= this.condTemp.RoleSkillCount {
				return false
			}
		}
		return true
	})
	if this.cond.DoneTimes != doneTimes {
		this.cond.DoneTimes = doneTimes
		this.setDirty()
		if this.cond.DoneTimes >= this.condTemp.RoleSkillCount {
			this.cond.DoneTimes = this.condTemp.RoleSkillCount
			return true
		}
	}
	return false
}

func (this *roleSkillLevel) Before() bool {
	return this.do()
}
func (this *roleSkillLevel) Trigger() {
	trigger := func(charaID, skillID, oldLevel, newLevel int32) {
		if oldLevel < this.condTemp.RoleSkillLevel && newLevel >= this.condTemp.RoleSkillLevel {
			if this.do() {
				this.Done()
			}
		}
	}
	this.trigger(event.EventCharacterSkillLevelUp, trigger)
	this.cond.condEvent = this
}

/***************************** 任意角色已装备束流器件数 **************************************************************************************/
// RoleEquipCount	任意角色已装备束流器件数	装备数量	5	查历史进度
type roleEquipCount struct {
	*condEventBase
}

func (this *roleEquipCount) do() bool {
	doneTimes := int32(0)
	charaModule := this.uq.userI.GetSubModule(module.Character).(*character.UserCharacter)
	charaModule.Range(func(c *character.Character) bool {
		equipCount := int32(0)
		for _, id := range c.EquipIDs {
			if id != 0 {
				equipCount++
			}
		}
		if equipCount > doneTimes {
			doneTimes = equipCount
			if doneTimes >= this.condTemp.EquipCount {
				return false
			}
		}

		return true
	})
	if this.cond.DoneTimes != doneTimes {
		this.cond.DoneTimes = doneTimes
		this.setDirty()
		if this.cond.DoneTimes >= this.condTemp.EquipCount {
			this.cond.DoneTimes = this.condTemp.EquipCount
			return true
		}
	}

	return false
}

func (this *roleEquipCount) Before() bool {
	return this.do()
}
func (this *roleEquipCount) Trigger() {
	trigger := func() {
		if this.do() {
			this.Done()
		}
	}
	this.trigger(event.EventEquipEquipped, trigger)
	this.cond.condEvent = this
}

/***************************** 通关任意秘境副本 **************************************************************************************/
// AnySecretInstanceSucceed	通关任意秘境副本	通关次数	10	需要执行
type anySecretInstanceSucceed struct {
	*condEventBase
}

func (this *anySecretInstanceSucceed) Before() bool { return false }
func (this *anySecretInstanceSucceed) Trigger() {
	trigger := func(instanceID, instanceStar int32, monster map[int32]int32) {
		def := Dungeon.GetID(instanceID)
		if def != nil && def.SystemTypeEnum == enumType.DungeonSystemType_Secret {
			this.cond.DoneTimes++
			this.setDirty()
			if this.cond.DoneTimes >= int32(this.condTemp.DoneTimes) {
				this.cond.DoneTimes = int32(this.condTemp.DoneTimes)
				this.Done()
			}
		}
	}
	this.trigger(event.EventInstanceSucceed, trigger)
	this.cond.condEvent = this
}

/***************************** 黑月探秘累计进度 **************************************************************************************/
// TrailCount	黑月探秘累计进度	累计进度	300	查历史进度
type trailCount struct {
	*condEventBase
}

func (this *trailCount) Before() bool {
	trailModule := this.uq.userI.GetSubModule(module.TrialDungeon).(*trialDungeon.TrialDungeon)
	if trailModule.GetTrialCount() != this.cond.DoneTimes {
		this.cond.DoneTimes = trailModule.GetTrialCount()
		this.setDirty()
		if this.cond.DoneTimes >= this.condTemp.TrailCount {
			this.cond.DoneTimes = this.condTemp.TrailCount
			return true
		}
	}
	return false
}
func (this *trailCount) Trigger() {
	trigger := func(oldCount, newCount int32) {
		if newCount != this.cond.DoneTimes {
			this.cond.DoneTimes = newCount
			this.setDirty()
			if this.cond.DoneTimes >= this.condTemp.TrailCount {
				this.cond.DoneTimes = this.condTemp.TrailCount
				this.Done()
			}
		}
	}
	this.trigger(event.EventTrailCount, trigger)
	this.cond.condEvent = this
}

/***************************** 通过黑月探秘涨潮期任意关卡 **************************************************************************************/
// AnyTrailInstanceSucceed	通过黑月探秘涨潮期任意关卡	通关次数	20	需要执行(可重复)
type anyTrailInstanceSucceed struct {
	*condEventBase
}

func (this *anyTrailInstanceSucceed) Before() bool { return false }
func (this *anyTrailInstanceSucceed) Trigger() {
	trigger := func(instanceID, instanceStar int32, monster map[int32]int32) {
		def := Dungeon.GetID(instanceID)
		if def != nil && def.SystemTypeEnum == enumType.DungeonSystemType_Trial {
			this.cond.DoneTimes++
			this.setDirty()
			if this.cond.DoneTimes >= int32(this.condTemp.DoneTimes) {
				this.cond.DoneTimes = int32(this.condTemp.DoneTimes)
				this.Done()
			}
		}
	}
	this.trigger(event.EventInstanceSucceed, trigger)
	this.cond.condEvent = this
}

/***************************** 通关战痕印刻副本指定难度指定次数 **************************************************************************************/
// ScarsIngrainDifficultInstanceSucceed	通关战痕印刻副本指定难度指定次数	难度枚举#通关次数	5#10	需要执行(不可重复)
type scarsIngrainDifficultInstanceSucceed struct {
	*condEventBase
}

func (this *scarsIngrainDifficultInstanceSucceed) do() bool {
	doneTimes := int32(0)
	scarsModule := this.uq.userI.GetSubModule(module.ScarsIngrain).(*scarsIngrain.ScarsIngrain)
	for _, bg := range scarsModule.GetData().BossGroup {
		_, _, score := scarsModule.BossDifficultScore(bg.GetBossID(), this.condTemp.ScarsDifficult)
		if score > 0 {
			doneTimes++
		}
	}
	if this.cond.DoneTimes != doneTimes {
		this.cond.DoneTimes = doneTimes
		this.setDirty()
		if this.cond.DoneTimes >= this.condTemp.ScarsPassTimes {
			this.cond.DoneTimes = this.condTemp.ScarsPassTimes
			return true
		}
	}
	return false
}

func (this *scarsIngrainDifficultInstanceSucceed) Before() bool {
	return this.do()
}
func (this *scarsIngrainDifficultInstanceSucceed) Trigger() {
	trigger := func(instanceID, instanceStar int32, monster map[int32]int32) {
		def := Dungeon.GetID(instanceID)
		if def != nil && def.SystemTypeEnum == enumType.DungeonSystemType_ScarsIngrain {
			scarsDef := ScarsIngrainBossInstance.GetID(def.SystemConfig)
			_, difficult := scarsDef.BossDifficult()
			if difficult >= this.condTemp.ScarsDifficult && this.do() {
				this.Done()
			}
		}
	}
	this.trigger(event.EventInstanceSucceed, trigger)
	this.cond.condEvent = this
}

/***************************** 战痕印刻积分达到指定分数 **************************************************************************************/
// ScarsIngrainScore	战痕印刻积分达到指定分数	积分	10000	需要执行
type scarsIngrainScore struct {
	*condEventBase
}

func (this *scarsIngrainScore) do() bool {
	scarsModule := this.uq.userI.GetSubModule(module.ScarsIngrain).(*scarsIngrain.ScarsIngrain)
	doneTimes := scarsModule.GetData().TotalScore
	if this.cond.DoneTimes != doneTimes {
		this.cond.DoneTimes = doneTimes
		this.setDirty()
		if this.cond.DoneTimes >= this.condTemp.ScarsScore {
			this.cond.DoneTimes = this.condTemp.ScarsScore
			return true
		}
	}
	return false
}

func (this *scarsIngrainScore) Before() bool {
	return this.do()
}
func (this *scarsIngrainScore) Trigger() {
	trigger := func(instanceID, instanceStar int32, monster map[int32]int32) {
		def := Dungeon.GetID(instanceID)
		if def != nil && def.SystemTypeEnum == enumType.DungeonSystemType_ScarsIngrain {
			if this.do() {
				this.Done()
			}
		}
	}
	this.trigger(event.EventInstanceSucceed, trigger)
	this.cond.condEvent = this
}

/***************************** 任意卡池抽卡count **************************************************************************************/
// AnyDrawcardTimes	任意卡池抽卡count	抽卡次数	10	需要执行
type anyDrawcardTimes struct {
	*condEventBase
}

func (this *anyDrawcardTimes) Before() bool { return false }
func (this *anyDrawcardTimes) Trigger() {
	trigger := func(id, count int32) {
		this.cond.DoneTimes += count
		this.setDirty()
		if this.cond.DoneTimes >= int32(this.condTemp.DoneTimes) {
			this.cond.DoneTimes = int32(this.condTemp.DoneTimes)
			this.Done()
		}
	}
	this.trigger(event.EventDrawCard, trigger)
	this.cond.condEvent = this
}

/***************************** 本周累积登录count **************************************************************************************/
// WeeklyLogin	本周累积登录count	登录次数	5	查历史进度
type weeklyLogin struct {
	*condEventBase
}

func (this *weeklyLogin) Check() bool {
	attrModule := this.uq.userI.GetSubModule(module.Attr).(*attr2.UserAttr)
	return attrModule.TimeEqual(module.WeeklyTimeName, this.uq.timedata[module.WeeklyTimeName])
}

func (this *weeklyLogin) Before() bool {
	attrModule := this.uq.userI.GetSubModule(module.Attr).(*attr2.UserAttr)
	val, _ := attrModule.GetAttr(attr.WeeklyLogin)
	if this.cond.DoneTimes != int32(val) {
		this.cond.DoneTimes = int32(val)
		this.setDirty()
		if this.cond.DoneTimes >= int32(this.condTemp.LoginTimes) {
			this.cond.DoneTimes = int32(this.condTemp.LoginTimes)
			return true
		}
	}
	return false
}
func (this *weeklyLogin) Trigger() {
	trigger := func(id int32, oldVal, newVal int64) {
		if id == attr.WeeklyLogin && this.cond.DoneTimes != int32(newVal) {
			this.cond.DoneTimes = int32(newVal)
			this.setDirty()
			if this.cond.DoneTimes >= int32(this.condTemp.LoginTimes) {
				this.cond.DoneTimes = int32(this.condTemp.LoginTimes)
				this.Done()
			}
		}
	}
	this.trigger(event.EventAttrChange, trigger)
	this.cond.condEvent = this
}

/***************************** 本周累积每日活跃度 **************************************************************************************/
// WeeklyActiveness	本周累积每日活跃度	本周每日活跃度之和	500	查历史进度
type weeklyActiveness struct {
	*condEventBase
}

func (this *weeklyActiveness) Check() bool {
	attrModule := this.uq.userI.GetSubModule(module.Attr).(*attr2.UserAttr)
	return attrModule.TimeEqual(module.WeeklyTimeName, this.uq.timedata[module.WeeklyTimeName])
}

func (this *weeklyActiveness) Before() bool {
	attrModule := this.uq.userI.GetSubModule(module.Attr).(*attr2.UserAttr)
	val, _ := attrModule.GetAttr(attr.WeeklyActiveness)
	if this.cond.DoneTimes != int32(val) {
		this.cond.DoneTimes = int32(val)
		this.setDirty()
		if this.cond.DoneTimes >= int32(this.condTemp.ActivenessCount) {
			this.cond.DoneTimes = int32(this.condTemp.ActivenessCount)
			return true
		}
	}
	return false
}
func (this *weeklyActiveness) Trigger() {
	trigger := func(id int32, oldVal, newVal int64) {
		if id == attr.WeeklyActiveness && this.cond.DoneTimes != int32(newVal) {
			this.cond.DoneTimes = int32(newVal)
			this.setDirty()
			if this.cond.DoneTimes >= int32(this.condTemp.ActivenessCount) {
				this.cond.DoneTimes = int32(this.condTemp.ActivenessCount)
				this.Done()
			}
		}
	}
	this.trigger(event.EventAttrChange, trigger)
	this.cond.condEvent = this
}

/***************************** 拥有稀有度rarity以上的装备count件 **************************************************************************************/
// EquipRarityOwned	拥有稀有度rarity以上的装备count件	稀有度#数量	Star5#4	查历史进度
type equipRarityOwned struct {
	*condEventBase
}

func (this *equipRarityOwned) do() bool {
	doneTimes := int32(0)
	equipModule := this.uq.userI.GetSubModule(module.Equip).(*equip.UserEquip)
	equipModule.Range(func(e *equip.Equip) bool {
		def := Equip.GetID(e.ConfigID)
		if def.QualityEnum >= this.condTemp.EquipRarity {
			doneTimes++
			if doneTimes >= this.condTemp.EquipCount {
				return false
			}
		}
		return true
	})

	if this.cond.DoneTimes != doneTimes {
		this.cond.DoneTimes = doneTimes
		this.setDirty()
		if this.cond.DoneTimes >= this.condTemp.EquipCount {
			this.cond.DoneTimes = this.condTemp.EquipCount
			return true
		}
	}
	return false
}

func (this *equipRarityOwned) Before() bool {
	return this.do()
}
func (this *equipRarityOwned) Trigger() {
	trigger := func(equipID uint32) {
		if this.do() {
			this.Done()
		}
	}
	this.trigger(event.EventEquipAdd, trigger)
	this.cond.condEvent = this
}

/***************************** 拥有稀有度rarity以上的武器count件 **************************************************************************************/
// WeaponRarityOwned	拥有稀有度rarity以上的武器count件	稀有度#数量	Star5#4	查历史进度
type weaponRarityOwned struct {
	*condEventBase
}

func (this *weaponRarityOwned) do() bool {
	doneTimes := int32(0)
	weaponModule := this.uq.userI.GetSubModule(module.Weapon).(*weapon.UserWeapon)
	weaponModule.Range(func(w *weapon.Weapon) bool {
		def := Weapon.GetID(w.ConfigID)
		if def.RarityTypeEnum >= this.condTemp.WeaponRarity {
			doneTimes++
			if doneTimes >= this.condTemp.WeaponCount {
				return false
			}
		}
		return true
	})

	if this.cond.DoneTimes != doneTimes {
		this.cond.DoneTimes = doneTimes
		this.setDirty()
		if this.cond.DoneTimes >= this.condTemp.WeaponCount {
			this.cond.DoneTimes = this.condTemp.WeaponCount
			return true
		}
	}
	return false
}

func (this *weaponRarityOwned) Before() bool {
	return this.do()
}
func (this *weaponRarityOwned) Trigger() {
	trigger := func(weapon uint32) {
		if this.do() {
			this.Done()
		}
	}
	this.trigger(event.EventWeaponAdd, trigger)
	this.cond.condEvent = this
}

/***************************** 大秘境成功挑战次数 **************************************************************************************/
// BigSecretPassTimes	大秘境成功挑战次数	次数	5	需要执行
type bigSecretPassTimes struct {
	*condEventBase
}

func (this *bigSecretPassTimes) Before() bool { return false }
func (this *bigSecretPassTimes) Trigger() {
	trigger := func(instanceID, instanceStar int32, monster map[int32]int32) {
		def := Dungeon.GetID(instanceID)
		if def != nil && def.SystemTypeEnum == enumType.DungeonSystemType_BigSecret {
			this.cond.DoneTimes++
			this.setDirty()
			if this.cond.DoneTimes >= int32(this.condTemp.DoneTimes) {
				this.cond.DoneTimes = int32(this.condTemp.DoneTimes)
				this.Done()
			}
		}
	}
	this.trigger(event.EventInstanceSucceed, trigger)
	this.cond.condEvent = this
}

/***************************** 大秘境通关最高层级 **************************************************************************************/
// BigSecretBestScore	大秘境通关最高层级	最高层级	100	查历史进度
type bigSecretBestScore struct {
	*condEventBase
}

func (this *bigSecretBestScore) Before() bool {
	m := this.uq.userI.GetSubModule(module.BigSecret).(*bigSecret.BigSecretDungeon)
	lv := m.BestPassed()
	if lv != this.cond.DoneTimes {
		this.cond.DoneTimes = lv
		this.setDirty()
		if this.cond.DoneTimes >= int32(this.condTemp.DoneTimes) {
			this.cond.DoneTimes = int32(this.condTemp.DoneTimes)
			return true
		}
	}
	return false
}
func (this *bigSecretBestScore) Trigger() {
	trigger := func(instanceID, instanceStar int32, monster map[int32]int32) {
		def := Dungeon.GetID(instanceID)
		if def != nil && def.SystemTypeEnum == enumType.DungeonSystemType_BigSecret {
			if def.SystemConfig > this.cond.DoneTimes {
				this.cond.DoneTimes = def.SystemConfig
				this.setDirty()
				if this.cond.DoneTimes >= int32(this.condTemp.DoneTimes) {
					this.cond.DoneTimes = int32(this.condTemp.DoneTimes)
					this.Done()
				}
			}
		}
	}
	this.trigger(event.EventInstanceSucceed, trigger)
	this.cond.condEvent = this
}

/***************************** end **************************************************************************************/

func makeCondEvent(cond *Condition, condTemp *Quest2.Condition, q *Quest, u *UserQuest) CondEvent {
	creator := condEventCreator[condTemp.Type]
	if creator == nil {
		panic(fmt.Sprintf("quest condition type %d creator is nil", condTemp.Type))
	}
	base := newCondEventBase(cond, condTemp, q, u)
	return creator(base)
}

func (this *UserQuest) registerEvent(q *Quest) {
	//log.GetLogger().Debugln("registerEvent", q.ID, q.State, q.registered)
	if q.State == message.QuestState_Running && !q.registered {
		cfg := q.GetConfig()
		conds := cfg.Conditions()

		condEvents := make([]CondEvent, 0, len(q.Conditions))
		for i, c := range q.Conditions {
			if !c.Complete {
				condTemp := conds[i]
				cond := c
				condEvent := makeCondEvent(cond, condTemp, q, this)
				if !condEvent.Check() {
					// zaplogger.GetSugar().Errorf("%s register quest %d condition %d check failed", this.userI.GetUserID(), q.ID, i)
					return
				}
				condEvents = append(condEvents, condEvent)
			}
		}

		q.registered = true
		delete(this.needRegister, q.ID)

		for _, condEvent := range condEvents {
			if condEvent.Before() {
				condEvent.Done()
			} else {
				condEvent.Trigger()
			}
		}

	} else {
		delete(this.needRegister, q.ID)
	}
}

var condEventCreator = map[int32]func(base *condEventBase) CondEvent{}

func init() {
	condEventCreator = map[int32]func(base *condEventBase) CondEvent{
		enumType.QuestConditionType_DailyLogin:                           func(base *condEventBase) CondEvent { return &dailyLogin{condEventBase: base} },
		enumType.QuestConditionType_AccumulateLogin:                      func(base *condEventBase) CondEvent { return &accumulateLogin{condEventBase: base} },
		enumType.QuestConditionType_OnlineTime:                           func(base *condEventBase) CondEvent { return &onlineTime{condEventBase: base} },
		enumType.QuestConditionType_ConsumeResource:                      func(base *condEventBase) CondEvent { return &consumeResource{condEventBase: base} },
		enumType.QuestConditionType_YuruInteract:                         func(base *condEventBase) CondEvent { return &yuruInteract{condEventBase: base} },
		enumType.QuestConditionType_DailyActiveness:                      func(base *condEventBase) CondEvent { return &dailyActiveness{condEventBase: base} },
		enumType.QuestConditionType_EnterDorm:                            func(base *condEventBase) CondEvent { return &enterDorm{condEventBase: base} },
		enumType.QuestConditionType_ExchangeFatigue:                      func(base *condEventBase) CondEvent { return &exchangeFatigue{condEventBase: base} },
		enumType.QuestConditionType_PlayerLevelUp:                        func(base *condEventBase) CondEvent { return &playerLevelUp{condEventBase: base} },
		enumType.QuestConditionType_AnyInstanceSucceed:                   func(base *condEventBase) CondEvent { return &anyInstanceSucceed{condEventBase: base} },
		enumType.QuestConditionType_InstanceSucceed:                      func(base *condEventBase) CondEvent { return &instanceSucceed{condEventBase: base} },
		enumType.QuestConditionType_AnyMaterialInstanceSucceed:           func(base *condEventBase) CondEvent { return &anyMaterialInstanceSucceed{condEventBase: base} },
		enumType.QuestConditionType_AnyScarsInstanceSucceed:              func(base *condEventBase) CondEvent { return &anyScarsInstanceSucceed{condEventBase: base} },
		enumType.QuestConditionType_AnyRoleLevelUp:                       func(base *condEventBase) CondEvent { return &anyRoleLevelUp{condEventBase: base} },
		enumType.QuestConditionType_AnyRoleBreak:                         func(base *condEventBase) CondEvent { return &anyRoleBreak{condEventBase: base} },
		enumType.QuestConditionType_RoleBreak:                            func(base *condEventBase) CondEvent { return &roleBreak{condEventBase: base} },
		enumType.QuestConditionType_RoleLevelUp:                          func(base *condEventBase) CondEvent { return &roleLevelUp{condEventBase: base} },
		enumType.QuestConditionType_AnyWeaponLevelUp:                     func(base *condEventBase) CondEvent { return &anyWeaponLevelUp{condEventBase: base} },
		enumType.QuestConditionType_AnyWeaponBreak:                       func(base *condEventBase) CondEvent { return &anyWeaponBreak{condEventBase: base} },
		enumType.QuestConditionType_AnyWeaponRefine:                      func(base *condEventBase) CondEvent { return &anyWeaponRefine{condEventBase: base} },
		enumType.QuestConditionType_WeaponOwned:                          func(base *condEventBase) CondEvent { return &weaponOwned{condEventBase: base} },
		enumType.QuestConditionType_WeaponEquipped:                       func(base *condEventBase) CondEvent { return &weaponEquipped{condEventBase: base} },
		enumType.QuestConditionType_AnyEquipLevelUp:                      func(base *condEventBase) CondEvent { return &anyEquipLevelUp{condEventBase: base} },
		enumType.QuestConditionType_AnyEquipRefine:                       func(base *condEventBase) CondEvent { return &anyEquipRefine{condEventBase: base} },
		enumType.QuestConditionType_EquipOwned:                           func(base *condEventBase) CondEvent { return &equipOwned{condEventBase: base} },
		enumType.QuestConditionType_EquipEquipped:                        func(base *condEventBase) CondEvent { return &equipEquipped{condEventBase: base} },
		enumType.QuestConditionType_QuestComplete:                        func(base *condEventBase) CondEvent { return &questComplete{condEventBase: base} },
		enumType.QuestConditionType_FatigueRegen:                         func(base *condEventBase) CondEvent { return &fatigueRegen{condEventBase: base} },
		enumType.QuestConditionType_AnyWorldQuestInstanceSucceed:         func(base *condEventBase) CondEvent { return &anyWorldQuestInstanceSucceed{condEventBase: base} },
		enumType.QuestConditionType_AnyShopPurchase:                      func(base *condEventBase) CondEvent { return &anyShopPurchase{condEventBase: base} },
		enumType.QuestConditionType_BranchInstancePass:                   func(base *condEventBase) CondEvent { return &branchInstancePass{condEventBase: base} },
		enumType.QuestConditionType_CharacterInstancePass:                func(base *condEventBase) CondEvent { return &characterInstancePass{condEventBase: base} },
		enumType.QuestConditionType_RoleLevel:                            func(base *condEventBase) CondEvent { return &roleLevel{condEventBase: base} },
		enumType.QuestConditionType_RoleSkillLevelUp:                     func(base *condEventBase) CondEvent { return &roleSkillLevelUp{condEventBase: base} },
		enumType.QuestConditionType_RoleSkillLevel:                       func(base *condEventBase) CondEvent { return &roleSkillLevel{condEventBase: base} },
		enumType.QuestConditionType_RoleEquipCount:                       func(base *condEventBase) CondEvent { return &roleEquipCount{condEventBase: base} },
		enumType.QuestConditionType_AnySecretInstanceSucceed:             func(base *condEventBase) CondEvent { return &anySecretInstanceSucceed{condEventBase: base} },
		enumType.QuestConditionType_TrailCount:                           func(base *condEventBase) CondEvent { return &trailCount{condEventBase: base} },
		enumType.QuestConditionType_AnyTrailInstanceSucceed:              func(base *condEventBase) CondEvent { return &anyTrailInstanceSucceed{condEventBase: base} },
		enumType.QuestConditionType_ScarsIngrainDifficultInstanceSucceed: func(base *condEventBase) CondEvent { return &scarsIngrainDifficultInstanceSucceed{condEventBase: base} },
		enumType.QuestConditionType_ScarsIngrainScore:                    func(base *condEventBase) CondEvent { return &scarsIngrainScore{condEventBase: base} },
		enumType.QuestConditionType_AnyDrawcardTimes:                     func(base *condEventBase) CondEvent { return &anyDrawcardTimes{condEventBase: base} },
		enumType.QuestConditionType_WeeklyLogin:                          func(base *condEventBase) CondEvent { return &weeklyLogin{condEventBase: base} },
		enumType.QuestConditionType_WeeklyActiveness:                     func(base *condEventBase) CondEvent { return &weeklyActiveness{condEventBase: base} },
		enumType.QuestConditionType_EquipRarityOwned:                     func(base *condEventBase) CondEvent { return &equipRarityOwned{condEventBase: base} },
		enumType.QuestConditionType_WeaponRarityOwned:                    func(base *condEventBase) CondEvent { return &weaponRarityOwned{condEventBase: base} },
		enumType.QuestConditionType_BigSecretBestScore:                   func(base *condEventBase) CondEvent { return &bigSecretBestScore{condEventBase: base} },
		enumType.QuestConditionType_BigSecretPassTimes:                   func(base *condEventBase) CondEvent { return &bigSecretPassTimes{condEventBase: base} },
	}

}
