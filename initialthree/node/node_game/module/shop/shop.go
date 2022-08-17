package shop

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	flyfish "github.com/sniperHW/flyfish/client"
	"initialthree/node/node_game/module"
	"initialthree/pkg/json"
	"initialthree/protocol/cs/message"
	"initialthree/zaplogger"
	"time"
)

const (
	shopField    = "shop"
	productField = "product"
	timeField    = "timedata"
)

var tableFields = []string{shopField, productField, timeField}

type ShopData struct {
	userI     module.UserI
	products  map[int32]int32 // id -> byt times
	shop      map[int32]int32 // id -> refresh times
	timedata  map[string]int64
	pDirty    map[int32]struct{}
	shopDirty map[int32]struct{}
	*module.ModuleSaveBase
}

func (this *ShopData) clearProductBuyTimes(products []int32) {
	for _, pid := range products {
		delete(this.products, pid)
		this.pDirty[pid] = struct{}{}
	}
}

func (this *ShopData) clearShopRefreshTimes() {
	for id := range this.shop {
		this.shopDirty[id] = struct{}{}
	}
	this.shop = map[int32]int32{}
}

func (this *ShopData) GetProductBuyTimes(id int32) int32 {
	return this.products[id]
}

func (this *ShopData) GetShopRefreshTimes(id int32) int32 {
	return this.shop[id]
}

func (this *ShopData) ShopRefresh(id int32, products []int32) {
	// 累加商店刷新次数
	times := this.shop[id]
	this.shop[id] = times + 1
	if _, ok := this.shopDirty[id]; !ok {
		this.shopDirty[id] = struct{}{}
	}
	this.SetDirty(shopField)

	// 商品已购次数清零
	this.clearProductBuyTimes(products)
	this.SetDirty(productField)

}

func (this *ShopData) ShopBuy(pid int32, count int32) {
	times := this.products[pid]
	this.products[pid] = times + count
	if _, ok := this.pDirty[pid]; !ok {
		this.pDirty[pid] = struct{}{}
	}
	this.SetDirty(productField)
}

func (this *ShopData) ModuleType() module.ModuleType {
	return module.Shop
}

func (this *ShopData) Init(fields map[string]*flyfish.Field) error {
	for _, name := range tableFields {
		field, ok := fields[name]
		if ok && len(field.GetBlob()) != 0 {
			var err error
			switch name {
			case shopField:
				err = json.Unmarshal(field.GetBlob(), &this.shop)
			case productField:
				err = json.Unmarshal(field.GetBlob(), &this.products)
			case timeField:
				err = json.Unmarshal(field.GetBlob(), &this.timedata)
			}
			if err != nil {
				zaplogger.GetSugar().Errorf("%s initShop name %s err %s ", this.userI.GetUserID(), name, err)
				return fmt.Errorf("unmarshal: %s", err)
			}
		}
	}
	return nil
}

func (this *ShopData) ReadOut() *module.ReadOutCommand {
	return &module.ReadOutCommand{
		Table:  this.ModuleType().String(),
		Key:    this.userI.GetIDStr(),
		Fields: tableFields,
		Module: this,
	}
}

func (this *ShopData) WriteCommand(fields map[interface{}]struct{}) *module.WriteBackCommand {
	cmd := &module.WriteBackCommand{
		Table:  this.ModuleType().String(),
		Key:    this.userI.GetIDStr(),
		Module: this,
		Fields: make([]*module.WriteBackFiled, 0, len(fields)),
	}

	for field := range fields {
		name := field.(string)
		var data []byte
		switch name {
		case shopField:
			data, _ = json.Marshal(this.shop)
		case productField:
			data, _ = json.Marshal(this.products)
		case timeField:
			data, _ = json.Marshal(this.timedata)
		default:
			continue
		}
		cmd.Fields = append(cmd.Fields, &module.WriteBackFiled{
			Name:  name,
			Value: data,
		})
	}

	return cmd
}

func (this *ShopData) Tick(now time.Time) {
	this.clockTimer()
}

func (this *ShopData) FlushDirtyToClient() {
	if len(this.pDirty) > 0 || len(this.shopDirty) > 0 {
		msg := &message.ShopSyncToC{
			IsAll:    proto.Bool(false),
			Shops:    make([]*message.Shop, 0, len(this.shopDirty)),
			Products: make([]*message.Product, 0, len(this.pDirty)),
		}

		for id := range this.pDirty {
			times := this.products[id]
			msg.Products = append(msg.Products, &message.Product{
				Id:              proto.Int32(id),
				AlreadyBuyTimes: proto.Int32(times),
			})
		}

		for id := range this.shopDirty {
			times := this.shop[id]
			msg.Shops = append(msg.Shops, &message.Shop{
				Id:                  proto.Int32(id),
				AlreadyRefreshTimes: proto.Int32(times),
			})
		}

		this.userI.Post(msg)
		this.pDirty = map[int32]struct{}{}
		this.shopDirty = map[int32]struct{}{}
	}
}

func (this *ShopData) FlushAllToClient(seqNo ...uint32) {
	msg := &message.ShopSyncToC{
		IsAll:    proto.Bool(true),
		Shops:    make([]*message.Shop, 0, len(this.shop)),
		Products: make([]*message.Product, 0, len(this.products)),
	}

	for id, times := range this.products {
		msg.Products = append(msg.Products, &message.Product{
			Id:              proto.Int32(id),
			AlreadyBuyTimes: proto.Int32(times),
		})
	}

	for id, times := range this.shop {
		msg.Shops = append(msg.Shops, &message.Shop{
			Id:                  proto.Int32(id),
			AlreadyRefreshTimes: proto.Int32(times),
		})
	}
	this.userI.Post(msg)
	this.pDirty = map[int32]struct{}{}
	this.shopDirty = map[int32]struct{}{}
}

/*
func (this *ShopData) ShopSyncToC() *message.ShopSyncToC {
	this.clock()

	return &message.ShopSyncToC{
		IsAll:    proto.Bool(true),
		Shops:    this.shopdata,
		Products: this.products,
	}
}
*/
func init() {
	module.RegisterModule(module.Shop, func(userI module.UserI) module.ModuleI {
		m := &ShopData{
			userI:     userI,
			shop:      map[int32]int32{},
			products:  map[int32]int32{},
			shopDirty: map[int32]struct{}{},
			pDirty:    map[int32]struct{}{},
			timedata:  map[string]int64{},
		}

		m.ModuleSaveBase = module.NewModuleSaveBase(m)
		return m
	})
}
