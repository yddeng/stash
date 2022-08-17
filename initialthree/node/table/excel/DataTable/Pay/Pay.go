package Pay

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	path2 "path"
	"strings"
	"sync/atomic"
	"initialthree/node/table/excel"
	"initialthree/zaplogger"
)

type Pay struct{
    ID int32 
    Desc string 
    Type string 
    CommonPayCount int32 
    CommonPresentedCount int32 
    FirstPayCount int32 
    FirstPresentedCount int32 
    BigIcon string 

}

var Table_ Table

type Table struct{
	xlsxName string
	nextPath string
	indexID  atomic.Value

}

func GetID(key int32)*Pay{
	v := Table_.indexID.Load().(map[int32]*Pay)[key]
	if v == nil{
		zaplogger.GetSugar().Errorf("Table %s GetID %d is nil", Table_.xlsxName, key)
	}
	return v
}

func GetIDMap() map[int32]*Pay {
	return Table_.indexID.Load().(map[int32]*Pay)
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
	Table_.xlsxName = "Pay.xlsx"
	excel.AddTable(&Table_)
}

func (this *Table)Load(path string){
	xlFile, err := excelize.OpenFile(path2.Join(path2.Join(path, "/DataTable/商店"), this.xlsxName))
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	tableRows := xlFile.GetRows(xlFile.GetSheetName(xlFile.GetActiveSheetIndex()))
	names := tableRows[1]
	rows := tableRows[4:]

	tmpIDMap := map[int32]*Pay{}

	for _,row := range rows{
		if row[0] != "" {
			e := &Pay{}
            e.ID = excel.ReadInt32(readCellValue(row, names, "ID"))
            e.Desc = excel.ReadStr(readCellValue(row, names, "Desc"))
            e.Type = excel.ReadStr(readCellValue(row, names, "Type"))
            e.CommonPayCount = excel.ReadInt32(readCellValue(row, names, "CommonPayCount"))
            e.CommonPresentedCount = excel.ReadInt32(readCellValue(row, names, "CommonPresentedCount"))
            e.FirstPayCount = excel.ReadInt32(readCellValue(row, names, "FirstPayCount"))
            e.FirstPresentedCount = excel.ReadInt32(readCellValue(row, names, "FirstPresentedCount"))
            e.BigIcon = excel.ReadStr(readCellValue(row, names, "bigIcon"))

			tmpIDMap[excel.ReadInt32(readCellValue(row, names, "ID"))] = e

		}
	}
	this.indexID.Store(tmpIDMap)	

}
