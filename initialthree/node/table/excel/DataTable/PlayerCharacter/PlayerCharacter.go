package PlayerCharacter

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	path2 "path"
	"strings"
	"sync/atomic"
	"initialthree/node/table/excel"
	"initialthree/zaplogger"
)

type PlayerCharacter struct{
    ID int32 
    Name string 
    Sex int32 
    Rarity string 
    GiftType string 
    BreakIDList string 
    ProfessionType int32 
    WeaponType int32 
    ResonanceConfig int32 
    DefaultWeapon int32 
    DefaultCharacterResourceID int32 
    PlayeSkills string 
    Camp string 
    DrawCardItemIDs string 
    DrawCardItemCounts string 
    DrawCardTimes int32 
    MaxGeneLvDrawCardItemIDs string 
    MaxGeneLvDrawCardItemCounts string 
    DrawCardTalk string 
    CharacterSkillCoefficientID int32 
    PassiveSkillCombatPower string 
    Backstory string 
    RarityEnum int32 
    BreakIDListArray []*BreakIDList_ 
    GiftTypeEnum int32 
    PlayeSkillsArray []*PlayeSkills_ 
    DrawCardItemIDsArray []*DrawCardItemIDs_ 
    DrawCardItemCountsArray []*DrawCardItemCounts_ 
    MaxGeneLvDrawCardItemIDsArray []*MaxGeneLvDrawCardItemIDs_ 
    MaxGeneLvDrawCardItemCountsArray []*MaxGeneLvDrawCardItemCounts_ 

}

type BreakIDList_ struct{
    ID int32 

}

func readBreakIDListArray(row, names []string)[]*BreakIDList_{
	value := readCellValue(row, names, "BreakIDList")
	if value == "" {
		return nil
	}

	r := excel.Split(value,",")
	if len(r) == 0 {
		return nil
	}

	ret := make([]*BreakIDList_,0)
	for _, v := range r {
		if len(v) == 1{
			e := &BreakIDList_{}
	        e.ID = excel.ReadInt32(v[0])

			ret = append(ret, e)
		}
	}

	return ret
}

type PlayeSkills_ struct{
    ID int32 

}

func readPlayeSkillsArray(row, names []string)[]*PlayeSkills_{
	value := readCellValue(row, names, "PlayeSkills")
	if value == "" {
		return nil
	}

	r := excel.Split(value,",")
	if len(r) == 0 {
		return nil
	}

	ret := make([]*PlayeSkills_,0)
	for _, v := range r {
		if len(v) == 1{
			e := &PlayeSkills_{}
	        e.ID = excel.ReadInt32(v[0])

			ret = append(ret, e)
		}
	}

	return ret
}

type DrawCardItemIDs_ struct{
    ID int32 

}

func readDrawCardItemIDsArray(row, names []string)[]*DrawCardItemIDs_{
	value := readCellValue(row, names, "DrawCardItemIDs")
	if value == "" {
		return nil
	}

	r := excel.Split(value,",")
	if len(r) == 0 {
		return nil
	}

	ret := make([]*DrawCardItemIDs_,0)
	for _, v := range r {
		if len(v) == 1{
			e := &DrawCardItemIDs_{}
	        e.ID = excel.ReadInt32(v[0])

			ret = append(ret, e)
		}
	}

	return ret
}

type DrawCardItemCounts_ struct{
    Count int32 

}

func readDrawCardItemCountsArray(row, names []string)[]*DrawCardItemCounts_{
	value := readCellValue(row, names, "DrawCardItemCounts")
	if value == "" {
		return nil
	}

	r := excel.Split(value,",")
	if len(r) == 0 {
		return nil
	}

	ret := make([]*DrawCardItemCounts_,0)
	for _, v := range r {
		if len(v) == 1{
			e := &DrawCardItemCounts_{}
	        e.Count = excel.ReadInt32(v[0])

			ret = append(ret, e)
		}
	}

	return ret
}

type MaxGeneLvDrawCardItemIDs_ struct{
    ID int32 

}

func readMaxGeneLvDrawCardItemIDsArray(row, names []string)[]*MaxGeneLvDrawCardItemIDs_{
	value := readCellValue(row, names, "MaxGeneLvDrawCardItemIDs")
	if value == "" {
		return nil
	}

	r := excel.Split(value,",")
	if len(r) == 0 {
		return nil
	}

	ret := make([]*MaxGeneLvDrawCardItemIDs_,0)
	for _, v := range r {
		if len(v) == 1{
			e := &MaxGeneLvDrawCardItemIDs_{}
	        e.ID = excel.ReadInt32(v[0])

			ret = append(ret, e)
		}
	}

	return ret
}

type MaxGeneLvDrawCardItemCounts_ struct{
    Count int32 

}

