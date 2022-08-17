package scarsIngrain

import (
	"initialthree/node/node_game/module"
	"time"
)

func (this *ScarsIngrain) tryClock(name string, nowUnix int64, fn func(nowUnix int64)) {
	if timestamp, ok := this.timedata[name]; ok {
		if nowUnix >= timestamp {
			fn(nowUnix)
		}
	} else {
		fn(nowUnix)
	}
}

func (this *ScarsIngrain) clockTimer() {
	now := time.Now().Unix()

	this.tryClock(module.WeeklyTimeName, now, this.weeklyClock)
}

func (this *ScarsIngrain) weeklyClock(now int64) {
	this.dirty = true

	this.timedata[module.WeeklyTimeName] = module.CalWeeklyTime().Unix()
	this.SetDirty(timeField)
}
