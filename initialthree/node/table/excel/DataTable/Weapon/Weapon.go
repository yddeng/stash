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
    Name string 
    Icon string 
    DrawCardPicture string 
    DrawCardSingleResultPicture string 
    DrawCardResultPicture string 
    Desc string 
    LeftWeapon string 
    RightWeapon string 
    LeftWeaponBoneName string 
    RightWeaponBoneName string 
    DisappearVfx string 
    WeaponType string 
    RarityType string 
    WeaponSkillConfig int32 
    AttribConfig_1 int32 
    AttribConfig_2 int32 
    AttribConfig_3 int32 
    BreakThroughConfig int32 
    LevelExpConfig int32 
    DecomposeConfig int32 
    SupplyExp int32 
    WeaponSkillCombatPower string 
    AccessWayList string 
    RarityTypeEnum int32 
    WeaponTypeEnum int32 
    AttrConfigs []*AttrConfigs_ 

}

type AttrConfigs_ struct{
    ID int32 

}

func readAttrConfigs(row, names []string)[]*AttrConfigs_{
	ret := make([]*AttrConfigs_, 0)
	base := excel.Split("AttribConfig_", ",")

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
	xlFile, err := excelize.OpenFile(path2.Join(path2.Join(path, "/DataTable/武器"), this.xlsxName))
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
            e.Name = excel.ReadStr(readCellValue(row, names, "Name"))
            e.Icon = excel.ReadStr(readCellValue(row, names, "Icon"))
            e.DrawCardPicture = excel.ReadStr(readCellValue(row, names, "DrawCardPicture"))
            e.DrawCardSingleResultPicture = excel.ReadStr(readCellValue(row, names, "DrawCardSingleResultPicture"))
            e.DrawCardResultPicture = excel.ReadStr(readCellValue(row, names, "DrawCardResultPicture"))
            e.Desc = excel.ReadStr(readCellValue(row, names, "Desc"))
            e.LeftWeapon = excel.ReadStr(readCellValue(row, names, "LeftWeapon"))
            e.RightWeapon = excel.ReadStr(readCellValue(row, names, "RightWeapon"))
            e.LeftWeaponBoneName = excel.ReadStr(readCellValue(row, names, "LeftWeaponBoneName"))
            e.RightWeaponBoneName = excel.ReadStr(readCellValue(row, names, "RightWeaponBoneName"))
            e.DisappearVfx = excel.ReadStr(readCellValue(row, names, "DisappearVfx"))
            e.WeaponType = excel.ReadStr(readCellValue(row, names, "WeaponType"))
            e.RarityType = excel.ReadStr(readCellValue(row, names, "RarityType"))
            e.WeaponSkillConfig = excel.ReadInt32(readCellValue(row, names, "WeaponSkillConfig"))
            e.AttribConfig_1 = excel.ReadInt32(readCellValue(row, names, "AttribConfig_1"))
            e.AttribConfig_2 = excel.ReadInt32(readCellValue(row, names, "AttribConfig_2"))
            e.AttribConfig_3 = excel.ReadInt32(readCellValue(row, names, "AttribConfig_3"))
            e.BreakThroughConfig = excel.ReadInt32(readCellValue(row, names, "BreakThroughConfig"))
            e.LevelExpConfig = excel.ReadInt32(readCellValue(row, names, "LevelExpConfig"))
            e.DecomposeConfig = excel.ReadInt32(readCellValue(row, names, "DecomposeConfig"))
            e.SupplyExp = excel.ReadInt32(readCellValue(row, names, "SupplyExp"))
            e.WeaponSkillCombatPower = excel.ReadStr(readCellValue(row, names, "WeaponSkillCombatPower"))
            e.AccessWayList = excel.ReadStr(readCellValue(row, names, "AccessWayList"))
            e.RarityTypeEnum = excel.ReadEnum(readCellValue(row, names, "RarityType"))
            e.WeaponTypeEnum = excel.ReadEnum(readCellValue(row, names, "WeaponType"))
            e.AttrConfigs = readAttrConfigs(row, names)

			tmpIDMap[excel.ReadInt32(readCellValue(row, names, "ID"))] = e

		}
	}
	this.indexID.Store(tmpIDMap)	

}
