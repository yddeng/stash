package main

import (
	"fmt"
	"initialthree/node/common/wordsFilter"
)

func main() {
	err := wordsFilter.Init("/Users/yidongdeng/go/src/initialthree/node/configs/wordsFilter/wordsFilter.txt")
	//err := wordsFilter.Init("wordsFilter.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println('\r', '\n')

	fmt.Println(wordsFilter.Replace("4564==习+-近平=+"), wordsFilter.Check("/*习=近平++"))
	fmt.Println(wordsFilter.Replace("4564==毛主席=+"), wordsFilter.Check("/*毛主席++"))
	fmt.Println(wordsFilter.Replace("fuc++k"), wordsFilter.Check("fu-+ck"))
	fmt.Println(wordsFilter.Replace("fuc++k you"), wordsFilter.Check("fu-+ck you"))
	fmt.Println(wordsFilter.Replace("你是==狗屎++"), wordsFilter.Check("你是狗+-/屎"))
	fmt.Println(wordsFilter.Replace("你是==狗屎//臭狗屎"), wordsFilter.Check("你是--狗+-/屎臭狗屎"))
}
