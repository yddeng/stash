package MaterialDungeon

import (
	"initialthree/node/common/timeDisposal"
	"initialthree/node/table/excel/ConstTable/Global"
	"time"
)

func (this *MaterialDungeon) DungeonOpen() bool {
	rt := Global.Get().GetDailyRefreshTime()
	for _, v := range this.OpenTimeArray {
		bt := timeDisposal.CalcLatestWeekTimeAfter(int(v.Weekday), rt.Hour, rt.Minute, 0).AddDate(0, 0, -7)
		et := bt.Add(time.Hour * 24)
		now := timeDisposal.Now()
		if now.After(bt) && now.Before(et) {
			return true
		}
	}
	return false
}
