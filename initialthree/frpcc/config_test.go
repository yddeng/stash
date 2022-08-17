package main

import (
	"fmt"
	"testing"
)

func TestConfig(t *testing.T) {

	conf, err := LoadConfigStr(`
Proxy = [
	{Name="pmp",SvrAddr="212.129.131.27",SvrPort=41098},
	{Name="sniper-ssh",SvrAddr="212.129.131.27",SvrPort=41097},
]		

`)

	fmt.Println(*conf, err)

}
