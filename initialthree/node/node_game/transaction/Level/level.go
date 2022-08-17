package Level

import (
	"initialthree/codec/cs"
	_ "initialthree/node/node_game/transaction/Level/BigSecret"
	_ "initialthree/node/node_game/transaction/Level/ScarsIngrain"
	_ "initialthree/node/node_game/transaction/Level/SecertDungeon"
	_ "initialthree/node/node_game/transaction/Level/TrialDungeon"
	"initialthree/node/node_game/user"
)

/*
type ScarsIngrainLevelArg struct {
	StartTime   int64     // 开始时间
	RoleMaxHP   []float64 // 角色队伍对应血量上限
	BossTotalHP float64   // 怪物总血量
}

*/

func CheckModuleOpen(user2 *user.User, msg *cs.Message) bool {
	level := user2.GetLevel()
	return level > 0
}

func init() {
	user.RegisterCheckModuleOpen("Level", CheckModuleOpen)
}
