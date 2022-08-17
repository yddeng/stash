package Rank

import (
	"initialthree/cluster"
	"initialthree/cluster/addr"
	"initialthree/node/common/serverType"
	"initialthree/node/node_game/global"
	"initialthree/node/node_game/global/bigSecret"
	"initialthree/node/node_game/global/scarsIngrain"
)

func RankLogicAddr(rankID int32) (logicAddr addr.LogicAddr, err error) {
	logic := ""

	tt := global.ParseRankID(rankID)
	switch tt {
	case global.ScarsIngrain:
		// 第二位表示区段
		area := (rankID / 10) % 10
		siData := scarsIngrain.GetIDData(area)
		if siData != nil && siData.RankID == rankID {
			logic = siData.RankLogic
		}
	case global.BigSecret:
		data := bigSecret.GetData()
		if data != nil {
			logic = data.RankLogic
		}
	default:

	}

	if logic != "" {
		logicAddr, err = addr.MakeLogicAddr(logic)
		if err != nil {
			logicAddr, err = cluster.Random(serverType.Rank)
		}
	} else {
		logicAddr, err = cluster.Random(serverType.Rank)
	}
	return
}
