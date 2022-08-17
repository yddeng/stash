package BossKilledReward

import (
	"initialthree/node/common/attr"
	"initialthree/node/common/droppool"
	"initialthree/node/common/enumType"
	"initialthree/node/table/excel/DataTable/DropPool"
	"initialthree/pkg/util"
	"time"
)

// 掉落池随机次数 = boss数

func GenDropAward(difficult, killed int32) []*droppool.Award {

	def := GetID(difficult)
	if def == nil {
		return nil
	}

	if def.EquipDropID == 0 && def.ExpDropID == 0 {
		return nil
	}

	times := killed
	if times <= 0 {
		return nil
	}

	var up []*DropPool.UpWeight
	weekDay := time.Now().Weekday()
	switch weekDay {
	case time.Sunday:
		up = make([]*DropPool.UpWeight, 0, len(def.WeekdayPoolUpArray))
		for _, v := range def.WeekdayPoolUpArray {
			up = append(up, &DropPool.UpWeight{
				Type: enumType.DropType_Pool,
				ID:   v.ID,
				Rate: def.WeekdaySUpRate,
			})
		}
	case time.Monday:
		up = make([]*DropPool.UpWeight, 0, len(def.WeekdaySUpID_1Array))
		for _, v := range def.WeekdaySUpID_1Array {
			up = append(up, &DropPool.UpWeight{
				Type: enumType.DropType_Equip,
				ID:   v.ID,
				Rate: def.WeekdaySUpRate,
			})
		}

	case time.Tuesday:
		up = make([]*DropPool.UpWeight, 0, len(def.WeekdaySUpID_2Array))
		for _, v := range def.WeekdaySUpID_2Array {
			up = append(up, &DropPool.UpWeight{
				Type: enumType.DropType_Equip,
				ID:   v.ID,
				Rate: def.WeekdaySUpRate,
			})
		}
	case time.Wednesday:
		up = make([]*DropPool.UpWeight, 0, len(def.WeekdaySUpID_3Array))
		for _, v := range def.WeekdaySUpID_3Array {
			up = append(up, &DropPool.UpWeight{
				Type: enumType.DropType_Equip,
				ID:   v.ID,
				Rate: def.WeekdaySUpRate,
			})
		}
	case time.Thursday:
		up = make([]*DropPool.UpWeight, 0, len(def.WeekdaySUpID_4Array))
		for _, v := range def.WeekdaySUpID_4Array {
			up = append(up, &DropPool.UpWeight{
				Type: enumType.DropType_Equip,
				ID:   v.ID,
				Rate: def.WeekdaySUpRate,
			})
		}
	case time.Friday:
		up = make([]*DropPool.UpWeight, 0, len(def.WeekdaySUpID_5Array))
		for _, v := range def.WeekdaySUpID_5Array {
			up = append(up, &DropPool.UpWeight{
				Type: enumType.DropType_Equip,
				ID:   v.ID,
				Rate: def.WeekdaySUpRate,
			})
		}
	case time.Saturday:
		up = make([]*DropPool.UpWeight, 0, len(def.WeekdaySUpID_6Array))
		for _, v := range def.WeekdaySUpID_6Array {
			up = append(up, &DropPool.UpWeight{
				Type: enumType.DropType_Equip,
				ID:   v.ID,
				Rate: def.WeekdaySUpRate,
			})
		}
	}

	awards := make([]*droppool.Award, 0, times*2)
	for i := int32(0); i < times; i++ {
		equip := droppool.DropWithID(def.EquipDropID, up...)
		awards = append(awards, equip)
		exp := droppool.DropWithID(def.ExpDropID)
		awards = append(awards, exp)
	}

	return awards
}

func GetSecretCoinAward(difficult int32) *droppool.Award {
	def := GetID(difficult)
	if def == nil {
		return nil
	}

	// SecretCoin
	award := droppool.NewAward()
	SecretCoin := util.Random(def.SecretCoinCount-def.SecretCoinWave, def.SecretCoinWave+def.SecretCoinCount)
	award.AddInfo(enumType.DropType_UsualAttribute, attr.SecretCoin, SecretCoin)

	return award
}
