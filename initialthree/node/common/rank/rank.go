package rank

const realRankCount int = 10000000
const maxItemCount int = 10000000

type rankItemBlock struct {
	items    []node
	nextFree int
}

func (rb *rankItemBlock) get() *node {
	if rb.nextFree >= cap(rb.items) {
		return nil
	} else {
		item := &rb.items[rb.nextFree]
		rb.nextFree++
		return item
	}
}
func (rb *rankItemBlock) reset() {
	for i, _ := range rb.items {
		item := &rb.items[i]
		item.sl = nil
		for j := 0; j < maxLevel; j++ {
			item.links[j].skip = 0
			item.links[j].pnext, item.links[j].pprev = nil, nil
		}
	}
	rb.nextFree = 0
}
func newRankItemBlock() *rankItemBlock {
	return &rankItemBlock{
		items: make([]node, 10000),
	}
}

type rankItemPool struct {
	blocks   []*rankItemBlock
	nextFree int
}

func newRankItemPool() *rankItemPool {
	return &rankItemPool{
		blocks: []*rankItemBlock{newRankItemBlock()},
	}
}
func (rp *rankItemPool) reset() {
	for _, v := range rp.blocks {
		v.reset()
	}
	rp.nextFree = 0
}
func (rp *rankItemPool) get() *node {
	item := rp.blocks[rp.nextFree].get()
	if nil == item {
		block := newRankItemBlock()
		rp.blocks = append(rp.blocks, block)
		rp.nextFree++
		item = block.get()
	}
	item.sl = nil
	return item
}

type Rank struct {
	id2Item   map[uint64]*node
	spans     []*skiplists
	nextShink int
	cc        int
	itemPool  *rankItemPool
}

func NewRank() *Rank {
	return &Rank{
		id2Item:  map[uint64]*node{},
		spans:    make([]*skiplists, 0, 8192),
		itemPool: newRankItemPool(),
	}
}

func (r *Rank) Reset() {
	r.id2Item = map[uint64]*node{}
	r.spans = make([]*skiplists, 0, 8192)
	r.itemPool.reset()
}

// 数据少于 10000，准确。大于10000 ，模糊。
// 数据少于 100，百分比为名次。

//func (r *Rank) getRankPercent(item *node, idxInSpan int) int {
//	if len(r.spans) < 100 {
//		rank := r.getRank(item)
//		size := r.Size()
//		if size <= 100 {
//			return rank
//		} else {
//			return rank * 100 / size
//		}
//		//return 100 - idxInSpan*100/maxItemCount*len(r.spans)
//	} else {
//		return 100 - maxItemCount*item.sl.idx/(len(r.spans)-1)
//	}
//}

func (r *Rank) getRankPercentByItem(item *node) int {
	if len(r.spans) < 100 {
		rank := r.getRank(item)
		size := r.Size()
		if size <= 100 {
			return rank
		} else {
			perc := rank * 100 / size
			// 0 默认为 1
			if perc <= 0 {
				return 1
			} else if perc > 100 {
				return 100
			}
			return perc
		}
	} else {
		return (item.sl.idx + 1) * 100 / len(r.spans)
	}
}

func (r *Rank) GetRankPercent(id uint64) int {
	item := r.getRankItem(id)
	if nil == item {
		return -1
	} else {
		return r.getRankPercentByItem(item)
	}
}

func (r *Rank) getFrontSpanItemCount(item *node) int {
	c := 0
	i := 0
	for ; i < len(r.spans); i++ {
		v := r.spans[i]
		if item.sl == v {
			break
		} else {
			c += v.size
			if c >= realRankCount {
				break
			}
		}
	}

	if i < item.sl.idx {
		c += maxItemCount * (item.sl.idx - i)
	}

	return c
}

func (r *Rank) getRank(item *node) int {
	return r.getFrontSpanItemCount(item) + item.sl.GetNodeRank(item)
}

func (r *Rank) GetRank(id uint64) int {

	r.cc++
	defer func() {
		if r.cc%100 == 0 {
			r.shrink(nil)
		}
	}()

	item := r.getRankItem(id)
	if nil == item {
		return -1
	} else {
		return r.getRank(item)
	}
}

func (r *Rank) Check() bool {
	if len(r.spans) == 0 {
		return true
	}
	vv := r.spans[0].max
	for i, v := range r.spans {
		vv = v.check(vv)
		if vv == -1 && i > 0 {
			r.spans[i-1].show()
			return false
		}
	}

	return true
}

func (r *Rank) Show() {
	for _, v := range r.spans {
		v.show()
	}
}

func (r *Rank) getRankItem(id uint64) *node {
	return r.id2Item[id]
}

func (r *Rank) binarySearch(score int, left int, right int) *skiplists {

	if left >= right {
		return r.spans[left]
	}

	mIdx := (right-left)/2 + left
	m := r.spans[mIdx]

	if m.max > score {
		nIdx := mIdx + 1
		if nIdx >= len(r.spans) || r.spans[nIdx].max < score {
			return m
		}
		return r.binarySearch(score, mIdx+1, right)
	} else {
		pIdx := mIdx - 1
		if pIdx < 0 || r.spans[pIdx].min > score {
			return m
		}
		return r.binarySearch(score, left, mIdx-1)
	}

}

func (r *Rank) findSpan(score int) *skiplists {
	var c *skiplists
	if len(r.spans) == 0 {
		c = newSkipLists(len(r.spans))
		r.spans = append(r.spans, c)
	} else {
		c = r.binarySearch(score, 0, len(r.spans)-1)
	}

	return c
}

