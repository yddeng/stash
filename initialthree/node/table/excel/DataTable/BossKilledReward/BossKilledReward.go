package BossKilledReward

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	path2 "path"
	"strings"
	"sync/atomic"
	"initialthree/node/table/excel"
	"initialthree/zaplogger"
)

type BossKilledReward struct{
    ID int32 
    WeekdaySUpID_1 string 
    WeekdaySUpID_2 string 
    WeekdaySUpID_3 string 
    WeekdaySUpID_4 string 
    WeekdaySUpID_5 string 
    WeekdaySUpID_6 string 
    WeekdaySUpRate int32 
    WeekdayPoolUp string 
    WeekdayPoolUpRate int32 
    EquipDropID int32 
    ExpDropID int32 
    SecretCoinCount int32 
    SecretCoinWave int32 
    WeekdaySUpID_1Array []*WeekdaySUpID_1_ 
    WeekdaySUpID_2Array []*WeekdaySUpID_2_ 
    WeekdaySUpID_3Array []*WeekdaySUpID_3_ 
    WeekdaySUpID_4Array []*WeekdaySUpID_4_ 
    WeekdaySUpID_5Array []*WeekdaySUpID_5_ 
    WeekdaySUpID_6Array []*WeekdaySUpID_6_ 
    WeekdayPoolUpArray []*WeekdayPoolUp_ 

}

type WeekdaySUpID_1_ struct{
    ID int32 

}

func readWeekdaySUpID_1Array(row, names []string)[]*WeekdaySUpID_1_{
	value := readCellValue(row, names, "WeekdaySUpID_1")
	if value == "" {
		return nil
	}

	r := excel.Split(value,",")
	if len(r) == 0 {
		return nil
	}

	ret := make([]*WeekdaySUpID_1_,0)
	for _, v := range r {
		if len(v) == 1{
			e := &WeekdaySUpID_1_{}
	        e.ID = excel.ReadInt32(v[0])

			ret = append(ret, e)
		}
	}

	return ret
}

type WeekdaySUpID_2_ struct{
    ID int32 

}

func readWeekdaySUpID_2Array(row, names []string)[]*WeekdaySUpID_2_{
	value := readCellValue(row, names, "WeekdaySUpID_2")
	if value == "" {
		return nil
	}

	r := excel.Split(value,",")
	if len(r) == 0 {
		return nil
	}

	ret := make([]*WeekdaySUpID_2_,0)
	for _, v := range r {
		if len(v) == 1{
			e := &WeekdaySUpID_2_{}
	        e.ID = excel.ReadInt32(v[0])

			ret = append(ret, e)
		}
	}

	return ret
}

type WeekdaySUpID_3_ struct{
    ID int32 

}

func readWeekdaySUpID_3Array(row, names []string)[]*WeekdaySUpID_3_{
	value := readCellValue(row, names, "WeekdaySUpID_3")
	if value == "" {
		return nil
	}

	r := excel.Split(value,",")
	if len(r) == 0 {
		return nil
	}

	ret := make([]*WeekdaySUpID_3_,0)
	for _, v := range r {
		if len(v) == 1{
			e := &WeekdaySUpID_3_{}
	        e.ID = excel.ReadInt32(v[0])

			ret = append(ret, e)
		}
	}

	return ret
}

type WeekdaySUpID_4_ struct{
    ID int32 

}

func readWeekdaySUpID_4Array(row, names []string)[]*WeekdaySUpID_4_{
	value := readCellValue(row, names, "WeekdaySUpID_4")
	if value == "" {
		return nil
	}

	r := excel.Split(value,",")
	if len(r) == 0 {
		return nil
	}

	ret := make([]*WeekdaySUpID_4_,0)
	for _, v := range r {
		if len(v) == 1{
			e := &WeekdaySUpID_4_{}
	        e.ID = excel.ReadInt32(v[0])

			ret = append(ret, e)
		}
	}

	return ret
}

type WeekdaySUpID_5_ struct{
    ID int32 

}

func readWeekdaySUpID_5Array(row, names []string)[]*WeekdaySUpID_5_{
	value := readCellValue(row, names, "WeekdaySUpID_5")
	if value == "" {
		return nil
	}

	r := excel.Split(value,",")
	if len(r) == 0 {
		return nil
	}

	ret := make([]*WeekdaySUpID_5_,0)
	for _, v := range r {
		if len(v) == 1{
			e := &WeekdaySUpID_5_{}
	        e.ID = excel.ReadInt32(v[0])

			ret = append(ret, e)
		}
	}

	return ret
}

type WeekdaySUpID_6_ struct{
    ID int32 

}

