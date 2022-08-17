package inoutput

//go test -covermode=count -v -coverprofile=coverage.out -run=.
//go tool cover -html=coverage.out

import (
	"fmt"
	"testing"
)

const (
	Type_Item       = 1
	Type_Equip      = 2
	Type_Character  = 3
	ItemCountLimit  = 10 //道具管理器上限
	EquipCountLimit = 10 //装备管理器上限
)

type Item struct {
	ID    int
	Count int
}

type Equip struct {
	UniID       int64
	ID          int
	Attrs       []int
	CharacterID int
}

type Character struct {
	ID      int
	EquipID int64
}

type ItemMgr struct {
	user  *User
	items map[int]*Item
}

func (this *ItemMgr) CheckRes(r ResDesc) bool {
	return true
}

func (this *ItemMgr) RemoveRes(r ResDesc) RollBackFunc {
	i, ok := this.items[r.ID]
	if !ok {
		return nil
	} else {
		if i.Count < r.Count {
			return nil
		} else {
			i.Count -= r.Count
			if i.Count == 0 {
				delete(this.items, i.ID)
			}
			return func() bool {
				if rb := this.AddRes(r); rb == nil {
					return false
				}
				return true
			}
		}
	}
}

func (this *ItemMgr) AddRes(r ResDesc) RollBackFunc {
	i, ok := this.items[r.ID]
	if !ok {
		if len(this.items) >= ItemCountLimit {
			return nil
		}

		this.items[r.ID] = &Item{
			ID:    r.ID,
			Count: r.Count,
		}

	} else {
		i.Count += r.Count
	}

	return func() bool {
		if rb := this.RemoveRes(r); rb == nil {
			return false
		}
		return true
	}
}

type EquipMgr struct {
	equips map[int64]*Equip
	user   *User
}

func (this *EquipMgr) CheckRes(r ResDesc) bool {
	return true
}

func (this *EquipMgr) RemoveRes(r ResDesc) RollBackFunc {
	id := int64(r.ID)
	e, ok := this.equips[id]
	if !ok {
		return nil
	} else {
		delete(this.equips, id)
		return func() bool {
			this.equips[id] = e
			return true
		}
	}
}

func (this *EquipMgr) AddRes(r ResDesc) RollBackFunc {
	if len(this.equips)+r.Count > EquipCountLimit {
		return nil
	}

	var ids []int64
	for i := 0; i < r.Count; i++ {
		id := this.user.genID()
		ids = append(ids, id)
		e := &Equip{
			UniID: id,
			ID:    r.ID,
		}
		this.equips[e.UniID] = e
	}

	return func() bool {
		for _, id := range ids {
			delete(this.equips, id)
		}
		return true
	}
}

type CharaMgr struct {
	chara map[int]*Character
	user  *User
}

func (this *CharaMgr) CheckRes(r ResDesc) bool {
	return true
}

func (this *CharaMgr) RemoveRes(r ResDesc) RollBackFunc {
	return nil
}

func (this *CharaMgr) AddRes(r ResDesc) RollBackFunc {

	c, ok := this.chara[r.ID]
	if ok {
		// item
	}

	if r.Count > 1 {
		// item
	}

	c = &Character{
		ID: r.ID,
	}

	weaponID := 1
	if weaponID != 0 {
		eInsID := this.user.genID()
		e := &Equip{
			UniID: eInsID,
			ID:    weaponID,
		}
		this.user.eqMgr.equips[eInsID] = e

		e.CharacterID = c.ID
		c.EquipID = eInsID
	}

	return func() bool {
		return true
	}
}

type User struct {
	eqMgr        EquipMgr
	itemMgr      ItemMgr
	charaMgr     CharaMgr
	containerMgr map[int]Container
	nextID       int64
}

func (this *User) genID() int64 {
	this.nextID++
	return this.nextID
}

func (this *User) RegisterContainer(tt int, c Container) {
	this.containerMgr[tt] = c
}

