package Function

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	path2 "path"
	"strings"
	"sync/atomic"
	"initialthree/node/table/excel"
	"initialthree/zaplogger"
)

type Function struct{
    ID int32 
    Comment string 
    UnlockType_1 string 
    UnlockArg_1 string 
    UnlockAlert_1 string 
    UnlockType_2 string 
    UnlockArg_2 string 
    UnlockAlert_2 string 
    UnlockType_3 string 
    UnlockArg_3 string 
    UnlockAlert_3 string 
    IsShowUnlockPopUp bool 
    UnlockImg string 
    UnlockDesc string 
    Unlock []*Unlock_ 

}

type Unlock_ struct{
    Type int32 
    Arg string 

}

func readUnlock(row, names []string)[]*Unlock_{
	ret := make([]*Unlock_, 0)
	base := excel.Split("UnlockType_,UnlockArg_", ",")

	for i := 1; ; i++ {
		if idx, find := hasCellName(names, fmt.Sprintf("%s%d", base[0][0], i));!find {
			break
		}else if _,ok := hasCellValue(row,idx);!ok{
			break
		}
		
		e := &Unlock_{}
        e.Type = excel.ReadEnum(readCellValue(row, names, fmt.Sprintf("%s%d", base[0][0], i)))
        e.Arg = excel.ReadStr(readCellValue(row, names, fmt.Sprintf("%s%d", base[1][0], i)))

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

func GetID(key int32)*Function{
	v := Table_.indexID.Load().(map[int32]*Function)[key]
	if v == nil{
		zaplogger.GetSugar().Errorf("Table %s GetID %d is nil", Table_.xlsxName, key)
	}
	return v
}

func GetIDMap() map[int32]*Function {
	return Table_.indexID.Load().(map[int32]*Function)
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
	Table_.xlsxName = "Function.xlsx"
	excel.AddTable(&Table_)
}

func (this *Table)Load(path string){
	xlFile, err := excelize.OpenFile(path2.Join(path2.Join(path, "/DataTable/功能解锁"), this.xlsxName))
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	tableRows := xlFile.GetRows(xlFile.GetSheetName(xlFile.GetActiveSheetIndex()))
	names := tableRows[1]
	rows := tableRows[4:]

	tmpIDMap := map[int32]*Function{}

	for _,row := range rows{
		if row[0] != "" {
			e := &Function{}
            e.ID = excel.ReadInt32(readCellValue(row, names, "ID"))
            e.Comment = excel.ReadStr(readCellValue(row, names, "Comment"))
            e.UnlockType_1 = excel.ReadStr(readCellValue(row, names, "UnlockType_1"))
            e.UnlockArg_1 = excel.ReadStr(readCellValue(row, names, "UnlockArg_1"))
            e.UnlockAlert_1 = excel.ReadStr(readCellValue(row, names, "UnlockAlert_1"))
            e.UnlockType_2 = excel.ReadStr(readCellValue(row, names, "UnlockType_2"))
            e.UnlockArg_2 = excel.ReadStr(readCellValue(row, names, "UnlockArg_2"))
            e.UnlockAlert_2 = excel.ReadStr(readCellValue(row, names, "UnlockAlert_2"))
            e.UnlockType_3 = excel.ReadStr(readCellValue(row, names, "UnlockType_3"))
            e.UnlockArg_3 = excel.ReadStr(readCellValue(row, names, "UnlockArg_3"))
            e.UnlockAlert_3 = excel.ReadStr(readCellValue(row, names, "UnlockAlert_3"))
            e.IsShowUnlockPopUp = excel.ReadBool(readCellValue(row, names, "IsShowUnlockPopUp"))
            e.UnlockImg = excel.ReadStr(readCellValue(row, names, "UnlockImg"))
            e.UnlockDesc = excel.ReadStr(readCellValue(row, names, "UnlockDesc"))
            e.Unlock = readUnlock(row, names)

			tmpIDMap[excel.ReadInt32(readCellValue(row, names, "ID"))] = e

		}
	}
	this.indexID.Store(tmpIDMap)	

}
