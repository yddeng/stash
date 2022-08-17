package MainDungeon

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	path2 "path"
	"strings"
	"sync/atomic"
	"initialthree/node/table/excel"
	"initialthree/zaplogger"
)

type MainDungeon struct{
    ID int32 
    ChapterID int32 
    NextMainDungeonID int32 
    DungeonNamePrefix string 
    DungeonName string 
    Desc string 
    QuestId int32 
    FightingCapacityAlert int32 
    DropDisplay int32 
    DungeonID int32 
    DungeonOrder int32 
    OldVersionConf string 
    DungeonOrderCn string 
    IsStory bool 
    NormalNodeIcon string 
    IconNodeIcon string 

}

var Table_ Table

type Table struct{
	xlsxName string
	nextPath string
	indexID  atomic.Value

}

func GetID(key int32)*MainDungeon{
	v := Table_.indexID.Load().(map[int32]*MainDungeon)[key]
	if v == nil{
		zaplogger.GetSugar().Errorf("Table %s GetID %d is nil", Table_.xlsxName, key)
	}
	return v
}

func GetIDMap() map[int32]*MainDungeon {
	return Table_.indexID.Load().(map[int32]*MainDungeon)
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
	Table_.xlsxName = "MainDungeon.xlsx"
	excel.AddTable(&Table_)
}

func (this *Table)Load(path string){
	xlFile, err := excelize.OpenFile(path2.Join(path2.Join(path, "/DataTable/关卡/主线"), this.xlsxName))
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	tableRows := xlFile.GetRows(xlFile.GetSheetName(xlFile.GetActiveSheetIndex()))
	names := tableRows[1]
	rows := tableRows[4:]

	tmpIDMap := map[int32]*MainDungeon{}

	for _,row := range rows{
		if row[0] != "" {
			e := &MainDungeon{}
            e.ID = excel.ReadInt32(readCellValue(row, names, "ID"))
            e.ChapterID = excel.ReadInt32(readCellValue(row, names, "ChapterID"))
            e.NextMainDungeonID = excel.ReadInt32(readCellValue(row, names, "NextMainDungeonID"))
            e.DungeonNamePrefix = excel.ReadStr(readCellValue(row, names, "DungeonNamePrefix"))
            e.DungeonName = excel.ReadStr(readCellValue(row, names, "DungeonName"))
            e.Desc = excel.ReadStr(readCellValue(row, names, "Desc"))
            e.QuestId = excel.ReadInt32(readCellValue(row, names, "QuestId"))
            e.FightingCapacityAlert = excel.ReadInt32(readCellValue(row, names, "FightingCapacityAlert"))
            e.DropDisplay = excel.ReadInt32(readCellValue(row, names, "DropDisplay"))
            e.DungeonID = excel.ReadInt32(readCellValue(row, names, "DungeonID"))
            e.DungeonOrder = excel.ReadInt32(readCellValue(row, names, "DungeonOrder"))
            e.OldVersionConf = excel.ReadStr(readCellValue(row, names, "OldVersionConf"))
            e.DungeonOrderCn = excel.ReadStr(readCellValue(row, names, "DungeonOrderCn"))
            e.IsStory = excel.ReadBool(readCellValue(row, names, "IsStory"))
            e.NormalNodeIcon = excel.ReadStr(readCellValue(row, names, "NormalNodeIcon"))
            e.IconNodeIcon = excel.ReadStr(readCellValue(row, names, "IconNodeIcon"))

			tmpIDMap[excel.ReadInt32(readCellValue(row, names, "ID"))] = e

		}
	}
	this.indexID.Store(tmpIDMap)	

}
