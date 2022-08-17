package main

import (
	"fmt"
	"initialthree/node/table/excel"
	"os"
	"path"

	"github.com/360EntSecGroup-Skylar/excelize"
)

func main() {
	if len(os.Args) < 2 {
		panic("need argument: excel_path")
	}

	excel_path := os.Args[1]
	if excel_path == "" {
		panic("need argument: excel_path")
	}

	xlFile, err := excelize.OpenFile(path.Join(excel_path, "EnumTable/BattleAttribute.xlsx"))
	if err != nil {
		panic(err)
	}

	rows := xlFile.GetRows(xlFile.GetSheetName(xlFile.GetActiveSheetIndex()))
	rows = rows[1:]

	str :=
		`
package battleAttr

type BattleAttr struct {
	ID  int32
	Val float64
}

type Info struct{
	Min         float64
	Max         float64
}

var idxToName map[int32]string
var nameToId  map[string]int32
var battleAttrInfo map[int32]*Info

const(
`

	initStr :=

		`
func init() {	
	idxToName = map[int32]string{}
	nameToId  = map[string]int32{}
	battleAttrInfo = map[int32]*Info{}
`

	c := 0
	initStr1 := ""
	initStr2 := ""
	initStr3 := ""
	info :=
		`&Info{
		Min: %f,
		Max: %f,
	}`
	for _, row := range rows {
		id := excel.ReadInt32(row[0])
		if id >= 0 {
			name := row[2]
			desc := row[1]
			min := excel.ReadFloat(row[3])
			max := excel.ReadFloat(row[4])
			str = str + fmt.Sprintf("	%s = %d //%s\n", name, id, desc)
			initStr1 = initStr1 + fmt.Sprintf(`	idxToName[%d] = "%s"`, id, name) + "\n"
			initStr2 = initStr2 + fmt.Sprintf(`	nameToId["%s"] = %d`, name, id) + "\n"
			initStr3 = initStr3 + fmt.Sprintf(`	battleAttrInfo[%d] = `, id) + fmt.Sprintf(info, min, max) + "\n"
			c++
		}
	}
	initStr = initStr + initStr1 + initStr2 + initStr3 + "}\n"
	str = str + fmt.Sprintf("	AttrMax = %d\n", c)
	str = str + ")\n"
	str = str + initStr

	str = str +
		`

func GetNameById(id int32) string {
	return idxToName[id]
}

func GetIdByName(name string) int32 {
	return nameToId[name]
}

func GetBattleAttrInfo(idx int32) *Info {
	return battleAttrInfo[idx]
}

/*
func TransFormToFloat64(id, value int32) float64 {
	info, ok := battleAttrInfo[id]
	if !ok {
		return 0
	}

	if info.NeedFixed {
		v := float64(float64(value) / float64(GlobalConst.Table_.IDMap[1].BattleAttrRate))
		return v
	}
	return float64(value)
}
 */
`

	out_path := "../../../../common/battleAttr/battleAttr.go"

	f, err := os.OpenFile(out_path, os.O_RDWR, os.ModePerm)
	if err != nil {
		if os.IsNotExist(err) {
			f, err = os.Create(out_path)
			if err != nil {
				fmt.Printf("create %s failed:%s", out_path, err.Error())
				return
			}
		} else {
			fmt.Printf("open %s failed:%s", out_path, err.Error())
			return
		}
	}

	err = os.Truncate(out_path, 0)

	_, err = f.WriteString(str)

	if nil != err {
		fmt.Printf("%s Write error:%s\n", out_path, err.Error())
	} else {
		fmt.Printf("%s Write ok\n", out_path)
	}

	f.Close()
}
