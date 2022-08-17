package PlayerLevel

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	path2 "path"
	"strings"
	"sync/atomic"
	"initialthree/node/table/excel"
	"initialthree/zaplogger"
)

type PlayerLevel struct{
    ID int32 
    MaxFatigueValue int32 
    MaxExpValue int32 
    GiveFatigue int32 

}

var Table_ Table

type Table struct{
	xlsxName string
	nextPath string
	indexID  atomic.Value

}

func GetID(key int32)*PlayerLevel{
	v := Table_.indexID.Load().(map[int32]*PlayerLevel)[key]
	if v == nil{
		zaplogger.GetSugar().Errorf("Table %s GetID %d is nil", Table_.xlsxName, key)
	}
	return v
}

func GetIDMap() map[int32]*PlayerLevel {
	return Table_.indexID.Load().(map[int32]*PlayerLevel)
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
	Table_.xlsxName = "PlayerLevel.xlsx"
	excel.AddTable(&Table_)
}

func (this *Table)Load(path string){
	xlFile, err := excelize.OpenFile(path2.Join(path2.Join(path, "/DataTable/玩家相关"), this.xlsxName))
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	tableRows := xlFile.GetRows(xlFile.GetSheetName(xlFile.GetActiveSheetIndex()))
	names := tableRows[1]
	rows := tableRows[4:]

	tmpIDMap := map[int32]*PlayerLevel{}

	for _,row := range rows{
		if row[0] != "" {
			e := &PlayerLevel{}
            e.ID = excel.ReadInt32(readCellValue(row, names, "ID"))
            e.MaxFatigueValue = excel.ReadInt32(readCellValue(row, names, "MaxFatigueValue"))
            e.MaxExpValue = excel.ReadInt32(readCellValue(row, names, "MaxExpValue"))
            e.GiveFatigue = excel.ReadInt32(readCellValue(row, names, "GiveFatigue"))

			tmpIDMap[excel.ReadInt32(readCellValue(row, names, "ID"))] = e

		}
	}
	this.indexID.Store(tmpIDMap)	

}
