package PlayerSkill

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	path2 "path"
	"strings"
	"sync/atomic"
	"initialthree/node/table/excel"
	"initialthree/zaplogger"
)

type PlayerSkill struct{
    ID int32 
    Comment string 
    SkillID int32 
    TagType string 
    LevelUpPopDescTitle string 
    RequiredGeneLevel int32 
    LevelUpCharacterLvLimit_1 int32 
    LevelUpPopUpDescValue_1 string 
    LevelUpCostGold_1 int32 
    LevelUpCostItem_1 string 
    LevelUpCharacterLvLimit_2 int32 
    LevelUpPopUpDescValue_2 string 
    LevelUpCostGold_2 int32 
    LevelUpCostItem_2 string 
    LevelUpCharacterLvLimit_3 int32 
    LevelUpPopUpDescValue_3 string 
    LevelUpCostGold_3 int32 
    LevelUpCostItem_3 string 
    LevelUpCharacterLvLimit_4 int32 
    LevelUpPopUpDescValue_4 string 
    LevelUpCostGold_4 int32 
    LevelUpCostItem_4 string 
    LevelUpCharacterLvLimit_5 int32 
    LevelUpPopUpDescValue_5 string 
    LevelUpCostGold_5 int32 
    LevelUpCostItem_5 string 
    LevelUpCharacterLvLimit_6 int32 
    LevelUpPopUpDescValue_6 string 
    LevelUpCostGold_6 int32 
    LevelUpCostItem_6 string 
    LevelUpCharacterLvLimit_7 int32 
    LevelUpPopUpDescValue_7 string 
    LevelUpCostGold_7 int32 
    LevelUpCostItem_7 string 
    LevelUpCharacterLvLimit_8 int32 
    LevelUpPopUpDescValue_8 string 
    LevelUpCostGold_8 int32 
    LevelUpCostItem_8 string 
    LevelUpCharacterLvLimit_9 int32 
    LevelUpPopUpDescValue_9 string 
    LevelUpCostGold_9 int32 
    LevelUpCostItem_9 string 
    LevelUpCharacterLvLimit_10 int32 
    LevelUpPopUpDescValue_10 string 
    LevelUpCostGold_10 int32 
    LevelUpCostItem_10 string 
    LevelUpCharacterLvLimit_11 int32 
    LevelUpPopUpDescValue_11 string 
    LevelUpCostGold_11 int32 
    LevelUpCostItem_11 string 
    LevelUpCharacterLvLimit_12 int32 
    LevelUpPopUpDescValue_12 string 
    LevelUpCostGold_12 int32 
    LevelUpCostItem_12 string 
    LevelUpCharacterLvLimit_13 int32 
    LevelUpPopUpDescValue_13 string 
    LevelUpCostGold_13 int32 
    LevelUpCostItem_13 string 
    LevelUpCharacterLvLimit_14 int32 
    LevelUpPopUpDescValue_14 string 
    LevelUpCostGold_14 int32 
    LevelUpCostItem_14 string 
    LevelUpCharacterLvLimit_15 int32 
    LevelUpPopUpDescValue_15 string 
    LevelUpCostGold_15 int32 
    LevelUpCostItem_15 string 
    Skill []*Skill_ 

}

type Skill_ struct{
    LimitLevel int32 
    Gold int32 
    ItemStr string 

}

