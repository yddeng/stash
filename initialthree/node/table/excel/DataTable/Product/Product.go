package Product

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	path2 "path"
	"strings"
	"sync/atomic"
	"initialthree/node/table/excel"
	"initialthree/zaplogger"
)

type Product struct{
    ID int32 
    PType string 
    PID int32 
    PCount int32 
    PriceType string 
    Price int32 
    LibID int32 
    ProductLimitType string 
    LimitDateStart string 
    LimitDateEnd string 
    LimitCount int32 
    Discount float64 
    PTypeEnum int32 
    ProductLimitTypeEnum int32 

}

var Table_ Table

type Table struct{
	xlsxName string
	nextPath string
	indexID  atomic.Value

}

func GetID(key int32)*Product{
	v := Table_.indexID.Load().(map[int32]*Product)[key]
	if v == nil{
		zaplogger.GetSugar().Errorf("Table %s GetID %d is nil", Table_.xlsxName, key)
	}
	return v
}

func GetIDMap() map[int32]*Product {
	return Table_.indexID.Load().(map[int32]*Product)
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
	Table_.xlsxName = "Product.xlsx"
	excel.AddTable(&Table_)
}

func (this *Table)Load(path string){
	xlFile, err := excelize.OpenFile(path2.Join(path2.Join(path, "/DataTable/商店"), this.xlsxName))
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	tableRows := xlFile.GetRows(xlFile.GetSheetName(xlFile.GetActiveSheetIndex()))
	names := tableRows[1]
	rows := tableRows[4:]

	tmpIDMap := map[int32]*Product{}

	for _,row := range rows{
		if row[0] != "" {
			e := &Product{}
            e.ID = excel.ReadInt32(readCellValue(row, names, "ID"))
            e.PType = excel.ReadStr(readCellValue(row, names, "PType"))
            e.PID = excel.ReadInt32(readCellValue(row, names, "PID"))
            e.PCount = excel.ReadInt32(readCellValue(row, names, "PCount"))
            e.PriceType = excel.ReadStr(readCellValue(row, names, "PriceType"))
            e.Price = excel.ReadInt32(readCellValue(row, names, "Price"))
            e.LibID = excel.ReadInt32(readCellValue(row, names, "LibID"))
            e.ProductLimitType = excel.ReadStr(readCellValue(row, names, "ProductLimitType"))
            e.LimitDateStart = excel.ReadStr(readCellValue(row, names, "LimitDateStart"))
            e.LimitDateEnd = excel.ReadStr(readCellValue(row, names, "LimitDateEnd"))
            e.LimitCount = excel.ReadInt32(readCellValue(row, names, "LimitCount"))
            e.Discount = excel.ReadFloat(readCellValue(row, names, "Discount"))
            e.PTypeEnum = excel.ReadEnum(readCellValue(row, names, "PType"))
            e.ProductLimitTypeEnum = excel.ReadEnum(readCellValue(row, names, "ProductLimitType"))

			tmpIDMap[excel.ReadInt32(readCellValue(row, names, "ID"))] = e

		}
	}
	this.indexID.Store(tmpIDMap)	

}
