package PlayerGene

func GetGene(characterID int32, geneLevel int32) *PlayerGene {
	id := characterID*1000 + geneLevel
	return GetID(id)
}

//func GeneToMax(characterID int32, currentGeneLevel int32) int32 {
//
//}
