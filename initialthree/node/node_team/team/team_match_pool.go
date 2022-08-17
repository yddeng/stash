package team

import (
	"initialthree/protocol/cs/message"
	"sync"
	"time"
)

/*
1.单人匹配：
选中活动目标后进行匹配，匹配相同目标副本且处于匹配状态的队伍
出现满足条件的队伍且队伍上限人数未满，不需要发送申请，直接加入队伍
2.队伍匹配：
创建队伍且确定队伍目标后进行匹配，匹配相同目标副本且处于匹配状态的单人玩家
未确定目标进行匹配，提示“请先调整目标”
队伍匹配成功匹配到1个玩家后，结束匹配
*/

type Matcher interface {
	Cancel()
}

// 用 map 随机匹配
type Pool struct {
	target        *message.TeamTarget
	teamEleList   map[uint32]*element // 队伍
	playerEleList map[uint64]*element // 单人玩家
}

func NewPool(target *message.TeamTarget) *Pool {
	return &Pool{
		target:        target,
		teamEleList:   map[uint32]*element{},
		playerEleList: map[uint64]*element{},
	}
}

func (this *Pool) Push(elt interface{}) Matcher {
	ele := &element{
		elt:        elt,
		createTime: time.Now(),
		p:          this,
	}
	switch elt.(type) {
	case *Team:
		insTeam := elt.(*Team)
		if _, ok := this.teamEleList[insTeam.TeamID()]; !ok {
			this.teamEleList[insTeam.TeamID()] = ele
			return ele
		}
	case *Player:
		player := elt.(*Player)
		if _, ok := this.playerEleList[player.PlayerID()]; !ok {
			this.playerEleList[player.PlayerID()] = ele
			return ele
		}
	default:
	}
	return nil
}

// 给每一个队伍添加一个玩家，玩家移除匹配
// 队伍满，移除匹配
func (this *Pool) Tick(now time.Time) {
	if len(this.playerEleList) == 0 || len(this.teamEleList) == 0 {
		return
	}

	teams := getSlice()
	for _, tEle := range this.teamEleList {
		insTeam := tEle.elt.(*Team)

		players := getSlice()
		for _, pEle := range this.playerEleList {
			player := pEle.elt.(*Player)
			players = append(players, pEle)
			if player.Team() == nil {
				insTeam.AddPlayer(player)
				break
			}
		}

		// 玩家取消
		for _, p := range players {
			p.elt.(*Player).MatcherCancel()
		}
		resetSlice(players)

		// 队伍满
		if insTeam.PlayerFull() {
			teams = append(teams, tEle)
		}

		// 没有可匹配玩家
		if len(this.playerEleList) == 0 {
			break
		}
	}

	// 队伍取消
	for _, p := range teams {
		p.elt.(*Team).MatcherCancel()
	}
	resetSlice(teams)

}

type element struct {
	elt        interface{}
	createTime time.Time
	p          *Pool
}

func (this *element) Cancel() {
	switch this.elt.(type) {
	case *Team:
		insTeam := this.elt.(*Team)
		delete(this.p.teamEleList, insTeam.TeamID())
	case *Player:
		player := this.elt.(*Player)
		delete(this.p.playerEleList, player.PlayerID())
	}
}

type List struct {
	data []*element
}

func (this *List) Push(elem *element) {
	this.data = append(this.data, elem)
}

func (this *List) Remove(elem *element) {
	for i, e := range this.data {
		if e == elem {
			this.data = append(this.data[:i], this.data[i+1:]...)
			break
		}
	}
}

func (this *List) Len() int {
	return len(this.data)
}

func (this *List) Front() *element {
	return this.data[0]
}

func (this *List) Range(fn func(elem *element) bool) {
	for _, e := range this.data {
		if !fn(e) {
			break
		}
	}
}

var slicePool = sync.Pool{
	New: func() interface{} {
		return make([]*element, 0, 4)
	},
}

func getSlice() []*element {
	return slicePool.Get().([]*element)
}

func resetSlice(s []*element) {
	s = s[:0]
	slicePool.Put(s)
}
