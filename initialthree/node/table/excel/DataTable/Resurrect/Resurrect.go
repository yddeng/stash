package Resurrect

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	path2 "path"
	"strings"
	"sync/atomic"
	"initialthree/node/table/excel"
	"initialthree/zaplogger"
)

type Resurrect struct{
    ID int32 
    Comment string 
    Effect int32 
    ResourceType_1 string 
    ResourceID_1 int32 
    ResourceCount_1 int32 
    ResourceType_2 string 
    ResourceID_2 int32 
    ResourceCount_2 int32 
    ResourceType_3 string 
    ResourceID_3 int32 
    ResourceCount_3 int32 
    ResourceType_4 string 
    ResourceID_4 int32 
    ResourceCount_4 int32 
    ResourceType_5 string 
    ResourceID_5 int32 
    ResourceCount_5 int32 
    ResourceType_6 string 
    ResourceID_6 int32 
    ResourceCount_6 int32 
    ResourceType_7 string 
    ResourceID_7 int32 
    ResourceCount_7 int32 
    ResourceType_8 string 
    ResourceID_8 int32 
    ResourceCount_8 int32 
    ResourceType_9 string 
    ResourceID_9 int32 
    ResourceCount_9 int32 
    ResourceType_10 string 
    ResourceID_10 int32 
    ResourceCount_10 int32 
    ResourceType_11 string 
    ResourceID_11 int32 
    ResourceCount_11 int32 
    ResourceType_12 string 
    ResourceID_12 int32 
    ResourceCount_12 int32 
    ResourceType_13 string 
    ResourceID_13 int32 
    ResourceCount_13 int32 
    ResourceType_14 string 
    ResourceID_14 int32 
    ResourceCount_14 int32 
    ResourceType_15 string 
    ResourceID_15 int32 
    ResourceCount_15 int32 
    ResourceType_16 string 
    ResourceID_16 int32 
    ResourceCount_16 int32 
    ResourceType_17 string 
    ResourceID_17 int32 
    ResourceCount_17 int32 
    ResourceType_18 string 
    ResourceID_18 int32 
    ResourceCount_18 int32 
    ResourceType_19 string 
    ResourceID_19 int32 
    ResourceCount_19 int32 
    ResourceType_20 string 
    ResourceID_20 int32 
    ResourceCount_20 int32 
    Resource []*Resource_ 

}

type Resource_ struct{
    ResourceType string 
    ResourceID int32 
    ResourceCount int32 

}

