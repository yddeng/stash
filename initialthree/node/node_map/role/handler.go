package role

import (
	"github.com/golang/protobuf/proto"
	"initialthree/cluster"
	"initialthree/cluster/addr"
	"initialthree/node/node_map/scene"
	"initialthree/protocol/cmdEnum"
	ss_rpc "initialthree/protocol/ss/rpc"
	"initialthree/protocol/ss/ssmessage"
	"initialthree/rpc/enterMap"
	"initialthree/rpc/leaveMap"
	"initialthree/rpc/move"
	"time"
)

type Move struct {
}

func (this *Move) OnCall(replyer *move.MoveReplyer, arg *ss_rpc.MoveReq) {

	scene_ := scene.GetScene(arg.GetSceneIdx())
	if nil == scene_ {
		replyer.Reply(&ss_rpc.MoveResp{
			Ok: proto.Bool(false),
		})
		return
	}

	scene_.PostTask(func() {

		sceneO := scene_.GetObjByID(arg.GetUserID())
		if nil == sceneO {
			replyer.Reply(&ss_rpc.MoveResp{
				Ok: proto.Bool(false),
			})
			return
		}

		scenePos := scene.FixPosition(scene.Position{
			X:     arg.GetPos().GetX(),
			Y:     arg.GetPos().GetY(),
			Z:     arg.GetPos().GetZ(),
			Angle: arg.GetAngle(),
		})

		if scene_.Move(sceneO) {
			sceneO.(*Role).ScenePos = scenePos
			sceneO.(*Role).isMoved = true

			replyer.Reply(&ss_rpc.MoveResp{
				Ok: proto.Bool(true),
				Pos: &ss_rpc.Position{
					X: proto.Int32(scenePos.X),
					Y: proto.Int32(scenePos.Y),
					Z: proto.Int32(scenePos.Z),
				},
				Angle: proto.Int32(scenePos.Angle),
			})
		} else {
			replyer.Reply(&ss_rpc.MoveResp{
				Ok: proto.Bool(false),
			})
		}
	})
}

type EnterMap struct {
}

func (this *EnterMap) OnCall(replyer *enterMap.EnterMapReplyer, arg *ss_rpc.EnterMapReq) {

	logger.Debugln("OnCall EnterMap", arg)

	scene_ := scene.GetScene(arg.GetSceneIdx())

	if nil == scene_ {
		logger.Debugln("OnCall EnterMap SceneIdx:", arg.GetSceneIdx(), "not found")
		replyer.Reply(&ss_rpc.EnterMapResp{
			Ok: proto.Bool(false),
		})
		return
	}

	scene_.PostTask(func() {
		sceneO := scene_.GetObjByID(arg.GetUserID())
		if sceneO != nil {
			logger.Debugf("sceneO:%s is not nil", arg.GetUserID())
			replyer.Reply(&ss_rpc.EnterMapResp{
				Ok: proto.Bool(false),
			})
			return
		}

		r := &Role{
			enter:     map[scene.SceneObject]struct{}{},
			leave:     map[scene.SceneObject]struct{}{},
			heartbeat: time.Now(),
		}
		r.ID = arg.GetID()
		r.UserID = arg.GetUserID()
		r.Game, _ = addr.MakeLogicAddr(arg.GetGameAddr())

		gateAddr, _ := addr.MakeLogicAddr(arg.GetGateAddr())
		r.Gate = &gate{
			GateAddr: gateAddr,
			GateUid:  arg.GetGateUid(),
		}

		r.ScenePos = scene.Position{
			X: arg.GetPos().GetX(),
			Y: arg.GetPos().GetY(),
			Z: arg.GetPos().GetZ(),
		}

		ok := scene_.EnterScene(r)
		replyer.Reply(&ss_rpc.EnterMapResp{
			Ok: proto.Bool(ok),
		})

		logger.Debugf("%s EnterMap %v", arg.GetUserID(), ok)

	})
}

type LeaveMap struct {
}

func (this *LeaveMap) OnCall(replyer *leaveMap.LeaveMapReplyer, arg *ss_rpc.LeaveMapReq) {

	logger.Debugln("OnCall LeaveMap", arg)

	scene_ := scene.GetScene(arg.GetSceneIdx())

	if nil == scene_ {
		replyer.Reply(&ss_rpc.LeaveMapResp{
			Ok: proto.Bool(true),
		})
		return
	}

	scene_.PostTask(func() {
		logger.Debugln("do leave map")
		sceneO := scene_.GetObjByID(arg.GetUserID())
		if nil == sceneO {
			replyer.Reply(&ss_rpc.LeaveMapResp{
				Ok: proto.Bool(true),
			})
			return
		}

		ok := scene_.LeaveScene(sceneO)
		replyer.Reply(&ss_rpc.LeaveMapResp{
			Ok: proto.Bool(ok),
		})

		logger.Debugln(sceneO.(*Role).GetUserID(), "do leave map", ok)
	})
}

func onStartAoi(form addr.LogicAddr, msg proto.Message) {
	arg := msg.(*ssmessage.StartAoi)
	scene_ := scene.GetScene(arg.GetSceneIdx())

	if scene_ == nil {
		return
	}
	scene_.PostTask(func() {
		o := scene_.GetObjByID(arg.GetUserID())

		if o.GetAoiEntity() == nil {
			scene_.StartAoi(o)
		}
	})
}

func onHeartbeat(form addr.LogicAddr, msg proto.Message) {
	arg := msg.(*ssmessage.MapHeartbeat)
	scene_ := scene.GetScene(arg.GetSceneIdx())
	if scene_ == nil {
		return
	}
	scene_.PostTask(func() {
		o := scene_.GetObjByID(arg.GetUserID()).(*Role)
		o.heartbeat = time.Now()
	})
}

func init() {
	enterMap.Register(&EnterMap{})
	leaveMap.Register(&LeaveMap{})
	move.Register(&Move{})
	cluster.Register(cmdEnum.SS_StartAoi, onStartAoi)
	cluster.Register(cmdEnum.SS_MapHeartbeat, onHeartbeat)
}
