package Dungeon

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	path2 "path"
	"strings"
	"sync/atomic"
	"initialthree/node/table/excel"
	"initialthree/zaplogger"
)

type Dungeon struct{
    ID int32 
    UnlockType_1 string 
    UnlockArgs_1 string 
    UnlockType_2 string 
    UnlockArgs_2 string 
    UnlockType_3 string 
    UnlockArgs_3 string 
    FirstCostType string 
    FirstCostArgs int32 
    CommonCostType string 
    CommonCostArgs int32 
    FirstDrop int32 
    CommonDrop int32 
    FirstPlayerExp int32 
    CommonPlayerExp int32 
    FirstCharacterExp int32 
    CommonCharacterExp int32 
    DungeonSweep int32 
    TeamLimitType string 
    TeamLimitValue int32 
    InstanceID int32 
    ClearTimeLimit int32 
    TimeLimitShow int32 
    AllowResurrect bool 
    ResurrectID int32 
    SystemType string 
    SystemConfig int32 
    PlaySucceedTimeline bool 
    DisplayType_1 string 
    DisplayArg_1 string 
    DisplayType_2 string 
    DisplayArg_2 string 
    DisplayType_3 string 
    DisplayArg_3 string 
    DisplayType_4 string 
    DisplayArg_4 string 
    LoadingID int32 
    IsWorldMap bool 
    WorldAreaID int32 
    LocationID string 
    DemoCharacters string 
    TeamLimitTypeEnum int32 
    Unlocks []*Unlocks_ 
    Rewards []*Rewards_ 
    FirstCostTypeEnum int32 
    CommonCostTypeEnum int32 
    SystemTypeEnum int32 
    DemoCharactersArray []*DemoCharacters_ 

}

type Unlocks_ struct{
    Type int32 
    Arg string 

}

func readUnlocks(row, names []string)[]*Unlocks_{
	ret := make([]*Unlocks_, 0)
	base := excel.Split("UnlockType_,UnlockArgs_", ",")

	for i := 1; ; i++ {
		if idx, find := hasCellName(names, fmt.Sprintf("%s%d", base[0][0], i));!find {
			break
		}else if _,ok := hasCellValue(row,idx);!ok{
			break
		}
		
		e := &Unlocks_{}
        e.Type = excel.ReadEnum(readCellValue(row, names, fmt.Sprintf("%s%d", base[0][0], i)))
        e.Arg = excel.ReadStr(readCellValue(row, names, fmt.Sprintf("%s%d", base[1][0], i)))

		ret = append(ret, e)
	}

	return ret
}

type Rewards_ struct{
    Type int32 
    Arg int32 

}

func readRewards(row, names []string)[]*Rewards_{
	ret := make([]*Rewards_, 0)
	base := excel.Split("RewardType_,RewardArgs_", ",")

	for i := 1; ; i++ {
		if idx, find := hasCellName(names, fmt.Sprintf("%s%d", base[0][0], i));!find {
			break
		}else if _,ok := hasCellValue(row,idx);!ok{
			break
		}
		
		e := &Rewards_{}
        e.Type = excel.ReadEnum(readCellValue(row, names, fmt.Sprintf("%s%d", base[0][0], i)))
        e.Arg = excel.ReadInt32(readCellValue(row, names, fmt.Sprintf("%s%d", base[1][0], i)))

		ret = append(ret, e)
	}

	return ret
}

type DemoCharacters_ struct{
    CharacterID int32 

}

