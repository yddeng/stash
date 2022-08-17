package obj

import (
	"initialthree/node/node_world/world"
	"time"
)

type Boss struct {
	UserID string
	ID     uint64
	Pos    *world.Position
	Status int32

	Attack   float64
	Defend   float64
	HitPoint float64
}

func (this *Boss) GetUserID() string {
	return this.UserID
}

func (this *Boss) GetID() uint64 {
	return this.ID
}

func (this *Boss) CompareAndSwapStatus(old, new int32) bool {
	if this.Status == old {
		this.Status = new
		return true
	} else {
		return false
	}
}

func (this *Boss) SetPosition(pos *world.Position) {
	this.Pos = pos
}

func (this *Boss) GetPosition() *world.Position {
	return this.Pos
}

// 玩家攻击世界boss 测试
func (this *Boss) DoSomething(req interface{}) bool {
	arg := req.(world.HitBoss)

	this.HitPoint -= arg.Attack - this.Defend
	if this.HitPoint <= 0 {
		this.Status = 0
		world.RemoveWorldObj(this.UserID)
	}
	return true
}

func (this *Boss) StartAi() {
	go func() {
		timer := time.NewTicker(1000 * time.Millisecond)
		_var := int32(50)
		for {
			<-timer.C
			this.Pos.X += _var
			if this.Pos.X >= 2000 || this.Pos.X < 0 {
				_var = -_var
			}
			world.MoveAllMap(this, nil)
		}
	}()
}
