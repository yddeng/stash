package bevtree

// List element pool, used to cache list element.
type elemPool struct {
	p *pool
}

func newElemPool() *elemPool {
	p := new(elemPool)
	p.p = newPool(func() interface{} { return new(element) })
	return p
}

func (p *elemPool) getElement() *element {
	return p.p.get().(*element)
}

func (p *elemPool) putElement(e *element) {
	p.p.put(e)
}

var _elemPool = newElemPool()

// List element.
type element struct {
	l     *list
	Value interface{}
	prev  *element
	next  *element
}

func (e *element) getNext() *element {
	if p := e.next; e.l != nil && p != &e.l.root {
		return p
	}
	return nil
}

func (e *element) getPrev() *element {
	if p := e.prev; e.l != nil && p != &e.l.root {
		return p
	}
	return nil
}

func (e *element) reset() {
	e.l = nil
	e.Value = nil
	e.prev = nil
	e.next = nil
}

// List is a two-way list, like the go built-in container/list.List,
// but it use element pool to maintains elements.
type list struct {
	root element
	len  int
}

func newList() *list {
	return new(list).init()
}

func (l *list) init() *list {
	if l.len > 0 {
		elem := l.front()
		for elem != nil {
			next := elem.getNext()
			elem.reset()
			_elemPool.putElement(elem)
			elem = next
		}
	}

	l.len = 0
	l.root.prev = &l.root
	l.root.next = &l.root

	return l
}

func (l *list) lazyInit() {
	if l.root.next == nil {
		l.init()
	}
}

func (l *list) getLen() int { return l.len }

func (l *list) front() *element {
	if l.len == 0 {
		return nil
	}

	return l.root.next
}

func (l *list) back() *element {
	if l.len == 0 {
		return nil
	}

	return l.root.prev
}

func (l *list) pushFront(v interface{}) *element {
	l.lazyInit()
	return l.insertValue(v, &l.root)
}

func (l *list) pushBack(v interface{}) *element {
	l.lazyInit()
	return l.insertValue(v, l.root.prev)
}

func (l *list) insertBefore(v interface{}, e *element) *element {
	if e.l != l {
		return nil
	}
	return l.insertValue(v, e.prev)
}

func (l *list) insertAfter(v interface{}, e *element) *element {
	if e.l != l {
		return nil
	}
	return l.insertValue(v, e)
}

func (l *list) insert(e, at *element) *element {
	e.prev = at
	e.next = at.next
	e.prev.next = e
	e.next.prev = e
	e.l = l
	l.len++
	return e
}

func (l *list) insertValue(v interface{}, at *element) *element {
	elem := _elemPool.getElement()
	elem.Value = v
	return l.insert(elem, at)
}

func (l *list) Move(e, at *element) *element {
	if e == at {
		return e
	}

	e.prev.next = e.next
	e.next.prev = e.prev
	e.prev = at
	e.next = at.next
	e.prev.next = e
	e.next.prev = e
	return e
}

func (l *list) moveToFront(e *element) {
	if e.l != l || e == l.root.next {
		return
	}

	l.Move(e, &l.root)
}

func (l *list) moveToBack(e *element) {
	if e.l != l || e == l.root.prev {
		return
	}

	l.Move(e, l.root.prev)
}

func (l *list) moveBefore(e, mark *element) {
	if e.l != l || e == mark || mark.l != l {
		return
	}

	l.Move(e, mark.prev)
}

func (l *list) moveAfter(e, mark *element) {
	if e.l != l || e == mark || mark.l != l {
		return
	}

	l.Move(e, mark)
}

func (l *list) _remove(e *element) {
	e.prev.next = e.next
	e.next.prev = e.prev
	e.prev = nil
	e.next = nil
	e.l = nil
	e.Value = nil
	l.len--
	_elemPool.putElement(e)
}

func (l *list) remove(e *element) interface{} {
	if e == nil || e.l != l {
		return nil
	}

	v := e.Value
	l._remove(e)
	return v
}
