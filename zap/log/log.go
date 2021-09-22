package log

import "go.uber.org/zap"

var logger *zap.Logger

func init() {
	l, err := zap.NewDevelopment(zap.AddCallerSkip(1))
	if err != nil {
		panic(err)
	}
	logger = l
}

func Info(msg string) {
	logger.Info(msg)
}

func Error(msg string) {
	logger.Error(msg)
}
