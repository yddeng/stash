package EquipRandomAttributePool

import (
	"initialthree/pkg/util"
)

func (this *EquipRandomAttributePool) RandomID() int32 {
	weight := int32(0)
	for _, v := range this.Random {
		weight += v.Weight
	}

	d := util.Random(1, weight)
	for _, v := range this.Random {
		d -= v.Weight
		if d <= 0 {
			return v.ID
		}
	}
	return 0
}
