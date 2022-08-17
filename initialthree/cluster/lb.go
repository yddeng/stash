package cluster

import (
	"errors"
	"initialthree/cluster/addr"
	"sort"
	"sync"
	"time"
)

var (
	errLbNone = errors.New("load balance none")

	LBReportDur = time.Second
)

var (
	lbGroups  = map[uint32]*lbGroup{}
	lbGroupMu = sync.Mutex{}
)

type lbGroup struct {
	nodes []*lbNode
	mu    sync.Mutex
}

type lbNode struct {
	logic    addr.LogicAddr
	priority int
	capacity int //容量值
	deadline time.Time
}

func setLB(logic addr.LogicAddr, priority, count int) {
	deadline := time.Now().Add(LBReportDur * 2)

	tt := logic.Type()

	lbGroupMu.Lock()
	g, ok := lbGroups[tt]
	if !ok {
		g = &lbGroup{}
		lbGroups[tt] = g
	}
	lbGroupMu.Unlock()

	g.mu.Lock()
	defer g.mu.Unlock()

	find := false
	for _, n := range g.nodes {
		if n.logic == logic {
			n.priority = priority
			n.capacity = count
			n.deadline = deadline
			find = true
			break
		}
	}

	if !find {
		g.nodes = append(g.nodes, &lbNode{
			logic:    logic,
			priority: priority,
			capacity: count,
			deadline: deadline,
		})
	}

	sort.Slice(g.nodes, func(i, j int) bool {
		return g.nodes[i].priority > g.nodes[j].priority
	})

}

func LBMod(tt uint32) (addr.LogicAddr, error) {
	lbGroupMu.Lock()
	g, ok := lbGroups[tt]
	lbGroupMu.Unlock()
	if !ok {
		return 0, errLbNone
	}

	now := time.Now()

	g.mu.Lock()
	defer g.mu.Unlock()

	sIndex := -1
	for i, n := range g.nodes {
		if n.capacity == 0 || now.After(n.deadline) {
			continue
		}
		sIndex = i
		break
	}

	if sIndex == -1 {
		g.nodes = []*lbNode{}
		return 0, errLbNone
	} else if sIndex != 0 {
		g.nodes = g.nodes[sIndex:]
	}

	n := g.nodes[0]
	n.capacity--
	if n.capacity == 0 {
		g.nodes = g.nodes[1:]
	}

	return n.logic, nil
}
