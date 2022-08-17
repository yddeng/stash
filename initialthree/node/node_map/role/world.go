package role

import (
	"github.com/golang/protobuf/proto"
	"initialthree/cluster/addr"
	"initialthree/node/node_map/scene"
	"initialthree/protocol/ss/rpc"
	"initialthree/rpc/worldObjPush"
)

type WorldObjBase struct {
	UserID   string
	ID       uint64
	World    addr.LogicAddr //world地址
	ScenePos scene.Position
}

var wobjs = map[string]*WorldObjBase{}

type WorldObjPush struct{}

func (this *WorldObjPush) OnCall(replyer *worldObjPush.WorldObjPushReplyer, arg *rpc.WorldObjPushReq) {
	//logger.Debugln("OnCall WorldObjPush", arg)
	switch arg.GetT() {
	case rpc.PushType_EnterMap:
		req := arg.GetEnter()
		_, ok := wobjs[req.GetUserID()]
		if ok {
			logger.Debugf("sceneO:%s is not nil", req.GetUserID())
			replyer.Reply(&rpc.WorldObjPushResp{
				Ok: proto.Bool(false),
			})
			return
		}

		base := &WorldObjBase{
			UserID: req.GetUserID(),
			ID:     req.GetID(),
		}
		base.World, _ = addr.MakeLogicAddr(req.GetWorldAddr())
		base.ScenePos = scene.Position{
			X: req.GetPos().GetX(),
			Y: req.GetPos().GetY(),
			Z: req.GetPos().GetZ(),
		}

		wobjs[req.GetUserID()] = base

		replyer.Reply(&rpc.WorldObjPushResp{
			Ok: proto.Bool(true),
		})

		for _, s := range scene.GetScenes() {
			scene_ := s
			r := &WorldObj{
				enter: map[scene.SceneObject]struct{}{},
				leave: map[scene.SceneObject]struct{}{},
			}
			r.ID = req.GetID()
			r.UserID = req.GetUserID()
			r.World, _ = addr.MakeLogicAddr(req.GetWorldAddr())

			r.ScenePos = scene.Position{
				X: req.GetPos().GetX(),
				Y: req.GetPos().GetY(),
				Z: req.GetPos().GetZ(),
			}

			if scene_.EnterScene(r) {
				scene_.StartAoi(r)
			}
		}
		logger.Debugf("%s EnterMap ok", req.GetUserID())

	case rpc.PushType_LeaveMap:

	case rpc.PushType_DoAction:
		req := arg.GetDo()
		o, ok := wobjs[req.GetUserID()]
		if !ok {
			logger.Debugf("sceneO:%s is nil", req.GetUserID())
			replyer.Reply(&rpc.WorldObjPushResp{
				Ok: proto.Bool(false),
			})
			return
		}

		scenePos := scene.FixPosition(scene.Position{
			X: req.GetUpdatePos().GetX(),
			Y: req.GetUpdatePos().GetY(),
			Z: req.GetUpdatePos().GetZ(),
		})
		o.ScenePos = scenePos

		replyer.Reply(&rpc.WorldObjPushResp{
			Ok: proto.Bool(true),
		})

		for _, s := range scene.GetScenes() {
			scene_ := s
			sceneO := scene_.GetObjByID(req.GetUserID())
			if scene_.Move(sceneO) {
				sceneO.(*WorldObj).ScenePos = scenePos
				sceneO.(*WorldObj).isMoved = true
			}
		}
		//logger.Debugf("%s DoAction ok", req.GetUserID())
	}
}

func init() {
	worldObjPush.Register(&WorldObjPush{})
}
