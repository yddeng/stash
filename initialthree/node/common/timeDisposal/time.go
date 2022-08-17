package timeDisposal

import (
	"errors"
	"fmt"
	flyfish "github.com/sniperHW/flyfish/client"
	time2 "initialthree/node/common/omm/time"
	"time"
)

var WeekDay = map[string]int{
	"Monday":    1,
	"Tuesday":   2,
	"Wednesday": 3,
	"Thursday":  4,
	"Friday":    5,
	"Saturday":  6,
	"Sunday":    7,
}

var offset time.Duration

func Init(cli *flyfish.Client) error {
	if cli == nil {
		return errors.New("flyfish client nil")
	}

	time2.RegisterUpdateNotify(cli, func(offset2 time.Duration) {
		offset = offset2
	})

	return nil
}

func Now() time.Time {
	now := time.Now()

	if offset != 0 {
		now = now.Add(offset)
	}

	return now
}

func NowUnix() int64 {
	return Now().Unix()
}

func NowUnixNano() int64 {
	return Now().UnixNano()
}

func NowUnixMS() int64 {
	return NowUnixNano() / int64(time.Millisecond)
}

// s.UnixNano()
func TimeStamp2Time(ts int64) time.Time {
	return time.Unix(0, ts)
}

// 现在时间到 时间戳的 间隔
func SubDuration(timestamp int64) time.Duration {
	t := TimeStamp2Time(timestamp)
	return t.Sub(time.Now())
}

func getBaseTime(bases ...time.Time) time.Time {
	if len(bases) > 0 {
		return bases[0]
	}

	return Now()
}

// 计算今天
func TodayTime(hour, min, sec int32, base ...time.Time) time.Time {
	now := getBaseTime(base...)
	return time.Date(now.Year(), now.Month(), now.Day(), int(hour), int(min), int(sec), 0, time.Local)
}

// 计算基于某时刻之后的最近一个时刻
func CalcLatestTimeAfter(hour, min, sec int32, base ...time.Time) time.Time {
	now := getBaseTime(base...)

	t := time.Date(now.Year(), now.Month(), now.Day(), int(hour), int(min), int(sec), 0, time.Local)
	if now.Before(t) {
		return t
	}

	return t.AddDate(0, 0, 1)
}

func CalcLatestWeekTimeAfter(weekday int, hour, min, sec int32, base ...time.Time) time.Time {
	now := getBaseTime(base...)

	day := now.Day()
	offset := weekday - ConvertTimeWeekday(now.Weekday())
	if offset != 0 {
		day += offset
	}

	rt := time.Date(now.Year(), now.Month(), day, int(hour), int(min), int(sec), 0, time.Local)
	if now.Before(rt) {
		return rt
	}

	return rt.AddDate(0, 0, 7)
}

func CalcLatestMonthlyTimeAfter(day int, hour, min, sec int32, base ...time.Time) time.Time {
	now := getBaseTime(base...)

	t := time.Date(now.Year(), now.Month(), day, int(hour), int(min), int(sec), 0, time.Local)
	if now.Before(t) {
		return t
	}

	return t.AddDate(0, 1, 0)
}

// 涉及到系统开启时间， 比如某系统周一开启，每日8点刷新 实际开始时间为周一8点到周二8点
func NowInWeekdayWithRefreshTime(weekday int, hour, min, sec int32, base ...time.Time) bool {
	now := getBaseTime(base...)

	nowWeekday := ConvertTimeWeekday(now.Weekday())
	if nowWeekday == weekday && now.After(time.Date(now.Year(), now.Month(), now.Day(), int(hour), int(min), int(sec), 0, time.Local)) {
		return true
	}

	nn := nowWeekday - 1
	if nn < 1 {
		nn += 7
	}
	if nn == weekday && now.Before(time.Date(now.Year(), now.Month(), now.Day(), int(hour), int(min), int(sec), 0, time.Local)) {
		return true
	}

	return false
}

// 将当前周字符转换为 （1-7）
func ParseWeekday(weekday string) int {
	if i, ok := WeekDay[weekday]; ok {
		return i
	} else {
		panic(fmt.Sprintf("weekday %s failed", weekday))
	}
	return 0
}

// 将time.Weekyday(0 - 6)，转换为(1-7)
func ConvertTimeWeekday(wd time.Weekday) int {
	if wd == time.Sunday {
		return 7
	}

	return int(wd)
}

// 将1-7转换为time.Weekday
func ConvertWeekday(wd int) time.Weekday {
	if wd < 1 || wd > 7 {
		panic(errors.New("invalid weekday"))
	}

	if wd == 7 {
		return time.Sunday
	}

	return time.Weekday(wd)
}
