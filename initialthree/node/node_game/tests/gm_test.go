package node_game

import (
	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
	codecs "initialthree/codec/cs"
	"initialthree/node/node_game"
	"initialthree/protocol/cmdEnum"
	"initialthree/protocol/cs/message"
	"testing"
)

//go test -tags=game_test -coverpkg=initialthree/node/node_game/transaction/GameMaster -covermode=count -v -coverprofile=coverage.out -run=TestGm
//go tool cover -html=coverage.out
func TestGm(t *testing.T) {
	node_game.GameInit()

	userID := "test123"

	gate := node_game.NewGate()
	defer gate.Close()

	_, err := node_game.LoginAndCreateRole(gate, userID)

	assert.Nil(t, err)

	r := make(chan interface{}, 1)

	gate.Register(cmdEnum.CS_GameMaster, func(ret interface{}) {
		r <- ret
	})

	node_game.ForwardMsg(gate, userID, &message.GameMasterToS{
		Cmds: []*message.GmCmd{{
			Type:  proto.Int32(1),
			ID:    proto.Int32(5),
			Count: proto.Int32(100),
		}},
	})

	ret := <-r

	assert.Equal(t, uint16(message.ErrCode_OK), ret.(*codecs.Message).GetErrCode())

}
