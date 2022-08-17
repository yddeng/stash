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
    InitMainStoryQuest string 
    DailyRewardQuest string 
    DailyQuest_1 string 
    DailyQuest_2 string 
    DailyQuest_3 string 
    DailyQuest_4 string 
    DailyQuest_5 string 
    DailyQuest_6 string 
    DailyQuest_7 string 
    WeeklyQuest string 
    WorldQuestEachCount int32 
    WorldQuestMaxCount int32 
    WorldQuestRefreshFreeMaxCount int32 
    WorldQuestRefreshItemID int32 
    WorldQuestRefreshDiamondCount int32 
    NewbieGiftUnlockLevel int32 
    NewbieGiftUnlockLessTime int32 
    InitMainStoryQuestArray []*InitMainStoryQuest_ 
    DailyRewardQuestArray []*DailyRewardQuest_ 
    DailyQuest_1Array []*DailyQuest_1_ 
    DailyQuest_2Array []*DailyQuest_2_ 
    DailyQuest_3Array []*DailyQuest_3_ 
    DailyQuest_4Array []*DailyQuest_4_ 
    DailyQuest_5Array []*DailyQuest_5_ 
    DailyQuest_6Array []*DailyQuest_6_ 
    DailyQuest_7Array []*DailyQuest_7_ 
    WeeklyQuestArray []*WeeklyQuest_ 

}

type InitMainStoryQuest_ struct{
    ID int32 

}

func readInitMainStoryQuestArray(row, names []string)[]*InitMainStoryQuest_{
	value := readCellValue(row, names, "InitMainStoryQuest")
	if value == "" {
		return nil
	}

	r := excel.Split(value,",")
	if len(r) == 0 {
		return nil
	}

	ret := make([]*InitMainStoryQuest_,0)
	for _, v := range r {
		if len(v) == 1{
			e := &InitMainStoryQuest_{}
	        e.ID = excel.ReadInt32(v[0])

			ret = append(ret, e)
		}
	}

	return ret
}

type DailyRewardQuest_ struct{
    ID int32 

}

func readDailyRewardQuestArray(row, names []string)[]*DailyRewardQuest_{
	value := readCellValue(row, names, "DailyRewardQuest")
	if value == "" {
		return nil
	}

	r := excel.Split(value,",")
	if len(r) == 0 {
		return nil
	}

	ret := make([]*DailyRewardQuest_,0)
	for _, v := range r {
		if len(v) == 1{
			e := &DailyRewardQuest_{}
	        e.ID = excel.ReadInt32(v[0])

			ret = append(ret, e)
		}
	}

	return ret
}

type DailyQuest_1_ struct{
    ID int32 

}

func readDailyQuest_1Array(row, names []string)[]*DailyQuest_1_{
	value := readCellValue(row, names, "DailyQuest_1")
	if value == "" {
		return nil
	}

	r := excel.Split(value,",")
	if len(r) == 0 {
		return nil
	}

	ret := make([]*DailyQuest_1_,0)
	for _, v := range r {
		if len(v) == 1{
			e := &DailyQuest_1_{}
	        e.ID = excel.ReadInt32(v[0])

			ret = append(ret, e)
		}
	}

	return ret
}

type DailyQuest_2_ struct{
    ID int32 

}

func readDailyQuest_2Array(row, names []string)[]*DailyQuest_2_{
	value := readCellValue(row, names, "DailyQuest_2")
	if value == "" {
		return nil
	}

	r := excel.Split(value,",")
	if len(r) == 0 {
		return nil
	}

	ret := make([]*DailyQuest_2_,0)
	for _, v := range r {
		if len(v) == 1{
			e := &DailyQuest_2_{}
	        e.ID = excel.ReadInt32(v[0])

			ret = append(ret, e)
		}
	}

	return ret
}

type DailyQuest_3_ struct{
    ID int32 

}

func readDailyQuest_3Array(row, names []string)[]*DailyQuest_3_{
	value := readCellValue(row, names, "DailyQuest_3")
	if value == "" {
		return nil
	}

	r := excel.Split(value,",")
	if len(r) == 0 {
		return nil
	}

	ret := make([]*DailyQuest_3_,0)
	for _, v := range r {
		if len(v) == 1{
			e := &DailyQuest_3_{}
	        e.ID = excel.ReadInt32(v[0])

			ret = append(ret, e)
		}
	}

	return ret
}

type DailyQuest_4_ struct{
    ID int32 

}

func readDailyQuest_4Array(row, names []string)[]*DailyQuest_4_{
	value := readCellValue(row, names, "DailyQuest_4")
	if value == "" {
		return nil
	}

	r := excel.Split(value,",")
	if len(r) == 0 {
		return nil
	}

	ret := make([]*DailyQuest_4_,0)
	for _, v := range r {
		if len(v) == 1{
			e := &DailyQuest_4_{}
	        e.ID = excel.ReadInt32(v[0])

			ret = append(ret, e)
		}
	}

	return ret
}

type DailyQuest_5_ struct{
    ID int32 

}

