package PlayerCharacter

import (
	"fmt"
	"initialthree/node/table/excel/DataTable/CharacterResource"
	"sync/atomic"
)

func (this *PlayerCharacter) GetBreakID(idx int32) int32 {
	if idx <= 0 || idx > int32(len(this.BreakIDListArray)) {
		return -1
	}
	return this.BreakIDListArray[idx-1].ID

}

func (this *PlayerCharacter) GetDamageElementType() int32 {
	return CharacterResource.GetID(this.DefaultCharacterResourceID).DamageElementTypeEnum
}

type Item struct {
	ID    int32
	Count int32
}

var (
	fragmentItems atomic.Value //map[*PlayerCharacter][]*Item
	maxLvItems    atomic.Value //map[*PlayerCharacter][]*Item
)

func (this *PlayerCharacter) GetFragment() []*Item {
	return fragmentItems.Load().(map[*PlayerCharacter][]*Item)[this]
}

func (this *PlayerCharacter) GetFragmentMax() []*Item {
	return maxLvItems.Load().(map[*PlayerCharacter][]*Item)[this]
}

func (this *Table) AfterLoadAll() {
	idMap := GetIDMap()
	v1 := make(map[*PlayerCharacter][]*Item, len(idMap))
	v2 := make(map[*PlayerCharacter][]*Item, len(idMap))

	for id, v := range idMap {
		if len(v.DrawCardItemIDsArray) != len(v.DrawCardItemCountsArray) ||
			len(v.MaxGeneLvDrawCardItemIDsArray) != len(v.MaxGeneLvDrawCardItemCountsArray) {
			panic(fmt.Sprintf("tmpIDMap id %d is drawcard array not equal", id))
		}

		s1 := make([]*Item, 0, len(v.DrawCardItemIDsArray))
		for i, e := range v.DrawCardItemIDsArray {
			e2 := v.DrawCardItemCountsArray[i]
			s1 = append(s1, &Item{ID: e.ID, Count: e2.Count})
		}

		s2 := make([]*Item, 0, len(v.MaxGeneLvDrawCardItemIDsArray))
		for i, e := range v.MaxGeneLvDrawCardItemIDsArray {
			e2 := v.MaxGeneLvDrawCardItemCountsArray[i]
			s2 = append(s2, &Item{ID: e.ID, Count: e2.Count})
		}

		v1[v] = s1
		v2[v] = s2
	}

	fragmentItems.Store(v1)
	maxLvItems.Store(v2)
}
