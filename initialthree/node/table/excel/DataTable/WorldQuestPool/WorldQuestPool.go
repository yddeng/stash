package WorldQuestPool

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	path2 "path"
	"strings"
	"sync/atomic"
	"initialthree/node/table/excel"
	"initialthree/zaplogger"
)

type WorldQuestPool struct{
    ID int32 
    MinLevel int32 
    MaxLevel int32 
    QuestList_1 string 
    QuestList_2 string 
    QuestList_3 string 
    QuestList_4 string 
    Weight int32 
    QuestList []*QuestList_ 

}

type QuestList_ struct{
    List string 

}

func readQuestList(row, names []string)[]*QuestList_{
	ret := make([]*QuestList_, 0)
	base := excel.Split("QuestList_", ",")

	for i := 1; ; i++ {
		if idx, find := hasCellName(names, fmt.Sprintf("%s%d", base[0][0], i));!find {
			break
		}else if _,ok := hasCellValue(row,idx);!ok{
			break
		}
		
		e := &QuestList_{}
        e.List = excel.ReadStr(readCellValue(row, names, fmt.Sprintf("%s%d", base[0][0], i)))

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

func GetID(key int32)*WorldQuestPool{
	v := Table_.indexID.Load().(map[int32]*WorldQuestPool)[key]
	if v == nil{
		zaplogger.GetSugar().Errorf("Table %s GetID %d is nil", Table_.xlsxName, key)
	}
	return v
}

func GetIDMap() map[int32]*WorldQuestPool {
	return Table_.indexID.Load().(map[int32]*WorldQuestPool)
}

func GetMinMaxLevelID(MinLevel, MaxLevel int32)*WorldQuestPool{
    key := MinLevel*1000 + MaxLevel*1
	return Table_.indexID.Load().(map[int32]*WorldQuestPool)[key]
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
	Table_.xlsxName = "WorldQuestPool.xlsx"
	excel.AddTable(&Table_)
}

func (this *Table)Load(path string){
	xlFile, err := excelize.OpenFile(path2.Join(path2.Join(path, "/DataTable/关卡/世界任务"), this.xlsxName))
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	tableRows := xlFile.GetRows(xlFile.GetSheetName(xlFile.GetActiveSheetIndex()))
	names := tableRows[1]
	rows := tableRows[4:]

	tmpIDMap := map[int32]*WorldQuestPool{}

	for _,row := range rows{
		if row[0] != "" {
			e := &WorldQuestPool{}
            e.ID = excel.ReadInt32(readCellValue(row, names, "ID"))
            e.MinLevel = excel.ReadInt32(readCellValue(row, names, "MinLevel"))
            e.MaxLevel = excel.ReadInt32(readCellValue(row, names, "MaxLevel"))
            e.QuestList_1 = excel.ReadStr(readCellValue(row, names, "QuestList_1"))
            e.QuestList_2 = excel.ReadStr(readCellValue(row, names, "QuestList_2"))
            e.QuestList_3 = excel.ReadStr(readCellValue(row, names, "QuestList_3"))
            e.QuestList_4 = excel.ReadStr(readCellValue(row, names, "QuestList_4"))
            e.Weight = excel.ReadInt32(readCellValue(row, names, "Weight"))
            e.QuestList = readQuestList(row, names)

			tmpIDMap[excel.ReadInt32(readCellValue(row, names, "ID"))] = e

		}
	}
	this.indexID.Store(tmpIDMap)	

}
