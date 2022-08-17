package team

import (
	"github.com/gogo/protobuf/proto"
	"initialthree/protocol/cs/message"
	"sync/atomic"
	"time"
)

const teamMaxPlayerCount = 5

var seq uint32 = 1

func genID() uint32 {
	return atomic.AddUint32(&seq, 1)
}

type Team struct {
	teamID  uint32
	status  message.TeamStatus  // 状态
	target  *message.TeamTarget // 目标
	header  uint64              // 队长ID
	players []*Player           // 所有成员

	tInfo *message.Team

	autoAgreePlayerJoin bool                // 自动同意入队
	playerJoin          map[uint64]struct{} // 玩家入队申请

	playerInvite map[uint64]struct{} // 玩家邀请

	matcher Matcher
}

func NewTeam() *Team {
	return &Team{
		teamID:  genID(),
		status:  message.TeamStatus_Standby,
		target:  nil,
		header:  0,
		players: make([]*Player, 0, teamMaxPlayerCount),

		autoAgreePlayerJoin: false,
		playerJoin:          map[uint64]struct{}{},
		playerInvite:        map[uint64]struct{}{},
	}
}

func (t *Team) Info() *message.Team {
	ret := &message.Team{
		TeamID:   proto.Uint32(t.teamID),
		TeamName: proto.String("name"),
		Status:   t.status.Enum(),
		Header:   proto.Uint64(t.header),
		Target:   t.target,
		Players:  make([]*message.TeamPlayer, 0, len(t.players)),
	}

	for _, p := range t.players {
		ret.Players = append(ret.Players, p.Info())
	}
	return ret
}

func (t *Team) TeamID() uint32             { return t.teamID }
func (t *Team) Header() uint64             { return t.header }
func (t *Team) Status() message.TeamStatus { return t.status }
func (t *Team) PlayerCount() int           { return len(t.players) }
func (t *Team) PlayerFull() bool           { return len(t.players) == teamMaxPlayerCount }

func (t *Team) SetMatcher(m Matcher) {
	t.matcher = m
}
func (t *Team) MatcherCancel() {
	if t.matcher != nil {
		t.matcher.Cancel()
		t.matcher = nil
	}
}

// return , if fn is false
func (t *Team) RangePlayer(fn func(p *Player) bool) {
	for _, v := range t.players {
		if !fn(v) {
			return
		}
	}
}

func (t *Team) AddPlayer(p *Player) {
	if t.PlayerFull() {
		return
	}

	if t.PlayerCount() == 0 {
		t.header = p.PlayerID()
	}

	// 通知所有人（不包括自己） 玩家入队
	t.NotifyAll(&message.TeamPlayerJoinNotifyToC{
		Player: p.Info(),
	})

	t.players = append(t.players, p)
	p.team = t

	// 同步team
	t.SyncToClients()
}

func (t *Team) RemovePlayer(p *Player) {
	for i, v := range t.players {
		if v == p {
			t.players = append(t.players[0:i], t.players[i+1:]...)
			p.team = nil

			if t.PlayerCount() == 0 {
				Mgr().RemoveTeam(t.teamID)
			} else {
				// 队长离队
				if t.header == p.PlayerID() {
					t.header = t.players[0].PlayerID()
					t.players[0].SendMsg(&message.TeamHeaderChangedNotifyToC{})
				}
				// 通知所有人 玩家离队
				t.NotifyAll(&message.TeamPlayerLeaveNotifyToC{
					Player: p.Info(),
				})

				// 同步team
				t.SyncToClients()
			}

			return
		}
	}
}

func (t *Team) PlayerJoinApply(player *Player) {
	if t.autoAgreePlayerJoin {
		t.AddPlayer(player)
		return
	}

	pid := player.PlayerID()
	if _, ok := t.playerJoin[pid]; !ok {
		t.playerJoin[pid] = struct{}{}
	}

	// 通知给队长
	header := t.players[0]
	header.SendMsg(&message.TeamJoinApplyNotifyToC{
		Player: player.Info(),
	})
}

func (t *Team) PlayerJoinReply(player *Player, agree bool) {
	delete(t.playerJoin, player.PlayerID())

	// 通知给玩家
	player.SendMsg(&message.TeamJoinReplyNotifyToC{
		TeamID: proto.Uint32(t.teamID),
		Agree:  proto.Bool(agree),
	})

	if agree {
		t.AddPlayer(player)
	}

}

// 移交队长
func (t *Team) SetHeader(pid uint64) {
	if pid == t.header {
		return
	}

	for i, p := range t.players {
		if p.PlayerID() == pid {
			t.players[0], t.players[i] = t.players[i], t.players[0]
			t.header = pid
			p.SendMsg(&message.TeamHeaderChangedNotifyToC{})

			// 清空申请列表
			t.playerJoin = map[uint64]struct{}{}
			t.SyncToClients()
			return
		}
	}
}

func (t *Team) SetTarget(target *message.TeamTarget) {
	if target == nil {
		return
	}
	t.target = target
	t.SyncToClients()
}

func (t *Team) Tick(now time.Time) {
	// if this.vote != nil {
	//		if now.Unix() > this.vote.destoryTime {
	//			this.status = com.TeamStatus_Safe
	//			this.vote = nil
	//		}
	//	}
	//
	//	//队伍状态不在战斗中，清理掉线玩家
	//	if this.GetStatus() != com.TeamStatus_Battle {
	//		for _, r := range this.GetRoles() {
	//			if r.GetStatus() == com.RoleStatus_Offline {
	//				this.DelRole(r.GetRoleID())
	//
	//			}
	//		}
	//
	//	}
}

func (t *Team) SyncToClients() {
	t.NotifyAll(&message.TeamSyncToC{
		UpdateTeam: t.Info(),
	})
}

func (t *Team) NotifyAll(msg proto.Message) {
	for _, p := range t.players {
		p.SendMsg(msg)
	}
}

func (t *Team) NotifyAllExceptMe(msg proto.Message, me *Player) {
	for _, p := range t.players {
		if p != me {
			p.SendMsg(msg)
		}
	}
}
