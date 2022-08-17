package CharacterLevelUpAttribute

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	path2 "path"
	"strings"
	"sync/atomic"
	"initialthree/node/table/excel"
	"initialthree/zaplogger"
)

type CharacterLevelUpAttribute struct{
    ID int32 
    CharacterInfo string 
    AttributeList string 
    AttributeListArray []*AttributeList_ 

}

type AttributeList_ struct{
    ID int32 
    Val float64 

}

func readAttributeListArray(row, names []string)[]*AttributeList_{
	value := readCellValue(row, names, "AttributeList")
	if value == "" {
		return nil
	}

	r := excel.Split(value,",#")
	if len(r) == 0 {
		return nil
	}

	ret := make([]*AttributeList_,0)
	for _, v := range r {
		if len(v) == 2{
			e := &AttributeList_{}
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

func GetID(key int32)*CharacterLevelUpAttribute{
	v := Table_.indexID.Load().(map[int32]*CharacterLevelUpAttribute)[key]
	if v == nil{
		zaplogger.GetSugar().Errorf("Table %s GetID %d is nil", Table_.xlsxName, key)
	}
	return v
}

func GetIDMap() map[int32]*CharacterLevelUpAttribute {
	return Table_.indexID.Load().(map[int32]*CharacterLevelUpAttribute)
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
	Table_.xlsxName = "CharacterLevelUpAttribute.xlsx"
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

	tmpIDMap := map[int32]*CharacterLevelUpAttribute{}

	for _,row := range rows{
		if row[0] != "" {
			e := &CharacterLevelUpAttribute{}
            e.ID = excel.ReadInt32(readCellValue(row, names, "ID"))
            e.CharacterInfo = excel.ReadStr(readCellValue(row, names, "CharacterInfo"))
            e.AttributeList = excel.ReadStr(readCellValue(row, names, "AttributeList"))
            e.AttributeListArray = readAttributeListArray(row, names)

			tmpIDMap[excel.ReadInt32(readCellValue(row, names, "ID"))] = e

		}
	}
	this.indexID.Store(tmpIDMap)	

}
