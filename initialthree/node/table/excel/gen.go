package excel

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"os"
	"path"
	"strings"
)

/*
	1, 自动生成表内字段
	2，文件的读取，生成路径
	3，index索引生成方式
	4，特殊类型生成
    5. 默认表第一列为ID
	6. 定义字段及类型，以自定义为准
*/
/*
//const (
	//	DefType_key       // celName@key@string //类型int32
	//	DefType_mutilKey1 // name@mutilKey1@celName,celName //类型string
	//	DefType_mutilKey2 // name@mutilKey2@celName:1000,celName:1 //类型int32
	//
	//  自定义字段及类型  // celName@int64 //是表中有的字段名，类型以这里为准
	//	DefType_enum   // celName@enum
	//	DefType_struct // celName@struct@name:type,name:type@;,
	//	DefType_array1 // celName@array1@name:type,name:type@:, //对当前格子的数据分
	//	DefType_array2 // name@array2@name:type,name:type@celName,celName@1,3 //将多个格子合并
	//)
*/

type EleType struct {
	name string
	tt   string
	step []string
}

type TabDef struct {
	Name      string
	ReadPath  string
	WritePath string
	Indexs    []EleType         //索引
	Process   []EleType         //特殊结构
	Define    map[string]string //自定义字段
}

func Split(value, sep string) [][]string {
	if value == "" {
		return nil
	}

	ret := make([][]string, 0)
	if len(sep) > 1 {
		sep1 := string([]byte(sep)[:1])
		sep2 := string([]byte(sep)[1:])
		r := strings.Split(value, sep1)
		for _, v := range r {
			if v != "" {
				ret = append(ret, strings.Split(v, sep2))
			}
		}
	} else {
		s1 := strings.Split(value, sep)
		for _, s2 := range s1 {
			ret = append(ret, []string{s2})
		}
	}

	return ret
}

func makeReadFunc(tt string) string {
	if tt == "int32" || tt == "int" {
		return "excel.ReadInt32"
	} else if tt == "int64" {
		return "excel.ReadInt64"
	} else if tt == "float" || tt == "float64" || tt == "double" {
		return "excel.ReadFloat"
	} else if tt == "bool" {
		return "excel.ReadBool"
	} else if tt == "enum" {
		return "excel.ReadEnum"
	} else {
		return "excel.ReadStr"
	}
}

func makeStruct(name, elems string) string {
	return fmt.Sprintf(`
type %s_ struct{
%s
}
`, name, elems)
}

func makeTT(tt string) string {
	if tt == "int" || tt == "int32" {
		return "int32"
	} else if tt == "int64" {
		return "int64"
	} else if tt == "float" || tt == "float64" || tt == "double" {
		return "float64"
	} else if tt == "bool" {
		return "bool"
	} else if tt == "enum" {
		return "int32"
	} else {
		return "string"
	}
}

func makeTableTT(rows [][]string, name string) string {
	tt := rows[0]
	names := rows[1]
	for i, n := range names {
		if n == name {
			return makeTT(tt[i])
		}
	}
	return "string"
}

func processCol(tab *TabDef, v string) {

	t1 := strings.Split(v, "@")
	if len(t1) < 2 {
		fmt.Printf(" %s 非法列定义 %s \n", tab.Name, v)
		os.Exit(0)
	}

	if t1[1] == "key" || t1[1] == "mutilKey1" || t1[1] == "mutilKey2" {
		tab.Indexs = append(tab.Indexs, EleType{
			name: t1[0],
			tt:   t1[1],
			step: t1,
		})
	} else if t1[1] == "enum" || t1[1] == "struct" || t1[1] == "array1" || t1[1] == "array2" {
		tab.Process = append(tab.Process, EleType{
			name: t1[0],
			tt:   t1[1],
			step: t1,
		})
	} else if checkTT(t1[1]) {
		tab.Define[t1[0]] = t1[1]
	} else {
		fmt.Printf(" %s 非法列定义 %s \n", tab.Name, v)
		os.Exit(0)
	}
}

func checkTT(str string) bool {
	if str == "int32" || str == "int64" ||
		str == "float64" || str == "bool" || str == "string" {
		return true
	}
	return false
}

func genHead(tab *TabDef, f *os.File) {
	template := `package %s

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	path2 "path"
	"strings"
	"sync/atomic"
	"initialthree/node/table/excel"
	"initialthree/zaplogger"
)
`
	out := fmt.Sprintf(template, tab.Name)
	f.WriteString(out)

}