func readSkill(row, names []string)[]*Skill_{
	ret := make([]*Skill_, 0)
	base := excel.Split("LevelUpCharacterLvLimit_,LevelUpCostGold_,LevelUpCostItem_", ",")

	for i := 1; ; i++ {
		if idx, find := hasCellName(names, fmt.Sprintf("%s%d", base[0][0], i));!find {
			break
		}else if _,ok := hasCellValue(row,idx);!ok{
			break
		}
		
		e := &Skill_{}
        e.LimitLevel = excel.ReadInt32(readCellValue(row, names, fmt.Sprintf("%s%d", base[0][0], i)))
        e.Gold = excel.ReadInt32(readCellValue(row, names, fmt.Sprintf("%s%d", base[1][0], i)))
        e.ItemStr = excel.ReadStr(readCellValue(row, names, fmt.Sprintf("%s%d", base[2][0], i)))

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

func GetID(key int32)*PlayerSkill{
	v := Table_.indexID.Load().(map[int32]*PlayerSkill)[key]
	if v == nil{
		zaplogger.GetSugar().Errorf("Table %s GetID %d is nil", Table_.xlsxName, key)
	}
	return v
}

func GetIDMap() map[int32]*PlayerSkill {
	return Table_.indexID.Load().(map[int32]*PlayerSkill)
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
	Table_.xlsxName = "PlayerSkill.xlsx"
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

	tmpIDMap := map[int32]*PlayerSkill{}

	for _,row := range rows{
		if row[0] != "" {
			e := &PlayerSkill{}
            e.ID = excel.ReadInt32(readCellValue(row, names, "ID"))
            e.Comment = excel.ReadStr(readCellValue(row, names, "Comment"))
            e.SkillID = excel.ReadInt32(readCellValue(row, names, "SkillID"))
            e.TagType = excel.ReadStr(readCellValue(row, names, "TagType"))
            e.LevelUpPopDescTitle = excel.ReadStr(readCellValue(row, names, "LevelUpPopDescTitle"))
            e.RequiredGeneLevel = excel.ReadInt32(readCellValue(row, names, "RequiredGeneLevel"))
            e.LevelUpCharacterLvLimit_1 = excel.ReadInt32(readCellValue(row, names, "LevelUpCharacterLvLimit_1"))
            e.LevelUpPopUpDescValue_1 = excel.ReadStr(readCellValue(row, names, "LevelUpPopUpDescValue_1"))
            e.LevelUpCostGold_1 = excel.ReadInt32(readCellValue(row, names, "LevelUpCostGold_1"))
            e.LevelUpCostItem_1 = excel.ReadStr(readCellValue(row, names, "LevelUpCostItem_1"))
            e.LevelUpCharacterLvLimit_2 = excel.ReadInt32(readCellValue(row, names, "LevelUpCharacterLvLimit_2"))
            e.LevelUpPopUpDescValue_2 = excel.ReadStr(readCellValue(row, names, "LevelUpPopUpDescValue_2"))
            e.LevelUpCostGold_2 = excel.ReadInt32(readCellValue(row, names, "LevelUpCostGold_2"))
            e.LevelUpCostItem_2 = excel.ReadStr(readCellValue(row, names, "LevelUpCostItem_2"))
            e.LevelUpCharacterLvLimit_3 = excel.ReadInt32(readCellValue(row, names, "LevelUpCharacterLvLimit_3"))
            e.LevelUpPopUpDescValue_3 = excel.ReadStr(readCellValue(row, names, "LevelUpPopUpDescValue_3"))
            e.LevelUpCostGold_3 = excel.ReadInt32(readCellValue(row, names, "LevelUpCostGold_3"))
            e.LevelUpCostItem_3 = excel.ReadStr(readCellValue(row, names, "LevelUpCostItem_3"))
            e.LevelUpCharacterLvLimit_4 = excel.ReadInt32(readCellValue(row, names, "LevelUpCharacterLvLimit_4"))
            e.LevelUpPopUpDescValue_4 = excel.ReadStr(readCellValue(row, names, "LevelUpPopUpDescValue_4"))
            e.LevelUpCostGold_4 = excel.ReadInt32(readCellValue(row, names, "LevelUpCostGold_4"))
            e.LevelUpCostItem_4 = excel.ReadStr(readCellValue(row, names, "LevelUpCostItem_4"))
            e.LevelUpCharacterLvLimit_5 = excel.ReadInt32(readCellValue(row, names, "LevelUpCharacterLvLimit_5"))
            e.LevelUpPopUpDescValue_5 = excel.ReadStr(readCellValue(row, names, "LevelUpPopUpDescValue_5"))
            e.LevelUpCostGold_5 = excel.ReadInt32(readCellValue(row, names, "LevelUpCostGold_5"))
            e.LevelUpCostItem_5 = excel.ReadStr(readCellValue(row, names, "LevelUpCostItem_5"))
            e.LevelUpCharacterLvLimit_6 = excel.ReadInt32(readCellValue(row, names, "LevelUpCharacterLvLimit_6"))
            e.LevelUpPopUpDescValue_6 = excel.ReadStr(readCellValue(row, names, "LevelUpPopUpDescValue_6"))
            e.LevelUpCostGold_6 = excel.ReadInt32(readCellValue(row, names, "LevelUpCostGold_6"))
            e.LevelUpCostItem_6 = excel.ReadStr(readCellValue(row, names, "LevelUpCostItem_6"))
            e.LevelUpCharacterLvLimit_7 = excel.ReadInt32(readCellValue(row, names, "LevelUpCharacterLvLimit_7"))
            e.LevelUpPopUpDescValue_7 = excel.ReadStr(readCellValue(row, names, "LevelUpPopUpDescValue_7"))
            e.LevelUpCostGold_7 = excel.ReadInt32(readCellValue(row, names, "LevelUpCostGold_7"))
            e.LevelUpCostItem_7 = excel.ReadStr(readCellValue(row, names, "LevelUpCostItem_7"))
            e.LevelUpCharacterLvLimit_8 = excel.ReadInt32(readCellValue(row, names, "LevelUpCharacterLvLimit_8"))
            e.LevelUpPopUpDescValue_8 = excel.ReadStr(readCellValue(row, names, "LevelUpPopUpDescValue_8"))
            e.LevelUpCostGold_8 = excel.ReadInt32(readCellValue(row, names, "LevelUpCostGold_8"))
            e.LevelUpCostItem_8 = excel.ReadStr(readCellValue(row, names, "LevelUpCostItem_8"))
            e.LevelUpCharacterLvLimit_9 = excel.ReadInt32(readCellValue(row, names, "LevelUpCharacterLvLimit_9"))
            e.LevelUpPopUpDescValue_9 = excel.ReadStr(readCellValue(row, names, "LevelUpPopUpDescValue_9"))
            e.LevelUpCostGold_9 = excel.ReadInt32(readCellValue(row, names, "LevelUpCostGold_9"))
            e.LevelUpCostItem_9 = excel.ReadStr(readCellValue(row, names, "LevelUpCostItem_9"))
            e.LevelUpCharacterLvLimit_10 = excel.ReadInt32(readCellValue(row, names, "LevelUpCharacterLvLimit_10"))
            e.LevelUpPopUpDescValue_10 = excel.ReadStr(readCellValue(row, names, "LevelUpPopUpDescValue_10"))
            e.LevelUpCostGold_10 = excel.ReadInt32(readCellValue(row, names, "LevelUpCostGold_10"))
            e.LevelUpCostItem_10 = excel.ReadStr(readCellValue(row, names, "LevelUpCostItem_10"))
            e.LevelUpCharacterLvLimit_11 = excel.ReadInt32(readCellValue(row, names, "LevelUpCharacterLvLimit_11"))
            e.LevelUpPopUpDescValue_11 = excel.ReadStr(readCellValue(row, names, "LevelUpPopUpDescValue_11"))
            e.LevelUpCostGold_11 = excel.ReadInt32(readCellValue(row, names, "LevelUpCostGold_11"))
            e.LevelUpCostItem_11 = excel.ReadStr(readCellValue(row, names, "LevelUpCostItem_11"))
            e.LevelUpCharacterLvLimit_12 = excel.ReadInt32(readCellValue(row, names, "LevelUpCharacterLvLimit_12"))
            e.LevelUpPopUpDescValue_12 = excel.ReadStr(readCellValue(row, names, "LevelUpPopUpDescValue_12"))
            e.LevelUpCostGold_12 = excel.ReadInt32(readCellValue(row, names, "LevelUpCostGold_12"))
            e.LevelUpCostItem_12 = excel.ReadStr(readCellValue(row, names, "LevelUpCostItem_12"))
            e.LevelUpCharacterLvLimit_13 = excel.ReadInt32(readCellValue(row, names, "LevelUpCharacterLvLimit_13"))
            e.LevelUpPopUpDescValue_13 = excel.ReadStr(readCellValue(row, names, "LevelUpPopUpDescValue_13"))
            e.LevelUpCostGold_13 = excel.ReadInt32(readCellValue(row, names, "LevelUpCostGold_13"))
            e.LevelUpCostItem_13 = excel.ReadStr(readCellValue(row, names, "LevelUpCostItem_13"))
            e.LevelUpCharacterLvLimit_14 = excel.ReadInt32(readCellValue(row, names, "LevelUpCharacterLvLimit_14"))
            e.LevelUpPopUpDescValue_14 = excel.ReadStr(readCellValue(row, names, "LevelUpPopUpDescValue_14"))
            e.LevelUpCostGold_14 = excel.ReadInt32(readCellValue(row, names, "LevelUpCostGold_14"))
            e.LevelUpCostItem_14 = excel.ReadStr(readCellValue(row, names, "LevelUpCostItem_14"))
            e.LevelUpCharacterLvLimit_15 = excel.ReadInt32(readCellValue(row, names, "LevelUpCharacterLvLimit_15"))
            e.LevelUpPopUpDescValue_15 = excel.ReadStr(readCellValue(row, names, "LevelUpPopUpDescValue_15"))
            e.LevelUpCostGold_15 = excel.ReadInt32(readCellValue(row, names, "LevelUpCostGold_15"))
            e.LevelUpCostItem_15 = excel.ReadStr(readCellValue(row, names, "LevelUpCostItem_15"))
            e.Skill = readSkill(row, names)

			tmpIDMap[excel.ReadInt32(readCellValue(row, names, "ID"))] = e

		}
	}
	this.indexID.Store(tmpIDMap)	

}
