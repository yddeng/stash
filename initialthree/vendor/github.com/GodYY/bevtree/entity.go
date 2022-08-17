package bevtree

import (
	"fmt"
	"log"
	"reflect"

	"github.com/GodYY/gutils/assert"
	"github.com/GodYY/gutils/finalize"
)

// NodeList interface.
type NodeList interface {
	// Push one node.
	PushNode(Node)

	// Push nodes.
	PushNodes(...Node)
}

type nodeList struct {
	l list
}

func newNodeList() *nodeList {
	return &nodeList{}
}

func (nl *nodeList) PushNode(node Node) {
	assert.Assert(node != nil, "node nil")
	nl.l.pushBack(node)
}

func (nl *nodeList) PushNodes(nodes ...Node) {
	assert.Assert(len(nodes) > 0, "no nodes")
	for _, node := range nodes {
		nl.l.pushBack(node)
	}
}

func (nl *nodeList) pop() Node {
	if elem := nl.l.front(); elem != nil {
		return nl.l.remove(elem).(Node)
	} else {
		return nil
	}
}

func (nl *nodeList) len() int { return nl.l.getLen() }

func (nl *nodeList) clear() { nl.l.init() }

// Entity used to run a behavior tree.
type Entity interface {
	// Get the context.
	Context() Context

	// Update behavior tree and get a result from this
	// update.
	Update() Result

	// Stops running the bahavior tree.
	Stop()

	// If the entity is no longer used, call Release to
	// release resource of it.
	Release()

	internal
}

// agent represents common parts of behavior tree node.
// agent links Node and Task, maintains status infomation
// and implements the workflow of behavior tree node.
// All ruuning agents form a run-time behavior tree.
type agent struct {
	// Corresponding Node.
	node Node

	// Corresponding Task.
	task Task

	// Parent agent.
	parent *agent

	// Child agent list.
	firstChild *agent

	// Previous, next agent.
	prev, next *agent

	// Store the serial number of the latest updating.
	latestUpdateSeri uint32

	// Store the current status.
	st status

	// Store the lazyStop type.
	lzStop lazyStop

	// agent placeholder int the work queue.
	elem *element

	internalImpl
}

// onCreate is called immediately after the agent is created.
func (a *agent) onCreate(node Node, task Task) {
	a.node = node
	a.task = task
	a.latestUpdateSeri = 0
	a.st = sNone
	a.lzStop = lzsNone
}

// onDestroy is called before the agent is destroyed.
func (a *agent) onDestroy() {
	a.node = nil
	a.task = nil
}

// Indicates whether the agent is persistent. That is the
// the update method of the agent must be called whenever
// the behavior tree update before it terminated.
func (a *agent) isPersistent() bool { return a.task.TaskType() == Single }

func (a *agent) getNext() *agent {
	if a.parent != nil && a.next != a.parent.firstChild {
		return a.next
	}
	return nil
}

func (a *agent) getPrev() *agent {
	if a.parent != nil && a.prev != a.parent.firstChild {
		return a.prev
	}
	return nil
}

// Add a running child for the agent.
func (a *agent) addChild(child *agent) {
	assert.Assert(child != nil && child.parent == nil, "child nil or has parent")

	if a.firstChild == nil {
		child.prev = child
		child.next = child
		a.firstChild = child
	} else {
		child.prev = a.firstChild.prev
		child.next = a.firstChild
		child.prev.next = child
		child.next.prev = child
	}

	child.parent = a
}

// Remove a child for the agent.
func (a *agent) removeChild(child *agent) {
	assert.Assert(child != nil && child.parent == a, "child nil or parent not match")

	if child == a.firstChild && a.firstChild.next == a.firstChild {
		a.firstChild = nil
	} else {
		if child == a.firstChild {
			a.firstChild = a.firstChild.next
		}

		child.prev.next = child.next
		child.next.prev = child.prev
	}

	child.prev = nil
	child.next = nil
	child.parent = nil
}

func (a *agent) getParent() *agent         { return a.parent }
func (a *agent) getStatus() status         { return a.st }
func (a *agent) setStatus(st status)       { a.st = st }
func (a *agent) getLZStop() lazyStop       { return a.lzStop }
func (a *agent) setLZStop(lzStop lazyStop) { a.lzStop = lzStop }
func (a *agent) getElem() *element         { return a.elem }
func (a *agent) setElem(elem *element)     { a.elem = elem }

