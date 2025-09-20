package observability

import (
	"log/slog"

	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/otel/log/global"
)

// SetupOTelSlog sets up slog to export logs via OpenTelemetry
func SetupOTelSlog() *slog.Logger {
	// Get the global log provider
	provider := global.GetLoggerProvider()

	// Create OTEL slog handler
	handler := otelslog.NewHandler("go-boilerplate", otelslog.WithLoggerProvider(provider))

	// Create and return logger
	logger := slog.New(handler)

	// Set as default logger
	slog.SetDefault(logger)

	return logger
}
