// +build debug

package bevtree

import (
	"math/rand"
	"runtime"
	"testing"
	"time"
)

func TestPool(t *testing.T) {
	test := newTest()

	tree := test.createTree("test pool")

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

	test.run(t, "test pool", Failure, 50)

	getTotal := _poolDebug.getTotal()
	putTotal := _poolDebug.putTotal()
	loseTotal := getTotal - putTotal
	if loseTotal > 0 {
		t.Fatalf("get: %d put: %d lose:%d", getTotal, putTotal, loseTotal)
	} else {
		t.Logf("get==put: %d", getTotal)
	}
}

func TestDebugReset(t *testing.T) {
	framework := newTestFramework()

	tree := NewTree("test debug reset")
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

	e, err := framework.CreateEntity("test debug reset", nil)
	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < 100; i++ {
		e.Update()
		e.Stop()

		getTotal := _poolDebug.getTotal()
		putTotal := _poolDebug.putTotal()
		loseTotal := getTotal - putTotal
		if loseTotal > 0 {
			t.Fatalf("get: %d put: %d lose:%d", getTotal, putTotal, loseTotal)
		} else {
			t.Logf("get==put: %d", getTotal)
		}
	}
}

func TestEntityFinalizer(t *testing.T) {
	framework := newTestFramework()

	tree := NewTree("test entity finalizer")
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

	for i := 0; i < 100; i++ {
		entity, err := framework.CreateEntity("test entity finalizer", nil)
		if err != nil {
			t.Fatal(err)
		}

		entity.Update()
		entity.Stop()
	}

	runtime.GC()
	time.Sleep(1 * time.Second)

	getTotal := _poolDebug.getTotal()
	putTotal := _poolDebug.putTotal()
	loseTotal := getTotal - putTotal
	if loseTotal > 0 {
		t.Fatalf("get: %d put: %d lose:%d", getTotal, putTotal, loseTotal)
	} else {
		t.Logf("get==put: %d", getTotal)
	}
}
