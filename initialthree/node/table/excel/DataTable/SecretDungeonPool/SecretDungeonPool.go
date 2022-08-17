package SecretDungeonPool

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	path2 "path"
	"strings"
	"sync/atomic"
	"initialthree/node/table/excel"
	"initialthree/zaplogger"
)

type SecretDungeonPool struct{
    ID int32 
    DifficultyType string 
    Pool string 
    AffixCount int32 
    LeaderAffixCount int32 
    BossAffixCount int32 
    AffixPoolID string 
    LeaderAffixPoolID string 
    BossAffixPoolID string 
    PoolArray []*Pool_ 

}

type Pool_ struct{
    ID int32 

}

func readPoolArray(row, names []string)[]*Pool_{
	value := readCellValue(row, names, "Pool")
	if value == "" {
		return nil
	}

	r := excel.Split(value,",")
	if len(r) == 0 {
		return nil
	}

	ret := make([]*Pool_,0)
	for _, v := range r {
		if len(v) == 1{
			e := &Pool_{}
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

func GetID(key int32)*SecretDungeonPool{
	v := Table_.indexID.Load().(map[int32]*SecretDungeonPool)[key]
	if v == nil{
		zaplogger.GetSugar().Errorf("Table %s GetID %d is nil", Table_.xlsxName, key)
	}
	return v
}

func GetIDMap() map[int32]*SecretDungeonPool {
	return Table_.indexID.Load().(map[int32]*SecretDungeonPool)
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
	Table_.xlsxName = "SecretDungeonPool.xlsx"
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

	tmpIDMap := map[int32]*SecretDungeonPool{}

	for _,row := range rows{
		if row[0] != "" {
			e := &SecretDungeonPool{}
            e.ID = excel.ReadInt32(readCellValue(row, names, "ID"))
            e.DifficultyType = excel.ReadStr(readCellValue(row, names, "DifficultyType"))
            e.Pool = excel.ReadStr(readCellValue(row, names, "Pool"))
            e.AffixCount = excel.ReadInt32(readCellValue(row, names, "AffixCount"))
            e.LeaderAffixCount = excel.ReadInt32(readCellValue(row, names, "LeaderAffixCount"))
            e.BossAffixCount = excel.ReadInt32(readCellValue(row, names, "BossAffixCount"))
            e.AffixPoolID = excel.ReadStr(readCellValue(row, names, "AffixPoolID"))
            e.LeaderAffixPoolID = excel.ReadStr(readCellValue(row, names, "LeaderAffixPoolID"))
            e.BossAffixPoolID = excel.ReadStr(readCellValue(row, names, "BossAffixPoolID"))
            e.PoolArray = readPoolArray(row, names)

			tmpIDMap[excel.ReadInt32(readCellValue(row, names, "ID"))] = e

		}
	}
	this.indexID.Store(tmpIDMap)	

}
