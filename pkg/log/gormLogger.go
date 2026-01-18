package log

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

var logger *Logger

type GormLogger struct {
	SlowThreshold         time.Duration
	SkipErrRecordNotFound bool
	Debug                 bool
}

func NewGormLogger(l *Logger) gormLogger.Interface {
	logger = l
	return &GormLogger{
		SlowThreshold:         10 * time.Second,
		SkipErrRecordNotFound: true,
		Debug:                 true,
	}
}

func (l *GormLogger) LogMode(gormLogger.LogLevel) gormLogger.Interface {
	return l
}

func (l *GormLogger) Info(ctx context.Context, s string, args ...interface{}) {
	logger.InfoContext(ctx, s, args...)
}

func (l *GormLogger) Warn(ctx context.Context, s string, args ...interface{}) {
	logger.WarnContext(ctx, s, args...)
}

func (l *GormLogger) Error(ctx context.Context, s string, args ...interface{}) {
	logger.ErrorContext(ctx, s, args...)
}

func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)
	sql, _ := fc()
	var fields []any
	
	if err != nil && !(errors.Is(err, gorm.ErrRecordNotFound) && l.SkipErrRecordNotFound) {
		fields = append(fields, slog.Any("error", err))
		logger.ErrorContext(ctx, sql, fields...)
		return
	}

	if l.SlowThreshold != 0 && elapsed > l.SlowThreshold {
		logger.WarnContext(ctx, sql, fields...)
		return
	}

	if l.Debug {
		logger.DebugContext(ctx, sql, fields...)
	}
}