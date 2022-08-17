package Skill

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	path2 "path"
	"strings"
	"sync/atomic"
	"initialthree/node/table/excel"
	"initialthree/zaplogger"
)

type Skill struct{
    ID int32 
    SkillName string 
    BehaviorIDs string 
    Icon string 
    Brief string 
    Desc_1 string 
    DescArg_1 string 
    Attri_1 string 
    Desc_2 string 
    DescArg_2 string 
    Attri_2 string 
    Desc_3 string 
    DescArg_3 string 
    Attri_3 string 
    Desc_4 string 
    DescArg_4 string 
    Attri_4 string 
    Desc_5 string 
    DescArg_5 string 
    Attri_5 string 
    Desc_6 string 
    DescArg_6 string 
    Attri_6 string 
    Desc_7 string 
    DescArg_7 string 
    Attri_7 string 
    Desc_8 string 
    DescArg_8 string 
    Attri_8 string 
    Desc_9 string 
    DescArg_9 string 
    Attri_9 string 
    Desc_10 string 
    DescArg_10 string 
    Attri_10 string 
    Desc_11 string 
    DescArg_11 string 
    Attri_11 string 
    Desc_12 string 
    DescArg_12 string 
    Attri_12 string 
    Desc_13 string 
    DescArg_13 string 
    Attri_13 string 
    Desc_14 string 
    DescArg_14 string 
    Attri_14 string 
    Desc_15 string 
    DescArg_15 string 
    Attri_15 string 
    Damage []*Damage_ 

}

type Damage_ struct{
    Desc string 

}

