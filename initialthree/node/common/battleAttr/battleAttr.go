
package battleAttr

type BattleAttr struct {
	ID  int32
	Val float64
}

type Info struct{
	Min         float64
	Max         float64
}

var idxToName map[int32]string
var nameToId  map[string]int32
var battleAttrInfo map[int32]*Info

const(
	PlaceHolder_1 = 1 //预留占位1
	PlaceHolder_2 = 2 //预留占位2
	PlaceHolder_3 = 3 //预留占位3
	PlaceHolder_4 = 4 //预留占位4
	PlaceHolder_5 = 5 //预留占位5
	PlaceHolder_6 = 6 //预留占位6
	PlaceHolder_7 = 7 //预留占位7
	PlaceHolder_8 = 8 //预留占位8
	PlaceHolder_9 = 9 //预留占位9
	PlaceHolder_10 = 10 //预留占位10
	Attack = 11 //攻击
	AttackPer = 12 //攻击百分比
	Defense = 13 //防御
	DefensePer = 14 //防御百分比
	HitPoint = 15 //生命
	HitPointPer = 16 //生命百分比
	PlaceHolder_17 = 17 //预留占位17
	CriticalChance = 18 //暴击率
	CriticalAmplification = 19 //暴击率增幅
	CriticalDamage = 20 //暴击伤害
	DamageAmplification = 21 //增伤
	DamageReduction = 22 //免伤
	HitPointMax = 23 //生命值上限
	ShieldMax = 24 //护盾值上限
	SuperArmorMax = 25 //霸体值上限
	BeheadDegeneration = 26 //预留占位26
	EXEnergy = 27 //大招能量值
	EXEnergyMax = 28 //大招能量值上限
	EXEnergyRechargeBackstage = 29 //大招后台充能效率
	ManaRegeneration = 30 //预留占位30
	Shield = 31 //护盾值
	ShieldRecoverTime = 32 //护盾恢复时间
	SuperArmor = 33 //霸体值
	SuperArmorRecoverTime = 34 //霸体值回复时间
	Vulnerable = 35 //易伤
	MentorEnergy = 36 //导师能量
	MentorEnergyRecharge = 37 //导师能量恢复速率
	BeheadMax = 38 //预留占位38
	BrokenShieldStrength = 39 //破盾强度
	SuperArmorRemove = 40 //破霸体强度
	SuperArmorPrevent = 41 //霸体强度
	EXEnergyRecharge = 42 //大招充能效率
	SkillEnergyRecharge = 43 //技能充能效率
	ShieldStrength = 44 //护盾强度
	RecoveryHPStrength = 45 //治疗强度
	ElementEffectAddition = 46 //元素附着增强
	Vigour = 47 //精力值
	VigourRecoveryRate = 48 //精力恢复速度
	VigourMax = 49 //精力值上限
	PlaceHolder_50 = 50 //预留占位50
	MatterWeakness = 51 //物质-弱点
	EnergyWeakness = 52 //异能-弱点
	ViodWeakness = 53 //虚空-弱点
	PlaceHolder_54 = 54 //预留占位54
	MatterResistance = 55 //物质-抗性
	EnergyResistance = 56 //异能-抗性
	ViodResistance = 57 //虚空-抗性
	AttackAmplificationCoefficient = 58 //进攻战力增幅系数
	PlaceHolder_59 = 59 //预留占位59
	PlaceHolder_60 = 60 //预留占位60
	PlaceHolder_61 = 61 //预留占位61
	PlaceHolder_62 = 62 //预留占位62
	PlaceHolder_63 = 63 //预留占位63
	PlaceHolder_64 = 64 //预留占位64
	PlaceHolder_65 = 65 //预留占位65
	PlaceHolder_66 = 66 //预留占位66
	PlaceHolder_67 = 67 //预留占位67
	PlaceHolder_68 = 68 //预留占位68
	PlaceHolder_69 = 69 //预留占位69
	PlaceHolder_70 = 70 //预留占位70
	AttrMax = 70
)

