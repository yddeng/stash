package EquipLevelMaxExp

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	path2 "path"
	"strings"
	"sync/atomic"
	"initialthree/node/table/excel"
	"initialthree/zaplogger"
)

type EquipLevelMaxExp struct{
    ID int32 
    LevelMaxExp_1 int32 
    LevelMaxExp_2 int32 
    LevelMaxExp_3 int32 
    LevelMaxExp_4 int32 
    LevelMaxExp_5 int32 
    LevelMaxExp_6 int32 
    LevelMaxExp_7 int32 
    LevelMaxExp_8 int32 
    LevelMaxExp_9 int32 
    LevelMaxExp_10 int32 
    LevelMaxExp_11 int32 
    LevelMaxExp_12 int32 
    LevelMaxExp_13 int32 
    LevelMaxExp_14 int32 
    LevelMaxExp_15 int32 
    LevelMaxExp_16 int32 
    LevelMaxExp_17 int32 
    LevelMaxExp_18 int32 
    LevelMaxExp_19 int32 
    LevelMaxExp_20 int32 
    LevelMaxExp_21 int32 
    LevelMaxExp_22 int32 
    LevelMaxExp_23 int32 
    LevelMaxExp_24 int32 
    LevelMaxExp_25 int32 
    LevelMaxExp_26 int32 
    LevelMaxExp_27 int32 
    LevelMaxExp_28 int32 
    LevelMaxExp_29 int32 
    LevelMaxExp_30 int32 
    LevelMaxExp_31 int32 
    LevelMaxExp_32 int32 
    LevelMaxExp_33 int32 
    LevelMaxExp_34 int32 
    LevelMaxExp_35 int32 
    LevelMaxExp_36 int32 
    LevelMaxExp_37 int32 
    LevelMaxExp_38 int32 
    LevelMaxExp_39 int32 
    LevelMaxExp_40 int32 
    LevelMaxExp_41 int32 
    LevelMaxExp_42 int32 
    LevelMaxExp_43 int32 
    LevelMaxExp_44 int32 
    LevelMaxExp_45 int32 
    LevelMaxExp_46 int32 
    LevelMaxExp_47 int32 
    LevelMaxExp_48 int32 
    LevelMaxExp_49 int32 
    LevelMaxExp_50 int32 
    LevelMaxExp_51 int32 
    LevelMaxExp_52 int32 
    LevelMaxExp_53 int32 
    LevelMaxExp_54 int32 
    LevelMaxExp_55 int32 
    LevelMaxExp_56 int32 
    LevelMaxExp_57 int32 
    LevelMaxExp_58 int32 
    LevelMaxExp_59 int32 
    LevelMaxExp_60 int32 
    LevelMaxExp_61 int32 
    LevelMaxExp_62 int32 
    LevelMaxExp_63 int32 
    LevelMaxExp_64 int32 
    LevelMaxExp_65 int32 
    LevelMaxExp_66 int32 
    LevelMaxExp_67 int32 
    LevelMaxExp_68 int32 
    LevelMaxExp_69 int32 
    LevelMaxExp_70 int32 
    LevelMaxExp_71 int32 
    LevelMaxExp_72 int32 
    LevelMaxExp_73 int32 
    LevelMaxExp_74 int32 
    LevelMaxExp_75 int32 
    LevelMaxExp_76 int32 
    LevelMaxExp_77 int32 
    LevelMaxExp_78 int32 
    LevelMaxExp_79 int32 
    LevelMaxExp_80 int32 
    LevelMaxExp_81 int32 
    LevelMaxExp_82 int32 
    LevelMaxExp_83 int32 
    LevelMaxExp_84 int32 
    LevelMaxExp_85 int32 
    LevelMaxExp_86 int32 
    LevelMaxExp_87 int32 
    LevelMaxExp_88 int32 
    LevelMaxExp_89 int32 
    LevelMaxExp_90 int32 
    LevelMaxExp_91 int32 
    LevelMaxExp_92 int32 
    LevelMaxExp_93 int32 
    LevelMaxExp_94 int32 
    LevelMaxExp_95 int32 
    LevelMaxExp_96 int32 
    LevelMaxExp_97 int32 
    LevelMaxExp_98 int32 
    LevelMaxExp_99 int32 
    LevelMaxExp_100 int32 
    MaxExp []*MaxExp_ 

}

type MaxExp_ struct{
    Exp int32 

}

