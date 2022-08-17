package main

import (
	"bytes"
	"fmt"
	"github.com/BurntSushi/toml"
	"initialthree/node/common/config"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

/*****/
type Servers struct {
	ServerList []*Server `toml:"Servers"`
}

type Server struct {
	ServerId   int32  `toml:"serverID"`
	ServerName string `toml:"serverName"`
}

var basePort int
var right int

func getPort() int {
	p := basePort
	if p == right {
		panic("端口超出分配范围")
	}
	basePort++
	return p
}

func makeStartSh(port int) {
	data, err := ioutil.ReadFile(starShPath)
	if err != nil {
		panic(err)
	}

	str := string(data)
	dataStr := strings.Replace(str, "8010", fmt.Sprintf("%d", port), -1)
	ioutil.WriteFile("start.sh", []byte(dataStr), os.ModePerm)
}

func makeIP(conf *config.Config, ip string, l, r, flyfishPort int) {
	basePort = l
	right = r

	var port int
	for i := range conf.Common.CenterAddr {
		port = getPort()
		conf.Common.CenterAddr[i] = fmt.Sprintf("localhost:%d", port)
		makeStartSh(port)
	}
	for i, str := range conf.Common.DbConfig {
		prefix := strings.Split(str, ":")[0]
		conf.Common.DbConfig[i] = fmt.Sprintf("%s:%d", prefix, flyfishPort)
	}
	for _, v := range conf.Dir {
		port = getPort()
		v.ClusterAddr = fmt.Sprintf("localhost:%d", port)
		port = getPort()
		v.ExternalAddr = fmt.Sprintf("%s:%d", ip, port)
	}
	for _, v := range conf.Login {
		port = getPort()
		v.ClusterAddr = fmt.Sprintf("localhost:%d", port)
		port = getPort()
		v.ExternalAddr = fmt.Sprintf("%s:%d", ip, port)
	}
	for _, v := range conf.Gate {
		port = getPort()
		v.ClusterAddr = fmt.Sprintf("localhost:%d", port)
		port = getPort()
		v.ExternalAddr = fmt.Sprintf("%s:%d", ip, port)
	}
	for _, v := range conf.Game {
		port = getPort()
		v.ClusterAddr = fmt.Sprintf("localhost:%d", port)
	}
	for _, v := range conf.Rank {
		port = getPort()
		v.ClusterAddr = fmt.Sprintf("localhost:%d", port)
	}
	for _, v := range conf.World {
		port = getPort()
		v.ClusterAddr = fmt.Sprintf("localhost:%d", port)
	}
	for _, v := range conf.Map {
		port = getPort()
		v.ClusterAddr = fmt.Sprintf("localhost:%d", port)
	}
	for _, v := range conf.Team {
		port = getPort()
		v.ClusterAddr = fmt.Sprintf("localhost:%d", port)
	}
	for _, v := range conf.GS {
		port = getPort()
		v.ClusterAddr = fmt.Sprintf("localhost:%d", port)
	}
	for _, v := range conf.GMMgr {
		port = getPort()
		v.ClusterAddr = fmt.Sprintf("localhost:%d", port)
		port = getPort()
		v.WebAddr = fmt.Sprintf("%s:%d", ip, port)
	}
	for _, v := range conf.WebService {
		port = getPort()
		v.ClusterAddr = fmt.Sprintf("localhost:%d", port)
		v.WebAddress = fmt.Sprintf("%s:%d", ip, port)
	}
}

const (
	nodeConfigPath = "config.toml.template"
	nodeDirPath    = "node_dir/config/config.toml.template"
	starShPath     = "start.sh.template"
)

var ()

func main() {
	if len(os.Args) < 4 {
		fmt.Println("need ip start_port flyfish_port, optional ")
		return
	}
	ip := os.Args[1]
	port, err := strconv.Atoi(os.Args[2])
	if err != nil {
		panic(err)
	}
	flyfish_port, err := strconv.Atoi(os.Args[3])
	if err != nil {
		panic(err)
	}

	conf := &config.Config{}
	_, _ = toml.DecodeFile(nodeConfigPath, conf)

	ser := &Servers{}
	_, _ = toml.DecodeFile(nodeDirPath, ser)

	ser.ServerList[0].ServerName = "服务器1"
	makeIP(conf, ip, port, port+100, flyfish_port)

	// install 编译
	if len(os.Args) >= 4 {
		conf.Common.CfgPathRoot = "../configs"
		for _, v := range conf.Dir {
			v.FilePath = "dir/config.toml"
		}
	}

	buf := new(bytes.Buffer)
	_ = toml.NewEncoder(buf).Encode(conf)
	_ = ioutil.WriteFile("config.toml", buf.Bytes(), os.ModePerm)

	buf.Reset()
	_ = toml.NewEncoder(buf).Encode(ser)
	_ = ioutil.WriteFile("node_dir/config/config.toml", buf.Bytes(), os.ModePerm)

}
