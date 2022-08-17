package BigSecret

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	path2 "path"
	"strings"
	"sync/atomic"
	"initialthree/node/table/excel"
	"initialthree/zaplogger"
)

type BigSecret struct{
    ID int32 
    CostFatigue int32 
    MaxKeyCount int32 
    MaxLv int32 
    NormalMonsterProgressValue int32 
    EliteMonsterProgressValue int32 
    MaxProgress int32 
    NormalAffixLib int32 
    AdvancedlAffixLib int32 
    HaloAffixLib int32 
    LinkageAffixLib int32 
    Weakness string 
    WeaknessRefreshCost string 
    PassTimeLimit_1 int32 
    UnlockLv_1 int32 
    PassTimeLimit_2 int32 
    UnlockLv_2 int32 
    PassTimeLimit_3 int32 
    UnlockLv_3 int32 
    PassTimeLimit_4 int32 
    UnlockLv_4 int32 
    BlessingUnlockLv int32 
    ResetLv int32 
    PassTimeUnlock []*PassTimeUnlock_ 
    WeaknessArray []*Weakness_ 
    WeaknessRefreshCostArray []*WeaknessRefreshCost_ 

}

type PassTimeUnlock_ struct{
    PassTime int32 
    Unlock int32 

}

func readPassTimeUnlock(row, names []string)[]*PassTimeUnlock_{
	ret := make([]*PassTimeUnlock_, 0)
	base := excel.Split("PassTimeLimit_,UnlockLv_", ",")

	for i := 1; ; i++ {
		if idx, find := hasCellName(names, fmt.Sprintf("%s%d", base[0][0], i));!find {
			break
		}else if _,ok := hasCellValue(row,idx);!ok{
			break
		}
		
		e := &PassTimeUnlock_{}
        e.PassTime = excel.ReadInt32(readCellValue(row, names, fmt.Sprintf("%s%d", base[0][0], i)))
        e.Unlock = excel.ReadInt32(readCellValue(row, names, fmt.Sprintf("%s%d", base[1][0], i)))

		ret = append(ret, e)
	}

	return ret
}

type Weakness_ struct{
    ID int32 

}

func readWeaknessArray(row, names []string)[]*Weakness_{
	value := readCellValue(row, names, "Weakness")
	if value == "" {
		return nil
	}

	r := excel.Split(value,",")
	if len(r) == 0 {
		return nil
	}

	ret := make([]*Weakness_,0)
	for _, v := range r {
		if len(v) == 1{
			e := &Weakness_{}
	        e.ID = excel.ReadInt32(v[0])

			ret = append(ret, e)
		}
	}

	return ret
}

type WeaknessRefreshCost_ struct{
    Cost int32 

}

func readWeaknessRefreshCostArray(row, names []string)[]*WeaknessRefreshCost_{
	value := readCellValue(row, names, "WeaknessRefreshCost")
	if value == "" {
		return nil
	}

	r := excel.Split(value,",")
	if len(r) == 0 {
		return nil
	}

	ret := make([]*WeaknessRefreshCost_,0)
	for _, v := range r {
		if len(v) == 1{
			e := &WeaknessRefreshCost_{}
	        e.Cost = excel.ReadInt32(v[0])

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

func GetID(key int32)*BigSecret{
	v := Table_.indexID.Load().(map[int32]*BigSecret)[key]
	if v == nil{
		zaplogger.GetSugar().Errorf("Table %s GetID %d is nil", Table_.xlsxName, key)
	}
	return v
}

func GetIDMap() map[int32]*BigSecret {
	return Table_.indexID.Load().(map[int32]*BigSecret)
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
	Table_.xlsxName = "BigSecret.xlsx"
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

	tmpIDMap := map[int32]*BigSecret{}

	for _,row := range rows{
		if row[0] != "" {
			e := &BigSecret{}
            e.ID = excel.ReadInt32(readCellValue(row, names, "ID"))
            e.CostFatigue = excel.ReadInt32(readCellValue(row, names, "CostFatigue"))
            e.MaxKeyCount = excel.ReadInt32(readCellValue(row, names, "MaxKeyCount"))
            e.MaxLv = excel.ReadInt32(readCellValue(row, names, "MaxLv"))
            e.NormalMonsterProgressValue = excel.ReadInt32(readCellValue(row, names, "NormalMonsterProgressValue"))
            e.EliteMonsterProgressValue = excel.ReadInt32(readCellValue(row, names, "EliteMonsterProgressValue"))
            e.MaxProgress = excel.ReadInt32(readCellValue(row, names, "MaxProgress"))
            e.NormalAffixLib = excel.ReadInt32(readCellValue(row, names, "NormalAffixLib"))
            e.AdvancedlAffixLib = excel.ReadInt32(readCellValue(row, names, "AdvancedlAffixLib"))
            e.HaloAffixLib = excel.ReadInt32(readCellValue(row, names, "HaloAffixLib"))
            e.LinkageAffixLib = excel.ReadInt32(readCellValue(row, names, "LinkageAffixLib"))
            e.Weakness = excel.ReadStr(readCellValue(row, names, "Weakness"))
            e.WeaknessRefreshCost = excel.ReadStr(readCellValue(row, names, "WeaknessRefreshCost"))
            e.PassTimeLimit_1 = excel.ReadInt32(readCellValue(row, names, "PassTimeLimit_1"))
            e.UnlockLv_1 = excel.ReadInt32(readCellValue(row, names, "UnlockLv_1"))
            e.PassTimeLimit_2 = excel.ReadInt32(readCellValue(row, names, "PassTimeLimit_2"))
            e.UnlockLv_2 = excel.ReadInt32(readCellValue(row, names, "UnlockLv_2"))
            e.PassTimeLimit_3 = excel.ReadInt32(readCellValue(row, names, "PassTimeLimit_3"))
            e.UnlockLv_3 = excel.ReadInt32(readCellValue(row, names, "UnlockLv_3"))
            e.PassTimeLimit_4 = excel.ReadInt32(readCellValue(row, names, "PassTimeLimit_4"))
            e.UnlockLv_4 = excel.ReadInt32(readCellValue(row, names, "UnlockLv_4"))
            e.BlessingUnlockLv = excel.ReadInt32(readCellValue(row, names, "BlessingUnlockLv"))
            e.ResetLv = excel.ReadInt32(readCellValue(row, names, "ResetLv"))
            e.PassTimeUnlock = readPassTimeUnlock(row, names)
            e.WeaknessArray = readWeaknessArray(row, names)
            e.WeaknessRefreshCostArray = readWeaknessRefreshCostArray(row, names)

			tmpIDMap[excel.ReadInt32(readCellValue(row, names, "ID"))] = e

		}
	}
	this.indexID.Store(tmpIDMap)	

}
