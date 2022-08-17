package robot

import (
	"go.uber.org/zap"
)

var zapLogger *zap.Logger
var sugaredLogger *zap.SugaredLogger

func initLogger(logger *zap.Logger) {
	zapLogger = logger.WithOptions(zap.AddCallerSkip(+1))
	sugaredLogger = zapLogger.Sugar()
}
