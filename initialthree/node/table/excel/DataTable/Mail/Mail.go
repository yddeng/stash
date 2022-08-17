package Mail

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	path2 "path"
	"strings"
	"sync/atomic"
	"initialthree/node/table/excel"
	"initialthree/zaplogger"
)

type Mail struct{
    ID int32 
    Describe string 
    Title string 
    Sender string 
    Expire int32 
    Content string 

}

var Table_ Table

type Table struct{
	xlsxName string
	nextPath string
	indexID  atomic.Value

}

func GetID(key int32)*Mail{
	v := Table_.indexID.Load().(map[int32]*Mail)[key]
	if v == nil{
		zaplogger.GetSugar().Errorf("Table %s GetID %d is nil", Table_.xlsxName, key)
	}
	return v
}

func GetIDMap() map[int32]*Mail {
	return Table_.indexID.Load().(map[int32]*Mail)
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
	Table_.xlsxName = "Mail.xlsx"
	excel.AddTable(&Table_)
}

func (this *Table)Load(path string){
	xlFile, err := excelize.OpenFile(path2.Join(path2.Join(path, "/DataTable/邮件"), this.xlsxName))
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	tableRows := xlFile.GetRows(xlFile.GetSheetName(xlFile.GetActiveSheetIndex()))
	names := tableRows[1]
	rows := tableRows[4:]

	tmpIDMap := map[int32]*Mail{}

	for _,row := range rows{
		if row[0] != "" {
			e := &Mail{}
            e.ID = excel.ReadInt32(readCellValue(row, names, "ID"))
            e.Describe = excel.ReadStr(readCellValue(row, names, "Describe"))
            e.Title = excel.ReadStr(readCellValue(row, names, "Title"))
            e.Sender = excel.ReadStr(readCellValue(row, names, "Sender"))
            e.Expire = excel.ReadInt32(readCellValue(row, names, "Expire"))
            e.Content = excel.ReadStr(readCellValue(row, names, "Content"))

			tmpIDMap[excel.ReadInt32(readCellValue(row, names, "ID"))] = e

		}
	}
	this.indexID.Store(tmpIDMap)	

}
