package battleAtt

import (
	"initialthree/node/common/battleAttr"
)

var (
	additionAttrs = map[int32][]int32{
		battleAttr.Attack:   {battleAttr.AttackPer},
		battleAttr.Defense:  {battleAttr.DefensePer},
		battleAttr.HitPoint: {battleAttr.HitPointPer},
	}

	ConstAttrs = map[int32]float64{
		//battleAttr.EXEnergyRecharge:    1,
		//battleAttr.SkillEnergyRecharge: 1,
		battleAttr.Vulnerable: 1,
	}

	// 特殊属性, 计算时用（1 + 物理增伤）
	AddAttrs = map[int32]struct{}{
		battleAttr.DamageAmplification: {},
		//battleAttr.PhysicsAmplification: {}, //物理增伤
		//battleAttr.IceAmplification:     {}, //冰系增伤
		//battleAttr.Fir:                  {}, //火系增伤
		//battleAttr.ThunderAmplification: {}, //雷系增伤
		//battleAttr.DarkAmplification:    {}, //暗系增伤
		//battleAttr.EXEnergyRecharge:    {},
		//battleAttr.SkillEnergyRecharge: {},
	}

	// 特殊属性, 计算时用（1 - 物理免疫）
	SubAttrs = map[int32]struct{}{
		battleAttr.DamageReduction: {},
	}

	// 物伤m = 伤害 d * 物理系数 * 物增伤 * （1 - 物抗） * （1 - 物免）
	// 这些属性，同步到客户端时，由于要经过公式（1 - 物抗），这里同步最终值需要用 1 - x。
	// 注意：服务器做战斗计算时，不用公式中的1- x,直接用改值参与计算。
	FinalAttrs = map[int32]struct{}{
		//battleAttr.PhysicsImmunity:   {}, //物理免疫
		//battleAttr.IceImmunity:       {}, //冰系免疫
		//battleAttr.FireImmunity:      {}, //火系免疫
		//battleAttr.ThunderImmunity:   {}, //雷系免疫
		//battleAttr.DarkImmunity:      {}, //暗系免疫
		//battleAttr.PhysicsResistance: {}, //物理抗性
		//battleAttr.IceResistance:     {}, //冰系抗性
		//battleAttr.DarkResistance:    {}, //暗系抗性
	}
)

// 系统内属性计算完成后，调用此接口将特殊属性转换
func transAttr(kvMap map[int32]float64) {
	for id := range AddAttrs {
		v := kvMap[id]
		kvMap[id] = 1 + v
	}
	for id := range SubAttrs {
		v := kvMap[id]
		kvMap[id] = 1 - v
	}
}

// 最终处理
func finalAttr(kvMap map[int32]float64) {
	for id := range FinalAttrs {
		v := kvMap[id]
		kvMap[id] = 1 - v
	}

	kvMap[battleAttr.HitPointMax] = kvMap[battleAttr.HitPoint]
}

// 系统间做计算，
func systemAttrCalculate(kvMap ...map[int32]float64) map[int32]float64 {
	length := len(kvMap)
	if length == 0 {
		return newTransAttr()
	}

	ret := kvMap[0]
	for i := 1; i < length; i++ {
		ret = systemTwoCal(ret, kvMap[i])
	}
	return ret
}

// 两个系统间计算,部分直接累加，部分相乘
func systemTwoCal(kvMap, source map[int32]float64) map[int32]float64 {
	// 根据战斗公式计算
	for id, val := range source {
		if _, ok := AddAttrs[id]; ok {
			multiplyAttr(kvMap, id, val)
			continue
		}
		if _, ok := SubAttrs[id]; ok {
			multiplyAttr(kvMap, id, val)
			continue
		}
		addAttrs(kvMap, id, val)
	}
	return kvMap
}

// 加
func addAttrs(kvMap map[int32]float64, id int32, val float64) {
	if v, ok := kvMap[id]; ok {
		kvMap[id] = val + v
	} else {
		kvMap[id] = val
	}
}

// 乘
func multiplyAttr(kvMap map[int32]float64, id int32, val float64) {
	if v, ok := kvMap[id]; ok {
		kvMap[id] = val * v
	} else {
		kvMap[id] = val
	}
}

// 返回规则后的属性值
func newTransAttr() map[int32]float64 {
	ret := map[int32]float64{}
	for id, v := range ConstAttrs {
		ret[id] = v
	}

	transAttr(ret)
	return ret
}

// 计算战斗属性用角色
type BattleCharacter struct {
	UserID      string
	GameID      uint64
	CharacterID int32
	Level       int32
	BreakLevel  int32
	GeneLevel   int32
	Weapon      *BattleWeapon
	Equips      []*BattleEquip
}

type BattleWeapon struct {
	ConfigID   int32
	Level      int32
	Refine     int32
	BreakTimes int32
}

// 装备
type BattleEquip struct {
	ConfigID     int32
	Level        int32
	RandomAttrID int32
	Refine       []int32
}
