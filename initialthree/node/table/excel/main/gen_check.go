package main

import (
	"fmt"
	"initialthree/node/table/excel"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		panic("need argument: excel_path")
	}

	excel_path := os.Args[1]
	if excel_path == "" {
		panic("need argument: excel_path")
	}

	fmt.Println("gen_check start")

	tmp := `package main

import (
    "fmt"
    "initialthree/node/table/excel"
    "initialthree/zaplogger"
    "os"
    "runtime"
%s
)

func main() {
    if len(os.Args) < 2 {
        panic("need argument: excel_path")
    }

    logger := zaplogger.NewZapLogger("check.log","log","debug", 100, 14, 10, true)
    zaplogger.InitLogger(logger)

    excel_path := os.Args[1]
    if excel_path == "" {
        panic("need argument: excel_path")
    }
	defer func() {
		if r := recover(); r != nil {
			buf := make([]byte, 65535)
			l := runtime.Stack(buf, false)
			fmt.Println("check failed :", r)
			fmt.Println(string(buf[:l]))
		} else {
			fmt.Println("check table ok")
		}
	}()
    
	excel.Load(excel_path)
}
`
	str := ""

	for _, k := range excel.GenCheck(excel_path, "ConstTable") {
		str += k
	}
	for _, k := range excel.GenCheck(excel_path, "DataTable") {
		str += k
	}

	out := fmt.Sprintf(tmp, str)
	out_path := "check/main.go"

	f, err := os.OpenFile(out_path, os.O_RDWR, os.ModePerm)
	if err != nil {
		if os.IsNotExist(err) {
			f, err = os.Create("check/main.go")
			if err != nil {
				fmt.Printf("create %s failed:%s", "check/main.go", err.Error())
				return
			}
		} else {
			fmt.Printf("open %s failed:%s", "check/main.go", err.Error())
			return
		}
	}
	_ = os.Truncate(out_path, 0)
	_, _ = f.WriteString(out)
	_ = f.Close()

	fmt.Println("gen_check end")
}
