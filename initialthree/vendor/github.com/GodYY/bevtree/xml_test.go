package bevtree

import (
	"math/rand"
	"testing"
	"time"
)

func TestTreeMarshalXML(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	framework := newTestFramework()

	oldTree := NewTree("XML测试")
	framework.addTree(oldTree)

	parallel := NewParallelNode()
	oldTree.Root().SetChild(parallel)

	key := "key"
	sum := 0
	unit := 1
	low := 5
	max := 10

	rand.Seed(time.Now().UnixNano())

	{

		subtree_a := NewTree("subtree_a")
		framework.addTree(subtree_a)
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
		framework.addTree(subtree_b)
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

	// paral.AddChild(NewRandSequence())
	// paral.AddChild(NewRandSelector())
	// paral.AddChild(NewParallel())

	entity, err := framework.CreateEntity("XML测试", nil)
	if err != nil {
		t.Fatal(err)
	}

	entity.Context().DataSet().Set(key, 0)
	entity.Update()
	v, _ := entity.Context().DataSet().GetInt(key)
	if v != sum {
		t.Fatalf("test Tree before marshal: sum(%d) != %d", v, sum)
	}
	entity.Release()

	data, err := framework.MarshalXMLTree(oldTree)
	if err != nil {
		t.Fatal("marshal Tree:", err)
	} else {
		t.Log("marshal Tree:", string(data))
	}

	newTree := new(tree)
	if err := framework.UnmarshalXMLTree(data, newTree); err != nil {
		t.Fatal("unmarshal previos Tree:", err)
	}
	newTree.SetName("XML测试2")
	framework.addTree(newTree)

	entity, err = framework.CreateEntity("XML测试2", nil)
	if err != nil {
		t.Fatal(err)
	}

	entity.Context().DataSet().Set(key, 0)
	entity.Update()

	v, _ = entity.Context().DataSet().GetInt(key)
	if v != sum {
		t.Fatalf("test Tree after unmarshal: sum(%d) != %d", v, sum)
	}

}
