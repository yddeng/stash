package Item

import (
	"fmt"
	"initialthree/node/common/enumType"
	"strconv"
	"sync/atomic"
	"time"
)

func (it *Item) IsTimeLimit() bool { return it.TimeLimitTypeEnum != 0 }

func (it *Item) GetTimeLimitType() int32 { return it.TimeLimitTypeEnum }

func (it *Item) GetTimeLimitDuration() int64 {
	if it.TimeLimitTypeEnum != enumType.ItemTimeLimitType_Duration {
		panic("TimeLimitType error")
	}

	d, ok := itemTimeLimitDuration.Load().(map[int32]int64)[it.ID]
	if !ok {
		panic("duration not found")
	}

	return d
}

func (it *Item) GetTimeLimitTime() time.Time {
	if it.TimeLimitTypeEnum != enumType.ItemTimeLimitType_Date {
		panic("TimeLimitType error")
	}

	t, ok := itemTimeLimitTime.Load().(map[int32]time.Time)[it.ID]
	if !ok {
		panic("time not found")
	}

	return t
}

func (it *Item) GetBackpackType() int32 {
	return it.TypeEnum
}

func parseTimeLimitDuration(timeLimit string) (int64, bool) {
	t, err := strconv.ParseInt(timeLimit, 10, 64)
	if err != nil {
		return 0, false
	}

	if t <= 0 {
		return 0, false
	}

	return t * 60, true
}

func parseTimeLimitTime(timeLimit string) (time.Time, bool) {
	var year, month, day, hour, min int

	n, err := fmt.Sscanf(timeLimit, "%d-%d-%d %d:%d", &year, &month, &day, &hour, &min)
	if n != 5 || err != nil || month < 1 || month > 12 {
		return time.Time{}, false
	}

	return time.Date(year, time.Month(month), day, hour, min, 0, 0, time.Local), true
}

var (
	itemTimeLimitDuration atomic.Value
	itemTimeLimitTime     atomic.Value
)

func (this *Table) AfterLoad() {
	items := this.indexID.Load().(map[int32]*Item)

	timeLimitDurations := make(map[int32]int64)
	timeLimitTimes := make(map[int32]time.Time)
	for k, v := range items {
		if v.TimeLimitTypeEnum == enumType.ItemTimeLimitType_Duration {
			d, ok := parseTimeLimitDuration(v.TimeLimit)
			if !ok {
				panic(fmt.Errorf("Item:%d invalid TimeLimit:%s", k, v.TimeLimit))
			}

			timeLimitDurations[k] = d
		} else if v.TimeLimitTypeEnum == enumType.ItemTimeLimitType_Date {
			t, ok := parseTimeLimitTime(v.TimeLimit)
			if !ok {
				panic(fmt.Errorf("Item:%d invalid TimeLimit:%s", k, v.TimeLimit))
			}

			timeLimitTimes[k] = t
		}
	}

	itemTimeLimitDuration.Store(timeLimitDurations)
	itemTimeLimitTime.Store(timeLimitTimes)
}
