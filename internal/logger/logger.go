package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/vellalasantosh/wound_iq_api_new/internal/config"
)

func New(cfg *config.Config) *zap.Logger {
	var level zapcore.Level
	switch cfg.LogLevel {
	case "debug":
		level = zapcore.DebugLevel
	case "warn":
		level = zapcore.WarnLevel
	case "error":
		level = zapcore.ErrorLevel
	default:
		level = zapcore.InfoLevel
	}
	cfgZap := zap.NewProductionConfig()
	cfgZap.Encoding = "json"
	cfgZap.Level = zap.NewAtomicLevelAt(level)
	logger, _ := cfgZap.Build()
	return logger
}
