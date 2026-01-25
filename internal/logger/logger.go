package logger

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLogger(lvl, env string) *zap.Logger {
	var logger *zap.Logger
	var err error
	var opt zap.Option
	defaultOpt := false
	defaultEnv := false
	switch lvl {
	case "debug":
		opt = zap.AddStacktrace(zapcore.DebugLevel)
	case "info":
		opt = zap.AddStacktrace(zapcore.InfoLevel)
	case "warn":
		opt = zap.AddStacktrace(zapcore.WarnLevel)
	case "error":
		opt = zap.AddStacktrace(zapcore.ErrorLevel)
	default:
		defaultOpt = true
		opt = zap.AddStacktrace(zapcore.WarnLevel)
	}

	switch env {
	case "prod":
		logger, err = zap.NewProduction(opt)
	case "dev":
		logger, err = zap.NewDevelopment(opt)
	default:
		defaultEnv = true
		logger, err = zap.NewProduction(opt)
	}

	if err != nil {
		panic(fmt.Errorf("fatal error when making logger: %w", err))
	}
	if defaultOpt {
		logger.Warn("wrong or empty lvl field")
	}
	if defaultEnv {
		logger.Warn("wrong or empty env field")
	}
	return logger
}
