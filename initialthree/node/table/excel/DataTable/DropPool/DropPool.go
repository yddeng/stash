package DropPool

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	path2 "path"
	"strings"
	"sync/atomic"
	"initialthree/node/table/excel"
	"initialthree/zaplogger"
)

type DropPool struct{
    ID int32 
    Desc string 
    Type string 
    MinCount int32 
    MaxCount int32 
    Repeatable bool 
    DropType_1 string 
    DropID_1 int32 
    DropCount_1 int32 
    Wave_1 int32 
    DropWeight_1 int32 
    DropType_2 string 
    DropID_2 int32 
    DropCount_2 int32 
    Wave_2 int32 
    DropWeight_2 int32 
    DropType_3 string 
    DropID_3 int32 
    DropCount_3 int32 
    Wave_3 int32 
    DropWeight_3 int32 
    DropType_4 string 
    DropID_4 int32 
    DropCount_4 int32 
    Wave_4 int32 
    DropWeight_4 int32 
    DropType_5 string 
    DropID_5 int32 
    DropCount_5 int32 
    Wave_5 int32 
    DropWeight_5 int32 
    DropType_6 string 
    DropID_6 int32 
    DropCount_6 int32 
    Wave_6 int32 
    DropWeight_6 int32 
    DropType_7 string 
    DropID_7 int32 
    DropCount_7 int32 
    Wave_7 int32 
    DropWeight_7 int32 
    DropType_8 string 
    DropID_8 int32 
    DropCount_8 int32 
    Wave_8 int32 
    DropWeight_8 int32 
    DropType_9 string 
    DropID_9 int32 
    DropCount_9 int32 
    Wave_9 int32 
    DropWeight_9 int32 
    DropType_10 string 
    DropID_10 int32 
    DropCount_10 int32 
    Wave_10 int32 
    DropWeight_10 int32 
    DropType_11 string 
    DropID_11 int32 
    DropCount_11 int32 
    Wave_11 int32 
    DropWeight_11 int32 
    DropType_12 string 
    DropID_12 int32 
    DropCount_12 int32 
    Wave_12 int32 
    DropWeight_12 int32 
    TypeEnum int32 
    DropList []*DropList_ 

}

type DropList_ struct{
    Type int32 
    ID int32 
    Count int32 
    Wave int32 
    Weight int32 

}

func readDropList(row, names []string)[]*DropList_{
	ret := make([]*DropList_, 0)
	base := excel.Split("DropType_,DropID_,DropCount_,Wave_,DropWeight_", ",")

	for i := 1; ; i++ {
		if idx, find := hasCellName(names, fmt.Sprintf("%s%d", base[0][0], i));!find {
			break
		}else if _,ok := hasCellValue(row,idx);!ok{
			break
		}
		
		e := &DropList_{}
        e.Type = excel.ReadEnum(readCellValue(row, names, fmt.Sprintf("%s%d", base[0][0], i)))
        e.ID = excel.ReadInt32(readCellValue(row, names, fmt.Sprintf("%s%d", base[1][0], i)))
        e.Count = excel.ReadInt32(readCellValue(row, names, fmt.Sprintf("%s%d", base[2][0], i)))
        e.Wave = excel.ReadInt32(readCellValue(row, names, fmt.Sprintf("%s%d", base[3][0], i)))
        e.Weight = excel.ReadInt32(readCellValue(row, names, fmt.Sprintf("%s%d", base[4][0], i)))

		ret = append(ret, e)
	}

	return ret
}

var Table_ Table

type Table struct{
	xlsxName string
	nextPath string
	indexID  atomic.Value

}

func GetID(key int32)*DropPool{
	v := Table_.indexID.Load().(map[int32]*DropPool)[key]
	if v == nil{
		zaplogger.GetSugar().Errorf("Table %s GetID %d is nil", Table_.xlsxName, key)
	}
	return v
}

func GetIDMap() map[int32]*DropPool {
	return Table_.indexID.Load().(map[int32]*DropPool)
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
	Table_.xlsxName = "DropPool.xlsx"
	excel.AddTable(&Table_)
}

