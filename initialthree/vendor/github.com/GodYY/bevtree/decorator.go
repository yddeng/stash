package bevtree

import "github.com/GodYY/gutils/assert"

// DecoratorNode interface indicates the functions that
// a decorator node in behavior tree must implement.
type DecoratorNode interface {
	Node

	// Get the child node.
	Child() Node

	// Set the child node.
	SetChild(Node)
}

// The common part of decoratoer node.
type decoratorNode struct {
	node
	child Node
}

func newDecoratorNode() decoratorNode {
	return decoratorNode{
		node: newNode(),
	}
}

func (d *decoratorNode) Child() Node { return d.child }

func (d *decoratorNode) setChild(child Node) bool {
	if child == nil || child.Parent() != nil {
		return false
	}

	if d.child != nil {
		d.child.SetParent(nil)
	}

	d.child = child

	return child != nil
}

// Inverter node runs child and returns the reverse value of
// child's result.
type InverterNode struct {
	decoratorNode
}

func NewInverterNode() *InverterNode {
	return &InverterNode{decoratorNode: newDecoratorNode()}
}

func (i *InverterNode) NodeType() NodeType { return inverter }

func (i *InverterNode) SetChild(child Node) {
	if i.decoratorNode.setChild(child) {
		child.SetParent(i)
	}
}

// The inverter node task.
type inverterTask struct {
	node *InverterNode
}

func (i *inverterTask) TaskType() TaskType { return Serial }
func (i *inverterTask) OnCreate(node Node) { i.node = node.(*InverterNode) }

func (i *inverterTask) OnInit(nextChildNodes NodeList, ctx Context) bool {
	if i.node.Child() == nil {
		return false
	} else {
		nextChildNodes.PushNode(i.node.Child())
		return true
	}
}

func (i *inverterTask) OnUpdate(ctx Context) Result { return Running }
func (i *inverterTask) OnTerminate(ctx Context)     { i.node = nil }

func (i *inverterTask) OnChildTerminated(result Result, _ NodeList, ctx Context) Result {
	if result == Success {
		return Failure
	} else {
		return Success
	}
}

// Succeeder node runs child node and always return success.
type SucceederNode struct {
	decoratorNode
}

func NewSucceederNode() *SucceederNode {
	return &SucceederNode{decoratorNode: newDecoratorNode()}
}

func (s *SucceederNode) NodeType() NodeType { return succeeder }

func (s *SucceederNode) SetChild(child Node) {
	if s.decoratorNode.setChild(child) {
		child.SetParent(s)
	}
}

// Succeeder node task.
type succeederTask struct {
	node *SucceederNode
}

func (s *succeederTask) TaskType() TaskType { return Serial }
func (s *succeederTask) OnCreate(node Node) { s.node = node.(*SucceederNode) }

func (s *succeederTask) OnInit(nextChildNodes NodeList, ctx Context) bool {
	if s.node.Child() == nil {
		return false
	} else {
		nextChildNodes.PushNode(s.node.Child())
		return true
	}
}

func (s *succeederTask) OnUpdate(ctx Context) Result { return Running }
func (s *succeederTask) OnTerminate(ctx Context)     { s.node = nil }

func (s *succeederTask) OnChildTerminated(result Result, _ NodeList, ctx Context) Result {
	return Success
}

// Repeater node runs child node in limited times until child
// returns failure. It returns the result of child directly.
type RepeaterNode struct {
	decoratorNode
	limited int
}

func NewRepeaterNode(limited int) *RepeaterNode {
	assert.Assert(limited > 0, "invalid limited")
	return &RepeaterNode{
		decoratorNode: newDecoratorNode(),
		limited:       limited,
	}
}

func (r *RepeaterNode) NodeType() NodeType { return repeater }

func (r *RepeaterNode) SetChild(child Node) {
	if r.decoratorNode.setChild(child) {
		child.SetParent(r)
	}
}

func (r *RepeaterNode) Limited() int { return r.limited }

// Repeater node task.
type repeaterTask struct {
	node  *RepeaterNode
	count int
}

func (r *repeaterTask) TaskType() TaskType { return Serial }
func (r *repeaterTask) OnCreate(node Node) { r.node = node.(*RepeaterNode); r.count = 0 }

func (r *repeaterTask) OnInit(nextChildNodes NodeList, ctx Context) bool {
	if r.node.Child() == nil {
		return false
	} else {
		nextChildNodes.PushNode(r.node.Child())
		return true
	}
}

func (r *repeaterTask) OnUpdate(ctx Context) Result { return Running }
func (r *repeaterTask) OnTerminate(ctx Context)     { r.node = nil }

func (r *repeaterTask) OnChildTerminated(result Result, nextChildNodes NodeList, ctx Context) Result {
	r.count++
	if result != Failure && r.count < r.node.limited {
		nextChildNodes.PushNode(r.node.Child())
		return Running
	} else {
		return result
	}
}

// RepeatUntilFail node runs child node until child returns
// failure. It returns success if successOnFail is true or
// failure.
type RepeatUntilFailNode struct {
	decoratorNode
	successOnFail bool
}

func NewRepeatUntilFailNode(successOnFail bool) *RepeatUntilFailNode {
	return &RepeatUntilFailNode{
		decoratorNode: newDecoratorNode(),
		successOnFail: successOnFail,
	}
}

func (r *RepeatUntilFailNode) NodeType() NodeType { return repeatUntilFail }

func (r *RepeatUntilFailNode) SetChild(child Node) {
	if r.decoratorNode.setChild(child) {
		child.SetParent(r)
	}
}

func (r *RepeatUntilFailNode) SuccessOnFail() bool { return r.successOnFail }

// RepeatUntilFail node task.
type repeatUntilFailTask struct {
	node *RepeatUntilFailNode
}

func (r *repeatUntilFailTask) TaskType() TaskType { return Serial }
func (r *repeatUntilFailTask) OnCreate(node Node) { r.node = node.(*RepeatUntilFailNode) }

func (r *repeatUntilFailTask) OnInit(nextChildNodes NodeList, ctx Context) bool {
	if r.node.Child() == nil {
		return false
	} else {
		nextChildNodes.PushNode(r.node.Child())
		return true
	}
}
func (r *repeatUntilFailTask) OnUpdate(ctx Context) Result { return Running }
func (r *repeatUntilFailTask) OnTerminate(ctx Context)     { r.node = nil }

func (r *repeatUntilFailTask) OnChildTerminated(result Result, nextChildNodes NodeList, ctx Context) Result {
	if result == Success {
		nextChildNodes.PushNode(r.node.Child())
		return Running
	} else if result == Failure && r.node.successOnFail {
		return Success
	} else {
		return result
	}
}
