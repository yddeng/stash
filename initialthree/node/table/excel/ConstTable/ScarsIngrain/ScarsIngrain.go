package ScarsIngrain

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	path2 "path"
	"strings"
	"sync/atomic"
	"initialthree/node/table/excel"
	"initialthree/zaplogger"
)

type ScarsIngrain struct{
    ID int32 
    RoleChallengeTimes int32 
    BeginTime string 
    EndTime string 
    BossOpenTime1 string 
    BossOpenTime2 string 
    BossOpenTime3 string 
    ScoreAddition float64 
    SIcon string 
    AIcon string 
    BIcon string 
    CIcon string 
    First string 
    Second string 
    Third string 
    BeginTimeStruct *BeginTime_ 
    EndTimeStruct *EndTime_ 
    BossOpenTime1Struct *BossOpenTime1_ 
    BossOpenTime2Struct *BossOpenTime2_ 
    BossOpenTime3Struct *BossOpenTime3_ 

}

type BeginTime_ struct{
    Weekly int32 
    Hour int32 
    Minute int32 

}

func readBeginTimeStruct(row, names []string)*BeginTime_{
	value := readCellValue(row, names, "BeginTime")
	r := excel.Split(value,",")
	
	e := &BeginTime_{}
    e.Weekly = excel.ReadInt32(r[0][0])
    e.Hour = excel.ReadInt32(r[1][0])
    e.Minute = excel.ReadInt32(r[2][0])

	return e
}

type EndTime_ struct{
    Weekly int32 
    Hour int32 
    Minute int32 

}

func readEndTimeStruct(row, names []string)*EndTime_{
	value := readCellValue(row, names, "EndTime")
	r := excel.Split(value,",")
	
	e := &EndTime_{}
    e.Weekly = excel.ReadInt32(r[0][0])
    e.Hour = excel.ReadInt32(r[1][0])
    e.Minute = excel.ReadInt32(r[2][0])

	return e
}

type BossOpenTime1_ struct{
    Weekly int32 
    Hour int32 
    Minute int32 

}

func readBossOpenTime1Struct(row, names []string)*BossOpenTime1_{
	value := readCellValue(row, names, "BossOpenTime1")
	r := excel.Split(value,",")
	
	e := &BossOpenTime1_{}
    e.Weekly = excel.ReadInt32(r[0][0])
    e.Hour = excel.ReadInt32(r[1][0])
    e.Minute = excel.ReadInt32(r[2][0])

	return e
}

type BossOpenTime2_ struct{
    Weekly int32 
    Hour int32 
    Minute int32 

}

func readBossOpenTime2Struct(row, names []string)*BossOpenTime2_{
	value := readCellValue(row, names, "BossOpenTime2")
	r := excel.Split(value,",")
	
	e := &BossOpenTime2_{}
    e.Weekly = excel.ReadInt32(r[0][0])
    e.Hour = excel.ReadInt32(r[1][0])
    e.Minute = excel.ReadInt32(r[2][0])

	return e
}

type BossOpenTime3_ struct{
    Weekly int32 
    Hour int32 
    Minute int32 

}

func readBossOpenTime3Struct(row, names []string)*BossOpenTime3_{
	value := readCellValue(row, names, "BossOpenTime3")
	r := excel.Split(value,",")
	
	e := &BossOpenTime3_{}
    e.Weekly = excel.ReadInt32(r[0][0])
    e.Hour = excel.ReadInt32(r[1][0])
    e.Minute = excel.ReadInt32(r[2][0])

	return e
}

var Table_ Table

type Table struct{
	xlsxName string
	nextPath string
	indexID  atomic.Value

}

func GetID(key int32)*ScarsIngrain{
	v := Table_.indexID.Load().(map[int32]*ScarsIngrain)[key]
	if v == nil{
		zaplogger.GetSugar().Errorf("Table %s GetID %d is nil", Table_.xlsxName, key)
	}
	return v
}

func GetIDMap() map[int32]*ScarsIngrain {
	return Table_.indexID.Load().(map[int32]*ScarsIngrain)
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
	Table_.xlsxName = "ScarsIngrain.xlsx"
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

	tmpIDMap := map[int32]*ScarsIngrain{}

	for _,row := range rows{
		if row[0] != "" {
			e := &ScarsIngrain{}
            e.ID = excel.ReadInt32(readCellValue(row, names, "ID"))
            e.RoleChallengeTimes = excel.ReadInt32(readCellValue(row, names, "RoleChallengeTimes"))
            e.BeginTime = excel.ReadStr(readCellValue(row, names, "BeginTime"))
            e.EndTime = excel.ReadStr(readCellValue(row, names, "EndTime"))
            e.BossOpenTime1 = excel.ReadStr(readCellValue(row, names, "BossOpenTime1"))
            e.BossOpenTime2 = excel.ReadStr(readCellValue(row, names, "BossOpenTime2"))
            e.BossOpenTime3 = excel.ReadStr(readCellValue(row, names, "BossOpenTime3"))
            e.ScoreAddition = excel.ReadFloat(readCellValue(row, names, "ScoreAddition"))
            e.SIcon = excel.ReadStr(readCellValue(row, names, "SIcon"))
            e.AIcon = excel.ReadStr(readCellValue(row, names, "AIcon"))
            e.BIcon = excel.ReadStr(readCellValue(row, names, "BIcon"))
            e.CIcon = excel.ReadStr(readCellValue(row, names, "CIcon"))
            e.First = excel.ReadStr(readCellValue(row, names, "First"))
            e.Second = excel.ReadStr(readCellValue(row, names, "Second"))
            e.Third = excel.ReadStr(readCellValue(row, names, "Third"))
            e.BeginTimeStruct = readBeginTimeStruct(row, names)
            e.EndTimeStruct = readEndTimeStruct(row, names)
            e.BossOpenTime1Struct = readBossOpenTime1Struct(row, names)
            e.BossOpenTime2Struct = readBossOpenTime2Struct(row, names)
            e.BossOpenTime3Struct = readBossOpenTime3Struct(row, names)

			tmpIDMap[excel.ReadInt32(readCellValue(row, names, "ID"))] = e

		}
	}
	this.indexID.Store(tmpIDMap)	

}
