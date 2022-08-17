package WeaponAttribute

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	path2 "path"
	"strings"
	"sync/atomic"
	"initialthree/node/table/excel"
	"initialthree/zaplogger"
)

type WeaponAttribute struct{
    ID int32 
    AttributeType string 
    LevelAttrib_1 float64 
    LevelAttrib_2 float64 
    LevelAttrib_3 float64 
    LevelAttrib_4 float64 
    LevelAttrib_5 float64 
    LevelAttrib_6 float64 
    LevelAttrib_7 float64 
    LevelAttrib_8 float64 
    LevelAttrib_9 float64 
    LevelAttrib_10 float64 
    LevelAttrib_11 float64 
    LevelAttrib_12 float64 
    LevelAttrib_13 float64 
    LevelAttrib_14 float64 
    LevelAttrib_15 float64 
    LevelAttrib_16 float64 
    LevelAttrib_17 float64 
    LevelAttrib_18 float64 
    LevelAttrib_19 float64 
    LevelAttrib_20 float64 
    LevelAttrib_21 float64 
    LevelAttrib_22 float64 
    LevelAttrib_23 float64 
    LevelAttrib_24 float64 
    LevelAttrib_25 float64 
    LevelAttrib_26 float64 
    LevelAttrib_27 float64 
    LevelAttrib_28 float64 
    LevelAttrib_29 float64 
    LevelAttrib_30 float64 
    LevelAttrib_31 float64 
    LevelAttrib_32 float64 
    LevelAttrib_33 float64 
    LevelAttrib_34 float64 
    LevelAttrib_35 float64 
    LevelAttrib_36 float64 
    LevelAttrib_37 float64 
    LevelAttrib_38 float64 
    LevelAttrib_39 float64 
    LevelAttrib_40 float64 
    LevelAttrib_41 float64 
    LevelAttrib_42 float64 
    LevelAttrib_43 float64 
    LevelAttrib_44 float64 
    LevelAttrib_45 float64 
    LevelAttrib_46 float64 
    LevelAttrib_47 float64 
    LevelAttrib_48 float64 
    LevelAttrib_49 float64 
    LevelAttrib_50 float64 
    LevelAttrib_51 float64 
    LevelAttrib_52 float64 
    LevelAttrib_53 float64 
    LevelAttrib_54 float64 
    LevelAttrib_55 float64 
    LevelAttrib_56 float64 
    LevelAttrib_57 float64 
    LevelAttrib_58 float64 
    LevelAttrib_59 float64 
    LevelAttrib_60 float64 
    LevelAttrib_61 float64 
    LevelAttrib_62 float64 
    LevelAttrib_63 float64 
    LevelAttrib_64 float64 
    LevelAttrib_65 float64 
    LevelAttrib_66 float64 
    LevelAttrib_67 float64 
    LevelAttrib_68 float64 
    LevelAttrib_69 float64 
    LevelAttrib_70 float64 
    LevelAttrib_71 float64 
    LevelAttrib_72 float64 
    LevelAttrib_73 float64 
    LevelAttrib_74 float64 
    LevelAttrib_75 float64 
    LevelAttrib_76 float64 
    LevelAttrib_77 float64 
    LevelAttrib_78 float64 
    LevelAttrib_79 float64 
    LevelAttrib_80 float64 
    LevelAttrib_81 float64 
    LevelAttrib_82 float64 
    LevelAttrib_83 float64 
    LevelAttrib_84 float64 
    LevelAttrib_85 float64 
    LevelAttrib_86 float64 
    LevelAttrib_87 float64 
    LevelAttrib_88 float64 
    LevelAttrib_89 float64 
    LevelAttrib_90 float64 
    LevelAttrib_91 float64 
    LevelAttrib_92 float64 
    LevelAttrib_93 float64 
    LevelAttrib_94 float64 
    LevelAttrib_95 float64 
    LevelAttrib_96 float64 
    LevelAttrib_97 float64 
    LevelAttrib_98 float64 
    LevelAttrib_99 float64 
    LevelAttrib_100 float64 
    LevelAttr []*LevelAttr_ 
    BreakLevelAttr []*BreakLevelAttr_ 

}

type LevelAttr_ struct{
    Val float64 

}

func readLevelAttr(row, names []string)[]*LevelAttr_{
	ret := make([]*LevelAttr_, 0)
	base := excel.Split("LevelAttrib_", ",")

	for i := 1; ; i++ {
		if idx, find := hasCellName(names, fmt.Sprintf("%s%d", base[0][0], i));!find {
			break
		}else if _,ok := hasCellValue(row,idx);!ok{
			break
		}
		
		e := &LevelAttr_{}
        e.Val = excel.ReadFloat(readCellValue(row, names, fmt.Sprintf("%s%d", base[0][0], i)))

		ret = append(ret, e)
	}

	return ret
}

