package EquipLevelMaxExp

func (this *EquipLevelMaxExp) GetMaxExp(level int32) int32 {
	if level < 1 || int(level) > len(this.MaxExp) {
		return -1
	}
	return this.MaxExp[int(level)-1].Exp
}

func (this *EquipLevelMaxExp) MaxLevel() int32 {
	return int32(len(this.MaxExp))
}

func (this *EquipLevelMaxExp) GetRangeTotalExp(left, right int32) int32 {
	exp := int32(0)
	max := this.MaxLevel()
	for i := left; i < right && i < max; i++ {
		exp += this.MaxExp[i-1].Exp
	}
	return exp
}
