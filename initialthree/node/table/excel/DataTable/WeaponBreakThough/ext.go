package WeaponBreakThough

import (
	"initialthree/node/common/battleAttr"
	"initialthree/node/table/excel"
	"sync/atomic"
)

var (
	bl2CostItems atomic.Value //map[*BreakLevel_][]*Items
	bl2Attribute atomic.Value //map[*BreakLevel_]map[int32]float64
)

type Items struct {
	ItemID int32
	Count  int32
}

func (this *BreakLevel_) CostItems() []*Items {
	return bl2CostItems.Load().(map[*BreakLevel_][]*Items)[this]
}

func (this *BreakLevel_) Attribute() map[int32]float64 {
	return bl2Attribute.Load().(map[*BreakLevel_]map[int32]float64)[this]
}

func (this *WeaponBreakThough) GetBreakTimesLevel(times int32) *BreakLevel_ {
	idx := times - 1
	if idx < 0 || idx >= int32(len(this.BreakLevel)) {
		return nil
	}
	return this.BreakLevel[idx]
}

func (this *Table) AfterLoad() {
	idMap := GetIDMap()
	values := make(map[*BreakLevel_][]*Items, len(idMap))
	values1 := make(map[*BreakLevel_]map[int32]float64, len(idMap))
	for _, v := range idMap {
		for _, bl := range v.BreakLevel {
			r := excel.Split(bl.ItemStr, ",#")
			ret := make([]*Items, 0, len(r))
			for _, v := range r {
				if len(v) == 2 {
					e := &Items{}
					e.ItemID = excel.ReadInt32(v[0])
					e.Count = excel.ReadInt32(v[1])
					ret = append(ret, e)
				}
			}
			values[bl] = ret

			r = excel.Split(bl.AttrStr, ",#")
			m := map[int32]float64{}
			for _, v := range r {
				if len(v) == 2 {
					id := battleAttr.GetIdByName(v[0])
					v := excel.ReadFloat(v[1])
					m[id] += v
				}
			}
			values1[bl] = m
		}
	}

	bl2CostItems.Store(values)
	bl2Attribute.Store(values1)
}
