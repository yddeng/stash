package backpack

import (
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/sniperHW/flyfish/client"
	"initialthree/node/common/enumType"
	"initialthree/node/node_game/module"
	Item2 "initialthree/node/table/excel/DataTable/Item"
	"initialthree/pkg/json"
	"initialthree/zaplogger"
	"time"

	cs_msg "initialthree/protocol/cs/message"
)

var (
	ErrItemNotFound       = errors.New("item not found")
	ErrItemAlreadyExist   = errors.New("item already exist")
	ErrItemAlreadyExpired = errors.New("item already expired")
)

const (
	dbTable   = "backpack"
	slotCount = 10
)

func slotName(slotIdx int) string {
	return fmt.Sprintf("slot%d", slotIdx)
}

func calcSlotIdx(id uint32) int {
	return int(id % slotCount)
}

type Item struct {
	ID          uint32 `json:"id"`
	TID         int32  `json:"tid"`
	Count       int32  `json:"cnt"`
	AcquireTime int64  `json:"at,omitempty"`
	TimeLimit   int64  `json:"tl,omitempty"`
}

func (it *Item) IsExpired(t time.Time) bool {
	return it.TimeLimit != 0 && t.Unix() >= it.TimeLimit
}

func CreateItem(id uint32, tid int32, count int32, at time.Time) *Item {
	def := Item2.GetID(tid)
	it := &Item{
		ID:          id,
		TID:         tid,
		Count:       count,
		AcquireTime: at.Unix(),
	}

	if def.GetTimeLimitType() == enumType.ItemTimeLimitType_Duration {
		it.TimeLimit = at.Unix() + def.GetTimeLimitDuration()
	} else if def.GetTimeLimitType() == enumType.ItemTimeLimitType_Date {
		it.TimeLimit = def.GetTimeLimitTime().Unix()
	}

	return it
}

type Backpack struct {
	user       module.UserI
	slotItems  []map[uint32]*Item  // 道具
	items      map[int32][]*Item   // tid -> []item
	dirtyItems map[uint32]struct{} // 已更新的物品
	*module.ModuleSaveBase
}

func (bp *Backpack) ModuleType() module.ModuleType {
	return module.Backpack
}

func (bp *Backpack) Init(fields map[string]*client.Field) error {
	for i := 0; i < slotCount; i++ {
		bp.slotItems[i] = map[uint32]*Item{}
		slotIdx := slotName(i)
		field, ok := fields[slotIdx]
		if !ok || len(field.GetBlob()) == 0 {
			bp.SetDirty(i)
		} else {
			err := json.Unmarshal(field.GetBlob(), &bp.slotItems[i])
			if err != nil {
				return fmt.Errorf("unmarshal:%s %s", string(field.GetBlob()), err)
			}

			for _, it := range bp.slotItems[i] {
				if s, ok := bp.items[it.TID]; ok {
					s = append(s, it)
					bp.items[it.TID] = s
				} else {
					bp.items[it.TID] = []*Item{it}
				}
			}
		}
	}
	return nil
}

func (bp *Backpack) ReadOut() *module.ReadOutCommand {
	cmd := &module.ReadOutCommand{
		Table:  bp.ModuleType().String(),
		Key:    bp.user.GetIDStr(),
		Module: bp,
	}

	cmd.Fields = make([]string, 0, slotCount)
	for i := 0; i < slotCount; i++ {
		cmd.Fields = append(cmd.Fields, slotName(i))
	}
	return cmd
}

func (bp *Backpack) WriteCommand(fields map[interface{}]struct{}) *module.WriteBackCommand {
	wdc := &module.WriteBackCommand{
		Table:  dbTable,
		Key:    bp.user.GetIDStr(),
		Module: bp,
	}
	wdc.Fields = make([]*module.WriteBackFiled, 0, len(fields))

	for field := range fields {
		switch slotIdx := field.(type) {
		case int:
			if slotIdx < 0 || slotIdx >= slotCount {
				panic("invalid dirty slotIdx")
			}

			bytes, err := json.Marshal(&bp.slotItems[slotIdx])
			if err != nil {
				zaplogger.GetSugar().Errorf("user(%s, %d) backpack write-back: marshal slotIdx %d: %s",
					bp.user.GetUserID(), bp.user.GetID(), slotIdx, err)
				return nil
			}

			wdc.Fields = append(wdc.Fields, &module.WriteBackFiled{
				Name:  slotName(slotIdx),
				Value: bytes,
			})
		default:
			panic("invalid dirty field")
		}
	}

	return wdc
}

