package bevtree

import (
	"math/rand"

	"github.com/GodYY/gutils/assert"
)

// The CompositeNode Interface represents the common functions that
// the composite nodes in behavior tree need implement.
type CompositeNode interface {
	Node

	// Get the number of child nodes.
	ChildCount() int

	// Get the child node with index idx.
	Child(idx int) Node

	// Add a child node.
	AddChild(child Node)

	// Remove a child node with index idx.
	RemoveChild(idx int) Node
}

// The common part of composite node.
type compositeNode struct {
	node

	// The child nodes.
	children []Node
}

func newCompositeNode() compositeNode {
	return compositeNode{}
}

func (c *compositeNode) ChildCount() int { return len(c.children) }

func (c *compositeNode) Child(idx int) Node {
	assert.Assert(idx >= 0 && idx < c.ChildCount(), "index out of range")

	return c.children[idx]
}

func (c *compositeNode) addChild(child Node) {
	assert.Assert(child != nil, "child nil")
	assert.Assert(child.Parent() == nil, "child already has parent")

	c.children = append(c.children, child)
}

func (c *compositeNode) RemoveChild(idx int) Node {
	assert.Assert(idx >= 0 && idx < c.ChildCount(), "index out of range")

	child := c.children[idx]
	child.SetParent(nil)
	c.children = append(c.children[:idx], c.children[idx+1:]...)
	return child
}

// Sequence node runs child node one bye one until a child
// returns failure. It returns the result of the last
// running node.
type SequenceNode struct {
	compositeNode
}

func NewSequenceNode() *SequenceNode {
	return &SequenceNode{
		compositeNode: newCompositeNode(),
	}
}

func (s *SequenceNode) NodeType() NodeType { return sequence }

func (s *SequenceNode) AddChild(child Node) {
	s.compositeNode.addChild(child)
	child.SetParent(s)
}

// The sequence node task.
type sequenceTask struct {
	node        *SequenceNode
	curChildIdx int
}

func (s *sequenceTask) TaskType() TaskType { return Serial }

func (s *sequenceTask) OnCreate(node Node) {
	s.node = node.(*SequenceNode)
	s.curChildIdx = 0
}

func (s *sequenceTask) OnInit(nextChildNodes NodeList, ctx Context) bool {
	if s.node.ChildCount() == 0 {
		return false
	}

	nextChildNodes.PushNode(s.node.Child(0))
	return true
}

func (s *sequenceTask) OnUpdate(ctx Context) Result { return Running }
func (s *sequenceTask) OnTerminate(ctx Context)     { s.node = nil }

func (s *sequenceTask) OnChildTerminated(result Result, nextChildNodes NodeList, ctx Context) Result {
	s.curChildIdx++
	if result == Success && s.curChildIdx < s.node.ChildCount() {
		nextChildNodes.PushNode(s.node.Child(s.curChildIdx))
		return Running
	} else {
		return result
	}
}

// Selector node runs child node one by one until a child
// returns success. It returns the result of the last
// running node.
type SelectorNode struct {
	compositeNode
}

func NewSelectorNode() *SelectorNode {
	return &SelectorNode{
		compositeNode: newCompositeNode(),
	}
}

func (s *SelectorNode) NodeType() NodeType { return selector }

func (s *SelectorNode) AddChild(child Node) {
	s.compositeNode.addChild(child)
	child.SetParent(s)
}

// The selector node task.
type selectorTask struct {
	node        *SelectorNode
	curChildIdx int
}

func (s *selectorTask) TaskType() TaskType { return Serial }

func (s *selectorTask) OnCreate(node Node) {
	s.node = node.(*SelectorNode)
	s.curChildIdx = 0
}

func (s *selectorTask) OnInit(nextChildNodes NodeList, ctx Context) bool {
	if s.node.ChildCount() == 0 {
		return false
	} else {
		nextChildNodes.PushNode(s.node.Child(0))
		return true
	}
}

func (s *selectorTask) OnUpdate(ctx Context) Result { return Running }
func (s *selectorTask) OnTerminate(ctx Context)     { s.node = nil }

func (s *selectorTask) OnChildTerminated(result Result, nextChildNodes NodeList, ctx Context) Result {
	s.curChildIdx++
	if result == Failure && s.curChildIdx < s.node.ChildCount() {
		nextChildNodes.PushNode(s.node.Child(s.curChildIdx))
		return Running
	} else {
		return result
	}
}