func readResource(row, names []string)[]*Resource_{
	ret := make([]*Resource_, 0)
	base := excel.Split("ResourceType_,ResourceID_,ResourceCount_", ",")

	for i := 1; ; i++ {
		if idx, find := hasCellName(names, fmt.Sprintf("%s%d", base[0][0], i));!find {
			break
		}else if _,ok := hasCellValue(row,idx);!ok{
			break
		}
		
		e := &Resource_{}
        e.ResourceType = excel.ReadStr(readCellValue(row, names, fmt.Sprintf("%s%d", base[0][0], i)))
        e.ResourceID = excel.ReadInt32(readCellValue(row, names, fmt.Sprintf("%s%d", base[1][0], i)))
        e.ResourceCount = excel.ReadInt32(readCellValue(row, names, fmt.Sprintf("%s%d", base[2][0], i)))

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

func GetID(key int32)*Resurrect{
	v := Table_.indexID.Load().(map[int32]*Resurrect)[key]
	if v == nil{
		zaplogger.GetSugar().Errorf("Table %s GetID %d is nil", Table_.xlsxName, key)
	}
	return v
}

func GetIDMap() map[int32]*Resurrect {
	return Table_.indexID.Load().(map[int32]*Resurrect)
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
	Table_.xlsxName = "Resurrect.xlsx"
	excel.AddTable(&Table_)
}

func (this *Table)Load(path string){
	xlFile, err := excelize.OpenFile(path2.Join(path2.Join(path, "/DataTable/副本（编辑器内关卡）"), this.xlsxName))
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	tableRows := xlFile.GetRows(xlFile.GetSheetName(xlFile.GetActiveSheetIndex()))
	names := tableRows[1]
	rows := tableRows[4:]

	tmpIDMap := map[int32]*Resurrect{}

	for _,row := range rows{
		if row[0] != "" {
			e := &Resurrect{}
            e.ID = excel.ReadInt32(readCellValue(row, names, "ID"))
            e.Comment = excel.ReadStr(readCellValue(row, names, "Comment"))
            e.Effect = excel.ReadInt32(readCellValue(row, names, "Effect"))
            e.ResourceType_1 = excel.ReadStr(readCellValue(row, names, "ResourceType_1"))
            e.ResourceID_1 = excel.ReadInt32(readCellValue(row, names, "ResourceID_1"))
            e.ResourceCount_1 = excel.ReadInt32(readCellValue(row, names, "ResourceCount_1"))
            e.ResourceType_2 = excel.ReadStr(readCellValue(row, names, "ResourceType_2"))
            e.ResourceID_2 = excel.ReadInt32(readCellValue(row, names, "ResourceID_2"))
            e.ResourceCount_2 = excel.ReadInt32(readCellValue(row, names, "ResourceCount_2"))
            e.ResourceType_3 = excel.ReadStr(readCellValue(row, names, "ResourceType_3"))
            e.ResourceID_3 = excel.ReadInt32(readCellValue(row, names, "ResourceID_3"))
            e.ResourceCount_3 = excel.ReadInt32(readCellValue(row, names, "ResourceCount_3"))
            e.ResourceType_4 = excel.ReadStr(readCellValue(row, names, "ResourceType_4"))
            e.ResourceID_4 = excel.ReadInt32(readCellValue(row, names, "ResourceID_4"))
            e.ResourceCount_4 = excel.ReadInt32(readCellValue(row, names, "ResourceCount_4"))
            e.ResourceType_5 = excel.ReadStr(readCellValue(row, names, "ResourceType_5"))
            e.ResourceID_5 = excel.ReadInt32(readCellValue(row, names, "ResourceID_5"))
            e.ResourceCount_5 = excel.ReadInt32(readCellValue(row, names, "ResourceCount_5"))
            e.ResourceType_6 = excel.ReadStr(readCellValue(row, names, "ResourceType_6"))
            e.ResourceID_6 = excel.ReadInt32(readCellValue(row, names, "ResourceID_6"))
            e.ResourceCount_6 = excel.ReadInt32(readCellValue(row, names, "ResourceCount_6"))
            e.ResourceType_7 = excel.ReadStr(readCellValue(row, names, "ResourceType_7"))
            e.ResourceID_7 = excel.ReadInt32(readCellValue(row, names, "ResourceID_7"))
            e.ResourceCount_7 = excel.ReadInt32(readCellValue(row, names, "ResourceCount_7"))
            e.ResourceType_8 = excel.ReadStr(readCellValue(row, names, "ResourceType_8"))
            e.ResourceID_8 = excel.ReadInt32(readCellValue(row, names, "ResourceID_8"))
            e.ResourceCount_8 = excel.ReadInt32(readCellValue(row, names, "ResourceCount_8"))
            e.ResourceType_9 = excel.ReadStr(readCellValue(row, names, "ResourceType_9"))
            e.ResourceID_9 = excel.ReadInt32(readCellValue(row, names, "ResourceID_9"))
            e.ResourceCount_9 = excel.ReadInt32(readCellValue(row, names, "ResourceCount_9"))
            e.ResourceType_10 = excel.ReadStr(readCellValue(row, names, "ResourceType_10"))
            e.ResourceID_10 = excel.ReadInt32(readCellValue(row, names, "ResourceID_10"))
            e.ResourceCount_10 = excel.ReadInt32(readCellValue(row, names, "ResourceCount_10"))
            e.ResourceType_11 = excel.ReadStr(readCellValue(row, names, "ResourceType_11"))
            e.ResourceID_11 = excel.ReadInt32(readCellValue(row, names, "ResourceID_11"))
            e.ResourceCount_11 = excel.ReadInt32(readCellValue(row, names, "ResourceCount_11"))
            e.ResourceType_12 = excel.ReadStr(readCellValue(row, names, "ResourceType_12"))
            e.ResourceID_12 = excel.ReadInt32(readCellValue(row, names, "ResourceID_12"))
            e.ResourceCount_12 = excel.ReadInt32(readCellValue(row, names, "ResourceCount_12"))
            e.ResourceType_13 = excel.ReadStr(readCellValue(row, names, "ResourceType_13"))
            e.ResourceID_13 = excel.ReadInt32(readCellValue(row, names, "ResourceID_13"))
            e.ResourceCount_13 = excel.ReadInt32(readCellValue(row, names, "ResourceCount_13"))
            e.ResourceType_14 = excel.ReadStr(readCellValue(row, names, "ResourceType_14"))
            e.ResourceID_14 = excel.ReadInt32(readCellValue(row, names, "ResourceID_14"))
            e.ResourceCount_14 = excel.ReadInt32(readCellValue(row, names, "ResourceCount_14"))
            e.ResourceType_15 = excel.ReadStr(readCellValue(row, names, "ResourceType_15"))
            e.ResourceID_15 = excel.ReadInt32(readCellValue(row, names, "ResourceID_15"))
            e.ResourceCount_15 = excel.ReadInt32(readCellValue(row, names, "ResourceCount_15"))
            e.ResourceType_16 = excel.ReadStr(readCellValue(row, names, "ResourceType_16"))
            e.ResourceID_16 = excel.ReadInt32(readCellValue(row, names, "ResourceID_16"))
            e.ResourceCount_16 = excel.ReadInt32(readCellValue(row, names, "ResourceCount_16"))
            e.ResourceType_17 = excel.ReadStr(readCellValue(row, names, "ResourceType_17"))
            e.ResourceID_17 = excel.ReadInt32(readCellValue(row, names, "ResourceID_17"))
            e.ResourceCount_17 = excel.ReadInt32(readCellValue(row, names, "ResourceCount_17"))
            e.ResourceType_18 = excel.ReadStr(readCellValue(row, names, "ResourceType_18"))
            e.ResourceID_18 = excel.ReadInt32(readCellValue(row, names, "ResourceID_18"))
            e.ResourceCount_18 = excel.ReadInt32(readCellValue(row, names, "ResourceCount_18"))
            e.ResourceType_19 = excel.ReadStr(readCellValue(row, names, "ResourceType_19"))
            e.ResourceID_19 = excel.ReadInt32(readCellValue(row, names, "ResourceID_19"))
            e.ResourceCount_19 = excel.ReadInt32(readCellValue(row, names, "ResourceCount_19"))
            e.ResourceType_20 = excel.ReadStr(readCellValue(row, names, "ResourceType_20"))
            e.ResourceID_20 = excel.ReadInt32(readCellValue(row, names, "ResourceID_20"))
            e.ResourceCount_20 = excel.ReadInt32(readCellValue(row, names, "ResourceCount_20"))
            e.Resource = readResource(row, names)

			tmpIDMap[excel.ReadInt32(readCellValue(row, names, "ID"))] = e

		}
	}
	this.indexID.Store(tmpIDMap)	

}
