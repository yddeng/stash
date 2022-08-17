package bevtree

import (
	"github.com/GodYY/gutils/assert"
)

// Subtree node is a kind of leaf node, used to run a subtree
// for reducing complexity.
type SubtreeNode struct {
	node

	// The subtree.
	subtree Tree

	// Whether to create a independent dataset. If set to true,
	// the behavior tree do not share dataset with subtree.
	independentDataSet bool
}

func NewSubtreeNode(subtree Tree, independentDataSet bool) *SubtreeNode {
	assert.Assert(subtree != nil, "subtree nil")
	return &SubtreeNode{
		node:               newNode(),
		subtree:            subtree,
		independentDataSet: independentDataSet,
	}
}

func (s *SubtreeNode) NodeType() NodeType { return subtree }

func (s *SubtreeNode) Subtree() Tree { return s.subtree }

func (s *SubtreeNode) IndependentDataSet() bool { return s.independentDataSet }

type subtreeTask struct {
	node   *SubtreeNode
	entity Entity
}

// Get the TaskType.
func (s *subtreeTask) TaskType() TaskType { return Single }

func (s *subtreeTask) OnCreate(node Node) {
	s.node = node.(*SubtreeNode)
}

// OnInit is called before the first update of the Task.
// nextChildNodes is used to return the child nodes that need
// to run next. ctx represents the running context of the
// behavior tree.
func (s *subtreeTask) OnInit(_ NodeList, ctx Context) bool {
	s.entity = newEntity(ctx.cloneWithTree(s.node.subtree, s.node.independentDataSet))
	return true
}

// OnUpdate is called until the Task is terminated.
func (s *subtreeTask) OnUpdate(ctx Context) Result {
	return s.entity.Update()
}

// OnTerminate is called after ths last update of the Task.
func (s *subtreeTask) OnTerminate(ctx Context) {
	if s.entity != nil {
		s.entity.Release()
		s.entity = nil
	}
	s.node = nil
}

// OnChildTerminated is called when a sub Task is terminated.
//
// result Indicates the running result of the subtask.
// nextChildNodes is used to return the child nodes that need to
// run next.
//
// OnChildTerminated returns the decision result.
func (s *subtreeTask) OnChildTerminated(result Result, _ NodeList, _ Context) Result {
	panic("shouldnt be invoked")
}