// Running logic of the agent.
func (a *agent) update(entity *entity) Result {
	if debug {
		log.Printf("agent nodetype:%v update %v %v", a.node.NodeType(), a.getStatus(), a.getLZStop())
	}

	st := a.getStatus()

	if debug {
		assert.NotEqualF(st, sDestroyed, "agent nodetype:%v already destroyed", a.node.NodeType())
	}

	// Update seri.
	a.latestUpdateSeri = entity.getUpdateSeri()

	// lazy Stop before Update.
	lzStop := a.getLZStop()
	if lzStop == lzsBeforeUpdate {
		return a.doLazyStop(entity)
	}

	// init.
	if st == sNone {
		if !a.task.OnInit(entity.getChildNodeList(), entity.Context()) {
			a.task.OnTerminate(entity.Context())
			a.setStatus(sTerminated)
			return Failure
		}

		if debug {
			switch a.task.TaskType() {
			case Single:
				assert.AssertF(entity.getChildNodeList().len() == 0, "node type \"%s\" has children", a.node.NodeType().String())

			case Serial:
				assert.AssertF(entity.getChildNodeList().len() == 1, "node type \"%s\" have no or more than one child", a.node.NodeType().String())

			case Parallel:
				assert.AssertF(entity.getChildNodeList().len() > 0, "node type \"%s\" have no children", a.node.NodeType().String())
			}
		}

		a.processNextChildren(entity)
	}

	// Update.
	result := a.task.OnUpdate(entity.Context())

	// lazy Stop after Update
	if lzStop == lzsAfterUpdate {
		return a.doLazyStop(entity)
	}

	if result == Running {
		a.setStatus(sRunning)
	} else {
		// terminate.
		a.task.OnTerminate(entity.Context())
		a.setStatus(sTerminated)
	}

	return result
}

// Procoess child nodes filtered by making decision. Child nodes
// are cached in Context.
func (a *agent) processNextChildren(entity *entity) {
	childNodeList := entity.getChildNodeList()
	for nextChildNode := childNodeList.pop(); nextChildNode != nil; nextChildNode = childNodeList.pop() {
		childAgent := entity.createAgent(nextChildNode)
		a.addChild(childAgent)
		entity.pushAgent(childAgent)
	}
}

// If the agent is running, stop it. remove all child agents,
// notify the task to terminate.
func (a *agent) stop(ctx Context) {
	if a.getStatus() != sRunning {
		return
	}

	if debug {
		log.Printf("agent nodetype:%v stop", a.node.NodeType())
	}

	child := a.firstChild
	for child != nil {
		agent := child
		child = child.getNext()
		a.removeChild(agent)
	}

	a.task.OnTerminate(ctx)
	a.setStatus(sStopped)
	a.setLZStop(lzsNone)
}

// Lazy-Stop the agent if it is running and not set with
// lazy-stop state yet.
func (a *agent) lazyStop(entity *entity) {
	if debug {
		log.Printf("agent nodetype:%v lazyStop", a.node.NodeType())
	}

	st := a.getStatus()
	if st == sStopped || st == sTerminated || a.getLZStop() != lzsNone {
		return
	}

	if a.latestUpdateSeri != entity.getUpdateSeri() {
		// Not updated on the latest updating.
		// Stop after update.
		a.setLZStop(lzsAfterUpdate)
	} else {
		// Updated on the latest updating.
		// Stop before update.
		a.setLZStop(lzsBeforeUpdate)
	}

	// Lazy-Stop need agent to update again.
	if a.elem == nil || a.getLZStop() == lzsBeforeUpdate {
		entity.pushAgent(a)
	}
}

// The implementation of Lazy-Stop on agent.
func (a *agent) doLazyStop(entity *entity) Result {
	a.lazyStopChildren(entity)
	a.task.OnTerminate(entity.Context())
	a.setStatus(sStopped)
	a.setLZStop(lzsNone)
	return Failure
}

func (a *agent) lazyStopChildren(entity *entity) {
	child := a.firstChild
	for child != nil {
		child.lazyStop(entity)
		node := child
		child = child.getNext()
		a.removeChild(node)
	}
}

