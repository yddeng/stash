package main

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"initialthree/protocol/cs/message"
	"io/ioutil"
	"strings"
)

var sheet = "Sheet1"

func genErrExcel(inPath, outPath string) {
	data, err := ioutil.ReadFile(inPath)
	if err != nil {
		fmt.Println(err)
		return
	}

	xFile, err := excelize.OpenFile(outPath)
	if err != nil {
		fmt.Println(err)
		return
	}

	names := map[string]struct{}{}

	dataStr := string(data)
	lines := strings.Split(dataStr, "\n")

	for _, line := range lines {
		//fmt.Println(line)
		s1 := strings.Split(line, "=")
		if len(s1) != 2 {
			//fmt.Printf("%s failed\n", line)
			continue
		}
		name := strings.TrimSpace(s1[0])
		s2 := strings.Split(s1[1], "//")
		if len(s2) != 2 {
			//fmt.Printf("%s failed\n", line)
			continue
		}
		id, ok := message.ErrCode_value[name]
		if !ok {
			//fmt.Printf("%s not find\n", name)
			continue
		}
		desc := strings.TrimSpace(s2[1])

		names[name] = struct{}{}

		set(xFile, name, int(id), desc)
	}

	clear(xFile, names)

	err = xFile.Save()
	if err != nil {
		fmt.Println(err)
	}
}

func set(xFile *excelize.File, name string, id int, desc string) {

	ok := false
	i := 5
	//row := 0
	for {
		aname := xFile.GetCellValue(sheet, fmt.Sprintf("B%d", i))
		if aname == "" {
			break
		} else {
			//bId, err := strconv.Atoi(xFile.GetCellValue(sheet, fmt.Sprintf("B%d", i)))
			//if err == nil && bId < id {
			//	row = i + 1
			//}
		}

		if aname == name {
			xFile.SetCellInt(sheet, fmt.Sprintf("A%d", i), id)
			xFile.SetCellStr(sheet, fmt.Sprintf("C%d", i), desc)
			ok = true
			break
		}

		i++
	}

	if !ok {
		fmt.Println("insert", name, id, desc)
		//xFile.InsertRow(sheet, row-1)
		xFile.SetCellInt(sheet, fmt.Sprintf("A%d", i), id)
		xFile.SetCellStr(sheet, fmt.Sprintf("B%d", i), name)
		xFile.SetCellStr(sheet, fmt.Sprintf("C%d", i), desc)
	}

}

// 清理不再使用的 错误码
func clear(xFile *excelize.File, names map[string]struct{}) {
	i := 5
	for {
		aname := xFile.GetCellValue(sheet, fmt.Sprintf("B%d", i))
		if aname == "" {
			break
		}

		if _, ok := names[aname]; !ok {
			fmt.Println("remove", aname)
			xFile.RemoveRow(sheet, i-1)
		}

		i++
	}
}

func main() {
	fmt.Println("--------------------- gen_er_excel start")
	genErrExcel("../proto/message/errerror.proto", "/Users/yidongdeng/svn/TheInitial3_Dev/TheInitial3_Dev_Config/Excel/DataTable/ErrCode.xlsx")
	fmt.Println("--------------------- gen_er_excel end")
}
