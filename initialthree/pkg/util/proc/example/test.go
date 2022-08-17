package main

import (
	"fmt"
	"initialthree/util/proc"
)

func main() {

	procList, _ := proc.GetProcs("initialthree")

	for _, v := range procList {
		fmt.Println(v.User, v.Pid, v.CommandName, v.FullCommand)
	}

	//fmt.Println(procList)

}
