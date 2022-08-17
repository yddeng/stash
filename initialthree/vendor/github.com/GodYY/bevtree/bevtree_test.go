package bevtree

import (
	"math"
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/GodYY/gutils/assert"
	"github.com/pkg/errors"
)

func (s *Framework) addTree(tree *tree) {
	if !s.initialized {
		panic(errors.New("bevtree framework uninitialized"))
	}

	if s.treeAssets == nil {
		s.treeAssets = make(map[string]*treeAsset)
	}

	_, ok := s.treeAssets[tree.Name()]
	assert.AssertF(!ok, "tree named \"%s\" already exist", tree.Name())

	var ta *treeAsset
	if s.loadAll {
		ta = &treeAsset{
			entry: &TreeEntry{Name: tree.Name()},
			tree:  tree,
		}
	} else {
		ta = &treeAsset{
			entry: &TreeEntry{Name: tree.Name()},
			once:  new(sync.Once),
			tree:  tree,
		}
	}

	s.treeAssets[tree.Name()] = ta
}

const (
	function       = BevType("func")
	increase       = BevType("incr")
	update         = BevType("update")
	blackboardIncr = BevType("blackboardIncr")
)

type bevFunc struct {
	f func(Context) Result
}

func newBevFunc(f func(Context) Result) *bevFunc {
	return &bevFunc{f: f}
}

func (bevFunc) BevType() BevType { return function }

func (b *bevFunc) CreateInstance() BevInstance {
	return &bevFuncEntity{f: b.f}
}

func (b *bevFunc) DestroyInstance(BevInstance) {}

type bevFuncEntity struct {
	f func(Context) Result
}

func (b *bevFuncEntity) BevType() BevType { return function }

func (b *bevFuncEntity) OnInit(e Context) bool { return true }

func (b *bevFuncEntity) OnUpdate(e Context) Result {
	// fmt.Println("behaviorFunc OnUpdate")
	return b.f(e)
}

func (b *bevFuncEntity) OnTerminate(e Context) {
}

type behaviorIncr struct {
	key     string
	limited int
}

func newBehaviorIncr(key string, limited int) *behaviorIncr {
	return &behaviorIncr{key: key, limited: limited}
}

func (behaviorIncr) BevType() BevType { return increase }

func (b *behaviorIncr) CreateInstance() BevInstance {
	return &behaviorIncrEntity{behaviorIncr: b}
}

func (b *behaviorIncr) DestroyInstance(BevInstance) {}

type behaviorIncrEntity struct {
	*behaviorIncr
	count int
}

func (b *behaviorIncrEntity) BevType() BevType { return increase }

func (b *behaviorIncrEntity) OnInit(e Context) bool { return true }

func (b *behaviorIncrEntity) OnUpdate(e Context) Result {
	if b.count >= b.limited {
		return Failure
	}

	b.count++
	e.DataSet().IncInt(b.key)
	if b.count >= b.limited {
		return Success
	}

	return Running
}

func (b *behaviorIncrEntity) OnTerminate(e Context) { b.count = 0 }

type behaviorUpdate struct {
	limited int
}

func newBehaviorUpdate(lmited int) *behaviorUpdate {
	return &behaviorUpdate{limited: lmited}
}

func (behaviorUpdate) BevType() BevType { return update }

func (b *behaviorUpdate) CreateInstance() BevInstance {
	return &behaviorUpdateEntity{behaviorUpdate: b}
}

func (b *behaviorUpdate) DestroyInstance(BevInstance) {}

type behaviorUpdateEntity struct {
	*behaviorUpdate
	count int
}

func (b *behaviorUpdateEntity) BevType() BevType { return update }

func (b *behaviorUpdateEntity) OnInit(e Context) bool { return true }

func (b *behaviorUpdateEntity) OnUpdate(e Context) Result {
	if b.count >= b.limited {
		return Success
	}

	b.count++
	if b.count >= b.limited {
		return Success
	}

	return Running
}

