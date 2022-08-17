package sign

import (
	"initialthree/node/common/enumType"
	"initialthree/node/node_game/module"
	Sign2 "initialthree/node/table/excel/DataTable/Sign"
	"time"
)

func (this *UserSign) tryClock(name string, nowUnix int64, fn func(nowUnix int64)) {
	if timestamp, ok := this.timedata[name]; ok {
		if nowUnix >= timestamp {
			fn(nowUnix)
		}
	} else {
		fn(nowUnix)
	}
}

func (this *UserSign) clockTimer() {
	now := time.Now().Unix()
	this.tryClock(module.MonthlyTimeName, now, this.monthlyClock)
}

// 月更新
func (this *UserSign) monthlyClock(now int64) {
	for id := range this.data {
		def := Sign2.GetID(id)
		if def.TypeEnum == enumType.SignType_Month {
			this.Reset(id)
		}
	}

	this.timedata[module.MonthlyTimeName] = module.CalMonthlyTime().Unix()
	this.SetDirty(timeField)
}
