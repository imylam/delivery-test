package logger

import (
	"go.uber.org/zap"
)

var Logger *zap.Logger

func Init() {
	zapLogger, _ := zap.NewProduction()
	defer zapLogger.Sync() // flushes buffer, if any

	Logger = zapLogger
}
