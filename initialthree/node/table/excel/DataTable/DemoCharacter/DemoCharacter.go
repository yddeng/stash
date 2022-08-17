package DemoCharacter

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	path2 "path"
	"strings"
	"sync/atomic"
	"initialthree/node/table/excel"
	"initialthree/zaplogger"
)

type DemoCharacter struct{
    ID int32 
    Comment string 
    CharacterID int32 
    CharacterLv int32 
    CharacterBreakLv int32 
    WeaponID int32 
    WeaponLv int32 
    WeaponBreakLv int32 
    WeaponRefineLv int32 
    GeneLevel int32 
    SkillLevel_1 int32 
    SkillLevel_2 int32 
    SkillLevel_3 int32 
    SkillLevel_4 int32 
    SkillLevel_5 int32 
    SkillLevel_6 int32 
    EquipID_1 int32 
    EquipLv_1 int32 
    EquipSkillLv_1 int32 
    EquipRandomAttrib_1 int32 
    EquipRandomAttribLv_1 int32 
    EquipID_2 int32 
    EquipLv_2 int32 
    EquipSkillLv_2 int32 
    EquipRandomAttrib_2 int32 
    EquipRandomAttribLv_2 int32 
    EquipID_3 int32 
    EquipLv_3 int32 
    EquipSkillLv_3 int32 
    EquipRandomAttrib_3 int32 
    EquipRandomAttribLv_3 int32 
    EquipID_4 int32 
    EquipLv_4 int32 
    EquipSkillLv_4 int32 
    EquipRandomAttrib_4 int32 
    EquipRandomAttribLv_4 int32 
    EquipID_5 int32 
    EquipLv_5 int32 
    EquipSkillLv_5 int32 
    EquipRandomAttrib_5 int32 
    EquipRandomAttribLv_5 int32 

}

var Table_ Table

type Table struct{
	xlsxName string
	nextPath string
	indexID  atomic.Value

}

func GetID(key int32)*DemoCharacter{
	v := Table_.indexID.Load().(map[int32]*DemoCharacter)[key]
	if v == nil{
		zaplogger.GetSugar().Errorf("Table %s GetID %d is nil", Table_.xlsxName, key)
	}
	return v
}

func GetIDMap() map[int32]*DemoCharacter {
	return Table_.indexID.Load().(map[int32]*DemoCharacter)
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
	Table_.xlsxName = "DemoCharacter.xlsx"
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

	tmpIDMap := map[int32]*DemoCharacter{}

	for _,row := range rows{
		if row[0] != "" {
			e := &DemoCharacter{}
            e.ID = excel.ReadInt32(readCellValue(row, names, "ID"))
            e.Comment = excel.ReadStr(readCellValue(row, names, "Comment"))
            e.CharacterID = excel.ReadInt32(readCellValue(row, names, "CharacterID"))
            e.CharacterLv = excel.ReadInt32(readCellValue(row, names, "CharacterLv"))
            e.CharacterBreakLv = excel.ReadInt32(readCellValue(row, names, "CharacterBreakLv"))
            e.WeaponID = excel.ReadInt32(readCellValue(row, names, "WeaponID"))
            e.WeaponLv = excel.ReadInt32(readCellValue(row, names, "WeaponLv"))
            e.WeaponBreakLv = excel.ReadInt32(readCellValue(row, names, "WeaponBreakLv"))
            e.WeaponRefineLv = excel.ReadInt32(readCellValue(row, names, "WeaponRefineLv"))
            e.GeneLevel = excel.ReadInt32(readCellValue(row, names, "GeneLevel"))
            e.SkillLevel_1 = excel.ReadInt32(readCellValue(row, names, "SkillLevel_1"))
            e.SkillLevel_2 = excel.ReadInt32(readCellValue(row, names, "SkillLevel_2"))
            e.SkillLevel_3 = excel.ReadInt32(readCellValue(row, names, "SkillLevel_3"))
            e.SkillLevel_4 = excel.ReadInt32(readCellValue(row, names, "SkillLevel_4"))
            e.SkillLevel_5 = excel.ReadInt32(readCellValue(row, names, "SkillLevel_5"))
            e.SkillLevel_6 = excel.ReadInt32(readCellValue(row, names, "SkillLevel_6"))
            e.EquipID_1 = excel.ReadInt32(readCellValue(row, names, "EquipID_1"))
            e.EquipLv_1 = excel.ReadInt32(readCellValue(row, names, "EquipLv_1"))
            e.EquipSkillLv_1 = excel.ReadInt32(readCellValue(row, names, "EquipSkillLv_1"))
            e.EquipRandomAttrib_1 = excel.ReadInt32(readCellValue(row, names, "EquipRandomAttrib_1"))
            e.EquipRandomAttribLv_1 = excel.ReadInt32(readCellValue(row, names, "EquipRandomAttribLv_1"))
            e.EquipID_2 = excel.ReadInt32(readCellValue(row, names, "EquipID_2"))
            e.EquipLv_2 = excel.ReadInt32(readCellValue(row, names, "EquipLv_2"))
            e.EquipSkillLv_2 = excel.ReadInt32(readCellValue(row, names, "EquipSkillLv_2"))
            e.EquipRandomAttrib_2 = excel.ReadInt32(readCellValue(row, names, "EquipRandomAttrib_2"))
            e.EquipRandomAttribLv_2 = excel.ReadInt32(readCellValue(row, names, "EquipRandomAttribLv_2"))
            e.EquipID_3 = excel.ReadInt32(readCellValue(row, names, "EquipID_3"))
            e.EquipLv_3 = excel.ReadInt32(readCellValue(row, names, "EquipLv_3"))
            e.EquipSkillLv_3 = excel.ReadInt32(readCellValue(row, names, "EquipSkillLv_3"))
            e.EquipRandomAttrib_3 = excel.ReadInt32(readCellValue(row, names, "EquipRandomAttrib_3"))
            e.EquipRandomAttribLv_3 = excel.ReadInt32(readCellValue(row, names, "EquipRandomAttribLv_3"))
            e.EquipID_4 = excel.ReadInt32(readCellValue(row, names, "EquipID_4"))
            e.EquipLv_4 = excel.ReadInt32(readCellValue(row, names, "EquipLv_4"))
            e.EquipSkillLv_4 = excel.ReadInt32(readCellValue(row, names, "EquipSkillLv_4"))
            e.EquipRandomAttrib_4 = excel.ReadInt32(readCellValue(row, names, "EquipRandomAttrib_4"))
            e.EquipRandomAttribLv_4 = excel.ReadInt32(readCellValue(row, names, "EquipRandomAttribLv_4"))
            e.EquipID_5 = excel.ReadInt32(readCellValue(row, names, "EquipID_5"))
            e.EquipLv_5 = excel.ReadInt32(readCellValue(row, names, "EquipLv_5"))
            e.EquipSkillLv_5 = excel.ReadInt32(readCellValue(row, names, "EquipSkillLv_5"))
            e.EquipRandomAttrib_5 = excel.ReadInt32(readCellValue(row, names, "EquipRandomAttrib_5"))
            e.EquipRandomAttribLv_5 = excel.ReadInt32(readCellValue(row, names, "EquipRandomAttribLv_5"))

			tmpIDMap[excel.ReadInt32(readCellValue(row, names, "ID"))] = e

		}
	}
	this.indexID.Store(tmpIDMap)	

}
