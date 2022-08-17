package DrawCardsPool

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	path2 "path"
	"strings"
	"sync/atomic"
	"initialthree/node/table/excel"
	"initialthree/zaplogger"
)

type DrawCardsPool struct{
    ID int32 
    AdditionalNote string 
    TitleBg string 
    TenTimesGuaranteePools string 
    GuaranteePool int32 
    IsTargetWeapon bool 
    TargetID string 
    FiveUpProb string 
    FiveUpID string 
    FourUpProb string 
    FourUpID string 
    CardsPoolDesc_1 string 
    CardsPoolID_1 int32 
    CardsPoolWeight_1 int32 
    CardsPoolDesc_2 string 
    CardsPoolID_2 int32 
    CardsPoolWeight_2 int32 
    CardsPoolDesc_3 string 
    CardsPoolID_3 int32 
    CardsPoolWeight_3 int32 
    CardsPoolDesc_4 string 
    CardsPoolID_4 int32 
    CardsPoolWeight_4 int32 
    CardsPoolDesc_5 string 
    CardsPoolID_5 int32 
    CardsPoolWeight_5 int32 
    CardsPoolDesc_6 string 
    CardsPoolID_6 int32 
    CardsPoolWeight_6 int32 
    CardsPoolDesc_7 string 
    CardsPoolID_7 int32 
    CardsPoolWeight_7 int32 
    CardsPoolDesc_8 string 
    CardsPoolID_8 int32 
    CardsPoolWeight_8 int32 
    CardsPoolDesc_9 string 
    CardsPoolID_9 int32 
    CardsPoolWeight_9 int32 
    CardsPoolDesc_10 string 
    CardsPoolID_10 int32 
    CardsPoolWeight_10 int32 
    TenTimesGuaranteePoolsArray []*TenTimesGuaranteePools_ 
    DropList []*DropList_ 

}

type TenTimesGuaranteePools_ struct{
    Idx int32 

}

func readTenTimesGuaranteePoolsArray(row, names []string)[]*TenTimesGuaranteePools_{
	value := readCellValue(row, names, "TenTimesGuaranteePools")
	if value == "" {
		return nil
	}

	r := excel.Split(value,",")
	if len(r) == 0 {
		return nil
	}

	ret := make([]*TenTimesGuaranteePools_,0)
	for _, v := range r {
		if len(v) == 1{
			e := &TenTimesGuaranteePools_{}
	        e.Idx = excel.ReadInt32(v[0])

			ret = append(ret, e)
		}
	}

	return ret
}

type DropList_ struct{
    ID int32 
    Weight int32 

}

