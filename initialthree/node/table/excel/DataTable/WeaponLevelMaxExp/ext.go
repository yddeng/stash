package WeaponLevelMaxExp

func (this *WeaponLevelMaxExp) GetMaxExp(level int32) int32 {
	if level < 1 || int(level) > len(this.MaxExp) {
		return -1
	}
	return this.MaxExp[int(level)-1].Exp
}

func (this *WeaponLevelMaxExp) MaxLevel() int32 {
	return int32(len(this.MaxExp))
}

// 获取 从left到 level 的等级经验
func (this *WeaponLevelMaxExp) GetRangeTotalExp(left, right int32) int32 {
	exp := int32(0)
	max := this.MaxLevel()
	for i := left; i < right && i < max; i++ {
		exp += this.MaxExp[i-1].Exp
	}
	return exp
}
