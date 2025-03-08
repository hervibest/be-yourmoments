package utils

import (
	"go.uber.org/zap"
)

func NewLogger() *zap.SugaredLogger {

	logger, _ := zap.NewProduction()
	sugar := logger.Sugar()

	return sugar
}