func (this *Table)Load(path string){
	xlFile, err := excelize.OpenFile(path2.Join(path2.Join(path, "/DataTable/掉落"), this.xlsxName))
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	tableRows := xlFile.GetRows(xlFile.GetSheetName(xlFile.GetActiveSheetIndex()))
	names := tableRows[1]
	rows := tableRows[4:]

	tmpIDMap := map[int32]*DropPool{}

	for _,row := range rows{
		if row[0] != "" {
			e := &DropPool{}
            e.ID = excel.ReadInt32(readCellValue(row, names, "ID"))
            e.Desc = excel.ReadStr(readCellValue(row, names, "Desc"))
            e.Type = excel.ReadStr(readCellValue(row, names, "Type"))
            e.MinCount = excel.ReadInt32(readCellValue(row, names, "MinCount"))
            e.MaxCount = excel.ReadInt32(readCellValue(row, names, "MaxCount"))
            e.Repeatable = excel.ReadBool(readCellValue(row, names, "Repeatable"))
            e.DropType_1 = excel.ReadStr(readCellValue(row, names, "DropType_1"))
            e.DropID_1 = excel.ReadInt32(readCellValue(row, names, "DropID_1"))
            e.DropCount_1 = excel.ReadInt32(readCellValue(row, names, "DropCount_1"))
            e.Wave_1 = excel.ReadInt32(readCellValue(row, names, "Wave_1"))
            e.DropWeight_1 = excel.ReadInt32(readCellValue(row, names, "DropWeight_1"))
            e.DropType_2 = excel.ReadStr(readCellValue(row, names, "DropType_2"))
            e.DropID_2 = excel.ReadInt32(readCellValue(row, names, "DropID_2"))
            e.DropCount_2 = excel.ReadInt32(readCellValue(row, names, "DropCount_2"))
            e.Wave_2 = excel.ReadInt32(readCellValue(row, names, "Wave_2"))
            e.DropWeight_2 = excel.ReadInt32(readCellValue(row, names, "DropWeight_2"))
            e.DropType_3 = excel.ReadStr(readCellValue(row, names, "DropType_3"))
            e.DropID_3 = excel.ReadInt32(readCellValue(row, names, "DropID_3"))
            e.DropCount_3 = excel.ReadInt32(readCellValue(row, names, "DropCount_3"))
            e.Wave_3 = excel.ReadInt32(readCellValue(row, names, "Wave_3"))
            e.DropWeight_3 = excel.ReadInt32(readCellValue(row, names, "DropWeight_3"))
            e.DropType_4 = excel.ReadStr(readCellValue(row, names, "DropType_4"))
            e.DropID_4 = excel.ReadInt32(readCellValue(row, names, "DropID_4"))
            e.DropCount_4 = excel.ReadInt32(readCellValue(row, names, "DropCount_4"))
            e.Wave_4 = excel.ReadInt32(readCellValue(row, names, "Wave_4"))
            e.DropWeight_4 = excel.ReadInt32(readCellValue(row, names, "DropWeight_4"))
            e.DropType_5 = excel.ReadStr(readCellValue(row, names, "DropType_5"))
            e.DropID_5 = excel.ReadInt32(readCellValue(row, names, "DropID_5"))
            e.DropCount_5 = excel.ReadInt32(readCellValue(row, names, "DropCount_5"))
            e.Wave_5 = excel.ReadInt32(readCellValue(row, names, "Wave_5"))
            e.DropWeight_5 = excel.ReadInt32(readCellValue(row, names, "DropWeight_5"))
            e.DropType_6 = excel.ReadStr(readCellValue(row, names, "DropType_6"))
            e.DropID_6 = excel.ReadInt32(readCellValue(row, names, "DropID_6"))
            e.DropCount_6 = excel.ReadInt32(readCellValue(row, names, "DropCount_6"))
            e.Wave_6 = excel.ReadInt32(readCellValue(row, names, "Wave_6"))
            e.DropWeight_6 = excel.ReadInt32(readCellValue(row, names, "DropWeight_6"))
            e.DropType_7 = excel.ReadStr(readCellValue(row, names, "DropType_7"))
            e.DropID_7 = excel.ReadInt32(readCellValue(row, names, "DropID_7"))
            e.DropCount_7 = excel.ReadInt32(readCellValue(row, names, "DropCount_7"))
            e.Wave_7 = excel.ReadInt32(readCellValue(row, names, "Wave_7"))
            e.DropWeight_7 = excel.ReadInt32(readCellValue(row, names, "DropWeight_7"))
            e.DropType_8 = excel.ReadStr(readCellValue(row, names, "DropType_8"))
            e.DropID_8 = excel.ReadInt32(readCellValue(row, names, "DropID_8"))
            e.DropCount_8 = excel.ReadInt32(readCellValue(row, names, "DropCount_8"))
            e.Wave_8 = excel.ReadInt32(readCellValue(row, names, "Wave_8"))
            e.DropWeight_8 = excel.ReadInt32(readCellValue(row, names, "DropWeight_8"))
            e.DropType_9 = excel.ReadStr(readCellValue(row, names, "DropType_9"))
            e.DropID_9 = excel.ReadInt32(readCellValue(row, names, "DropID_9"))
            e.DropCount_9 = excel.ReadInt32(readCellValue(row, names, "DropCount_9"))
            e.Wave_9 = excel.ReadInt32(readCellValue(row, names, "Wave_9"))
            e.DropWeight_9 = excel.ReadInt32(readCellValue(row, names, "DropWeight_9"))
            e.DropType_10 = excel.ReadStr(readCellValue(row, names, "DropType_10"))
            e.DropID_10 = excel.ReadInt32(readCellValue(row, names, "DropID_10"))
            e.DropCount_10 = excel.ReadInt32(readCellValue(row, names, "DropCount_10"))
            e.Wave_10 = excel.ReadInt32(readCellValue(row, names, "Wave_10"))
            e.DropWeight_10 = excel.ReadInt32(readCellValue(row, names, "DropWeight_10"))
            e.DropType_11 = excel.ReadStr(readCellValue(row, names, "DropType_11"))
            e.DropID_11 = excel.ReadInt32(readCellValue(row, names, "DropID_11"))
            e.DropCount_11 = excel.ReadInt32(readCellValue(row, names, "DropCount_11"))
            e.Wave_11 = excel.ReadInt32(readCellValue(row, names, "Wave_11"))
            e.DropWeight_11 = excel.ReadInt32(readCellValue(row, names, "DropWeight_11"))
            e.DropType_12 = excel.ReadStr(readCellValue(row, names, "DropType_12"))
            e.DropID_12 = excel.ReadInt32(readCellValue(row, names, "DropID_12"))
            e.DropCount_12 = excel.ReadInt32(readCellValue(row, names, "DropCount_12"))
            e.Wave_12 = excel.ReadInt32(readCellValue(row, names, "Wave_12"))
            e.DropWeight_12 = excel.ReadInt32(readCellValue(row, names, "DropWeight_12"))
            e.TypeEnum = excel.ReadEnum(readCellValue(row, names, "Type"))
            e.DropList = readDropList(row, names)

			tmpIDMap[excel.ReadInt32(readCellValue(row, names, "ID"))] = e

		}
	}
	this.indexID.Store(tmpIDMap)	

}
