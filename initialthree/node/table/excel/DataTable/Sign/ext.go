package Sign

import (
	"sync/atomic"
	"time"
)

type LimitDate struct {
	StartTime int64
	EndTime   int64
}

var limitDate atomic.Value

func (this *Sign) LimitDate() *LimitDate {
	return limitDate.Load().(map[*Sign]*LimitDate)[this]
}

func (this *Table) AfterLoad() {
	idMap := GetIDMap()
	dates := make(map[*Sign]*LimitDate)
	for _, v := range idMap {
		ld := new(LimitDate)

		if v.StartTime != "" {
			t, err := time.ParseInLocation("2006-01-02 15:04:05", v.StartTime, time.Local)
			if err != nil {
				panic(err)
			}
			ld.StartTime = t.Unix()
		}

		if v.EndTime != "" {
			t, err := time.ParseInLocation("2006-01-02 15:04:05", v.EndTime, time.Local)
			if err != nil {
				panic(err)
			}
			ld.EndTime = t.Unix()
		}

		dates[v] = ld
	}
	limitDate.Store(dates)
}
