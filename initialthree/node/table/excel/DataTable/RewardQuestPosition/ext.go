package RewardQuestPosition

import "math/rand"

// 随机个数 不包括存在的，不能随机存在的位置
func RandPositions(count int, existPos map[int32]struct{}) []int32 {
	pos := []int32{}
	cfgPos := []int32{}
	idMap := GetIDMap()
	for id := range idMap {
		if _, ok := existPos[id]; !ok {
			cfgPos = append(cfgPos, id)
		}
	}

	if count >= len(cfgPos) {
		return cfgPos
	}

	for {
		if count == 0 {
			return pos
		}
		idx := rand.Intn(len(cfgPos))
		pos = append(pos, cfgPos[idx])
		cfgPos = append(cfgPos[:idx], cfgPos[idx+1:]...)
		count--
	}
}
