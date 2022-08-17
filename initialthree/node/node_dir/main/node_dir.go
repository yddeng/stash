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
	"initialthree/node/common/uniLocker"
	"initialthree/node/node_dir"
	"initialthree/node/node_dir/server"
	"initialthree/zaplogger"
	"os"
)

func main() {
	// 是否启动pprof
	otherStart.PProfHasAndRun(os.Args)
	if len(os.Args) < 4 {
		fmt.Printf("usage ./logger config\n")
		return
	}

	conf := hutils.Must(config.LoadConfig(os.Args[1])).(*config.Config)
	ret := hutils.Must(config.GetConfig(conf, os.Args[2], os.Args[3])).(*config.Dir)
	_addr := hutils.Must(ret.MakeAddr()).(addr.Addr)

	filename := fmt.Sprintf("node_dir_%s.log", _addr.Logic.String())
	logger := zaplogger.NewZapLogger(filename, conf.Log.Path, conf.Log.Level, conf.Log.MaxSize, conf.Log.MaxAge, conf.Log.MaxBackups, conf.Log.EnableLogStdout)
	cluster.InitLogger(logger)
	hutils.Must(nil, db.FlyfishInit(conf.Common.DbConfig, logger))

	config.MustStartCluster(conf, _addr, serverType.Dir, uniLocker.New(db.GetFlyfishClient("nodelock")))

	hutils.Must(nil, server.LoadConfig(ret))

	server.InitDBQueryMax(ret.BBQueryMax)
	hutils.Must(nil, node_dir.Start(ret.ExternalAddr))

	zaplogger.GetSugar().Infof("receive signal:%s to shutdown", <-signal.ListenStop())

	cluster.Stop(func() {}, true)
}
