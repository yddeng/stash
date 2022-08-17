package PlayerCampReputationShopItem

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	path2 "path"
	"strings"
	"sync/atomic"
	"initialthree/node/table/excel"
	"initialthree/zaplogger"
)

type PlayerCampReputationShopItem struct{
    ID int32 
    CampType string 
    ReputationLevelType string 
    ItemID int32 
    ItemCount int32 
    ProductLimitType string 
    ItemCapacity int32 
    CostItemCount int32 
    CampTypeEnum int32 
    ReputationLevelTypeEnum int32 
    ProductLimitTypeEnum int32 

}

var Table_ Table

type Table struct{
	xlsxName string
	nextPath string
	indexID  atomic.Value

}

func GetID(key int32)*PlayerCampReputationShopItem{
	v := Table_.indexID.Load().(map[int32]*PlayerCampReputationShopItem)[key]
	if v == nil{
		zaplogger.GetSugar().Errorf("Table %s GetID %d is nil", Table_.xlsxName, key)
	}
	return v
}

func GetIDMap() map[int32]*PlayerCampReputationShopItem {
	return Table_.indexID.Load().(map[int32]*PlayerCampReputationShopItem)
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
	Table_.xlsxName = "PlayerCampReputationShopItem.xlsx"
	excel.AddTable(&Table_)
}

func (this *Table)Load(path string){
	xlFile, err := excelize.OpenFile(path2.Join(path2.Join(path, "/DataTable/玩家阵营及声望"), this.xlsxName))
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	tableRows := xlFile.GetRows(xlFile.GetSheetName(xlFile.GetActiveSheetIndex()))
	names := tableRows[1]
	rows := tableRows[4:]

	tmpIDMap := map[int32]*PlayerCampReputationShopItem{}

	for _,row := range rows{
		if row[0] != "" {
			e := &PlayerCampReputationShopItem{}
            e.ID = excel.ReadInt32(readCellValue(row, names, "ID"))
            e.CampType = excel.ReadStr(readCellValue(row, names, "CampType"))
            e.ReputationLevelType = excel.ReadStr(readCellValue(row, names, "ReputationLevelType"))
            e.ItemID = excel.ReadInt32(readCellValue(row, names, "ItemID"))
            e.ItemCount = excel.ReadInt32(readCellValue(row, names, "ItemCount"))
            e.ProductLimitType = excel.ReadStr(readCellValue(row, names, "ProductLimitType"))
            e.ItemCapacity = excel.ReadInt32(readCellValue(row, names, "ItemCapacity"))
            e.CostItemCount = excel.ReadInt32(readCellValue(row, names, "CostItemCount"))
            e.CampTypeEnum = excel.ReadEnum(readCellValue(row, names, "CampType"))
            e.ReputationLevelTypeEnum = excel.ReadEnum(readCellValue(row, names, "ReputationLevelType"))
            e.ProductLimitTypeEnum = excel.ReadEnum(readCellValue(row, names, "ProductLimitType"))

			tmpIDMap[excel.ReadInt32(readCellValue(row, names, "ID"))] = e

		}
	}
	this.indexID.Store(tmpIDMap)	

}
