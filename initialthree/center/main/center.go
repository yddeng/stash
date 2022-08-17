package main

import (
	"fmt"
	fnet "github.com/sniperHW/flyfish/pkg/net"
	"go.uber.org/zap"
	"initialthree/center"
	"initialthree/zaplogger"
	"os"
	"os/signal"
	"syscall"
)

var logger *zap.Logger

func main() {
	filename := fmt.Sprintf("center_%s.log", os.Args[1])
	logger = zaplogger.NewZapLogger(filename, "log", "debug", 100, 14, 10, true)
	zaplogger.InitLogger(logger)
	fnet.InitLogger(logger)
	if err := center.Start(os.Args[1], logger); err != nil {
		panic(err)
	}

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)
	zaplogger.GetSugar().Infof("receive signal:%s to stop.", <-stopChan)
}
