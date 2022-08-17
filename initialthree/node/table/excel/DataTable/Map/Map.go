package Map

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	path2 "path"
	"strings"
	"sync/atomic"
	"initialthree/node/table/excel"
	"initialthree/zaplogger"
)

type Map struct{
    ID int32 
    Length int32 
    Width int32 
    DefaultPosition string 
    DefaultRotation int32 
    DefaultPositionStruct *DefaultPosition_ 

}

type DefaultPosition_ struct{
    X int32 
    Y int32 
    Z int32 

}

func readDefaultPositionStruct(row, names []string)*DefaultPosition_{
	value := readCellValue(row, names, "DefaultPosition")
	r := excel.Split(value,",")
	
	e := &DefaultPosition_{}
    e.X = excel.ReadInt32(r[0][0])
    e.Y = excel.ReadInt32(r[1][0])
    e.Z = excel.ReadInt32(r[2][0])

	return e
}

var Table_ Table

type Table struct{
	xlsxName string
	nextPath string
	indexID  atomic.Value

}

func GetID(key int32)*Map{
	v := Table_.indexID.Load().(map[int32]*Map)[key]
	if v == nil{
		zaplogger.GetSugar().Errorf("Table %s GetID %d is nil", Table_.xlsxName, key)
	}
	return v
}

func GetIDMap() map[int32]*Map {
	return Table_.indexID.Load().(map[int32]*Map)
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
	Table_.xlsxName = "Map.xlsx"
	excel.AddTable(&Table_)
}

func (this *Table)Load(path string){
	xlFile, err := excelize.OpenFile(path2.Join(path2.Join(path, "/DataTable"), this.xlsxName))
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	tableRows := xlFile.GetRows(xlFile.GetSheetName(xlFile.GetActiveSheetIndex()))
	names := tableRows[1]
	rows := tableRows[4:]

	tmpIDMap := map[int32]*Map{}

	for _,row := range rows{
		if row[0] != "" {
			e := &Map{}
            e.ID = excel.ReadInt32(readCellValue(row, names, "ID"))
            e.Length = excel.ReadInt32(readCellValue(row, names, "Length"))
            e.Width = excel.ReadInt32(readCellValue(row, names, "Width"))
            e.DefaultPosition = excel.ReadStr(readCellValue(row, names, "DefaultPosition"))
            e.DefaultRotation = excel.ReadInt32(readCellValue(row, names, "DefaultRotation"))
            e.DefaultPositionStruct = readDefaultPositionStruct(row, names)

			tmpIDMap[excel.ReadInt32(readCellValue(row, names, "ID"))] = e

		}
	}
	this.indexID.Store(tmpIDMap)	

}
