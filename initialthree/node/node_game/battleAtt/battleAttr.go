package battleAtt

import (
	"initialthree/node/table/excel/DataTable/CharacterBreakThrough"
	"initialthree/node/table/excel/DataTable/CharacterLevelUpAttribute"
	"initialthree/node/table/excel/DataTable/Equip"
	"initialthree/node/table/excel/DataTable/EquipAttribute"
	"initialthree/node/table/excel/DataTable/PlayerCharacter"
	"initialthree/node/table/excel/DataTable/PlayerGene"
	"initialthree/node/table/excel/DataTable/Weapon"
	"initialthree/node/table/excel/DataTable/WeaponAttribute"
	"initialthree/node/table/excel/DataTable/WeaponBreakThough"
	"initialthree/zaplogger"
)

/*
  用于计算玩家战斗属性。
*/

/*
最终增伤的计算方式与攻击汇总不同，攻击是各个系统间值相加，而最终增伤是（各个系统间值+1）后相乘。
例如，坐骑系统有增伤属性50%，时装系统50%，则最终增伤=（1+坐骑50%）*（1+时装50%）=1.5*1.5=2.25

最终免伤的计算方式与攻击汇总不同，攻击是各个系统间值相加，而最终免伤是（各个系统间值）相乘（注意不+1）。
例如，坐骑系统有免伤属性50%，时装系统50%，则最终免伤=(1-坐骑50%)*(1-时装50%)=0.5*0.5=0.25

元素抗性 与之前计算伤害D的免伤方式一致，系统内的免伤相加，系统间的由（1-系统内免伤和）后相乘。
*/

func RecalculateBattleAttr(battleChara *BattleCharacter) map[int32]float64 {
	if battleChara == nil {
		return nil
	}

	characterAttrs := calculateCharacter(battleChara)
	weaponAttrs := calculateWeapon(battleChara)
	equipAttrs := calculateEquip(battleChara)

	finalAttrs := systemAttrCalculate(characterAttrs, weaponAttrs, equipAttrs)
	finalAttr(finalAttrs)
	return finalAttrs
}

// 角色属性（基础属性直接相加，特殊属性处理）
// 增伤举例：(1+升级)*(1+晋升)*（1+技能）*（1+进化）*（1+解放），算作不同的系统
func calculateCharacter(chara *BattleCharacter) map[int32]float64 {
	def := PlayerCharacter.GetID(chara.CharacterID)
	// 等级属性
	levelAttr := map[int32]float64{}
	levelDef := CharacterLevelUpAttribute.GetAttribute(chara.CharacterID, chara.Level)
	if levelDef == nil {
		return newTransAttr()
	}

	for _, v := range levelDef.AttributeListArray {
		addAttrs(levelAttr, v.ID, v.Val)
	}
	transAttr(levelAttr)

	// 突破属性
	breakAttr := map[int32]float64{}
	if chara.BreakLevel > 0 {
		breakDef := CharacterBreakThrough.GetID(def.GetBreakID(chara.BreakLevel))
		for _, v := range breakDef.AttributeBonusArray {
			addAttrs(breakAttr, v.ID, v.Val)
		}
	}
	transAttr(breakAttr)

	// 命座属性
	geneAttr := map[int32]float64{}
	if chara.GeneLevel > 0 {
		geneDef := PlayerGene.GetGene(chara.CharacterID, chara.GeneLevel)
		for _, v := range geneDef.AttriArray {
			addAttrs(geneAttr, v.ID, v.Val)
		}
	}
	transAttr(geneAttr)

	return systemAttrCalculate(levelAttr, breakAttr, geneAttr)
}

// 武器属性
// 增伤举例：(1+武器强化+武器突破+武器基本属性）*（1+武器共鸣1属性+武器共鸣2+……）
func calculateWeapon(chara *BattleCharacter) map[int32]float64 {

	doWeapon := chara.Weapon
	if doWeapon == nil {
		zaplogger.GetSugar().Errorf("%s %d weapon is nil", chara.UserID, chara.GameID)
		return newTransAttr()
	}

	def := Weapon.GetID(doWeapon.ConfigID)
	if def == nil {
		zaplogger.GetSugar().Errorf("weapon def config %d  not found", doWeapon.ConfigID)
		return newTransAttr()
	}
	breakDef := WeaponBreakThough.GetID(def.BreakThroughConfig)
	if breakDef == nil {
		zaplogger.GetSugar().Errorf("weapon breakDef config %d  not found", def.BreakThroughConfig)
		return newTransAttr()
	}

	// 强化等级属性
	baseAttr := map[int32]float64{}
	for _, v := range def.AttrConfigs {
		attrDef := WeaponAttribute.GetID(v.ID)
		if attrDef == nil {
			zaplogger.GetSugar().Errorf("weapon attr def config %d not found", v.ID)
			return newTransAttr()
		}

		val := attrDef.LevelAttr[doWeapon.Level-1].Val
		addAttrs(baseAttr, attrDef.AttrID(), val)
	}
	transAttr(baseAttr)

	// 突破属性
	breakAttr := map[int32]float64{}
	if doWeapon.BreakTimes > 0 {
		breakDef := WeaponBreakThough.GetID(def.BreakThroughConfig)
		bl := breakDef.BreakLevel[doWeapon.BreakTimes-1]
		for id, v := range bl.Attribute() {
			addAttrs(breakAttr, id, v)
		}
	}
	transAttr(breakAttr)

	return systemAttrCalculate(baseAttr, breakAttr)

}

// 意识属性，1-6号位
// 1+意识基础属性（意识1+意识2+……）+意识强化（意识强化1+意识强化2+……）+意识突破（意识突破1+意识突破2+……））*（1+意识共鸣1+意识共鸣2+……）
func calculateEquip(chara *BattleCharacter) map[int32]float64 {
	baseAttr := map[int32]float64{} // 6意识的 基础、强化、突破 和
	for _, doEquip := range chara.Equips {
		// 基础、强化
		def := Equip.GetID(doEquip.ConfigID)
		if def == nil {
			zaplogger.GetSugar().Errorf("equipDef config %d  not found", doEquip.ConfigID)
			continue
		}

		for _, v := range def.AttrConfigs {
			attrDef := EquipAttribute.GetID(v.ID)
			if attrDef == nil {
				zaplogger.GetSugar().Errorf("equip attr def config %d not found", v.ID)
				return newTransAttr()
			}

			val := attrDef.Attr[doEquip.Level-1].Val
			addAttrs(baseAttr, attrDef.AttrID(), val)
		}

		// 随机属性
		attrDef := EquipAttribute.GetID(doEquip.RandomAttrID)
		if attrDef != nil {
			val := attrDef.Attr[doEquip.Refine[1]].Val
			addAttrs(baseAttr, attrDef.AttrID(), val)
		}
	}
	transAttr(baseAttr)

	return systemAttrCalculate(baseAttr)
}
