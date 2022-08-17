package main

import (
	"fmt"
	"os"
	"path"
	"strconv"

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

	xlFile, err := excelize.OpenFile(path.Join(excel_path, "EnumTable/UsualAttribute.xlsx"))
	if err != nil {
		panic(err)
	}

	rows := xlFile.GetRows(xlFile.GetSheetName(xlFile.GetActiveSheetIndex()))
	rows = rows[1:]

	str :=
		`
package attr

type Attr struct {
	ID  int32
	Val int64
}

type AttrInfo struct {
	Min int64
	Max int64
}

var nameToIdx map[string]int32
var idxToName map[int32]string
var infoMap   map[int32]*AttrInfo

const(
`

	initStr :=
		`
func init() {	
	nameToIdx = map[string]int32{}
	idxToName = map[int32]string{}
	infoMap   = map[int32]*AttrInfo{}

`
	initInfoStr := ""
	info :=
		`&AttrInfo{
		Min: %d,
		Max: %d,
	}`

	c := 0
	for _, row := range rows {
		id, _ := strconv.Atoi(row[0])
		if id >= 0 {
			name := row[2]
			desc := row[1]
			min, _ := strconv.Atoi(row[3])
			max, _ := strconv.Atoi(row[4])
			str = str + fmt.Sprintf("	%s = %d //%s\n", name, id, desc)
			initStr = initStr + fmt.Sprintf(`	nameToIdx["%s"] = %d`, name, id) + "\n"
			initStr = initStr + fmt.Sprintf(`	idxToName[%d] = "%s"`, id, name) + "\n"
			initInfoStr = initInfoStr + fmt.Sprintf(`	infoMap[%d] = `, id) + fmt.Sprintf(info, min, max) + "\n"
			c++
		}
	}
	initStr = initStr + "\n" + initInfoStr + "}\n"
	str = str + fmt.Sprintf("	AttrMax = %d\n", c)
	str = str + ")\n"
	str = str + initStr

	str = str +
		`
func GetIdByName(name string) int32 {
	return nameToIdx[name]
}

func GetNameById(id int32) string {
	return idxToName[id]
}

func GetAttrInfo(idx int32) *AttrInfo {
	return infoMap[idx]
}
`

	out_path := "../../../../common/attr/attr.go"

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
