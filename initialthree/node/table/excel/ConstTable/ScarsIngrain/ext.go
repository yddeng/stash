package ScarsIngrain

import (
	"initialthree/node/common/timeDisposal"
	"time"
)

func GetBeginTime() time.Time {
	bt := GetID(1).BeginTimeStruct
	return timeDisposal.CalcLatestWeekTimeAfter(int(bt.Weekly), bt.Hour, bt.Minute, 0)
}

func GetEndTime() time.Time {
	bt := GetID(1).EndTimeStruct
	return timeDisposal.CalcLatestWeekTimeAfter(int(bt.Weekly), bt.Hour, bt.Minute, 0)
}

func BossOpen(idx int) bool {
	var weekday, hour, minute int32
	if idx == 0 {
		bt := GetID(1).BossOpenTime1Struct
		weekday, hour, minute = bt.Weekly, bt.Hour, bt.Minute
	} else if idx == 1 {
		bt := GetID(1).BossOpenTime2Struct
		weekday, hour, minute = bt.Weekly, bt.Hour, bt.Minute
	} else if idx == 2 {
		bt := GetID(1).BossOpenTime3Struct
		weekday, hour, minute = bt.Weekly, bt.Hour, bt.Minute
	} else {
		return false
	}

	now := timeDisposal.Now()
	nowWeekday := int32(timeDisposal.ConvertTimeWeekday(now.Weekday()))
	nowH, nowM := int32(now.Hour()), int32(now.Minute())
	if nowWeekday > weekday {
		return true
	} else if nowWeekday == weekday && nowH > hour {
		return true
	} else if nowWeekday == weekday && nowH == hour && nowM >= minute {
		return true
	}

	return false
}
