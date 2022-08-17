package main

import (
	"fmt"
	"github.com/hqpko/hutils"
	"initialthree/cluster"
	"initialthree/cluster/addr"
	"initialthree/node/common/config"
	"initialthree/node/common/otherStart"
	"initialthree/pkg"
	"os"

	"initialthree/node/common/serverType"
	"initialthree/node/common/signal"
	"initialthree/node/node_world/world"
	"initialthree/util"
	"initialthree/zaplogger"
)

func main() {
	// 是否启动pprof
	otherStart.PProfHasAndRun(os.Args)
	if len(os.Args) < 4 {
		fmt.Printf("usage ./node_world config\n")
		return
	}

	conf := hutils.Must(config.LoadConfig(os.Args[1])).(*config.Config)
	ret := hutils.Must(config.GetConfig(conf, os.Args[2], os.Args[3])).(*config.World)
	_addr := hutils.Must(ret.MakeAddr()).(addr.Addr)

	filename := fmt.Sprintf("node_world_%s.log", _addr.Logic.String())
	logger := zaplogger.NewZapLoggerWithConfig(filename, conf.Log)
	zaplogger.InitLogger(logger)
	kendynet.InitLogger(zaplogger.GetSugar())
	cluster.InitLogger(zaplogger.GetSugar())
	world.InitLogger(zaplogger.GetSugar())

	world.InitWorld(_addr.Logic.Server())
	//obj.Init(logger, _addr.Logic.Server())

	config.MustStartCluster(conf, _addr, serverType.World)

	zaplogger.GetSugar().Infof("receive signal:%s to shutdown", <-signal.ListenStop())

	zaplogger.GetSugar().Info("stop ok")
}
