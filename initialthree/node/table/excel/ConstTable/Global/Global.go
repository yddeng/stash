package Global

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	path2 "path"
	"strings"
	"sync/atomic"
	"initialthree/node/table/excel"
	"initialthree/zaplogger"
)

type Global struct{
    ID int32 
    MapID int32 
    PositionRate float64 
    DailyRefreshTime string 
    WeeklyRefreshTime string 
    MonthlyRefreshTime int32 
    MaxCharacterLevel int32 
    LevepUpCostItemsAndItOfferedExp string 
    MaxStarLevel int32 
    DefaultYuru int32 
    FavorExp string 
    DefaultUnlockPortrait int32 
    DefaultPortrait int32 
    DefaultUnlockPortraitFrame int32 
    DefaultPortraitFrame int32 
    DefaultUnlockPlayerCard int32 
    DefaultPlayerCard int32 
    DrawCardTenthGuaranteeRarity string 
    DrawCardGuaranteeRarity string 
    CanNotSkipUI string 
    DefaultCharacterTeam string 
    DrawCardDailyLimit int32 
    CharacterLvUpExpCostGold int32 
    DailyRefreshTimeStruct *DailyRefreshTime_ 
    LevepUpCostItemsAndItOfferedExpArray []*LevepUpCostItemsAndItOfferedExp_ 
    FavorExpArray []*FavorExp_ 
    DrawCardTenthGuaranteeRarityEnum int32 
    DrawCardGuaranteeRarityEnum int32 
    DefaultCharacterTeamArray []*DefaultCharacterTeam_ 

}

type DailyRefreshTime_ struct{
    Hour int32 
    Minute int32 

}

func readDailyRefreshTimeStruct(row, names []string)*DailyRefreshTime_{
	value := readCellValue(row, names, "DailyRefreshTime")
	r := excel.Split(value,":")
	
	e := &DailyRefreshTime_{}
    e.Hour = excel.ReadInt32(r[0][0])
    e.Minute = excel.ReadInt32(r[1][0])

	return e
}

type LevepUpCostItemsAndItOfferedExp_ struct{
    ItemID int32 
    Exp int32 

}

func readLevepUpCostItemsAndItOfferedExpArray(row, names []string)[]*LevepUpCostItemsAndItOfferedExp_{
	value := readCellValue(row, names, "LevepUpCostItemsAndItOfferedExp")
	if value == "" {
		return nil
	}

	r := excel.Split(value,",#")
	if len(r) == 0 {
		return nil
	}

	ret := make([]*LevepUpCostItemsAndItOfferedExp_,0)
	for _, v := range r {
		if len(v) == 2{
			e := &LevepUpCostItemsAndItOfferedExp_{}
	        e.ItemID = excel.ReadInt32(v[0])
        e.Exp = excel.ReadInt32(v[1])

			ret = append(ret, e)
		}
	}

	return ret
}

type FavorExp_ struct{
    Exp int32 

}

func readFavorExpArray(row, names []string)[]*FavorExp_{
	value := readCellValue(row, names, "FavorExp")
	if value == "" {
		return nil
	}

	r := excel.Split(value,",")
	if len(r) == 0 {
		return nil
	}

	ret := make([]*FavorExp_,0)
	for _, v := range r {
		if len(v) == 1{
			e := &FavorExp_{}
	        e.Exp = excel.ReadInt32(v[0])

			ret = append(ret, e)
		}
	}

	return ret
}

type DefaultCharacterTeam_ struct{
    ID int32 

}

