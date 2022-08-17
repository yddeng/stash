package RewardQuest

import "initialthree/node/node_game/module/rewardQuest"

func getExist(userQuest *rewardQuest.RewardQuest, removed int32) map[int32]struct{} {
	exist := map[int32]struct{}{}
	for id := range userQuest.GetData() {
		if id != removed {
			pos := id / 10
			exist[pos] = struct{}{}
		}
	}
	return exist
}
