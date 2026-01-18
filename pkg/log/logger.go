package log

import (
	"context"
	"log/slog"
	"os"
	"strings"

	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/otel/log/global"
)

func ParseLogLevel(level string) slog.Level {
	switch strings.ToUpper(strings.TrimSpace(level)) {
	case "DEBUG":
		return slog.LevelDebug
	case "INFO":
		return slog.LevelInfo
	case "WARN", "WARNING":
		return slog.LevelWarn
	case "ERROR":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

type Logger struct {
	*slog.Logger
}

func NewLogger(serviceName string, stdoutLevel, otelLevel slog.Level) *Logger {
	stdoutHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:     stdoutLevel,
		AddSource: true,
	})

	loggerProvider := global.GetLoggerProvider()
	if loggerProvider == nil {
		return &Logger{slog.New(stdoutHandler)}
	}

	baseOtelHandler := otelslog.NewHandler(serviceName,
		otelslog.WithLoggerProvider(loggerProvider),
		otelslog.WithSource(true),
	)

	otelHandler := &LeveledHandler{
		handler: baseOtelHandler,
		level:   otelLevel,
	}

	multiHandler := &MultiHandler{
		handlers: []slog.Handler{stdoutHandler, otelHandler},
	}

	return &Logger{slog.New(multiHandler)}
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

type LeveledHandler struct {
	handler slog.Handler
	level   slog.Level
}

func (h *LeveledHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return level >= h.level && h.handler.Enabled(ctx, level)
}

func (h *LeveledHandler) Handle(ctx context.Context, r slog.Record) error {
	if r.Level >= h.level {
		return h.handler.Handle(ctx, r)
	}
	return nil
}

func (h *LeveledHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &LeveledHandler{
		handler: h.handler.WithAttrs(attrs),
		level:   h.level,
	}
}

func (h *LeveledHandler) WithGroup(name string) slog.Handler {
	return &LeveledHandler{
		handler: h.handler.WithGroup(name),
		level:   h.level,
	}
}
