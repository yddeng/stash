package NewbieGift

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	path2 "path"
	"strings"
	"sync/atomic"
	"initialthree/node/table/excel"
	"initialthree/zaplogger"
)

type NewbieGift struct{
    ID int32 
    QuestTitle string 
    QuestIDList string 
    GroupRewardQuest int32 
    ImportantRewardImage string 
    ImportantRewardTitle string 
    ImportantRewardContent string 
    QuestIDListArray []*QuestIDList_ 

}

type QuestIDList_ struct{
    QuestID int32 

}

func readQuestIDListArray(row, names []string)[]*QuestIDList_{
	value := readCellValue(row, names, "QuestIDList")
	if value == "" {
		return nil
	}

	r := excel.Split(value,",")
	if len(r) == 0 {
		return nil
	}

	ret := make([]*QuestIDList_,0)
	for _, v := range r {
		if len(v) == 1{
			e := &QuestIDList_{}
	        e.QuestID = excel.ReadInt32(v[0])

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

func GetID(key int32)*NewbieGift{
	v := Table_.indexID.Load().(map[int32]*NewbieGift)[key]
	if v == nil{
		zaplogger.GetSugar().Errorf("Table %s GetID %d is nil", Table_.xlsxName, key)
	}
	return v
}

func GetIDMap() map[int32]*NewbieGift {
	return Table_.indexID.Load().(map[int32]*NewbieGift)
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
	Table_.xlsxName = "NewbieGift.xlsx"
	excel.AddTable(&Table_)
}

func (this *Table)Load(path string){
	xlFile, err := excelize.OpenFile(path2.Join(path2.Join(path, "/DataTable/任务/七日任务"), this.xlsxName))
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	tableRows := xlFile.GetRows(xlFile.GetSheetName(xlFile.GetActiveSheetIndex()))
	names := tableRows[1]
	rows := tableRows[4:]

	tmpIDMap := map[int32]*NewbieGift{}

	for _,row := range rows{
		if row[0] != "" {
			e := &NewbieGift{}
            e.ID = excel.ReadInt32(readCellValue(row, names, "ID"))
            e.QuestTitle = excel.ReadStr(readCellValue(row, names, "QuestTitle"))
            e.QuestIDList = excel.ReadStr(readCellValue(row, names, "QuestIDList"))
            e.GroupRewardQuest = excel.ReadInt32(readCellValue(row, names, "GroupRewardQuest"))
            e.ImportantRewardImage = excel.ReadStr(readCellValue(row, names, "ImportantRewardImage"))
            e.ImportantRewardTitle = excel.ReadStr(readCellValue(row, names, "ImportantRewardTitle"))
            e.ImportantRewardContent = excel.ReadStr(readCellValue(row, names, "ImportantRewardContent"))
            e.QuestIDListArray = readQuestIDListArray(row, names)

			tmpIDMap[excel.ReadInt32(readCellValue(row, names, "ID"))] = e

		}
	}
	this.indexID.Store(tmpIDMap)	

}
