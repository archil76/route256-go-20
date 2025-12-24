package logger

import (
	"fmt"

	"go.uber.org/zap"
)

const (
	serviceName = "loms"
)

var InfoLogger, FatalLogger *zap.Logger

func Info(format string, args ...interface{}) {
	if InfoLogger == nil {
		panic("InfoLogger is not initialized")
	}

	msg := fmt.Sprintf(format, args...)
	InfoLogger.With(
		zap.String("service", serviceName),
	).Info(msg)
}

func Error(format string, args ...interface{}) {
	if InfoLogger == nil {
		panic("InfoLogger is not initialized")
	}

	msg := fmt.Sprintf(format, args...)
	InfoLogger.With(
		zap.String("service", serviceName),
	).Error(msg)
}

func Fatal(format string, args ...interface{}) {
	if FatalLogger == nil {
		panic("FatalLogger is not initialized")
	}

	msg := fmt.Sprintf(format, args...)
	FatalLogger.With(
		zap.String("service", serviceName),
	).Fatal(msg)
}
