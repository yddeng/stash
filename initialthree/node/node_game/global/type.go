package global

import (
	"initialthree/node/table/excel/DataTable/ScarsIngrainArea"
	"initialthree/node/table/excel/DataTable/ScarsIngrainRankReward"
)

const (
	ScarsIngrain = 1 // 战痕印刻
	BigSecret    = 2 // 大秘境
)

// 返回类型
func ParseRankID(rankID int32) int32 {
	return rankID % 10
}

func RankPercent(rank, total int32) float64 {
	if total == 0 {
		return 0
	}
	if total < 100 {
		return float64(rank) / 100
	} else {
		return float64(rank) / float64(total)
	}
}

func RankAwardPoolID(rankID, rank, total int32) int32 {
	tt := ParseRankID(rankID)
	var poolID int32
	switch tt {
	case ScarsIngrain:
		// 第二位表示区段
		area := (rankID / 10) % 10
		siArea := ScarsIngrainArea.GetID(area)
		rankAward := ScarsIngrainRankReward.GetID(siArea.RankReward)
		perc := RankPercent(rank, total)
		poolID = rankAward.GetRankPercentPoolID(perc)
	case BigSecret:
	default:

	}

	return poolID
}
