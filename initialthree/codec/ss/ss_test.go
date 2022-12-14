package ss

//go test -covermode=count -v -coverprofile=coverage.out -run=.
import (
	"fmt"
	"github.com/golang/protobuf/proto"
	buffer "github.com/sniperHW/kendynet/buffer"
	"github.com/stretchr/testify/assert"
	codecs "initialthree/codec/cs"
	"initialthree/codec/relaysc"
	"initialthree/protocol/cs/message"
	"initialthree/protocol/ss/ssmessage"
	"testing"
)

func TestSS(t *testing.T) {
	forwordMsg := codecs.NewMessage(0, &message.BaseSyncToC{
		UserID: proto.String("sniperHW"),
		GameID: proto.Uint64(1),
		Name:   proto.String("test"),
	})

	relaymsg := relaysc.NewMessage(forwordMsg, codecs.NewEncoder("sc"), false, 1001)

	ssEncoder := NewEncoder("ss", "rpc_req", "rpc_resp")

	m := buffer.New()

	ssEncoder.EnCode(relaymsg, m)

	r := NewReceiver("ss", "rpc_req", "rpc_resp")

	ret, _, _ := r.unpack(m.Bytes(), 0, len(m.Bytes()))

	ss2gate := ret.(*Message).GetData().(*ssmessage.SsToGate)

	assert.Equal(t, uint64(1001), ss2gate.GetGateUsers()[0])

	var err error

	fmt.Println(len(ss2gate.GetMessage()[0]))

	ret, err = codecs.NewReceiver("sc").DirectUnpack(ss2gate.GetMessage()[0])

	if nil != err {
		fmt.Println(err)
	} else {

		assert.Equal(t, ret.(*codecs.Message).GetData().(*message.BaseSyncToC).GetUserID(), "sniperHW")

	}
}
