package ScarsIngrainBossInstance

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	path2 "path"
	"strings"
	"sync/atomic"
	"initialthree/node/table/excel"
	"initialthree/zaplogger"
)

type ScarsIngrainBossInstance struct{
    ID int32 
    AllScore int32 
    DamageScorePercent float64 
    TimeScorePercent float64 
    DungeonID int32 
    BossID int32 

}

var Table_ Table

type Table struct{
	xlsxName string
	nextPath string
	indexID  atomic.Value

}

func GetID(key int32)*ScarsIngrainBossInstance{
	v := Table_.indexID.Load().(map[int32]*ScarsIngrainBossInstance)[key]
	if v == nil{
		zaplogger.GetSugar().Errorf("Table %s GetID %d is nil", Table_.xlsxName, key)
	}
	return v
}

func GetIDMap() map[int32]*ScarsIngrainBossInstance {
	return Table_.indexID.Load().(map[int32]*ScarsIngrainBossInstance)
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
	Table_.xlsxName = "ScarsIngrainBossInstance.xlsx"
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

	tmpIDMap := map[int32]*ScarsIngrainBossInstance{}

	for _,row := range rows{
		if row[0] != "" {
			e := &ScarsIngrainBossInstance{}
            e.ID = excel.ReadInt32(readCellValue(row, names, "ID"))
            e.AllScore = excel.ReadInt32(readCellValue(row, names, "AllScore"))
            e.DamageScorePercent = excel.ReadFloat(readCellValue(row, names, "DamageScorePercent"))
            e.TimeScorePercent = excel.ReadFloat(readCellValue(row, names, "TimeScorePercent"))
            e.DungeonID = excel.ReadInt32(readCellValue(row, names, "DungeonID"))
            e.BossID = excel.ReadInt32(readCellValue(row, names, "BossID"))

			tmpIDMap[excel.ReadInt32(readCellValue(row, names, "ID"))] = e

		}
	}
	this.indexID.Store(tmpIDMap)	

}
