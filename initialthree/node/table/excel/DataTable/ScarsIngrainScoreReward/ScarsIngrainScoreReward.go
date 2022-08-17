package ScarsIngrainScoreReward

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	path2 "path"
	"strings"
	"sync/atomic"
	"initialthree/node/table/excel"
	"initialthree/zaplogger"
)

type ScarsIngrainScoreReward struct{
    ID int32 
    Score_1 int32 
    DropPool_1 int32 
    Score_2 int32 
    DropPool_2 int32 
    Score_3 int32 
    DropPool_3 int32 
    Score_4 int32 
    DropPool_4 int32 
    Score_5 int32 
    DropPool_5 int32 
    ScoreReward []*ScoreReward_ 

}

type ScoreReward_ struct{
    Score int32 
    DropPoolID int32 

}

func readScoreReward(row, names []string)[]*ScoreReward_{
	ret := make([]*ScoreReward_, 0)
	base := excel.Split("Score_,DropPool_", ",")

	for i := 1; ; i++ {
		if idx, find := hasCellName(names, fmt.Sprintf("%s%d", base[0][0], i));!find {
			break
		}else if _,ok := hasCellValue(row,idx);!ok{
			break
		}
		
		e := &ScoreReward_{}
        e.Score = excel.ReadInt32(readCellValue(row, names, fmt.Sprintf("%s%d", base[0][0], i)))
        e.DropPoolID = excel.ReadInt32(readCellValue(row, names, fmt.Sprintf("%s%d", base[1][0], i)))

		ret = append(ret, e)
	}

	return ret
}

var Table_ Table

type Table struct{
	xlsxName string
	nextPath string
	indexID  atomic.Value

}

func GetID(key int32)*ScarsIngrainScoreReward{
	v := Table_.indexID.Load().(map[int32]*ScarsIngrainScoreReward)[key]
	if v == nil{
		zaplogger.GetSugar().Errorf("Table %s GetID %d is nil", Table_.xlsxName, key)
	}
	return v
}

func GetIDMap() map[int32]*ScarsIngrainScoreReward {
	return Table_.indexID.Load().(map[int32]*ScarsIngrainScoreReward)
}


func hasCellName(names []string, name string) (int, bool) {
	for i, n := range names {
		if n == name {
			return i, true
		}
	}
	return 0, false
}

func hasCellValue(row []string, idx int) (string, bool) {
	val := strings.TrimSpace(row[idx])
	return val, val != ""
}

func readCellValue(row, names []string, name string) string {
	for i, n := range names {
		if n == name {
			return strings.TrimSpace(row[i])
		}
	}

	panic("cell " + name + " not found")

	return ""
}

func init(){
	Table_.xlsxName = "ScarsIngrainScoreReward.xlsx"
	excel.AddTable(&Table_)
}

func (this *Table)Load(path string){
	xlFile, err := excelize.OpenFile(path2.Join(path2.Join(path, "/DataTable/关卡/战痕"), this.xlsxName))
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	tableRows := xlFile.GetRows(xlFile.GetSheetName(xlFile.GetActiveSheetIndex()))
	names := tableRows[1]
	rows := tableRows[4:]

	tmpIDMap := map[int32]*ScarsIngrainScoreReward{}

	for _,row := range rows{
		if row[0] != "" {
			e := &ScarsIngrainScoreReward{}
            e.ID = excel.ReadInt32(readCellValue(row, names, "ID"))
            e.Score_1 = excel.ReadInt32(readCellValue(row, names, "Score_1"))
            e.DropPool_1 = excel.ReadInt32(readCellValue(row, names, "DropPool_1"))
            e.Score_2 = excel.ReadInt32(readCellValue(row, names, "Score_2"))
            e.DropPool_2 = excel.ReadInt32(readCellValue(row, names, "DropPool_2"))
            e.Score_3 = excel.ReadInt32(readCellValue(row, names, "Score_3"))
            e.DropPool_3 = excel.ReadInt32(readCellValue(row, names, "DropPool_3"))
            e.Score_4 = excel.ReadInt32(readCellValue(row, names, "Score_4"))
            e.DropPool_4 = excel.ReadInt32(readCellValue(row, names, "DropPool_4"))
            e.Score_5 = excel.ReadInt32(readCellValue(row, names, "Score_5"))
            e.DropPool_5 = excel.ReadInt32(readCellValue(row, names, "DropPool_5"))
            e.ScoreReward = readScoreReward(row, names)

			tmpIDMap[excel.ReadInt32(readCellValue(row, names, "ID"))] = e

		}
	}
	this.indexID.Store(tmpIDMap)	

}
