package node_game

//go test -tags=game_test -coverpkg=initialthree/node/node_game/user -covermode=count -v -coverprofile=coverage.out -run=Test1
//go tool cover -html=coverage.out

import (
	//"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/sniperHW/flyfish/server/mock/kvnode"
	"github.com/stretchr/testify/assert"
	codecs "initialthree/codec/cs"
	"initialthree/node/node_game"
	"initialthree/protocol/cmdEnum"
	"initialthree/protocol/cs/message"
	ss_rpc "initialthree/protocol/ss/rpc"
	"initialthree/protocol/ss/ssmessage"
	"testing"
	"time"
)

func init() {
	node_game.GameInit()
}

func Test1(t *testing.T) {

	userID := "test1"

	gate := node_game.NewGate()
	defer gate.Close()

	time.Sleep(time.Second)

	_, err := node_game.LoginUser(gate, userID)

	assert.Nil(t, err)

	{

		r := make(chan interface{}, 1)

		gate.Register(cmdEnum.CS_CreateRole, func(ret interface{}) {
			r <- ret
		})

		node_game.ForwardMsg(gate, userID, &message.CreateRoleToS{
			Name: proto.String(userID),
		})

		ret := <-r

		assert.Equal(t, uint16(message.ErrCode_OK), ret.(*codecs.Message).GetErrCode())

	}

	{

		r := make(chan interface{}, 1)

		gate.RegisterMsgHandlerOnce(cmdEnum.SS_KickGateUser, func(ret interface{}) {
			r <- ret
		})

		ret, _ := node_game.LoginUser(gate, userID)

		assert.Equal(t, uint16(message.ErrCode_OK), uint16(ret.(*ss_rpc.GateUserLoginResp).GetCode()))

		ret = <-r

		assert.Equal(t, userID, ret.(*ssmessage.KickGateUser).GetUserID())

		gate.PostMessageToGame(&ssmessage.GateUserDisconnect{
			UserID:     proto.String(userID),
			GateUserID: proto.Uint64(1),
		})

	}

	{
		//kick game user
		gate.PostMessageToGame(&ssmessage.KickGameUser{
			UserID: proto.String(userID),
		})

		ret, _ := node_game.LoginUser(gate, userID)

		assert.Equal(t, uint16(message.ErrCode_OK), uint16(ret.(*ss_rpc.GateUserLoginResp).GetCode()))

	}

	{
		//kick game user
		gate.PostMessageToGame(&ssmessage.KickGameUser{
			UserID: proto.String(userID),
		})

		mock_kvnode.SetProcessDelay(time.Second * 2)

		ret, _ := node_game.LoginUser(gate, userID)

		assert.Equal(t, uint16(message.ErrCode_RETRY), uint16(ret.(*ss_rpc.GateUserLoginResp).GetCode()))

		time.Sleep(time.Second * 3)

		mock_kvnode.SetProcessDelay(time.Second * 7)

		ret, _ = node_game.LoginUser(gate, userID)

		assert.Equal(t, uint16(message.ErrCode_RETRY), uint16(ret.(*ss_rpc.GateUserLoginResp).GetCode()))

		mock_kvnode.SetProcessDelay(0)

	}

}
