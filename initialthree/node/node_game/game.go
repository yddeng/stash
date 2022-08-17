package node_game

import (
	"github.com/hqpko/hutils"
	"initialthree/cluster"
	"initialthree/cluster/addr"
	"initialthree/node/common/config"
	"initialthree/node/common/db"
	"initialthree/node/common/idGenerator"
	"initialthree/node/common/otherStart"
	"initialthree/node/common/serverType"
	"initialthree/node/common/signal"
	"initialthree/node/common/timeDisposal"
	"initialthree/node/common/uniLocker"
	"initialthree/node/common/wordsFilter"
	"initialthree/node/node_game/global/bigSecret"
	"initialthree/node/node_game/global/mail"
	"initialthree/node/node_game/global/scarsIngrain"
	_ "initialthree/node/node_game/module"
	"initialthree/node/node_game/monitor"
	_ "initialthree/node/node_game/transaction"
	"initialthree/node/node_game/user"
	"initialthree/node/table/excel"
	"initialthree/zaplogger"
	"math/rand"
	"time"
)

var conf *config.Config

func ReloadTable() {
	excel.Reload(conf.Common.GetExcelPath())
}

func Start(address addr.Addr, conf_ *config.Config) error {

	conf = conf_

	logger := zaplogger.GetSugar()

	hutils.Must(nil, db.FlyfishInit(conf.Common.DbConfig, zaplogger.GetLogger()))

	//加载配置文件
	excel.Load(conf.Common.GetExcelPath())

	config.MustStartCluster(conf, address, serverType.Game, uniLocker.New(db.GetFlyfishClient("nodelock")))

	hutils.Must(nil, timeDisposal.Init(db.GetFlyfishClient("global")))

	hutils.Must(nil, wordsFilter.Init(conf.Common.GetWordsFilterPath()))

	//初始化随机种子
	rand.Seed(time.Now().Unix())

	scarsIngrain.Launch()
	bigSecret.Launch()

	hutils.Must(nil, mail.InitGlobalMail(db.GetFlyfishClient("global")))

	user.StartLBCollector()

	// gameID 生成规则
	idGenerator.Register("game_id", db.GetFlyfishClient("global"), func(i int64) int64 {
		return 100000 + i
	})

	// prometheus
	if address, ok := otherStart.Has(otherStart.PROMETHEUS); ok {
		logger.Infof("has prometheus, address %s", address)
		monitor.StartPrometheus(address)
	}

	return nil
}

func WaitStop() {

	logger := zaplogger.GetSugar()

	defer func() {
		logger.Info("stop ok")
	}()

	for {
		select {
		case <-signal.ListenStop():
			logger.Infof("receive shutdown")
			//服务要下线,完成停机前通告center移除当前节点
			cluster.Stop(func() {
				user.StopAndWaitAllUserLogout(false)
			}, true)
			return
		case <-signal.ListenSigUpdateProcess():
			logger.Infof("receive update server")
			cluster.Stop(func() {
				user.StopAndWaitAllUserLogout(true)
			})
			return
		case <-signal.ListenSigUpdateRes():
			ReloadTable()
			logger.Infof("receive update res")
		}
	}

}
