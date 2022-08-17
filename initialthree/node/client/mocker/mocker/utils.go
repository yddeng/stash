package mocker

import (
	"fmt"
	codecs "initialthree/codec/cs"

	"github.com/gogo/protobuf/proto"
)

func newMessage(req proto.Message) *codecs.Message {
	return codecs.NewMessage(newSeq(), req)
}

var seq = uint32(0)

func newSeq() uint32 {
	seq++
	return seq
}

func mustBeZero(errCode uint16) {
	if errCode != 0 {
		panic(fmt.Sprintf("error code is %d ,not 0", errCode))
	}
}

func must(i interface{}, e error) interface{} {
	if e != nil {
		panic("error not be nil")
	}
	return i
}
