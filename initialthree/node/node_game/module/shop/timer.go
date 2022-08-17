package shop

import (
	"initialthree/node/common/enumType"
	"initialthree/node/node_game/module"
	"initialthree/node/table/excel/DataTable/Product"
	"time"
)

func (this *ShopData) tryClock(name string, nowUnix int64, fn func(nowUnix int64)) {
	if timestamp, ok := this.timedata[name]; ok {
		if nowUnix >= timestamp {
			fn(nowUnix)
		}
	} else {
		fn(nowUnix)
	}
}

func (this *ShopData) clockTimer() {
	now := time.Now().Unix()
	this.tryClock(module.MonthlyTimeName, now, this.monthlyClock)
	this.tryClock(module.WeeklyTimeName, now, this.weeklyClock)
	this.tryClock(module.DailyTimeName, now, this.dailyClock)
}

// 日更新
func (this *ShopData) dailyClock(now int64) {
	// 商品次数
	this.refresh(enumType.ProductLimitType_Daily)

	// 商店次数
	this.clearShopRefreshTimes()
	this.SetDirty(shopField)

	this.timedata[module.DailyTimeName] = module.CalDailyTime().Unix()
	this.SetDirty(timeField)
}

// 周更新
func (this *ShopData) weeklyClock(now int64) {
	// 商品次数
	this.refresh(enumType.ProductLimitType_Weekly)

	this.timedata[module.WeeklyTimeName] = module.CalWeeklyTime().Unix()
	this.SetDirty(timeField)
}

// 月更新
func (this *ShopData) monthlyClock(now int64) {
	// 商品次数
	this.refresh(enumType.ProductLimitType_Monthly)

	this.timedata[module.MonthlyTimeName] = module.CalMonthlyTime().Unix()
	this.SetDirty(timeField)
}

func (this *ShopData) refresh(limit int32) {
	products := make([]int32, 0, len(this.products))
	for id := range this.products {
		def := Product.GetID(id)
		if def == nil {
			products = append(products, id)
		} else {
			if def.ProductLimitTypeEnum == limit {
				products = append(products, id)
			}
		}
	}
	this.clearProductBuyTimes(products)
	this.SetDirty(productField)
}
