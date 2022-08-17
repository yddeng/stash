package main

import (
	"bytes"
	"fmt"
	"github.com/BurntSushi/toml"
	"initialthree/node/common/config"
	"io/ioutil"
	"os"
	"strings"
)

func splitPort(address string) string {
	return strings.Split(address, ":")[1]
}

func replaceIP(conf *config.Config, ip string) {
	for _, v := range conf.Dir {
		v.ExternalAddr = fmt.Sprintf("%s:%s", ip, splitPort(v.ExternalAddr))
	}
	for _, v := range conf.Login {
		v.ExternalAddr = fmt.Sprintf("%s:%s", ip, splitPort(v.ExternalAddr))
	}
	for _, v := range conf.Gate {
		v.ExternalAddr = fmt.Sprintf("%s:%s", ip, splitPort(v.ExternalAddr))
	}
	for _, v := range conf.GMMgr {
		v.WebAddr = fmt.Sprintf("%s:%s", ip, splitPort(v.WebAddr))
	}
	for _, v := range conf.WebService {
		v.WebAddress = fmt.Sprintf("%s:%s", ip, splitPort(v.WebAddress))
	}
}

func main() {
	filename := os.Args[1]
	ip := os.Args[2]

	conf := &config.Config{}
	_, _ = toml.DecodeFile(filename, conf)

	replaceIP(conf, ip)

	buf := new(bytes.Buffer)
	_ = toml.NewEncoder(buf).Encode(conf)
	_ = ioutil.WriteFile(filename, buf.Bytes(), os.ModePerm)
}
