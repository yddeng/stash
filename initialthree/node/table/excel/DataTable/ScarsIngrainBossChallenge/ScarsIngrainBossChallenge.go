package ScarsIngrainBossChallenge

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	path2 "path"
	"strings"
	"sync/atomic"
	"initialthree/node/table/excel"
	"initialthree/zaplogger"
)

type ScarsIngrainBossChallenge struct{
    ID int32 
    ChallengeType string 
    ArgLeft_1 string 
    ArgRight_1 string 
    Score_1 int32 
    Buff_1 int32 
    ArgLeft_2 string 
    ArgRight_2 string 
    Score_2 int32 
    Buff_2 int32 
    ArgLeft_3 string 
    ArgRight_3 string 
    Score_3 int32 
    Buff_3 int32 
    ArgLeft_4 string 
    ArgRight_4 string 
    Score_4 int32 
    Buff_4 int32 
    ChallengeTypeEnum int32 
    Challenge []*Challenge_ 

}

type Challenge_ struct{
    ArgLeft int32 
    ArgRight int32 
    Score int32 
    Buff int32 

}

func readChallenge(row, names []string)[]*Challenge_{
	ret := make([]*Challenge_, 0)
	base := excel.Split("ArgLeft_,ArgRight_,Score_,Buff_", ",")

	for i := 1; ; i++ {
		if idx, find := hasCellName(names, fmt.Sprintf("%s%d", base[0][0], i));!find {
			break
		}else if _,ok := hasCellValue(row,idx);!ok{
			break
		}
		
		e := &Challenge_{}
        e.ArgLeft = excel.ReadInt32(readCellValue(row, names, fmt.Sprintf("%s%d", base[0][0], i)))
        e.ArgRight = excel.ReadInt32(readCellValue(row, names, fmt.Sprintf("%s%d", base[1][0], i)))
        e.Score = excel.ReadInt32(readCellValue(row, names, fmt.Sprintf("%s%d", base[2][0], i)))
        e.Buff = excel.ReadInt32(readCellValue(row, names, fmt.Sprintf("%s%d", base[3][0], i)))

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

func GetID(key int32)*ScarsIngrainBossChallenge{
	v := Table_.indexID.Load().(map[int32]*ScarsIngrainBossChallenge)[key]
	if v == nil{
		zaplogger.GetSugar().Errorf("Table %s GetID %d is nil", Table_.xlsxName, key)
	}
	return v
}

func GetIDMap() map[int32]*ScarsIngrainBossChallenge {
	return Table_.indexID.Load().(map[int32]*ScarsIngrainBossChallenge)
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
	Table_.xlsxName = "ScarsIngrainBossChallenge.xlsx"
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

	tmpIDMap := map[int32]*ScarsIngrainBossChallenge{}

	for _,row := range rows{
		if row[0] != "" {
			e := &ScarsIngrainBossChallenge{}
            e.ID = excel.ReadInt32(readCellValue(row, names, "ID"))
            e.ChallengeType = excel.ReadStr(readCellValue(row, names, "ChallengeType"))
            e.ArgLeft_1 = excel.ReadStr(readCellValue(row, names, "ArgLeft_1"))
            e.ArgRight_1 = excel.ReadStr(readCellValue(row, names, "ArgRight_1"))
            e.Score_1 = excel.ReadInt32(readCellValue(row, names, "Score_1"))
            e.Buff_1 = excel.ReadInt32(readCellValue(row, names, "Buff_1"))
            e.ArgLeft_2 = excel.ReadStr(readCellValue(row, names, "ArgLeft_2"))
            e.ArgRight_2 = excel.ReadStr(readCellValue(row, names, "ArgRight_2"))
            e.Score_2 = excel.ReadInt32(readCellValue(row, names, "Score_2"))
            e.Buff_2 = excel.ReadInt32(readCellValue(row, names, "Buff_2"))
            e.ArgLeft_3 = excel.ReadStr(readCellValue(row, names, "ArgLeft_3"))
            e.ArgRight_3 = excel.ReadStr(readCellValue(row, names, "ArgRight_3"))
            e.Score_3 = excel.ReadInt32(readCellValue(row, names, "Score_3"))
            e.Buff_3 = excel.ReadInt32(readCellValue(row, names, "Buff_3"))
            e.ArgLeft_4 = excel.ReadStr(readCellValue(row, names, "ArgLeft_4"))
            e.ArgRight_4 = excel.ReadStr(readCellValue(row, names, "ArgRight_4"))
            e.Score_4 = excel.ReadInt32(readCellValue(row, names, "Score_4"))
            e.Buff_4 = excel.ReadInt32(readCellValue(row, names, "Buff_4"))
            e.ChallengeTypeEnum = excel.ReadEnum(readCellValue(row, names, "ChallengeType"))
            e.Challenge = readChallenge(row, names)

			tmpIDMap[excel.ReadInt32(readCellValue(row, names, "ID"))] = e

		}
	}
	this.indexID.Store(tmpIDMap)	

}
