package MainChapter

import (
	"errors"
	"sync/atomic"
	"time"

	MainChapterTable "initialthree/node/table/excel/DataTable/MainChapter"
)

type activityTime struct {
	c     *MainChapter
	start time.Time
	end   time.Time
}

var (
	activityTimeValue atomic.Value
)

func init() {
	activityTimeValue.Store(&activityTime{c: nil})
}

func Get() *MainChapter { return GetID(1) }

func (this *MainChapter) GetActivityChapterID() int32 { return this.SpecialUnlockChapter }

func (this *MainChapter) IsActivityChapterOpen(t time.Time) bool {
	activityTime := activityTimeValue.Load().(*activityTime)

	if activityTime.c == this {
		return !t.Before(activityTime.start) && t.Before(activityTime.end)
	} else {
		st := this.SpecialUnlockStarDateStruct
		startTime := time.Date(int(st.Year), time.Month(st.Month), int(st.Day), int(st.Hour), int(st.Min), 0, 0, time.Local)

		et := this.SpecialUnlockEndDateStruct
		endTime := time.Date(int(et.Year), time.Month(et.Month), int(et.Day), int(et.Hour), int(et.Min), 0, 0, time.Local)

		return !t.Before(startTime) && t.Before(endTime)
	}
}

//func (this *Table) BeforeLoad() {
//	activityTimeValue.Store(invalidActivityTime)
//}

func (this *Table) AfterLoad() {
	c := Get()

	st := c.SpecialUnlockStarDateStruct
	et := c.SpecialUnlockEndDateStruct

	activityTime := &activityTime{
		c:     c,
		start: time.Date(int(st.Year), time.Month(st.Month), int(st.Day), int(st.Hour), int(st.Min), 0, 0, time.Local),
		end:   time.Date(int(et.Year), time.Month(et.Month), int(et.Day), int(et.Hour), int(et.Min), 0, 0, time.Local),
	}

	activityTimeValue.Store(activityTime)
}

func (this *Table) AfterLoadAll() {
	c := Get()

	if c.SpecialUnlockChapter != 0 && MainChapterTable.GetID(c.SpecialUnlockChapter) == nil {
		panic(errors.New("special unlock chapter not exist"))
	}
}
