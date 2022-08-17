package node_login

import (
	"initialthree/node/common/db"
	"initialthree/node/common/omm/serverstatus"
	"initialthree/zaplogger"
)

var serverStatus = serverstatus.DefaultStatue

func initServerStatus() error {
	serverstatus.RegisterUpdateNotify(db.GetFlyfishClient("global"), func(status serverstatus.Status) {
		zaplogger.GetSugar().Info("status update", serverStatus, status)
		serverStatus = status
	})
	return nil
}

func isServerOpen() bool {
	return serverStatus == serverstatus.StatusOpen
}
