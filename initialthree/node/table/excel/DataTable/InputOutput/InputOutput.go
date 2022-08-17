package InputOutput

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	path2 "path"
	"strings"
	"sync/atomic"
	"initialthree/node/table/excel"
	"initialthree/zaplogger"
)

type InputOutput struct{
    ID int32 
    Describe string 
    InputType_1 string 
    InputID_1 int32 
    InputCount_1 int32 
    InputType_2 string 
    InputID_2 int32 
    InputCount_2 int32 
    InputType_3 string 
    InputID_3 int32 
    InputCount_3 int32 
    InputType_4 string 
    InputID_4 int32 
    InputCount_4 int32 
    OutputType_1 string 
    OutputID_1 int32 
    OutputCount_1 int32 
    OutputType_2 string 
    OutputID_2 int32 
    OutputCount_2 int32 
    OutputType_3 string 
    OutputID_3 int32 
    OutputCount_3 int32 
    OutputType_4 string 
    OutputID_4 int32 
    OutputCount_4 int32 
    Input []*Input_ 
    Output []*Output_ 

}

type Input_ struct{
    Type int32 
    ID int32 
    Count int32 

}

func readInput(row, names []string)[]*Input_{
	ret := make([]*Input_, 0)
	base := excel.Split("InputType_,InputID_,InputCount_", ",")

	for i := 1; ; i++ {
		if idx, find := hasCellName(names, fmt.Sprintf("%s%d", base[0][0], i));!find {
			break
		}else if _,ok := hasCellValue(row,idx);!ok{
			break
		}
		
		e := &Input_{}
        e.Type = excel.ReadEnum(readCellValue(row, names, fmt.Sprintf("%s%d", base[0][0], i)))
        e.ID = excel.ReadInt32(readCellValue(row, names, fmt.Sprintf("%s%d", base[1][0], i)))
        e.Count = excel.ReadInt32(readCellValue(row, names, fmt.Sprintf("%s%d", base[2][0], i)))

		ret = append(ret, e)
	}

	return ret
}

type Output_ struct{
    Type int32 
    ID int32 
    Count int32 

}

func readOutput(row, names []string)[]*Output_{
	ret := make([]*Output_, 0)
	base := excel.Split("OutputType_,OutputID_,OutputCount_", ",")

	for i := 1; ; i++ {
		if idx, find := hasCellName(names, fmt.Sprintf("%s%d", base[0][0], i));!find {
			break
		}else if _,ok := hasCellValue(row,idx);!ok{
			break
		}
		
		e := &Output_{}
        e.Type = excel.ReadEnum(readCellValue(row, names, fmt.Sprintf("%s%d", base[0][0], i)))
        e.ID = excel.ReadInt32(readCellValue(row, names, fmt.Sprintf("%s%d", base[1][0], i)))
        e.Count = excel.ReadInt32(readCellValue(row, names, fmt.Sprintf("%s%d", base[2][0], i)))

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

func GetID(key int32)*InputOutput{
	v := Table_.indexID.Load().(map[int32]*InputOutput)[key]
	if v == nil{
		zaplogger.GetSugar().Errorf("Table %s GetID %d is nil", Table_.xlsxName, key)
	}
	return v
}

func GetIDMap() map[int32]*InputOutput {
	return Table_.indexID.Load().(map[int32]*InputOutput)
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
	Table_.xlsxName = "InputOutput.xlsx"
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

	tmpIDMap := map[int32]*InputOutput{}

	for _,row := range rows{
		if row[0] != "" {
			e := &InputOutput{}
            e.ID = excel.ReadInt32(readCellValue(row, names, "ID"))
            e.Describe = excel.ReadStr(readCellValue(row, names, "Describe"))
            e.InputType_1 = excel.ReadStr(readCellValue(row, names, "InputType_1"))
            e.InputID_1 = excel.ReadInt32(readCellValue(row, names, "InputID_1"))
            e.InputCount_1 = excel.ReadInt32(readCellValue(row, names, "InputCount_1"))
            e.InputType_2 = excel.ReadStr(readCellValue(row, names, "InputType_2"))
            e.InputID_2 = excel.ReadInt32(readCellValue(row, names, "InputID_2"))
            e.InputCount_2 = excel.ReadInt32(readCellValue(row, names, "InputCount_2"))
            e.InputType_3 = excel.ReadStr(readCellValue(row, names, "InputType_3"))
            e.InputID_3 = excel.ReadInt32(readCellValue(row, names, "InputID_3"))
            e.InputCount_3 = excel.ReadInt32(readCellValue(row, names, "InputCount_3"))
            e.InputType_4 = excel.ReadStr(readCellValue(row, names, "InputType_4"))
            e.InputID_4 = excel.ReadInt32(readCellValue(row, names, "InputID_4"))
            e.InputCount_4 = excel.ReadInt32(readCellValue(row, names, "InputCount_4"))
            e.OutputType_1 = excel.ReadStr(readCellValue(row, names, "OutputType_1"))
            e.OutputID_1 = excel.ReadInt32(readCellValue(row, names, "OutputID_1"))
            e.OutputCount_1 = excel.ReadInt32(readCellValue(row, names, "OutputCount_1"))
            e.OutputType_2 = excel.ReadStr(readCellValue(row, names, "OutputType_2"))
            e.OutputID_2 = excel.ReadInt32(readCellValue(row, names, "OutputID_2"))
            e.OutputCount_2 = excel.ReadInt32(readCellValue(row, names, "OutputCount_2"))
            e.OutputType_3 = excel.ReadStr(readCellValue(row, names, "OutputType_3"))
            e.OutputID_3 = excel.ReadInt32(readCellValue(row, names, "OutputID_3"))
            e.OutputCount_3 = excel.ReadInt32(readCellValue(row, names, "OutputCount_3"))
            e.OutputType_4 = excel.ReadStr(readCellValue(row, names, "OutputType_4"))
            e.OutputID_4 = excel.ReadInt32(readCellValue(row, names, "OutputID_4"))
            e.OutputCount_4 = excel.ReadInt32(readCellValue(row, names, "OutputCount_4"))
            e.Input = readInput(row, names)
            e.Output = readOutput(row, names)

			tmpIDMap[excel.ReadInt32(readCellValue(row, names, "ID"))] = e

		}
	}
	this.indexID.Store(tmpIDMap)	

}
