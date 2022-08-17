package Global

import (
	"initialthree/node/common/timeDisposal"
	"sync/atomic"
)

func Get() *Global { return GetID(1) }

func (this *Global) GetDailyRefreshTime() DailyRefreshTime_ { return *this.DailyRefreshTimeStruct }

func (this *Global) GetWeeklyRefreshTime() *WeeklyRefreshTimeStruct {
	return weekly.Load().(*WeeklyRefreshTimeStruct)
}

func (this *Global) GetMonthlyRefreshTime() *MonthlyRefreshTimeStruct {
	return monthly.Load().(*MonthlyRefreshTimeStruct)
}

type WeeklyRefreshTimeStruct struct {
	Weekday int
	DailyRefreshTime_
}

type MonthlyRefreshTimeStruct struct {
	Day int32
	DailyRefreshTime_
}

var weekly atomic.Value
var monthly atomic.Value

func (this *Table) AfterLoad() {

	daily := Get().GetDailyRefreshTime()

	weekday := Get().WeeklyRefreshTime
	wd := &WeeklyRefreshTimeStruct{}
	wd.Weekday = timeDisposal.ParseWeekday(weekday)
	wd.Hour = daily.Hour
	wd.Minute = daily.Minute
	weekly.Store(wd)

	day := Get().MonthlyRefreshTime
	md := &MonthlyRefreshTimeStruct{
		Day: day,
		DailyRefreshTime_: DailyRefreshTime_{
			Hour:   daily.Hour,
			Minute: daily.Minute,
		},
	}
	monthly.Store(md)
}
