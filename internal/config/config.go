package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/ouz/goauthboilerplate/pkg/errors"
)

const (
	envPort                 = "PORT"
	envV1Prefix             = "V1_PREFIX"
	envPGHost               = "PG_DB_HOST"
	envPGUser               = "PG_DB_USER"
	envPGPassword           = "PG_DB_PASSWORD"
	envPGName               = "PG_DB_NAME"
	envPGPort               = "PG_DB_PORT"
	envPGTimeZone           = "PG_DB_TIMEZONE"
	envPGMaxOpenConns       = "PG_DB_MAX_OPEN_CONNS"
	envPGMaxIdleConns       = "PG_DB_MAX_IDLE_CONNS"
	envPGConnMaxLifetimeMin = "PG_DB_CONN_MAX_LIFETIME_MINUTES"
	envPGCloseTimeoutSec    = "PG_DB_CLOSE_TIMEOUT_SECONDS"
	envValkeyHost           = "VALKEY_HOST"
	envValkeyPort           = "VALKEY_PORT"
	envJWTSecret            = "JWT_SECRET"
	envJWTAccessExpiration  = "JWT_ACCESS_EXPIRATION"
	envJWTRefreshExpiration = "JWT_REFRESH_EXPIRATION"
	envMailHost             = "MAIL_HOST"
	envMailPort             = "MAIL_PORT"
	envMailUsername         = "MAIL_USERNAME"
	envMailPassword         = "MAIL_PASSWORD"
	envCacheSizeMB          = "CACHE_SIZE_MB"
	envOtelServiceName      = "OTEL_SERVICE_NAME"
	envOtelExporterEndpoint = "OTEL_EXPORTER_OTLP_ENDPOINT"

	// Default values
	defaultPort                 = "8080"
	defaultV1Prefix             = "/api/v1"
	defaultPGTimeZone           = "Europe/Istanbul"
	defaultPGMaxOpenConns       = 20
	defaultPGMaxIdleConns       = 25
	defaultPGConnMaxLifetimeMin = 5
	defaultPGCloseTimeoutSec    = 5
	defaultCacheSizeMB          = 100
	defaultOtelServiceName      = "go-auth-boilerplate"
	defaultOtelExporterEndpoint = "otel-collector:4317"

	// Validation constants
	minJWTSecretLength = 32
	minCacheSizeMB     = 10
	maxCacheSizeMB     = 1024
	minDBConnections   = 1
	maxDBConnections   = 100
)

type Config struct {
	App      AppConfig
	Postgres PostgresDatabaseConfig
	Valkey   ValkeyConfig
	JWT      JWTConfig
	Mail     MailConfig
	Cache    CacheConfig
	Otel     OtelConfig
}

type AppConfig struct {
	Port     string
	V1Prefix string
}

type OtelConfig struct {
	ServiceName      string
	ExporterEndpoint string
}

type PostgresDatabaseConfig struct {
	Host                   string
	User                   string
	Password               string
	Name                   string
	Port                   string
	TimeZone               string
	MaxOpenConns           int
	MaxIdleConns           int
	ConnMaxLifetimeMinutes int
	CloseTimeoutSeconds    int
}

type ValkeyConfig struct {
	Host string
	Port string
}

type JWTConfig struct {
	Secret            string
	AccessExpiration  time.Duration
	RefreshExpiration time.Duration
}

type MailConfig struct {
	Host     string
	Port     int
	Username string
	Password string
}

type CacheConfig struct {
	SizeMB int
}

var (
	conf *Config
)

func Load(logger *Logger) error {
	if err := loadEnv(logger); err != nil {
		logger.Error("Loading environment variables", "error", err)
	}

	config, err := parseConfig()
	if err != nil {
		return errors.GenericError("Parsing config", err)
	}

	if err := validate(config); err != nil {
		return err
	}

	conf = config
	return nil
}

func Get() *Config {
	return conf
}

func loadEnv(logger *Logger) error {
	envPath := filepath.Join(".", ".env")
	if err := godotenv.Load(envPath); err != nil {
		logger.Info("No .env file found, using environment variables")
	}
	return nil
}

