package event

/*、
参1：事件类型， 参2：参数
func (this *User) EmitEvent(event interface{}, args ...interface{}) {
	this.evHandler.Emit(event, args...)
}*/

/*
 定义事件类型，整型默认使用 int32。其他情况指明类型
 新增参数，末尾添加。不影响原有逻辑
*/

type EventType int

const (
	EventAttrChange            EventType = iota // 属性变化事件，参数2（属性ID; oldVal:int64; newVal:int64）
	EventInstanceSucceed                        // 关卡完成，参数2（关卡ID; 通关星级; 击杀怪物ID、数量:map[int32]int32 ）
	EventQuestComplete                          // 任务完成，参数2（任务ID）
	EventCharacterLevelUp                       // 角色升级，参数2（角色ID; oldLevel; newLevel）
	EventCharacterGeneLevelUp                   // 角色命座升级，参数2（角色ID; oldLevel; newLevel）
	EventCharacterBreak                         // 角色突破，参数2（角色ID; oldLevel; newLevel）
	EventUseItem                                // 道具资源消耗，参数2（资源ID、数量:map[int32]int32 ） // 特指道具 item 。（金币、砖石等由属性变化触发）
	EventExchangeFatigue                        // 属性兑换体力，参数2（兑换次数:int32）
	EventWeaponLevelUp                          // 武器升级，参数2（weapon:uint32; oldLevel; newLevel）
	EventWeaponBreak                            // 武器突破，参数2（weapon:uint32; oldLevel; newLevel）
	EventWeaponRefine                           // 武器精炼，参数2（weapon:uint32; oldLevel; newLevel）
	EventEquipLevelUp                           // 装备精炼，参数2（equipID:uint32; oldLevel; newLevel）
	EventEquipRefine                            // 装备精炼，参数2（equipID:uint32; oldLevel; newLevel）
	EventEquipEquipped                          // 装备，无参数
	EventWeaponEquipped                         // 装备，参数2（newLevel; rarity）
	EventFatigueRegen                           // 使用道具恢复体力，参数2(times;)
	EventShopBuy                                // 商城购买，参数2(ID;count)
	EventMainChapter                            // 主线章节完成，参数2(chapterID)
	EventCharacterSkillLevelUp                  // 角色技能升级，参数2(characterID;skillID; oldLevel; newLevel)
	EventTrailCount                             // 秘境点数，参数2( oldCount; newCount)
	EventDrawCard                               // 抽卡，参数2(ID; count)
	EventEquipAdd                               // 装备获取，参数2(ID:uint32)
	EventWeaponAdd                              // 武器获取，参数2(ID:uint32)
	EventBigSecret                              // 大秘境重置，参数2(ID:uint32)
)
