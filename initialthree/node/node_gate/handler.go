package node_gate

import (
	"errors"
	"github.com/golang/protobuf/proto"
	"initialthree/cluster"
	"initialthree/cluster/addr"
	"initialthree/protocol/cmdEnum"
	"initialthree/protocol/ss/ssmessage"
	"initialthree/zaplogger"
)

func onKickIPGate(from addr.LogicAddr, msg proto.Message) {
	kickMsg := msg.(*ssmessage.KickIPGate)
	zaplogger.GetSugar().Infof("onKickIPGate %v", kickMsg)

	ipMap := map[string]struct{}{}
	for _, ip := range kickMsg.GetRegexpIPs() {
		ipMap[ip] = struct{}{}
	}

	// todo 通配符

	closeChannel := make([]*channel, 0, 48)

	gameLock.Lock()
	for _, g := range gameMap {
		g.channelLock.Lock()
		for _, c := range g.channels {
			ip := c.tcpAddr.IP.String()
			if _, ok := ipMap[ip]; ok {
				closeChannel = append(closeChannel, c)
			}
		}
		g.channelLock.Unlock()
	}
	gameLock.Unlock()

	for _, c := range closeChannel {
		c.close(errors.New("ip kicked"))
	}

}

func init() {
	cluster.Register(cmdEnum.SS_KickIPGate, onKickIPGate)
}
