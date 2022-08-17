package AccountInitializeAssets

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	path2 "path"
	"strings"
	"sync/atomic"
	"initialthree/node/table/excel"
	"initialthree/zaplogger"
)

type AccountInitializeAssets struct{
    ID int32 
    Comment string 
    AssetType string 
    AssetID int32 
    Count int32 
    AssetTypeEnum int32 

}

var Table_ Table

type Table struct{
	xlsxName string
	nextPath string
	indexID  atomic.Value

}

func GetID(key int32)*AccountInitializeAssets{
	v := Table_.indexID.Load().(map[int32]*AccountInitializeAssets)[key]
	if v == nil{
		zaplogger.GetSugar().Errorf("Table %s GetID %d is nil", Table_.xlsxName, key)
	}
	return v
}

func GetIDMap() map[int32]*AccountInitializeAssets {
	return Table_.indexID.Load().(map[int32]*AccountInitializeAssets)
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
	Table_.xlsxName = "AccountInitializeAssets.xlsx"
	excel.AddTable(&Table_)
}

func (this *Table)Load(path string){
	xlFile, err := excelize.OpenFile(path2.Join(path2.Join(path, "/DataTable/账号初始资源"), this.xlsxName))
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	tableRows := xlFile.GetRows(xlFile.GetSheetName(xlFile.GetActiveSheetIndex()))
	names := tableRows[1]
	rows := tableRows[4:]

	tmpIDMap := map[int32]*AccountInitializeAssets{}

	for _,row := range rows{
		if row[0] != "" {
			e := &AccountInitializeAssets{}
            e.ID = excel.ReadInt32(readCellValue(row, names, "ID"))
            e.Comment = excel.ReadStr(readCellValue(row, names, "Comment"))
            e.AssetType = excel.ReadStr(readCellValue(row, names, "AssetType"))
            e.AssetID = excel.ReadInt32(readCellValue(row, names, "AssetID"))
            e.Count = excel.ReadInt32(readCellValue(row, names, "Count"))
            e.AssetTypeEnum = excel.ReadEnum(readCellValue(row, names, "AssetType"))

			tmpIDMap[excel.ReadInt32(readCellValue(row, names, "ID"))] = e

		}
	}
	this.indexID.Store(tmpIDMap)	

}
