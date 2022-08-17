package world

import (
	"github.com/golang/protobuf/proto"
	"initialthree/protocol/ss/rpc"
	"initialthree/rpc/enterMap"
	"initialthree/rpc/enterWorld"
)

type EnterWorld struct{}

func (this *EnterWorld) OnCall(replyer *enterWorld.EnterWorldReplyer, arg *rpc.EnterWorldReq) {

	logger.Debugln("OnCall EnterWorld", arg)
	resp := &rpc.EnterWorldResp{
		Ok: proto.Bool(false),
	}

	nodeMap, scene, err := randomNodeMap()
	if err != nil {
		logger.Infoln(err)
		replyer.Reply(resp)
		return
	}

	if !scene.enterScene() {
		logger.Infoln("enterScene false")
		replyer.Reply(resp)
		return
	}

	enter := &rpc.EnterMapReq{
		SceneIdx: proto.Int32(scene.SceneID),
		UserID:   proto.String(arg.GetUserID()),
		ID:       proto.Uint64(arg.GetID()),
		Pos: &rpc.Position{
			X: proto.Int32(arg.GetPos().GetX()),
			Y: proto.Int32(arg.GetPos().GetY()),
			Z: proto.Int32(arg.GetPos().GetZ()),
		},
		GameAddr: proto.String(arg.GetGameAddr()),
		GateAddr: proto.String(arg.GetGateAddr()),
		GateUid:  proto.Uint64(arg.GetGateUid()),
	}

	// 调用enterMap
	enterMap.AsynCall(nodeMap.LogicAddr, enter, time.Second*5, func(resp2 *rpc.EnterMapResp, e error) {
		if e != nil {
			logger.Errorln(e)
			replyer.Reply(resp)
		} else {
			if resp2.GetOk() {
				resp.Ok = proto.Bool(true)
				resp.MapLogicAddr = proto.Uint32(uint32(nodeMap.LogicAddr))
				resp.SceneIdx = proto.Int32(scene.SceneID)

				scene.AOICurr++
			}
			scene.AOIReq--
			replyer.Reply(resp)
		}
	})

}

func init() {
	enterWorld.Register(&EnterWorld{})
}
