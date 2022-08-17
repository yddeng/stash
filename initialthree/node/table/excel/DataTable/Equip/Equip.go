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
    Name string 
    Vfx string 
    Icon string 
    Desc string 
    EquipType string 
    TendencyIcon string 
    Quality string 
    Cost float64 
    SupplyExp int32 
    LevelMaxExp int32 
    Decompose int32 
    Attrib_1 int32 
    Attrib_2 int32 
    Attrib_3 int32 
    Attrib_4 int32 
    Attrib_5 int32 
    SkillID int32 
    RandomAttribPool int32 
    EquipSkillCombatPower string 
    AccessWayList string 
    QualityEnum int32 
    AttrConfigs []*AttrConfigs_ 

}

type AttrConfigs_ struct{
    ID int32 

}

func readAttrConfigs(row, names []string)[]*AttrConfigs_{
	ret := make([]*AttrConfigs_, 0)
	base := excel.Split("Attrib_", ",")

	for i := 1; ; i++ {
		if idx, find := hasCellName(names, fmt.Sprintf("%s%d", base[0][0], i));!find {
			break
		}else if _,ok := hasCellValue(row,idx);!ok{
			break
		}
		
		e := &AttrConfigs_{}
        e.ID = excel.ReadInt32(readCellValue(row, names, fmt.Sprintf("%s%d", base[0][0], i)))

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
	xlFile, err := excelize.OpenFile(path2.Join(path2.Join(path, "/DataTable/装备"), this.xlsxName))
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
            e.Name = excel.ReadStr(readCellValue(row, names, "Name"))
            e.Vfx = excel.ReadStr(readCellValue(row, names, "Vfx"))
            e.Icon = excel.ReadStr(readCellValue(row, names, "Icon"))
            e.Desc = excel.ReadStr(readCellValue(row, names, "Desc"))
            e.EquipType = excel.ReadStr(readCellValue(row, names, "EquipType"))
            e.TendencyIcon = excel.ReadStr(readCellValue(row, names, "TendencyIcon"))
            e.Quality = excel.ReadStr(readCellValue(row, names, "Quality"))
            e.Cost = excel.ReadFloat(readCellValue(row, names, "Cost"))
            e.SupplyExp = excel.ReadInt32(readCellValue(row, names, "SupplyExp"))
            e.LevelMaxExp = excel.ReadInt32(readCellValue(row, names, "LevelMaxExp"))
            e.Decompose = excel.ReadInt32(readCellValue(row, names, "Decompose"))
            e.Attrib_1 = excel.ReadInt32(readCellValue(row, names, "Attrib_1"))
            e.Attrib_2 = excel.ReadInt32(readCellValue(row, names, "Attrib_2"))
            e.Attrib_3 = excel.ReadInt32(readCellValue(row, names, "Attrib_3"))
            e.Attrib_4 = excel.ReadInt32(readCellValue(row, names, "Attrib_4"))
            e.Attrib_5 = excel.ReadInt32(readCellValue(row, names, "Attrib_5"))
            e.SkillID = excel.ReadInt32(readCellValue(row, names, "SkillID"))
            e.RandomAttribPool = excel.ReadInt32(readCellValue(row, names, "RandomAttribPool"))
            e.EquipSkillCombatPower = excel.ReadStr(readCellValue(row, names, "EquipSkillCombatPower"))
            e.AccessWayList = excel.ReadStr(readCellValue(row, names, "AccessWayList"))
            e.QualityEnum = excel.ReadEnum(readCellValue(row, names, "Quality"))
            e.AttrConfigs = readAttrConfigs(row, names)

			tmpIDMap[excel.ReadInt32(readCellValue(row, names, "ID"))] = e

		}
	}
	this.indexID.Store(tmpIDMap)	

}
