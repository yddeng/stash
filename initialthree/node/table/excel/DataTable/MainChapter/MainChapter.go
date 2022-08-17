package MainChapter

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	path2 "path"
	"strings"
	"sync/atomic"
	"initialthree/node/table/excel"
	"initialthree/zaplogger"
)

type MainChapter struct{
    ID int32 
    ChapterType string 
    ChapterNamePrefix string 
    ChapterName string 
    PlayerLevelLimit int32 
    UnlockByMainDungeons string 
    BgSprite string 
    MapIcon string 
    Dungeons string 
    Award string 
    UnlockDescribe string 
    OldVersionConf string 
    ChapterOrder int32 
    ChapterOrderCn string 
    UINodeConfig int32 
    ChapterIcon string 
    ChapterTypeEnum int32 
    DungeonsArray []*Dungeons_ 
    UnlockByMainDungeonsArray []*UnlockByMainDungeons_ 

}

type Dungeons_ struct{
    ID int32 

}

func readDungeonsArray(row, names []string)[]*Dungeons_{
	value := readCellValue(row, names, "Dungeons")
	if value == "" {
		return nil
	}

	r := excel.Split(value,",")
	if len(r) == 0 {
		return nil
	}

	ret := make([]*Dungeons_,0)
	for _, v := range r {
		if len(v) == 1{
			e := &Dungeons_{}
	        e.ID = excel.ReadInt32(v[0])

			ret = append(ret, e)
		}
	}

	return ret
}

type UnlockByMainDungeons_ struct{
    DungeonID int32 

}

func readUnlockByMainDungeonsArray(row, names []string)[]*UnlockByMainDungeons_{
	value := readCellValue(row, names, "UnlockByMainDungeons")
	if value == "" {
		return nil
	}

	r := excel.Split(value,",")
	if len(r) == 0 {
		return nil
	}

	ret := make([]*UnlockByMainDungeons_,0)
	for _, v := range r {
		if len(v) == 1{
			e := &UnlockByMainDungeons_{}
	        e.DungeonID = excel.ReadInt32(v[0])

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

func GetID(key int32)*MainChapter{
	v := Table_.indexID.Load().(map[int32]*MainChapter)[key]
	if v == nil{
		zaplogger.GetSugar().Errorf("Table %s GetID %d is nil", Table_.xlsxName, key)
	}
	return v
}

func GetIDMap() map[int32]*MainChapter {
	return Table_.indexID.Load().(map[int32]*MainChapter)
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
	Table_.xlsxName = "MainChapter.xlsx"
	excel.AddTable(&Table_)
}

func (this *Table)Load(path string){
	xlFile, err := excelize.OpenFile(path2.Join(path2.Join(path, "/DataTable/关卡/主线"), this.xlsxName))
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	tableRows := xlFile.GetRows(xlFile.GetSheetName(xlFile.GetActiveSheetIndex()))
	names := tableRows[1]
	rows := tableRows[4:]

	tmpIDMap := map[int32]*MainChapter{}

	for _,row := range rows{
		if row[0] != "" {
			e := &MainChapter{}
            e.ID = excel.ReadInt32(readCellValue(row, names, "ID"))
            e.ChapterType = excel.ReadStr(readCellValue(row, names, "ChapterType"))
            e.ChapterNamePrefix = excel.ReadStr(readCellValue(row, names, "ChapterNamePrefix"))
            e.ChapterName = excel.ReadStr(readCellValue(row, names, "ChapterName"))
            e.PlayerLevelLimit = excel.ReadInt32(readCellValue(row, names, "PlayerLevelLimit"))
            e.UnlockByMainDungeons = excel.ReadStr(readCellValue(row, names, "UnlockByMainDungeons"))
            e.BgSprite = excel.ReadStr(readCellValue(row, names, "BgSprite"))
            e.MapIcon = excel.ReadStr(readCellValue(row, names, "MapIcon"))
            e.Dungeons = excel.ReadStr(readCellValue(row, names, "Dungeons"))
            e.Award = excel.ReadStr(readCellValue(row, names, "Award"))
            e.UnlockDescribe = excel.ReadStr(readCellValue(row, names, "UnlockDescribe"))
            e.OldVersionConf = excel.ReadStr(readCellValue(row, names, "OldVersionConf"))
            e.ChapterOrder = excel.ReadInt32(readCellValue(row, names, "chapterOrder"))
            e.ChapterOrderCn = excel.ReadStr(readCellValue(row, names, "ChapterOrderCn"))
            e.UINodeConfig = excel.ReadInt32(readCellValue(row, names, "UINodeConfig"))
            e.ChapterIcon = excel.ReadStr(readCellValue(row, names, "ChapterIcon"))
            e.ChapterTypeEnum = excel.ReadEnum(readCellValue(row, names, "ChapterType"))
            e.DungeonsArray = readDungeonsArray(row, names)
            e.UnlockByMainDungeonsArray = readUnlockByMainDungeonsArray(row, names)

			tmpIDMap[excel.ReadInt32(readCellValue(row, names, "ID"))] = e

		}
	}
	this.indexID.Store(tmpIDMap)	

}
