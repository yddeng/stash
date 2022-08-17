package Weapon

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	path2 "path"
	"strings"
	"sync/atomic"
	"initialthree/node/table/excel"
	"initialthree/zaplogger"
)

type Weapon struct{
    ID int32 
    ExpToGoldRate float64 
    WeaponSupplyExpItem string 
    DecomposeReturnEquipRate float64 
    WeaponDisplayPrefab_1 string 
    WeaponDisplayPrefab_2 string 
    WeaponDisplayPrefab_3 string 
    WeaponDisplayPrefab_4 string 
    WeaponDisplayPrefab_5 string 
    WeaponDisplayPrefab_6 string 
    WeaponDisplayPrefab_7 string 
    WeaponDisplayPrefab_8 string 
    WeaponDisplayPrefab_9 string 
    WeaponDisplayPrefab_10 string 
    WeaponDisplayPrefab_11 string 
    WeaponDisplayPrefab_12 string 
    WeaponSupplyExpItemArray []*WeaponSupplyExpItem_ 

}

type WeaponSupplyExpItem_ struct{
    ItemID int32 
    Exp int32 

}

func readWeaponSupplyExpItemArray(row, names []string)[]*WeaponSupplyExpItem_{
	value := readCellValue(row, names, "WeaponSupplyExpItem")
	if value == "" {
		return nil
	}

	r := excel.Split(value,",#")
	if len(r) == 0 {
		return nil
	}

	ret := make([]*WeaponSupplyExpItem_,0)
	for _, v := range r {
		if len(v) == 2{
			e := &WeaponSupplyExpItem_{}
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

func GetID(key int32)*Weapon{
	v := Table_.indexID.Load().(map[int32]*Weapon)[key]
	if v == nil{
		zaplogger.GetSugar().Errorf("Table %s GetID %d is nil", Table_.xlsxName, key)
	}
	return v
}

func GetIDMap() map[int32]*Weapon {
	return Table_.indexID.Load().(map[int32]*Weapon)
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
	Table_.xlsxName = "Weapon.xlsx"
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

	tmpIDMap := map[int32]*Weapon{}

	for _,row := range rows{
		if row[0] != "" {
			e := &Weapon{}
            e.ID = excel.ReadInt32(readCellValue(row, names, "ID"))
            e.ExpToGoldRate = excel.ReadFloat(readCellValue(row, names, "ExpToGoldRate"))
            e.WeaponSupplyExpItem = excel.ReadStr(readCellValue(row, names, "WeaponSupplyExpItem"))
            e.DecomposeReturnEquipRate = excel.ReadFloat(readCellValue(row, names, "DecomposeReturnEquipRate"))
            e.WeaponDisplayPrefab_1 = excel.ReadStr(readCellValue(row, names, "WeaponDisplayPrefab_1"))
            e.WeaponDisplayPrefab_2 = excel.ReadStr(readCellValue(row, names, "WeaponDisplayPrefab_2"))
            e.WeaponDisplayPrefab_3 = excel.ReadStr(readCellValue(row, names, "WeaponDisplayPrefab_3"))
            e.WeaponDisplayPrefab_4 = excel.ReadStr(readCellValue(row, names, "WeaponDisplayPrefab_4"))
            e.WeaponDisplayPrefab_5 = excel.ReadStr(readCellValue(row, names, "WeaponDisplayPrefab_5"))
            e.WeaponDisplayPrefab_6 = excel.ReadStr(readCellValue(row, names, "WeaponDisplayPrefab_6"))
            e.WeaponDisplayPrefab_7 = excel.ReadStr(readCellValue(row, names, "WeaponDisplayPrefab_7"))
            e.WeaponDisplayPrefab_8 = excel.ReadStr(readCellValue(row, names, "WeaponDisplayPrefab_8"))
            e.WeaponDisplayPrefab_9 = excel.ReadStr(readCellValue(row, names, "WeaponDisplayPrefab_9"))
            e.WeaponDisplayPrefab_10 = excel.ReadStr(readCellValue(row, names, "WeaponDisplayPrefab_10"))
            e.WeaponDisplayPrefab_11 = excel.ReadStr(readCellValue(row, names, "WeaponDisplayPrefab_11"))
            e.WeaponDisplayPrefab_12 = excel.ReadStr(readCellValue(row, names, "WeaponDisplayPrefab_12"))
            e.WeaponSupplyExpItemArray = readWeaponSupplyExpItemArray(row, names)

			tmpIDMap[excel.ReadInt32(readCellValue(row, names, "ID"))] = e

		}
	}
	this.indexID.Store(tmpIDMap)	

}
