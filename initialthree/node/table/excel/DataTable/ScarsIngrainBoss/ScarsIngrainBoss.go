package ScarsIngrainBoss

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	path2 "path"
	"strings"
	"sync/atomic"
	"initialthree/node/table/excel"
	"initialthree/zaplogger"
)

type ScarsIngrainBoss struct{
    ID int32 
    DisplayID int32 
    ChallengeConfig string 
    Instance_1 int32 
    Skills_1 string 
    Buffs_1 string 
    Instance_2 int32 
    Skills_2 string 
    Buffs_2 string 
    Instance_3 int32 
    Skills_3 string 
    Buffs_3 string 
    Instance_4 int32 
    Skills_4 string 
    Buffs_4 string 
    Instance_5 int32 
    Skills_5 string 
    Buffs_5 string 
    ChallengeConfigArray []*ChallengeConfig_ 
    BossDifficulty []*BossDifficulty_ 

}

type ChallengeConfig_ struct{
    ID int32 

}

func readChallengeConfigArray(row, names []string)[]*ChallengeConfig_{
	value := readCellValue(row, names, "ChallengeConfig")
	if value == "" {
		return nil
	}

	r := excel.Split(value,",")
	if len(r) == 0 {
		return nil
	}

	ret := make([]*ChallengeConfig_,0)
	for _, v := range r {
		if len(v) == 1{
			e := &ChallengeConfig_{}
	        e.ID = excel.ReadInt32(v[0])

			ret = append(ret, e)
		}
	}

	return ret
}

type BossDifficulty_ struct{
    InstanceID int32 
    SkillStr string 
    BuffStr string 

}

func readBossDifficulty(row, names []string)[]*BossDifficulty_{
	ret := make([]*BossDifficulty_, 0)
	base := excel.Split("Instance_,skills_,buffs_", ",")

	for i := 1; ; i++ {
		if idx, find := hasCellName(names, fmt.Sprintf("%s%d", base[0][0], i));!find {
			break
		}else if _,ok := hasCellValue(row,idx);!ok{
			break
		}
		
		e := &BossDifficulty_{}
        e.InstanceID = excel.ReadInt32(readCellValue(row, names, fmt.Sprintf("%s%d", base[0][0], i)))
        e.SkillStr = excel.ReadStr(readCellValue(row, names, fmt.Sprintf("%s%d", base[1][0], i)))
        e.BuffStr = excel.ReadStr(readCellValue(row, names, fmt.Sprintf("%s%d", base[2][0], i)))

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

func GetID(key int32)*ScarsIngrainBoss{
	v := Table_.indexID.Load().(map[int32]*ScarsIngrainBoss)[key]
	if v == nil{
		zaplogger.GetSugar().Errorf("Table %s GetID %d is nil", Table_.xlsxName, key)
	}
	return v
}

func GetIDMap() map[int32]*ScarsIngrainBoss {
	return Table_.indexID.Load().(map[int32]*ScarsIngrainBoss)
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
	Table_.xlsxName = "ScarsIngrainBoss.xlsx"
	excel.AddTable(&Table_)
}

func (this *Table)Load(path string){
	xlFile, err := excelize.OpenFile(path2.Join(path2.Join(path, "/DataTable/关卡/战痕"), this.xlsxName))
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	tableRows := xlFile.GetRows(xlFile.GetSheetName(xlFile.GetActiveSheetIndex()))
	names := tableRows[1]
	rows := tableRows[4:]

	tmpIDMap := map[int32]*ScarsIngrainBoss{}

	for _,row := range rows{
		if row[0] != "" {
			e := &ScarsIngrainBoss{}
            e.ID = excel.ReadInt32(readCellValue(row, names, "ID"))
            e.DisplayID = excel.ReadInt32(readCellValue(row, names, "DisplayID"))
            e.ChallengeConfig = excel.ReadStr(readCellValue(row, names, "ChallengeConfig"))
            e.Instance_1 = excel.ReadInt32(readCellValue(row, names, "Instance_1"))
            e.Skills_1 = excel.ReadStr(readCellValue(row, names, "skills_1"))
            e.Buffs_1 = excel.ReadStr(readCellValue(row, names, "buffs_1"))
            e.Instance_2 = excel.ReadInt32(readCellValue(row, names, "Instance_2"))
            e.Skills_2 = excel.ReadStr(readCellValue(row, names, "skills_2"))
            e.Buffs_2 = excel.ReadStr(readCellValue(row, names, "buffs_2"))
            e.Instance_3 = excel.ReadInt32(readCellValue(row, names, "Instance_3"))
            e.Skills_3 = excel.ReadStr(readCellValue(row, names, "skills_3"))
            e.Buffs_3 = excel.ReadStr(readCellValue(row, names, "buffs_3"))
            e.Instance_4 = excel.ReadInt32(readCellValue(row, names, "Instance_4"))
            e.Skills_4 = excel.ReadStr(readCellValue(row, names, "skills_4"))
            e.Buffs_4 = excel.ReadStr(readCellValue(row, names, "buffs_4"))
            e.Instance_5 = excel.ReadInt32(readCellValue(row, names, "Instance_5"))
            e.Skills_5 = excel.ReadStr(readCellValue(row, names, "skills_5"))
            e.Buffs_5 = excel.ReadStr(readCellValue(row, names, "buffs_5"))
            e.ChallengeConfigArray = readChallengeConfigArray(row, names)
            e.BossDifficulty = readBossDifficulty(row, names)

			tmpIDMap[excel.ReadInt32(readCellValue(row, names, "ID"))] = e

		}
	}
	this.indexID.Store(tmpIDMap)	

}
