// Package logger provides functions to set up a new logger
package logger

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	logFormatText    = "text"
	logFormatConsole = "console"
)

// Logger is a small wrapper around a zap.Logger.
type Logger struct {
	*zap.Logger
}

// New creates a new Logger with given logLevel and logFormat as part of a permanent field of the logger.
func New(logLevel, logFormat string) (*Logger, error) {
	if logFormat == logFormatText {
		logFormat = logFormatConsole
	}

	zapConfig := zap.NewProductionConfig()
	zapConfig.Encoding = logFormat

	var level zapcore.Level
	err := level.UnmarshalText([]byte(logLevel))
	if err != nil {
		return nil, err
	}
	zapConfig.Level = zap.NewAtomicLevelAt(level)
	zapConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	logger, err := zapConfig.Build()
	if err != nil {
		return nil, fmt.Errorf("could not build logger: %w", err)
	}

	zap.ReplaceGlobals(logger)

	return &Logger{Logger: logger}, nil
}
