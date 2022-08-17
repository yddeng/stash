package CharacterLevelUpAttribute

//ID计算方式：玩家角色ID*1000 + 等级
func GetAttribute(characterID, level int32) *CharacterLevelUpAttribute {
	id := characterID*1000 + level
	return GetID(id)
}
