package CharacterLevelUpExp

//ID计算方式：rarity*1000 + 等级
func GetMaxExp(rarity, level int32) (int32, bool) {
	id := rarity*1000 + level
	def := GetID(id)
	if def == nil {
		return 0, false
	}
	return def.ExpToCurrentLevel, true
}
