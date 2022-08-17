package user

import (
	"fmt"
	attr2 "initialthree/node/common/attr"
	"initialthree/node/common/enumType"
	"initialthree/node/common/inoutput"
	"initialthree/node/common/timeDisposal"
	"initialthree/node/node_game/event"
	"initialthree/node/node_game/module"
	"initialthree/node/node_game/module/attr"
	"initialthree/node/node_game/module/backpack"
	"initialthree/node/node_game/module/character"
	"initialthree/node/node_game/module/equip"
	"initialthree/node/node_game/module/weapon"
	"initialthree/node/table/excel/DataTable/Equip"
	ItemTable "initialthree/node/table/excel/DataTable/Item"
	"initialthree/node/table/excel/DataTable/PlayerCharacter"
	"initialthree/node/table/excel/DataTable/Weapon"
)

type containerCreateFunc func(user *User) inoutput.Container

func (this *User) GetContainer(tt int) inoutput.Container {
	if create := containerCreate[tt]; create != nil {
		return create(this)
	}
	panic(fmt.Sprintf("container %d is nil", tt))
}

func (this *User) RegisterContainer(tt int, c inoutput.Container) {}

// 提前实例化，返回实例化后的配置。用于前端展示
func (this *User) OutputIns(out []inoutput.ResDesc) []inoutput.ResDesc {
	insOut := make([]inoutput.ResDesc, 0, len(out))
	for _, v := range out {
		switch v.Type {
		case enumType.IOType_Equip, enumType.IOType_Weapon, enumType.IOType_Character:
			for i := int32(0); i < v.Count; i++ {
				insID := this.GenUID()
				insOut = append(insOut, inoutput.ResDesc{Type: v.Type, ID: v.ID, Count: 1, InsID: insID})
			}
		default:
			insOut = append(insOut, v)
		}
	}
	return insOut
}

/*       ********************** attr **************************        */

type attrContainer struct {
	user *User
}

func (this *attrContainer) RemoveRes(desc inoutput.ResDesc) (error, inoutput.RollBackFunc, inoutput.CommitFunc) {
	attrModule := this.user.GetSubModule(module.Attr).(*attr.UserAttr)
	v, err := attrModule.GetAttr(desc.ID)
	if err != nil {
		return inoutput.ErrInvalidResType, nil, nil
	}
	if v < int64(desc.Count) {
		return inoutput.ErrInputNotEnough, nil, nil
	}

	_, _ = attrModule.AddAttr(desc.ID, -int64(desc.Count))
	return nil, func() {
		_, _ = attrModule.AddAttr(desc.ID, int64(desc.Count))
	}, nil
}

func (this *attrContainer) AddRes(desc inoutput.ResDesc) (error, inoutput.RollBackFunc, inoutput.CommitFunc, []inoutput.ResDesc) {
	attrModule := this.user.GetSubModule(module.Attr).(*attr.UserAttr)
	_, _ = attrModule.AddAttr(desc.ID, int64(desc.Count))
	return nil, func() {
		_, _ = attrModule.AddAttr(desc.ID, -int64(desc.Count))
	}, nil, nil
}

/*       ********************** equip **************************        */

type equipContainer struct {
	user *User
}

func (this *equipContainer) RemoveRes(desc inoutput.ResDesc) (error, inoutput.RollBackFunc, inoutput.CommitFunc) {
	equipModule := this.user.GetSubModule(module.Equip).(*equip.UserEquip)
	e := equipModule.GetEquip(uint32(desc.ID))
	if e == nil {
		return inoutput.ErrInputNotEnough, nil, nil
	}

	equipModule.Remove(e.InsID)

	return nil, func() {
		equipModule.AddEquip(e)
	}, nil
}

func (this *equipContainer) AddRes(desc inoutput.ResDesc) (error, inoutput.RollBackFunc, inoutput.CommitFunc, []inoutput.ResDesc) {
	// 无配置
	def := Equip.GetID(desc.ID)
	if def == nil {
		return inoutput.ErrCfgNotFound, nil, nil, nil
	}

	if desc.Count > 1 {
		out := []inoutput.ResDesc{}
		for i := int32(0); i < desc.Count; i++ {
			out = append(out, inoutput.ResDesc{Type: enumType.IOType_Equip, ID: desc.ID, Count: 1})
		}
		return nil, nil, nil, out
	} else {

		equipModule := this.user.GetSubModule(module.Equip).(*equip.UserEquip)
		userAttr := this.user.GetSubModule(module.Attr).(*attr.UserAttr)
		eCap := equipModule.GetUseCap()
		attrCap, _ := userAttr.GetAttr(attr2.EquipCap)

		if int64(eCap+1) > attrCap {
			return inoutput.ErrSpaceNotEnough, nil, nil, nil
		}

		var insId uint32
		if nil != desc.InsID {
			insId = desc.InsID.(uint32)
		}
		if 0 == insId {
			insId = this.user.GenUID()
		}

		e := equipModule.NewEquip(desc.ID, insId)
		equipModule.AddEquip(e)

		return nil, func() {
				equipModule.Remove(insId)
			}, func() {
				this.user.EmitEvent(event.EventEquipAdd, insId)
			}, nil
	}
}

