package ScarsIngrainBossChallenge

func (this *ScarsIngrainBossChallenge) GetScoreBuff(arg int32) (int32, int32) {
	length := len(this.Challenge)
	for idx, v := range this.Challenge {
		if idx == 0 {
			if arg <= v.ArgRight {
				return v.Score, v.Buff
			}
		}
		if idx == length-1 {
			return v.Score, v.Buff
		}

		if arg > v.ArgLeft && arg <= v.ArgRight {
			return v.Score, v.Buff
		}
	}

	return 0, 0
}
