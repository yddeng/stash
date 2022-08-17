package world

import (
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	"initialthree/cluster"
	"initialthree/cluster/addr"
	"initialthree/pkg/golog"
	"initialthree/pkg/timer"
	"time"
)

type World struct {
	WorldID   uint32 // 地图地址
	XLength   int32
	YLength   int32
	NodeMaps  map[addr.LogicAddr]*NodeMap
	WorldObjs map[string]WorldObj
}

var (
	myWorld *World
	logger  golog.LoggerI
	timeout = time.Second * 5 // 心跳超时时间
)

// 获取一个有容量的map节点
func randomNodeMap() (*NodeMap, *Scene, error) {
	for _, m := range myWorld.NodeMaps {
		if m.status {
			scene_ := m.getFreeScene()
			if scene_ != nil {
				return m, scene_, nil
			}
		}
	}

	return nil, nil, errors.New("no free nodeMap or Scene")
}

func getNodeMap(_addr addr.LogicAddr) *NodeMap {
	return myWorld.NodeMaps[_addr]
}

// 广播给所有Map
func BroadcastToMap(msg proto.Message) {
	for _addr, m := range myWorld.NodeMaps {
		fmt.Println(m)
		cluster.PostMessage(_addr, msg)
	}
}

func InitLogger(baseLogger golog.LoggerI) {
	logger = baseLogger
}

func InitWorld(worldID uint32) {
	myWorld = &World{
		WorldID:   worldID,
		NodeMaps:  map[addr.LogicAddr]*NodeMap{},
		WorldObjs: map[string]WorldObj{},
	}

	cluster.RegisterTimer(time.Second, func(t *timer.Timer, _ interface{}) {
		now := time.Now()
		for _, m := range myWorld.NodeMaps {
			if m.status && now.After(m.SyncTime.Add(timeout)) {
				m.status = false
			}
		}
	}, nil)
}

func Tick(now time.Time) {

}
