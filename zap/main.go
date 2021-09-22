package main

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/jun06t/go-sample/zap/log"
)

func main() {
	{
		logger := NewMultiOutputLogger()
		logger.Error("error message")
		logger.Info("info message")
		logger.Debug("debug message")
	}

	{
		var buf bytes.Buffer
		logger := NewWithWriter(&buf)
		logger.Error("error message")
		logger.Info("info message")
		logger.Debug("debug message")
		fmt.Println(buf.String())
	}

	{
		logger := NewCloudLoggingLogger()
		logger.Error("error message")
		logger.Info("info message")
		logger.Debug("debug message")
	}

	{
		log.Info("info message")
		log.Error("error message")
	}
}

func NewWithWriter(w io.Writer) *zap.Logger {
	sink := zapcore.AddSync(w)
	lsink := zapcore.Lock(sink)

	enc := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())

	core := zapcore.NewCore(enc, lsink, zapcore.InfoLevel)

	logger := zap.New(core)
	return logger
}

func NewMultiOutputLogger() *zap.Logger {
	highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})
	lowPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl < zapcore.ErrorLevel
	})
	stdoutSink := zapcore.Lock(os.Stdout)
	stderrSink := zapcore.Lock(os.Stderr)

	enc := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())

	core := zapcore.NewTee(
		zapcore.NewCore(enc, stderrSink, highPriority),
		zapcore.NewCore(enc, stdoutSink, lowPriority),
	)

	logger := zap.New(core)
	return logger
}

func NewCloudLoggingLogger() *zap.Logger {
	cfg := zap.NewProductionConfig()
	cfg.EncoderConfig = newEncoderConfig()
	logger, _ := cfg.Build()
	return logger
}

var logLevelSeverity = map[zapcore.Level]string{
	zapcore.DebugLevel:  "DEBUG",
	zapcore.InfoLevel:   "INFO",
	zapcore.WarnLevel:   "WARNING",
	zapcore.ErrorLevel:  "ERROR",
	zapcore.DPanicLevel: "CRITICAL",
	zapcore.PanicLevel:  "ALERT",
	zapcore.FatalLevel:  "EMERGENCY",
}

func EncodeLevel(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(logLevelSeverity[l])
}

func newEncoderConfig() zapcore.EncoderConfig {
	cfg := zap.NewProductionEncoderConfig()
	cfg.TimeKey = "timestamp"
	cfg.LevelKey = "severity"
	cfg.MessageKey = "message"
	cfg.EncodeLevel = EncodeLevel
	cfg.EncodeTime = zapcore.RFC3339NanoTimeEncoder

	return cfg
}
