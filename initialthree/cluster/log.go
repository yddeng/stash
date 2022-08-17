package cluster

import (
	fnet "github.com/sniperHW/flyfish/pkg/net"
	"go.uber.org/zap"
	"initialthree/zaplogger"
)

func InitLogger(l *zap.Logger) {
	logger = l
	zaplogger.InitLogger(l)
	fnet.InitLogger(l)
}