// Get a random sequence of nodes.
func genRandNodes(nodes []Node) []Node {
	count := len(nodes)
	if count == 0 {
		return nil
	}

	result := make([]Node, count)
	for i := count - 1; i > 0; i-- {
		if result[i] == nil {
			result[i] = nodes[i]
		}

		k := rand.Intn(i + 1)
		if k != i {
			if result[k] == nil {
				result[k] = nodes[k]
			}

			result[k], result[i] = result[i], result[k]
		}
	}

	if result[0] == nil {
		result[0] = nodes[0]
	}

	return result
}

// Random sequence runs child nodes one by one in a
// random sequence until a child returns failure. It
// returns the result of the last running node.
type RandSequenceNode struct {
	compositeNode
}

func NewRandSequenceNode() *RandSequenceNode {
	return &RandSequenceNode{
		compositeNode: newCompositeNode(),
	}
}

func (s *RandSequenceNode) NodeType() NodeType { return randSequence }

func (s *RandSequenceNode) AddChild(child Node) {
	s.compositeNode.addChild(child)
	child.SetParent(s)
}

// The randome sequence node task.
type randSequenceTask struct {
	node        *RandSequenceNode
	childs      []Node
	curChildIdx int
}

func (s *randSequenceTask) TaskType() TaskType { return Serial }

func (s *randSequenceTask) OnCreate(node Node) {
	s.node = node.(*RandSequenceNode)
	s.curChildIdx = 0
}

func (s *randSequenceTask) OnInit(nextChildNodes NodeList, ctx Context) bool {
	if s.childs = genRandNodes(s.node.children); len(s.childs) == 0 {
		return false
	} else {
		nextChildNodes.PushNode(s.childs[s.curChildIdx])
		return true
	}
}

func (s *randSequenceTask) OnUpdate(ctx Context) Result { return Running }
func (s *randSequenceTask) OnTerminate(ctx Context) {
	s.node = nil
	s.childs = nil
}

func (s *randSequenceTask) OnChildTerminated(result Result, nextChildNodes NodeList, ctx Context) Result {
	s.curChildIdx++

	if result == Success && s.curChildIdx < s.node.ChildCount() {
		nextChildNodes.PushNode(s.childs[s.curChildIdx])
		return Running
	} else {
		return result
	}
}

// Random selector node runs child nodes one by one in a
// random sequence until a child returns success. It returns
// the result of the last running node.
type RandSelectorNode struct {
	compositeNode
}

func NewRandSelectorNode() *RandSelectorNode {
	return &RandSelectorNode{
		compositeNode: newCompositeNode(),
	}
}

func (s *RandSelectorNode) NodeType() NodeType { return randSelector }

func (s *RandSelectorNode) AddChild(child Node) {
	s.compositeNode.addChild(child)
	child.SetParent(s)
}

// The random selector task.
type randSelectorTask struct {
	node        *RandSelectorNode
	childs      []Node
	curChildIdx int
}

func (s *randSelectorTask) TaskType() TaskType { return Serial }

func (s *randSelectorTask) OnCreate(node Node) {
	s.node = node.(*RandSelectorNode)
	s.curChildIdx = 0
}

func (s *randSelectorTask) OnInit(nextChildNodes NodeList, ctx Context) bool {
	s.childs = genRandNodes(s.node.children)
	if len(s.childs) == 0 {
		return false
	} else {
		nextChildNodes.PushNode(s.childs[s.curChildIdx])
		return true
	}
}

func (s *randSelectorTask) OnUpdate(ctx Context) Result { return Running }
func (s *randSelectorTask) OnTerminate(ctx Context) {
	s.node = nil
	s.childs = nil
}

func (s *randSelectorTask) OnChildTerminated(result Result, nextChildNodes NodeList, ctx Context) Result {
	s.curChildIdx++

	if result == Failure && s.curChildIdx < s.node.ChildCount() {
		nextChildNodes.PushNode(s.childs[s.curChildIdx])
		return Running
	} else {
		return result
	}
}

type weightNode struct {
	node   Node
	weight float32
}

type WeightSelectorNode struct {
	node
	children []*weightNode
}