func init() {	
	idxToName = map[int32]string{}
	nameToId  = map[string]int32{}
	battleAttrInfo = map[int32]*Info{}
	idxToName[1] = "PlaceHolder_1"
	idxToName[2] = "PlaceHolder_2"
	idxToName[3] = "PlaceHolder_3"
	idxToName[4] = "PlaceHolder_4"
	idxToName[5] = "PlaceHolder_5"
	idxToName[6] = "PlaceHolder_6"
	idxToName[7] = "PlaceHolder_7"
	idxToName[8] = "PlaceHolder_8"
	idxToName[9] = "PlaceHolder_9"
	idxToName[10] = "PlaceHolder_10"
	idxToName[11] = "Attack"
	idxToName[12] = "AttackPer"
	idxToName[13] = "Defense"
	idxToName[14] = "DefensePer"
	idxToName[15] = "HitPoint"
	idxToName[16] = "HitPointPer"
	idxToName[17] = "PlaceHolder_17"
	idxToName[18] = "CriticalChance"
	idxToName[19] = "CriticalAmplification"
	idxToName[20] = "CriticalDamage"
	idxToName[21] = "DamageAmplification"
	idxToName[22] = "DamageReduction"
	idxToName[23] = "HitPointMax"
	idxToName[24] = "ShieldMax"
	idxToName[25] = "SuperArmorMax"
	idxToName[26] = "BeheadDegeneration"
	idxToName[27] = "EXEnergy"
	idxToName[28] = "EXEnergyMax"
	idxToName[29] = "EXEnergyRechargeBackstage"
	idxToName[30] = "ManaRegeneration"
	idxToName[31] = "Shield"
	idxToName[32] = "ShieldRecoverTime"
	idxToName[33] = "SuperArmor"
	idxToName[34] = "SuperArmorRecoverTime"
	idxToName[35] = "Vulnerable"
	idxToName[36] = "MentorEnergy"
	idxToName[37] = "MentorEnergyRecharge"
	idxToName[38] = "BeheadMax"
	idxToName[39] = "BrokenShieldStrength"
	idxToName[40] = "SuperArmorRemove"
	idxToName[41] = "SuperArmorPrevent"
	idxToName[42] = "EXEnergyRecharge"
	idxToName[43] = "SkillEnergyRecharge"
	idxToName[44] = "ShieldStrength"
	idxToName[45] = "RecoveryHPStrength"
	idxToName[46] = "ElementEffectAddition"
	idxToName[47] = "Vigour"
	idxToName[48] = "VigourRecoveryRate"
	idxToName[49] = "VigourMax"
	idxToName[50] = "PlaceHolder_50"
	idxToName[51] = "MatterWeakness"
	idxToName[52] = "EnergyWeakness"
	idxToName[53] = "ViodWeakness"
	idxToName[54] = "PlaceHolder_54"
	idxToName[55] = "MatterResistance"
	idxToName[56] = "EnergyResistance"
	idxToName[57] = "ViodResistance"
	idxToName[58] = "AttackAmplificationCoefficient"
	idxToName[59] = "PlaceHolder_59"
	idxToName[60] = "PlaceHolder_60"
	idxToName[61] = "PlaceHolder_61"
	idxToName[62] = "PlaceHolder_62"
	idxToName[63] = "PlaceHolder_63"
	idxToName[64] = "PlaceHolder_64"
	idxToName[65] = "PlaceHolder_65"
	idxToName[66] = "PlaceHolder_66"
	idxToName[67] = "PlaceHolder_67"
	idxToName[68] = "PlaceHolder_68"
	idxToName[69] = "PlaceHolder_69"
	idxToName[70] = "PlaceHolder_70"
	nameToId["PlaceHolder_1"] = 1
	nameToId["PlaceHolder_2"] = 2
	nameToId["PlaceHolder_3"] = 3
	nameToId["PlaceHolder_4"] = 4
	nameToId["PlaceHolder_5"] = 5
	nameToId["PlaceHolder_6"] = 6
	nameToId["PlaceHolder_7"] = 7
	nameToId["PlaceHolder_8"] = 8
	nameToId["PlaceHolder_9"] = 9
	nameToId["PlaceHolder_10"] = 10
	nameToId["Attack"] = 11
	nameToId["AttackPer"] = 12
	nameToId["Defense"] = 13
	nameToId["DefensePer"] = 14
	nameToId["HitPoint"] = 15
	nameToId["HitPointPer"] = 16
	nameToId["PlaceHolder_17"] = 17
	nameToId["CriticalChance"] = 18
	nameToId["CriticalAmplification"] = 19
	nameToId["CriticalDamage"] = 20
	nameToId["DamageAmplification"] = 21
	nameToId["DamageReduction"] = 22
	nameToId["HitPointMax"] = 23
	nameToId["ShieldMax"] = 24
	nameToId["SuperArmorMax"] = 25
	nameToId["BeheadDegeneration"] = 26
	nameToId["EXEnergy"] = 27
	nameToId["EXEnergyMax"] = 28
	nameToId["EXEnergyRechargeBackstage"] = 29
	nameToId["ManaRegeneration"] = 30
	nameToId["Shield"] = 31
	nameToId["ShieldRecoverTime"] = 32
	nameToId["SuperArmor"] = 33
	nameToId["SuperArmorRecoverTime"] = 34
	nameToId["Vulnerable"] = 35
	nameToId["MentorEnergy"] = 36
	nameToId["MentorEnergyRecharge"] = 37
	nameToId["BeheadMax"] = 38
	nameToId["BrokenShieldStrength"] = 39
	nameToId["SuperArmorRemove"] = 40
	nameToId["SuperArmorPrevent"] = 41
	nameToId["EXEnergyRecharge"] = 42
	nameToId["SkillEnergyRecharge"] = 43
	nameToId["ShieldStrength"] = 44
	nameToId["RecoveryHPStrength"] = 45
	nameToId["ElementEffectAddition"] = 46
	nameToId["Vigour"] = 47
	nameToId["VigourRecoveryRate"] = 48
	nameToId["VigourMax"] = 49
	nameToId["PlaceHolder_50"] = 50
	nameToId["MatterWeakness"] = 51
	nameToId["EnergyWeakness"] = 52
	nameToId["ViodWeakness"] = 53
	nameToId["PlaceHolder_54"] = 54
	nameToId["MatterResistance"] = 55
	nameToId["EnergyResistance"] = 56
	nameToId["ViodResistance"] = 57
	nameToId["AttackAmplificationCoefficient"] = 58
	nameToId["PlaceHolder_59"] = 59
	nameToId["PlaceHolder_60"] = 60
	nameToId["PlaceHolder_61"] = 61
	nameToId["PlaceHolder_62"] = 62
	nameToId["PlaceHolder_63"] = 63
	nameToId["PlaceHolder_64"] = 64
	nameToId["PlaceHolder_65"] = 65
	nameToId["PlaceHolder_66"] = 66
	nameToId["PlaceHolder_67"] = 67
	nameToId["PlaceHolder_68"] = 68
	nameToId["PlaceHolder_69"] = 69
	nameToId["PlaceHolder_70"] = 70
	battleAttrInfo[1] = &Info{
		Min: 0.000000,
		Max: 0.000000,
	}
	battleAttrInfo[2] = &Info{
		Min: 0.000000,
		Max: 0.000000,
	}
	battleAttrInfo[3] = &Info{
		Min: 0.000000,
		Max: 0.000000,
	}
	battleAttrInfo[4] = &Info{
		Min: 0.000000,
		Max: 0.000000,
	}
	battleAttrInfo[5] = &Info{
		Min: 0.000000,
		Max: 0.000000,
	}
	battleAttrInfo[6] = &Info{
		Min: 0.000000,
		Max: 0.000000,
	}
	battleAttrInfo[7] = &Info{
		Min: 0.000000,
		Max: 0.000000,
	}
	battleAttrInfo[8] = &Info{
		Min: 0.000000,
		Max: 0.000000,
	}
	battleAttrInfo[9] = &Info{
		Min: 0.000000,
		Max: 0.000000,
	}
	battleAttrInfo[10] = &Info{
		Min: 0.000000,
		Max: 0.000000,
	}
	battleAttrInfo[11] = &Info{
		Min: 0.000000,
		Max: 10000000.000000,
	}
	battleAttrInfo[12] = &Info{
		Min: 0.000000,
		Max: 10.000000,
	}
	battleAttrInfo[13] = &Info{
		Min: 0.000000,
		Max: 10000000.000000,
	}
	battleAttrInfo[14] = &Info{
		Min: 0.000000,
		Max: 10.000000,
	}
	battleAttrInfo[15] = &Info{
		Min: 0.000000,
		Max: 10000000.000000,
	}
	battleAttrInfo[16] = &Info{
		Min: 0.000000,
		Max: 10.000000,
	}
	battleAttrInfo[17] = &Info{
		Min: 0.000000,
		Max: 0.000000,
	}
	battleAttrInfo[18] = &Info{
		Min: 0.000000,
		Max: 1.000000,
	}
	battleAttrInfo[19] = &Info{
		Min: 0.000000,
		Max: 1.000000,
	}
	battleAttrInfo[20] = &Info{
		Min: 0.000000,
		Max: 10.000000,
	}
	battleAttrInfo[21] = &Info{
		Min: 0.000000,
		Max: 10.000000,
	}
	battleAttrInfo[22] = &Info{
		Min: 0.000000,
		Max: 1.000000,
	}
	battleAttrInfo[23] = &Info{
		Min: 0.000000,
		Max: 10000000.000000,
	}
	battleAttrInfo[24] = &Info{
		Min: 0.000000,
		Max: 10000000.000000,
	}
	battleAttrInfo[25] = &Info{
		Min: 0.000000,
		Max: 10000000.000000,
	}
	battleAttrInfo[26] = &Info{
		Min: 0.000000,
		Max: 0.000000,
	}
	battleAttrInfo[27] = &Info{
		Min: 0.000000,
		Max: 10000000.000000,
	}
	battleAttrInfo[28] = &Info{
		Min: 0.000000,
		Max: 10000000.000000,
	}
	battleAttrInfo[29] = &Info{
		Min: 0.000000,
		Max: 10000.000000,
	}
	battleAttrInfo[30] = &Info{
		Min: 0.000000,
		Max: 10000000.000000,
	}
	battleAttrInfo[31] = &Info{
		Min: 0.000000,
		Max: 10000000.000000,
	}
	battleAttrInfo[32] = &Info{
		Min: 0.000000,
		Max: 10000000.000000,
	}
	battleAttrInfo[33] = &Info{
		Min: 0.000000,
		Max: 10000000.000000,
	}
	battleAttrInfo[34] = &Info{
		Min: 0.000000,
		Max: 10000000.000000,
	}
	battleAttrInfo[35] = &Info{
		Min: 0.000000,
		Max: 10.000000,
	}
	battleAttrInfo[36] = &Info{
		Min: 0.000000,
		Max: 100.000000,
	}
	battleAttrInfo[37] = &Info{
		Min: 0.000000,
		Max: 100.000000,
	}
	battleAttrInfo[38] = &Info{
		Min: 100.000000,
		Max: 10000000.000000,
	}
	battleAttrInfo[39] = &Info{
		Min: 0.000000,
		Max: 100.000000,
	}
	battleAttrInfo[40] = &Info{
		Min: 0.000000,
		Max: 100.000000,
	}
	battleAttrInfo[41] = &Info{
		Min: 0.000000,
		Max: 100.000000,
	}
	battleAttrInfo[42] = &Info{
		Min: 0.000000,
		Max: 10000.000000,
	}
	battleAttrInfo[43] = &Info{
		Min: 0.000000,
		Max: 10000.000000,
	}
	battleAttrInfo[44] = &Info{
		Min: 0.000000,
		Max: 10000.000000,
	}
	battleAttrInfo[45] = &Info{
		Min: 0.000000,
		Max: 10000.000000,
	}
	battleAttrInfo[46] = &Info{
		Min: 0.000000,
		Max: 10000.000000,
	}
	battleAttrInfo[47] = &Info{
		Min: 0.000000,
		Max: 10000.000000,
	}
	battleAttrInfo[48] = &Info{
		Min: 0.000000,
		Max: 10000.000000,
	}
	battleAttrInfo[49] = &Info{
		Min: 0.000000,
		Max: 10000.000000,
	}
	battleAttrInfo[50] = &Info{
		Min: 0.000000,
		Max: 0.000000,
	}
	battleAttrInfo[51] = &Info{
		Min: 0.000000,
		Max: 10.000000,
	}
	battleAttrInfo[52] = &Info{
		Min: 0.000000,
		Max: 10.000000,
	}
	battleAttrInfo[53] = &Info{
		Min: 0.000000,
		Max: 10.000000,
	}
	battleAttrInfo[54] = &Info{
		Min: 0.000000,
		Max: 10.000000,
	}
	battleAttrInfo[55] = &Info{
		Min: 0.000000,
		Max: 1.000000,
	}
	battleAttrInfo[56] = &Info{
		Min: 0.000000,
		Max: 1.000000,
	}
	battleAttrInfo[57] = &Info{
		Min: 0.000000,
		Max: 1.000000,
	}
	battleAttrInfo[58] = &Info{
		Min: 0.000000,
		Max: 100.000000,
	}
	battleAttrInfo[59] = &Info{
		Min: 0.000000,
		Max: 0.000000,
	}
	battleAttrInfo[60] = &Info{
		Min: 0.000000,
		Max: 0.000000,
	}
	battleAttrInfo[61] = &Info{
		Min: 0.000000,
		Max: 0.000000,
	}
	battleAttrInfo[62] = &Info{
		Min: 0.000000,
		Max: 0.000000,
	}
	battleAttrInfo[63] = &Info{
		Min: 0.000000,
		Max: 0.000000,
	}
	battleAttrInfo[64] = &Info{
		Min: 0.000000,
		Max: 0.000000,
	}
	battleAttrInfo[65] = &Info{
		Min: 0.000000,
		Max: 0.000000,
	}
	battleAttrInfo[66] = &Info{
		Min: 0.000000,
		Max: 0.000000,
	}
	battleAttrInfo[67] = &Info{
		Min: 0.000000,
		Max: 0.000000,
	}
	battleAttrInfo[68] = &Info{
		Min: 0.000000,
		Max: 0.000000,
	}
	battleAttrInfo[69] = &Info{
		Min: 0.000000,
		Max: 0.000000,
	}
	battleAttrInfo[70] = &Info{
		Min: 0.000000,
		Max: 0.000000,
	}
}


func GetNameById(id int32) string {
	return idxToName[id]
}

func GetIdByName(name string) int32 {
	return nameToId[name]
}

func GetBattleAttrInfo(idx int32) *Info {
	return battleAttrInfo[idx]
}

/*
func TransFormToFloat64(id, value int32) float64 {
	info, ok := battleAttrInfo[id]
	if !ok {
		return 0
	}

	if info.NeedFixed {
		v := float64(float64(value) / float64(GlobalConst.Table_.IDMap[1].BattleAttrRate))
		return v
	}
	return float64(value)
}
 */