func (this *User) GetContainer(tt int) Container {
	c, ok := this.containerMgr[tt]
	if ok {
		return c
	} else {
		return nil
	}
}

func (this *User) Print() {
	s := fmt.Sprintf("item length %d \n", len(this.itemMgr.items))
	for _, v := range this.itemMgr.items {
		s += fmt.Sprintf("< itemID %d count %d > ", v.ID, v.Count)
	}
	fmt.Println(s)

	s = fmt.Sprintf("equip length %d \n", len(this.eqMgr.equips))
	for _, v := range this.eqMgr.equips {
		s += fmt.Sprintf("< UniID %d ID %d CharacterID %d > ", v.UniID, v.ID, v.CharacterID)
	}
	fmt.Println(s)

	s = fmt.Sprintf("character length %d \n", len(this.charaMgr.chara))
	for _, v := range this.charaMgr.chara {
		s += fmt.Sprintf("< ID %d EquipID %d > ", v.ID, v.EquipID)
	}
	fmt.Println(s)
}

func NewUser() *User {
	u := &User{
		containerMgr: map[int]Container{},
	}

	u.eqMgr = EquipMgr{equips: map[int64]*Equip{}, user: u}
	u.itemMgr = ItemMgr{items: map[int]*Item{}, user: u}
	u.charaMgr = CharaMgr{chara: map[int]*Character{}, user: u}

	u.RegisterContainer(Type_Item, &u.itemMgr)
	u.RegisterContainer(Type_Equip, &u.eqMgr)
	u.RegisterContainer(Type_Character, &u.charaMgr)
	return u
}

func TestInOutPut(t *testing.T) {
	u := NewUser()

	//添加道具
	out := []ResDesc{
		ResDesc{
			ID:    1,
			Type:  Type_Item,
			Count: 10,
		},
		ResDesc{
			ID:    2,
			Type:  Type_Item,
			Count: 10,
		},
		ResDesc{
			ID:    3,
			Type:  Type_Item,
			Count: 10,
		},
		ResDesc{
			ID:    4,
			Type:  Type_Item,
			Count: 10,
		},
		ResDesc{
			ID:    5,
			Type:  Type_Item,
			Count: 10,
		},
		ResDesc{
			ID:    6,
			Type:  Type_Item,
			Count: 10,
		},
		ResDesc{
			ID:    7,
			Type:  Type_Item,
			Count: 10,
		},
		ResDesc{
			ID:    8,
			Type:  Type_Item,
			Count: 10,
		},
		ResDesc{
			ID:    9,
			Type:  Type_Item,
			Count: 10,
		},
	}

	err := DoInputOutput(u, nil, out)
	if err != nil {
		fmt.Println(err)
	}
	u.Print()

	//out = []ResDesc{
	//	ResDesc{
	//		ID:    10,
	//		Type:  Type_Item,
	//		Count: 10,
	//	},
	//	ResDesc{
	//		ID:    11,
	//		Type:  Type_Item,
	//		Count: 10,
	//	},
	//}
	//
	//err = DoInputOutput(u, nil, out)
	//if err != nil {
	//	fmt.Println(err)
	//}
	//
	//for i := 1; i <= 9; i++ {
	//	item, ok := u.itemMgr.items[i]
	//	if ok {
	//		fmt.Println(item.ID, item.Count)
	//	}
	//}

	in := []ResDesc{
		{
			ID:    1,
			Type:  Type_Item,
			Count: 3,
		},
		{
			ID:    2,
			Type:  Type_Item,
			Count: 5,
		},
	}

	out = []ResDesc{
		{
			ID:    10,
			Type:  Type_Item,
			Count: 10,
		},
		{
			ID:    3,
			Type:  Type_Equip,
			Count: 2,
		},
		{
			ID:    1,
			Type:  Type_Character,
			Count: 2,
		},
	}

	err = DoInputOutput(u, in, out)
	if err != nil {
		fmt.Println(err)
	}
	u.Print()
	fmt.Println("end")
}
