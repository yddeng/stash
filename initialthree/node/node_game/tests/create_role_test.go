package node_game

import (
	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
	"initialthree/protocol/cs/message"
	//ss_rpc "initialthree/protocol/ss/rpc"
	"github.com/sniperHW/flyfish/server/mock/kvnode"
	codecs "initialthree/codec/cs"
	"initialthree/node/node_game"
	"initialthree/protocol/cmdEnum"
	"testing"
	"time"
	//"fmt"
)

//go test -tags=game_test -coverpkg=initialthree/node/node_game/transaction/CreateRole -covermode=count -v -coverprofile=coverage.out -run=TestCreateRole
//go tool cover -html=coverage.out
func TestCreateRole(t *testing.T) {
	node_game.GameInit()

	userID := "test123"

	gate := node_game.NewGate()
	defer gate.Close()

	_, err := node_game.LoginUser(gate, userID)

	assert.Nil(t, err)

	r := make(chan interface{}, 1)

	gate.Register(cmdEnum.CS_CreateRole, func(ret interface{}) {
		r <- ret
	})

	node_game.ForwardMsg(gate, userID, &message.CreateRoleToS{
		Name: proto.String(userID),
	})

	ret := <-r

	assert.Equal(t, uint16(message.ErrCode_OK), ret.(*codecs.Message).GetErrCode())

	node_game.ForwardMsg(gate, userID, &message.CreateRoleToS{
		Name: proto.String(userID),
	})

	ret = <-r

	assert.Equal(t, uint16(message.ErrCode_OK), ret.(*codecs.Message).GetErrCode())

	_, err = node_game.LoginUser(gate, "test567")

	node_game.ForwardMsg(gate, "test567", &message.CreateRoleToS{
		Name: proto.String("test123"),
	})

	ret = <-r

	assert.Equal(t, uint16(message.ErrCode_Create_Role_Name_Repeat), ret.(*codecs.Message).GetErrCode())

	_, err = node_game.LoginUser(gate, "test10")

	mock_kvnode.SetProcessDelay(time.Second * 3)

	node_game.ForwardMsg(gate, "test10", &message.CreateRoleToS{
		Name: proto.String("test10"),
	})

	ret = <-r

	assert.Equal(t, uint16(message.ErrCode_RETRY), ret.(*codecs.Message).GetErrCode())

	mock_kvnode.SetProcessDelay(0)

}
