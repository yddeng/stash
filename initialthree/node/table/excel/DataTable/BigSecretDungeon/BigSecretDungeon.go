package BigSecretDungeon

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	path2 "path"
	"strings"
	"sync/atomic"
	"initialthree/node/table/excel"
	"initialthree/zaplogger"
)

type BigSecretDungeon struct{
    ID int32 
    DungeonIDPool string 
    AffixCountMin int32 
    AffixCountMax int32 
    DungeonIDPoolArray []*DungeonIDPool_ 

}

type DungeonIDPool_ struct{
    DungeonID int32 

}

func readDungeonIDPoolArray(row, names []string)[]*DungeonIDPool_{
	value := readCellValue(row, names, "DungeonIDPool")
	if value == "" {
		return nil
	}

	r := excel.Split(value,",")
	if len(r) == 0 {
		return nil
	}

	ret := make([]*DungeonIDPool_,0)
	for _, v := range r {
		if len(v) == 1{
			e := &DungeonIDPool_{}
	        e.DungeonID = excel.ReadInt32(v[0])

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

func GetID(key int32)*BigSecretDungeon{
	v := Table_.indexID.Load().(map[int32]*BigSecretDungeon)[key]
	if v == nil{
		zaplogger.GetSugar().Errorf("Table %s GetID %d is nil", Table_.xlsxName, key)
	}
	return v
}

func GetIDMap() map[int32]*BigSecretDungeon {
	return Table_.indexID.Load().(map[int32]*BigSecretDungeon)
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
	Table_.xlsxName = "BigSecretDungeon.xlsx"
	excel.AddTable(&Table_)
}

func (this *Table)Load(path string){
	xlFile, err := excelize.OpenFile(path2.Join(path2.Join(path, "/DataTable/关卡/大秘境"), this.xlsxName))
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	tableRows := xlFile.GetRows(xlFile.GetSheetName(xlFile.GetActiveSheetIndex()))
	names := tableRows[1]
	rows := tableRows[4:]

	tmpIDMap := map[int32]*BigSecretDungeon{}

	for _,row := range rows{
		if row[0] != "" {
			e := &BigSecretDungeon{}
            e.ID = excel.ReadInt32(readCellValue(row, names, "ID"))
            e.DungeonIDPool = excel.ReadStr(readCellValue(row, names, "DungeonIDPool"))
            e.AffixCountMin = excel.ReadInt32(readCellValue(row, names, "AffixCountMin"))
            e.AffixCountMax = excel.ReadInt32(readCellValue(row, names, "AffixCountMax"))
            e.DungeonIDPoolArray = readDungeonIDPoolArray(row, names)

			tmpIDMap[excel.ReadInt32(readCellValue(row, names, "ID"))] = e

		}
	}
	this.indexID.Store(tmpIDMap)	

}
