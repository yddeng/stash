package main

import (
	"fmt"
	"github.com/hqpko/hutils"
	"initialthree/cluster"
	addr2 "initialthree/cluster/addr"
	"initialthree/node/common/config"
	"initialthree/node/common/db"
	"initialthree/node/common/otherStart"
	"initialthree/node/common/serverType"
	"initialthree/node/common/signal"
	"initialthree/node/common/uniLocker"
	_ "initialthree/node/node_webservice/handler"
	"initialthree/node/node_webservice/webservice"
	"initialthree/zaplogger"
	"os"
)

func main() {
	// 是否启动pprof
	otherStart.PProfHasAndRun(os.Args)
	if len(os.Args) < 4 {
		fmt.Printf("usage: %s config.toml WebService idx ", os.Args[0])
	}

	cfg := hutils.Must(config.LoadConfig(os.Args[1])).(*config.Config)
	ret := hutils.Must(config.GetConfig(cfg, os.Args[2], os.Args[3])).(*config.WebService)
	addr := hutils.Must(ret.MakeAddr()).(addr2.Addr)

	filename := fmt.Sprintf("node_webservice_%s.log", addr.Logic.String())
	logger := zaplogger.NewZapLogger(filename, cfg.Log.Path, cfg.Log.Level, cfg.Log.MaxSize, cfg.Log.MaxAge, cfg.Log.MaxBackups, cfg.Log.EnableLogStdout)
	cluster.InitLogger(logger)

	hutils.Must(nil, db.FlyfishInit(cfg.Common.DbConfig, zaplogger.GetLogger()))

	config.MustStartCluster(cfg, addr, serverType.WebService, uniLocker.New(db.GetFlyfishClient("nodelock")))

	zaplogger.GetSugar().Infof("web service run.")

	hutils.Must(nil, webservice.StartWebService(ret))

	zaplogger.GetSugar().Infof("web service stop: %s", <-signal.ListenStop())

	cluster.Stop(func() {}, true)
	zaplogger.GetSugar().Info("shutdown.")
}
