package FragmentChange

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	path2 "path"
	"strings"
	"sync/atomic"
	"initialthree/node/table/excel"
	"initialthree/zaplogger"
)

type FragmentChange struct{
    ID int32 
    ItemID1 int32 
    ItemCount1 int32 
    ItemID2 int32 
    ItemCount2 int32 
    ItemID3 int32 
    ItemCount3 int32 
    Item []*Item_ 

}

type Item_ struct{
    ID int32 
    Count int32 

}

func readItem(row, names []string)[]*Item_{
	ret := make([]*Item_, 0)
	base := excel.Split("ItemID,ItemCount", ",")

	for i := 1; ; i++ {
		if idx, find := hasCellName(names, fmt.Sprintf("%s%d", base[0][0], i));!find {
			break
		}else if _,ok := hasCellValue(row,idx);!ok{
			break
		}
		
		e := &Item_{}
        e.ID = excel.ReadInt32(readCellValue(row, names, fmt.Sprintf("%s%d", base[0][0], i)))
        e.Count = excel.ReadInt32(readCellValue(row, names, fmt.Sprintf("%s%d", base[1][0], i)))

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

func GetID(key int32)*FragmentChange{
	v := Table_.indexID.Load().(map[int32]*FragmentChange)[key]
	if v == nil{
		zaplogger.GetSugar().Errorf("Table %s GetID %d is nil", Table_.xlsxName, key)
	}
	return v
}

func GetIDMap() map[int32]*FragmentChange {
	return Table_.indexID.Load().(map[int32]*FragmentChange)
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
	Table_.xlsxName = "FragmentChange.xlsx"
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

	tmpIDMap := map[int32]*FragmentChange{}

	for _,row := range rows{
		if row[0] != "" {
			e := &FragmentChange{}
            e.ID = excel.ReadInt32(readCellValue(row, names, "ID"))
            e.ItemID1 = excel.ReadInt32(readCellValue(row, names, "ItemID1"))
            e.ItemCount1 = excel.ReadInt32(readCellValue(row, names, "ItemCount1"))
            e.ItemID2 = excel.ReadInt32(readCellValue(row, names, "ItemID2"))
            e.ItemCount2 = excel.ReadInt32(readCellValue(row, names, "ItemCount2"))
            e.ItemID3 = excel.ReadInt32(readCellValue(row, names, "ItemID3"))
            e.ItemCount3 = excel.ReadInt32(readCellValue(row, names, "ItemCount3"))
            e.Item = readItem(row, names)

			tmpIDMap[excel.ReadInt32(readCellValue(row, names, "ID"))] = e

		}
	}
	this.indexID.Store(tmpIDMap)	

}