type BreakLevelAttr_ struct{
    Val float64 

}

func readBreakLevelAttr(row, names []string)[]*BreakLevelAttr_{
	ret := make([]*BreakLevelAttr_, 0)
	base := excel.Split("BreakLevelAttrib_", ",")

	for i := 1; ; i++ {
		if idx, find := hasCellName(names, fmt.Sprintf("%s%d", base[0][0], i));!find {
			break
		}else if _,ok := hasCellValue(row,idx);!ok{
			break
		}
		
		e := &BreakLevelAttr_{}
        e.Val = excel.ReadFloat(readCellValue(row, names, fmt.Sprintf("%s%d", base[0][0], i)))

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

func GetID(key int32)*WeaponAttribute{
	v := Table_.indexID.Load().(map[int32]*WeaponAttribute)[key]
	if v == nil{
		zaplogger.GetSugar().Errorf("Table %s GetID %d is nil", Table_.xlsxName, key)
	}
	return v
}

func GetIDMap() map[int32]*WeaponAttribute {
	return Table_.indexID.Load().(map[int32]*WeaponAttribute)
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
	Table_.xlsxName = "WeaponAttribute.xlsx"
	excel.AddTable(&Table_)
}

func (this *Table)Load(path string){
	xlFile, err := excelize.OpenFile(path2.Join(path2.Join(path, "/DataTable/武器"), this.xlsxName))
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	tableRows := xlFile.GetRows(xlFile.GetSheetName(xlFile.GetActiveSheetIndex()))
	names := tableRows[1]
	rows := tableRows[4:]

	tmpIDMap := map[int32]*WeaponAttribute{}

	for _,row := range rows{
		if row[0] != "" {
			e := &WeaponAttribute{}
            e.ID = excel.ReadInt32(readCellValue(row, names, "ID"))
            e.AttributeType = excel.ReadStr(readCellValue(row, names, "AttributeType"))
            e.LevelAttrib_1 = excel.ReadFloat(readCellValue(row, names, "LevelAttrib_1"))
            e.LevelAttrib_2 = excel.ReadFloat(readCellValue(row, names, "LevelAttrib_2"))
            e.LevelAttrib_3 = excel.ReadFloat(readCellValue(row, names, "LevelAttrib_3"))
            e.LevelAttrib_4 = excel.ReadFloat(readCellValue(row, names, "LevelAttrib_4"))
            e.LevelAttrib_5 = excel.ReadFloat(readCellValue(row, names, "LevelAttrib_5"))
            e.LevelAttrib_6 = excel.ReadFloat(readCellValue(row, names, "LevelAttrib_6"))
            e.LevelAttrib_7 = excel.ReadFloat(readCellValue(row, names, "LevelAttrib_7"))
            e.LevelAttrib_8 = excel.ReadFloat(readCellValue(row, names, "LevelAttrib_8"))
            e.LevelAttrib_9 = excel.ReadFloat(readCellValue(row, names, "LevelAttrib_9"))
            e.LevelAttrib_10 = excel.ReadFloat(readCellValue(row, names, "LevelAttrib_10"))
            e.LevelAttrib_11 = excel.ReadFloat(readCellValue(row, names, "LevelAttrib_11"))
            e.LevelAttrib_12 = excel.ReadFloat(readCellValue(row, names, "LevelAttrib_12"))
            e.LevelAttrib_13 = excel.ReadFloat(readCellValue(row, names, "LevelAttrib_13"))
            e.LevelAttrib_14 = excel.ReadFloat(readCellValue(row, names, "LevelAttrib_14"))
            e.LevelAttrib_15 = excel.ReadFloat(readCellValue(row, names, "LevelAttrib_15"))
            e.LevelAttrib_16 = excel.ReadFloat(readCellValue(row, names, "LevelAttrib_16"))
            e.LevelAttrib_17 = excel.ReadFloat(readCellValue(row, names, "LevelAttrib_17"))
            e.LevelAttrib_18 = excel.ReadFloat(readCellValue(row, names, "LevelAttrib_18"))
            e.LevelAttrib_19 = excel.ReadFloat(readCellValue(row, names, "LevelAttrib_19"))
            e.LevelAttrib_20 = excel.ReadFloat(readCellValue(row, names, "LevelAttrib_20"))
            e.LevelAttrib_21 = excel.ReadFloat(readCellValue(row, names, "LevelAttrib_21"))
            e.LevelAttrib_22 = excel.ReadFloat(readCellValue(row, names, "LevelAttrib_22"))
            e.LevelAttrib_23 = excel.ReadFloat(readCellValue(row, names, "LevelAttrib_23"))
            e.LevelAttrib_24 = excel.ReadFloat(readCellValue(row, names, "LevelAttrib_24"))
            e.LevelAttrib_25 = excel.ReadFloat(readCellValue(row, names, "LevelAttrib_25"))
            e.LevelAttrib_26 = excel.ReadFloat(readCellValue(row, names, "LevelAttrib_26"))
            e.LevelAttrib_27 = excel.ReadFloat(readCellValue(row, names, "LevelAttrib_27"))
            e.LevelAttrib_28 = excel.ReadFloat(readCellValue(row, names, "LevelAttrib_28"))
            e.LevelAttrib_29 = excel.ReadFloat(readCellValue(row, names, "LevelAttrib_29"))
            e.LevelAttrib_30 = excel.ReadFloat(readCellValue(row, names, "LevelAttrib_30"))
            e.LevelAttrib_31 = excel.ReadFloat(readCellValue(row, names, "LevelAttrib_31"))
            e.LevelAttrib_32 = excel.ReadFloat(readCellValue(row, names, "LevelAttrib_32"))
            e.LevelAttrib_33 = excel.ReadFloat(readCellValue(row, names, "LevelAttrib_33"))
            e.LevelAttrib_34 = excel.ReadFloat(readCellValue(row, names, "LevelAttrib_34"))
            e.LevelAttrib_35 = excel.ReadFloat(readCellValue(row, names, "LevelAttrib_35"))
            e.LevelAttrib_36 = excel.ReadFloat(readCellValue(row, names, "LevelAttrib_36"))
            e.LevelAttrib_37 = excel.ReadFloat(readCellValue(row, names, "LevelAttrib_37"))
            e.LevelAttrib_38 = excel.ReadFloat(readCellValue(row, names, "LevelAttrib_38"))
            e.LevelAttrib_39 = excel.ReadFloat(readCellValue(row, names, "LevelAttrib_39"))
            e.LevelAttrib_40 = excel.ReadFloat(readCellValue(row, names, "LevelAttrib_40"))
            e.LevelAttrib_41 = excel.ReadFloat(readCellValue(row, names, "LevelAttrib_41"))
            e.LevelAttrib_42 = excel.ReadFloat(readCellValue(row, names, "LevelAttrib_42"))
            e.LevelAttrib_43 = excel.ReadFloat(readCellValue(row, names, "LevelAttrib_43"))
            e.LevelAttrib_44 = excel.ReadFloat(readCellValue(row, names, "LevelAttrib_44"))
            e.LevelAttrib_45 = excel.ReadFloat(readCellValue(row, names, "LevelAttrib_45"))
            e.LevelAttrib_46 = excel.ReadFloat(readCellValue(row, names, "LevelAttrib_46"))
            e.LevelAttrib_47 = excel.ReadFloat(readCellValue(row, names, "LevelAttrib_47"))
            e.LevelAttrib_48 = excel.ReadFloat(readCellValue(row, names, "LevelAttrib_48"))
            e.LevelAttrib_49 = excel.ReadFloat(readCellValue(row, names, "LevelAttrib_49"))
            e.LevelAttrib_50 = excel.ReadFloat(readCellValue(row, names, "LevelAttrib_50"))
            e.LevelAttrib_51 = excel.ReadFloat(readCellValue(row, names, "LevelAttrib_51"))
            e.LevelAttrib_52 = excel.ReadFloat(readCellValue(row, names, "LevelAttrib_52"))
            e.LevelAttrib_53 = excel.ReadFloat(readCellValue(row, names, "LevelAttrib_53"))
            e.LevelAttrib_54 = excel.ReadFloat(readCellValue(row, names, "LevelAttrib_54"))
            e.LevelAttrib_55 = excel.ReadFloat(readCellValue(row, names, "LevelAttrib_55"))
            e.LevelAttrib_56 = excel.ReadFloat(readCellValue(row, names, "LevelAttrib_56"))
            e.LevelAttrib_57 = excel.ReadFloat(readCellValue(row, names, "LevelAttrib_57"))
            e.LevelAttrib_58 = excel.ReadFloat(readCellValue(row, names, "LevelAttrib_58"))
            e.LevelAttrib_59 = excel.ReadFloat(readCellValue(row, names, "LevelAttrib_59"))
            e.LevelAttrib_60 = excel.ReadFloat(readCellValue(row, names, "LevelAttrib_60"))
            e.LevelAttrib_61 = excel.ReadFloat(readCellValue(row, names, "LevelAttrib_61"))
            e.LevelAttrib_62 = excel.ReadFloat(readCellValue(row, names, "LevelAttrib_62"))
            e.LevelAttrib_63 = excel.ReadFloat(readCellValue(row, names, "LevelAttrib_63"))
            e.LevelAttrib_64 = excel.ReadFloat(readCellValue(row, names, "LevelAttrib_64"))
            e.LevelAttrib_65 = excel.ReadFloat(readCellValue(row, names, "LevelAttrib_65"))
            e.LevelAttrib_66 = excel.ReadFloat(readCellValue(row, names, "LevelAttrib_66"))
            e.LevelAttrib_67 = excel.ReadFloat(readCellValue(row, names, "LevelAttrib_67"))
            e.LevelAttrib_68 = excel.ReadFloat(readCellValue(row, names, "LevelAttrib_68"))
            e.LevelAttrib_69 = excel.ReadFloat(readCellValue(row, names, "LevelAttrib_69"))
            e.LevelAttrib_70 = excel.ReadFloat(readCellValue(row, names, "LevelAttrib_70"))
            e.LevelAttrib_71 = excel.ReadFloat(readCellValue(row, names, "LevelAttrib_71"))
            e.LevelAttrib_72 = excel.ReadFloat(readCellValue(row, names, "LevelAttrib_72"))
            e.LevelAttrib_73 = excel.ReadFloat(readCellValue(row, names, "LevelAttrib_73"))
            e.LevelAttrib_74 = excel.ReadFloat(readCellValue(row, names, "LevelAttrib_74"))
            e.LevelAttrib_75 = excel.ReadFloat(readCellValue(row, names, "LevelAttrib_75"))
            e.LevelAttrib_76 = excel.ReadFloat(readCellValue(row, names, "LevelAttrib_76"))
            e.LevelAttrib_77 = excel.ReadFloat(readCellValue(row, names, "LevelAttrib_77"))
            e.LevelAttrib_78 = excel.ReadFloat(readCellValue(row, names, "LevelAttrib_78"))
            e.LevelAttrib_79 = excel.ReadFloat(readCellValue(row, names, "LevelAttrib_79"))
            e.LevelAttrib_80 = excel.ReadFloat(readCellValue(row, names, "LevelAttrib_80"))
            e.LevelAttrib_81 = excel.ReadFloat(readCellValue(row, names, "LevelAttrib_81"))
            e.LevelAttrib_82 = excel.ReadFloat(readCellValue(row, names, "LevelAttrib_82"))
            e.LevelAttrib_83 = excel.ReadFloat(readCellValue(row, names, "LevelAttrib_83"))
            e.LevelAttrib_84 = excel.ReadFloat(readCellValue(row, names, "LevelAttrib_84"))
            e.LevelAttrib_85 = excel.ReadFloat(readCellValue(row, names, "LevelAttrib_85"))
            e.LevelAttrib_86 = excel.ReadFloat(readCellValue(row, names, "LevelAttrib_86"))
            e.LevelAttrib_87 = excel.ReadFloat(readCellValue(row, names, "LevelAttrib_87"))
            e.LevelAttrib_88 = excel.ReadFloat(readCellValue(row, names, "LevelAttrib_88"))
            e.LevelAttrib_89 = excel.ReadFloat(readCellValue(row, names, "LevelAttrib_89"))
            e.LevelAttrib_90 = excel.ReadFloat(readCellValue(row, names, "LevelAttrib_90"))
            e.LevelAttrib_91 = excel.ReadFloat(readCellValue(row, names, "LevelAttrib_91"))
            e.LevelAttrib_92 = excel.ReadFloat(readCellValue(row, names, "LevelAttrib_92"))
            e.LevelAttrib_93 = excel.ReadFloat(readCellValue(row, names, "LevelAttrib_93"))
            e.LevelAttrib_94 = excel.ReadFloat(readCellValue(row, names, "LevelAttrib_94"))
            e.LevelAttrib_95 = excel.ReadFloat(readCellValue(row, names, "LevelAttrib_95"))
            e.LevelAttrib_96 = excel.ReadFloat(readCellValue(row, names, "LevelAttrib_96"))
            e.LevelAttrib_97 = excel.ReadFloat(readCellValue(row, names, "LevelAttrib_97"))
            e.LevelAttrib_98 = excel.ReadFloat(readCellValue(row, names, "LevelAttrib_98"))
            e.LevelAttrib_99 = excel.ReadFloat(readCellValue(row, names, "LevelAttrib_99"))
            e.LevelAttrib_100 = excel.ReadFloat(readCellValue(row, names, "LevelAttrib_100"))
            e.LevelAttr = readLevelAttr(row, names)
            e.BreakLevelAttr = readBreakLevelAttr(row, names)

			tmpIDMap[excel.ReadInt32(readCellValue(row, names, "ID"))] = e

		}
	}
	this.indexID.Store(tmpIDMap)	

}
