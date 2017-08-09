package logger

import (
	"go.uber.org/zap"
)

var L *zap.SugaredLogger

func init() {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()
	L = logger.Sugar()
}