func (bp *Backpack) FlushDirtyToClient() {
	if len(bp.dirtyItems) > 0 {
		msg := &cs_msg.BackpackSyncToC{
			All:      proto.Bool(false),
			Entities: make([]*cs_msg.BackpackEntity, 0, len(bp.dirtyItems)),
		}

		for id := range bp.dirtyItems {
			slotIdx := calcSlotIdx(id)
			it := bp.slotItems[slotIdx][id]
			if it == nil {
				msg.Entities = append(msg.Entities, &cs_msg.BackpackEntity{
					Id:    proto.Uint32(id),
					Count: proto.Int32(0),
				})
			} else {
				msg.Entities = append(msg.Entities, packItem(it))
			}

		}
		bp.dirtyItems = map[uint32]struct{}{}
		bp.user.Post(msg)

	}
}

func (bp *Backpack) FlushAllToClient(seqNo ...uint32) {
	msg := &cs_msg.BackpackSyncToC{
		All:      proto.Bool(true),
		Entities: make([]*cs_msg.BackpackEntity, 0, len(bp.items)),
	}

	for _, slot := range bp.slotItems {
		for _, it := range slot {
			msg.Entities = append(msg.Entities, packItem(it))
		}
	}

	bp.dirtyItems = map[uint32]struct{}{}
	bp.user.Post(msg)
}

func (bp *Backpack) Tick(time.Time) {}

func (bp *Backpack) AddItem(it *Item) {
	slotIdx := calcSlotIdx(it.ID)
	if _, ok := bp.slotItems[slotIdx][it.ID]; ok {
		return
	}

	bp.slotItems[slotIdx][it.ID] = it
	if items, ok := bp.items[it.TID]; ok {
		items = append(items, it)
		bp.items[it.TID] = items
	} else {
		bp.items[it.TID] = []*Item{it}
	}

	bp.dirtyItems[it.ID] = struct{}{}
	bp.SetDirty(slotIdx)
}

func (bp *Backpack) GetItem(id uint32) *Item {
	slotIdx := calcSlotIdx(id)
	return bp.slotItems[slotIdx][id]
}

func (bp *Backpack) GetItemsByTID(tid int32) []*Item {
	return bp.items[tid]
}

func (bp *Backpack) GetItemCountByTID(tid int32) (count int32) {
	for _, it := range bp.GetItemsByTID(tid) {
		count += it.Count
	}
	return
}

func (bp *Backpack) AddItemCount(id uint32, dt int32) {
	slotIdx := calcSlotIdx(id)
	it := bp.slotItems[slotIdx][id]

	it.Count += dt
	if it.Count <= 0 {
		bp.RemItem(id)
		return
	}

	bp.dirtyItems[it.ID] = struct{}{}
	bp.SetDirty(slotIdx)
}

func (bp *Backpack) RemItem(id uint32) {
	slotIdx := calcSlotIdx(id)
	if it, ok := bp.slotItems[slotIdx][id]; ok {
		delete(bp.slotItems[slotIdx], id)
		items := bp.items[it.TID]
		if len(items) <= 1 {
			delete(bp.items, it.TID)
		} else {
			for idx, it2 := range items {
				if it2.ID == id {
					items = append(items[:idx], items[idx+1:]...)
					bp.items[it.TID] = items
				}
			}
		}
		bp.dirtyItems[id] = struct{}{}
		bp.SetDirty(slotIdx)
	}

}

func packItem(item *Item) *cs_msg.BackpackEntity {
	return &cs_msg.BackpackEntity{
		Tid:         proto.Int32(item.TID),
		Id:          proto.Uint32(item.ID),
		Count:       proto.Int32(item.Count),
		ExpireTime:  proto.Int64(item.TimeLimit),
		AcquireTime: proto.Int64(item.AcquireTime),
	}
}

func init() {
	module.RegisterModule(module.Backpack, func(user module.UserI) (m module.ModuleI) {
		ret := &Backpack{
			user:       user,
			slotItems:  make([]map[uint32]*Item, slotCount),
			items:      map[int32][]*Item{},
			dirtyItems: map[uint32]struct{}{},
		}

		ret.ModuleSaveBase = module.NewModuleSaveBase(ret)
		return ret
	})
}
