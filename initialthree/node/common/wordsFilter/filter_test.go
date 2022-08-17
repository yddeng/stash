package wordsFilter

import (
	"fmt"
	"testing"
)

func TestNewFilter(t *testing.T) {
	fmt.Println(Replace("4564==毛主席=+"), Check("/*毛主席++"))

	err := Init("/Users/yidongdeng/go/src/initialthree/node/configs/wordsFilter/wordsFilter.txt")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(Replace("4564==习+-近平=+"), Check("/*习=近平++"))
	fmt.Println(Replace("4564==毛主席=+"), Check("/*毛主席++"))
	fmt.Println(Replace("fuc++k"), Check("fu-+ck"))
	fmt.Println(Replace("fuc++k you"), Check("fu-+ck you"))
	fmt.Println(Replace("hello"), Check("hello"))
	fmt.Println(Replace("你是==狗屎++"), Check("你是狗+-/屎"))
	fmt.Println(Replace("88草66"), Check("说草ss"))
	fmt.Println(Replace(""), Check(""))

	err = ReLoad("test/wordsFilter.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(Replace("4564==毛主席=+"), Check("/*毛主席++"))
	fmt.Println(Replace("你是==狗屎++"), Check("你是狗+-/屎"))
}
