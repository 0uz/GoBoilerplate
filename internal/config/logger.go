package config

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"time"

	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
)

const (
	envAppEnvironment = "APP_ENV"
	envProdValue      = "PROD"

	defaultSlowThreshold = time.Second // 1 second
	defaultSourceField   = "file"
)

type Logger struct {
	*slog.Logger
}

var logger *Logger

func NewLogger() *Logger {
	if logger != nil {
		return logger
	}

	var level slog.Level
	if os.Getenv(envAppEnvironment) == envProdValue {
		level = slog.LevelInfo
	} else {
		level = slog.LevelDebug
	}

	l := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     level,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.SourceKey {
				a.Key = "caller"
			}
			return a
		},
	}))

	logger = &Logger{l}
	return logger
}

type GormLogger struct {
	SlowThreshold         time.Duration
	SourceField           string
	SkipErrRecordNotFound bool
	Debug                 bool
}

func NewGormLogger() gormlogger.Interface {
	return &GormLogger{
		SlowThreshold:         defaultSlowThreshold,
		SourceField:           defaultSourceField,
		SkipErrRecordNotFound: true,
		Debug:                 true,
	}
}

func (l *GormLogger) LogMode(gormlogger.LogLevel) gormlogger.Interface {
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
	var fields []interface{}
	if l.SourceField != "" {
		fields = append(fields, slog.String(l.SourceField, utils.FileWithLineNum()))
	}
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
