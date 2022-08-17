package SecretDungeon

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	path2 "path"
	"strings"
	"sync/atomic"
	"initialthree/node/table/excel"
	"initialthree/zaplogger"
)

type SecretDungeon struct{
    ID int32 
    DifficultyType string 
    DungeonID int32 
    MonsterPool string 
    CommonMonsterKillCounts string 
    EliteMonsterID string 
    DifficultyTypeEnum int32 

}

var Table_ Table

type Table struct{
	xlsxName string
	nextPath string
	indexID  atomic.Value

}

func GetID(key int32)*SecretDungeon{
	v := Table_.indexID.Load().(map[int32]*SecretDungeon)[key]
	if v == nil{
		zaplogger.GetSugar().Errorf("Table %s GetID %d is nil", Table_.xlsxName, key)
	}
	return v
}

func GetIDMap() map[int32]*SecretDungeon {
	return Table_.indexID.Load().(map[int32]*SecretDungeon)
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
	Table_.xlsxName = "SecretDungeon.xlsx"
	excel.AddTable(&Table_)
}

func (this *Table)Load(path string){
	xlFile, err := excelize.OpenFile(path2.Join(path2.Join(path, "/DataTable/关卡/秘境"), this.xlsxName))
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	tableRows := xlFile.GetRows(xlFile.GetSheetName(xlFile.GetActiveSheetIndex()))
	names := tableRows[1]
	rows := tableRows[4:]

	tmpIDMap := map[int32]*SecretDungeon{}

	for _,row := range rows{
		if row[0] != "" {
			e := &SecretDungeon{}
            e.ID = excel.ReadInt32(readCellValue(row, names, "ID"))
            e.DifficultyType = excel.ReadStr(readCellValue(row, names, "DifficultyType"))
            e.DungeonID = excel.ReadInt32(readCellValue(row, names, "DungeonID"))
            e.MonsterPool = excel.ReadStr(readCellValue(row, names, "MonsterPool"))
            e.CommonMonsterKillCounts = excel.ReadStr(readCellValue(row, names, "CommonMonsterKillCounts"))
            e.EliteMonsterID = excel.ReadStr(readCellValue(row, names, "EliteMonsterID"))
            e.DifficultyTypeEnum = excel.ReadEnum(readCellValue(row, names, "DifficultyType"))

			tmpIDMap[excel.ReadInt32(readCellValue(row, names, "ID"))] = e

		}
	}
	this.indexID.Store(tmpIDMap)	

}
