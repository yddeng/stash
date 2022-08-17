package ProductLibrary

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	path2 "path"
	"strings"
	"sync/atomic"
	"initialthree/node/table/excel"
	"initialthree/zaplogger"
)

type ProductLibrary struct{
    ID int32 
    ShopType string 
    FunctionType string 
    ShopProductLibraryLabelID int32 
    Products string 
    PriceType_1 string 
    PriceType_2 string 
    PriceType_3 string 
    RefreshTime string 
    RefreshPriceType string 
    RefreshPrice_1 int32 
    RefreshPrice_2 int32 
    RefreshPrice_3 int32 
    RefreshPrice_4 int32 
    RefreshPrice_5 int32 
    RedreshPrice []*RedreshPrice_ 
    ProductsArray []*Products_ 

}

type RedreshPrice_ struct{
    Price int32 

}

func readRedreshPrice(row, names []string)[]*RedreshPrice_{
	ret := make([]*RedreshPrice_, 0)
	base := excel.Split("RefreshPrice_", ",")

	for i := 1; ; i++ {
		if idx, find := hasCellName(names, fmt.Sprintf("%s%d", base[0][0], i));!find {
			break
		}else if _,ok := hasCellValue(row,idx);!ok{
			break
		}
		
		e := &RedreshPrice_{}
        e.Price = excel.ReadInt32(readCellValue(row, names, fmt.Sprintf("%s%d", base[0][0], i)))

		ret = append(ret, e)
	}

	return ret
}

type Products_ struct{
    ID int32 

}

func readProductsArray(row, names []string)[]*Products_{
	value := readCellValue(row, names, "Products")
	if value == "" {
		return nil
	}

	r := excel.Split(value,",")
	if len(r) == 0 {
		return nil
	}

	ret := make([]*Products_,0)
	for _, v := range r {
		if len(v) == 1{
			e := &Products_{}
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

func GetID(key int32)*ProductLibrary{
	v := Table_.indexID.Load().(map[int32]*ProductLibrary)[key]
	if v == nil{
		zaplogger.GetSugar().Errorf("Table %s GetID %d is nil", Table_.xlsxName, key)
	}
	return v
}

func GetIDMap() map[int32]*ProductLibrary {
	return Table_.indexID.Load().(map[int32]*ProductLibrary)
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
	Table_.xlsxName = "ProductLibrary.xlsx"
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

	tmpIDMap := map[int32]*ProductLibrary{}

	for _,row := range rows{
		if row[0] != "" {
			e := &ProductLibrary{}
            e.ID = excel.ReadInt32(readCellValue(row, names, "ID"))
            e.ShopType = excel.ReadStr(readCellValue(row, names, "ShopType"))
            e.FunctionType = excel.ReadStr(readCellValue(row, names, "FunctionType"))
            e.ShopProductLibraryLabelID = excel.ReadInt32(readCellValue(row, names, "ShopProductLibraryLabelID"))
            e.Products = excel.ReadStr(readCellValue(row, names, "Products"))
            e.PriceType_1 = excel.ReadStr(readCellValue(row, names, "PriceType_1"))
            e.PriceType_2 = excel.ReadStr(readCellValue(row, names, "PriceType_2"))
            e.PriceType_3 = excel.ReadStr(readCellValue(row, names, "PriceType_3"))
            e.RefreshTime = excel.ReadStr(readCellValue(row, names, "RefreshTime"))
            e.RefreshPriceType = excel.ReadStr(readCellValue(row, names, "RefreshPriceType"))
            e.RefreshPrice_1 = excel.ReadInt32(readCellValue(row, names, "RefreshPrice_1"))
            e.RefreshPrice_2 = excel.ReadInt32(readCellValue(row, names, "RefreshPrice_2"))
            e.RefreshPrice_3 = excel.ReadInt32(readCellValue(row, names, "RefreshPrice_3"))
            e.RefreshPrice_4 = excel.ReadInt32(readCellValue(row, names, "RefreshPrice_4"))
            e.RefreshPrice_5 = excel.ReadInt32(readCellValue(row, names, "RefreshPrice_5"))
            e.RedreshPrice = readRedreshPrice(row, names)
            e.ProductsArray = readProductsArray(row, names)

			tmpIDMap[excel.ReadInt32(readCellValue(row, names, "ID"))] = e

		}
	}
	this.indexID.Store(tmpIDMap)	

}
