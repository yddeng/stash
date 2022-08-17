package timeDisposal

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestGetDailyRefreshTime(t *testing.T) {
	//t1 := GetNextDailyRefreshTime(12, 0, 0)
	//t2 := GetNextDailyRefreshTime(11, 0, 0)
	//
	//fmt.Println(t1.String(), t2.String())
}

func TestGetWeeklyRefreshTime(t *testing.T) {
	t1 := CalcLatestWeekTimeAfter(1, 12, 0, 0)
	fmt.Println(t1.String())

	t2 := CalcLatestWeekTimeAfter(2, 11, 0, 0)
	fmt.Println(t2.String())

	t3 := CalcLatestWeekTimeAfter(2, 13, 0, 0)
	fmt.Println(t3.String())

	t4 := CalcLatestWeekTimeAfter(2, 11, 0, 0)
	fmt.Println(t4.String())
}

func TestSubDuration(t *testing.T) {
	s := time.Now().Add(time.Second)
	fmt.Println(s.UnixNano())
	e := time.Unix(0, s.UnixNano())
	fmt.Println(e.UnixNano())
	d := SubDuration(s.UnixNano())
	fmt.Println(d.Nanoseconds(), d.String(), d.Milliseconds())
}

func TestTimer(t *testing.T) {
	initTimerMgr()

	onceExpired := sync.WaitGroup{}
	onceExpired.Add(1)
	tm.createOnceTimer(5*time.Second, func(t *Timer, ctx interface{}) {
		fmt.Println("once timer:", Now())
		onceExpired.Done()
	}, nil, nil)

	repeatTimer, _ := tm.createRepeatTimer(1*time.Second, func(t *Timer, ctx interface{}) {
		fmt.Println("repeat timer:", Now())
	}, nil, nil)

	taExpired := sync.WaitGroup{}
	taExpired.Add(1)
	timeAt := Now().Add(10 * time.Second)
	tm.createTimerAt(timeAt, func(t *Timer, ctx interface{}) {
		fmt.Println("time at:", timeAt, "expired at:", Now())
		taExpired.Done()
	}, nil, nil)

	onceExpired.Wait()
	repeatTimer.Stop()

	taExpired.Wait()
}

func TestCalcLatestTimeAfter(t *testing.T) {
	d := CalcLatestTimeAfter(18, 0, 0)
	fmt.Println(d.String())
}

func TestCalcLatestMonthlyTimeAfter(t *testing.T) {
	d := CalcLatestMonthlyTimeAfter(12, 18, 0, 0)
	fmt.Println(d.String())
}

func TestNowInWeekdayWithRefreshTime(t *testing.T) {
	b := NowInWeekdayWithRefreshTime(1, 15, 0, 0)
	fmt.Println(b)
	b = NowInWeekdayWithRefreshTime(1, 16, 0, 0)
	fmt.Println(b)
	b = NowInWeekdayWithRefreshTime(7, 15, 0, 0)
	fmt.Println(b)
	b = NowInWeekdayWithRefreshTime(7, 16, 0, 0)
	fmt.Println(b)
	b = NowInWeekdayWithRefreshTime(2, 15, 0, 0)
	fmt.Println(b)
	b = NowInWeekdayWithRefreshTime(2, 16, 0, 0)
	fmt.Println(b)
}
