package WeaponDecompose

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	path2 "path"
	"strings"
	"sync/atomic"
	"initialthree/node/table/excel"
	"initialthree/zaplogger"
)

type WeaponDecompose struct{
    ID int32 
    ExpToOneGoldRate float64 
    ExpToItemID int32 
    ReturnGoldCount int32 
    ReturnItemID_1 int32 
    ReturnItemCount_1 int32 
    ReturnItemID_2 int32 
    ReturnItemCount_2 int32 
    ReturnItemID_3 int32 
    ReturnItemCount_3 int32 
    ReturnItemID_4 int32 
    ReturnItemCount_4 int32 
    ReturnItemID_5 int32 
    ReturnItemCount_5 int32 
    ReturnItemID_6 int32 
    ReturnItemCount_6 int32 
    ReturnItemID_7 int32 
    ReturnItemCount_7 int32 
    ReturnItemID_8 int32 
    ReturnItemCount_8 int32 
    ReturnItemID_9 int32 
    ReturnItemCount_9 int32 
    ReturnItemID_10 int32 
    ReturnItemCount_10 int32 
    ReturnItemID_11 int32 
    ReturnItemCount_11 int32 
    ReturnItemID_12 int32 
    ReturnItemCount_12 int32 
    ReturnItemID_13 int32 
    ReturnItemCount_13 int32 
    ReturnItemID_14 int32 
    ReturnItemCount_14 int32 
    ReturnItemID_15 int32 
    ReturnItemCount_15 int32 
    ReturnItemID_16 int32 
    ReturnItemCount_16 int32 
    ReturnItemID_17 int32 
    ReturnItemCount_17 int32 
    ReturnItemID_18 int32 
    ReturnItemCount_18 int32 
    Items []*Items_ 

}

type Items_ struct{
    Id int32 
    Count int32 

}

