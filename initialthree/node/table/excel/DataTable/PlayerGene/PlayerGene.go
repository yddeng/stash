package PlayerGene

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	path2 "path"
	"strings"
	"sync/atomic"
	"initialthree/node/table/excel"
	"initialthree/zaplogger"
)

type PlayerGene struct{
    ID int32 
    RoleName string 
    Name string 
    IconName string 
    ItemID int32 
    ItemCount int32 
    EffectDesc string 
    PlayerSkillID string 
    PlayerSkillLevelUp int32 
    Attri string 
    GeneCombatPower string 
    AttriArray []*Attri_ 

}

type Attri_ struct{
    ID int32 
    Val float64 

}

func readAttriArray(row, names []string)[]*Attri_{
	value := readCellValue(row, names, "Attri")
	if value == "" {
		return nil
	}

	r := excel.Split(value,",#")
	if len(r) == 0 {
		return nil
	}

	ret := make([]*Attri_,0)
	for _, v := range r {
		if len(v) == 2{
			e := &Attri_{}
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

func GetID(key int32)*PlayerGene{
	v := Table_.indexID.Load().(map[int32]*PlayerGene)[key]
	if v == nil{
		zaplogger.GetSugar().Errorf("Table %s GetID %d is nil", Table_.xlsxName, key)
	}
	return v
}

func GetIDMap() map[int32]*PlayerGene {
	return Table_.indexID.Load().(map[int32]*PlayerGene)
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
	Table_.xlsxName = "PlayerGene.xlsx"
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

	tmpIDMap := map[int32]*PlayerGene{}

	for _,row := range rows{
		if row[0] != "" {
			e := &PlayerGene{}
            e.ID = excel.ReadInt32(readCellValue(row, names, "ID"))
            e.RoleName = excel.ReadStr(readCellValue(row, names, "RoleName"))
            e.Name = excel.ReadStr(readCellValue(row, names, "Name"))
            e.IconName = excel.ReadStr(readCellValue(row, names, "IconName"))
            e.ItemID = excel.ReadInt32(readCellValue(row, names, "ItemID"))
            e.ItemCount = excel.ReadInt32(readCellValue(row, names, "ItemCount"))
            e.EffectDesc = excel.ReadStr(readCellValue(row, names, "EffectDesc"))
            e.PlayerSkillID = excel.ReadStr(readCellValue(row, names, "PlayerSkillID"))
            e.PlayerSkillLevelUp = excel.ReadInt32(readCellValue(row, names, "PlayerSkillLevelUp"))
            e.Attri = excel.ReadStr(readCellValue(row, names, "Attri"))
            e.GeneCombatPower = excel.ReadStr(readCellValue(row, names, "GeneCombatPower"))
            e.AttriArray = readAttriArray(row, names)

			tmpIDMap[excel.ReadInt32(readCellValue(row, names, "ID"))] = e

		}
	}
	this.indexID.Store(tmpIDMap)	

}
