package MainDungeon

import "sync/atomic"

var (
	quest atomic.Value // []int32
)

func GetInstanceQuest() []int32 {
	return quest.Load().([]int32)
}

func (this *Table) AfterLoad() {
	idMap := GetIDMap()
	values := make([]int32, 0, len(idMap)*4)
	for _, v := range idMap {
		if v.QuestId != 0 {
			values = append(values, v.QuestId)
		}
	}

	quest.Store(values)
}
