package PlayerPortrait

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	path2 "path"
	"strings"
	"sync/atomic"
	"initialthree/node/table/excel"
	"initialthree/zaplogger"
)

type PlayerPortrait struct{
    ID int32 
    Name string 
    Desc string 
    AccessWays string 
    Path string 

}

var Table_ Table

type Table struct{
	xlsxName string
	nextPath string
	indexID  atomic.Value

}

func GetID(key int32)*PlayerPortrait{
	v := Table_.indexID.Load().(map[int32]*PlayerPortrait)[key]
	if v == nil{
		zaplogger.GetSugar().Errorf("Table %s GetID %d is nil", Table_.xlsxName, key)
	}
	return v
}

func GetIDMap() map[int32]*PlayerPortrait {
	return Table_.indexID.Load().(map[int32]*PlayerPortrait)
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
	Table_.xlsxName = "PlayerPortrait.xlsx"
	excel.AddTable(&Table_)
}

func (this *Table)Load(path string){
	xlFile, err := excelize.OpenFile(path2.Join(path2.Join(path, "/DataTable/玩家相关/头像、头像框、名片"), this.xlsxName))
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	tableRows := xlFile.GetRows(xlFile.GetSheetName(xlFile.GetActiveSheetIndex()))
	names := tableRows[1]
	rows := tableRows[4:]

	tmpIDMap := map[int32]*PlayerPortrait{}

	for _,row := range rows{
		if row[0] != "" {
			e := &PlayerPortrait{}
            e.ID = excel.ReadInt32(readCellValue(row, names, "ID"))
            e.Name = excel.ReadStr(readCellValue(row, names, "Name"))
            e.Desc = excel.ReadStr(readCellValue(row, names, "Desc"))
            e.AccessWays = excel.ReadStr(readCellValue(row, names, "AccessWays"))
            e.Path = excel.ReadStr(readCellValue(row, names, "Path"))

			tmpIDMap[excel.ReadInt32(readCellValue(row, names, "ID"))] = e

		}
	}
	this.indexID.Store(tmpIDMap)	

}
