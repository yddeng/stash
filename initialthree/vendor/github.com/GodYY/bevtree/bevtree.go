package bevtree

import (
	"path"
	"sync"
	"sync/atomic"
	"unsafe"

	"github.com/GodYY/gutils/assert"
	"github.com/pkg/errors"
)

// Node type.
type NodeType string

func (t NodeType) Valid() bool { return string(t) != "" }

func (t NodeType) String() string { return string(t) }

// Node represents the structure portion of behavior tree.
// Node defines the basic function of node of behavior tree.
// Nodes must implement the Node interface to be considered
// behavior tree nodes.
type Node interface {
	// Node type.
	NodeType() NodeType

	// Get the parent node.
	Parent() Node

	// Set the parent node.
	SetParent(Node)

	// Get the comment.
	Comment() string

	// Set the comment.
	SetComment(string)
}

// The common part of node.
type node struct {
	parent  Node
	comment string
}

func newNode() node {
	return node{}
}

func (n *node) Parent() Node { return n.parent }

func (n *node) SetParent(parent Node) {
	n.parent = parent
}

func (n *node) Comment() string           { return n.comment }
func (n *node) SetComment(comment string) { n.comment = comment }

// status indicate the status of node's runtime.
type status int8

const (
	// The initail state.
	sNone = status(iota)

	// Running.
	sRunning

	// Terminated.
	sTerminated

	// Was stopped.
	sStopped

	// Was destroyed.
	sDestroyed
)

// The strings represent the status values.
var statusStrings = [...]string{
	sNone:       "none",
	sRunning:    "running",
	sTerminated: "terminated",
	sStopped:    "stopped",
	sDestroyed:  "destroyed",
}

func (s status) String() string { return statusStrings[s] }

// lazyStop indicate node's runtime how to stop.
type lazyStop int8

const (
	// Don't need to stop.
	lzsNone = lazyStop(iota)

	// Stop before update.
	lzsBeforeUpdate

	// Stop after update.
	lzsAfterUpdate
)

// The strings represent the lazyStop values.
var lazyStopStrings = [...]string{
	lzsNone:         "none",
	lzsBeforeUpdate: "before-Update",
	lzsAfterUpdate:  "after-Update",
}

func (l lazyStop) String() string { return lazyStopStrings[l] }

// Result represents the running results of node's runtime and even behavior trees.
type Result int8

const (
	// Success can indicate that the behavior ran successfully,
	// or the node made a decision successfully, or the behavior
	// tree ran successfully.
	Success = Result(iota)

	// Failure can indicate that the behavior fails to run, or
	// the node fails to make a decision, or the behavior tree
	// fails to run.
	Failure

	// Running can indicate that a behavior run is running, or
	// that a node is making a decision, or that the behavior
	// tree is running.
	Running
)

// The strings repesents the Result values.
var resultStrings = [...]string{
	Success: "success",
	Failure: "failure",
	Running: "running",
}

func (r Result) String() string { return resultStrings[r] }

// TaskType indicate how the task will run.
type TaskType int8

const (
	// Single task, no any subtask.
	Single = TaskType(iota)

	// Serial task, there are subtasks and the subtasks run one
	// by one.
	Serial

	// Parallel task, there are subtasks and the subtasks run
	// together.
	Parallel
)

var taskTypeStrings = [...]string{
	Single:   "single",
	Serial:   "serial",
	Parallel: "parallel",
}

func (tt TaskType) String() string { return taskTypeStrings[tt] }

// Task represents the independent parts of behavir tree node.
// Task maintains runtime data and implements the logic of the
// corresponding node.
type Task interface {
	// Get the TaskType.
	TaskType() TaskType

	// OnCreate is called immediately after the Task is created.
	// node indicates the node on which the Task is created.
	OnCreate(node Node)

	// OnInit is called before the first update of the Task.
	// nextChildNodes is used to return the child nodes that need
	// to run next. ctx represents the running context of the
	// behavior tree.
	OnInit(nextChildNodes NodeList, ctx Context) bool

	// OnUpdate is called until the Task is terminated.
	OnUpdate(ctx Context) Result

	// OnTerminate is called either after ths last update or when
	// stopping it.
	OnTerminate(ctx Context)

	// OnChildTerminated is called when a sub Task is terminated.
	//
	// result Indicates the running result of the subtask.
	// nextChildNodes is used to return the child nodes that need to
	// run next.
	//
	// OnChildTerminated returns the decision result.
	OnChildTerminated(result Result, nextChildNodes NodeList, ctx Context) Result
}

// Root node, a special node in behavior tree. it has
// only one child and no parent. It returns result of
// child directly.
type rootNode struct {
	child Node
}

func newRootNode() *rootNode {
	return &rootNode{}
}

func (rootNode) NodeType() NodeType { return root }
func (rootNode) Parent() Node       { return nil }
func (rootNode) SetParent(Node)     {}
func (rootNode) Comment() string    { return "" }
func (rootNode) SetComment(string)  {}
func (r *rootNode) Child() Node     { return r.child }

func (r *rootNode) SetChild(child Node) {
	assert.Assert(child == nil || child.Parent() == nil, "child already has parent")

	if r.child != nil {
		r.child.SetParent(nil)
		r.child = nil
	}

	if child != nil {
		child.SetParent(r)
		r.child = child
	}
}

// rootNode Task.
type rootTask struct {
	node *rootNode
}

// rootNode Task is serail task.
func (r *rootTask) TaskType() TaskType { return Serial }
func (r *rootTask) OnCreate(node Node) { r.node = node.(*rootNode) }

func (r *rootTask) OnInit(nextChildNodes NodeList, ctx Context) bool {
	if r.node.Child() == nil {
		return false
	} else {
		nextChildNodes.PushNode(r.node.Child())
		return true
	}
}

