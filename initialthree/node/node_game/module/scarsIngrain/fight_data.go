package scarsIngrain

import "initialthree/protocol/cs/message"

// 手动战斗的零时数据
// 需要玩家确认战斗结束，才扣除消耗
// 即使一场打完后，没有这里确认，仍然可以进行正常战斗，战斗开始清理
type FightData struct {
	ScoreID int32
	Score   int32
	BossDie bool // boss是否死亡，用于解锁下一难度
	Tos     *message.LevelFightToS
}

func (this *ScarsIngrain) AddFightData(score int32, bossDie bool, tos *message.LevelFightToS) int32 {
	this.scoreId++

	v := &FightData{
		ScoreID: this.scoreId,
		Score:   score,
		Tos:     tos,
		BossDie: bossDie,
	}
	this.fightData[v.ScoreID] = v
	return v.ScoreID
}

func (this *ScarsIngrain) GetFightData(scoreID int32) (*FightData, bool) {
	v, ok := this.fightData[scoreID]
	return v, ok
}

func (this *ScarsIngrain) ClearFightData() {
	this.fightData = map[int32]*FightData{}
}
