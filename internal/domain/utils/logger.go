package utils

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

const (
	defaultLogLevelDevelopment = zap.InfoLevel
	defaultLogLevelProduction  = zap.InfoLevel
)

type LoggerOption struct {
	Tag      string
	LogLevel *zapcore.Level
}

func NewLogger(option LoggerOption) (*zap.Logger, error) {
	if option.LogLevel == nil {
		var defaultLogLevel zapcore.Level
		if os.Getenv("ENV") == "production" {
			defaultLogLevel = defaultLogLevelProduction
		} else {
			defaultLogLevel = defaultLogLevelDevelopment
		}
		option.LogLevel = &defaultLogLevel
	}

	var logger *zap.Logger
	if os.Getenv("ENV") == "production" {
		l, err := zap.NewProduction()
		if err != nil {
			return nil, err
		}

		logger = l
	} else {
		dc := zap.NewDevelopmentConfig()
		dc.Encoding = "console"
		dc.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		dc.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		dc.DisableStacktrace = true
		dc.Level = zap.NewAtomicLevelAt(*option.LogLevel)

		l, err := dc.Build()
		if err != nil {
			return nil, err
		}

		logger = l
	}

	if option.Tag != "" {
		logger = logger.With(zap.String("tag", option.Tag))
	}

	return logger, nil
}

func LoggerLevelPointer(level zapcore.Level) *zapcore.Level {
	return &level
}
