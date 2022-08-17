package bevtree

import "github.com/GodYY/gutils/assert"

// Behavior type.
type BevType string

func (t BevType) Valid() bool { return string(t) != "" }

func (t BevType) String() string { return string(t) }

// Bev is the structure data of Behavior.
type Bev interface {
	BevType() BevType
	CreateInstance() BevInstance
	DestroyInstance(BevInstance)
}

// BevInstance is the entity of Bev for running.
type BevInstance interface {
	// Behavior type.
	BevType() BevType

	// OnInit is called before the first update of the behavior.
	OnInit(Context) bool

	// OnUpdate is called when the behavior tree update before the
	// behavior terminate.
	OnUpdate(Context) Result

	// OnTerminate is called after the last update of the behavior.
	OnTerminate(Context)
}

// The behavior node of behavior tree, a kind of leaf node.
type BevNode struct {
	// Common part of node.
	node

	// Behavior.
	bev Bev
}

func NewBevNode(bev Bev) *BevNode {
	assert.Assert(bev != nil, "bev nil")
	return &BevNode{
		node: newNode(),
		bev:  bev,
	}
}

func (BevNode) NodeType() NodeType { return behavior }

func (b *BevNode) Bev() Bev { return b.bev }

func (b *BevNode) SetBev(bev Bev) {
	assert.Assert(bev != nil, "bev nil")
	b.bev = bev
}

// Behavior task, the runtime of BevNode.
type bevTask struct {
	bev     Bev
	bevInst BevInstance
}

func (b *bevTask) TaskType() TaskType { return Single }

func (b *bevTask) OnCreate(node Node) {
	b.bev = node.(*BevNode).Bev()
	b.bevInst = b.bev.CreateInstance()
}

func (b *bevTask) OnInit(_ NodeList, ctx Context) bool {
	return b.bevInst.OnInit(ctx)
}

func (b *bevTask) OnUpdate(ctx Context) Result {
	return b.bevInst.OnUpdate(ctx)
}

func (b *bevTask) OnTerminate(ctx Context) {
	b.bevInst.OnTerminate(ctx)
	b.bev.DestroyInstance(b.bevInst)
	b.bevInst = nil
	b.bev = nil
}

func (b *bevTask) OnChildTerminated(Result, NodeList, Context) Result { panic("shouldnt be invoked") }