func parseConfig() (*Config, error) {
	accessExp, err := time.ParseDuration(os.Getenv(envJWTAccessExpiration))
	if err != nil {
		return nil, errors.ValidationError("invalid JWT_ACCESS_EXPIRATION", err)
	}

	refreshExp, err := time.ParseDuration(os.Getenv(envJWTRefreshExpiration))
	if err != nil {
		return nil, errors.ValidationError("invalid JWT_REFRESH_EXPIRATION", err)
	}

	mailPort, err := strconv.Atoi(os.Getenv(envMailPort))
	if err != nil {
		return nil, errors.ValidationError("invalid MAIL_PORT", err)
	}

	maxOpenConns, _ := strconv.Atoi(getEnvOrDefault(envPGMaxOpenConns, strconv.Itoa(defaultPGMaxOpenConns)))
	maxIdleConns, _ := strconv.Atoi(getEnvOrDefault(envPGMaxIdleConns, strconv.Itoa(defaultPGMaxIdleConns)))
	connMaxLifetimeMinutes, _ := strconv.Atoi(getEnvOrDefault(envPGConnMaxLifetimeMin, strconv.Itoa(defaultPGConnMaxLifetimeMin)))
	closeTimeoutSeconds, _ := strconv.Atoi(getEnvOrDefault(envPGCloseTimeoutSec, strconv.Itoa(defaultPGCloseTimeoutSec)))
	cacheSizeMB, _ := strconv.Atoi(getEnvOrDefault(envCacheSizeMB, strconv.Itoa(defaultCacheSizeMB)))

	return &Config{
		App: AppConfig{
			Port:     getEnvOrDefault(envPort, defaultPort),
			V1Prefix: getEnvOrDefault(envV1Prefix, defaultV1Prefix),
		},
		Postgres: PostgresDatabaseConfig{
			Host:                   os.Getenv(envPGHost),
			User:                   os.Getenv(envPGUser),
			Password:               os.Getenv(envPGPassword),
			Name:                   os.Getenv(envPGName),
			Port:                   os.Getenv(envPGPort),
			TimeZone:               getEnvOrDefault(envPGTimeZone, defaultPGTimeZone),
			MaxOpenConns:           maxOpenConns,
			MaxIdleConns:           maxIdleConns,
			ConnMaxLifetimeMinutes: connMaxLifetimeMinutes,
			CloseTimeoutSeconds:    closeTimeoutSeconds,
		},
		Valkey: ValkeyConfig{
			Host: os.Getenv(envValkeyHost),
			Port: os.Getenv(envValkeyPort),
		},
		JWT: JWTConfig{
			Secret:            os.Getenv(envJWTSecret),
			AccessExpiration:  accessExp,
			RefreshExpiration: refreshExp,
		},
		Mail: MailConfig{
			Host:     os.Getenv(envMailHost),
			Port:     mailPort,
			Username: os.Getenv(envMailUsername),
			Password: os.Getenv(envMailPassword),
		},
		Cache: CacheConfig{
			SizeMB: cacheSizeMB,
		},
		Otel: OtelConfig{
			ServiceName:      getEnvOrDefault(envOtelServiceName, defaultOtelServiceName),
			ExporterEndpoint: getEnvOrDefault(envOtelExporterEndpoint, defaultOtelExporterEndpoint),
		},
	}, nil
}

func validate(c *Config) error {
	checks := []struct {
		value string
		name  string
	}{
		{c.App.Port, envPort},
		{c.App.V1Prefix, envV1Prefix},
		{c.Postgres.Host, envPGHost},
		{c.Postgres.User, envPGUser},
		{c.Postgres.Password, envPGPassword},
		{c.Postgres.Name, envPGName},
		{c.Postgres.Port, envPGPort},
		{c.Valkey.Host, envValkeyHost},
		{c.Valkey.Port, envValkeyPort},
		{c.JWT.Secret, envJWTSecret},
		{c.Mail.Host, envMailHost},
		{c.Mail.Username, envMailUsername},
		{c.Mail.Password, envMailPassword},
		{c.Otel.ServiceName, envOtelServiceName},
		{c.Otel.ExporterEndpoint, envOtelExporterEndpoint},
	}

	for _, check := range checks {
		if check.value == "" {
			return errors.ValidationError(fmt.Sprintf("%s is not set", check.name), nil)
		}
	}

	if c.Mail.Port == 0 {
		return errors.ValidationError("MAIL_PORT is not set", nil)
	}

	if _, err := strconv.Atoi(c.App.Port); err != nil {
		return errors.ValidationError("invalid PORT", err)
	}

	if _, err := strconv.Atoi(c.Postgres.Port); err != nil {
		return errors.ValidationError("invalid PG_DB_PORT", err)
	}

	if _, err := strconv.Atoi(c.Valkey.Port); err != nil {
		return errors.ValidationError("invalid VALKEY_PORT", err)
	}

	if len(c.JWT.Secret) < minJWTSecretLength {
		return errors.ValidationError(fmt.Sprintf("JWT_SECRET must be at least %d characters long", minJWTSecretLength), nil)
	}

	if c.Cache.SizeMB < minCacheSizeMB || c.Cache.SizeMB > maxCacheSizeMB {
		return errors.ValidationError(fmt.Sprintf("CACHE_SIZE_MB must be between %d and %d", minCacheSizeMB, maxCacheSizeMB), nil)
	}

	if c.Postgres.MaxOpenConns < minDBConnections || c.Postgres.MaxOpenConns > maxDBConnections {
		return errors.ValidationError(fmt.Sprintf("PG_DB_MAX_OPEN_CONNS must be between %d and %d", minDBConnections, maxDBConnections), nil)
	}

	if c.Postgres.MaxIdleConns < minDBConnections || c.Postgres.MaxIdleConns > c.Postgres.MaxOpenConns {
		return errors.ValidationError(fmt.Sprintf("PG_DB_MAX_IDLE_CONNS must be between %d and %d", minDBConnections, c.Postgres.MaxOpenConns), nil)
	}

	return nil
}

func getEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
