package world

import (
	"initialthree/cluster/addr"
	"time"
)

type NodeMap struct {
	LogicAddr addr.LogicAddr
	Scenes    []*Scene

	status   bool
	SyncTime time.Time
}

type Scene struct {
	SceneID  int32
	AOITotal int32 // 人数上线
	AOICurr  int32 // 同步过来的当前人数
	AOIReq   int32 // 请求的人数
}

func (this *NodeMap) getFreeScene() *Scene {
	for _, s := range this.Scenes {
		if s.AOIReq+s.AOICurr < s.AOITotal {
			return s
		}
	}
	return nil
}

func (this *NodeMap) getScene(idx int) *Scene {
	if len(this.Scenes) < idx {
		return nil
	}
	return this.Scenes[idx]
}

func (scene *Scene) enterScene() bool {
	total := scene.AOICurr + scene.AOIReq
	if total+1 <= scene.AOITotal {
		return true
	}
	scene.AOIReq += 1
	return false
}
