package TalentGroup

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	path2 "path"
	"strings"
	"sync/atomic"
	"initialthree/node/table/excel"
	"initialthree/zaplogger"
)

type TalentGroup struct{
    ID int32 
    Name string 
    Talents string 
    Desc string 
    TalentsArray []*Talents_ 

}

type Talents_ struct{
    ID int32 

}

func readTalentsArray(row, names []string)[]*Talents_{
	value := readCellValue(row, names, "Talents")
	if value == "" {
		return nil
	}

	r := excel.Split(value,",")
	if len(r) == 0 {
		return nil
	}

	ret := make([]*Talents_,0)
	for _, v := range r {
		if len(v) == 1{
			e := &Talents_{}
	        e.ID = excel.ReadInt32(v[0])

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

func GetID(key int32)*TalentGroup{
	v := Table_.indexID.Load().(map[int32]*TalentGroup)[key]
	if v == nil{
		zaplogger.GetSugar().Errorf("Table %s GetID %d is nil", Table_.xlsxName, key)
	}
	return v
}

func GetIDMap() map[int32]*TalentGroup {
	return Table_.indexID.Load().(map[int32]*TalentGroup)
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
	Table_.xlsxName = "TalentGroup.xlsx"
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

	tmpIDMap := map[int32]*TalentGroup{}

	for _,row := range rows{
		if row[0] != "" {
			e := &TalentGroup{}
            e.ID = excel.ReadInt32(readCellValue(row, names, "ID"))
            e.Name = excel.ReadStr(readCellValue(row, names, "Name"))
            e.Talents = excel.ReadStr(readCellValue(row, names, "Talents"))
            e.Desc = excel.ReadStr(readCellValue(row, names, "Desc"))
            e.TalentsArray = readTalentsArray(row, names)

			tmpIDMap[excel.ReadInt32(readCellValue(row, names, "ID"))] = e

		}
	}
	this.indexID.Store(tmpIDMap)	

}
