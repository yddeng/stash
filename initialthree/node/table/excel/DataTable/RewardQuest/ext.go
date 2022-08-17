package RewardQuest

import (
	"fmt"
	"initialthree/node/common/enumType"
	"sort"
	"sync/atomic"
)

// 所有条件按照 从大到小排序
// 品质判断失败 1，提供品质数据不够。 2， 同位置小于需求品质
// 攻击属性失败 1，提供数量不够      2，同位置不同
func (this *RewardQuest) IsAccept(attacks, qualitys []int32) bool {
	cond := conditions.Load().(map[int32]*condition)[this.ID]
	if cond == nil {
		return false
	}

	if len(cond.attacks) > len(attacks) || len(cond.qualitys) > len(qualitys) {
		return false
	}

	sort.Slice(attacks, func(i, j int) bool { return attacks[i] > attacks[j] })
	sort.Slice(qualitys, func(i, j int) bool { return qualitys[i] > qualitys[j] })

	for i, a := range cond.attacks {
		if attacks[i] != a {
			return false
		}
	}

	for i, q := range cond.qualitys {
		if qualitys[i] < q {
			return false
		}
	}

	return true
}

func (this *RewardQuest) PosQuality() (int32, int32) {
	quality := this.ID % 10
	pos := this.ID / 10
	return pos, quality
}

// pos两位，quality一位
func GetPosQuality(pos, quality int32) *RewardQuest {
	id := pos*10 + quality
	return GetID(id)
}

// 条件按照 从大到小排序
type condition struct {
	attacks  []int32
	qualitys []int32
}

var (
	conditions atomic.Value //map[int32]*condition
)

func (this *Table) AfterLoad() {
	idMap := GetIDMap()
	values := make(map[int32]*condition, len(idMap))
	for id, v := range idMap {

		_, quality := v.PosQuality()
		switch quality {
		case enumType.RarityType_Star2, enumType.RarityType_Star3, enumType.RarityType_Star4, enumType.RarityType_Star5:
		default:
			panic(fmt.Sprintf("RewardQuest ID %d is qualiy %d failed", v.ID, quality))
		}

		attacks := []int32{}
		qualitys := []int32{}
		if v.AttackType_1 != 0 {
			attacks = append(attacks, v.AttackType_1)
		}
		if v.AttackType_2 != 0 {
			attacks = append(attacks, v.AttackType_2)
		}
		if v.AttackType_3 != 0 {
			attacks = append(attacks, v.AttackType_3)
		}

		if v.QualityType_1 != 0 {
			qualitys = append(qualitys, v.QualityType_1)
		}
		if v.QualityType_2 != 0 {
			qualitys = append(qualitys, v.QualityType_2)
		}
		if v.QualityType_3 != 0 {
			qualitys = append(qualitys, v.QualityType_3)
		}

		sort.Slice(attacks, func(i, j int) bool { return attacks[i] > attacks[j] })
		sort.Slice(qualitys, func(i, j int) bool { return qualitys[i] > qualitys[j] })
		values[id] = &condition{
			attacks:  attacks,
			qualitys: qualitys,
		}
	}

	conditions.Store(values)
}
