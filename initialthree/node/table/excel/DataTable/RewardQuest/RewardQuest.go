package RewardQuest

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	path2 "path"
	"strings"
	"sync/atomic"
	"initialthree/node/table/excel"
	"initialthree/zaplogger"
)

type RewardQuest struct{
    ID int32 
    Position int32 
    DungeonName string 
    AttackType_1 int32 
    AttackType_2 int32 
    AttackType_3 int32 
    QualityType_1 int32 
    QualityType_2 int32 
    QualityType_3 int32 
    ExecutionTime int32 
    DroppoolID int32 

}

var Table_ Table

type Table struct{
	xlsxName string
	nextPath string
	indexID  atomic.Value

}

func GetID(key int32)*RewardQuest{
	v := Table_.indexID.Load().(map[int32]*RewardQuest)[key]
	if v == nil{
		zaplogger.GetSugar().Errorf("Table %s GetID %d is nil", Table_.xlsxName, key)
	}
	return v
}

func GetIDMap() map[int32]*RewardQuest {
	return Table_.indexID.Load().(map[int32]*RewardQuest)
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
	Table_.xlsxName = "RewardQuest.xlsx"
	excel.AddTable(&Table_)
}

func (this *Table)Load(path string){
	xlFile, err := excelize.OpenFile(path2.Join(path2.Join(path, "/DataTable/世界地图及任务/悬赏任务(无用)"), this.xlsxName))
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	tableRows := xlFile.GetRows(xlFile.GetSheetName(xlFile.GetActiveSheetIndex()))
	names := tableRows[1]
	rows := tableRows[4:]

	tmpIDMap := map[int32]*RewardQuest{}

	for _,row := range rows{
		if row[0] != "" {
			e := &RewardQuest{}
            e.ID = excel.ReadInt32(readCellValue(row, names, "ID"))
            e.Position = excel.ReadInt32(readCellValue(row, names, "Position"))
            e.DungeonName = excel.ReadStr(readCellValue(row, names, "DungeonName"))
            e.AttackType_1 = excel.ReadInt32(readCellValue(row, names, "AttackType_1"))
            e.AttackType_2 = excel.ReadInt32(readCellValue(row, names, "AttackType_2"))
            e.AttackType_3 = excel.ReadInt32(readCellValue(row, names, "AttackType_3"))
            e.QualityType_1 = excel.ReadInt32(readCellValue(row, names, "QualityType_1"))
            e.QualityType_2 = excel.ReadInt32(readCellValue(row, names, "QualityType_2"))
            e.QualityType_3 = excel.ReadInt32(readCellValue(row, names, "QualityType_3"))
            e.ExecutionTime = excel.ReadInt32(readCellValue(row, names, "ExecutionTime"))
            e.DroppoolID = excel.ReadInt32(readCellValue(row, names, "DroppoolID"))

			tmpIDMap[excel.ReadInt32(readCellValue(row, names, "ID"))] = e

		}
	}
	this.indexID.Store(tmpIDMap)	

}
