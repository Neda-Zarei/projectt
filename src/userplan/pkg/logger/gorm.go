package logger

import (
	"context"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm/logger"
)

type zapGormLogger struct {
	log *zap.Logger
}

func NewZapGormLogger(log *zap.Logger) logger.Interface {
	return &zapGormLogger{log: log}
}

func (l *zapGormLogger) LogMode(level logger.LogLevel) logger.Interface {
	// You can implement different log levels if you want
	return l
}

func (l *zapGormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	l.log.Sugar().Infow(msg, "data", data)
}

func (l *zapGormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	l.log.Sugar().Warnw(msg, "data", data)
}

func (l *zapGormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	l.log.Sugar().Errorw(msg, "data", data)
}

func (l *zapGormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	sql, rows := fc()
	latency := time.Since(begin)

	fields := []zap.Field{
		zap.Duration("latency", latency),
		zap.String("latency_human", latency.String()),
		zap.String("sql", sql),
		zap.Int64("rows", rows),
	}

	if err != nil {
		l.log.Error("gorm query failed", append(fields, zap.Error(err))...)
	} else {
		l.log.Info("gorm query executed", fields...)
	}
}