// onChildTerminated is called when a child agent is terminated.
func (a *agent) onChildTerminated(child *agent, result Result, entity *entity) Result {
	if debug {
		log.Printf("agent nodetype:%v onChildTerminated %v", a.node.NodeType(), result)
		assert.Assert(a.task.TaskType() != Single, "shouldnt be singletask")
		assert.Assert(child.getParent() == a, "invalid child")
		assert.NotEqual(result, Running, "child terminated with running")
	}

	// Remove child.
	a.removeChild(child)

	// Not running, Failure.
	if a.getStatus() != sRunning {
		return Failure
	}

	// Lazy-Stopping, Running.
	if a.getLZStop() != lzsNone {
		return Running
	}

	// Invoke task.OnChildTerminated to make decision.
	if result = a.task.OnChildTerminated(result, entity.getChildNodeList(), entity.Context()); result == Running {
		if debug {
			switch a.task.TaskType() {
			case Serial:
				assert.AssertF(entity.getChildNodeList().len() == 1, "node type \"%s\" has no or more than one next child", a.node.NodeType().String())

			case Parallel:
				assert.AssertF(entity.getChildNodeList().len() == 0, "node type \"%s\" has next children", a.node.NodeType().String())
			}
		}

		a.processNextChildren(entity)
	} else {
		if debug {
			assert.AssertF(entity.getChildNodeList().len() == 0, "node type \"%s\" has next children on terminating.", a.node.NodeType())
		}

		// Lazy-Stop children, avoid nested calls.
		a.lazyStopChildren(entity)

		a.task.OnTerminate(entity.Context())
		a.setStatus(sTerminated)
		a.setLZStop(lzsNone)
	}

	return result
}

// The pool to cache destroyed agent.
var agentPool = newPool(func() interface{} { return &agent{} })

// Entity implementation.
type entity struct {
	// The context.
	ctx Context

	// Agent list.
	agentList *list

	// Agent updating boundary. Agents behind the boundary
	// will update at next updateing.
	agentUpdateBoundary *element

	// Child node list. It is used to temporarily store
	// subsequent child nodes.
	childNodeList *nodeList

	internalImpl
}

func newEntity(ctx Context) *entity {
	assert.Assert(ctx != nil, "ctx nil")

	entity := &entity{
		ctx:           ctx,
		agentList:     newList(),
		childNodeList: newNodeList(),
	}

	finalize.SetFinalizer(entity)

	return entity
}

// If the entity is no longer used, call Release to
// release resource of it.
func (e *entity) Release() {
	finalize.UnsetFinalizer(e)
	e.release()
}

func (e *entity) release() {
	e.clearAgent()
	e.agentList = nil
	e.childNodeList.clear()
	e.childNodeList = nil
	e.ctx.release()
	e.ctx = nil
	e.agentUpdateBoundary = nil
}

// Finalizer will be called by GC if there is no explicitly
// call Release.
func (e *entity) Finalizer() {
	if debug {
		log.Println("Env.Finalizer")
	}
	e.release()
}

func (e *entity) Context() Context { return e.ctx }

func (e *entity) getUpdateSeri() uint32 { return e.ctx.UpdateSeri() }

func (e *entity) getChildNodeList() *nodeList { return e.childNodeList }

func (e *entity) noAgents() bool {
	return e.agentList.getLen() == 0 || (e.agentList.getLen() == 1 && e.agentList.front() == e.agentUpdateBoundary)
}

func (e *entity) lazyPushUpdateBoundary() {
	if e.agentUpdateBoundary == nil {
		e.agentUpdateBoundary = e.agentList.pushBack(nil)
	}
}

func (e *entity) pushAgent_(agent *agent, nextRound bool) {
	assert.Assert(agent != nil, "agent nil")

	e.lazyPushUpdateBoundary()

	elem := agent.getElem()

	if elem == nil {
		if nextRound {
			elem = e.agentList.pushBack(agent)
		} else {
			elem = e.agentList.insertBefore(agent, e.agentUpdateBoundary)
		}

		agent.setElem(elem)
	} else {
		if nextRound {
			e.agentList.moveToBack(elem)
		} else {
			e.agentList.moveBefore(elem, e.agentUpdateBoundary)
		}
	}
}

// Push a agent that need to run in current update。
func (e *entity) pushAgent(agent *agent) {
	e.pushAgent_(agent, false)
}

// Pop a agent that need to run in current update.
func (e *entity) popAgent() *agent {
	e.lazyPushUpdateBoundary()

	elem := e.agentList.front()
	if elem == e.agentUpdateBoundary {
		e.agentList.moveToBack(elem)
		return nil
	}

	agent := elem.Value.(*agent)
	agent.setElem(nil)
	e.agentList.remove(elem)

	return agent
}

