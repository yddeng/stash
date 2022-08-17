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

	xlFile, err := excelize.OpenFile(path.Join(excel_path, "EnumTable/CommonEnums.xlsx"))
	if err != nil {
		panic(err)
	}

	rows := xlFile.GetRows(xlFile.GetSheetName(xlFile.GetActiveSheetIndex()))
	rows = rows[1:]

	str :=
		`//所有项目的类型枚举
package enumType

import "fmt"

var nameToIdx map[string]int32

const(
`

	initStr :=
		`
func init() {	
	nameToIdx = map[string]int32{}
`
	for _, row := range rows {
		name := row[0]
		type_ := row[1]
		name_type := fmt.Sprintf("%s_%s", name, type_)
		enum, _ := strconv.Atoi(row[2])
		desc := row[3]
		str = str + fmt.Sprintf(`
	%s = %d //%s`, name_type, enum, desc)
		initStr = initStr + fmt.Sprintf(`	
	nameToIdx["%s"] = %d`, name_type, enum)
	}
	str = str + `
)
`
	initStr = initStr + `
}

func GetEnumType(name string) (int32,error) {
	if idx, ok := nameToIdx[name]; ok {
		return idx,nil
	}else {
		return 0,fmt.Errorf("%s is not find",name)
	}
}
`
	str = str + initStr

	out_path := "../../../../common/enumType/enumType.go"

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
