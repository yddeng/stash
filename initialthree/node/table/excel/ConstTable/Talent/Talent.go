package Talent

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	path2 "path"
	"strings"
	"sync/atomic"
	"initialthree/node/table/excel"
	"initialthree/zaplogger"
)

type Talent struct{
    ID int32 
    ItemID int32 
    ResetGoldCost int32 
    ResetItemCost int32 
    InfiniteTalentGoldCost int32 
    InfiniteTalentItemCost int32 
    InfiniteTalentAttr_1 string 
    InfiniteTalentAttrValue_1 string 
    InfiniteTalentAttr_2 string 
    InfiniteTalentAttrValue_2 string 
    InfiniteTalentAttr_3 string 
    InfiniteTalentAttrValue_3 string 
    Attr []*Attr_ 

}

type Attr_ struct{
    Attr int32 
    Value int32 

}

func readAttr(row, names []string)[]*Attr_{
	ret := make([]*Attr_, 0)
	base := excel.Split("InfiniteTatlentAttr_,InfiniteTalentAttrValue_", ",")

	for i := 1; ; i++ {
		if idx, find := hasCellName(names, fmt.Sprintf("%s%d", base[0][0], i));!find {
			break
		}else if _,ok := hasCellValue(row,idx);!ok{
			break
		}
		
		e := &Attr_{}
        e.Attr = excel.ReadInt32(readCellValue(row, names, fmt.Sprintf("%s%d", base[0][0], i)))
        e.Value = excel.ReadInt32(readCellValue(row, names, fmt.Sprintf("%s%d", base[1][0], i)))

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

func GetID(key int32)*Talent{
	v := Table_.indexID.Load().(map[int32]*Talent)[key]
	if v == nil{
		zaplogger.GetSugar().Errorf("Table %s GetID %d is nil", Table_.xlsxName, key)
	}
	return v
}

func GetIDMap() map[int32]*Talent {
	return Table_.indexID.Load().(map[int32]*Talent)
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
	Table_.xlsxName = "Talent.xlsx"
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

	tmpIDMap := map[int32]*Talent{}

	for _,row := range rows{
		if row[0] != "" {
			e := &Talent{}
            e.ID = excel.ReadInt32(readCellValue(row, names, "ID"))
            e.ItemID = excel.ReadInt32(readCellValue(row, names, "ItemID"))
            e.ResetGoldCost = excel.ReadInt32(readCellValue(row, names, "ResetGoldCost"))
            e.ResetItemCost = excel.ReadInt32(readCellValue(row, names, "ResetItemCost"))
            e.InfiniteTalentGoldCost = excel.ReadInt32(readCellValue(row, names, "InfiniteTalentGoldCost"))
            e.InfiniteTalentItemCost = excel.ReadInt32(readCellValue(row, names, "InfiniteTalentItemCost"))
            e.InfiniteTalentAttr_1 = excel.ReadStr(readCellValue(row, names, "InfiniteTalentAttr_1"))
            e.InfiniteTalentAttrValue_1 = excel.ReadStr(readCellValue(row, names, "InfiniteTalentAttrValue_1"))
            e.InfiniteTalentAttr_2 = excel.ReadStr(readCellValue(row, names, "InfiniteTalentAttr_2"))
            e.InfiniteTalentAttrValue_2 = excel.ReadStr(readCellValue(row, names, "InfiniteTalentAttrValue_2"))
            e.InfiniteTalentAttr_3 = excel.ReadStr(readCellValue(row, names, "InfiniteTalentAttr_3"))
            e.InfiniteTalentAttrValue_3 = excel.ReadStr(readCellValue(row, names, "InfiniteTalentAttrValue_3"))
            e.Attr = readAttr(row, names)

			tmpIDMap[excel.ReadInt32(readCellValue(row, names, "ID"))] = e

		}
	}
	this.indexID.Store(tmpIDMap)	

}
