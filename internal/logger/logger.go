package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	logger *zap.SugaredLogger
	level  zapcore.Level
)

func SetLevel(l zapcore.Level) {
	level = l
}

func GetLogger() *zap.SugaredLogger {
	if logger == nil {
		logger = newLogger(level)
	}

	return logger
}

func newLogger(level zapcore.Level) *zap.SugaredLogger {
	writer := zapcore.AddSync(os.Stdout)

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoder := zapcore.NewConsoleEncoder(encoderConfig)

	core := zapcore.NewCore(encoder, writer, level)
	logger := zap.New(core, zap.AddCaller())
	return logger.Sugar()
}
