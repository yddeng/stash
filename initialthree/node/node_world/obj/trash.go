package obj

import "initialthree/node/node_world/world"

// 垃圾桶
type TrashCan struct {
	UserID string
	ID     uint64
	Pos    *world.Position
	Status int32
}

func (this *TrashCan) GetUserID() string {
	return this.UserID
}

func (this *TrashCan) GetID() uint64 {
	return this.ID
}

func (this *TrashCan) CompareAndSwapStatus(old, new int32) bool {
	if this.Status == old {
		this.Status = new
		return true
	} else {
		return false
	}
}

func (this *TrashCan) SetPosition(pos *world.Position) {
	this.Pos = pos
}

func (this *TrashCan) GetPosition() *world.Position {
	return this.Pos
}

func (this *TrashCan) DoSomething(req interface{}) bool {
	return true
}
