package Query

import (
	"container/list"
	"initialthree/node/node_game/module"
	"time"
)

var (
	roleHotLru = New(100 + 500)
	timeout    = time.Minute * 5
)

type Cache struct {
	maxEntries int
	ll         *list.List
	cache      map[interface{}]*list.Element
}

type entry struct {
	key    uint64
	value  *roleModuleData
	expire time.Time
}

type roleModuleData struct {
	modules map[module.ModuleType]module.ModuleI
}

func New(maxEntries int) *Cache {
	return &Cache{
		maxEntries: maxEntries,
		ll:         list.New(),
		cache:      make(map[interface{}]*list.Element),
	}
}

func (this *Cache) Get(key uint64) (*roleModuleData, bool) {
	defer this.doRemoveOldest()

	if ele, hit := this.cache[key]; hit {
		this.ll.MoveToFront(ele)
		e := ele.Value.(*entry)
		if !e.expire.IsZero() {
			e.expire = time.Time{}
		}
		return e.value, true
	}
	return nil, false
}

func (this *Cache) Set(key uint64, value *roleModuleData) {
	defer this.doRemoveOldest()

	if ele, ok := this.cache[key]; ok {
		this.ll.MoveToFront(ele)
		e := ele.Value.(*entry)
		if !e.expire.IsZero() {
			e.expire = time.Time{}
		}
		e.value = value
		return
	}

	ele := this.ll.PushFront(&entry{key, value, time.Time{}})
	this.cache[key] = ele
}

func (this *Cache) doRemoveOldest() {
	length := this.ll.Len()
	if length <= this.maxEntries {
		return
	}

	offset := length
	now := time.Now()
	for ele := this.ll.Back(); ele != nil && offset > this.maxEntries; ele, offset = ele.Prev(), offset-1 {
		e := ele.Value.(*entry)
		if e.expire.IsZero() {
			e.expire = now.Add(timeout)
		} else if now.After(e.expire) {
			this.ll.Remove(ele)
			delete(this.cache, e.key)
		}
	}
}

func newRoleModuleData() *roleModuleData {
	return &roleModuleData{
		modules: map[module.ModuleType]module.ModuleI{},
	}
}
