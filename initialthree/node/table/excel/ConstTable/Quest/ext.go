package Quest

func GetDailyQuest(weekday int) []int32 {
	ids := make([]int32, 0, 4)
	def := GetID(1)

	switch weekday {
	case 0: // 周末
		for _, v := range def.DailyQuest_7Array {
			ids = append(ids, v.ID)
		}
	case 1:
		for _, v := range def.DailyQuest_1Array {
			ids = append(ids, v.ID)
		}
	case 2:
		for _, v := range def.DailyQuest_2Array {
			ids = append(ids, v.ID)
		}
	case 3:
		for _, v := range def.DailyQuest_3Array {
			ids = append(ids, v.ID)
		}
	case 4:
		for _, v := range def.DailyQuest_4Array {
			ids = append(ids, v.ID)
		}
	case 5:
		for _, v := range def.DailyQuest_5Array {
			ids = append(ids, v.ID)
		}
	case 6:
		for _, v := range def.DailyQuest_6Array {
			ids = append(ids, v.ID)
		}
	}

	return ids
}

func GetWeekdayQuest() []int32 {
	ids := make([]int32, 0, 4)
	def := GetID(1)

	for _, v := range def.WeeklyQuestArray {
		ids = append(ids, v.ID)
	}
	return ids
}

func GetMainQuest() []int32 {
	ids := make([]int32, 0, 4)
	def := GetID(1)

	for _, v := range def.InitMainStoryQuestArray {
		ids = append(ids, v.ID)
	}
	return ids
}

func GetDailyRewardQuest() []int32 {
	ids := make([]int32, 0, 4)
	def := GetID(1)

	for _, v := range def.DailyRewardQuestArray {
		ids = append(ids, v.ID)
	}
	return ids
}
