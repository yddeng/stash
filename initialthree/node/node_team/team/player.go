package team

import (
	"fmt"
	"github.com/gogo/protobuf/proto"
	"initialthree/cluster/addr"
	"initialthree/node/node_team/net"
	"initialthree/protocol/cs/message"
)

type Player struct {
	pInfo    *message.TeamPlayer
	gameAddr addr.LogicAddr
	team     *Team

	matcher Matcher
}

func NewPlayer(game addr.LogicAddr, pInfo *message.TeamPlayer) *Player {
	return &Player{
		pInfo:    pInfo,
		gameAddr: game,
	}
}

func (p *Player) PlayerCount() int {
	return 1
}

func (p *Player) Info() *message.TeamPlayer {
	return p.pInfo
}

func (p *Player) LogStr() string {
	return fmt.Sprintf("player(%s %d %s)", p.UserID(), p.PlayerID(), p.gameAddr.String())
}

func (p *Player) UserID() string {
	return p.pInfo.GetUserID()
}

func (p *Player) Game() addr.LogicAddr {
	return p.gameAddr
}

func (p *Player) Name() string {
	return p.pInfo.GetName()
}

func (p *Player) PlayerID() uint64 {
	return p.pInfo.GetPlayerID()
}

func (p *Player) Team() *Team {
	return p.team
}

func (p *Player) SetMatcher(m Matcher) {
	p.matcher = m
}
func (p *Player) MatcherCancel() {
	if p.matcher != nil {
		p.matcher.Cancel()
		p.matcher = nil
	}
}

func (p *Player) UpdateInfo(pInfo *message.TeamPlayer) {
	p.pInfo = pInfo
	if p.team != nil {
		p.team.SyncToClients()
	}
}

func (p *Player) SendMsg(msg proto.Message) {
	net.SendRelayMessage(p.gameAddr, p.UserID(), p.PlayerID(), msg)
}
