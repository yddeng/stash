package main

import (
	"fmt"
	"github.com/hqpko/hutils"
	"initialthree/cluster"
	"initialthree/cluster/addr"
	"initialthree/node/common/config"
	"initialthree/node/common/otherStart"
	"initialthree/node/node_game"
	_ "initialthree/node/node_game/handler"
	"initialthree/zaplogger"
	"os"
)

func main() {
	if len(os.Args) < 4 {
		fmt.Printf("usage ./node_game config")
		return
	}

	// 是否启动pprof
	otherStart.PProfHasAndRun(os.Args)

	conf := hutils.Must(config.LoadConfig(os.Args[1])).(*config.Config)
	ret := hutils.Must(config.GetConfig(conf, os.Args[2], os.Args[3])).(*config.Game)
	address := hutils.Must(ret.MakeAddr()).(addr.Addr)

	filename := fmt.Sprintf("node_game_%s.log", address.Logic.String())
	logger := zaplogger.NewZapLogger(filename, conf.Log.Path, conf.Log.Level, conf.Log.MaxSize, conf.Log.MaxAge, conf.Log.MaxBackups, conf.Log.EnableLogStdout)
	cluster.InitLogger(logger)

	hutils.Must(nil, node_game.Start(address, conf))
	node_game.WaitStop()
}