func readDropList(row, names []string)[]*DropList_{
	ret := make([]*DropList_, 0)
	base := excel.Split("CardsPoolID_,CardsPoolWeight_", ",")

	for i := 1; ; i++ {
		if idx, find := hasCellName(names, fmt.Sprintf("%s%d", base[0][0], i));!find {
			break
		}else if _,ok := hasCellValue(row,idx);!ok{
			break
		}
		
		e := &DropList_{}
        e.ID = excel.ReadInt32(readCellValue(row, names, fmt.Sprintf("%s%d", base[0][0], i)))
        e.Weight = excel.ReadInt32(readCellValue(row, names, fmt.Sprintf("%s%d", base[1][0], i)))

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

func GetID(key int32)*DrawCardsPool{
	v := Table_.indexID.Load().(map[int32]*DrawCardsPool)[key]
	if v == nil{
		zaplogger.GetSugar().Errorf("Table %s GetID %d is nil", Table_.xlsxName, key)
	}
	return v
}

func GetIDMap() map[int32]*DrawCardsPool {
	return Table_.indexID.Load().(map[int32]*DrawCardsPool)
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
	Table_.xlsxName = "DrawCardsPool.xlsx"
	excel.AddTable(&Table_)
}

func (this *Table)Load(path string){
	xlFile, err := excelize.OpenFile(path2.Join(path2.Join(path, "/DataTable/抽卡"), this.xlsxName))
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	tableRows := xlFile.GetRows(xlFile.GetSheetName(xlFile.GetActiveSheetIndex()))
	names := tableRows[1]
	rows := tableRows[4:]

	tmpIDMap := map[int32]*DrawCardsPool{}

	for _,row := range rows{
		if row[0] != "" {
			e := &DrawCardsPool{}
            e.ID = excel.ReadInt32(readCellValue(row, names, "ID"))
            e.AdditionalNote = excel.ReadStr(readCellValue(row, names, "AdditionalNote"))
            e.TitleBg = excel.ReadStr(readCellValue(row, names, "TitleBg"))
            e.TenTimesGuaranteePools = excel.ReadStr(readCellValue(row, names, "TenTimesGuaranteePools"))
            e.GuaranteePool = excel.ReadInt32(readCellValue(row, names, "GuaranteePool"))
            e.IsTargetWeapon = excel.ReadBool(readCellValue(row, names, "IsTargetWeapon"))
            e.TargetID = excel.ReadStr(readCellValue(row, names, "TargetID"))
            e.FiveUpProb = excel.ReadStr(readCellValue(row, names, "FiveUpProb"))
            e.FiveUpID = excel.ReadStr(readCellValue(row, names, "FiveUpID"))
            e.FourUpProb = excel.ReadStr(readCellValue(row, names, "FourUpProb"))
            e.FourUpID = excel.ReadStr(readCellValue(row, names, "FourUpID"))
            e.CardsPoolDesc_1 = excel.ReadStr(readCellValue(row, names, "CardsPoolDesc_1"))
            e.CardsPoolID_1 = excel.ReadInt32(readCellValue(row, names, "CardsPoolID_1"))
            e.CardsPoolWeight_1 = excel.ReadInt32(readCellValue(row, names, "CardsPoolWeight_1"))
            e.CardsPoolDesc_2 = excel.ReadStr(readCellValue(row, names, "CardsPoolDesc_2"))
            e.CardsPoolID_2 = excel.ReadInt32(readCellValue(row, names, "CardsPoolID_2"))
            e.CardsPoolWeight_2 = excel.ReadInt32(readCellValue(row, names, "CardsPoolWeight_2"))
            e.CardsPoolDesc_3 = excel.ReadStr(readCellValue(row, names, "CardsPoolDesc_3"))
            e.CardsPoolID_3 = excel.ReadInt32(readCellValue(row, names, "CardsPoolID_3"))
            e.CardsPoolWeight_3 = excel.ReadInt32(readCellValue(row, names, "CardsPoolWeight_3"))
            e.CardsPoolDesc_4 = excel.ReadStr(readCellValue(row, names, "CardsPoolDesc_4"))
            e.CardsPoolID_4 = excel.ReadInt32(readCellValue(row, names, "CardsPoolID_4"))
            e.CardsPoolWeight_4 = excel.ReadInt32(readCellValue(row, names, "CardsPoolWeight_4"))
            e.CardsPoolDesc_5 = excel.ReadStr(readCellValue(row, names, "CardsPoolDesc_5"))
            e.CardsPoolID_5 = excel.ReadInt32(readCellValue(row, names, "CardsPoolID_5"))
            e.CardsPoolWeight_5 = excel.ReadInt32(readCellValue(row, names, "CardsPoolWeight_5"))
            e.CardsPoolDesc_6 = excel.ReadStr(readCellValue(row, names, "CardsPoolDesc_6"))
            e.CardsPoolID_6 = excel.ReadInt32(readCellValue(row, names, "CardsPoolID_6"))
            e.CardsPoolWeight_6 = excel.ReadInt32(readCellValue(row, names, "CardsPoolWeight_6"))
            e.CardsPoolDesc_7 = excel.ReadStr(readCellValue(row, names, "CardsPoolDesc_7"))
            e.CardsPoolID_7 = excel.ReadInt32(readCellValue(row, names, "CardsPoolID_7"))
            e.CardsPoolWeight_7 = excel.ReadInt32(readCellValue(row, names, "CardsPoolWeight_7"))
            e.CardsPoolDesc_8 = excel.ReadStr(readCellValue(row, names, "CardsPoolDesc_8"))
            e.CardsPoolID_8 = excel.ReadInt32(readCellValue(row, names, "CardsPoolID_8"))
            e.CardsPoolWeight_8 = excel.ReadInt32(readCellValue(row, names, "CardsPoolWeight_8"))
            e.CardsPoolDesc_9 = excel.ReadStr(readCellValue(row, names, "CardsPoolDesc_9"))
            e.CardsPoolID_9 = excel.ReadInt32(readCellValue(row, names, "CardsPoolID_9"))
            e.CardsPoolWeight_9 = excel.ReadInt32(readCellValue(row, names, "CardsPoolWeight_9"))
            e.CardsPoolDesc_10 = excel.ReadStr(readCellValue(row, names, "CardsPoolDesc_10"))
            e.CardsPoolID_10 = excel.ReadInt32(readCellValue(row, names, "CardsPoolID_10"))
            e.CardsPoolWeight_10 = excel.ReadInt32(readCellValue(row, names, "CardsPoolWeight_10"))
            e.TenTimesGuaranteePoolsArray = readTenTimesGuaranteePoolsArray(row, names)
            e.DropList = readDropList(row, names)

			tmpIDMap[excel.ReadInt32(readCellValue(row, names, "ID"))] = e

		}
	}
	this.indexID.Store(tmpIDMap)	

}
