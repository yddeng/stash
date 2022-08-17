package Equip

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	path2 "path"
	"strings"
	"sync/atomic"
	"initialthree/node/table/excel"
	"initialthree/zaplogger"
)

type Equip struct{
    ID int32 
    ExpToGoldRate float64 
    EquipSupplyExpItem string 
    DecomposeReturnEquipRate float64 
    MaxEquipWeight float64 
    EquipSupplyExpItemArray []*EquipSupplyExpItem_ 

}

type EquipSupplyExpItem_ struct{
    ItemID int32 
    Exp int32 

}

func readEquipSupplyExpItemArray(row, names []string)[]*EquipSupplyExpItem_{
	value := readCellValue(row, names, "EquipSupplyExpItem")
	if value == "" {
		return nil
	}

	r := excel.Split(value,",#")
	if len(r) == 0 {
		return nil
	}

	ret := make([]*EquipSupplyExpItem_,0)
	for _, v := range r {
		if len(v) == 2{
			e := &EquipSupplyExpItem_{}
	        e.ItemID = excel.ReadInt32(v[0])
        e.Exp = excel.ReadInt32(v[1])

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

func GetID(key int32)*Equip{
	v := Table_.indexID.Load().(map[int32]*Equip)[key]
	if v == nil{
		zaplogger.GetSugar().Errorf("Table %s GetID %d is nil", Table_.xlsxName, key)
	}
	return v
}

func GetIDMap() map[int32]*Equip {
	return Table_.indexID.Load().(map[int32]*Equip)
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
	Table_.xlsxName = "Equip.xlsx"
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

	tmpIDMap := map[int32]*Equip{}

	for _,row := range rows{
		if row[0] != "" {
			e := &Equip{}
            e.ID = excel.ReadInt32(readCellValue(row, names, "ID"))
            e.ExpToGoldRate = excel.ReadFloat(readCellValue(row, names, "ExpToGoldRate"))
            e.EquipSupplyExpItem = excel.ReadStr(readCellValue(row, names, "EquipSupplyExpItem"))
            e.DecomposeReturnEquipRate = excel.ReadFloat(readCellValue(row, names, "DecomposeReturnEquipRate"))
            e.MaxEquipWeight = excel.ReadFloat(readCellValue(row, names, "MaxEquipWeight"))
            e.EquipSupplyExpItemArray = readEquipSupplyExpItemArray(row, names)

			tmpIDMap[excel.ReadInt32(readCellValue(row, names, "ID"))] = e

		}
	}
	this.indexID.Store(tmpIDMap)	

}
