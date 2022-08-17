package Product

import (
	"fmt"
	"initialthree/node/common/enumType"
	"sync/atomic"
	"time"
)

type LimitDate struct {
	StartTime int64
	EndTime   int64
}

var limitDate atomic.Value

func (this *Product) LimitDate() *LimitDate {
	return limitDate.Load().(map[*Product]*LimitDate)[this]
}

func (this *Table) AfterLoad() {
	idMap := GetIDMap()
	var y, mon, d, h, min int
	var err error
	dates := make(map[*Product]*LimitDate)
	for _, v := range idMap {
		if v.ProductLimitTypeEnum == enumType.ProductLimitType_Date {
			if v.LimitDateStart == "" || v.LimitDateEnd == "" {
				panic(fmt.Sprintf("Product %d is ProductLimitType_Date,but limit date is nil", v.ID))
			}

			ld := new(LimitDate)
			if _, err = fmt.Sscanf(v.LimitDateStart, "%d-%d-%d %d:%d", &y, &mon, &d, &h, &min); err != nil {
				panic(err)
			}
			t := time.Date(y, time.Month(mon), d, h, min, 0, 0, time.Local)
			ld.StartTime = t.Unix()

			if _, err = fmt.Sscanf(v.LimitDateEnd, "%d-%d-%d %d:%d", &y, &mon, &d, &h, &min); err != nil {
				panic(err)
			}
			t = time.Date(y, time.Month(mon), d, h, min, 0, 0, time.Local)
			ld.EndTime = t.Unix()

			dates[v] = ld
		}
	}
	limitDate.Store(dates)
}
