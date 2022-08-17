package Equip

func GetSupplyExp(itemID int32) int32 {
	def := Get()
	for _, v := range def.EquipSupplyExpItemArray {
		if v.ItemID == itemID {
			return v.Exp
		}
	}
	return 0
}

func Get() *Equip {
	return GetID(1)
}