func genEntity(rows [][]string, tab *TabDef, f *os.File) {
	template := `
type %s struct{
%s
}
`
	elem := ""
	for i := 0; i < len(rows[0]); i++ {
		tt := rows[0][i]
		name := rows[1][i]

		if name != "" {
			if i == 0 {
				elem += fmt.Sprintf("    %s int32 \n", strings.Title(name))
			} else {
				//处理自定义的字段
				if v, ok := tab.Define[name]; ok {
					elem += fmt.Sprintf("    %s %s \n", strings.Title(name), makeTT(v))
				} else {
					elem += fmt.Sprintf("    %s %s \n", strings.Title(name), makeTT(tt))
				}
			}
		}
	}

	for _, e := range tab.Process {
		switch e.tt {
		case "enum":
			elem += fmt.Sprintf("    %sEnum int32 \n", e.name)
		case "struct":
			elem += fmt.Sprintf("    %sStruct *%s_ \n", e.name, e.name)
		case "array1":
			elem += fmt.Sprintf("    %sArray []*%s_ \n", e.name, e.name)
		case "array2":
			elem += fmt.Sprintf("    %s []*%s_ \n", e.name, e.name)
		}
	}

	out := fmt.Sprintf(template, tab.Name, elem)
	f.WriteString(out)
}

func genProcess(tab *TabDef, f *os.File) {

	out := ""
	for _, e := range tab.Process {
		if e.tt == "struct" {
			elem1, elem2 := "", ""
			defs := strings.Split(e.step[2], ",")
			for k, def := range defs {
				s := strings.Split(def, ":")
				elem1 += fmt.Sprintf("    %s %s \n", strings.Title(s[0]), makeTT(s[1]))
				elem2 += fmt.Sprintf("    e.%s = %s(r[%d][0])\n", strings.Title(s[0]), makeReadFunc(s[1]), k)
			}
			out += makeStruct(e.name, elem1)
			out += fmt.Sprintf(`
func read%sStruct(row, names []string)*%s_{
	value := readCellValue(row, names, "%s")
	r := excel.Split(value,"%s")
	
	e := &%s_{}
%s
	return e
}
`, e.name, e.name, e.name, e.step[3], e.name, elem2)
		} else if e.tt == "array1" {
			elem1, elem2 := "", ""
			defs := strings.Split(e.step[2], ",")
			for k, def := range defs {
				s := strings.Split(def, ":")
				elem1 += fmt.Sprintf("    %s %s \n", strings.Title(s[0]), makeTT(s[1]))
				elem2 += fmt.Sprintf("        e.%s = %s(v[%d])\n", strings.Title(s[0]), makeReadFunc(s[1]), k)
			}
			out += makeStruct(e.name, elem1)
			out += fmt.Sprintf(`
func read%sArray(row, names []string)[]*%s_{
	value := readCellValue(row, names, "%s")
	if value == "" {
		return nil
	}

	r := excel.Split(value,"%s")
	if len(r) == 0 {
		return nil
	}

	ret := make([]*%s_,0)
	for _, v := range r {
		if len(v) == %d{
			e := &%s_{}
	%s
			ret = append(ret, e)
		}
	}

	return ret
}
`, e.name, e.name, e.name, e.step[3], e.name, len(defs), e.name, elem2)
		} else if e.tt == "array2" {
			elem1, elem2 := "", ""
			defs := strings.Split(e.step[2], ",")
			for k, def := range defs {
				s := strings.Split(def, ":")
				elem1 += fmt.Sprintf("    %s %s \n", strings.Title(s[0]), makeTT(s[1]))
				elem2 += fmt.Sprintf("        e.%s = %s(readCellValue(row, names, fmt.Sprintf(", strings.Title(s[0]), makeReadFunc(s[1])) +
					`"%s%d"` + fmt.Sprintf(", base[%d][0], i)))\n", k)
			}
			out += makeStruct(e.name, elem1)
			out += fmt.Sprintf(`
func read%s(row, names []string)[]*%s_{
	ret := make([]*%s_, 0)
	base := excel.Split("%s", ",")

	for i := 1; ; i++ {
		if idx, find := hasCellName(names, fmt.Sprintf("%s%s", base[0][0], i));!find {
			break
		}else if _,ok := hasCellValue(row,idx);!ok{
			break
		}
		
		e := &%s_{}
%s
		ret = append(ret, e)
	}

	return ret
}
`, e.name, e.name, e.name, e.step[3], "%s", "%d", e.name, elem2)
		}

	}

	f.WriteString(out)
}

