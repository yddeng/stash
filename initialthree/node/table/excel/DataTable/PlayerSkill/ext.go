package PlayerSkill

import (
	"initialthree/node/table/excel"
	"sync/atomic"
)

func (this *PlayerSkill) GetMaxLevel() int32 {
	return int32(len(this.Skill))
}

// 默认解锁
func (this *PlayerSkill) GetDefaultUnlockSkillCond() (*Skill_, bool) {
	if len(this.Skill) == 0 {
		panic("player skill length is 0")
	}
	s := this.Skill[0]
	// 默认解锁，不消耗道具
	if s.Gold == 0 && len(s.CostItems()) == 0 {
		return s, true
	}
	return this.Skill[0], false
}

var skillCost atomic.Value

type Items struct {
	ItemID int32
	Count  int32
}

func (this *Skill_) CostItems() []*Items {
	return skillCost.Load().(map[*Skill_][]*Items)[this]
}

func (this *Table) AfterLoad() {
	idMap := GetIDMap()
	values := make(map[*Skill_][]*Items, len(idMap))
	for _, v := range idMap {
		for _, bl := range v.Skill {
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
		}
	}

	skillCost.Store(values)
}
