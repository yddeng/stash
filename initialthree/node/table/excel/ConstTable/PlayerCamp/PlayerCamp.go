package PlayerCamp

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	path2 "path"
	"strings"
	"sync/atomic"
	"initialthree/node/table/excel"
	"initialthree/zaplogger"
)

type PlayerCamp struct{
    ID int32 
    UnlockDungeon_1 int32 
    UnlockDungeon_2 int32 
    UnlockDungeon_3 int32 
    UnlockDungeon_4 int32 
    ReputationItem_1 int32 
    ReputationItem_2 int32 
    ReputationItem_3 int32 
    ReputationItem_4 int32 
    ReputationRefresh string 
    ReputationValueItem_1 int32 
    ReputationValueItem_2 int32 
    ReputationValueItem_3 int32 
    ReputationValueItem_4 int32 
    UnlockDungeon []*UnlockDungeon_ 
    ReputationItem []*ReputationItem_ 
    ReputationRefreshArray []*ReputationRefresh_ 

}

type UnlockDungeon_ struct{
    ID int32 

}

func readUnlockDungeon(row, names []string)[]*UnlockDungeon_{
	ret := make([]*UnlockDungeon_, 0)
	base := excel.Split("UnlockDungeon_", ",")

	for i := 1; ; i++ {
		if idx, find := hasCellName(names, fmt.Sprintf("%s%d", base[0][0], i));!find {
			break
		}else if _,ok := hasCellValue(row,idx);!ok{
			break
		}
		
		e := &UnlockDungeon_{}
        e.ID = excel.ReadInt32(readCellValue(row, names, fmt.Sprintf("%s%d", base[0][0], i)))

		ret = append(ret, e)
	}

	return ret
}

type ReputationItem_ struct{
    ID int32 

}

func readReputationItem(row, names []string)[]*ReputationItem_{
	ret := make([]*ReputationItem_, 0)
	base := excel.Split("ReputationItem_", ",")

	for i := 1; ; i++ {
		if idx, find := hasCellName(names, fmt.Sprintf("%s%d", base[0][0], i));!find {
			break
		}else if _,ok := hasCellValue(row,idx);!ok{
			break
		}
		
		e := &ReputationItem_{}
        e.ID = excel.ReadInt32(readCellValue(row, names, fmt.Sprintf("%s%d", base[0][0], i)))

		ret = append(ret, e)
	}

	return ret
}

type ReputationRefresh_ struct{
    Cost int32 

}

func readReputationRefreshArray(row, names []string)[]*ReputationRefresh_{
	value := readCellValue(row, names, "ReputationRefresh")
	if value == "" {
		return nil
	}

	r := excel.Split(value,",")
	if len(r) == 0 {
		return nil
	}

	ret := make([]*ReputationRefresh_,0)
	for _, v := range r {
		if len(v) == 1{
			e := &ReputationRefresh_{}
	        e.Cost = excel.ReadInt32(v[0])

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

func GetID(key int32)*PlayerCamp{
	v := Table_.indexID.Load().(map[int32]*PlayerCamp)[key]
	if v == nil{
		zaplogger.GetSugar().Errorf("Table %s GetID %d is nil", Table_.xlsxName, key)
	}
	return v
}

func GetIDMap() map[int32]*PlayerCamp {
	return Table_.indexID.Load().(map[int32]*PlayerCamp)
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
	Table_.xlsxName = "PlayerCamp.xlsx"
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

	tmpIDMap := map[int32]*PlayerCamp{}

	for _,row := range rows{
		if row[0] != "" {
			e := &PlayerCamp{}
            e.ID = excel.ReadInt32(readCellValue(row, names, "ID"))
            e.UnlockDungeon_1 = excel.ReadInt32(readCellValue(row, names, "UnlockDungeon_1"))
            e.UnlockDungeon_2 = excel.ReadInt32(readCellValue(row, names, "UnlockDungeon_2"))
            e.UnlockDungeon_3 = excel.ReadInt32(readCellValue(row, names, "UnlockDungeon_3"))
            e.UnlockDungeon_4 = excel.ReadInt32(readCellValue(row, names, "UnlockDungeon_4"))
            e.ReputationItem_1 = excel.ReadInt32(readCellValue(row, names, "ReputationItem_1"))
            e.ReputationItem_2 = excel.ReadInt32(readCellValue(row, names, "ReputationItem_2"))
            e.ReputationItem_3 = excel.ReadInt32(readCellValue(row, names, "ReputationItem_3"))
            e.ReputationItem_4 = excel.ReadInt32(readCellValue(row, names, "ReputationItem_4"))
            e.ReputationRefresh = excel.ReadStr(readCellValue(row, names, "ReputationRefresh"))
            e.ReputationValueItem_1 = excel.ReadInt32(readCellValue(row, names, "ReputationValueItem_1"))
            e.ReputationValueItem_2 = excel.ReadInt32(readCellValue(row, names, "ReputationValueItem_2"))
            e.ReputationValueItem_3 = excel.ReadInt32(readCellValue(row, names, "ReputationValueItem_3"))
            e.ReputationValueItem_4 = excel.ReadInt32(readCellValue(row, names, "ReputationValueItem_4"))
            e.UnlockDungeon = readUnlockDungeon(row, names)
            e.ReputationItem = readReputationItem(row, names)
            e.ReputationRefreshArray = readReputationRefreshArray(row, names)

			tmpIDMap[excel.ReadInt32(readCellValue(row, names, "ID"))] = e

		}
	}
	this.indexID.Store(tmpIDMap)	

}
