package Dungeon

import (
	"initialthree/node/table/excel"
	"strings"
	"sync/atomic"
)

var unlock atomic.Value

type Unlock struct {
	Type int32
	Args []int32
}

func GetUnlock(id int32) []*Unlock {
	return unlock.Load().(map[int32][]*Unlock)[id]
}

func (this *Table) AfterLoad() {
	idMap := GetIDMap()
	values := make(map[int32][]*Unlock, len(idMap))
	for id, v := range idMap {
		list := make([]*Unlock, 0, len(v.Unlocks))
		for _, l := range v.Unlocks {
			un := &Unlock{Type: l.Type}
			ss := strings.Split(l.Arg, ",")
			ids := make([]int32, 0, len(ss))
			for _, s := range ss {
				ids = append(ids, excel.ReadInt32(s))
			}
			un.Args = ids
			list = append(list, un)
		}
		values[id] = list
	}

	unlock.Store(values)
}
