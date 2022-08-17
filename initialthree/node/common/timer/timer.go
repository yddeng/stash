package timer

import (
	"sync"
)

/*
 非线程安全
 单位 秒
*/
var (
	tMgr *TimeMgr
	once = sync.Once{}
)

type TimeMgr struct {
	minHeap  *Heap
	lastTime int64 // s
}

func (mgr *TimeMgr) addTimer(t *Timer) bool {
	// 已经过期
	lastTime := mgr.lastTime
	if lastTime != 0 && t.expiredTime <= lastTime {
		t.do(lastTime)
		return false
	}
	mgr.minHeap.Push(t)
	return true
}

type Timer struct {
	f           func(now int64)
	expiredTime int64
	stop        bool
	tMgr        *TimeMgr
}

func (t *Timer) Less(e Element) bool {
	return t.expiredTime < e.(*Timer).expiredTime
}

func (t *Timer) do(now int64) {
	if t.stop {
		return
	}
	t.f(now)
}

func (t *Timer) Stop() {
	t.stop = true
	t.tMgr.minHeap.Remove(t)
}

func (t *Timer) Reset(expiredTime int64) {
	t.stop = false
	t.expiredTime = expiredTime
	t.tMgr.minHeap.Remove(t)
	t.tMgr.addTimer(t)
}

func NewTimerMgr(startTime int64) *TimeMgr {
	return &TimeMgr{
		minHeap:  NewHeap(),
		lastTime: startTime,
	}
}

func (mgr *TimeMgr) Loop(unixNow int64) {
	mgr.lastTime = unixNow
	e := mgr.minHeap.Peek()
	for e != nil && e.(*Timer).expiredTime <= unixNow {
		mgr.minHeap.Pop()
		t := e.(*Timer)
		t.do(unixNow)
		e = mgr.minHeap.Peek()
	}
}

func (mgr *TimeMgr) Once(expiredTime int64, f func(now int64)) *Timer {
	t := &Timer{
		f:           f,
		expiredTime: expiredTime,
		tMgr:        mgr,
	}

	mgr.addTimer(t)
	return t
}

func Once(expiredTime int64, f func(now int64)) *Timer {
	once.Do(func() {
		tMgr = &TimeMgr{
			minHeap: NewHeap(),
		}
	})
	return tMgr.Once(expiredTime, f)
}

func Loop(now int64) {
	if tMgr == nil {
		return
	}
	tMgr.Loop(now)
}
