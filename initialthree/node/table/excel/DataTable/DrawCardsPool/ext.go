package DrawCardsPool

import (
	"initialthree/node/common/enumType"
	"initialthree/node/table/excel/DataTable/DropPool"
	"sync/atomic"
)

var (
	drops     atomic.Value //map[int32]*DropPool.DropPool
	tenDrop   atomic.Value
	guarantee atomic.Value
)

func GetDrawCardsPool(id int32) *DropPool.DropPool {
	return drops.Load().(map[int32]*DropPool.DropPool)[id]
}

func GetTenGuaranteePool(id int32) *DropPool.DropPool {
	return tenDrop.Load().(map[int32]*DropPool.DropPool)[id]
}

func GetGuaranteePool(id int32) *DropPool.DropPool {
	return guarantee.Load().(map[int32]*DropPool.DropPool)[id]
}

func (this *Table) AfterLoad() {
	idMap := GetIDMap()

	values1 := make(map[int32]*DropPool.DropPool, len(idMap))
	values2 := make(map[int32]*DropPool.DropPool, len(idMap))
	values3 := make(map[int32]*DropPool.DropPool, len(idMap))

	for id, v := range idMap {
		list1 := make([]*DropPool.DropList_, 0, len(v.DropList))
		list2 := make([]*DropPool.DropList_, 0, len(v.TenTimesGuaranteePoolsArray))
		list3 := make([]*DropPool.DropList_, 0, 1)

		for i, l := range v.DropList {
			list1 = append(list1, &DropPool.DropList_{
				Type:   enumType.DropType_Pool,
				ID:     l.ID,
				Count:  1,
				Wave:   0,
				Weight: l.Weight,
			})

			for _, arr := range v.TenTimesGuaranteePoolsArray {
				if arr.Idx == int32(i) {
					list2 = append(list2, &DropPool.DropList_{
						Type:   enumType.DropType_Pool,
						ID:     l.ID,
						Count:  1,
						Wave:   0,
						Weight: l.Weight,
					})
					break
				}
			}

			if v.GuaranteePool == int32(i) {
				list3 = append(list3, &DropPool.DropList_{
					Type:   enumType.DropType_Pool,
					ID:     l.ID,
					Count:  1,
					Wave:   0,
					Weight: l.Weight,
				})

			}
		}

		values1[id] = &DropPool.DropPool{
			ID:         id,
			MinCount:   1,
			MaxCount:   1,
			Repeatable: false,
			TypeEnum:   enumType.DropPoolType_Rand,
			DropList:   list1,
		}

		values2[id] = &DropPool.DropPool{
			ID:         id,
			MinCount:   1,
			MaxCount:   1,
			Repeatable: false,
			TypeEnum:   enumType.DropPoolType_Rand,
			DropList:   list2,
		}

		values3[id] = &DropPool.DropPool{
			ID:         id,
			MinCount:   1,
			MaxCount:   1,
			Repeatable: false,
			TypeEnum:   enumType.DropPoolType_Rand,
			DropList:   list3,
		}

	}

	drops.Store(values1)
	tenDrop.Store(values2)
	guarantee.Store(values3)
}
