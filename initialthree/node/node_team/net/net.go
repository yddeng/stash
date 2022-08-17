package net

import (
	"github.com/golang/protobuf/proto"
	"initialthree/cluster"
	"initialthree/cluster/addr"
	"initialthree/codec/cs"
	"initialthree/protocol/ss/ssmessage"
	"reflect"
)

var (
	csEncoder *cs.Encoder
)

func SendRelayMessage(serverAddr addr.LogicAddr, userID string, roleID uint64, msg proto.Message) {
	relayMsg := &ssmessage.Relay{
		Targets: []*ssmessage.Target{{
			UserID: proto.String(userID),
			RoleID: proto.Uint64(roleID),
		}},
		MsgType: proto.String(reflect.TypeOf(msg).String()),
	}
	bytes, err := csEncoder.EnCode(cs.NewMessage(0, msg))
	if nil == err {
		relayMsg.Msg = bytes.Bytes()
		cluster.PostMessage(serverAddr, relayMsg)
	}
}

func init() {
	csEncoder = cs.NewEncoder("sc")
}
