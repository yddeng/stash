package main

import (
	"fmt"
	"github.com/hqpko/hutils"
	"initialthree/cluster/addr"
	"initialthree/node/common/otherStart"
	"os"

	"initialthree/node/common/serverType"
	"initialthree/node/common/signal"
	table "initialthree/node/table/excel"

	"initialthree/cluster"
	"initialthree/node/common/config"
	"initialthree/node/node_map/role"
	_ "initialthree/node/node_map/role"
	"initialthree/node/node_map/scene"
	_ "initialthree/protocol/cs"
	"initialthree/util"

	"initialthree/pkg"
)

func main() {
	// 是否启动pprof
	otherStart.PProfHasAndRun(os.Args)
	if len(os.Args) < 4 {
		fmt.Printf("usage ./node_map config\n")
		return
	}

	conf := hutils.Must(config.LoadConfig(os.Args[1])).(*config.Config)
	ret := hutils.Must(config.GetConfig(conf, os.Args[2], os.Args[3])).(*config.Map)
	_addr := hutils.Must(ret.MakeAddr()).(addr.Addr)

	filename := fmt.Sprintf("node_map_%s", _addr.Logic.String())
	logger := util.NewLogger("log", filename, 1024*1024*50)
	kendynet.InitLogger(logger)
	cluster.InitLogger(logger)
	role.InitLogger(logger)

	config.MustStartCluster(conf, _addr, serverType.Map)

	table.Load(conf.Common.GetExcelPath())

	if err := scene.Init(logger, ret); err != nil {
		logger.Errorf("scene.Init() failed:%s\n", err)
		return
	}

	logger.Infof("receive signal:%s to shutdown", <-signal.ListenStop())
	logger.Infoln("stop ok")
}