func (b *behaviorUpdateEntity) OnTerminate(e Context) { b.count = 0 }

func init() {

}

type bevBBIncr struct {
	Key     string
	Limited int
}

func newBevBBIncr(key string, limited int) *bevBBIncr {
	return &bevBBIncr{
		Key:     key,
		Limited: limited,
	}
}

func (bevBBIncr) BevType() BevType { return blackboardIncr }

func (b *bevBBIncr) CreateInstance() BevInstance {
	return &bevBBIncrEntity{bevBBIncr: b}
}

func (b *bevBBIncr) DestroyInstance(BevInstance) {}

type bevBBIncrEntity struct {
	*bevBBIncr
	count int
}

func (b *bevBBIncrEntity) BevType() BevType      { return blackboardIncr }
func (b *bevBBIncrEntity) OnInit(_ Context) bool { return true }

func (b *bevBBIncrEntity) OnUpdate(e Context) Result {
	e.DataSet().IncInt(b.Key)
	b.count++

	if b.count >= b.Limited {
		return Success
	} else {
		return Running
	}

}

func (b *bevBBIncrEntity) OnTerminate(_ Context) {}

func newTestFramework() *Framework {
	framework := NewFramework()
	framework.RegisterBevType(function, func() Bev { return new(bevFunc) })
	framework.RegisterBevType(increase, func() Bev { return new(behaviorIncr) })
	framework.RegisterBevType(update, func() Bev { return new(behaviorUpdate) })
	framework.RegisterBevType(blackboardIncr, func() Bev { return &bevBBIncr{} })
	framework.initialized = true
	framework.loadAll = true
	return framework
}

type test struct {
	framework *Framework
}

func newTest() *test {
	t := new(test)
	t.framework = newTestFramework()
	return t
}

func (t *test) createTree(name string) *tree {
	tree := NewTree(name)
	t.framework.addTree(tree)
	return tree
}

type keyValue struct {
	key      string
	def      interface{}
	expected interface{}
}

func (t *test) run(tt *testing.T, tree string, expectedResult Result, tick int, expectedKeyValues ...keyValue) {
	entity, err := t.framework.CreateEntity(tree, nil)
	if err != nil {
		tt.Fatalf("create entity: %s", err)
	}

	for _, v := range expectedKeyValues {
		entity.Context().DataSet().Set(v.key, v.def)
	}

	result := Running

	for i := 0; i < tick; i++ {
		tt.Log("run", i, "start")
		result = Running
		k := 0
		for result == Running {
			tt.Log("run", i, "update", k)
			k++
			result = entity.Update()
			time.Sleep(1 * time.Millisecond)
		}
		tt.Log("run", i, "end", result)
	}

	if result != expectedResult {
		tt.Fatalf("should return %v but get %v", expectedResult, result)
	}

	for _, v := range expectedKeyValues {
		if entity.Context().DataSet().Get(v.key) != v.expected {
			tt.Fatalf("%s = %v(%v)", v.key, entity.Context().DataSet().Get(v.key), v.expected)
		}
	}

	entity.Release()
}

func TestRoot(t *testing.T) {
	test := newTest()
	test.createTree("test root")
	test.run(t, "test root", Failure, 1)
}

func TestSequence(t *testing.T) {
	test := newTest()

	tree := test.createTree("test sequence")

	seq := NewSequenceNode()

	tree.Root().SetChild(seq)

	key := "counter"
	n := 2
	for i := 0; i < n; i++ {
		seq.AddChild(NewBevNode(newBevFunc(func(e Context) Result {
			e.DataSet().IncInt(key)
			return Success
		})))
	}

	test.run(t, "test sequence", Success, 1, keyValue{key: "counter", def: 0, expected: 2})
}

