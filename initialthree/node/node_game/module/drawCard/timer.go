package drawCard

import (
	"initialthree/node/common/timeDisposal"
	"initialthree/node/node_game/module"
	"initialthree/node/table/excel/ConstTable/Global"
	"time"
)

func (this *DrawCard) tryClock(name string, nowUnix int64, fn func(nowUnix int64)) {
	if timestamp, ok := this.timedata[name]; ok {
		if nowUnix >= timestamp {
			fn(nowUnix)
		}
	} else {
		fn(nowUnix)
	}
}

func (this *DrawCard) clockTimer() {
	now := time.Now().Unix()
	this.tryClock(module.DailyTimeName, now, this.dailyClock)
}

// 日更新
func (this *DrawCard) dailyClock(now int64) {
	for id := range this.dailyTimes {
		this.dailyTimesDirty[id] = struct{}{}
	}

	this.dailyTimes = map[int32]int32{}
	this.SetDirty(dailyTimesField)

	rt := Global.Get().GetDailyRefreshTime()
	dailyTime := timeDisposal.CalcLatestTimeAfter(rt.Hour, rt.Minute, 0)
	this.timedata[module.DailyTimeName] = dailyTime.Unix()
	this.SetDirty(timeField)
}
