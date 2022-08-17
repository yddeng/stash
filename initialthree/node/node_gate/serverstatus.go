package node_gate

import (
	"initialthree/node/common/db"
	"initialthree/node/common/omm/serverstatus"
)

var serverStatus = serverstatus.DefaultStatue

func initServerStatus() error {
	serverstatus.RegisterUpdateNotify(db.GetFlyfishClient("global"), func(status serverstatus.Status) {
		serverStatus = status
	})
	return nil
}

func isServerOpen() bool {
	return serverStatus == serverstatus.StatusOpen
}
