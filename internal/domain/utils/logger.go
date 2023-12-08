package utils

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

type LoggerOption struct {
	Tag string
}

func NewLogger(option LoggerOption) (*zap.Logger, error) {
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