func TestSelector(t *testing.T) {
	test := newTest()

	tree := test.createTree("test selector")

	selc := NewSelectorNode()

	tree.Root().SetChild(selc)

	key := "selected"
	var selected int
	n := 10
	for i := 0; i < n; i++ {
		k := i
		selc.AddChild(NewBevNode(newBevFunc((func(e Context) Result {
			if k == selected {
				e.DataSet().SetInt(key, selected)
				return Success
			} else {
				return Failure
			}
		}))))
	}

	rand.Seed(time.Now().Unix())
	selected = rand.Intn(n)

	test.run(t, "test selector", Success, 1, keyValue{key: key, def: -1, expected: selected})
}

func TestRandomSequence(t *testing.T) {
	rand.Seed(time.Now().Unix())

	test := newTest()

	tree := test.createTree("test random sequence")

	seq := NewRandSequenceNode()

	tree.Root().SetChild(seq)

	key := "counter"
	n := 2
	for i := 0; i < n; i++ {
		k := i
		seq.AddChild(NewBevNode(newBevFunc(func(e Context) Result {
			t.Log("seq", k, "update")
			e.DataSet().IncInt(key)
			return Success
		})))
	}

	test.run(t, "test random sequence", Success, 1, keyValue{key: key, def: 0, expected: n})
}

func TestRandomSelector(t *testing.T) {
	rand.Seed(time.Now().Unix())

	test := newTest()

	tree := test.createTree("test random selector")

	selc := NewRandSelectorNode()

	tree.Root().SetChild(selc)

	key := "selected"
	var selected int
	n := 10
	for i := 0; i < n; i++ {
		k := i
		selc.AddChild(NewBevNode(newBevFunc((func(e Context) Result {
			t.Log("seq", k, "update")
			if k == selected {
				e.DataSet().SetInt(key, selected)
				return Success
			} else {
				return Failure
			}
		}))))
	}

	rand.Seed(time.Now().Unix())
	selected = rand.Intn(n)

	test.run(t, "test random selector", Success, 1, keyValue{key: key, def: -1, expected: selected})
}

func TestParallel(t *testing.T) {
	test := newTest()

	tree := test.createTree("test parallel")

	paral := NewParallelNode()

	tree.Root().SetChild(paral)

	rand.Seed(time.Now().Unix())

	n := 2
	for i := 0; i < n; i++ {
		k := i + 1
		timer := time.NewTimer(1000 * time.Millisecond * time.Duration(k))
		paral.AddChild(NewBevNode(newBevFunc(func(e Context) Result {
			select {
			case <-timer.C:
				t.Logf("timer No.%d up", k)
				return Success
			default:
				t.Logf("timer No.%d update", k)
				return Running
			}
		})))
	}

	test.run(t, "test parallel", Success, 1)
}

func TestParallelLazyStop(t *testing.T) {
	test := newTest()

	tree := test.createTree("test parallel lazy stop")

	paral := NewParallelNode()

	tree.Root().SetChild(paral)

	rand.Seed(time.Now().Unix())

	lowUpdate, maxUpdate := 2, 10
	n := 10
	lowDepth, maxDepth := 5, 10
	for i := 0; i < n; i++ {
		k := i + 1
		ut := lowUpdate + rand.Intn(maxUpdate-lowUpdate+1)
		c := DecoratorNode(NewInverterNode())
		paral.AddChild(c)

		depth := 5 + rand.Intn(maxDepth-lowDepth)
		for d := 0; d < depth; d++ {
			cc := NewSucceederNode()
			c.SetChild(cc)
			c = cc
		}

		c.SetChild(NewBevNode(newBevFunc(func(e Context) Result {
			t.Logf("No.%d update", k)
			ut--
			if ut <= 0 {
				t.Logf("No.%d over", k)
				return Success
			} else {
				return Running
			}
		})))
	}

	test.run(t, "test parallel lazy stop", Failure, 1)
}

func TestRepeater(t *testing.T) {
	test := newTest()

	tree := test.createTree("test repeater")

	n := 10
	repeater := NewRepeaterNode(n)

	tree.Root().SetChild(repeater)

	key := "counter"

	repeater.SetChild(NewBevNode(newBevFunc((func(e Context) Result {
		e.DataSet().IncInt(key)
		return Success
	}))))

	test.run(t, "test repeater", Success, 1, keyValue{key: key, def: 0, expected: n})
}