func readDefaultCharacterTeamArray(row, names []string)[]*DefaultCharacterTeam_{
	value := readCellValue(row, names, "DefaultCharacterTeam")
	if value == "" {
		return nil
	}

	r := excel.Split(value,",")
	if len(r) == 0 {
		return nil
	}

	ret := make([]*DefaultCharacterTeam_,0)
	for _, v := range r {
		if len(v) == 1{
			e := &DefaultCharacterTeam_{}
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

func GetID(key int32)*Global{
	v := Table_.indexID.Load().(map[int32]*Global)[key]
	if v == nil{
		zaplogger.GetSugar().Errorf("Table %s GetID %d is nil", Table_.xlsxName, key)
	}
	return v
}

func GetIDMap() map[int32]*Global {
	return Table_.indexID.Load().(map[int32]*Global)
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
	Table_.xlsxName = "Global.xlsx"
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

	tmpIDMap := map[int32]*Global{}

	for _,row := range rows{
		if row[0] != "" {
			e := &Global{}
            e.ID = excel.ReadInt32(readCellValue(row, names, "ID"))
            e.MapID = excel.ReadInt32(readCellValue(row, names, "MapID"))
            e.PositionRate = excel.ReadFloat(readCellValue(row, names, "PositionRate"))
            e.DailyRefreshTime = excel.ReadStr(readCellValue(row, names, "DailyRefreshTime"))
            e.WeeklyRefreshTime = excel.ReadStr(readCellValue(row, names, "WeeklyRefreshTime"))
            e.MonthlyRefreshTime = excel.ReadInt32(readCellValue(row, names, "MonthlyRefreshTime"))
            e.MaxCharacterLevel = excel.ReadInt32(readCellValue(row, names, "MaxCharacterLevel"))
            e.LevepUpCostItemsAndItOfferedExp = excel.ReadStr(readCellValue(row, names, "LevepUpCostItemsAndItOfferedExp"))
            e.MaxStarLevel = excel.ReadInt32(readCellValue(row, names, "MaxStarLevel"))
            e.DefaultYuru = excel.ReadInt32(readCellValue(row, names, "DefaultYuru"))
            e.FavorExp = excel.ReadStr(readCellValue(row, names, "FavorExp"))
            e.DefaultUnlockPortrait = excel.ReadInt32(readCellValue(row, names, "DefaultUnlockPortrait"))
            e.DefaultPortrait = excel.ReadInt32(readCellValue(row, names, "DefaultPortrait"))
            e.DefaultUnlockPortraitFrame = excel.ReadInt32(readCellValue(row, names, "DefaultUnlockPortraitFrame"))
            e.DefaultPortraitFrame = excel.ReadInt32(readCellValue(row, names, "DefaultPortraitFrame"))
            e.DefaultUnlockPlayerCard = excel.ReadInt32(readCellValue(row, names, "DefaultUnlockPlayerCard"))
            e.DefaultPlayerCard = excel.ReadInt32(readCellValue(row, names, "DefaultPlayerCard"))
            e.DrawCardTenthGuaranteeRarity = excel.ReadStr(readCellValue(row, names, "DrawCardTenthGuaranteeRarity"))
            e.DrawCardGuaranteeRarity = excel.ReadStr(readCellValue(row, names, "DrawCardGuaranteeRarity"))
            e.CanNotSkipUI = excel.ReadStr(readCellValue(row, names, "CanNotSkipUI"))
            e.DefaultCharacterTeam = excel.ReadStr(readCellValue(row, names, "DefaultCharacterTeam"))
            e.DrawCardDailyLimit = excel.ReadInt32(readCellValue(row, names, "DrawCardDailyLimit"))
            e.CharacterLvUpExpCostGold = excel.ReadInt32(readCellValue(row, names, "CharacterLvUpExpCostGold"))
            e.DailyRefreshTimeStruct = readDailyRefreshTimeStruct(row, names)
            e.LevepUpCostItemsAndItOfferedExpArray = readLevepUpCostItemsAndItOfferedExpArray(row, names)
            e.FavorExpArray = readFavorExpArray(row, names)
            e.DrawCardTenthGuaranteeRarityEnum = excel.ReadEnum(readCellValue(row, names, "DrawCardTenthGuaranteeRarity"))
            e.DrawCardGuaranteeRarityEnum = excel.ReadEnum(readCellValue(row, names, "DrawCardGuaranteeRarity"))
            e.DefaultCharacterTeamArray = readDefaultCharacterTeamArray(row, names)

			tmpIDMap[excel.ReadInt32(readCellValue(row, names, "ID"))] = e

		}
	}
	this.indexID.Store(tmpIDMap)	

}
