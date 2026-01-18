package config

import (
	"fmt"
	"time"

	pkgconfig "github.com/ouz/goauthboilerplate/pkg/config"
	"github.com/ouz/goauthboilerplate/pkg/errors"
)

const (
	// Validation constants
	minJWTSecretLength = 32
	minCacheSizeMB     = 10
	maxCacheSizeMB     = 1024
	minDBConnections   = 1
	maxDBConnections   = 100
)

// Config holds all configuration for the application
type Config struct {
	App      AppConfig      `mapstructure:"app"`
	Postgres PostgresConfig `mapstructure:"postgres"`
	Valkey   ValkeyConfig   `mapstructure:"valkey"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	Mail     MailConfig     `mapstructure:"mail"`
	Cache    CacheConfig    `mapstructure:"cache"`
	Otel     OtelConfig     `mapstructure:"otel"`
}

type AppConfig struct {
	Port        string `mapstructure:"port"`
	V1Prefix    string `mapstructure:"v1Prefix"`
	Environment string `mapstructure:"environment"`
	LogLevel    string `mapstructure:"logLevel"`
}

type PostgresConfig struct {
	Host                   string `mapstructure:"host"`
	User                   string `mapstructure:"user"`
	Password               string `mapstructure:"password"`
	Name                   string `mapstructure:"name"`
	Port                   string `mapstructure:"port"`
	TimeZone               string `mapstructure:"timeZone"`
	MaxOpenConns           int    `mapstructure:"maxOpenConns"`
	MaxIdleConns           int    `mapstructure:"maxIdleConns"`
	ConnMaxLifetimeMinutes int    `mapstructure:"connMaxLifetimeMinutes"`
	CloseTimeoutSeconds    int    `mapstructure:"closeTimeoutSeconds"`
}

type ValkeyConfig struct {
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
}

type JWTConfig struct {
	Secret            string        `mapstructure:"secret"`
	AccessExpiration  time.Duration `mapstructure:"accessExpiration"`
	RefreshExpiration time.Duration `mapstructure:"refreshExpiration"`
}

type MailConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

type CacheConfig struct {
	SizeMB int `mapstructure:"sizeMB"`
}

type OtelConfig struct {
	ServiceName       string `mapstructure:"serviceName"`
	ExporterEndpoint  string `mapstructure:"exporterEndpoint"`
	MonitoringEnabled bool   `mapstructure:"monitoringEnabled"`
}

var conf *Config

// Load loads the configuration from yaml files and environment variables
func Load() error {
	var cfg Config

	if err := pkgconfig.LoadConfig("", &cfg, ""); err != nil {
		return errors.GenericError("loading config from yaml", err)
	}

	if err := validate(&cfg); err != nil {
		return err
	}

	conf = &cfg
	return nil
}

// Get returns the loaded configuration
func Get() *Config {
	return conf
}

func validate(c *Config) error {
	// Required string fields
	checks := []struct {
		value string
		name  string
	}{
		{c.App.Port, "app.port"},
		{c.App.V1Prefix, "app.v1Prefix"},
		{c.Postgres.Host, "postgres.host"},
		{c.Postgres.User, "postgres.user"},
		{c.Postgres.Name, "postgres.name"},
		{c.Postgres.Port, "postgres.port"},
		{c.Valkey.Host, "valkey.host"},
		{c.Valkey.Port, "valkey.port"},
		{c.JWT.Secret, "jwt.secret"},
		{c.Otel.ServiceName, "otel.serviceName"},
		{c.Otel.ExporterEndpoint, "otel.exporterEndpoint"},
	}

	for _, check := range checks {
		if check.value == "" {
			return errors.ValidationError(fmt.Sprintf("%s is not set", check.name), nil)
		}
	}

	// JWT secret length validation
	if len(c.JWT.Secret) < minJWTSecretLength {
		return errors.ValidationError(
			fmt.Sprintf("jwt.secret must be at least %d characters long", minJWTSecretLength),
			nil,
		)
	}

	// JWT expiration validation
	if c.JWT.AccessExpiration <= 0 {
		return errors.ValidationError("jwt.accessExpiration must be greater than 0", nil)
	}

	if c.JWT.RefreshExpiration <= 0 {
		return errors.ValidationError("jwt.refreshExpiration must be greater than 0", nil)
	}

	// Cache size validation
	if c.Cache.SizeMB < minCacheSizeMB || c.Cache.SizeMB > maxCacheSizeMB {
		return errors.ValidationError(
			fmt.Sprintf("cache.sizeMB must be between %d and %d", minCacheSizeMB, maxCacheSizeMB),
			nil,
		)
	}

	// Database connection pool validation
	if c.Postgres.MaxOpenConns < minDBConnections || c.Postgres.MaxOpenConns > maxDBConnections {
		return errors.ValidationError(
			fmt.Sprintf("postgres.maxOpenConns must be between %d and %d", minDBConnections, maxDBConnections),
			nil,
		)
	}

	if c.Postgres.MaxIdleConns < minDBConnections || c.Postgres.MaxIdleConns > c.Postgres.MaxOpenConns {
		return errors.ValidationError(
			fmt.Sprintf("postgres.maxIdleConns must be between %d and %d", minDBConnections, c.Postgres.MaxOpenConns),
			nil,
		)
	}

	return nil
}
