package DrawCardsLib

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	path2 "path"
	"strings"
	"sync/atomic"
	"initialthree/node/table/excel"
	"initialthree/zaplogger"
)

type DrawCardsLib struct{
    ID int32 
    AdditionalNote string 
    Name string 
    UnselectTabBg string 
    DropDetail int32 
    RuleDesc int32 
    ConsumeTokenType int32 
    ConsumeTokenCount int32 
    LibOpenType int32 
    LibOpenTime string 
    LibEndTime string 
    GuaranteeID int32 
    GuaranteeThreshold int32 
    GuaranteeIncrease int32 
    IsShowTargetSelector bool 
    DrawCardsPool string 
    DrawCardsPoolArray []*DrawCardsPool_ 

}

type DrawCardsPool_ struct{
    PoolID int32 

}

func readDrawCardsPoolArray(row, names []string)[]*DrawCardsPool_{
	value := readCellValue(row, names, "DrawCardsPool")
	if value == "" {
		return nil
	}

	r := excel.Split(value,",")
	if len(r) == 0 {
		return nil
	}

	ret := make([]*DrawCardsPool_,0)
	for _, v := range r {
		if len(v) == 1{
			e := &DrawCardsPool_{}
	        e.PoolID = excel.ReadInt32(v[0])

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

func GetID(key int32)*DrawCardsLib{
	v := Table_.indexID.Load().(map[int32]*DrawCardsLib)[key]
	if v == nil{
		zaplogger.GetSugar().Errorf("Table %s GetID %d is nil", Table_.xlsxName, key)
	}
	return v
}

func GetIDMap() map[int32]*DrawCardsLib {
	return Table_.indexID.Load().(map[int32]*DrawCardsLib)
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
	Table_.xlsxName = "DrawCardsLib.xlsx"
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

	tmpIDMap := map[int32]*DrawCardsLib{}

	for _,row := range rows{
		if row[0] != "" {
			e := &DrawCardsLib{}
            e.ID = excel.ReadInt32(readCellValue(row, names, "ID"))
            e.AdditionalNote = excel.ReadStr(readCellValue(row, names, "AdditionalNote"))
            e.Name = excel.ReadStr(readCellValue(row, names, "Name"))
            e.UnselectTabBg = excel.ReadStr(readCellValue(row, names, "UnselectTabBg"))
            e.DropDetail = excel.ReadInt32(readCellValue(row, names, "DropDetail"))
            e.RuleDesc = excel.ReadInt32(readCellValue(row, names, "RuleDesc"))
            e.ConsumeTokenType = excel.ReadInt32(readCellValue(row, names, "ConsumeTokenType"))
            e.ConsumeTokenCount = excel.ReadInt32(readCellValue(row, names, "ConsumeTokenCount"))
            e.LibOpenType = excel.ReadInt32(readCellValue(row, names, "LibOpenType"))
            e.LibOpenTime = excel.ReadStr(readCellValue(row, names, "LibOpenTime"))
            e.LibEndTime = excel.ReadStr(readCellValue(row, names, "LibEndTime"))
            e.GuaranteeID = excel.ReadInt32(readCellValue(row, names, "GuaranteeID"))
            e.GuaranteeThreshold = excel.ReadInt32(readCellValue(row, names, "GuaranteeThreshold"))
            e.GuaranteeIncrease = excel.ReadInt32(readCellValue(row, names, "GuaranteeIncrease"))
            e.IsShowTargetSelector = excel.ReadBool(readCellValue(row, names, "IsShowTargetSelector"))
            e.DrawCardsPool = excel.ReadStr(readCellValue(row, names, "DrawCardsPool"))
            e.DrawCardsPoolArray = readDrawCardsPoolArray(row, names)

			tmpIDMap[excel.ReadInt32(readCellValue(row, names, "ID"))] = e

		}
	}
	this.indexID.Store(tmpIDMap)	

}
