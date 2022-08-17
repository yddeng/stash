package MaterialDungeon

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	path2 "path"
	"strings"
	"sync/atomic"
	"initialthree/node/table/excel"
	"initialthree/zaplogger"
)

type MaterialDungeon struct{
    ID int32 
    DungeonGroupName string 
    DungeonName string 
    DungeonIcon string 
    OpenTime string 
    CombatPrompt string 
    Order int32 
    Type string 
    DropDisplay int32 
    DungeonID int32 
    MultipleChallengeDesc string 
    OpenTimeArray []*OpenTime_ 

}

type OpenTime_ struct{
    Weekday int32 

}

func readOpenTimeArray(row, names []string)[]*OpenTime_{
	value := readCellValue(row, names, "OpenTime")
	if value == "" {
		return nil
	}

	r := excel.Split(value,",")
	if len(r) == 0 {
		return nil
	}

	ret := make([]*OpenTime_,0)
	for _, v := range r {
		if len(v) == 1{
			e := &OpenTime_{}
	        e.Weekday = excel.ReadInt32(v[0])

			ret = append(ret, e)
		}
	}

	return ret
}

var Table_ Table

type Table struct{
	xlsxName string
	nextPath string
	indexID  atomic.Value

}

func GetID(key int32)*MaterialDungeon{
	v := Table_.indexID.Load().(map[int32]*MaterialDungeon)[key]
	if v == nil{
		zaplogger.GetSugar().Errorf("Table %s GetID %d is nil", Table_.xlsxName, key)
	}
	return v
}

func GetIDMap() map[int32]*MaterialDungeon {
	return Table_.indexID.Load().(map[int32]*MaterialDungeon)
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
	Table_.xlsxName = "MaterialDungeon.xlsx"
	excel.AddTable(&Table_)
}

func (this *Table)Load(path string){
	xlFile, err := excelize.OpenFile(path2.Join(path2.Join(path, "/DataTable/关卡/材料"), this.xlsxName))
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	tableRows := xlFile.GetRows(xlFile.GetSheetName(xlFile.GetActiveSheetIndex()))
	names := tableRows[1]
	rows := tableRows[4:]

	tmpIDMap := map[int32]*MaterialDungeon{}

	for _,row := range rows{
		if row[0] != "" {
			e := &MaterialDungeon{}
            e.ID = excel.ReadInt32(readCellValue(row, names, "ID"))
            e.DungeonGroupName = excel.ReadStr(readCellValue(row, names, "DungeonGroupName"))
            e.DungeonName = excel.ReadStr(readCellValue(row, names, "DungeonName"))
            e.DungeonIcon = excel.ReadStr(readCellValue(row, names, "DungeonIcon"))
            e.OpenTime = excel.ReadStr(readCellValue(row, names, "OpenTime"))
            e.CombatPrompt = excel.ReadStr(readCellValue(row, names, "CombatPrompt"))
            e.Order = excel.ReadInt32(readCellValue(row, names, "Order"))
            e.Type = excel.ReadStr(readCellValue(row, names, "Type"))
            e.DropDisplay = excel.ReadInt32(readCellValue(row, names, "DropDisplay"))
            e.DungeonID = excel.ReadInt32(readCellValue(row, names, "DungeonID"))
            e.MultipleChallengeDesc = excel.ReadStr(readCellValue(row, names, "MultipleChallengeDesc"))
            e.OpenTimeArray = readOpenTimeArray(row, names)

			tmpIDMap[excel.ReadInt32(readCellValue(row, names, "ID"))] = e

		}
	}
	this.indexID.Store(tmpIDMap)	

}
