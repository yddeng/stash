package Sign

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	path2 "path"
	"strings"
	"sync/atomic"
	"initialthree/node/table/excel"
	"initialthree/zaplogger"
)

type Sign struct{
    ID int32 
    Type string 
    StartTime string 
    EndTime string 
    SignTimes int32 
    AwardRule int32 
    TypeEnum int32 

}

var Table_ Table

type Table struct{
	xlsxName string
	nextPath string
	indexID  atomic.Value

}

func GetID(key int32)*Sign{
	v := Table_.indexID.Load().(map[int32]*Sign)[key]
	if v == nil{
		zaplogger.GetSugar().Errorf("Table %s GetID %d is nil", Table_.xlsxName, key)
	}
	return v
}

func GetIDMap() map[int32]*Sign {
	return Table_.indexID.Load().(map[int32]*Sign)
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
	Table_.xlsxName = "Sign.xlsx"
	excel.AddTable(&Table_)
}

func (this *Table)Load(path string){
	xlFile, err := excelize.OpenFile(path2.Join(path2.Join(path, "/DataTable/签到"), this.xlsxName))
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	tableRows := xlFile.GetRows(xlFile.GetSheetName(xlFile.GetActiveSheetIndex()))
	names := tableRows[1]
	rows := tableRows[4:]

	tmpIDMap := map[int32]*Sign{}

	for _,row := range rows{
		if row[0] != "" {
			e := &Sign{}
            e.ID = excel.ReadInt32(readCellValue(row, names, "ID"))
            e.Type = excel.ReadStr(readCellValue(row, names, "Type"))
            e.StartTime = excel.ReadStr(readCellValue(row, names, "StartTime"))
            e.EndTime = excel.ReadStr(readCellValue(row, names, "EndTime"))
            e.SignTimes = excel.ReadInt32(readCellValue(row, names, "SignTimes"))
            e.AwardRule = excel.ReadInt32(readCellValue(row, names, "AwardRule"))
            e.TypeEnum = excel.ReadEnum(readCellValue(row, names, "Type"))

			tmpIDMap[excel.ReadInt32(readCellValue(row, names, "ID"))] = e

		}
	}
	this.indexID.Store(tmpIDMap)	

}
