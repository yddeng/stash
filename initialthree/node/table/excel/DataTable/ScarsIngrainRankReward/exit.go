package ScarsIngrainRankReward

func (this *ScarsIngrainRankReward) GetRankPercentPoolID(percent float64) int32 {
	if percent == 0 {
		return 0
	}
	for _, v := range this.RankReward {
		if percent > v.Start && percent <= v.End {
			return v.DropPoolID
		}
	}
	return 0
}
