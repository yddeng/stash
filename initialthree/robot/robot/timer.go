package robot

import (
	externlTimer "initialthree/pkg/timer"
	"initialthree/robot/types"
	"time"
)

type TimerID = types.TimerID

type timer struct {
	*externlTimer.Timer
	id  TimerID
	cb  func(types.Robot, interface{})
	ctx interface{}
}

func (r *Robot) AddTimer(id TimerID, d time.Duration, ctx interface{}, cb func(r RobotI, ctx interface{})) {
	if _, ok := r.timers[id]; ok {
		r.Panicf("add duplicate timer \"%s\"", id.String())
	}

	timer := &timer{id: id, cb: cb, ctx: ctx}
	timer.Timer = externlTimer.Once(d, r.onTimer, timer)
	r.timers[id] = timer
}

func (r *Robot) RemTimer(id TimerID) {
	if t, ok := r.timers[id]; ok {
		delete(r.timers, id)
		t.Cancel()
	}
}

func (r *Robot) onTimer(t *externlTimer.Timer, ctx interface{}) {
	r.PostNoWait(r.onTimer_, ctx.(*timer))
}

func (r *Robot) onTimer_(t *timer) {
	if timer, ok := r.timers[t.id]; ok && timer == t {
		delete(r.timers, t.id)
		timer.cb(r, timer.ctx)
	}
}

func (r *Robot) clearTimer() {
	for k := range r.timers {
		r.RemTimer(k)
	}
}