func readDailyQuest_5Array(row, names []string)[]*DailyQuest_5_{
	value := readCellValue(row, names, "DailyQuest_5")
	if value == "" {
		return nil
	}

	r := excel.Split(value,",")
	if len(r) == 0 {
		return nil
	}

	ret := make([]*DailyQuest_5_,0)
	for _, v := range r {
		if len(v) == 1{
			e := &DailyQuest_5_{}
	        e.ID = excel.ReadInt32(v[0])

			ret = append(ret, e)
		}
	}

	return ret
}

type DailyQuest_6_ struct{
    ID int32 

}

func readDailyQuest_6Array(row, names []string)[]*DailyQuest_6_{
	value := readCellValue(row, names, "DailyQuest_6")
	if value == "" {
		return nil
	}

	r := excel.Split(value,",")
	if len(r) == 0 {
		return nil
	}

	ret := make([]*DailyQuest_6_,0)
	for _, v := range r {
		if len(v) == 1{
			e := &DailyQuest_6_{}
	        e.ID = excel.ReadInt32(v[0])

			ret = append(ret, e)
		}
	}

	return ret
}

type DailyQuest_7_ struct{
    ID int32 

}

func readDailyQuest_7Array(row, names []string)[]*DailyQuest_7_{
	value := readCellValue(row, names, "DailyQuest_7")
	if value == "" {
		return nil
	}

	r := excel.Split(value,",")
	if len(r) == 0 {
		return nil
	}

	ret := make([]*DailyQuest_7_,0)
	for _, v := range r {
		if len(v) == 1{
			e := &DailyQuest_7_{}
	        e.ID = excel.ReadInt32(v[0])

			ret = append(ret, e)
		}
	}

	return ret
}

type WeeklyQuest_ struct{
    ID int32 

}

func readWeeklyQuestArray(row, names []string)[]*WeeklyQuest_{
	value := readCellValue(row, names, "WeeklyQuest")
	if value == "" {
		return nil
	}

	r := excel.Split(value,",")
	if len(r) == 0 {
		return nil
	}

	ret := make([]*WeeklyQuest_,0)
	for _, v := range r {
		if len(v) == 1{
			e := &WeeklyQuest_{}
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
	xlFile, err := excelize.OpenFile(path2.Join(path2.Join(path, "/ConstTable"), this.xlsxName))
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
            e.InitMainStoryQuest = excel.ReadStr(readCellValue(row, names, "InitMainStoryQuest"))
            e.DailyRewardQuest = excel.ReadStr(readCellValue(row, names, "DailyRewardQuest"))
            e.DailyQuest_1 = excel.ReadStr(readCellValue(row, names, "DailyQuest_1"))
            e.DailyQuest_2 = excel.ReadStr(readCellValue(row, names, "DailyQuest_2"))
            e.DailyQuest_3 = excel.ReadStr(readCellValue(row, names, "DailyQuest_3"))
            e.DailyQuest_4 = excel.ReadStr(readCellValue(row, names, "DailyQuest_4"))
            e.DailyQuest_5 = excel.ReadStr(readCellValue(row, names, "DailyQuest_5"))
            e.DailyQuest_6 = excel.ReadStr(readCellValue(row, names, "DailyQuest_6"))
            e.DailyQuest_7 = excel.ReadStr(readCellValue(row, names, "DailyQuest_7"))
            e.WeeklyQuest = excel.ReadStr(readCellValue(row, names, "WeeklyQuest"))
            e.WorldQuestEachCount = excel.ReadInt32(readCellValue(row, names, "WorldQuestEachCount"))
            e.WorldQuestMaxCount = excel.ReadInt32(readCellValue(row, names, "WorldQuestMaxCount"))
            e.WorldQuestRefreshFreeMaxCount = excel.ReadInt32(readCellValue(row, names, "WorldQuestRefreshFreeMaxCount"))
            e.WorldQuestRefreshItemID = excel.ReadInt32(readCellValue(row, names, "WorldQuestRefreshItemID"))
            e.WorldQuestRefreshDiamondCount = excel.ReadInt32(readCellValue(row, names, "WorldQuestRefreshDiamondCount"))
            e.NewbieGiftUnlockLevel = excel.ReadInt32(readCellValue(row, names, "NewbieGiftUnlockLevel"))
            e.NewbieGiftUnlockLessTime = excel.ReadInt32(readCellValue(row, names, "NewbieGiftUnlockLessTime"))
            e.InitMainStoryQuestArray = readInitMainStoryQuestArray(row, names)
            e.DailyRewardQuestArray = readDailyRewardQuestArray(row, names)
            e.DailyQuest_1Array = readDailyQuest_1Array(row, names)
            e.DailyQuest_2Array = readDailyQuest_2Array(row, names)
            e.DailyQuest_3Array = readDailyQuest_3Array(row, names)
            e.DailyQuest_4Array = readDailyQuest_4Array(row, names)
            e.DailyQuest_5Array = readDailyQuest_5Array(row, names)
            e.DailyQuest_6Array = readDailyQuest_6Array(row, names)
            e.DailyQuest_7Array = readDailyQuest_7Array(row, names)
            e.WeeklyQuestArray = readWeeklyQuestArray(row, names)

			tmpIDMap[excel.ReadInt32(readCellValue(row, names, "ID"))] = e

		}
	}
	this.indexID.Store(tmpIDMap)	

}
