package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewZapLogger(isDevEnv bool) *zap.Logger {
	logLevel := zap.InfoLevel
	if isDevEnv {
		logLevel = zap.DebugLevel
	}

	level := zap.NewAtomicLevelAt(logLevel)
	encoder := zapcore.NewJSONEncoder(jsonEncoderConfig())
	stdout := zapcore.AddSync(os.Stdout)

	core := zapcore.NewCore(encoder, stdout, level)

	opts := []zap.Option{zap.AddCaller()}
	if isDevEnv {
		opts = append(opts, zap.Development())
	}

	return zap.New(core, opts...)
}

func consoleEncoderConfig() zapcore.EncoderConfig {
	cfg := zap.NewDevelopmentEncoderConfig()
	cfg.EncodeTime = zapcore.ISO8601TimeEncoder
	cfg.TimeKey = "timestamp"
	cfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
	return cfg
}

func jsonEncoderConfig() zapcore.EncoderConfig {
	cfg := zap.NewProductionEncoderConfig()
	cfg.EncodeTime = zapcore.ISO8601TimeEncoder
	cfg.TimeKey = "timestamp"
	return cfg
}