/*       ********************** item **************************        */

type itemContainer struct {
	user *User
}

func (this *itemContainer) RemoveRes(desc inoutput.ResDesc) (error, inoutput.RollBackFunc, inoutput.CommitFunc) {
	backpackModule := this.user.GetSubModule(module.Backpack).(*backpack.Backpack)

	useItems := map[*backpack.Item]int32{}
	removeItems := map[*backpack.Item]struct{}{}

	if desc.InsID != nil {
		it := backpackModule.GetItem(desc.InsID.(uint32))
		if it.Count < desc.Count {
			return inoutput.ErrInputNotEnough, nil, nil
		}

		if it.Count == desc.Count {
			backpackModule.RemItem(it.ID)
			removeItems[it] = struct{}{}
		} else {
			backpackModule.AddItemCount(it.ID, -desc.Count)
			useItems[it] = desc.Count
		}

	} else {

		count := backpackModule.GetItemCountByTID(desc.ID)
		if count < desc.Count {
			return inoutput.ErrInputNotEnough, nil, nil
		}

		needUseCount := desc.Count
		items := backpackModule.GetItemsByTID(desc.ID)
		for _, v := range items {
			if v.Count <= needUseCount {
				backpackModule.RemItem(v.ID)
				needUseCount -= v.Count
				removeItems[v] = struct{}{}
			} else {
				backpackModule.AddItemCount(v.ID, -needUseCount)
				needUseCount = 0
				useItems[v] = needUseCount
			}

			if needUseCount == 0 {
				break
			}
		}
	}

	return nil, func() {
		// 移除的添加回来
		for it := range removeItems {
			backpackModule.AddItem(it)
		}
		// 扣除的数量加回来
		for it, count := range useItems {
			backpackModule.AddItemCount(it.ID, count)
		}
	}, nil
}

func (this *itemContainer) AddRes(desc inoutput.ResDesc) (error, inoutput.RollBackFunc, inoutput.CommitFunc, []inoutput.ResDesc) {
	itCfg := ItemTable.GetID(desc.ID)
	if itCfg == nil {
		return inoutput.ErrCfgNotFound, nil, nil, nil
	}

	backpackModule := this.user.GetSubModule(module.Backpack).(*backpack.Backpack)
	addItems := map[*backpack.Item]int32{} // 新创建的 数量为0

	if itCfg.GetTimeLimitType() == enumType.ItemTimeLimitType_Duration {
		it := backpack.CreateItem(this.user.GenUID(), desc.ID, desc.Count, timeDisposal.Now())
		backpackModule.AddItem(it)
		addItems[it] = 0
	} else {
		items := backpackModule.GetItemsByTID(desc.ID)
		if len(items) == 0 {
			it := backpack.CreateItem(this.user.GenUID(), desc.ID, desc.Count, timeDisposal.Now())
			backpackModule.AddItem(it)
			addItems[it] = 0
		} else {
			for _, it := range items {
				backpackModule.AddItemCount(it.ID, desc.Count)
				addItems[it] = desc.Count
				break
			}
		}
	}

	return nil, func() {
		for it, count := range addItems {
			if count == 0 {
				backpackModule.RemItem(it.ID)
			} else {
				backpackModule.AddItemCount(it.ID, -count)
			}
		}
	}, nil, nil
}

/*       ********************** character **************************        */

type characterContainer struct {
	user *User
}

func (this *characterContainer) RemoveRes(desc inoutput.ResDesc) (error, inoutput.RollBackFunc, inoutput.CommitFunc) {
	// 不允许移除
	return inoutput.ErrInvalidResType, nil, nil
}