func readMaxGeneLvDrawCardItemCountsArray(row, names []string)[]*MaxGeneLvDrawCardItemCounts_{
	value := readCellValue(row, names, "MaxGeneLvDrawCardItemCounts")
	if value == "" {
		return nil
	}

	r := excel.Split(value,",")
	if len(r) == 0 {
		return nil
	}

	ret := make([]*MaxGeneLvDrawCardItemCounts_,0)
	for _, v := range r {
		if len(v) == 1{
			e := &MaxGeneLvDrawCardItemCounts_{}
	        e.Count = excel.ReadInt32(v[0])

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

func GetID(key int32)*PlayerCharacter{
	v := Table_.indexID.Load().(map[int32]*PlayerCharacter)[key]
	if v == nil{
		zaplogger.GetSugar().Errorf("Table %s GetID %d is nil", Table_.xlsxName, key)
	}
	return v
}

func GetIDMap() map[int32]*PlayerCharacter {
	return Table_.indexID.Load().(map[int32]*PlayerCharacter)
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
	Table_.xlsxName = "PlayerCharacter.xlsx"
	excel.AddTable(&Table_)
}

func (this *Table)Load(path string){
	xlFile, err := excelize.OpenFile(path2.Join(path2.Join(path, "/DataTable/角色"), this.xlsxName))
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	tableRows := xlFile.GetRows(xlFile.GetSheetName(xlFile.GetActiveSheetIndex()))
	names := tableRows[1]
	rows := tableRows[4:]

	tmpIDMap := map[int32]*PlayerCharacter{}

	for _,row := range rows{
		if row[0] != "" {
			e := &PlayerCharacter{}
            e.ID = excel.ReadInt32(readCellValue(row, names, "ID"))
            e.Name = excel.ReadStr(readCellValue(row, names, "Name"))
            e.Sex = excel.ReadInt32(readCellValue(row, names, "Sex"))
            e.Rarity = excel.ReadStr(readCellValue(row, names, "Rarity"))
            e.GiftType = excel.ReadStr(readCellValue(row, names, "GiftType"))
            e.BreakIDList = excel.ReadStr(readCellValue(row, names, "BreakIDList"))
            e.ProfessionType = excel.ReadInt32(readCellValue(row, names, "ProfessionType"))
            e.WeaponType = excel.ReadInt32(readCellValue(row, names, "WeaponType"))
            e.ResonanceConfig = excel.ReadInt32(readCellValue(row, names, "ResonanceConfig"))
            e.DefaultWeapon = excel.ReadInt32(readCellValue(row, names, "DefaultWeapon"))
            e.DefaultCharacterResourceID = excel.ReadInt32(readCellValue(row, names, "DefaultCharacterResourceID"))
            e.PlayeSkills = excel.ReadStr(readCellValue(row, names, "PlayeSkills"))
            e.Camp = excel.ReadStr(readCellValue(row, names, "Camp"))
            e.DrawCardItemIDs = excel.ReadStr(readCellValue(row, names, "DrawCardItemIDs"))
            e.DrawCardItemCounts = excel.ReadStr(readCellValue(row, names, "DrawCardItemCounts"))
            e.DrawCardTimes = excel.ReadInt32(readCellValue(row, names, "DrawCardTimes"))
            e.MaxGeneLvDrawCardItemIDs = excel.ReadStr(readCellValue(row, names, "MaxGeneLvDrawCardItemIDs"))
            e.MaxGeneLvDrawCardItemCounts = excel.ReadStr(readCellValue(row, names, "MaxGeneLvDrawCardItemCounts"))
            e.DrawCardTalk = excel.ReadStr(readCellValue(row, names, "DrawCardTalk"))
            e.CharacterSkillCoefficientID = excel.ReadInt32(readCellValue(row, names, "CharacterSkillCoefficientID"))
            e.PassiveSkillCombatPower = excel.ReadStr(readCellValue(row, names, "PassiveSkillCombatPower"))
            e.Backstory = excel.ReadStr(readCellValue(row, names, "Backstory"))
            e.RarityEnum = excel.ReadEnum(readCellValue(row, names, "Rarity"))
            e.BreakIDListArray = readBreakIDListArray(row, names)
            e.GiftTypeEnum = excel.ReadEnum(readCellValue(row, names, "GiftType"))
            e.PlayeSkillsArray = readPlayeSkillsArray(row, names)
            e.DrawCardItemIDsArray = readDrawCardItemIDsArray(row, names)
            e.DrawCardItemCountsArray = readDrawCardItemCountsArray(row, names)
            e.MaxGeneLvDrawCardItemIDsArray = readMaxGeneLvDrawCardItemIDsArray(row, names)
            e.MaxGeneLvDrawCardItemCountsArray = readMaxGeneLvDrawCardItemCountsArray(row, names)

			tmpIDMap[excel.ReadInt32(readCellValue(row, names, "ID"))] = e

		}
	}
	this.indexID.Store(tmpIDMap)	

}
