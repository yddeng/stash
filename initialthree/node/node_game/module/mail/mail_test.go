package mail

import (
	"github.com/golang/protobuf/proto"
	"initialthree/protocol/cs/message"
	"testing"
)

func TestSort(t *testing.T) {
	mails := []*message.Mail{
		{
			Id:         proto.Uint32(1),
			CreateTime: proto.Int64(1),
			Read:       proto.Bool(false),
		},
		{
			Id:         proto.Uint32(2),
			CreateTime: proto.Int64(2),
			Read:       proto.Bool(true),
		},
		{
			Id:         proto.Uint32(3),
			CreateTime: proto.Int64(5),
			Read:       proto.Bool(false),
		},
		{
			Id:         proto.Uint32(4),
			CreateTime: proto.Int64(4),
			Read:       proto.Bool(true),
		},
		{
			Id:         proto.Uint32(5),
			CreateTime: proto.Int64(3),
			Read:       proto.Bool(true),
		},
	}

	mails = sortMail(mails)

	count := len(mails) - 3
	mails = mails[count:]

	for _, m := range mails {
		t.Log(m)
	}
}
