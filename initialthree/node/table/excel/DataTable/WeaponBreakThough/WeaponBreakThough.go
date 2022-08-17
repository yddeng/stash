package WeaponBreakThough

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	path2 "path"
	"strings"
	"sync/atomic"
	"initialthree/node/table/excel"
	"initialthree/zaplogger"
)

type WeaponBreakThough struct{
    ID int32 
    Level_1 int32 
    GoldCost_1 int32 
    ItemCost_1 string 
    Attribute_1 string 
    Level_2 int32 
    GoldCost_2 int32 
    ItemCost_2 string 
    Attribute_2 string 
    Level_3 int32 
    GoldCost_3 int32 
    ItemCost_3 string 
    Attribute_3 string 
    Level_4 int32 
    GoldCost_4 int32 
    ItemCost_4 string 
    Attribute_4 string 
    Level_5 int32 
    GoldCost_5 int32 
    ItemCost_5 string 
    Attribute_5 string 
    Level_6 int32 
    GoldCost_6 int32 
    ItemCost_6 string 
    Attribute_6 string 
    Level_7 int32 
    GoldCost_7 int32 
    ItemCost_7 string 
    Attribute_7 string 
    BreakLevel []*BreakLevel_ 

}

type BreakLevel_ struct{
    Level int32 
    Gold int32 
    ItemStr string 
    AttrStr string 

}

func readBreakLevel(row, names []string)[]*BreakLevel_{
	ret := make([]*BreakLevel_, 0)
	base := excel.Split("Level_,GoldCost_,ItemCost_,Attribute_", ",")

	for i := 1; ; i++ {
		if idx, find := hasCellName(names, fmt.Sprintf("%s%d", base[0][0], i));!find {
			break
		}else if _,ok := hasCellValue(row,idx);!ok{
			break
		}
		
		e := &BreakLevel_{}
        e.Level = excel.ReadInt32(readCellValue(row, names, fmt.Sprintf("%s%d", base[0][0], i)))
        e.Gold = excel.ReadInt32(readCellValue(row, names, fmt.Sprintf("%s%d", base[1][0], i)))
        e.ItemStr = excel.ReadStr(readCellValue(row, names, fmt.Sprintf("%s%d", base[2][0], i)))
        e.AttrStr = excel.ReadStr(readCellValue(row, names, fmt.Sprintf("%s%d", base[3][0], i)))

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

func GetID(key int32)*WeaponBreakThough{
	v := Table_.indexID.Load().(map[int32]*WeaponBreakThough)[key]
	if v == nil{
		zaplogger.GetSugar().Errorf("Table %s GetID %d is nil", Table_.xlsxName, key)
	}
	return v
}

func GetIDMap() map[int32]*WeaponBreakThough {
	return Table_.indexID.Load().(map[int32]*WeaponBreakThough)
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
	Table_.xlsxName = "WeaponBreakThough.xlsx"
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

	tmpIDMap := map[int32]*WeaponBreakThough{}

	for _,row := range rows{
		if row[0] != "" {
			e := &WeaponBreakThough{}
            e.ID = excel.ReadInt32(readCellValue(row, names, "ID"))
            e.Level_1 = excel.ReadInt32(readCellValue(row, names, "Level_1"))
            e.GoldCost_1 = excel.ReadInt32(readCellValue(row, names, "GoldCost_1"))
            e.ItemCost_1 = excel.ReadStr(readCellValue(row, names, "ItemCost_1"))
            e.Attribute_1 = excel.ReadStr(readCellValue(row, names, "Attribute_1"))
            e.Level_2 = excel.ReadInt32(readCellValue(row, names, "Level_2"))
            e.GoldCost_2 = excel.ReadInt32(readCellValue(row, names, "GoldCost_2"))
            e.ItemCost_2 = excel.ReadStr(readCellValue(row, names, "ItemCost_2"))
            e.Attribute_2 = excel.ReadStr(readCellValue(row, names, "Attribute_2"))
            e.Level_3 = excel.ReadInt32(readCellValue(row, names, "Level_3"))
            e.GoldCost_3 = excel.ReadInt32(readCellValue(row, names, "GoldCost_3"))
            e.ItemCost_3 = excel.ReadStr(readCellValue(row, names, "ItemCost_3"))
            e.Attribute_3 = excel.ReadStr(readCellValue(row, names, "Attribute_3"))
            e.Level_4 = excel.ReadInt32(readCellValue(row, names, "Level_4"))
            e.GoldCost_4 = excel.ReadInt32(readCellValue(row, names, "GoldCost_4"))
            e.ItemCost_4 = excel.ReadStr(readCellValue(row, names, "ItemCost_4"))
            e.Attribute_4 = excel.ReadStr(readCellValue(row, names, "Attribute_4"))
            e.Level_5 = excel.ReadInt32(readCellValue(row, names, "Level_5"))
            e.GoldCost_5 = excel.ReadInt32(readCellValue(row, names, "GoldCost_5"))
            e.ItemCost_5 = excel.ReadStr(readCellValue(row, names, "ItemCost_5"))
            e.Attribute_5 = excel.ReadStr(readCellValue(row, names, "Attribute_5"))
            e.Level_6 = excel.ReadInt32(readCellValue(row, names, "Level_6"))
            e.GoldCost_6 = excel.ReadInt32(readCellValue(row, names, "GoldCost_6"))
            e.ItemCost_6 = excel.ReadStr(readCellValue(row, names, "ItemCost_6"))
            e.Attribute_6 = excel.ReadStr(readCellValue(row, names, "Attribute_6"))
            e.Level_7 = excel.ReadInt32(readCellValue(row, names, "Level_7"))
            e.GoldCost_7 = excel.ReadInt32(readCellValue(row, names, "GoldCost_7"))
            e.ItemCost_7 = excel.ReadStr(readCellValue(row, names, "ItemCost_7"))
            e.Attribute_7 = excel.ReadStr(readCellValue(row, names, "Attribute_7"))
            e.BreakLevel = readBreakLevel(row, names)

			tmpIDMap[excel.ReadInt32(readCellValue(row, names, "ID"))] = e

		}
	}
	this.indexID.Store(tmpIDMap)	

}
