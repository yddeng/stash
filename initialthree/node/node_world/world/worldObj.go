package world

import (
	"github.com/golang/protobuf/proto"
	"initialthree/cluster"
	"initialthree/protocol/ss/rpc"
	"initialthree/rpc/worldObjPush"
)

type WorldObj interface {
	GetUserID() string
	GetID() uint64
	CompareAndSwapStatus(old, new int32) bool
	SetPosition(pos *Position)
	GetPosition() *Position
	DoSomething(req interface{}) bool
}

type Position struct {
	X     int32
	Y     int32
	Z     int32
	Angle int32
}

type HitBoss struct {
	Attack float64
}

func AddWorldObj(obj WorldObj) bool {
	if _, ok := myWorld.WorldObjs[obj.GetUserID()]; !ok {
		myWorld.WorldObjs[obj.GetUserID()] = obj
		return true
	}
	return false
}

func RemoveWorldObj(id string) {
	delete(myWorld.WorldObjs, id)
}

func EnterAllMap(obj WorldObj, cb func(err error)) {
	pos := obj.GetPosition()
	enter := &rpc.WorldObjPushReq{
		T: rpc.PushType_EnterMap.Enum(),
		Enter: &rpc.WOEnterMap{
			UserID: proto.String(obj.GetUserID()),
			ID:     proto.Uint64(obj.GetID()),
			Pos: &rpc.Position{
				X: proto.Int32(pos.X),
				Y: proto.Int32(pos.Y),
				Z: proto.Int32(pos.Z),
			},
			WorldAddr: proto.String(cluster.SelfAddr().Logic.String()),
		},
	}

	for addr_ := range myWorld.NodeMaps {
		// 调用enterMap
		worldObjPush.AsynCall(addr_, enter, time.Second*5, func(resp2 *rpc.WorldObjPushResp, e error) {
			if e != nil {
				logger.Errorln(e)
			} else {
				if !resp2.GetOk() {
					logger.Errorf("world obj:%s enterMap:%s failed\n", obj.GetUserID(), addr_.String())
				}
			}
		})
	}

}

func LeaveAllMap(obj WorldObj, cb func(err error)) {

	leave := &rpc.WorldObjPushReq{
		T: rpc.PushType_LeaveMap.Enum(),
		Leave: &rpc.WOLeaveMap{
			UserID: proto.String(obj.GetUserID()),
		},
	}

	for addr_ := range myWorld.NodeMaps {
		// 调用enterMap
		worldObjPush.AsynCall(addr_, leave, time.Second*5, func(resp2 *rpc.WorldObjPushResp, e error) {
			if e != nil {
				logger.Errorln(e)
			} else {
				if !resp2.GetOk() {
					logger.Infof("world obj:%s leaveMap:%s failed\n", obj.GetUserID(), addr_.String())
				}
			}
		})
	}
}

func MoveAllMap(obj WorldObj, cb func(err error)) {

	pos := obj.GetPosition()
	move := &rpc.WorldObjPushReq{
		T: rpc.PushType_DoAction.Enum(),
		Do: &rpc.WODoAction{
			UserID: proto.String(obj.GetUserID()),
			UpdatePos: &rpc.Position{
				X: proto.Int32(pos.X),
				Y: proto.Int32(pos.Y),
				Z: proto.Int32(pos.Z),
			},
		},
	}

	for addr_ := range myWorld.NodeMaps {
		// 调用enterMap
		worldObjPush.AsynCall(addr_, move, time.Second*5, func(resp2 *rpc.WorldObjPushResp, e error) {
			if e != nil {
				logger.Errorln(e)
			} else {
				if !resp2.GetOk() {
					logger.Infof("world obj:%s move:%s failed\n", obj.GetUserID(), addr_.String())
				}
			}
		})
	}

}
