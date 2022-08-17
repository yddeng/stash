package main

import (
	"fmt"
	"github.com/hqpko/hutils"
	"initialthree/cluster"
	"initialthree/cluster/addr"
	"initialthree/node/common/config"
	"initialthree/node/common/db"
	"initialthree/node/common/otherStart"
	"initialthree/node/common/serverType"
	"initialthree/node/common/signal"
	"initialthree/node/common/timeDisposal"
	"initialthree/node/common/uniLocker"
	_ "initialthree/node/node_rank/net"
	"initialthree/node/node_rank/rank"
	"initialthree/zaplogger"
	"os"
)

func main() {
	// 是否启动pprof
	otherStart.PProfHasAndRun(os.Args)
	if len(os.Args) < 4 {
		fmt.Printf("usage ./node_rank config\n")
		return
	}

	conf := hutils.Must(config.LoadConfig(os.Args[1])).(*config.Config)
	ret := hutils.Must(config.GetConfig(conf, os.Args[2], os.Args[3])).(*config.Rank)
	_addr := hutils.Must(ret.MakeAddr()).(addr.Addr)

	filename := fmt.Sprintf("node_rank_%s.log", _addr.Logic.String())
	logger := zaplogger.NewZapLogger(filename, conf.Log.Path, conf.Log.Level, conf.Log.MaxSize, conf.Log.MaxAge, conf.Log.MaxBackups, conf.Log.EnableLogStdout)
	zaplogger.InitLogger(logger)
	cluster.InitLogger(logger)

	db.FlyfishInit(conf.Common.DbConfig, logger)
	hutils.Must(nil, timeDisposal.Init(db.GetFlyfishClient("global")))

	config.MustStartCluster(conf, _addr, serverType.Rank, uniLocker.New(db.GetFlyfishClient("nodelock")))

	hutils.Must(nil, rank.InitDBClient(ret))
	hutils.Must(nil, rank.InitManager())

	//testCmd()

	zaplogger.GetSugar().Infof("receive signal:%s to shutdown", <-signal.ListenStop())
	rank.Shutdown()
	cluster.Stop(func() {}, true)
	zaplogger.GetSugar().Info("rank stop ok")
}
