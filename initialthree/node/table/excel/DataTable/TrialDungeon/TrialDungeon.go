package TrialDungeon

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	path2 "path"
	"strings"
	"sync/atomic"
	"initialthree/node/table/excel"
	"initialthree/zaplogger"
)

type TrialDungeon struct{
    ID int32 
    Distance int32 
    Fighting int32 
    Title_1 string 
    Desc_1 string 
    Title_2 string 
    Desc_2 string 
    Title_3 string 
    Desc_3 string 
    TaskTitle string 
    Difficult bool 
    NextTrialDungeonID int32 
    TrialDungeonID int32 
    TrialCount int32 
    DropPoolID int32 
    DungeonID int32 

}

var Table_ Table

type Table struct{
	xlsxName string
	nextPath string
	indexID  atomic.Value

}

func GetID(key int32)*TrialDungeon{
	v := Table_.indexID.Load().(map[int32]*TrialDungeon)[key]
	if v == nil{
		zaplogger.GetSugar().Errorf("Table %s GetID %d is nil", Table_.xlsxName, key)
	}
	return v
}

func GetIDMap() map[int32]*TrialDungeon {
	return Table_.indexID.Load().(map[int32]*TrialDungeon)
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
	Table_.xlsxName = "TrialDungeon.xlsx"
	excel.AddTable(&Table_)
}

func (this *Table)Load(path string){
	xlFile, err := excelize.OpenFile(path2.Join(path2.Join(path, "/DataTable/关卡/试炼塔"), this.xlsxName))
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	tableRows := xlFile.GetRows(xlFile.GetSheetName(xlFile.GetActiveSheetIndex()))
	names := tableRows[1]
	rows := tableRows[4:]

	tmpIDMap := map[int32]*TrialDungeon{}

	for _,row := range rows{
		if row[0] != "" {
			e := &TrialDungeon{}
            e.ID = excel.ReadInt32(readCellValue(row, names, "ID"))
            e.Distance = excel.ReadInt32(readCellValue(row, names, "Distance"))
            e.Fighting = excel.ReadInt32(readCellValue(row, names, "Fighting"))
            e.Title_1 = excel.ReadStr(readCellValue(row, names, "Title_1"))
            e.Desc_1 = excel.ReadStr(readCellValue(row, names, "Desc_1"))
            e.Title_2 = excel.ReadStr(readCellValue(row, names, "Title_2"))
            e.Desc_2 = excel.ReadStr(readCellValue(row, names, "Desc_2"))
            e.Title_3 = excel.ReadStr(readCellValue(row, names, "Title_3"))
            e.Desc_3 = excel.ReadStr(readCellValue(row, names, "Desc_3"))
            e.TaskTitle = excel.ReadStr(readCellValue(row, names, "TaskTitle"))
            e.Difficult = excel.ReadBool(readCellValue(row, names, "Difficult"))
            e.NextTrialDungeonID = excel.ReadInt32(readCellValue(row, names, "NextTrialDungeonID"))
            e.TrialDungeonID = excel.ReadInt32(readCellValue(row, names, "TrialDungeonID"))
            e.TrialCount = excel.ReadInt32(readCellValue(row, names, "TrialCount"))
            e.DropPoolID = excel.ReadInt32(readCellValue(row, names, "DropPoolID"))
            e.DungeonID = excel.ReadInt32(readCellValue(row, names, "DungeonID"))

			tmpIDMap[excel.ReadInt32(readCellValue(row, names, "ID"))] = e

		}
	}
	this.indexID.Store(tmpIDMap)	

}
