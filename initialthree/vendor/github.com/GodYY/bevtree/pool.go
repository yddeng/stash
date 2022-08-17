package bevtree

import (
	"sync"
)

// A packaging of go built-in sync.Pool.
type pool struct {
	p sync.Pool
}

func newPool(new func() interface{}) *pool {
	p := &pool{}
	p.p.New = new
	return p
}

type taskPool struct {
	pool
}

func newTaskPool(new func() Task) *taskPool {
	p := &taskPool{}
	p.p.New = func() interface{} { return new() }
	return p
}

func (p *taskPool) get() Task {
	return p.pool.get().(Task)
}

func (p *taskPool) put(t Task) {
	p.pool.put(t)
}
