package Skill

func (this *Skill) GetMaxLevel() int32 {
	return int32(len(this.Damage))
}