func readItems(row, names []string)[]*Items_{
	ret := make([]*Items_, 0)
	base := excel.Split("ReturnItemID_,ReturnItemCount_", ",")

	for i := 1; ; i++ {
		if idx, find := hasCellName(names, fmt.Sprintf("%s%d", base[0][0], i));!find {
			break
		}else if _,ok := hasCellValue(row,idx);!ok{
			break
		}
		
		e := &Items_{}
        e.Id = excel.ReadInt32(readCellValue(row, names, fmt.Sprintf("%s%d", base[0][0], i)))
        e.Count = excel.ReadInt32(readCellValue(row, names, fmt.Sprintf("%s%d", base[1][0], i)))

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

func GetID(key int32)*WeaponDecompose{
	v := Table_.indexID.Load().(map[int32]*WeaponDecompose)[key]
	if v == nil{
		zaplogger.GetSugar().Errorf("Table %s GetID %d is nil", Table_.xlsxName, key)
	}
	return v
}

func GetIDMap() map[int32]*WeaponDecompose {
	return Table_.indexID.Load().(map[int32]*WeaponDecompose)
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
	Table_.xlsxName = "WeaponDecompose.xlsx"
	excel.AddTable(&Table_)
}

func (this *Table)Load(path string){
	xlFile, err := excelize.OpenFile(path2.Join(path2.Join(path, "/DataTable/武器"), this.xlsxName))
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	tableRows := xlFile.GetRows(xlFile.GetSheetName(xlFile.GetActiveSheetIndex()))
	names := tableRows[1]
	rows := tableRows[4:]

	tmpIDMap := map[int32]*WeaponDecompose{}

	for _,row := range rows{
		if row[0] != "" {
			e := &WeaponDecompose{}
            e.ID = excel.ReadInt32(readCellValue(row, names, "ID"))
            e.ExpToOneGoldRate = excel.ReadFloat(readCellValue(row, names, "ExpToOneGoldRate"))
            e.ExpToItemID = excel.ReadInt32(readCellValue(row, names, "ExpToItemID"))
            e.ReturnGoldCount = excel.ReadInt32(readCellValue(row, names, "ReturnGoldCount"))
            e.ReturnItemID_1 = excel.ReadInt32(readCellValue(row, names, "ReturnItemID_1"))
            e.ReturnItemCount_1 = excel.ReadInt32(readCellValue(row, names, "ReturnItemCount_1"))
            e.ReturnItemID_2 = excel.ReadInt32(readCellValue(row, names, "ReturnItemID_2"))
            e.ReturnItemCount_2 = excel.ReadInt32(readCellValue(row, names, "ReturnItemCount_2"))
            e.ReturnItemID_3 = excel.ReadInt32(readCellValue(row, names, "ReturnItemID_3"))
            e.ReturnItemCount_3 = excel.ReadInt32(readCellValue(row, names, "ReturnItemCount_3"))
            e.ReturnItemID_4 = excel.ReadInt32(readCellValue(row, names, "ReturnItemID_4"))
            e.ReturnItemCount_4 = excel.ReadInt32(readCellValue(row, names, "ReturnItemCount_4"))
            e.ReturnItemID_5 = excel.ReadInt32(readCellValue(row, names, "ReturnItemID_5"))
            e.ReturnItemCount_5 = excel.ReadInt32(readCellValue(row, names, "ReturnItemCount_5"))
            e.ReturnItemID_6 = excel.ReadInt32(readCellValue(row, names, "ReturnItemID_6"))
            e.ReturnItemCount_6 = excel.ReadInt32(readCellValue(row, names, "ReturnItemCount_6"))
            e.ReturnItemID_7 = excel.ReadInt32(readCellValue(row, names, "ReturnItemID_7"))
            e.ReturnItemCount_7 = excel.ReadInt32(readCellValue(row, names, "ReturnItemCount_7"))
            e.ReturnItemID_8 = excel.ReadInt32(readCellValue(row, names, "ReturnItemID_8"))
            e.ReturnItemCount_8 = excel.ReadInt32(readCellValue(row, names, "ReturnItemCount_8"))
            e.ReturnItemID_9 = excel.ReadInt32(readCellValue(row, names, "ReturnItemID_9"))
            e.ReturnItemCount_9 = excel.ReadInt32(readCellValue(row, names, "ReturnItemCount_9"))
            e.ReturnItemID_10 = excel.ReadInt32(readCellValue(row, names, "ReturnItemID_10"))
            e.ReturnItemCount_10 = excel.ReadInt32(readCellValue(row, names, "ReturnItemCount_10"))
            e.ReturnItemID_11 = excel.ReadInt32(readCellValue(row, names, "ReturnItemID_11"))
            e.ReturnItemCount_11 = excel.ReadInt32(readCellValue(row, names, "ReturnItemCount_11"))
            e.ReturnItemID_12 = excel.ReadInt32(readCellValue(row, names, "ReturnItemID_12"))
            e.ReturnItemCount_12 = excel.ReadInt32(readCellValue(row, names, "ReturnItemCount_12"))
            e.ReturnItemID_13 = excel.ReadInt32(readCellValue(row, names, "ReturnItemID_13"))
            e.ReturnItemCount_13 = excel.ReadInt32(readCellValue(row, names, "ReturnItemCount_13"))
            e.ReturnItemID_14 = excel.ReadInt32(readCellValue(row, names, "ReturnItemID_14"))
            e.ReturnItemCount_14 = excel.ReadInt32(readCellValue(row, names, "ReturnItemCount_14"))
            e.ReturnItemID_15 = excel.ReadInt32(readCellValue(row, names, "ReturnItemID_15"))
            e.ReturnItemCount_15 = excel.ReadInt32(readCellValue(row, names, "ReturnItemCount_15"))
            e.ReturnItemID_16 = excel.ReadInt32(readCellValue(row, names, "ReturnItemID_16"))
            e.ReturnItemCount_16 = excel.ReadInt32(readCellValue(row, names, "ReturnItemCount_16"))
            e.ReturnItemID_17 = excel.ReadInt32(readCellValue(row, names, "ReturnItemID_17"))
            e.ReturnItemCount_17 = excel.ReadInt32(readCellValue(row, names, "ReturnItemCount_17"))
            e.ReturnItemID_18 = excel.ReadInt32(readCellValue(row, names, "ReturnItemID_18"))
            e.ReturnItemCount_18 = excel.ReadInt32(readCellValue(row, names, "ReturnItemCount_18"))
            e.Items = readItems(row, names)

			tmpIDMap[excel.ReadInt32(readCellValue(row, names, "ID"))] = e

		}
	}
	this.indexID.Store(tmpIDMap)	

}