func genIndex(rows [][]string, tab *TabDef, f *os.File) {
	template := `
var Table_ Table

type Table struct{
	xlsxName string
	nextPath string
	indexID  atomic.Value
%s
}

func GetID(key int32)*%s{
	v := Table_.indexID.Load().(map[int32]*%s)[key]
	if v == nil{
		zaplogger.GetSugar().Errorf("Table %s GetID %s is nil", Table_.xlsxName, key)
	}
	return v
}

func GetIDMap() map[int32]*%s {
	return Table_.indexID.Load().(map[int32]*%s)
}
%s
`

	elem1, elem2 := "", ""
	for _, v := range tab.Indexs {
		if v.tt == "key" && (v.name != "ID" || v.step[2] != "int32") {
			elem1 += fmt.Sprintf(`    index%s%s atomic.Value
`, v.name, v.step[2])
			elem2 += fmt.Sprintf(`
func Get%s%s(key %s)*%s{
	return Table_.index%s%s.Load().(map[%s]*%s)[key]
}
`, v.name, v.step[2], v.step[2], tab.Name, v.name, v.step[2], v.step[2], tab.Name)

		} else if v.tt == "mutilKey1" {
			elem1 += fmt.Sprintf(`    index%s atomic.Value
`, v.name)
			s := strings.Split(v.step[2], ",")
			keys := fmt.Sprintf("%s %s", s[0], makeTableTT(rows, s[0]))
			key := `    key := fmt.Sprintf("%v,",` + fmt.Sprintf(`%s)`, s[0])
			for _, v1 := range s[1:] {
				keys = fmt.Sprintf(`, %s %s`, v1, makeTableTT(rows, v1))
				key += `+fmt.Sprintf("%v,",` + fmt.Sprintf(`%s)`, v1)
			}
			elem2 += fmt.Sprintf(`
func Get%s(%s)*%s{
%s
	return Table_.index%s.Load().(map[string]*%s)[key]
}
`, v.name, keys, tab.Name, key, v.name, tab.Name)

		} else if v.tt == "mutilKey2" {
			s := strings.Split(v.step[2], ",")
			elem := strings.Split(s[0], ":")
			keys := fmt.Sprintf("%s", elem[0])
			key := fmt.Sprintf("    key := %s*%s", elem[0], elem[1])
			for _, v1 := range s[1:] {
				elem = strings.Split(v1, ":")
				keys += fmt.Sprintf(", %s", elem[0])
				key += fmt.Sprintf(" + %s*%s", elem[0], elem[1])
			}
			elem2 += fmt.Sprintf(`
func Get%s(%s int32)*%s{
%s
	return Table_.indexID.Load().(map[int32]*%s)[key]
}
`, v.name, keys, tab.Name, key, tab.Name)
		}
	}

	out := fmt.Sprintf(template, elem1, tab.Name, tab.Name, "%s", "%d", tab.Name, tab.Name, elem2)
	f.WriteString(out)
}

