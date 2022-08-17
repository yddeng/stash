package world

import (
	"github.com/golang/protobuf/proto"
	"initialthree/cluster"
	"initialthree/cluster/addr"
	"initialthree/protocol/cmdEnum"
	"initialthree/protocol/ss/ssmessage"
	"time"
)

func mapToWorld(from addr.LogicAddr, msg proto.Message) {
	arg := msg.(*ssmessage.MapToWorld)

	now := time.Now()
	nodeMap := getNodeMap(from)
	if nodeMap == nil {
		nodeMap = &NodeMap{
			LogicAddr: from,
			status:    true,
			SyncTime:  now,
			Scenes:    make([]*Scene, 0, len(arg.GetScenes())),
		}
		for _, s := range arg.GetScenes() {
			nodeMap.Scenes = append(nodeMap.Scenes, &Scene{
				SceneID:  s.GetSceneID(),
				AOITotal: s.GetAOITotal(),
			})
		}
		myWorld.NodeMaps[from] = nodeMap

		mapLogic, err := addr.MakeLogicAddr(arg.GetMapAddr())
		logger.Infoln("register nodeMap:", arg.GetMapAddr(), from, mapLogic, err)
	} else {
		nodeMap.status = true
		nodeMap.SyncTime = now
		for i, s := range nodeMap.Scenes {
			s.AOICurr = arg.GetScenes()[i].GetAOICurr()
		}
	}

}

func init() {
	cluster.Register(cmdEnum.SS_MapToWorld, mapToWorld)
}
