package TalentLevel

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	path2 "path"
	"strings"
	"sync/atomic"
	"initialthree/node/table/excel"
	"initialthree/zaplogger"
)

type TalentLevel struct{
    ID int32 
    Values string 
    BehaviorID int32 
    GoldCost int32 
    ItemCost int32 
    RoleLevel int32 
    TotalTalentLLevel int32 
    PreTalentLevel string 
    PreTalentLevelArray []*PreTalentLevel_ 

}

type PreTalentLevel_ struct{
    ID int32 
    Level int32 

}

func readPreTalentLevelArray(row, names []string)[]*PreTalentLevel_{
	value := readCellValue(row, names, "PreTalentLevel")
	if value == "" {
		return nil
	}

	r := excel.Split(value,",#")
	if len(r) == 0 {
		return nil
	}

	ret := make([]*PreTalentLevel_,0)
	for _, v := range r {
		if len(v) == 2{
			e := &PreTalentLevel_{}
	        e.ID = excel.ReadInt32(v[0])
        e.Level = excel.ReadInt32(v[1])

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

func GetID(key int32)*TalentLevel{
	v := Table_.indexID.Load().(map[int32]*TalentLevel)[key]
	if v == nil{
		zaplogger.GetSugar().Errorf("Table %s GetID %d is nil", Table_.xlsxName, key)
	}
	return v
}

func GetIDMap() map[int32]*TalentLevel {
	return Table_.indexID.Load().(map[int32]*TalentLevel)
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
	Table_.xlsxName = "TalentLevel.xlsx"
	excel.AddTable(&Table_)
}

func (this *Table)Load(path string){
	xlFile, err := excelize.OpenFile(path2.Join(path2.Join(path, "/DataTable/天赋"), this.xlsxName))
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	tableRows := xlFile.GetRows(xlFile.GetSheetName(xlFile.GetActiveSheetIndex()))
	names := tableRows[1]
	rows := tableRows[4:]

	tmpIDMap := map[int32]*TalentLevel{}

	for _,row := range rows{
		if row[0] != "" {
			e := &TalentLevel{}
            e.ID = excel.ReadInt32(readCellValue(row, names, "ID"))
            e.Values = excel.ReadStr(readCellValue(row, names, "Values"))
            e.BehaviorID = excel.ReadInt32(readCellValue(row, names, "BehaviorID"))
            e.GoldCost = excel.ReadInt32(readCellValue(row, names, "GoldCost"))
            e.ItemCost = excel.ReadInt32(readCellValue(row, names, "ItemCost"))
            e.RoleLevel = excel.ReadInt32(readCellValue(row, names, "RoleLevel"))
            e.TotalTalentLLevel = excel.ReadInt32(readCellValue(row, names, "TotalTalentLLevel"))
            e.PreTalentLevel = excel.ReadStr(readCellValue(row, names, "PreTalentLevel"))
            e.PreTalentLevelArray = readPreTalentLevelArray(row, names)

			tmpIDMap[excel.ReadInt32(readCellValue(row, names, "ID"))] = e

		}
	}
	this.indexID.Store(tmpIDMap)	

}
