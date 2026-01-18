package observability

import (
	"context"

	redisCache "github.com/ouz/goauthboilerplate/pkg/cache/redis"
	sharedconfig "github.com/ouz/goauthboilerplate/pkg/config"
	"github.com/ouz/goauthboilerplate/pkg/log"
	sharedotel "github.com/ouz/goauthboilerplate/pkg/otel"
	"github.com/ouz/goauthboilerplate/internal/adapters/repo/postgres"
	"github.com/ouz/goauthboilerplate/internal/config"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func InitTelemetry(ctx context.Context) (func(context.Context) error, error) {
	cfg := config.Get()
	otelConfig := sharedotel.OtelConfig{
		ServiceName:       cfg.Otel.ServiceName,
		ServiceVersion:    sharedconfig.GetAppVersion(),
		ServiceNamespace:  "production",
		ExporterEndpoint:  cfg.Otel.ExporterEndpoint,
		MonitoringEnabled: cfg.Otel.MonitoringEnabled,
	}
	return sharedotel.SetupOTelSDK(ctx, otelConfig)
}

func InitLogger() *log.Logger {
	cfg := config.Get()
	stdoutLevel := log.ParseLogLevel(cfg.App.LogLevel)
	otelLevel := log.ParseLogLevel(cfg.App.LogLevel)
	return log.NewLogger("api-service", stdoutLevel, otelLevel)
}

func InitDatabase(logger *log.Logger) (*gorm.DB, error) {
	return postgres.ConnectDB(logger)
}

func InitRedis(logger *log.Logger) (*redis.Client, error) {
	cfg := config.Get()
	return redisCache.ConnectRedis(logger, cfg.Valkey.Host, cfg.Valkey.Port, cfg.Otel.MonitoringEnabled)
}
