package DropPool

import (
	"fmt"
)

type UpWeight struct {
	Type int32
	ID   int32
	Rate int32
}

func FindUpRate(t, id int32, up []*UpWeight) (int32, bool) {
	for _, v := range up {
		if v.Type == t && v.ID == id {
			return v.Rate, true
		}
	}
	return 1, false
}

func (this *Table) AfterLoadAll() {
	selfMap := GetIDMap()
	columnMap := getColumnIDMap()

	tmpIDMap := make(map[int32]*DropPool, len(selfMap)+len(columnMap))
	for id, p := range selfMap {
		tmpIDMap[id] = p
	}

	for id, p := range columnMap {
		if _, ok := tmpIDMap[id]; ok {
			panic(fmt.Sprintf("tmpIDMap id %d is exist", id))
		}
		tmpIDMap[id] = p
	}

	this.indexID.Store(tmpIDMap)
}