func NewWeightSelectorNode() *WeightSelectorNode {
	return &WeightSelectorNode{node: newNode()}
}

func (n *WeightSelectorNode) NodeType() NodeType { return weightSelector }

func (n *WeightSelectorNode) ChildCount() int { return len(n.children) }

func (n *WeightSelectorNode) Child(idx int) (Node, float32) {
	assert.Assert(idx >= 0 && idx < n.ChildCount(), "index out of range")

	wnode := n.children[idx]
	return wnode.node, wnode.weight
}

func (n *WeightSelectorNode) AddChild(child Node, weight float32) {
	assert.Assert(child != nil, "child nil")
	assert.Assert(child.Parent() == nil, "child already has parent")
	assert.Assert(weight > 0, "weight <= 0")

	var w float32
	for _, child := range n.children {
		w += child.weight
	}

	assert.Assert(w+weight <= 1, "total weight > 1")

	child.SetParent(n)
	n.children = append(n.children, &weightNode{node: child, weight: weight})
}

type weightSelectorTask struct {
	node *WeightSelectorNode
}

// Get the TaskType.
func (t *weightSelectorTask) TaskType() TaskType {
	return Serial
}

// OnCreate is called immediately after the Task is created.
// node indicates the node on which the Task is created.
func (t *weightSelectorTask) OnCreate(node Node) {
	t.node = node.(*WeightSelectorNode)
}

// OnInit is called before the first update of the Task.
// nextChildNodes is used to return the child nodes that need
// to run next. ctx represents the running context of the
// behavior tree.
func (t *weightSelectorTask) OnInit(nextChildNodes NodeList, ctx Context) bool {
	if t.node.ChildCount() == 0 {
		return false
	}

	r := rand.Float32()
	w := float32(0)
	for i := 0; i < t.node.ChildCount(); i++ {
		node, weight := t.node.Child(i)
		w += weight
		if w > r {
			nextChildNodes.PushNode(node)
			return true
		}
	}

	if 1-w < 1e-4 {
		node, _ := t.node.Child(t.node.ChildCount() - 1)
		nextChildNodes.PushNode(node)
		return true
	}

	return false
}

// OnUpdate is called until the Task is terminated.
func (t *weightSelectorTask) OnUpdate(ctx Context) Result {
	return Running
}

// OnTerminate is called after ths last update of the Task.
func (t *weightSelectorTask) OnTerminate(ctx Context) { t.node = nil }

// OnChildTerminated is called when a sub Task is terminated.
//
// result Indicates the running result of the subtask.
// nextChildNodes is used to return the child nodes that need to
// run next.
//
// OnChildTerminated returns the decision result.
func (t *weightSelectorTask) OnChildTerminated(result Result, nextChildNodes NodeList, ctx Context) Result {
	return result
}

// The parrallel node runs child nodes together until a
// child returns failure. It returns success if all child
// nodes return success, or returns failure.
type ParallelNode struct {
	compositeNode
}

func NewParallelNode() *ParallelNode {
	return &ParallelNode{
		compositeNode: newCompositeNode(),
	}
}

func (p *ParallelNode) NodeType() NodeType { return parallel }

func (p *ParallelNode) AddChild(child Node) {
	p.compositeNode.addChild(child)
	child.SetParent(p)
}

// The parallel node task.
type parallelTask struct {
	node      *ParallelNode
	completed int
}

func (p *parallelTask) TaskType() TaskType { return Parallel }

func (p *parallelTask) OnCreate(node Node) {
	p.node = node.(*ParallelNode)
	p.completed = 0
}

func (p *parallelTask) OnInit(nextChildNodes NodeList, ctx Context) bool {
	childCount := p.node.ChildCount()
	if childCount == 0 {
		return false
	} else {
		nextChildNodes.PushNodes(p.node.children...)
		// for i := 0; i < childCount; i++ {
		// 	nextChildNodes.PushNode(p.node.Child(i))
		// }
		return true
	}
}

func (p *parallelTask) OnUpdate(ctx Context) Result { return Running }
func (p *parallelTask) OnTerminate(ctx Context)     { p.node = nil }
func (p *parallelTask) OnChildTerminated(result Result, nextChildNodes NodeList, ctx Context) Result {
	p.completed++

	if result == Success && p.completed < p.node.ChildCount() {
		return Running
	} else {
		return result
	}
}
