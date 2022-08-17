package bevtree

import (
	"github.com/GodYY/gutils/assert"
)

// Default node types.
const (
	root            = NodeType("root")            // The root node of behavior tree.
	inverter        = NodeType("inverter")        // The inverter node.
	succeeder       = NodeType("succeeder")       // The succeeder node.
	repeater        = NodeType("repeater")        // The repeater node.
	repeatUntilFail = NodeType("repeatuntilfail") // The repeat-until-fail node.
	sequence        = NodeType("sequence")        // The sequence node.
	randSequence    = NodeType("randSequence")    // The random sequence node.
	selector        = NodeType("selector")        // The selector node.
	randSelector    = NodeType("randSelector")    // The random selector node.
	weightSelector  = NodeType("weightselector")  // The weight selector node.
	parallel        = NodeType("parallel")        // The parallel node.
	behavior        = NodeType("behavior")        // The behavior node.
	subtree         = NodeType("subtree")         // The subtree node.
)

// Node metadata.
type nodeMeta struct {
	// Node type.
	typ NodeType

	// The creator of node.
	creator func() Node

	// task pool, used to cache unused tasks of node type type.
	taskPool *taskPool
}

// Use creator to create node.
func (meta *nodeMeta) createNode() Node { return meta.creator() }

// Create a task of this type of node. First, get a cached task or
// create a new task. Then, call the OnCreate method of it.
func (meta *nodeMeta) createTask(node Node) Task {
	assert.Assert(node != nil, "node nil")
	task := meta.taskPool.get()
	task.OnCreate(node)
	return task
}

// Destroy a task of this type of node. First, call the OnDestroy
// method of the task. Then, put it to be cached to the pool.
func (meta *nodeMeta) destroyTask(task Task) {
	meta.taskPool.put(task)
}

// The metadata of behavior.
type bevMeta struct {
	// Behaivor type.
	typ BevType

	// The creator of behavior.
	creator func() Bev
}

// Use creator to create behavior.
func (meta *bevMeta) createBev() Bev {
	return meta.creator()
}

// Metadata of behavior tree system.
type meta struct {
	// Stores all node metas.
	nodeMetas map[NodeType]*nodeMeta

	// Stores all bev metas.
	bevMetas map[BevType]*bevMeta
}

func newMeta() *meta {

	m := &meta{
		nodeMetas: map[NodeType]*nodeMeta{},
		bevMetas:  map[BevType]*bevMeta{},
	}

	m.RegisterNodeType(root, func() Node { return newRootNode() }, func() Task { return &rootTask{} })
	m.RegisterNodeType(inverter, func() Node { return NewInverterNode() }, func() Task { return &inverterTask{} })
	m.RegisterNodeType(succeeder, func() Node { return NewSucceederNode() }, func() Task { return &succeederTask{} })
	m.RegisterNodeType(repeater, func() Node { return NewRepeaterNode(1) }, func() Task { return &repeaterTask{} })
	m.RegisterNodeType(repeatUntilFail, func() Node { return NewRepeatUntilFailNode(false) }, func() Task { return &repeatUntilFailTask{} })
	m.RegisterNodeType(sequence, func() Node { return NewSequenceNode() }, func() Task { return &sequenceTask{} })
	m.RegisterNodeType(selector, func() Node { return NewSelectorNode() }, func() Task { return &selectorTask{} })
	m.RegisterNodeType(randSequence, func() Node { return NewRandSequenceNode() }, func() Task { return &randSequenceTask{} })
	m.RegisterNodeType(randSelector, func() Node { return NewRandSelectorNode() }, func() Task { return &randSelectorTask{} })
	m.RegisterNodeType(parallel, func() Node { return NewParallelNode() }, func() Task { return &parallelTask{} })
	m.RegisterNodeType(behavior, func() Node { return new(BevNode) }, func() Task { return &bevTask{} })
	m.RegisterNodeType(subtree, func() Node { return new(SubtreeNode) }, func() Task { return &subtreeTask{} })
	m.RegisterNodeType(weightSelector, func() Node { return new(WeightSelectorNode) }, func() Task { return &weightSelectorTask{} })

	return m
}

// Register a type of node. It create metadata of the type of node.
func (m *meta) RegisterNodeType(nodeType NodeType, nodeCreator func() Node, taskCreator func() Task) {
	assert.AssertF(nodeType.Valid(), "invalid node type %s", nodeType.String())
	assert.AssertF(m.nodeMetas[nodeType] == nil, "node type \"%s\" registered", nodeType.String())
	assert.AssertF(nodeCreator != nil, "nodeCreator of node type \"%s\" nil", nodeType.String())
	assert.AssertF(taskCreator != nil, "taskCreator of node type \"%s\" nil", nodeType.String())

	node := nodeCreator()
	assert.AssertF(node != nil, "nodeCreator of node type \"%s\" create nil node", nodeType.String())
	assert.EqualF(node.NodeType(), nodeType, "node created of type \"%s\" has different type \"%s\"", nodeType.String(), node.NodeType().String())

	task := taskCreator()
	assert.AssertF(task != nil, "taskCreator of node type \"%s\" create nil task", nodeType.String())
	assert.AssertF(task.TaskType() == Single || task.TaskType() == Serial || task.TaskType() == Parallel,
		"node type \"%s\" create invalid type %d task", node.NodeType().String(), task.TaskType())

	meta := &nodeMeta{
		typ:      nodeType,
		creator:  nodeCreator,
		taskPool: newTaskPool(taskCreator),
	}

	m.nodeMetas[nodeType] = meta
}

func (m *meta) getNodeMeta(nodeType NodeType) *nodeMeta { return m.nodeMetas[nodeType] }

// Register a type of behavior. It create the metadata of the
// behavior.
func (m *meta) RegisterBevType(bevType BevType, creator func() Bev) {
	assert.AssertF(bevType.Valid(), "invalid bev type %s", bevType.String())
	assert.AssertF(m.bevMetas[bevType] == nil, "bev type \"%s\" already registered", bevType.String())
	assert.AssertF(creator != nil, "creator of bev type \"%s\" nil", bevType.String())

	bev := creator()
	assert.AssertF(bev != nil, "creator of bev type \"%s\" create nil bev", bevType.String())
	assert.AssertF(bev.BevType() == bevType, "bev created of type \"%s\" has different type \"%s\"", bevType.String(), bev.BevType().String())

	meta := &bevMeta{
		typ:     bevType,
		creator: creator,
	}

	m.bevMetas[bevType] = meta
}

func (m *meta) getBevMeta(bevType BevType) *bevMeta { return m.bevMetas[bevType] }
