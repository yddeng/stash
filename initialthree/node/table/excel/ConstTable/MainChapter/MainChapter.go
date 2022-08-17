package MainChapter

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	path2 "path"
	"strings"
	"sync/atomic"
	"initialthree/node/table/excel"
	"initialthree/zaplogger"
)

type MainChapter struct{
    ID int32 
    SpecialUnlockChapter int32 
    SpecialUnlockStarDate string 
    SpecialUnlockEndDate string 
    SpecialUnlockStarDateStruct *SpecialUnlockStarDate_ 
    SpecialUnlockEndDateStruct *SpecialUnlockEndDate_ 

}

type SpecialUnlockStarDate_ struct{
    Year int32 
    Month int32 
    Day int32 
    Hour int32 
    Min int32 

}

func readSpecialUnlockStarDateStruct(row, names []string)*SpecialUnlockStarDate_{
	value := readCellValue(row, names, "SpecialUnlockStarDate")
	r := excel.Split(value,".")
	
	e := &SpecialUnlockStarDate_{}
    e.Year = excel.ReadInt32(r[0][0])
    e.Month = excel.ReadInt32(r[1][0])
    e.Day = excel.ReadInt32(r[2][0])
    e.Hour = excel.ReadInt32(r[3][0])
    e.Min = excel.ReadInt32(r[4][0])

	return e
}

type SpecialUnlockEndDate_ struct{
    Year int32 
    Month int32 
    Day int32 
    Hour int32 
    Min int32 

}

func readSpecialUnlockEndDateStruct(row, names []string)*SpecialUnlockEndDate_{
	value := readCellValue(row, names, "SpecialUnlockEndDate")
	r := excel.Split(value,".")
	
	e := &SpecialUnlockEndDate_{}
    e.Year = excel.ReadInt32(r[0][0])
    e.Month = excel.ReadInt32(r[1][0])
    e.Day = excel.ReadInt32(r[2][0])
    e.Hour = excel.ReadInt32(r[3][0])
    e.Min = excel.ReadInt32(r[4][0])

	return e
}

var Table_ Table

type Table struct{
	xlsxName string
	nextPath string
	indexID  atomic.Value

}

func GetID(key int32)*MainChapter{
	v := Table_.indexID.Load().(map[int32]*MainChapter)[key]
	if v == nil{
		zaplogger.GetSugar().Errorf("Table %s GetID %d is nil", Table_.xlsxName, key)
	}
	return v
}

func GetIDMap() map[int32]*MainChapter {
	return Table_.indexID.Load().(map[int32]*MainChapter)
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
	Table_.xlsxName = "MainChapter.xlsx"
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

	tmpIDMap := map[int32]*MainChapter{}

	for _,row := range rows{
		if row[0] != "" {
			e := &MainChapter{}
            e.ID = excel.ReadInt32(readCellValue(row, names, "ID"))
            e.SpecialUnlockChapter = excel.ReadInt32(readCellValue(row, names, "SpecialUnlockChapter"))
            e.SpecialUnlockStarDate = excel.ReadStr(readCellValue(row, names, "SpecialUnlockStarDate"))
            e.SpecialUnlockEndDate = excel.ReadStr(readCellValue(row, names, "SpecialUnlockEndDate"))
            e.SpecialUnlockStarDateStruct = readSpecialUnlockStarDateStruct(row, names)
            e.SpecialUnlockEndDateStruct = readSpecialUnlockEndDateStruct(row, names)

			tmpIDMap[excel.ReadInt32(readCellValue(row, names, "ID"))] = e

		}
	}
	this.indexID.Store(tmpIDMap)	

}
