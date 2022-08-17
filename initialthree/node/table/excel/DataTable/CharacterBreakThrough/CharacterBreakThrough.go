package CharacterBreakThrough

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	path2 "path"
	"strings"
	"sync/atomic"
	"initialthree/node/table/excel"
	"initialthree/zaplogger"
)

type CharacterBreakThrough struct{
    ID int32 
    GoldCost int32 
    SpecifiedIDList string 
    LevelRequirement int32 
    AttributeBonus string 
    AccountLevelRequirement int32 
    SpecifiedIDListArray []*SpecifiedIDList_ 
    AttributeBonusArray []*AttributeBonus_ 

}

type SpecifiedIDList_ struct{
    ID int32 
    Count int32 

}

func readSpecifiedIDListArray(row, names []string)[]*SpecifiedIDList_{
	value := readCellValue(row, names, "SpecifiedIDList")
	if value == "" {
		return nil
	}

	r := excel.Split(value,",#")
	if len(r) == 0 {
		return nil
	}

	ret := make([]*SpecifiedIDList_,0)
	for _, v := range r {
		if len(v) == 2{
			e := &SpecifiedIDList_{}
	        e.ID = excel.ReadInt32(v[0])
        e.Count = excel.ReadInt32(v[1])

			ret = append(ret, e)
		}
	}

	return ret
}

type AttributeBonus_ struct{
    ID int32 
    Val float64 

}

func readAttributeBonusArray(row, names []string)[]*AttributeBonus_{
	value := readCellValue(row, names, "AttributeBonus")
	if value == "" {
		return nil
	}

	r := excel.Split(value,",#")
	if len(r) == 0 {
		return nil
	}

	ret := make([]*AttributeBonus_,0)
	for _, v := range r {
		if len(v) == 2{
			e := &AttributeBonus_{}
	        e.ID = excel.ReadInt32(v[0])
        e.Val = excel.ReadFloat(v[1])

			ret = append(ret, e)
		}
	}

	return ret
}

var Table_ Table

type Table struct{
	xlsxName string
	nextPath string
	indexID  atomic.Value

}

func GetID(key int32)*CharacterBreakThrough{
	v := Table_.indexID.Load().(map[int32]*CharacterBreakThrough)[key]
	if v == nil{
		zaplogger.GetSugar().Errorf("Table %s GetID %d is nil", Table_.xlsxName, key)
	}
	return v
}

func GetIDMap() map[int32]*CharacterBreakThrough {
	return Table_.indexID.Load().(map[int32]*CharacterBreakThrough)
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
	Table_.xlsxName = "CharacterBreakThrough.xlsx"
	excel.AddTable(&Table_)
}

func (this *Table)Load(path string){
	xlFile, err := excelize.OpenFile(path2.Join(path2.Join(path, "/DataTable/角色"), this.xlsxName))
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	tableRows := xlFile.GetRows(xlFile.GetSheetName(xlFile.GetActiveSheetIndex()))
	names := tableRows[1]
	rows := tableRows[4:]

	tmpIDMap := map[int32]*CharacterBreakThrough{}

	for _,row := range rows{
		if row[0] != "" {
			e := &CharacterBreakThrough{}
            e.ID = excel.ReadInt32(readCellValue(row, names, "ID"))
            e.GoldCost = excel.ReadInt32(readCellValue(row, names, "GoldCost"))
            e.SpecifiedIDList = excel.ReadStr(readCellValue(row, names, "SpecifiedIDList"))
            e.LevelRequirement = excel.ReadInt32(readCellValue(row, names, "LevelRequirement"))
            e.AttributeBonus = excel.ReadStr(readCellValue(row, names, "AttributeBonus"))
            e.AccountLevelRequirement = excel.ReadInt32(readCellValue(row, names, "AccountLevelRequirement"))
            e.SpecifiedIDListArray = readSpecifiedIDListArray(row, names)
            e.AttributeBonusArray = readAttributeBonusArray(row, names)

			tmpIDMap[excel.ReadInt32(readCellValue(row, names, "ID"))] = e

		}
	}
	this.indexID.Store(tmpIDMap)	

}
