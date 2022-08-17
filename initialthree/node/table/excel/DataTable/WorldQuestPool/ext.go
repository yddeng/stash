package WorldQuestPool

import (
	"initialthree/node/table/excel"
	"initialthree/node/table/excel/DataTable/WorldQuest"
	"math/rand"
	"sync/atomic"
)

var list atomic.Value

type QuestList struct {
	List []int32
}

func (this *WorldQuestPool) GetQuestList() []*QuestList {
	return list.Load().(map[*WorldQuestPool][]*QuestList)[this]
}

func (this *Table) AfterLoad() {
	idMap := GetIDMap()
	values := make(map[*WorldQuestPool][]*QuestList, len(idMap))

	for _, v := range idMap {
		list := make([]*QuestList, 0, len(v.QuestList))
		for _, bl := range v.QuestList {
			r := excel.Split(bl.List, ",")
			ret := make([]int32, 0, len(r))
			for _, v := range r {
				id := excel.ReadInt32(v[0])
				ret = append(ret, id)

			}
			list = append(list, &QuestList{List: ret})
		}
		values[v] = list
	}

	list.Store(values)
}

func GetByLevel(level int32) *WorldQuestPool {
	for _, v := range GetIDMap() {
		if level > v.MinLevel && level <= v.MaxLevel {
			return v
		}
	}
	return nil
}

func (this *WorldQuestPool) RandomWorldQuest(count int32, unlock map[int]struct{}, dones map[int32]struct{}) map[int32]bool {

	var def *WorldQuest.WorldQuest
	pool := this.GetQuestList()
	quest := map[int32]struct{}{} // 去重后的
	for idx, p := range pool {
		if _, ok := unlock[idx]; ok && len(p.List) > 0 {
			for _, id := range p.List {
				def = WorldQuest.GetID(id)
				if _, done := dones[id]; !done && def != nil {
					if _, has := quest[id]; !has {
						quest[id] = struct{}{}
					}
				}
			}
		}
	}

	questSlice := make([]int32, 0, len(quest))
	for id := range quest {
		questSlice = append(questSlice, id)
	}

	// 已经完成的不在出现， 如果少于也可出现

	total := int32(len(quest)) // 防止策划配置错误
	ret := map[int32]bool{}

	if total > count {
		for count > 0 {
			idx := rand.Int() % len(questSlice)
			id := questSlice[idx]
			questSlice = append(questSlice[:idx], questSlice[idx+1:]...)
			count--
			ret[id] = false
		}
	} else if total == count {
		for id := range quest {
			ret[id] = false
		}
	} else {
		for id := range quest {
			ret[id] = false
		}
		if int32(len(ret)+len(dones)) <= count {
			for id := range dones {
				ret[id] = false
			}
		} else {
			for id := range dones {
				ret[id] = false
				if int32(len(ret)) >= count {
					break
				}
			}
		}
	}

	return ret

}
