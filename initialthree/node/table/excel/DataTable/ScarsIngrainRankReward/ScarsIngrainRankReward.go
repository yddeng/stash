package ScarsIngrainRankReward

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	path2 "path"
	"strings"
	"sync/atomic"
	"initialthree/node/table/excel"
	"initialthree/zaplogger"
)

type ScarsIngrainRankReward struct{
    ID int32 
    RankStartPerc_1 float64 
    RankEndPerc_1 float64 
    DropPoolID_1 int32 
    RankStartPerc_2 float64 
    RankEndPerc_2 float64 
    DropPoolID_2 int32 
    RankStartPerc_3 float64 
    RankEndPerc_3 float64 
    DropPoolID_3 int32 
    RankStartPerc_4 float64 
    RankEndPerc_4 float64 
    DropPoolID_4 int32 
    RankStartPerc_5 float64 
    RankEndPerc_5 float64 
    DropPoolID_5 int32 
    RankReward []*RankReward_ 

}

type RankReward_ struct{
    Start float64 
    End float64 
    DropPoolID int32 

}

func readRankReward(row, names []string)[]*RankReward_{
	ret := make([]*RankReward_, 0)
	base := excel.Split("RankStartPerc_,RankEndPerc_,DropPoolID_", ",")

	for i := 1; ; i++ {
		if idx, find := hasCellName(names, fmt.Sprintf("%s%d", base[0][0], i));!find {
			break
		}else if _,ok := hasCellValue(row,idx);!ok{
			break
		}
		
		e := &RankReward_{}
        e.Start = excel.ReadFloat(readCellValue(row, names, fmt.Sprintf("%s%d", base[0][0], i)))
        e.End = excel.ReadFloat(readCellValue(row, names, fmt.Sprintf("%s%d", base[1][0], i)))
        e.DropPoolID = excel.ReadInt32(readCellValue(row, names, fmt.Sprintf("%s%d", base[2][0], i)))

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

func GetID(key int32)*ScarsIngrainRankReward{
	v := Table_.indexID.Load().(map[int32]*ScarsIngrainRankReward)[key]
	if v == nil{
		zaplogger.GetSugar().Errorf("Table %s GetID %d is nil", Table_.xlsxName, key)
	}
	return v
}

func GetIDMap() map[int32]*ScarsIngrainRankReward {
	return Table_.indexID.Load().(map[int32]*ScarsIngrainRankReward)
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
	Table_.xlsxName = "ScarsIngrainRankReward.xlsx"
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

	tmpIDMap := map[int32]*ScarsIngrainRankReward{}

	for _,row := range rows{
		if row[0] != "" {
			e := &ScarsIngrainRankReward{}
            e.ID = excel.ReadInt32(readCellValue(row, names, "ID"))
            e.RankStartPerc_1 = excel.ReadFloat(readCellValue(row, names, "RankStartPerc_1"))
            e.RankEndPerc_1 = excel.ReadFloat(readCellValue(row, names, "RankEndPerc_1"))
            e.DropPoolID_1 = excel.ReadInt32(readCellValue(row, names, "DropPoolID_1"))
            e.RankStartPerc_2 = excel.ReadFloat(readCellValue(row, names, "RankStartPerc_2"))
            e.RankEndPerc_2 = excel.ReadFloat(readCellValue(row, names, "RankEndPerc_2"))
            e.DropPoolID_2 = excel.ReadInt32(readCellValue(row, names, "DropPoolID_2"))
            e.RankStartPerc_3 = excel.ReadFloat(readCellValue(row, names, "RankStartPerc_3"))
            e.RankEndPerc_3 = excel.ReadFloat(readCellValue(row, names, "RankEndPerc_3"))
            e.DropPoolID_3 = excel.ReadInt32(readCellValue(row, names, "DropPoolID_3"))
            e.RankStartPerc_4 = excel.ReadFloat(readCellValue(row, names, "RankStartPerc_4"))
            e.RankEndPerc_4 = excel.ReadFloat(readCellValue(row, names, "RankEndPerc_4"))
            e.DropPoolID_4 = excel.ReadInt32(readCellValue(row, names, "DropPoolID_4"))
            e.RankStartPerc_5 = excel.ReadFloat(readCellValue(row, names, "RankStartPerc_5"))
            e.RankEndPerc_5 = excel.ReadFloat(readCellValue(row, names, "RankEndPerc_5"))
            e.DropPoolID_5 = excel.ReadInt32(readCellValue(row, names, "DropPoolID_5"))
            e.RankReward = readRankReward(row, names)

			tmpIDMap[excel.ReadInt32(readCellValue(row, names, "ID"))] = e

		}
	}
	this.indexID.Store(tmpIDMap)	

}