//返回排名和排名百分比
func (r *Rank) UpdateScore(id uint64, score int) (int, int) {

	r.cc++

	defer func() {
		if r.cc%100 == 0 {
			r.shrink(nil)
		}
	}()

	rank := 0

	item := r.getRankItem(id)
	if nil == item {
		item = r.itemPool.get()
		r.id2Item[id] = item
	} else {
		if item.value == score {
			return r.getRank(item), r.getRankPercentByItem(item)
		}
	}

	item.key = id
	item.value = score

	c := r.findSpan(score)

	oldC := item.sl

	if item.sl != nil {
		sl := item.sl
		sl.DeleteNode(item)
		sl.fixMinMax()
	}

	rank = c.InsertNode(item)

	if c.size > maxItemCount+maxItemCount/2 {

		if o := c.split(); nil != o {

			c.fixMinMax()
			o.fixMinMax()

			o.idx = c.idx + 1
			if o.idx >= len(r.spans) {
				r.spans = append(r.spans, o)
			} else {
				if len(r.spans) < cap(r.spans) {
					//还有空间,扩张containers,将downIdx开始的元素往后挪一个位置，空出downIdx所在位置
					l := len(r.spans)
					r.spans = r.spans[:len(r.spans)+1]
					for i := l - 1; i >= o.idx; i-- {
						r.spans[i+1] = r.spans[i]
						r.spans[i+1].idx = i + 1
					}
					r.spans[o.idx] = o

				} else {

					//下一个container满了，新建一个
					spans := make([]*skiplists, 0, len(r.spans)+1)
					for i := 0; i <= o.idx-1; i++ {
						spans = append(spans, r.spans[i])
					}

					spans = append(spans, o)

					for i := o.idx; i < len(r.spans); i++ {
						c := r.spans[i]
						c.idx = len(spans)
						spans = append(spans, c)
					}

					r.spans = spans
				}
			}
		} else {
			c.fixMinMax()
		}
	} else {
		c.fixMinMax()
	}

	if nil != oldC {
		if oldC.size == 0 {
			//oldC已经被清空，需要删除
			for i := oldC.idx + 1; i < len(r.spans); i++ {
				c := r.spans[i]
				r.spans[i-1] = c
				c.idx = i - 1
			}

			r.spans[len(r.spans)-1] = nil
			r.spans = r.spans[:len(r.spans)-1]
		} else if oldC.idx != item.sl.idx && oldC.idx+1 < len(r.spans) && oldC.size+r.spans[oldC.idx+1].size <= maxItemCount {
			r.shrink(oldC)
		}
	}

	return rank + r.getFrontSpanItemCount(item), r.getRankPercentByItem(item) // r.getRankPercent(item, rank)
}

func (r *Rank) shrink(s *skiplists) {
	if nil == s {
		if r.nextShink >= len(r.spans)-1 {
			r.nextShink = 0
			return
		} else {
			s = r.spans[r.nextShink]
			r.nextShink++
		}
	}

	if s.idx+1 < len(r.spans) && s.size+r.spans[s.idx+1].size <= maxItemCount /*+maxItemCount/5*/ {
		n := r.spans[s.idx+1]
		s.merge(n)
		s.fixMinMax()
		for i := n.idx + 1; i < len(r.spans); i++ {
			c := r.spans[i]
			r.spans[i-1] = c
			c.idx = i - 1
		}
		r.spans[len(r.spans)-1] = nil
		r.spans = r.spans[:len(r.spans)-1]
	}
}

func (r *Rank) GetScore(id uint64) int {
	item := r.getRankItem(id)
	if nil == item {
		return -1
	} else {
		return item.value
	}
}

// 获取排名前 n
func (r *Rank) GetTopN(n int) []uint64 {
	ids := make([]uint64, 0, n)
	for _, sl := range r.spans {
		node := &sl.head
		for &sl.tail != node.links[0].pnext && n > 0 {
			node = node.links[0].pnext
			n--
			ids = append(ids, node.key)
		}
		if n <= 0 {
			break
		}
	}

	return ids
}

// 获取整个排名数据
func (r *Rank) GetRankList() []uint64 {
	return r.GetTopN(r.Size())
}

// 获取指定名次的分数 1...
func (r *Rank) GetScoreByIdx(idx int) (uint64, int) {
	curIdx := 0
	var sl *skiplists

	for _, s := range r.spans {
		if curIdx+s.size < idx {
			curIdx += s.size
		} else {
			sl = s
			break
		}
	}

	if sl == nil {
		return 0, -1
	}

	node := &sl.head
	for i := sl.level; i >= 0; i-- {
		for &sl.tail != node.links[i].pnext && node.links[i].skip+curIdx <= idx {
			curIdx += node.links[i].skip
			node = node.links[i].pnext
		}
	}

	//node := &sl.head
	//for &sl.tail != node.links[0].pnext && curIdx < idx {
	//	node = node.links[0].pnext
	//	curIdx++
	//}

	if curIdx != idx {
		return 0, -1
	}

	return node.key, node.value
}

func (r *Rank) Size() int {
	return len(r.id2Item)
}

func (r *Rank) Delete(id uint64) {
	item := r.getRankItem(id)
	if nil == item {
		return
	}

	sl := item.sl
	if sl == nil {
		delete(r.id2Item, id)
	} else {
		sl.DeleteNode(item)
	}
}
