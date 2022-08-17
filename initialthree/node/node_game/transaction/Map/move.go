package Map

import (
	"github.com/golang/protobuf/proto"
	codecs "initialthree/codec/cs"
	"initialthree/node/common/transaction"
	"initialthree/node/node_game/module"
	"initialthree/node/node_game/module/mapData"
	"initialthree/node/node_game/temporary"
	"initialthree/node/node_game/user"
	"initialthree/protocol/cmdEnum"
	cs_message "initialthree/protocol/cs/message"
	ss_rpc "initialthree/protocol/ss/rpc"
	"initialthree/rpc/move"
	"time"
)

type transactionMove struct {
	transaction.TransactionBase
	user *user.User
	req  *codecs.Message
}

func (this *transactionMove) GetModuleName() string {
	return "Map"
}

func (this *transactionMove) Begin() {
	mapInfoCache := this.user.GetTemporary(temporary.TempMapInfo)
	if mapInfoCache == nil {
		this.EndTrans()
		return
	}

	mapInfo := mapInfoCache.(*temporary.MapInfo)

	req := this.req.GetData().(*cs_message.MoveToS)
	arg := &ss_rpc.MoveReq{}
	arg.SceneIdx = proto.Int32(mapInfo.SceneIdx)
	arg.MapID = proto.Int32(mapInfo.MapID)
	arg.UserID = proto.String(this.user.GetUserID())
	arg.RoleID = proto.Uint64(this.user.GetID())
	arg.Pos = &ss_rpc.Position{
		X: proto.Int32(req.GetPos().GetX()),
		Y: proto.Int32(req.GetPos().GetY()),
		Z: proto.Int32(req.GetPos().GetZ()),
	}
	arg.Angle = proto.Int32(req.GetAngle())

	this.AsynWrap(move.AsynCall)(mapInfo.MapAddr, arg, 8*time.Second, func(resp *ss_rpc.MoveResp, err error) {
		if nil == err && resp.GetOk() {
			pos := resp.GetPos()
			mapModule := this.user.GetSubModule(module.MapData).(*mapData.UserMapData)
			mapModule.SetPos(pos.GetX(), pos.GetY(), pos.GetZ(), req.GetAngle())
		}
		this.EndTrans()
	})

}

func (this *transactionMove) End() {}

func (this *transactionMove) Timeout() {}

func init() {
	user.RegisterTransFunc(cmdEnum.CS_Move, func(user *user.User, msg *codecs.Message) transaction.Transaction {
		return &transactionMove{
			user: user,
			req:  msg,
		}
	})
}