func readDemoCharactersArray(row, names []string)[]*DemoCharacters_{
	value := readCellValue(row, names, "DemoCharacters")
	if value == "" {
		return nil
	}

	r := excel.Split(value,",")
	if len(r) == 0 {
		return nil
	}

	ret := make([]*DemoCharacters_,0)
	for _, v := range r {
		if len(v) == 1{
			e := &DemoCharacters_{}
	        e.CharacterID = excel.ReadInt32(v[0])

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

func GetID(key int32)*Dungeon{
	v := Table_.indexID.Load().(map[int32]*Dungeon)[key]
	if v == nil{
		zaplogger.GetSugar().Errorf("Table %s GetID %d is nil", Table_.xlsxName, key)
	}
	return v
}

func GetIDMap() map[int32]*Dungeon {
	return Table_.indexID.Load().(map[int32]*Dungeon)
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
	Table_.xlsxName = "Dungeon.xlsx"
	excel.AddTable(&Table_)
}

func (this *Table)Load(path string){
	xlFile, err := excelize.OpenFile(path2.Join(path2.Join(path, "/DataTable/关卡"), this.xlsxName))
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	tableRows := xlFile.GetRows(xlFile.GetSheetName(xlFile.GetActiveSheetIndex()))
	names := tableRows[1]
	rows := tableRows[4:]

	tmpIDMap := map[int32]*Dungeon{}

	for _,row := range rows{
		if row[0] != "" {
			e := &Dungeon{}
            e.ID = excel.ReadInt32(readCellValue(row, names, "ID"))
            e.UnlockType_1 = excel.ReadStr(readCellValue(row, names, "UnlockType_1"))
            e.UnlockArgs_1 = excel.ReadStr(readCellValue(row, names, "UnlockArgs_1"))
            e.UnlockType_2 = excel.ReadStr(readCellValue(row, names, "UnlockType_2"))
            e.UnlockArgs_2 = excel.ReadStr(readCellValue(row, names, "UnlockArgs_2"))
            e.UnlockType_3 = excel.ReadStr(readCellValue(row, names, "UnlockType_3"))
            e.UnlockArgs_3 = excel.ReadStr(readCellValue(row, names, "UnlockArgs_3"))
            e.FirstCostType = excel.ReadStr(readCellValue(row, names, "FirstCostType"))
            e.FirstCostArgs = excel.ReadInt32(readCellValue(row, names, "FirstCostArgs"))
            e.CommonCostType = excel.ReadStr(readCellValue(row, names, "CommonCostType"))
            e.CommonCostArgs = excel.ReadInt32(readCellValue(row, names, "CommonCostArgs"))
            e.FirstDrop = excel.ReadInt32(readCellValue(row, names, "FirstDrop"))
            e.CommonDrop = excel.ReadInt32(readCellValue(row, names, "CommonDrop"))
            e.FirstPlayerExp = excel.ReadInt32(readCellValue(row, names, "FirstPlayerExp"))
            e.CommonPlayerExp = excel.ReadInt32(readCellValue(row, names, "CommonPlayerExp"))
            e.FirstCharacterExp = excel.ReadInt32(readCellValue(row, names, "FirstCharacterExp"))
            e.CommonCharacterExp = excel.ReadInt32(readCellValue(row, names, "CommonCharacterExp"))
            e.DungeonSweep = excel.ReadInt32(readCellValue(row, names, "DungeonSweep"))
            e.TeamLimitType = excel.ReadStr(readCellValue(row, names, "TeamLimitType"))
            e.TeamLimitValue = excel.ReadInt32(readCellValue(row, names, "TeamLimitValue"))
            e.InstanceID = excel.ReadInt32(readCellValue(row, names, "InstanceID"))
            e.ClearTimeLimit = excel.ReadInt32(readCellValue(row, names, "ClearTimeLimit"))
            e.TimeLimitShow = excel.ReadInt32(readCellValue(row, names, "TimeLimitShow"))
            e.AllowResurrect = excel.ReadBool(readCellValue(row, names, "AllowResurrect"))
            e.ResurrectID = excel.ReadInt32(readCellValue(row, names, "ResurrectID"))
            e.SystemType = excel.ReadStr(readCellValue(row, names, "SystemType"))
            e.SystemConfig = excel.ReadInt32(readCellValue(row, names, "SystemConfig"))
            e.PlaySucceedTimeline = excel.ReadBool(readCellValue(row, names, "PlaySucceedTimeline"))
            e.DisplayType_1 = excel.ReadStr(readCellValue(row, names, "DisplayType_1"))
            e.DisplayArg_1 = excel.ReadStr(readCellValue(row, names, "DisplayArg_1"))
            e.DisplayType_2 = excel.ReadStr(readCellValue(row, names, "DisplayType_2"))
            e.DisplayArg_2 = excel.ReadStr(readCellValue(row, names, "DisplayArg_2"))
            e.DisplayType_3 = excel.ReadStr(readCellValue(row, names, "DisplayType_3"))
            e.DisplayArg_3 = excel.ReadStr(readCellValue(row, names, "DisplayArg_3"))
            e.DisplayType_4 = excel.ReadStr(readCellValue(row, names, "DisplayType_4"))
            e.DisplayArg_4 = excel.ReadStr(readCellValue(row, names, "DisplayArg_4"))
            e.LoadingID = excel.ReadInt32(readCellValue(row, names, "LoadingID"))
            e.IsWorldMap = excel.ReadBool(readCellValue(row, names, "IsWorldMap"))
            e.WorldAreaID = excel.ReadInt32(readCellValue(row, names, "WorldAreaID"))
            e.LocationID = excel.ReadStr(readCellValue(row, names, "LocationID"))
            e.DemoCharacters = excel.ReadStr(readCellValue(row, names, "DemoCharacters"))
            e.TeamLimitTypeEnum = excel.ReadEnum(readCellValue(row, names, "TeamLimitType"), "")
            e.Unlocks = readUnlocks(row, names)
            e.Rewards = readRewards(row, names)
            e.FirstCostTypeEnum = excel.ReadEnum(readCellValue(row, names, "FirstCostType"), "")
            e.CommonCostTypeEnum = excel.ReadEnum(readCellValue(row, names, "CommonCostType"), "")
            e.SystemTypeEnum = excel.ReadEnum(readCellValue(row, names, "SystemType"))
            e.DemoCharactersArray = readDemoCharactersArray(row, names)

			tmpIDMap[excel.ReadInt32(readCellValue(row, names, "ID"))] = e

		}
	}
	this.indexID.Store(tmpIDMap)	

}