func TestInverter(t *testing.T) {
	test := newTest()

	tree := test.createTree("test inverter")

	inverter := NewInverterNode()

	tree.Root().SetChild(inverter)

	inverter.SetChild(NewBevNode(newBevFunc(func(e Context) Result {
		return Failure
	})))

	test.run(t, "test inverter", Success, 1)
}

func TestSucceeder(t *testing.T) {
	test := newTest()

	tree := test.createTree("test succeeder")

	succeeder := NewSucceederNode()
	tree.Root().SetChild(succeeder)

	succeeder.SetChild(NewBevNode(newBevFunc(func(e Context) Result { return Failure })))

	test.run(t, "test succeeder", Success, 1)
}

func TestRepeatUntilFail(t *testing.T) {
	test := newTest()

	tree := test.createTree("test repeat until fail")

	repeat := NewRepeatUntilFailNode(false)
	tree.Root().SetChild(repeat)

	n := 4
	repeat.SetChild(NewBevNode(newBevFunc(func(e Context) Result {
		t.Log("decr 1")

		n--

		if n <= 0 {
			return Failure
		}

		return Success
	})))

	test.run(t, "test repeat until fail", Failure, 1)
}

func TestShareTree(t *testing.T) {
	framework := newTestFramework()

	tree := NewTree("test share tree")
	framework.addTree(tree)

	paral := NewParallelNode()
	tree.Root().SetChild(paral)

	expectedResult := Success
	singleSum := 0
	key := "sum"
	numEntities := 100
	low, max := 5, 50
	n := 100

	rand.Seed(time.Now().UnixNano())

	for i := 0; i < n; i++ {
		limited := low + rand.Intn(max-low+1)
		singleSum += limited
		t.Logf("singleSum add %d to %d", limited, singleSum)

		paral.AddChild(NewBevNode(newBehaviorIncr(key, limited)))
	}

	entities := make([]Entity, numEntities)
	for i := 0; i < numEntities; i++ {
		if entity, err := framework.CreateEntity("test share tree", nil); err == nil {
			entities[i] = entity
			entities[i].Context().DataSet().SetInt(key, 0)
		} else {
			t.Fatal(err)
		}
	}

	result := Running
	for result == Running {
		for i := 0; i < numEntities; i++ {
			if i > 0 {
				r := entities[i].Update()
				if r != result {
					t.Fatal("invalid result", result, r)
				}
			} else {
				result = entities[i].Update()
			}
		}

		time.Sleep(1 * time.Millisecond)
	}

	if result != expectedResult {
		t.Fatalf("expected %v get %v", expectedResult, result)
	}

	sum := 0
	for i := 0; i < numEntities; i++ {
		v, _ := entities[i].Context().DataSet().GetInt(key)
		sum += v
	}

	if sum != singleSum*numEntities {
		t.Fatalf("expected sum %d get %d", singleSum*numEntities, sum)
	}
}

func TestReset(t *testing.T) {
	framework := newTestFramework()

	tree := NewTree("test reset")
	framework.addTree(tree)

	paral := NewParallelNode()

	tree.Root().SetChild(paral)

	rand.Seed(time.Now().Unix())

	lowUpdate, maxUpdate := 2, 10
	n := 10
	lowDepth, maxDepth := 5, 10
	for i := 0; i < n; i++ {
		ut := lowUpdate + rand.Intn(maxUpdate-lowUpdate+1)
		c := DecoratorNode(NewInverterNode())
		paral.AddChild(c)

		depth := 5 + rand.Intn(maxDepth-lowDepth)
		for d := 0; d < depth; d++ {
			cc := NewSucceederNode()
			c.SetChild(cc)
			c = cc
		}

		c.SetChild(NewBevNode(newBehaviorUpdate(ut)))
	}

	e, err := framework.CreateEntity("test reset", nil)
	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < 100; i++ {
		e.Update()
		e.Stop()
	}
}

