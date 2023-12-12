package utils

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

type LoggerOption struct {
	Tag      string
	LogLevel *zapcore.Level
}

func NewLogger(option LoggerOption) (*zap.Logger, error) {
	if option.LogLevel == nil {
		if os.Getenv("ENV") == "production" {
			*option.LogLevel = zap.InfoLevel
		} else {
			*option.LogLevel = zap.DebugLevel
		}
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
