package logger

import (
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	serviceName = "cart"

	global       *zap.SugaredLogger
	defaultLevel = zap.NewAtomicLevelAt(zap.InfoLevel)

	once = new(sync.Once)
)

func init() {
	initLogger(defaultLevel.Level())
}

func initLogger(level zapcore.Level) {
	once.Do(func() {
		config := zap.NewProductionConfig()
		config.OutputPaths = []string{"stdout"}
		config.ErrorOutputPaths = []string{"stdout"}
		config.Level.SetLevel(level)
		l, err := config.Build(zap.AddCallerSkip(1))
		if err != nil {
			panic("failed to build logger")
		}
		global = l.Sugar()
	})
}

func getLogger() *zap.SugaredLogger {
	if global == nil {
		initLogger(defaultLevel.Level())
	}
	return global
}

func Fatalw(msg string, keysAndValues ...interface{}) {
	getLogger().With(
		zap.String("service", serviceName),
	).Fatalw(msg, keysAndValues...)
}

func Errorw(msg string, keysAndValues ...interface{}) {
	getLogger().With(
		zap.String("service", serviceName),
	).Errorw(msg, keysAndValues...)
}

func Infow(msg string, keysAndValues ...interface{}) {
	getLogger().With(
		zap.String("service", serviceName),
	).Infow(msg, keysAndValues...)
}

func Sync() error {
	return getLogger().Sync()
}
