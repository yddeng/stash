package ScarsIngrainScoreReward

func (this *ScarsIngrainScoreReward) GetScoreAward(score int32) (int32, bool) {
	for _, v := range this.ScoreReward {
		if v.Score == score {
			return v.DropPoolID, true
		}
	}
	return 0, false
}
