package team

import (
	"fmt"
	"initialthree/cluster"
	"initialthree/pkg/timer"
	"initialthree/protocol/cs/message"
	"time"
)

type teamMgr struct {
	players         map[uint64]*Player
	teams           map[uint32]*Team
	targetMatchPool map[*message.TeamTarget]*Pool
}

var (
	mgr *teamMgr
)

func Mgr() *teamMgr {
	return mgr
}

func (mgr *teamMgr) GetTeam(tid uint32) *Team {
	return mgr.teams[tid]
}

func (mgr *teamMgr) AddTeam(t *Team) {
	if _, ok := mgr.teams[t.teamID]; ok {
		panic(fmt.Sprintf("team id %d repeated ", t.teamID))
	} else {
		mgr.teams[t.teamID] = t
	}
}

func (mgr *teamMgr) RemoveTeam(tId uint32) {
	delete(mgr.teams, tId)
}

func (mgr *teamMgr) GetPlayer(pid uint64) *Player {
	return mgr.players[pid]
}

func (mgr *teamMgr) AddPlayer(p *Player) {
	if _, ok := mgr.players[p.PlayerID()]; ok {
		panic(fmt.Sprintf("player id %d repeated ", p.PlayerID()))
	} else {
		mgr.players[p.PlayerID()] = p
	}
}

func (mgr *teamMgr) RemovePlayer(pid uint64) {
	delete(mgr.players, pid)
}

func (mgr *teamMgr) GetPool(target *message.TeamTarget) *Pool {
	for k, v := range mgr.targetMatchPool {
		if CheckTarget(k, target) {
			return v
		}
	}

	// new
	p := NewPool(target)
	mgr.targetMatchPool[target] = p
	return p
}

func (mgr *teamMgr) Tick(t *timer.Timer, _ interface{}) {
	now := time.Now()
	for _, v := range mgr.teams {
		v.Tick(now)
	}
	//util.Logger().Debugf("team length %d,player length %d", len(mgr.teams), len(mgr.players))

	for _, p := range mgr.targetMatchPool {
		p.Tick(now)
	}
}

func CheckTarget(t1, t2 *message.TeamTarget) bool {
	return t1.GetLevelID() == t2.GetLevelID()
}

func init() {
	mgr = &teamMgr{
		players:         map[uint64]*Player{},
		teams:           map[uint32]*Team{},
		targetMatchPool: map[*message.TeamTarget]*Pool{},
	}

	cluster.RegisterTimer(time.Second, mgr.Tick, nil)
}