func TestRemoveChild(t *testing.T) {
	key := "key"
	sum := 0

	unit := 1

	rand.Seed(time.Now().UnixNano())

	framework := newTestFramework()

	tree := NewTree("test remove child")
	framework.addTree(tree)

	paral := NewParallelNode()
	tree.Root().SetChild(paral)

	bd := newBevBBIncr(key, unit)

	sc := NewSucceederNode()
	sc.SetChild(NewBevNode(bd))
	paral.AddChild(sc)
	sum += unit

	low := 5
	max := 10
	rtimes := low + rand.Intn(max-low+1)
	r := NewRepeaterNode(rtimes)
	r.SetChild(NewBevNode(bd))
	paral.AddChild(r)
	sum += rtimes * unit

	iv_sc := NewSucceederNode()
	iv := NewInverterNode()
	iv.SetChild(NewBevNode(bd))
	iv_sc.SetChild(iv)
	paral.AddChild(iv_sc)
	sum += unit

	ruf := NewRepeatUntilFailNode(true)
	ruf_iv := NewInverterNode()
	ruf.SetChild(ruf_iv)
	ruf_iv.SetChild(NewBevNode(bd))
	paral.AddChild(ruf)
	sum += unit

	seqTimes := low + rand.Intn(max-low+1)
	seq := NewSequenceNode()
	for i := 0; i < seqTimes; i++ {
		seq.AddChild(NewBevNode(bd))
	}
	paral.AddChild(seq)
	sum += seqTimes * unit

	selcTimes := low + rand.Intn(max-low+1)
	selc := NewSelectorNode()
	selcSuccN := rand.Intn(selcTimes)
	for i := 0; i < selcTimes; i++ {
		if selcSuccN == i {
			selc.AddChild(NewBevNode(bd))
		} else {
			iv := NewInverterNode()
			iv.SetChild(NewBevNode(bd))
			selc.AddChild(iv)
		}
	}
	paral.AddChild(selc)
	sum += (selcSuccN + 1) * unit

	for paral.ChildCount() > 0 {
		paral.RemoveChild(0)
	}

	if paral.ChildCount() > 0 {
		t.FailNow()
	}

	entity, err := framework.CreateEntity("test remove child", nil)
	if err != nil {
		t.Fatal(err)
	}

	if entity.Update() != Failure {
		t.FailNow()
	}
}