func readMaxExp(row, names []string)[]*MaxExp_{
	ret := make([]*MaxExp_, 0)
	base := excel.Split("LevelMaxExp_", ",")

	for i := 1; ; i++ {
		if idx, find := hasCellName(names, fmt.Sprintf("%s%d", base[0][0], i));!find {
			break
		}else if _,ok := hasCellValue(row,idx);!ok{
			break
		}
		
		e := &MaxExp_{}
        e.Exp = excel.ReadInt32(readCellValue(row, names, fmt.Sprintf("%s%d", base[0][0], i)))

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

func GetID(key int32)*EquipLevelMaxExp{
	v := Table_.indexID.Load().(map[int32]*EquipLevelMaxExp)[key]
	if v == nil{
		zaplogger.GetSugar().Errorf("Table %s GetID %d is nil", Table_.xlsxName, key)
	}
	return v
}

func GetIDMap() map[int32]*EquipLevelMaxExp {
	return Table_.indexID.Load().(map[int32]*EquipLevelMaxExp)
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
	Table_.xlsxName = "EquipLevelMaxExp.xlsx"
	excel.AddTable(&Table_)
}

func (this *Table)Load(path string){
	xlFile, err := excelize.OpenFile(path2.Join(path2.Join(path, "/DataTable/装备"), this.xlsxName))
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	tableRows := xlFile.GetRows(xlFile.GetSheetName(xlFile.GetActiveSheetIndex()))
	names := tableRows[1]
	rows := tableRows[4:]

	tmpIDMap := map[int32]*EquipLevelMaxExp{}

	for _,row := range rows{
		if row[0] != "" {
			e := &EquipLevelMaxExp{}
            e.ID = excel.ReadInt32(readCellValue(row, names, "ID"))
            e.LevelMaxExp_1 = excel.ReadInt32(readCellValue(row, names, "LevelMaxExp_1"))
            e.LevelMaxExp_2 = excel.ReadInt32(readCellValue(row, names, "LevelMaxExp_2"))
            e.LevelMaxExp_3 = excel.ReadInt32(readCellValue(row, names, "LevelMaxExp_3"))
            e.LevelMaxExp_4 = excel.ReadInt32(readCellValue(row, names, "LevelMaxExp_4"))
            e.LevelMaxExp_5 = excel.ReadInt32(readCellValue(row, names, "LevelMaxExp_5"))
            e.LevelMaxExp_6 = excel.ReadInt32(readCellValue(row, names, "LevelMaxExp_6"))
            e.LevelMaxExp_7 = excel.ReadInt32(readCellValue(row, names, "LevelMaxExp_7"))
            e.LevelMaxExp_8 = excel.ReadInt32(readCellValue(row, names, "LevelMaxExp_8"))
            e.LevelMaxExp_9 = excel.ReadInt32(readCellValue(row, names, "LevelMaxExp_9"))
            e.LevelMaxExp_10 = excel.ReadInt32(readCellValue(row, names, "LevelMaxExp_10"))
            e.LevelMaxExp_11 = excel.ReadInt32(readCellValue(row, names, "LevelMaxExp_11"))
            e.LevelMaxExp_12 = excel.ReadInt32(readCellValue(row, names, "LevelMaxExp_12"))
            e.LevelMaxExp_13 = excel.ReadInt32(readCellValue(row, names, "LevelMaxExp_13"))
            e.LevelMaxExp_14 = excel.ReadInt32(readCellValue(row, names, "LevelMaxExp_14"))
            e.LevelMaxExp_15 = excel.ReadInt32(readCellValue(row, names, "LevelMaxExp_15"))
            e.LevelMaxExp_16 = excel.ReadInt32(readCellValue(row, names, "LevelMaxExp_16"))
            e.LevelMaxExp_17 = excel.ReadInt32(readCellValue(row, names, "LevelMaxExp_17"))
            e.LevelMaxExp_18 = excel.ReadInt32(readCellValue(row, names, "LevelMaxExp_18"))
            e.LevelMaxExp_19 = excel.ReadInt32(readCellValue(row, names, "LevelMaxExp_19"))
            e.LevelMaxExp_20 = excel.ReadInt32(readCellValue(row, names, "LevelMaxExp_20"))
            e.LevelMaxExp_21 = excel.ReadInt32(readCellValue(row, names, "LevelMaxExp_21"))
            e.LevelMaxExp_22 = excel.ReadInt32(readCellValue(row, names, "LevelMaxExp_22"))
            e.LevelMaxExp_23 = excel.ReadInt32(readCellValue(row, names, "LevelMaxExp_23"))
            e.LevelMaxExp_24 = excel.ReadInt32(readCellValue(row, names, "LevelMaxExp_24"))
            e.LevelMaxExp_25 = excel.ReadInt32(readCellValue(row, names, "LevelMaxExp_25"))
            e.LevelMaxExp_26 = excel.ReadInt32(readCellValue(row, names, "LevelMaxExp_26"))
            e.LevelMaxExp_27 = excel.ReadInt32(readCellValue(row, names, "LevelMaxExp_27"))
            e.LevelMaxExp_28 = excel.ReadInt32(readCellValue(row, names, "LevelMaxExp_28"))
            e.LevelMaxExp_29 = excel.ReadInt32(readCellValue(row, names, "LevelMaxExp_29"))
            e.LevelMaxExp_30 = excel.ReadInt32(readCellValue(row, names, "LevelMaxExp_30"))
            e.LevelMaxExp_31 = excel.ReadInt32(readCellValue(row, names, "LevelMaxExp_31"))
            e.LevelMaxExp_32 = excel.ReadInt32(readCellValue(row, names, "LevelMaxExp_32"))
            e.LevelMaxExp_33 = excel.ReadInt32(readCellValue(row, names, "LevelMaxExp_33"))
            e.LevelMaxExp_34 = excel.ReadInt32(readCellValue(row, names, "LevelMaxExp_34"))
            e.LevelMaxExp_35 = excel.ReadInt32(readCellValue(row, names, "LevelMaxExp_35"))
            e.LevelMaxExp_36 = excel.ReadInt32(readCellValue(row, names, "LevelMaxExp_36"))
            e.LevelMaxExp_37 = excel.ReadInt32(readCellValue(row, names, "LevelMaxExp_37"))
            e.LevelMaxExp_38 = excel.ReadInt32(readCellValue(row, names, "LevelMaxExp_38"))
            e.LevelMaxExp_39 = excel.ReadInt32(readCellValue(row, names, "LevelMaxExp_39"))
            e.LevelMaxExp_40 = excel.ReadInt32(readCellValue(row, names, "LevelMaxExp_40"))
            e.LevelMaxExp_41 = excel.ReadInt32(readCellValue(row, names, "LevelMaxExp_41"))
            e.LevelMaxExp_42 = excel.ReadInt32(readCellValue(row, names, "LevelMaxExp_42"))
            e.LevelMaxExp_43 = excel.ReadInt32(readCellValue(row, names, "LevelMaxExp_43"))
            e.LevelMaxExp_44 = excel.ReadInt32(readCellValue(row, names, "LevelMaxExp_44"))
            e.LevelMaxExp_45 = excel.ReadInt32(readCellValue(row, names, "LevelMaxExp_45"))
            e.LevelMaxExp_46 = excel.ReadInt32(readCellValue(row, names, "LevelMaxExp_46"))
            e.LevelMaxExp_47 = excel.ReadInt32(readCellValue(row, names, "LevelMaxExp_47"))
            e.LevelMaxExp_48 = excel.ReadInt32(readCellValue(row, names, "LevelMaxExp_48"))
            e.LevelMaxExp_49 = excel.ReadInt32(readCellValue(row, names, "LevelMaxExp_49"))
            e.LevelMaxExp_50 = excel.ReadInt32(readCellValue(row, names, "LevelMaxExp_50"))
            e.LevelMaxExp_51 = excel.ReadInt32(readCellValue(row, names, "LevelMaxExp_51"))
            e.LevelMaxExp_52 = excel.ReadInt32(readCellValue(row, names, "LevelMaxExp_52"))
            e.LevelMaxExp_53 = excel.ReadInt32(readCellValue(row, names, "LevelMaxExp_53"))
            e.LevelMaxExp_54 = excel.ReadInt32(readCellValue(row, names, "LevelMaxExp_54"))
            e.LevelMaxExp_55 = excel.ReadInt32(readCellValue(row, names, "LevelMaxExp_55"))
            e.LevelMaxExp_56 = excel.ReadInt32(readCellValue(row, names, "LevelMaxExp_56"))
            e.LevelMaxExp_57 = excel.ReadInt32(readCellValue(row, names, "LevelMaxExp_57"))
            e.LevelMaxExp_58 = excel.ReadInt32(readCellValue(row, names, "LevelMaxExp_58"))
            e.LevelMaxExp_59 = excel.ReadInt32(readCellValue(row, names, "LevelMaxExp_59"))
            e.LevelMaxExp_60 = excel.ReadInt32(readCellValue(row, names, "LevelMaxExp_60"))
            e.LevelMaxExp_61 = excel.ReadInt32(readCellValue(row, names, "LevelMaxExp_61"))
            e.LevelMaxExp_62 = excel.ReadInt32(readCellValue(row, names, "LevelMaxExp_62"))
            e.LevelMaxExp_63 = excel.ReadInt32(readCellValue(row, names, "LevelMaxExp_63"))
            e.LevelMaxExp_64 = excel.ReadInt32(readCellValue(row, names, "LevelMaxExp_64"))
            e.LevelMaxExp_65 = excel.ReadInt32(readCellValue(row, names, "LevelMaxExp_65"))
            e.LevelMaxExp_66 = excel.ReadInt32(readCellValue(row, names, "LevelMaxExp_66"))
            e.LevelMaxExp_67 = excel.ReadInt32(readCellValue(row, names, "LevelMaxExp_67"))
            e.LevelMaxExp_68 = excel.ReadInt32(readCellValue(row, names, "LevelMaxExp_68"))
            e.LevelMaxExp_69 = excel.ReadInt32(readCellValue(row, names, "LevelMaxExp_69"))
            e.LevelMaxExp_70 = excel.ReadInt32(readCellValue(row, names, "LevelMaxExp_70"))
            e.LevelMaxExp_71 = excel.ReadInt32(readCellValue(row, names, "LevelMaxExp_71"))
            e.LevelMaxExp_72 = excel.ReadInt32(readCellValue(row, names, "LevelMaxExp_72"))
            e.LevelMaxExp_73 = excel.ReadInt32(readCellValue(row, names, "LevelMaxExp_73"))
            e.LevelMaxExp_74 = excel.ReadInt32(readCellValue(row, names, "LevelMaxExp_74"))
            e.LevelMaxExp_75 = excel.ReadInt32(readCellValue(row, names, "LevelMaxExp_75"))
            e.LevelMaxExp_76 = excel.ReadInt32(readCellValue(row, names, "LevelMaxExp_76"))
            e.LevelMaxExp_77 = excel.ReadInt32(readCellValue(row, names, "LevelMaxExp_77"))
            e.LevelMaxExp_78 = excel.ReadInt32(readCellValue(row, names, "LevelMaxExp_78"))
            e.LevelMaxExp_79 = excel.ReadInt32(readCellValue(row, names, "LevelMaxExp_79"))
            e.LevelMaxExp_80 = excel.ReadInt32(readCellValue(row, names, "LevelMaxExp_80"))
            e.LevelMaxExp_81 = excel.ReadInt32(readCellValue(row, names, "LevelMaxExp_81"))
            e.LevelMaxExp_82 = excel.ReadInt32(readCellValue(row, names, "LevelMaxExp_82"))
            e.LevelMaxExp_83 = excel.ReadInt32(readCellValue(row, names, "LevelMaxExp_83"))
            e.LevelMaxExp_84 = excel.ReadInt32(readCellValue(row, names, "LevelMaxExp_84"))
            e.LevelMaxExp_85 = excel.ReadInt32(readCellValue(row, names, "LevelMaxExp_85"))
            e.LevelMaxExp_86 = excel.ReadInt32(readCellValue(row, names, "LevelMaxExp_86"))
            e.LevelMaxExp_87 = excel.ReadInt32(readCellValue(row, names, "LevelMaxExp_87"))
            e.LevelMaxExp_88 = excel.ReadInt32(readCellValue(row, names, "LevelMaxExp_88"))
            e.LevelMaxExp_89 = excel.ReadInt32(readCellValue(row, names, "LevelMaxExp_89"))
            e.LevelMaxExp_90 = excel.ReadInt32(readCellValue(row, names, "LevelMaxExp_90"))
            e.LevelMaxExp_91 = excel.ReadInt32(readCellValue(row, names, "LevelMaxExp_91"))
            e.LevelMaxExp_92 = excel.ReadInt32(readCellValue(row, names, "LevelMaxExp_92"))
            e.LevelMaxExp_93 = excel.ReadInt32(readCellValue(row, names, "LevelMaxExp_93"))
            e.LevelMaxExp_94 = excel.ReadInt32(readCellValue(row, names, "LevelMaxExp_94"))
            e.LevelMaxExp_95 = excel.ReadInt32(readCellValue(row, names, "LevelMaxExp_95"))
            e.LevelMaxExp_96 = excel.ReadInt32(readCellValue(row, names, "LevelMaxExp_96"))
            e.LevelMaxExp_97 = excel.ReadInt32(readCellValue(row, names, "LevelMaxExp_97"))
            e.LevelMaxExp_98 = excel.ReadInt32(readCellValue(row, names, "LevelMaxExp_98"))
            e.LevelMaxExp_99 = excel.ReadInt32(readCellValue(row, names, "LevelMaxExp_99"))
            e.LevelMaxExp_100 = excel.ReadInt32(readCellValue(row, names, "LevelMaxExp_100"))
            e.MaxExp = readMaxExp(row, names)

			tmpIDMap[excel.ReadInt32(readCellValue(row, names, "ID"))] = e

		}
	}
	this.indexID.Store(tmpIDMap)	

}
