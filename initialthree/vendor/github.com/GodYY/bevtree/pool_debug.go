// +build debug

package bevtree

import "sync/atomic"

type poolDebug struct {
	totalGet int64
	totalPut int64
}

func (p *poolDebug) get() {
	atomic.AddInt64(&p.totalGet, 1)
}

func (p *poolDebug) put() {
	atomic.AddInt64(&p.totalPut, 1)
}

func (p *poolDebug) getTotal() int64 {
	return atomic.LoadInt64(&p.totalGet)
}

func (p *poolDebug) putTotal() int64 {
	return atomic.LoadInt64(&p.totalPut)
}

var _poolDebug poolDebug

func (p *pool) get() interface{} {
	_poolDebug.get()
	return p.p.Get()
}

func (p *pool) put(i interface{}) {
	_poolDebug.put()
	p.p.Put(i)
}
