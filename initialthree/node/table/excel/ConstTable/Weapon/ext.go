package Weapon

func Get() *Weapon {
	return GetID(1)
}

func GetSupplyExp(itemID int32) int32 {
	def := Get()
	for _, v := range def.WeaponSupplyExpItemArray {
		if v.ItemID == itemID {
			return v.Exp
		}
	}
	return 0
}
