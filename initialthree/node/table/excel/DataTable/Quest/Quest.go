package Quest

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	path2 "path"
	"strings"
	"sync/atomic"
	"initialthree/node/table/excel"
	"initialthree/zaplogger"
)

type Quest struct{
    ID int32 
    Type string 
    UnlockQuests string 
    ShowOnQuestFinished int32 
    Title string 
    Desc string 
    UnreceivedSendMail bool 
    Reward int32 
    SkipID int32 
    CondType_1 string 
    CondArg_1 string 
    CondType_2 string 
    CondArg_2 string 
    CondType_3 string 
    CondArg_3 string 
    CondType_4 string 
    CondArg_4 string 
    CondType_5 string 
    CondArg_5 string 
    CondType_6 string 
    CondArg_6 string 
    CondType_7 string 
    CondArg_7 string 
    CondType_8 string 
    CondArg_8 string 
    CondType_9 string 
    CondArg_9 string 
    TypeEnum int32 
    Condition []*Condition_ 
    UnlockQuestsArray []*UnlockQuests_ 

}

type Condition_ struct{
    Type int32 
    Arg string 

}

func readCondition(row, names []string)[]*Condition_{
	ret := make([]*Condition_, 0)
	base := excel.Split("CondType_,CondArg_", ",")

	for i := 1; ; i++ {
		if idx, find := hasCellName(names, fmt.Sprintf("%s%d", base[0][0], i));!find {
			break
		}else if _,ok := hasCellValue(row,idx);!ok{
			break
		}
		
		e := &Condition_{}
        e.Type = excel.ReadEnum(readCellValue(row, names, fmt.Sprintf("%s%d", base[0][0], i)))
        e.Arg = excel.ReadStr(readCellValue(row, names, fmt.Sprintf("%s%d", base[1][0], i)))

		ret = append(ret, e)
	}

	return ret
}

type UnlockQuests_ struct{
    QuestID int32 

}

func readUnlockQuestsArray(row, names []string)[]*UnlockQuests_{
	value := readCellValue(row, names, "UnlockQuests")
	if value == "" {
		return nil
	}

	r := excel.Split(value,",")
	if len(r) == 0 {
		return nil
	}

	ret := make([]*UnlockQuests_,0)
	for _, v := range r {
		if len(v) == 1{
			e := &UnlockQuests_{}
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

func GetID(key int32)*Quest{
	v := Table_.indexID.Load().(map[int32]*Quest)[key]
	if v == nil{
		zaplogger.GetSugar().Errorf("Table %s GetID %d is nil", Table_.xlsxName, key)
	}
	return v
}

func GetIDMap() map[int32]*Quest {
	return Table_.indexID.Load().(map[int32]*Quest)
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
	Table_.xlsxName = "Quest.xlsx"
	excel.AddTable(&Table_)
}

func (this *Table)Load(path string){
	xlFile, err := excelize.OpenFile(path2.Join(path2.Join(path, "/DataTable/任务"), this.xlsxName))
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	tableRows := xlFile.GetRows(xlFile.GetSheetName(xlFile.GetActiveSheetIndex()))
	names := tableRows[1]
	rows := tableRows[4:]

	tmpIDMap := map[int32]*Quest{}

	for _,row := range rows{
		if row[0] != "" {
			e := &Quest{}
            e.ID = excel.ReadInt32(readCellValue(row, names, "ID"))
            e.Type = excel.ReadStr(readCellValue(row, names, "Type"))
            e.UnlockQuests = excel.ReadStr(readCellValue(row, names, "UnlockQuests"))
            e.ShowOnQuestFinished = excel.ReadInt32(readCellValue(row, names, "ShowOnQuestFinished"))
            e.Title = excel.ReadStr(readCellValue(row, names, "Title"))
            e.Desc = excel.ReadStr(readCellValue(row, names, "Desc"))
            e.UnreceivedSendMail = excel.ReadBool(readCellValue(row, names, "UnreceivedSendMail"))
            e.Reward = excel.ReadInt32(readCellValue(row, names, "Reward"))
            e.SkipID = excel.ReadInt32(readCellValue(row, names, "SkipID"))
            e.CondType_1 = excel.ReadStr(readCellValue(row, names, "CondType_1"))
            e.CondArg_1 = excel.ReadStr(readCellValue(row, names, "CondArg_1"))
            e.CondType_2 = excel.ReadStr(readCellValue(row, names, "CondType_2"))
            e.CondArg_2 = excel.ReadStr(readCellValue(row, names, "CondArg_2"))
            e.CondType_3 = excel.ReadStr(readCellValue(row, names, "CondType_3"))
            e.CondArg_3 = excel.ReadStr(readCellValue(row, names, "CondArg_3"))
            e.CondType_4 = excel.ReadStr(readCellValue(row, names, "CondType_4"))
            e.CondArg_4 = excel.ReadStr(readCellValue(row, names, "CondArg_4"))
            e.CondType_5 = excel.ReadStr(readCellValue(row, names, "CondType_5"))
            e.CondArg_5 = excel.ReadStr(readCellValue(row, names, "CondArg_5"))
            e.CondType_6 = excel.ReadStr(readCellValue(row, names, "CondType_6"))
            e.CondArg_6 = excel.ReadStr(readCellValue(row, names, "CondArg_6"))
            e.CondType_7 = excel.ReadStr(readCellValue(row, names, "CondType_7"))
            e.CondArg_7 = excel.ReadStr(readCellValue(row, names, "CondArg_7"))
            e.CondType_8 = excel.ReadStr(readCellValue(row, names, "CondType_8"))
            e.CondArg_8 = excel.ReadStr(readCellValue(row, names, "CondArg_8"))
            e.CondType_9 = excel.ReadStr(readCellValue(row, names, "CondType_9"))
            e.CondArg_9 = excel.ReadStr(readCellValue(row, names, "CondArg_9"))
            e.TypeEnum = excel.ReadEnum(readCellValue(row, names, "Type"))
            e.Condition = readCondition(row, names)
            e.UnlockQuestsArray = readUnlockQuestsArray(row, names)

			tmpIDMap[excel.ReadInt32(readCellValue(row, names, "ID"))] = e

		}
	}
	this.indexID.Store(tmpIDMap)	

}
