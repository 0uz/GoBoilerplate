package config

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"time"

	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/otel/log/global"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
)

const (
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

	env := "PROD"

	if cfg := Get(); cfg != nil {
		env = cfg.App.Environment
	}

	var l *slog.Logger

	if env == "DEV" {
		stdoutHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})

		otelHandler := otelslog.NewHandler("go-boilerplate",
			otelslog.WithLoggerProvider(global.GetLoggerProvider()),
		)

		multiHandler := &MultiHandler{
			handlers: []slog.Handler{stdoutHandler, otelHandler},
		}

		l = slog.New(multiHandler)
	} else {
		l = slog.New(otelslog.NewHandler("go-boilerplate",
			otelslog.WithLoggerProvider(global.GetLoggerProvider()),
		))
	}

	logger = &Logger{l}
	return logger
}

func ReinitializeLogger() {
	logger = nil
	NewLogger()
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

type MultiHandler struct {
	handlers []slog.Handler
}

func (h *MultiHandler) Enabled(ctx context.Context, level slog.Level) bool {
	for _, handler := range h.handlers {
		if handler.Enabled(ctx, level) {
			return true
		}
	}
	return false
}

func (h *MultiHandler) Handle(ctx context.Context, r slog.Record) error {
	for _, handler := range h.handlers {
		if handler.Enabled(ctx, r.Level) {
			if err := handler.Handle(ctx, r.Clone()); err != nil {
				continue
			}
		}
	}
	return nil
}

func (h *MultiHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	newHandlers := make([]slog.Handler, len(h.handlers))
	for i, handler := range h.handlers {
		newHandlers[i] = handler.WithAttrs(attrs)
	}
	return &MultiHandler{handlers: newHandlers}
}

func (h *MultiHandler) WithGroup(name string) slog.Handler {
	newHandlers := make([]slog.Handler, len(h.handlers))
	for i, handler := range h.handlers {
		newHandlers[i] = handler.WithGroup(name)
	}
	return &MultiHandler{handlers: newHandlers}
}
