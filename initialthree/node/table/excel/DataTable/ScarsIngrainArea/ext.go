package ScarsIngrainArea

import "math/rand"

func GetPrisonByLevel(level int32) *ScarsIngrainArea {
	for _, v := range GetIDMap() {
		if level >= v.MinLevel && level <= v.MaxLevel {
			return v
		}
	}
	return nil
}

// 本来只会产出上期以外的id，但是防止策划失误， 如果下期不够从上期的选
func (this *ScarsIngrainArea) RandomBoss(exits map[int32]struct{}, count int) []int32 {
	arr := make([]int32, 0, len(this.BossIDsArray))
	for _, v := range this.BossIDsArray {
		if _, ok := exits[v.ID]; !ok {
			arr = append(arr, v.ID)
		}
	}

	if len(arr) < count {
		for id := range exits {
			arr = append(arr, id)
			if len(arr) == count {
				break
			}
		}
	}

	ret := make([]int32, 0, count)
	for n := 0; n < count && len(arr) > 0; n++ {
		idx := rand.Intn(len(arr))
		ret = append(ret, arr[idx])
		arr = append(arr[:idx], arr[idx+1:]...)
	}

	return ret
}
