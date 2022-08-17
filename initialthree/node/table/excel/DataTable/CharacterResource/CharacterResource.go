package CharacterResource

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	path2 "path"
	"strings"
	"sync/atomic"
	"initialthree/node/table/excel"
	"initialthree/zaplogger"
)

type CharacterResource struct{
    ID int32 
    Comment string 
    DefaultName string 
    Prefab string 
    BehaviorProfile string 
    PrefabHD string 
    YuruProfile string 
    DrawCardPicture string 
    DrawCardSingleResultPicture string 
    DrawCardResultPicture string 
    Portrait string 
    UltimateHeadIcon string 
    HeadIcon string 
    UltimateBottomIcon string 
    UltimateIcon string 
    UltimateBottomVfx string 
    UltimateTopVfx string 
    ChainSkillIsBind string 
    EmotionTexture string 
    DefaultEmotion string 
    BattleEndShowTimeline string 
    BattleEndShowEmotion string 
    DamageElementType string 
    WeaponScale string 
    DissolveDelayTime float64 
    DissolveStartValue float64 
    DissolveEndValue float64 
    DissolveDuration float64 
    ShowEquipVfx string 
    ShowEquipVfxRoot string 
    DamageElementTypeEnum int32 

}

var Table_ Table

type Table struct{
	xlsxName string
	nextPath string
	indexID  atomic.Value

}

func GetID(key int32)*CharacterResource{
	v := Table_.indexID.Load().(map[int32]*CharacterResource)[key]
	if v == nil{
		zaplogger.GetSugar().Errorf("Table %s GetID %d is nil", Table_.xlsxName, key)
	}
	return v
}

func GetIDMap() map[int32]*CharacterResource {
	return Table_.indexID.Load().(map[int32]*CharacterResource)
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
	Table_.xlsxName = "CharacterResource.xlsx"
	excel.AddTable(&Table_)
}

func (this *Table)Load(path string){
	xlFile, err := excelize.OpenFile(path2.Join(path2.Join(path, "/DataTable/角色"), this.xlsxName))
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	tableRows := xlFile.GetRows(xlFile.GetSheetName(xlFile.GetActiveSheetIndex()))
	names := tableRows[1]
	rows := tableRows[4:]

	tmpIDMap := map[int32]*CharacterResource{}

	for _,row := range rows{
		if row[0] != "" {
			e := &CharacterResource{}
            e.ID = excel.ReadInt32(readCellValue(row, names, "ID"))
            e.Comment = excel.ReadStr(readCellValue(row, names, "Comment"))
            e.DefaultName = excel.ReadStr(readCellValue(row, names, "DefaultName"))
            e.Prefab = excel.ReadStr(readCellValue(row, names, "Prefab"))
            e.BehaviorProfile = excel.ReadStr(readCellValue(row, names, "BehaviorProfile"))
            e.PrefabHD = excel.ReadStr(readCellValue(row, names, "PrefabHD"))
            e.YuruProfile = excel.ReadStr(readCellValue(row, names, "YuruProfile"))
            e.DrawCardPicture = excel.ReadStr(readCellValue(row, names, "DrawCardPicture"))
            e.DrawCardSingleResultPicture = excel.ReadStr(readCellValue(row, names, "DrawCardSingleResultPicture"))
            e.DrawCardResultPicture = excel.ReadStr(readCellValue(row, names, "DrawCardResultPicture"))
            e.Portrait = excel.ReadStr(readCellValue(row, names, "Portrait"))
            e.UltimateHeadIcon = excel.ReadStr(readCellValue(row, names, "UltimateHeadIcon"))
            e.HeadIcon = excel.ReadStr(readCellValue(row, names, "HeadIcon"))
            e.UltimateBottomIcon = excel.ReadStr(readCellValue(row, names, "UltimateBottomIcon"))
            e.UltimateIcon = excel.ReadStr(readCellValue(row, names, "UltimateIcon"))
            e.UltimateBottomVfx = excel.ReadStr(readCellValue(row, names, "UltimateBottomVfx"))
            e.UltimateTopVfx = excel.ReadStr(readCellValue(row, names, "UltimateTopVfx"))
            e.ChainSkillIsBind = excel.ReadStr(readCellValue(row, names, "ChainSkillIsBind"))
            e.EmotionTexture = excel.ReadStr(readCellValue(row, names, "EmotionTexture"))
            e.DefaultEmotion = excel.ReadStr(readCellValue(row, names, "DefaultEmotion"))
            e.BattleEndShowTimeline = excel.ReadStr(readCellValue(row, names, "BattleEndShowTimeline"))
            e.BattleEndShowEmotion = excel.ReadStr(readCellValue(row, names, "BattleEndShowEmotion"))
            e.DamageElementType = excel.ReadStr(readCellValue(row, names, "DamageElementType"))
            e.WeaponScale = excel.ReadStr(readCellValue(row, names, "WeaponScale"))
            e.DissolveDelayTime = excel.ReadFloat(readCellValue(row, names, "DissolveDelayTime"))
            e.DissolveStartValue = excel.ReadFloat(readCellValue(row, names, "DissolveStartValue"))
            e.DissolveEndValue = excel.ReadFloat(readCellValue(row, names, "DissolveEndValue"))
            e.DissolveDuration = excel.ReadFloat(readCellValue(row, names, "DissolveDuration"))
            e.ShowEquipVfx = excel.ReadStr(readCellValue(row, names, "ShowEquipVfx"))
            e.ShowEquipVfxRoot = excel.ReadStr(readCellValue(row, names, "ShowEquipVfxRoot"))
            e.DamageElementTypeEnum = excel.ReadEnum(readCellValue(row, names, "DamageElementType"), "")

			tmpIDMap[excel.ReadInt32(readCellValue(row, names, "ID"))] = e

		}
	}
	this.indexID.Store(tmpIDMap)	

}
