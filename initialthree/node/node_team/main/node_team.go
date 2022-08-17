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
	_ "initialthree/node/node_team/net"
	_ "initialthree/node/node_team/net/rpc"
	_ "initialthree/node/node_team/net/ss"
	_ "initialthree/node/node_team/team"
	"initialthree/util"
)

func main() {
	// 是否启动pprof
	otherStart.PProfHasAndRun(os.Args)
	if len(os.Args) < 4 {
		fmt.Printf("usage ./node_team config\n")
		return
	}

	conf := hutils.Must(config.LoadConfig(os.Args[1])).(*config.Config)
	ret := hutils.Must(config.GetConfig(conf, os.Args[2], os.Args[3])).(*config.Team)
	_addr := hutils.Must(ret.MakeAddr()).(addr.Addr)

	filename := fmt.Sprintf("node_team_%s", _addr.Logic.String())
	logger := util.NewLogger("log", filename, 1024*1024*50)
	kendynet.InitLogger(logger)
	cluster.InitLogger(logger)

	config.MustStartCluster(conf, _addr, serverType.Team)

	logger.Infof("receive signal:%s to shutdown", <-signal.ListenStop())
	logger.Infoln("stop ok")
}
