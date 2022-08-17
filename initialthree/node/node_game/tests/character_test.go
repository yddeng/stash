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

//go test -tags=game_test -coverpkg=initialthree/node/node_game/transaction/Character -covermode=count -v -coverprofile=coverage.out -run=TestCharacter
//go tool cover -html=coverage.out
func TestCharacter(t *testing.T) {
	node_game.GameInit()

	userID := "test123"

	gate := node_game.NewGate()
	defer gate.Close()

	_, err := node_game.LoginAndCreateRole(gate, userID)

	assert.Nil(t, err)

	gmEW := gate.RegisterWait(cmdEnum.CS_GameMaster)

	node_game.ForwardMsg(gate, userID, &message.GameMasterToS{
		Cmds: []*message.GmCmd{
			{
				Type:  proto.Int32(1),
				ID:    proto.Int32(5),
				Count: proto.Int32(10000), // 金币
			}, {
				Type:  proto.Int32(2),
				ID:    proto.Int32(1),
				Count: proto.Int32(1), // 角色

			},
		},
	})

	assert.Equal(t, uint16(message.ErrCode_OK), gmEW.Pop().(*codecs.Message).GetErrCode())

	ctsEW := gate.RegisterWait(cmdEnum.CS_CharacterTeamSet)

	node_game.ForwardMsg(gate, userID, &message.CharacterTeamSetToS{
		CharacterTeam: &message.CharacterTeam{CharacterList: []int32{1, 0, 0}},
	})

	assert.Equal(t, uint16(message.ErrCode_OK), ctsEW.Pop().(*codecs.Message).GetErrCode())

}
