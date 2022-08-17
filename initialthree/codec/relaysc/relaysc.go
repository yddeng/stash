package relaysc

import (
	"github.com/sniperHW/flyfish/pkg/buffer"
	fnet "github.com/sniperHW/flyfish/pkg/net"
	codecs "initialthree/codec/cs"
)

type Message struct {
	gateUsers []uint64
	msg       *codecs.Message
	scEncoder fnet.Encoder
	toAll     bool
}

func NewMessage(msg *codecs.Message, e fnet.Encoder, toAll bool, gateUsers ...uint64) *Message {
	return &Message{
		gateUsers: gateUsers,
		msg:       msg,
		scEncoder: e,
		toAll:     toAll,
	}
}

func (this *Message) ToAll() bool {
	return this.toAll
}

func (this *Message) GetGateUsers() []uint64 {
	return this.gateUsers
}

func (this *Message) Encode(buff *buffer.Buffer) error {
	return this.scEncoder.EnCode(this.msg, buff)
}