// Push a agent that need to run in the next update.
func (e *entity) pushPendingAgent(agent *agent) {
	e.pushAgent_(agent, true)
}

func (e *entity) removeAgent(agent *agent) {
	elem := agent.getElem()
	if elem != nil {
		e.agentList.remove(elem)
		agent.setElem(nil)
	}
}

func (e *entity) clearAgent() {
	elem := e.agentList.front()
	for elem != nil {
		next := elem.getNext()
		agent, ok := e.agentList.remove(elem).(*agent)
		elem = next

		if ok && agent != nil {
			assert.Assert(agent.isPersistent(), "agent is not persistent")

			for agent != nil {
				agent.setElem(nil)
				parent := agent.getParent()
				agent.stop(e.ctx)
				e.destroyAgent(agent)
				agent = parent
			}
		}
	}
}

func (e *entity) createAgent(node Node) *agent {
	nodeMeta := e.ctx.framework().getNodeMeta(node.NodeType())
	if nodeMeta == nil {
		panic(fmt.Sprintf("node type \"%s\" meta not found, %s", node.NodeType(), reflect.TypeOf(node).Elem().Name()))
	}

	task := nodeMeta.createTask(node)
	switch task.TaskType() {
	case Single, Serial, Parallel:
	default:
		panic(fmt.Sprintf("node type \"%s\" create invalid type %d task", node.NodeType().String(), task.TaskType()))
	}

	agent := agentPool.get().(*agent)
	agent.onCreate(node, task)

	return agent
}

func (e *entity) destroyAgent(agent *agent) {
	if debug {
		assert.AssertF(agent.getElem() == nil, "agent node type \"%s\" still in list on destroy", agent.node.NodeType().String())
	}

	node := agent.node
	nodeMETA := e.ctx.framework().getNodeMeta(node.NodeType())
	if nodeMETA == nil {
		panic(fmt.Sprintf("node type \"%s\" meta not found, %s", node.NodeType(), reflect.TypeOf(node).Elem().Name()))
	}

	nodeMETA.destroyTask(agent.task)
	agent.onDestroy()
	agentPool.put(agent)
}

// Update used to update the behavior tree and get a result
// from this update.
func (e *entity) Update() Result {
	e.lazyPushUpdateBoundary()
	e.ctx.update()

	if e.noAgents() {
		// No agents indicate the behavior tree was not run yet
		// or it had completed a traversal from root to root node.
		// Need to start a new traversal from the root node.
		e.pushAgent(e.createAgent(e.ctx.Tree().root()))
	}

	// The default result.
	result := Running

	// Run agent one by one until there are no agents at current
	// updating or back to root node.
	for agent := e.popAgent(); agent != nil; agent = e.popAgent() {
		r := agent.update(e)
		st := agent.getStatus()
		if st == sStopped {
			e.destroyAgent(agent)
			continue
		}

		if st == sTerminated {
			// agent terminated, submit result to parent for
			// making decision.

			// The flag indicating whether to back to the root
			// node.
			isBackToRoot := true

			// Submit result to parent until no parent.
			for agent.getParent() != nil {
				parent := agent.getParent()
				parentTerminated := parent.getStatus() != sRunning

				r = parent.onChildTerminated(agent, r, e)
				if parentTerminated || r == Running {
					// Parent already terminated or still running, stop.
					isBackToRoot = false
					break
				}

				assert.Assert(parent.getElem() == nil, "parent is still in work list")

				// Destroy the child agent.
				e.destroyAgent(agent)

				agent = parent
			}

			// Destroy the last terminated agent.
			e.destroyAgent(agent)

			if isBackToRoot {
				// Back to root node, update result.

				assert.Equal(result, Running, "Update terminated reapeatedly")
				assert.NotEqual(r, Running, "Update terminated with RRunning")

				result = r
			}
		} else if agent.isPersistent() {
			// agent still running and persistent, set it to
			// update at the next updating.

			e.pushPendingAgent(agent)
		}
	}

	assert.Assert(result == Running || e.noAgents(), "Update terminated but already has agents")

	return result
}

// Stop stops running the behavior tree.
func (e *entity) Stop() {
	e.ctx.reset()
	e.clearAgent()
	e.agentUpdateBoundary = nil
}
