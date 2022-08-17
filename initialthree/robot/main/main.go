package main

import (
	"fmt"
	"initialthree/robot"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usage: %s configfile\n", os.Args[0])
		os.Exit(1)
	}

	robot.Start(os.Args[1])
}
