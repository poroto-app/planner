package utils

import "go.uber.org/zap"

type LoggerOption struct {
	Tag        string
	Production bool
}

func NewLogger(option LoggerOption) (*zap.Logger, error) {
	var logger *zap.Logger
	if option.Production {
		l, err := zap.NewProduction()
		if err != nil {
			return nil, err
		}

		logger = l
	} else {
		l, err := zap.NewDevelopment()
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
