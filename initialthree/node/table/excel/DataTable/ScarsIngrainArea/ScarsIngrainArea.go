package ScarsIngrainArea

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	path2 "path"
	"strings"
	"sync/atomic"
	"initialthree/node/table/excel"
	"initialthree/zaplogger"
)

type ScarsIngrainArea struct{
    ID int32 
    MinLevel int32 
    MaxLevel int32 
    AreaName string 
    Icon string 
    DanIcon string 
    BossIDs string 
    BossCount int32 
    ScoreReward int32 
    RankReward int32 
    BossIDsArray []*BossIDs_ 

}

type BossIDs_ struct{
    ID int32 

}

func readBossIDsArray(row, names []string)[]*BossIDs_{
	value := readCellValue(row, names, "BossIDs")
	if value == "" {
		return nil
	}

	r := excel.Split(value,",")
	if len(r) == 0 {
		return nil
	}

	ret := make([]*BossIDs_,0)
	for _, v := range r {
		if len(v) == 1{
			e := &BossIDs_{}
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

func GetID(key int32)*ScarsIngrainArea{
	v := Table_.indexID.Load().(map[int32]*ScarsIngrainArea)[key]
	if v == nil{
		zaplogger.GetSugar().Errorf("Table %s GetID %d is nil", Table_.xlsxName, key)
	}
	return v
}

func GetIDMap() map[int32]*ScarsIngrainArea {
	return Table_.indexID.Load().(map[int32]*ScarsIngrainArea)
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
	Table_.xlsxName = "ScarsIngrainArea.xlsx"
	excel.AddTable(&Table_)
}

func (this *Table)Load(path string){
	xlFile, err := excelize.OpenFile(path2.Join(path2.Join(path, "/DataTable/关卡/战痕"), this.xlsxName))
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	tableRows := xlFile.GetRows(xlFile.GetSheetName(xlFile.GetActiveSheetIndex()))
	names := tableRows[1]
	rows := tableRows[4:]

	tmpIDMap := map[int32]*ScarsIngrainArea{}

	for _,row := range rows{
		if row[0] != "" {
			e := &ScarsIngrainArea{}
            e.ID = excel.ReadInt32(readCellValue(row, names, "ID"))
            e.MinLevel = excel.ReadInt32(readCellValue(row, names, "MinLevel"))
            e.MaxLevel = excel.ReadInt32(readCellValue(row, names, "MaxLevel"))
            e.AreaName = excel.ReadStr(readCellValue(row, names, "AreaName"))
            e.Icon = excel.ReadStr(readCellValue(row, names, "Icon"))
            e.DanIcon = excel.ReadStr(readCellValue(row, names, "DanIcon"))
            e.BossIDs = excel.ReadStr(readCellValue(row, names, "BossIDs"))
            e.BossCount = excel.ReadInt32(readCellValue(row, names, "BossCount"))
            e.ScoreReward = excel.ReadInt32(readCellValue(row, names, "ScoreReward"))
            e.RankReward = excel.ReadInt32(readCellValue(row, names, "RankReward"))
            e.BossIDsArray = readBossIDsArray(row, names)

			tmpIDMap[excel.ReadInt32(readCellValue(row, names, "ID"))] = e

		}
	}
	this.indexID.Store(tmpIDMap)	

}
