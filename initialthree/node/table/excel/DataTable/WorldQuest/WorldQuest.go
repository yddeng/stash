package WorldQuest

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	path2 "path"
	"strings"
	"sync/atomic"
	"initialthree/node/table/excel"
	"initialthree/zaplogger"
)

type WorldQuest struct{
    ID int32 
    QuestName string 
    QuestPosition int32 
    QuestPositionImage string 
    QuestQualityType string 
    CampType string 
    QuestTypeDesc string 
    QuestDetail string 
    DungeonID int32 
    DroppoolID int32 
    ReputationReward int32 
    CampTypeEnum int32 

}

var Table_ Table

type Table struct{
	xlsxName string
	nextPath string
	indexID  atomic.Value

}

func GetID(key int32)*WorldQuest{
	v := Table_.indexID.Load().(map[int32]*WorldQuest)[key]
	if v == nil{
		zaplogger.GetSugar().Errorf("Table %s GetID %d is nil", Table_.xlsxName, key)
	}
	return v
}

func GetIDMap() map[int32]*WorldQuest {
	return Table_.indexID.Load().(map[int32]*WorldQuest)
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
	Table_.xlsxName = "WorldQuest.xlsx"
	excel.AddTable(&Table_)
}

func (this *Table)Load(path string){
	xlFile, err := excelize.OpenFile(path2.Join(path2.Join(path, "/DataTable/关卡/世界任务"), this.xlsxName))
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	tableRows := xlFile.GetRows(xlFile.GetSheetName(xlFile.GetActiveSheetIndex()))
	names := tableRows[1]
	rows := tableRows[4:]

	tmpIDMap := map[int32]*WorldQuest{}

	for _,row := range rows{
		if row[0] != "" {
			e := &WorldQuest{}
            e.ID = excel.ReadInt32(readCellValue(row, names, "ID"))
            e.QuestName = excel.ReadStr(readCellValue(row, names, "QuestName"))
            e.QuestPosition = excel.ReadInt32(readCellValue(row, names, "QuestPosition"))
            e.QuestPositionImage = excel.ReadStr(readCellValue(row, names, "QuestPositionImage"))
            e.QuestQualityType = excel.ReadStr(readCellValue(row, names, "QuestQualityType"))
            e.CampType = excel.ReadStr(readCellValue(row, names, "CampType"))
            e.QuestTypeDesc = excel.ReadStr(readCellValue(row, names, "QuestTypeDesc"))
            e.QuestDetail = excel.ReadStr(readCellValue(row, names, "QuestDetail"))
            e.DungeonID = excel.ReadInt32(readCellValue(row, names, "DungeonID"))
            e.DroppoolID = excel.ReadInt32(readCellValue(row, names, "DroppoolID"))
            e.ReputationReward = excel.ReadInt32(readCellValue(row, names, "ReputationReward"))
            e.CampTypeEnum = excel.ReadEnum(readCellValue(row, names, "CampType"))

			tmpIDMap[excel.ReadInt32(readCellValue(row, names, "ID"))] = e

		}
	}
	this.indexID.Store(tmpIDMap)	

}
