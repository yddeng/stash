// +build !debug

package bevtree

func (p *pool) get() interface{} {
	return p.p.Get()
}

func (p *pool) put(i interface{}) {
	p.p.Put(i)
}
