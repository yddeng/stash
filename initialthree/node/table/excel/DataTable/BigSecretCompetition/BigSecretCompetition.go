package BigSecretCompetition

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	path2 "path"
	"strings"
	"sync/atomic"
	"initialthree/node/table/excel"
	"initialthree/zaplogger"
)

type BigSecretCompetition struct{
    ID int32 
    StartTime string 
    EndTime string 
    Quest string 
    Name string 
    Story string 
    Rule string 
    Affix string 
    QuestArray []*Quest_ 

}

type Quest_ struct{
    QuestID int32 

}

func readQuestArray(row, names []string)[]*Quest_{
	value := readCellValue(row, names, "Quest")
	if value == "" {
		return nil
	}

	r := excel.Split(value,",")
	if len(r) == 0 {
		return nil
	}

	ret := make([]*Quest_,0)
	for _, v := range r {
		if len(v) == 1{
			e := &Quest_{}
	        e.QuestID = excel.ReadInt32(v[0])

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

func GetID(key int32)*BigSecretCompetition{
	v := Table_.indexID.Load().(map[int32]*BigSecretCompetition)[key]
	if v == nil{
		zaplogger.GetSugar().Errorf("Table %s GetID %d is nil", Table_.xlsxName, key)
	}
	return v
}

func GetIDMap() map[int32]*BigSecretCompetition {
	return Table_.indexID.Load().(map[int32]*BigSecretCompetition)
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
	Table_.xlsxName = "BigSecretCompetition.xlsx"
	excel.AddTable(&Table_)
}

func (this *Table)Load(path string){
	xlFile, err := excelize.OpenFile(path2.Join(path2.Join(path, "/DataTable/关卡/大秘境"), this.xlsxName))
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	tableRows := xlFile.GetRows(xlFile.GetSheetName(xlFile.GetActiveSheetIndex()))
	names := tableRows[1]
	rows := tableRows[4:]

	tmpIDMap := map[int32]*BigSecretCompetition{}

	for _,row := range rows{
		if row[0] != "" {
			e := &BigSecretCompetition{}
            e.ID = excel.ReadInt32(readCellValue(row, names, "ID"))
            e.StartTime = excel.ReadStr(readCellValue(row, names, "StartTime"))
            e.EndTime = excel.ReadStr(readCellValue(row, names, "EndTime"))
            e.Quest = excel.ReadStr(readCellValue(row, names, "Quest"))
            e.Name = excel.ReadStr(readCellValue(row, names, "Name"))
            e.Story = excel.ReadStr(readCellValue(row, names, "Story"))
            e.Rule = excel.ReadStr(readCellValue(row, names, "Rule"))
            e.Affix = excel.ReadStr(readCellValue(row, names, "Affix"))
            e.QuestArray = readQuestArray(row, names)

			tmpIDMap[excel.ReadInt32(readCellValue(row, names, "ID"))] = e

		}
	}
	this.indexID.Store(tmpIDMap)	

}
