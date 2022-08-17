package role

import (
	"github.com/golang/protobuf/proto"
	"initialthree/cluster/addr"
	"initialthree/node/common/aoi"
	"initialthree/node/node_map/scene"
	"initialthree/protocol/cs/message"
	"time"
)

type WorldObj struct {
	UserID string
	ID     uint64

	entity aoi.Entity
	Scene  *scene.Scene
	World  addr.LogicAddr //world地址

	ScenePos  scene.Position
	isMoved   bool
	enter     map[scene.SceneObject]struct{} // 视野对象，这里只存用户角色
	leave     map[scene.SceneObject]struct{}
	heartbeat time.Time
}

func (this *WorldObj) Heartbeat() time.Time {
	return this.heartbeat
}

func (this *WorldObj) PackRoleData() *message.ViewObj {
	return &message.ViewObj{
		UserID: proto.String(this.UserID),
		RoleID: proto.Uint64(this.ID),
		Pos: &message.Position{
			X: proto.Int32(this.ScenePos.X),
			Y: proto.Int32(this.ScenePos.Y),
			Z: proto.Int32(this.ScenePos.Z),
		},
		Angle: proto.Int32(this.ScenePos.Angle),
	}
}

func (this *WorldObj) GetUserID() string {
	return this.UserID
}

func (this *WorldObj) GetID() uint64 {
	return this.ID
}

func (this *WorldObj) SetScene(scene_ *scene.Scene) {
	this.Scene = scene_
}

func (this *WorldObj) GetScene() *scene.Scene {
	return this.Scene
}

func (this *WorldObj) GetAoiPosition() aoi.Position {
	pos := aoi.Position{
		X: this.ScenePos.X,
		Y: this.ScenePos.Y,
	}
	return pos
}

func (this *WorldObj) GetAoiEntity() aoi.Entity {
	return this.entity
}

func (this *WorldObj) SetAoiEntity(en aoi.Entity) {
	this.entity = en
}

func (this *WorldObj) OnAOIUpdate(enter, level []aoi.User) {

	if enter != nil {
		for _, u := range enter {
			r, ok := u.(*Role)
			if ok {
				this.enter[r] = struct{}{}
			}
		}
	}

	if level != nil {
		for _, u := range level {
			r, ok := u.(*Role)
			if ok {
				this.leave[r] = struct{}{}
			}
		}
	}
}

func (this *WorldObj) ClearMoved() {
	this.isMoved = false
}

func (this *WorldObj) Tick(now time.Time) {
	//this.NotifyPosUpdate()
}

func (this *WorldObj) NotifyPosUpdate() {

}
