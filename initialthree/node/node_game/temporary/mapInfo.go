package temporary

import (
	"github.com/golang/protobuf/proto"
	"initialthree/cluster"
	"initialthree/cluster/addr"
	"initialthree/zaplogger"

	"initialthree/protocol/ss/rpc"
	"initialthree/protocol/ss/ssmessage"
	"initialthree/rpc/leaveMap"
	"time"
)

const heartbeatDur = time.Second * 5

type MapInfo struct {
	UserI         UserI
	MapID         int32
	MapAddr       addr.LogicAddr
	SceneIdx      int32
	nextHeartbeat time.Time
}

func NewMapInfo(userI UserI, logicAddr addr.LogicAddr, mapID, sceneIdx int32) *MapInfo {
	return &MapInfo{
		UserI:         userI,
		MapID:         mapID,
		MapAddr:       logicAddr,
		SceneIdx:      sceneIdx,
		nextHeartbeat: time.Now().Add(heartbeatDur),
	}
}

func (this *MapInfo) LeaveMap(cb func(ok bool)) {
	leaveMap.AsynCall(this.MapAddr, &rpc.LeaveMapReq{
		SceneIdx: proto.Int32(this.SceneIdx),
		UserID:   proto.String(this.UserI.GetUserID()),
		RoleID:   proto.Uint64(this.UserI.GetID()),
	}, time.Second*5, func(resp *rpc.LeaveMapResp, e error) {
		if e == nil && resp.GetOk() {
			zaplogger.GetSugar().Debugf("leavemap ok %v %v", this.UserI.GetUserID(), this.UserI.GetID())
			this.UserI.ClearTemporary(TempMapInfo)
			if cb != nil {
				cb(true)
			}
		} else {
			zaplogger.GetSugar().Debugf("%v %v leavemap fail:%v %v", this.UserI.GetUserID(), this.UserI.GetID(), e, resp.GetOk())
			if cb != nil {
				cb(false)
			}
		}
	})
}

func (this *MapInfo) UserDisconnect() {
	this.LeaveMap(nil)
}

func (this *MapInfo) UserLogout() {
	this.LeaveMap(nil)
}

func (this *MapInfo) Tick(now time.Time) {
	if now.After(this.nextHeartbeat) {
		cluster.PostMessage(this.MapAddr, &ssmessage.MapHeartbeat{
			MapID:    proto.Int32(this.MapID),
			SceneIdx: proto.Int32(this.SceneIdx),
			UserID:   proto.String(this.UserI.GetUserID()),
			RoleID:   proto.Uint64(this.UserI.GetID()),
		})
		this.nextHeartbeat = now.Add(heartbeatDur)
	}
}
