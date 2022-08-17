package role

import (
	"initialthree/cluster"
	"initialthree/cluster/addr"
	"initialthree/codec/cs"
	codecs "initialthree/codec/cs"
	ss_message "initialthree/protocol/ss/ssmessage"

	"github.com/golang/protobuf/proto"
)

var encoder = cs.NewEncoder("sc")

func (this *Role) SendToClient(msg proto.Message) {
	if nil == this.Gate {
		return
	}
	//logger.Debugln("SendToClient", msg)

	forwordMsg := codecs.NewMessage(0, msg)
	buffer, err := encoder.EnCode(forwordMsg)
	if nil == err {
		ssMsg := &ss_message.SsToGate{
			GateUsers: []uint64{this.Gate.GateUid},
			Message:   [][]byte{buffer.Bytes()},
		}
		cluster.PostMessage(this.Gate.GateAddr, ssMsg)
	}
}

func SendMsgToClients(gate addr.LogicAddr, users []uint64, msg []byte, count int) {

	ssMsg := &ss_message.SsToGate{
		GateUsers: make([]uint64, 0, len(users)),
		Message:   [][]byte{msg},
	}

	for _, u := range users {
		ssMsg.GateUsers = append(ssMsg.GateUsers, u)
		if len(ssMsg.GateUsers) == count {
			cluster.PostMessage(gate, ssMsg)
			ssMsg.GateUsers = ssMsg.GateUsers[0:0]
		}
	}

	if len(ssMsg.GateUsers) > 0 {
		cluster.PostMessage(gate, ssMsg)
	}
}