func TestSubtree(t *testing.T) {
	test := newTest()

	tree := test.createTree("test subtree")

	parallel := NewParallelNode()
	tree.Root().SetChild(parallel)

	key := "key"
	sum := 0
	unit := 1
	low := 5
	max := 10

	rand.Seed(time.Now().UnixNano())

	{

		subtree_a := NewTree("subtree_a")
		test.framework.addTree(subtree_a)
		parallel.AddChild(NewSubtreeNode(subtree_a, false))
		paral := NewParallelNode()
		subtree_a.Root().SetChild(paral)

		bd := newBevBBIncr(key, unit)

		sc := NewSucceederNode()
		sc.SetChild(NewBevNode(bd))
		paral.AddChild(sc)
		sum += unit

		rtimes := low + rand.Intn(max-low+1)
		r := NewRepeaterNode(rtimes)
		r.SetChild(NewBevNode(bd))
		paral.AddChild(r)
		sum += rtimes * unit

		iv_sc := NewSucceederNode()
		iv := NewInverterNode()
		iv.SetChild(NewBevNode(bd))
		iv_sc.SetChild(iv)
		paral.AddChild(iv_sc)
		sum += unit

		ruf := NewRepeatUntilFailNode(true)
		ruf_iv := NewInverterNode()
		ruf.SetChild(ruf_iv)
		ruf_iv.SetChild(NewBevNode(bd))
		paral.AddChild(ruf)
		sum += unit

		seqTimes := low + rand.Intn(max-low+1)
		seq := NewSequenceNode()
		for i := 0; i < seqTimes; i++ {
			seq.AddChild(NewBevNode(bd))
		}
		paral.AddChild(seq)
		sum += seqTimes * unit

		selcTimes := low + rand.Intn(max-low+1)
		selc := NewSelectorNode()
		selcSuccN := rand.Intn(selcTimes)
		for i := 0; i < selcTimes; i++ {
			if selcSuccN == i {
				selc.AddChild(NewBevNode(bd))
			} else {
				iv := NewInverterNode()
				iv.SetChild(NewBevNode(bd))
				selc.AddChild(iv)
			}
		}
		paral.AddChild(selc)
		sum += (selcSuccN + 1) * unit
	}

	{
		subtree_b := NewTree("subtree_b")
		test.framework.addTree(subtree_b)
		parallel.AddChild(NewSubtreeNode(subtree_b, false))
		paral := NewParallelNode()
		subtree_b.Root().SetChild(paral)

		bd := newBevBBIncr(key, unit)

		sc := NewSucceederNode()
		sc.SetChild(NewBevNode(bd))
		paral.AddChild(sc)
		sum += unit

		rtimes := low + rand.Intn(max-low+1)
		r := NewRepeaterNode(rtimes)
		r.SetChild(NewBevNode(bd))
		paral.AddChild(r)
		sum += rtimes * unit

		iv_sc := NewSucceederNode()
		iv := NewInverterNode()
		iv.SetChild(NewBevNode(bd))
		iv_sc.SetChild(iv)
		paral.AddChild(iv_sc)
		sum += unit

		ruf := NewRepeatUntilFailNode(true)
		ruf_iv := NewInverterNode()
		ruf.SetChild(ruf_iv)
		ruf_iv.SetChild(NewBevNode(bd))
		paral.AddChild(ruf)
		sum += unit

		seqTimes := low + rand.Intn(max-low+1)
		seq := NewSequenceNode()
		for i := 0; i < seqTimes; i++ {
			seq.AddChild(NewBevNode(bd))
		}
		paral.AddChild(seq)
		sum += seqTimes * unit

		selcTimes := low + rand.Intn(max-low+1)
		selc := NewSelectorNode()
		selcSuccN := rand.Intn(selcTimes)
		for i := 0; i < selcTimes; i++ {
			if selcSuccN == i {
				selc.AddChild(NewBevNode(bd))
			} else {
				iv := NewInverterNode()
				iv.SetChild(NewBevNode(bd))
				selc.AddChild(iv)
			}
		}
		paral.AddChild(selc)
		sum += (selcSuccN + 1) * unit
	}

	test.run(t, "test subtree", Success, 1, keyValue{key: key, def: 0, expected: sum})
}

func TestWeightSelector(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	framework := newTestFramework()

	tree := NewTree("test weight selector")
	framework.addTree(tree)

	ws := NewWeightSelectorNode()
	tree.Root().SetChild(ws)

	n := 10000
	key := "key"
	valWeights := map[int]float32{
		1: 0.3,
		2: 0.4,
		3: 0.15,
		4: 0.05,
		5: 0.1,
	}
	tolerance := 0.015

	for v, w := range valWeights {
		vv := v
		ws.AddChild(NewBevNode(newBevFunc(func(c Context) Result {
			c.DataSet().Set(key, vv)
			return Success
		})), w)
	}

	entity, _ := framework.CreateEntity("test weight selector", nil)
	results := map[int]int{}
	for i := 0; i < n; i++ {
		entity.Update()
		results[entity.Context().DataSet().Get(key).(int)] += 1
	}

	t.Log("probability statistics")
	for v, count := range results {
		p := float64(count) / float64(n)
		diff := math.Abs(p - float64(valWeights[v]))
		if diff < tolerance {
			t.Logf("\t%d: %f", v, p)
		} else {
			t.Fatalf("\t%d: %f, taget: %f, diff(%f) > tolerance(%f)", v, p, valWeights[v], diff, tolerance)
		}
	}
}
