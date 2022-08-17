package bigSecret

import (
	"initialthree/node/node_game/module"
	"initialthree/node/table/excel/ConstTable/BigSecret"
	"time"
)

func (this *BigSecretDungeon) tryClock(name string, nowUnix int64, fn func(nowUnix int64)) {
	if timestamp, ok := this.timedata[name]; ok {
		if nowUnix >= timestamp {
			fn(nowUnix)
		}
	} else {
		fn(nowUnix)
	}
}

func (this *BigSecretDungeon) clockTimer() {
	now := time.Now().Unix()

	this.tryClock(module.DailyTimeName, now, this.dailyClock)
}

func (this *BigSecretDungeon) dailyClock(now int64) {
	year, mon, day := time.Now().Date()
	nowDay := time.Date(year, mon, day, 0, 0, 0, 0, time.Local).Unix()
	if this.resetDate == nowDay {
		return
	}

	def := BigSecret.GetID(1)

	this.data.KeyCount += 1
	if this.data.KeyCount > def.MaxKeyCount {
		this.data.KeyCount = def.MaxKeyCount
	}
	this.data.WeaknessRefreshTimes = 0
	this.dirtyData = true
	this.SetDirty(dataField)

	this.timedata[module.DailyTimeName] = module.CalDailyTime().Unix()
	this.SetDirty(timeField)
}