func (this *characterContainer) AddRes(desc inoutput.ResDesc) (error, inoutput.RollBackFunc, inoutput.CommitFunc, []inoutput.ResDesc) {
	def := PlayerCharacter.GetID(desc.ID)
	if def == nil {
		return inoutput.ErrCfgNotFound, nil, nil, nil
	}

	if desc.Count > 1 {
		out := make([]inoutput.ResDesc, 0, desc.Count)
		for i := int32(0); i < desc.Count; i++ {
			out = append(out, inoutput.ResDesc{Type: enumType.IOType_Character, ID: desc.ID, Count: 1})
		}
		return nil, nil, nil, out
	} else {

		characterModule := this.user.GetSubModule(module.Character).(*character.UserCharacter)
		c := characterModule.GetCharacter(desc.ID)
		if nil != c {
			// 角色未满命座时，抽取该重复角色将转换为该角色命座材料
			// 角色已满命座时，抽取该重复角色将转换为特殊代币
			// 满命座的检测应为：角色当前命座强化等级+当前该角色的命座材料=6

			var items []*PlayerCharacter.Item
			if c.HitTimes >= def.DrawCardTimes {
				items = def.GetFragmentMax()
			} else {
				items = def.GetFragment()
			}
			characterModule.AddHitTimes(c, 1)
			itemDesc := make([]inoutput.ResDesc, 0, len(items))
			for _, v := range items {
				itemDesc = append(itemDesc, inoutput.ResDesc{Type: enumType.IOType_Item, ID: v.ID, Count: v.Count})
			}
			return nil, func() {
				characterModule.AddHitTimes(c, -1)
			}, nil, itemDesc
		}

		c = characterModule.NewCharacter(desc.ID, def)

		var out []inoutput.ResDesc

		var commitFunc func()

		if def.DefaultWeapon != 0 {
			// 初始化武器
			defWeapon := Weapon.GetID(def.DefaultWeapon)
			if defWeapon == nil {
				return inoutput.ErrCfgNotFound, nil, nil, nil
			}

			// 用于给武器实例化
			var insId uint32
			if nil != desc.InsID {
				insId = desc.InsID.(uint32)
			}
			if 0 == insId {
				insId = this.user.GenUID()
			}
			out = []inoutput.ResDesc{{Type: enumType.IOType_Weapon, ID: def.DefaultWeapon, Count: 1, InsID: insId}}

			commitFunc = func() {
				userWeapon := this.user.GetSubModule(module.Weapon).(*weapon.UserWeapon)
				e := userWeapon.GetWeapon(insId)
				if nil != e {
					userWeapon.Equip(e, c.CharacterID)
					characterModule.WeaponReplace(c, e.InsID)
				}
			}
		}

		characterModule.AddCharacter(c)

		return nil, func() {
			characterModule.RemoveCharacter(desc.ID)
		}, commitFunc, out
	}
}

/*       ********************** weapon **************************        */

type weaponContainer struct {
	user *User
}

func (this *weaponContainer) RemoveRes(desc inoutput.ResDesc) (error, inoutput.RollBackFunc, inoutput.CommitFunc) {
	weaponModule := this.user.GetSubModule(module.Weapon).(*weapon.UserWeapon)
	e := weaponModule.GetWeapon(uint32(desc.ID))
	if e == nil {
		return inoutput.ErrInputNotEnough, nil, nil
	}

	weaponModule.Remove(e.InsID)

	return nil, func() {
		weaponModule.AddWeapon(e)
	}, nil
}

func (this *weaponContainer) AddRes(desc inoutput.ResDesc) (error, inoutput.RollBackFunc, inoutput.CommitFunc, []inoutput.ResDesc) {
	// 无配置
	def := Weapon.GetID(desc.ID)
	if def == nil {
		return inoutput.ErrCfgNotFound, nil, nil, nil
	}

	if desc.Count > 1 {
		out := []inoutput.ResDesc{}
		for i := int32(0); i < desc.Count; i++ {
			out = append(out, inoutput.ResDesc{Type: enumType.IOType_Weapon, ID: desc.ID, Count: 1})
		}
		return nil, nil, nil, out
	} else {

		weaponModule := this.user.GetSubModule(module.Weapon).(*weapon.UserWeapon)
		userAttr := this.user.GetSubModule(module.Attr).(*attr.UserAttr)
		eCap := weaponModule.GetUseCap()
		attrCap, _ := userAttr.GetAttr(attr2.WeaponCap)

		if int64(eCap+1) > attrCap {
			return inoutput.ErrSpaceNotEnough, nil, nil, nil
		}

		var insId uint32
		if nil != desc.InsID {
			insId = desc.InsID.(uint32)
		}
		if 0 == insId {
			insId = this.user.GenUID()
		}

		e := weaponModule.NewWeapon(desc.ID, insId)
		weaponModule.AddWeapon(e)

		return nil, func() {
				weaponModule.Remove(insId)
			}, func() {
				this.user.EmitEvent(event.EventWeaponAdd, insId)
			}, nil
	}
}

var containerCreate map[int]containerCreateFunc

func init() {
	containerCreate = map[int]containerCreateFunc{
		enumType.IOType_UsualAttribute: func(user *User) inoutput.Container {
			return &attrContainer{
				user: user,
			}
		},
		enumType.IOType_Equip: func(user *User) inoutput.Container {
			return &equipContainer{
				user: user,
			}
		},
		enumType.IOType_Weapon: func(user *User) inoutput.Container {
			return &weaponContainer{
				user: user,
			}
		},
		enumType.IOType_Character: func(user *User) inoutput.Container {
			return &characterContainer{
				user: user,
			}
		},
		enumType.IOType_Item: func(user *User) inoutput.Container {
			return &itemContainer{
				user: user,
			}
		},
	}
}
