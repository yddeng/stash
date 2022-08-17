package ScarsIngrainBossInstance

func (this *ScarsIngrainBossInstance) BossDifficult() (int32, int32) {
	difficult := this.ID % 10
	bossId := this.ID / 10
	return bossId, difficult
}
