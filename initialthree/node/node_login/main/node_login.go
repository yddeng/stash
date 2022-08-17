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
	"initialthree/node/node_login"
	"initialthree/zaplogger"
	"os"
)

func main() {
	// 是否启动pprof
	otherStart.PProfHasAndRun(os.Args)
	if len(os.Args) < 4 {
		fmt.Printf("usage ./node_login config\n")
		return
	}

	conf := hutils.Must(config.LoadConfig(os.Args[1])).(*config.Config)
	ret := hutils.Must(config.GetConfig(conf, os.Args[2], os.Args[3])).(*config.Login)
	_addr := hutils.Must(ret.MakeAddr()).(addr.Addr)

	filename := fmt.Sprintf("node_login_%s.log", _addr.Logic.String())
	logger := zaplogger.NewZapLogger(filename, conf.Log.Path, conf.Log.Level, conf.Log.MaxSize, conf.Log.MaxAge, conf.Log.MaxBackups, conf.Log.EnableLogStdout)
	cluster.InitLogger(logger)
	db.FlyfishInit(conf.Common.DbConfig, logger)

	//hutils.Must(nil, logictimesys.Launch(conf.Common.LogicTimeSysCfgPath, db.GetFlyfishClient("game"), zaplogger.GetSugar()))

	config.MustStartCluster(conf, _addr, serverType.Login, uniLocker.New(db.GetFlyfishClient("nodelock")))

	hutils.Must(nil, node_login.Start(ret.ExternalAddr, conf.Common))

	zaplogger.GetSugar().Infof("receive signal:%s to shutdown", <-signal.ListenStop())

	cluster.Stop(func() {}, true)
	zaplogger.GetSugar().Info("stop ok")
}