func readDamage(row, names []string)[]*Damage_{
	ret := make([]*Damage_, 0)
	base := excel.Split("Desc_", ",")

	for i := 1; ; i++ {
		if idx, find := hasCellName(names, fmt.Sprintf("%s%d", base[0][0], i));!find {
			break
		}else if _,ok := hasCellValue(row,idx);!ok{
			break
		}
		
		e := &Damage_{}
        e.Desc = excel.ReadStr(readCellValue(row, names, fmt.Sprintf("%s%d", base[0][0], i)))

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

func GetID(key int32)*Skill{
	v := Table_.indexID.Load().(map[int32]*Skill)[key]
	if v == nil{
		zaplogger.GetSugar().Errorf("Table %s GetID %d is nil", Table_.xlsxName, key)
	}
	return v
}

func GetIDMap() map[int32]*Skill {
	return Table_.indexID.Load().(map[int32]*Skill)
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
	Table_.xlsxName = "Skill.xlsx"
	excel.AddTable(&Table_)
}

func (this *Table)Load(path string){
	xlFile, err := excelize.OpenFile(path2.Join(path2.Join(path, "/DataTable"), this.xlsxName))
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	tableRows := xlFile.GetRows(xlFile.GetSheetName(xlFile.GetActiveSheetIndex()))
	names := tableRows[1]
	rows := tableRows[4:]

	tmpIDMap := map[int32]*Skill{}

	for _,row := range rows{
		if row[0] != "" {
			e := &Skill{}
            e.ID = excel.ReadInt32(readCellValue(row, names, "ID"))
            e.SkillName = excel.ReadStr(readCellValue(row, names, "SkillName"))
            e.BehaviorIDs = excel.ReadStr(readCellValue(row, names, "BehaviorIDs"))
            e.Icon = excel.ReadStr(readCellValue(row, names, "Icon"))
            e.Brief = excel.ReadStr(readCellValue(row, names, "Brief"))
            e.Desc_1 = excel.ReadStr(readCellValue(row, names, "Desc_1"))
            e.DescArg_1 = excel.ReadStr(readCellValue(row, names, "DescArg_1"))
            e.Attri_1 = excel.ReadStr(readCellValue(row, names, "Attri_1"))
            e.Desc_2 = excel.ReadStr(readCellValue(row, names, "Desc_2"))
            e.DescArg_2 = excel.ReadStr(readCellValue(row, names, "DescArg_2"))
            e.Attri_2 = excel.ReadStr(readCellValue(row, names, "Attri_2"))
            e.Desc_3 = excel.ReadStr(readCellValue(row, names, "Desc_3"))
            e.DescArg_3 = excel.ReadStr(readCellValue(row, names, "DescArg_3"))
            e.Attri_3 = excel.ReadStr(readCellValue(row, names, "Attri_3"))
            e.Desc_4 = excel.ReadStr(readCellValue(row, names, "Desc_4"))
            e.DescArg_4 = excel.ReadStr(readCellValue(row, names, "DescArg_4"))
            e.Attri_4 = excel.ReadStr(readCellValue(row, names, "Attri_4"))
            e.Desc_5 = excel.ReadStr(readCellValue(row, names, "Desc_5"))
            e.DescArg_5 = excel.ReadStr(readCellValue(row, names, "DescArg_5"))
            e.Attri_5 = excel.ReadStr(readCellValue(row, names, "Attri_5"))
            e.Desc_6 = excel.ReadStr(readCellValue(row, names, "Desc_6"))
            e.DescArg_6 = excel.ReadStr(readCellValue(row, names, "DescArg_6"))
            e.Attri_6 = excel.ReadStr(readCellValue(row, names, "Attri_6"))
            e.Desc_7 = excel.ReadStr(readCellValue(row, names, "Desc_7"))
            e.DescArg_7 = excel.ReadStr(readCellValue(row, names, "DescArg_7"))
            e.Attri_7 = excel.ReadStr(readCellValue(row, names, "Attri_7"))
            e.Desc_8 = excel.ReadStr(readCellValue(row, names, "Desc_8"))
            e.DescArg_8 = excel.ReadStr(readCellValue(row, names, "DescArg_8"))
            e.Attri_8 = excel.ReadStr(readCellValue(row, names, "Attri_8"))
            e.Desc_9 = excel.ReadStr(readCellValue(row, names, "Desc_9"))
            e.DescArg_9 = excel.ReadStr(readCellValue(row, names, "DescArg_9"))
            e.Attri_9 = excel.ReadStr(readCellValue(row, names, "Attri_9"))
            e.Desc_10 = excel.ReadStr(readCellValue(row, names, "Desc_10"))
            e.DescArg_10 = excel.ReadStr(readCellValue(row, names, "DescArg_10"))
            e.Attri_10 = excel.ReadStr(readCellValue(row, names, "Attri_10"))
            e.Desc_11 = excel.ReadStr(readCellValue(row, names, "Desc_11"))
            e.DescArg_11 = excel.ReadStr(readCellValue(row, names, "DescArg_11"))
            e.Attri_11 = excel.ReadStr(readCellValue(row, names, "Attri_11"))
            e.Desc_12 = excel.ReadStr(readCellValue(row, names, "Desc_12"))
            e.DescArg_12 = excel.ReadStr(readCellValue(row, names, "DescArg_12"))
            e.Attri_12 = excel.ReadStr(readCellValue(row, names, "Attri_12"))
            e.Desc_13 = excel.ReadStr(readCellValue(row, names, "Desc_13"))
            e.DescArg_13 = excel.ReadStr(readCellValue(row, names, "DescArg_13"))
            e.Attri_13 = excel.ReadStr(readCellValue(row, names, "Attri_13"))
            e.Desc_14 = excel.ReadStr(readCellValue(row, names, "Desc_14"))
            e.DescArg_14 = excel.ReadStr(readCellValue(row, names, "DescArg_14"))
            e.Attri_14 = excel.ReadStr(readCellValue(row, names, "Attri_14"))
            e.Desc_15 = excel.ReadStr(readCellValue(row, names, "Desc_15"))
            e.DescArg_15 = excel.ReadStr(readCellValue(row, names, "DescArg_15"))
            e.Attri_15 = excel.ReadStr(readCellValue(row, names, "Attri_15"))
            e.Damage = readDamage(row, names)

			tmpIDMap[excel.ReadInt32(readCellValue(row, names, "ID"))] = e

		}
	}
	this.indexID.Store(tmpIDMap)	

}