func readWeekdaySUpID_6Array(row, names []string)[]*WeekdaySUpID_6_{
	value := readCellValue(row, names, "WeekdaySUpID_6")
	if value == "" {
		return nil
	}

	r := excel.Split(value,",")
	if len(r) == 0 {
		return nil
	}

	ret := make([]*WeekdaySUpID_6_,0)
	for _, v := range r {
		if len(v) == 1{
			e := &WeekdaySUpID_6_{}
	        e.ID = excel.ReadInt32(v[0])

			ret = append(ret, e)
		}
	}

	return ret
}

type WeekdayPoolUp_ struct{
    ID int32 

}

func readWeekdayPoolUpArray(row, names []string)[]*WeekdayPoolUp_{
	value := readCellValue(row, names, "WeekdayPoolUp")
	if value == "" {
		return nil
	}

	r := excel.Split(value,",")
	if len(r) == 0 {
		return nil
	}

	ret := make([]*WeekdayPoolUp_,0)
	for _, v := range r {
		if len(v) == 1{
			e := &WeekdayPoolUp_{}
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

func GetID(key int32)*BossKilledReward{
	v := Table_.indexID.Load().(map[int32]*BossKilledReward)[key]
	if v == nil{
		zaplogger.GetSugar().Errorf("Table %s GetID %d is nil", Table_.xlsxName, key)
	}
	return v
}

func GetIDMap() map[int32]*BossKilledReward {
	return Table_.indexID.Load().(map[int32]*BossKilledReward)
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
	Table_.xlsxName = "BossKilledReward.xlsx"
	excel.AddTable(&Table_)
}

func (this *Table)Load(path string){
	xlFile, err := excelize.OpenFile(path2.Join(path2.Join(path, "/DataTable/关卡/秘境"), this.xlsxName))
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	tableRows := xlFile.GetRows(xlFile.GetSheetName(xlFile.GetActiveSheetIndex()))
	names := tableRows[1]
	rows := tableRows[4:]

	tmpIDMap := map[int32]*BossKilledReward{}

	for _,row := range rows{
		if row[0] != "" {
			e := &BossKilledReward{}
            e.ID = excel.ReadInt32(readCellValue(row, names, "ID"))
            e.WeekdaySUpID_1 = excel.ReadStr(readCellValue(row, names, "WeekdaySUpID_1"))
            e.WeekdaySUpID_2 = excel.ReadStr(readCellValue(row, names, "WeekdaySUpID_2"))
            e.WeekdaySUpID_3 = excel.ReadStr(readCellValue(row, names, "WeekdaySUpID_3"))
            e.WeekdaySUpID_4 = excel.ReadStr(readCellValue(row, names, "WeekdaySUpID_4"))
            e.WeekdaySUpID_5 = excel.ReadStr(readCellValue(row, names, "WeekdaySUpID_5"))
            e.WeekdaySUpID_6 = excel.ReadStr(readCellValue(row, names, "WeekdaySUpID_6"))
            e.WeekdaySUpRate = excel.ReadInt32(readCellValue(row, names, "WeekdaySUpRate"))
            e.WeekdayPoolUp = excel.ReadStr(readCellValue(row, names, "WeekdayPoolUp"))
            e.WeekdayPoolUpRate = excel.ReadInt32(readCellValue(row, names, "WeekdayPoolUpRate"))
            e.EquipDropID = excel.ReadInt32(readCellValue(row, names, "EquipDropID"))
            e.ExpDropID = excel.ReadInt32(readCellValue(row, names, "ExpDropID"))
            e.SecretCoinCount = excel.ReadInt32(readCellValue(row, names, "SecretCoinCount"))
            e.SecretCoinWave = excel.ReadInt32(readCellValue(row, names, "SecretCoinWave"))
            e.WeekdaySUpID_1Array = readWeekdaySUpID_1Array(row, names)
            e.WeekdaySUpID_2Array = readWeekdaySUpID_2Array(row, names)
            e.WeekdaySUpID_3Array = readWeekdaySUpID_3Array(row, names)
            e.WeekdaySUpID_4Array = readWeekdaySUpID_4Array(row, names)
            e.WeekdaySUpID_5Array = readWeekdaySUpID_5Array(row, names)
            e.WeekdaySUpID_6Array = readWeekdaySUpID_6Array(row, names)
            e.WeekdayPoolUpArray = readWeekdayPoolUpArray(row, names)

			tmpIDMap[excel.ReadInt32(readCellValue(row, names, "ID"))] = e

		}
	}
	this.indexID.Store(tmpIDMap)	

}