func genLoad(rows [][]string, tab *TabDef, f *os.File) {
	template := `
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
	Table_.xlsxName = "%s.xlsx"
	excel.AddTable(&Table_)
}

func (this *Table)Load(path string){
	xlFile, err := excelize.OpenFile(path2.Join(path2.Join(path, "%s"), this.xlsxName))
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	tableRows := xlFile.GetRows(xlFile.GetSheetName(xlFile.GetActiveSheetIndex()))
	names := tableRows[1]
	rows := tableRows[4:]

	tmpIDMap := map[int32]*%s{}
%s
	for _,row := range rows{
		if row[0] != "" {
			e := &%s{}
%s
			tmpIDMap[excel.ReadInt32(readCellValue(row, names, "%s"))] = e
%s
		}
	}
	this.indexID.Store(tmpIDMap)	
%s
}
`

	elem1 := ""
	tmp, app, cp := "", "", ""
	tt := rows[0]
	for k := 0; k < len(rows[1]); k++ {
		name := rows[1][k]
		if name != "" {
			if k == 0 {
				elem1 += fmt.Sprintf(`            e.%s = excel.ReadInt32(readCellValue(row, names, "%s"))
`, strings.Title(name), name)
			} else {
				//处理自定字段
				if v, ok := tab.Define[name]; ok {
					elem1 += fmt.Sprintf(`            e.%s = %s(readCellValue(row, names, "%s"))
`, strings.Title(name), makeReadFunc(v), name)
				} else {
					elem1 += fmt.Sprintf(`            e.%s = %s(readCellValue(row, names, "%s"))
`, strings.Title(name), makeReadFunc(tt[k]), name)
				}
			}
		}
	}
	for _, v := range tab.Indexs {
		if v.tt == "key" && (v.name != "ID" || v.step[2] != "int32") {
			tmp += fmt.Sprintf(`    tmp%s%s := map[%s]*%s{}
`, v.name, v.step[2], v.step[2], tab.Name)
			app += fmt.Sprintf(`            tmp%s%s[%s(readCellValue(row, names, "%s"))] = e
`, v.name, v.step[2], makeReadFunc(v.step[2]), v.name)
			cp += fmt.Sprintf(`    this.index%s%s.Store(tmp%s%s)
`, v.name, v.step[2], v.name, v.step[2])

		} else if v.tt == "mutilKey1" {
			tmp += fmt.Sprintf(`    tmp%s := map[string]*%s{}
`, v.name, tab.Name)
			s := strings.Split(v.step[2], ",")
			key := ""
			for _, v1 := range s {
				key += `fmt.Sprintf("%s,",` + fmt.Sprintf("e.%s)", v1)
			}
			app += fmt.Sprintf(`            key := %s
`, key)
			app += fmt.Sprintf(`            tmp%s[key] = e
`, v.name)
			cp += fmt.Sprintf(`    this.index%s.Store(tmp%s)
`, v.name, v.name)
		}

	}

	for _, v := range tab.Process {
		if v.tt == "enum" {
			if len(v.step) > 2 {
				elem1 += fmt.Sprintf(`            e.%sEnum = excel.ReadEnum(readCellValue(row, names, "%s"), "%s")
`, v.name, v.name, v.step[2])
			} else {
				elem1 += fmt.Sprintf(`            e.%sEnum = excel.ReadEnum(readCellValue(row, names, "%s"))
`, v.name, v.name)
			}

		} else if v.tt == "struct" {
			elem1 += fmt.Sprintf(`            e.%sStruct = read%sStruct(row, names)
`, v.name, v.name)
		} else if v.tt == "array1" {
			elem1 += fmt.Sprintf(`            e.%sArray = read%sArray(row, names)
`, v.name, v.name)
		} else if v.tt == "array2" {
			elem1 += fmt.Sprintf(`            e.%s = read%s(row, names)
`, v.name, v.name)
		}
	}

	out := fmt.Sprintf(template, tab.Name, tab.ReadPath, tab.Name, tmp, tab.Name, elem1, rows[1][0], app, cp)
	f.WriteString(out)
}

func gen(tab *TabDef, loadPath string) {

	tabRows, err := loadTable(path.Join(loadPath, tab.ReadPath, tab.Name))
	if err != nil {
		fmt.Printf("load table path:%s failed:%s\n", loadPath, err.Error())
		return
	}

	os.MkdirAll(tab.WritePath+"/"+tab.Name, os.ModePerm)
	out_path := fmt.Sprintf("%s/%s/%s.go", tab.WritePath, tab.Name, tab.Name)

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

	os.Truncate(out_path, 0)

	genHead(tab, f)
	genEntity(tabRows, tab, f)
	genProcess(tab, f)
	genIndex(tabRows, tab, f)
	genLoad(tabRows, tab, f)

	f.Close()
}

func loadTable(path string) (ret [][]string, err error) {
	var xlsx *excelize.File
	xlsx, err = excelize.OpenFile(path + ".xlsx")
	if err != nil {
		return
	} else {
		ret = xlsx.GetRows(xlsx.GetSheetName(xlsx.GetActiveSheetIndex()))
		return ret[:2], nil
	}
}

func newTable(sheetName, readPath, writePath string, defTypes ...string) *TabDef {
	tab := &TabDef{
		Name:      sheetName,
		ReadPath:  readPath,
		WritePath: writePath,
		Indexs:    make([]EleType, 0),
		Process:   make([]EleType, 0),
		Define:    map[string]string{},
	}

	for _, v := range defTypes {
		processCol(tab, v)
	}

	return tab
}
