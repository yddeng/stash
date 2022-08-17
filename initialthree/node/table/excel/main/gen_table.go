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
	fmt.Println("genTable start")
	excel.Gen(excel_path)
	fmt.Println("genTable end")
}
