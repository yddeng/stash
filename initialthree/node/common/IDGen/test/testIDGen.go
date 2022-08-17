package main

import (
	"fmt"
	"initialthree/node/common/IDGen"
	"time"
)

func main() {
	g := IDGen.New(1, 2)
	startTime := time.Now().UnixNano() / 1e6
	for i := 0; i < 10000; i++ {
		id, err := g.Gen()
		if nil == err {
			fmt.Println(id)
		} else {
			time.Sleep(1)
		}
	}
	fmt.Println(time.Now().UnixNano()/1e6 - startTime)
}
