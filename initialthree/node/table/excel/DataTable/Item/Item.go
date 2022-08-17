package Item

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	path2 "path"
	"strings"
	"sync/atomic"
	"initialthree/node/table/excel"
	"initialthree/zaplogger"
)

type Item struct{
    ID int32 
    Name string 
    Icon string 
    BackpackMenuTag string 
    Summary string 
    Description string 
    AdditionalNote string 
    Type string 
    Rarity string 
    IsConsumable bool 
    AllowUse bool 
    UseRuleID int32 
    AllowUseInBag bool 
    AllowBatchUse bool 
    SkipToUse int32 
    UseCD int32 
    AllowDestroy bool 
    AllowTrade bool 
    AllowSell bool 
    SellCurrencyType string 
    SellPrice int32 
    AllowShowOff bool 
    AllowSplit bool 
    HasTimeLimit bool 
    TimeLimitType string 
    TimeLimit string 
    AccessWayList string 
    DrawCardResultPicture string 
    DrawCardSingleResultPicture string 
    RarityEnum int32 
    TypeEnum int32 
    SellCurrencyTypeEnum int32 
    TimeLimitTypeEnum int32 

}

var Table_ Table

type Table struct{
	xlsxName string
	nextPath string
	indexID  atomic.Value

}

func GetID(key int32)*Item{
	v := Table_.indexID.Load().(map[int32]*Item)[key]
	if v == nil{
		zaplogger.GetSugar().Errorf("Table %s GetID %d is nil", Table_.xlsxName, key)
	}
	return v
}

func GetIDMap() map[int32]*Item {
	return Table_.indexID.Load().(map[int32]*Item)
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
	Table_.xlsxName = "Item.xlsx"
	excel.AddTable(&Table_)
}

func (this *Table)Load(path string){
	xlFile, err := excelize.OpenFile(path2.Join(path2.Join(path, "/DataTable/道具"), this.xlsxName))
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	tableRows := xlFile.GetRows(xlFile.GetSheetName(xlFile.GetActiveSheetIndex()))
	names := tableRows[1]
	rows := tableRows[4:]

	tmpIDMap := map[int32]*Item{}

	for _,row := range rows{
		if row[0] != "" {
			e := &Item{}
            e.ID = excel.ReadInt32(readCellValue(row, names, "ID"))
            e.Name = excel.ReadStr(readCellValue(row, names, "Name"))
            e.Icon = excel.ReadStr(readCellValue(row, names, "Icon"))
            e.BackpackMenuTag = excel.ReadStr(readCellValue(row, names, "BackpackMenuTag"))
            e.Summary = excel.ReadStr(readCellValue(row, names, "Summary"))
            e.Description = excel.ReadStr(readCellValue(row, names, "Description"))
            e.AdditionalNote = excel.ReadStr(readCellValue(row, names, "AdditionalNote"))
            e.Type = excel.ReadStr(readCellValue(row, names, "Type"))
            e.Rarity = excel.ReadStr(readCellValue(row, names, "Rarity"))
            e.IsConsumable = excel.ReadBool(readCellValue(row, names, "IsConsumable"))
            e.AllowUse = excel.ReadBool(readCellValue(row, names, "AllowUse"))
            e.UseRuleID = excel.ReadInt32(readCellValue(row, names, "UseRuleID"))
            e.AllowUseInBag = excel.ReadBool(readCellValue(row, names, "AllowUseInBag"))
            e.AllowBatchUse = excel.ReadBool(readCellValue(row, names, "AllowBatchUse"))
            e.SkipToUse = excel.ReadInt32(readCellValue(row, names, "SkipToUse"))
            e.UseCD = excel.ReadInt32(readCellValue(row, names, "UseCD"))
            e.AllowDestroy = excel.ReadBool(readCellValue(row, names, "AllowDestroy"))
            e.AllowTrade = excel.ReadBool(readCellValue(row, names, "AllowTrade"))
            e.AllowSell = excel.ReadBool(readCellValue(row, names, "AllowSell"))
            e.SellCurrencyType = excel.ReadStr(readCellValue(row, names, "SellCurrencyType"))
            e.SellPrice = excel.ReadInt32(readCellValue(row, names, "SellPrice"))
            e.AllowShowOff = excel.ReadBool(readCellValue(row, names, "AllowShowOff"))
            e.AllowSplit = excel.ReadBool(readCellValue(row, names, "AllowSplit"))
            e.HasTimeLimit = excel.ReadBool(readCellValue(row, names, "HasTimeLimit"))
            e.TimeLimitType = excel.ReadStr(readCellValue(row, names, "TimeLimitType"))
            e.TimeLimit = excel.ReadStr(readCellValue(row, names, "TimeLimit"))
            e.AccessWayList = excel.ReadStr(readCellValue(row, names, "AccessWayList"))
            e.DrawCardResultPicture = excel.ReadStr(readCellValue(row, names, "DrawCardResultPicture"))
            e.DrawCardSingleResultPicture = excel.ReadStr(readCellValue(row, names, "DrawCardSingleResultPicture"))
            e.RarityEnum = excel.ReadEnum(readCellValue(row, names, "Rarity"))
            e.TypeEnum = excel.ReadEnum(readCellValue(row, names, "Type"))
            e.SellCurrencyTypeEnum = excel.ReadEnum(readCellValue(row, names, "SellCurrencyType"))
            e.TimeLimitTypeEnum = excel.ReadEnum(readCellValue(row, names, "TimeLimitType"), "")

			tmpIDMap[excel.ReadInt32(readCellValue(row, names, "ID"))] = e

		}
	}
	this.indexID.Store(tmpIDMap)	

}
