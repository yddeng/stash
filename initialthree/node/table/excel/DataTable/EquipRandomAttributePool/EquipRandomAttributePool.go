package EquipRandomAttributePool

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	path2 "path"
	"strings"
	"sync/atomic"
	"initialthree/node/table/excel"
	"initialthree/zaplogger"
)

type EquipRandomAttributePool struct{
    ID int32 
    ID_1 int32 
    Weight_1 int32 
    ID_2 int32 
    Weight_2 int32 
    ID_3 int32 
    Weight_3 int32 
    ID_4 int32 
    Weight_4 int32 
    ID_5 int32 
    Weight_5 int32 
    ID_6 int32 
    Weight_6 int32 
    ID_7 int32 
    Weight_7 int32 
    ID_8 int32 
    Weight_8 int32 
    ID_9 int32 
    Weight_9 int32 
    ID_10 int32 
    Weight_10 int32 
    Random []*Random_ 

}

type Random_ struct{
    ID int32 
    Weight int32 

}

func readRandom(row, names []string)[]*Random_{
	ret := make([]*Random_, 0)
	base := excel.Split("ID_,Weight_", ",")

	for i := 1; ; i++ {
		if idx, find := hasCellName(names, fmt.Sprintf("%s%d", base[0][0], i));!find {
			break
		}else if _,ok := hasCellValue(row,idx);!ok{
			break
		}
		
		e := &Random_{}
        e.ID = excel.ReadInt32(readCellValue(row, names, fmt.Sprintf("%s%d", base[0][0], i)))
        e.Weight = excel.ReadInt32(readCellValue(row, names, fmt.Sprintf("%s%d", base[1][0], i)))

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

func GetID(key int32)*EquipRandomAttributePool{
	v := Table_.indexID.Load().(map[int32]*EquipRandomAttributePool)[key]
	if v == nil{
		zaplogger.GetSugar().Errorf("Table %s GetID %d is nil", Table_.xlsxName, key)
	}
	return v
}

func GetIDMap() map[int32]*EquipRandomAttributePool {
	return Table_.indexID.Load().(map[int32]*EquipRandomAttributePool)
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
	Table_.xlsxName = "EquipRandomAttributePool.xlsx"
	excel.AddTable(&Table_)
}

func (this *Table)Load(path string){
	xlFile, err := excelize.OpenFile(path2.Join(path2.Join(path, "/DataTable/装备"), this.xlsxName))
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	tableRows := xlFile.GetRows(xlFile.GetSheetName(xlFile.GetActiveSheetIndex()))
	names := tableRows[1]
	rows := tableRows[4:]

	tmpIDMap := map[int32]*EquipRandomAttributePool{}

	for _,row := range rows{
		if row[0] != "" {
			e := &EquipRandomAttributePool{}
            e.ID = excel.ReadInt32(readCellValue(row, names, "ID"))
            e.ID_1 = excel.ReadInt32(readCellValue(row, names, "ID_1"))
            e.Weight_1 = excel.ReadInt32(readCellValue(row, names, "Weight_1"))
            e.ID_2 = excel.ReadInt32(readCellValue(row, names, "ID_2"))
            e.Weight_2 = excel.ReadInt32(readCellValue(row, names, "Weight_2"))
            e.ID_3 = excel.ReadInt32(readCellValue(row, names, "ID_3"))
            e.Weight_3 = excel.ReadInt32(readCellValue(row, names, "Weight_3"))
            e.ID_4 = excel.ReadInt32(readCellValue(row, names, "ID_4"))
            e.Weight_4 = excel.ReadInt32(readCellValue(row, names, "Weight_4"))
            e.ID_5 = excel.ReadInt32(readCellValue(row, names, "ID_5"))
            e.Weight_5 = excel.ReadInt32(readCellValue(row, names, "Weight_5"))
            e.ID_6 = excel.ReadInt32(readCellValue(row, names, "ID_6"))
            e.Weight_6 = excel.ReadInt32(readCellValue(row, names, "Weight_6"))
            e.ID_7 = excel.ReadInt32(readCellValue(row, names, "ID_7"))
            e.Weight_7 = excel.ReadInt32(readCellValue(row, names, "Weight_7"))
            e.ID_8 = excel.ReadInt32(readCellValue(row, names, "ID_8"))
            e.Weight_8 = excel.ReadInt32(readCellValue(row, names, "Weight_8"))
            e.ID_9 = excel.ReadInt32(readCellValue(row, names, "ID_9"))
            e.Weight_9 = excel.ReadInt32(readCellValue(row, names, "Weight_9"))
            e.ID_10 = excel.ReadInt32(readCellValue(row, names, "ID_10"))
            e.Weight_10 = excel.ReadInt32(readCellValue(row, names, "Weight_10"))
            e.Random = readRandom(row, names)

			tmpIDMap[excel.ReadInt32(readCellValue(row, names, "ID"))] = e

		}
	}
	this.indexID.Store(tmpIDMap)	

}
