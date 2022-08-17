package role

import (
	"initialthree/cluster/addr"
	"initialthree/node/common/aoi"
	"initialthree/node/node_map/scene"
	cs_message "initialthree/protocol/cs/message"
	"time"

	"github.com/golang/protobuf/proto"
	"initialthree/pkg/golog"
)

var logger golog.LoggerI

func InitLogger(baseLogger golog.LoggerI) {
	logger = baseLogger
}

type RoleData struct {
	ID   uint64 //角色ID
	Name string //角色名
}

type gate struct {
	GateAddr addr.LogicAddr //gate地址
	GateUid  uint64         //gate用户id
}

type Role struct {
	RoleData
	entity aoi.Entity
	UserID string
	Scene  *scene.Scene
	Game   addr.LogicAddr
	Gate   *gate

	ScenePos  scene.Position
	isMoved   bool
	enter     map[scene.SceneObject]struct{}
	leave     map[scene.SceneObject]struct{}
	heartbeat time.Time
}

func (this *Role) Heartbeat() time.Time {
	return this.heartbeat
}

func (this *Role) PackRoleData() *cs_message.ViewObj {
	return &cs_message.ViewObj{
		UserID: proto.String(this.UserID),
		RoleID: proto.Uint64(this.ID),
		Pos: &cs_message.Position{
			X: proto.Int32(this.ScenePos.X),
			Y: proto.Int32(this.ScenePos.Y),
			Z: proto.Int32(this.ScenePos.Z),
		},
		Angle: proto.Int32(this.ScenePos.Angle),
	}
}

func (this *Role) GetUserID() string {
	return this.UserID
}

func (this *Role) GetID() uint64 {
	return this.ID
}

func (this *Role) SetScene(scene_ *scene.Scene) {
	this.Scene = scene_
}

func (this *Role) GetScene() *scene.Scene {
	return this.Scene
}

func (this *Role) GetAoiPosition() aoi.Position {
	pos := aoi.Position{
		X: this.ScenePos.X,
		Y: this.ScenePos.Y,
	}
	return pos
}

func (this *Role) GetAoiEntity() aoi.Entity {
	return this.entity
}

func (this *Role) SetAoiEntity(en aoi.Entity) {
	this.entity = en
}

func (this *Role) OnAOIUpdate(enter, level []aoi.User) {

	if enter != nil {
		for _, u := range enter {
			r := u.(scene.SceneObject)
			if _, ok := this.leave[r]; ok {
				delete(this.leave, r)
			} else {
				this.enter[r] = struct{}{}
			}
		}
	}

	if level != nil {
		for _, u := range level {
			r := u.(scene.SceneObject)
			if _, ok := this.enter[r]; ok {
				delete(this.enter, r)
			} else {
				this.leave[r] = struct{}{}
			}
		}
	}
}

func (this *Role) ClearMoved() {
	this.isMoved = false
}

func (this *Role) Tick(now time.Time) {
	if this.entity != nil {
		this.SendEnterSee()
		this.SendLeaveSee()
		this.PackPosUpdate()
	}
}

func (this *Role) SendEnterSee() {

	if len(this.enter) > 0 {
		EnterSeeToC := &cs_message.EnterSeeToC{
			Objs: make([]*cs_message.ViewObj, 0, len(this.enter)),
		}

		for o := range this.enter {
			if o.GetUserID() != this.UserID {
				EnterSeeToC.Objs = append(EnterSeeToC.Objs, o.PackRoleData())
				if len(EnterSeeToC.Objs) == 128 {
					this.SendToClient(EnterSeeToC)
					EnterSeeToC.Objs = []*cs_message.ViewObj{}
				}
			}
		}

		if len(EnterSeeToC.Objs) > 0 {
			this.SendToClient(EnterSeeToC)
		}

		this.enter = map[scene.SceneObject]struct{}{}
	}
}

func (this *Role) SendLeaveSee() {
	if len(this.leave) > 0 {
		LeaveSeeToC := &cs_message.LeaveSeeToC{
			RoleID: make([]uint64, 0, len(this.leave)),
		}

		for o := range this.leave {
			if o.GetUserID() != this.UserID {
				LeaveSeeToC.RoleID = append(LeaveSeeToC.RoleID, o.GetID())
				if len(LeaveSeeToC.RoleID) == 256 {
					this.SendToClient(LeaveSeeToC)
					LeaveSeeToC.RoleID = LeaveSeeToC.RoleID[0:0]
				}
			}
		}

		if len(LeaveSeeToC.RoleID) > 0 {
			this.SendToClient(LeaveSeeToC)
		}

		this.leave = map[scene.SceneObject]struct{}{}
	}
}

// 将自己能看到的玩家打包给自己（不包括自己)
func (this *Role) PackPosUpdate() {

	msg := &cs_message.UpdatePosToC{
		Objs: make([]*cs_message.UpdateObj, 0, 64),
	}

	// 能看见我的，我也能看见他
	this.entity.TraverseAOI(func(u aoi.User) error {
		r, ok := u.(*Role)
		if ok && r.isMoved {
			msg.Objs = append(msg.Objs, &cs_message.UpdateObj{
				RoleID: proto.Uint64(r.ID),
				Pos: &cs_message.Position{
					X: proto.Int32(r.ScenePos.X),
					Y: proto.Int32(r.ScenePos.Y),
					Z: proto.Int32(r.ScenePos.Z),
				},
				Angle: proto.Int32(r.ScenePos.Angle),
			})
			if len(msg.Objs) == 128 {
				this.SendToClient(msg)
				msg.Objs = msg.Objs[0:0]
			}
		}
		return nil
	})

	if len(msg.Objs) > 0 {
		this.SendToClient(msg)
	}
}
