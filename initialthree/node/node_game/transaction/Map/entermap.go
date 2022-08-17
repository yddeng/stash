package Map

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"initialthree/cluster"
	"initialthree/cluster/addr"
	codecs "initialthree/codec/cs"
	"initialthree/node/common/serverType"
	"initialthree/node/common/transaction"
	"initialthree/zaplogger"
	"time"

	"initialthree/node/node_game/module"
	"initialthree/node/node_game/module/mapData"
	"initialthree/node/node_game/temporary"
	"initialthree/node/node_game/user"
	"initialthree/protocol/cmdEnum"
	cs_message "initialthree/protocol/cs/message"
	"initialthree/protocol/ss/rpc"
	"initialthree/rpc/enterWorld"
)

type transactionEnterMap struct {
	transaction.TransactionBase
	user      *user.User
	req       *codecs.Message
	resp      *cs_message.EnterMapToC
	errcode   cs_message.ErrCode
	worldAddr *addr.LogicAddr
}

func (this *transactionEnterMap) GetModuleName() string {
	return "Map"
}

func (this *transactionEnterMap) Begin() {

	req := this.req.GetData().(*cs_message.EnterMapToS)
	this.resp = &cs_message.EnterMapToC{}
	zaplogger.GetSugar().Debugf("%s %d %s %v", this.user.GetUserID(), this.user.GetID(), "call entermap", req)

	wAddr, err := addr.MakeLogicAddr(fmt.Sprintf("%d.%d.%d", cluster.SelfAddr().Logic.Group(), serverType.World, req.GetMapID()))
	if err != nil {
		zaplogger.GetSugar().Debugf("no worldServer %d", req.GetMapID())
		this.errcode = cs_message.ErrCode_ERROR
		this.EndTrans()
		return
	}
	this.worldAddr = &wAddr

	mapInfoCache := this.user.GetTemporary(temporary.TempMapInfo)
	if mapInfoCache != nil {
		mapInfoCache.(*temporary.MapInfo).LeaveMap(func(ok bool) {
			if this.IsTimeout() {
				return
			}
			if ok {
				this.enterMap()
			} else {
				this.errcode = cs_message.ErrCode_ERROR
				this.EndTrans()
			}
		})
	} else {
		this.enterMap()
	}
}

func (this *transactionEnterMap) enterMap() {
	var X, Y, Z, Angle int32
	req := this.req.GetData().(*cs_message.EnterMapToS)
	reqMap := req.GetMapID()
	mapModule := this.user.GetSubModule(module.MapData).(*mapData.UserMapData)
	if reqMap == 0 {
		reqMap = mapModule.GetMapID()
		X, Y, Z, Angle = mapModule.GetPos()
	}

	gate := this.user.GetGate()
	enterReq := &rpc.EnterWorldReq{
		UserID: proto.String(this.user.GetUserID()),
		Pos: &rpc.Position{
			X: proto.Int32(X),
			Y: proto.Int32(Y),
			Z: proto.Int32(Z),
		},
		GameAddr: proto.String(cluster.SelfAddr().Logic.String()),
		GateAddr: proto.String(gate.GateAddr.String()),
		GateUid:  proto.Uint64(gate.GateUid),
		ID:       proto.Uint64(this.user.GetID()),
	}
	this.AsynWrap(enterWorld.AsynCall)(*this.worldAddr, enterReq, 8*time.Second, func(resp *rpc.EnterWorldResp, e error) {
		if e != nil {
			zaplogger.GetSugar().Debugf("%s %d %s %v", this.user.GetUserID(), this.user.GetID(), "entermap fail:", e)
			this.errcode = cs_message.ErrCode_ERROR
			this.EndTrans()
		} else {
			if resp.GetOk() {
				mapModule.SetMapID(reqMap)
				mapModule.SetPos(X, Y, Z, Angle)

				this.user.SetTemporary(temporary.TempMapInfo, temporary.NewMapInfo(
					this.user, addr.LogicAddr(resp.GetMapLogicAddr()), reqMap, resp.GetSceneIdx()))

				this.resp.Pos = &cs_message.Position{
					X: proto.Int32(X),
					Y: proto.Int32(Y),
					Z: proto.Int32(Z),
				}
				this.resp.MapID = proto.Int32(reqMap)
				this.resp.Angle = proto.Int32(Angle)
				this.resp.Scene = proto.Int32(resp.GetSceneIdx())
				this.resp.MapLogic = proto.Uint32(resp.GetMapLogicAddr())

				zaplogger.GetSugar().Infof("%s %d %s", this.user.GetUserID(), this.user.GetID(), "entermap ok")
				this.errcode = cs_message.ErrCode_OK
			} else {
				zaplogger.GetSugar().Infof("%s %d %s", this.user.GetUserID(), this.user.GetID(), "entermap failed")
				this.errcode = cs_message.ErrCode_RETRY
			}
			this.EndTrans()
		}
	})
}

func init() {
	user.RegisterTransFunc(cmdEnum.CS_EnterMap, func(user *user.User, msg *codecs.Message) transaction.Transaction {
		return &transactionEnterMap{
			user: user,
			req:  msg,
		}
	})
}
