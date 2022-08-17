package EquipDecompose

func Get(rarity, pos, breakLevel int32) *EquipDecompose {
	id := rarity*100 + pos*10 + breakLevel
	return GetID(id)
}