func (r *rootTask) OnUpdate(ctx Context) Result { return Running }
func (r *rootTask) OnTerminate(ctx Context)     { r.node = nil }

func (r *rootTask) OnChildTerminated(result Result, nextChildNodes NodeList, ctx Context) Result {
	// Returns result of child directly.
	return result
}

// Tree interface, make the tree readonly while in use.
type Tree interface {
	// Get tree anme.
	Name() string

	// Get tree comment.
	Comment() string

	// Get root node.
	root() *rootNode

	internal
}

// tree, contains the structure data
type tree struct {
	// The name of the behavior tree.
	name string

	// The comment of the behavior tree.
	comment string

	// The _root node of behavior tree.
	_root *rootNode

	internalImpl
}

func NewTree(name string) *tree {
	assert.AssertF(name != "", "invalid name \"%s\"", name)

	tree := &tree{
		name:  name,
		_root: newRootNode(),
	}

	return tree
}

func (t *tree) Name() string { return t.name }

func (t *tree) SetName(name string) {
	assert.AssertF(name != "", "invalid name \"%s\"", name)
	t.name = name
}

func (t *tree) Comment() string           { return t.comment }
func (t *tree) SetComment(comment string) { t.comment = comment }

func (t *tree) Root() *rootNode { return t._root }

func (t *tree) root() *rootNode { return t._root }

type treeAsset struct {
	entry *TreeEntry
	once  *sync.Once
	tree  *tree
}

type Framework struct {
	*meta
	initialized    bool
	loadAll        bool
	configPathRoot string
	treeAssets     map[string]*treeAsset
}

func NewFramework() *Framework {
	return &Framework{
		meta: newMeta(),
	}
}

func (s *Framework) RegsiterNodeType(nodeType NodeType, nodeCreator func() Node, taskCreator func() Task) {
	if s.initialized {
		panic("bevtree framework initialized")
	}

	s.meta.RegisterNodeType(nodeType, nodeCreator, taskCreator)
}

func (s *Framework) RegisterBevType(bevType BevType, creator func() Bev) {
	if s.initialized {
		panic("bevtree framework initialized")
	}
	s.meta.RegisterBevType(bevType, creator)
}

func (s *Framework) Init(cfgPath string) error {
	if s.initialized {
		return errors.New("bevtree framework repeated initialization")
	}

	config, err := loadConfig(cfgPath)
	if err != nil {
		return errors.WithMessagef(err, "bevtree framework init")
	}

	s.configPathRoot = path.Dir(cfgPath)
	s.treeAssets = make(map[string]*treeAsset, len(config.TreeEntries))
	s.loadAll = config.LoadAll

	if len(config.TreeEntries) > 0 {
		for _, entry := range config.TreeEntries {
			var ta *treeAsset
			if s.loadAll {
				tree := new(tree)
				path := path.Join(s.configPathRoot, entry.Path)
				if err = s.DecodeXMLTreeFile(path, tree); err == nil {
					if tree.Name() != entry.Name {
						err = errors.Errorf("load tree \"%s\": name don't match config name \"%s\"", tree.Name(), entry.Name)
					}
				}

				if err != nil {
					return errors.WithMessage(err, "bevtree framework init")
				}

				ta = &treeAsset{entry: entry, tree: tree}
			} else {
				ta = &treeAsset{entry: entry, once: new(sync.Once)}
			}

			if s.treeAssets[entry.Name] != nil {
				return errors.Errorf("bevtree framework init: duplcate tree name \"%s\"", entry.Name)
			}

			s.treeAssets[entry.Name] = ta
		}
	}

	s.initialized = true

	return nil
}

func (s *Framework) loadTree(ta *treeAsset) (*tree, error) {
	var err error

	ta.once.Do(func() {
		tree := new(tree)

		path := path.Join(s.configPathRoot, ta.entry.Path)
		if err = s.DecodeXMLTreeFile(path, tree); err != nil {
			return
		}

		if tree.Name() != ta.entry.Name {
			err = errors.Errorf("loadTree \"%s\": tree name don't match config name \"%s\"", tree.Name(), ta.entry.Name)
		}

		atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&ta.tree)), unsafe.Pointer(tree))
	})

	if err != nil {
		return nil, err
	} else if ta.tree == nil {
		return nil, errors.Errorf("loadTree \"%s\" failed", ta.entry.Name)
	} else {
		return ta.tree, nil
	}
}

func (s *Framework) GetOrLoadTree(name string) (Tree, error) {
	if !s.initialized {
		return nil, errors.New("bevtree framework uninitialized")
	}

	return s.getOrLoadTree(name)
}

func (s *Framework) getOrLoadTree(name string) (*tree, error) {
	ta := s.treeAssets[name]
	if ta == nil {
		return nil, nil
	}

	if s.loadAll {
		return ta.tree, nil
	} else {
		tree := (*tree)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&ta.tree))))
		if tree != nil {
			return tree, nil
		}

		var err error
		if tree, err = s.loadTree(ta); err != nil {
			return nil, errors.WithMessage(err, "bevtree framework GetOrLoadTree")
		}

		return tree, nil
	}
}

func (s *Framework) CreateEntity(treeName string, userData interface{}) (Entity, error) {
	if !s.initialized {
		return nil, errors.New("bevtree framework uninitialized")
	}

	tree, err := s.getOrLoadTree(treeName)
	if err != nil {
		return nil, err
	} else if tree == nil {
		return nil, errors.Errorf("bevtree framework CreateEntity: tree \"%s\" not exist", treeName)
	} else {
		return newEntity(newContext(s, tree, userData)), nil
	}
}
