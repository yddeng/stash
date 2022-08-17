package Player

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	path2 "path"
	"strings"
	"sync/atomic"
	"initialthree/node/table/excel"
	"initialthree/zaplogger"
)

type Player struct{
    ID int32 
    DefaultLevel int32 
    MaxLevel int32 
    FatigueValueRecoverValueEveryTime int32 
    FatigueValueRecoverTimeUnit int32 
    FatigueValueToExpRatio int32 
    FatigueSupplyCountEveryDay int32 
    GoldSupplyCountEveryDay int32 
    RenameCostItemID int32 

}

var Table_ Table

type Table struct{
	xlsxName string
	nextPath string
	indexID  atomic.Value

}

func GetID(key int32)*Player{
	v := Table_.indexID.Load().(map[int32]*Player)[key]
	if v == nil{
		zaplogger.GetSugar().Errorf("Table %s GetID %d is nil", Table_.xlsxName, key)
	}
	return v
}

func GetIDMap() map[int32]*Player {
	return Table_.indexID.Load().(map[int32]*Player)
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
	Table_.xlsxName = "Player.xlsx"
	excel.AddTable(&Table_)
}

func (this *Table)Load(path string){
	xlFile, err := excelize.OpenFile(path2.Join(path2.Join(path, "/ConstTable"), this.xlsxName))
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	tableRows := xlFile.GetRows(xlFile.GetSheetName(xlFile.GetActiveSheetIndex()))
	names := tableRows[1]
	rows := tableRows[4:]

	tmpIDMap := map[int32]*Player{}

	for _,row := range rows{
		if row[0] != "" {
			e := &Player{}
            e.ID = excel.ReadInt32(readCellValue(row, names, "ID"))
            e.DefaultLevel = excel.ReadInt32(readCellValue(row, names, "DefaultLevel"))
            e.MaxLevel = excel.ReadInt32(readCellValue(row, names, "MaxLevel"))
            e.FatigueValueRecoverValueEveryTime = excel.ReadInt32(readCellValue(row, names, "FatigueValueRecoverValueEveryTime"))
            e.FatigueValueRecoverTimeUnit = excel.ReadInt32(readCellValue(row, names, "FatigueValueRecoverTimeUnit"))
            e.FatigueValueToExpRatio = excel.ReadInt32(readCellValue(row, names, "FatigueValueToExpRatio"))
            e.FatigueSupplyCountEveryDay = excel.ReadInt32(readCellValue(row, names, "FatigueSupplyCountEveryDay"))
            e.GoldSupplyCountEveryDay = excel.ReadInt32(readCellValue(row, names, "GoldSupplyCountEveryDay"))
            e.RenameCostItemID = excel.ReadInt32(readCellValue(row, names, "RenameCostItemID"))

			tmpIDMap[excel.ReadInt32(readCellValue(row, names, "ID"))] = e

		}
	}
	this.indexID.Store(tmpIDMap)	

}
