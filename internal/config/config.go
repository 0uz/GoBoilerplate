package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/ouz/goauthboilerplate/pkg/errors"
	"github.com/sirupsen/logrus"
)

type Config struct {
	App        AppConfig
	Postgres   PostgresDatabaseConfig
	Valkey     ValkeyConfig
	JWT        JWTConfig
	Mail       MailConfig
	Cache      CacheConfig
}

type AppConfig struct {
	Port     string
	V1Prefix string
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

func Load(logger *logrus.Logger) error {
	if err := loadEnv(logger); err != nil {
		logger.Error("Loading environment variables", "error", err)
	}

	config, err := parseConfig()
	if err != nil {
		return errors.GenericError("Parsing config", err)
	}

	if err := validate(config); err != nil {
		return fmt.Errorf("validating config: %w", err)
	}

	conf = config
	return nil
}

func Get() *Config {
	return conf
}

func loadEnv(logger *logrus.Logger) error {
	if err := godotenv.Load("./.env"); err != nil {
		logger.Info("No .env file found, using environment variables")
	}
	return nil
}

func parseConfig() (*Config, error) {
	accessExp, err := time.ParseDuration(os.Getenv("JWT_ACCESS_EXPIRATION"))
	if err != nil {
		return nil, errors.ValidationError("invalid JWT_ACCESS_EXPIRATION", err)
	}

	refreshExp, err := time.ParseDuration(os.Getenv("JWT_REFRESH_EXPIRATION"))
	if err != nil {
		return nil, errors.ValidationError("invalid JWT_REFRESH_EXPIRATION", err)
	}

	mailPort, err := strconv.Atoi(os.Getenv("MAIL_PORT"))
	if err != nil {
		return nil, errors.ValidationError("invalid MAIL_PORT", err)
	}

	maxOpenConns, _ := strconv.Atoi(getEnvOrDefault("PG_DB_MAX_OPEN_CONNS", "20"))
	maxIdleConns, _ := strconv.Atoi(getEnvOrDefault("PG_DB_MAX_IDLE_CONNS", "25"))
	connMaxLifetimeMinutes, _ := strconv.Atoi(getEnvOrDefault("PG_DB_CONN_MAX_LIFETIME_MINUTES", "5"))
	closeTimeoutSeconds, _ := strconv.Atoi(getEnvOrDefault("PG_DB_CLOSE_TIMEOUT_SECONDS", "5"))
	cacheSizeMB, _ := strconv.Atoi(getEnvOrDefault("CACHE_SIZE_MB", "100"))
	return &Config{
		App: AppConfig{
			Port:     getEnvOrDefault("PORT", "8080"),
			V1Prefix: getEnvOrDefault("V1_PREFIX", "/api/v1"),
		},
		Postgres: PostgresDatabaseConfig{
			Host:                   os.Getenv("PG_DB_HOST"),
			User:                   os.Getenv("PG_DB_USER"),
			Password:               os.Getenv("PG_DB_PASSWORD"),
			Name:                   os.Getenv("PG_DB_NAME"),
			Port:                   os.Getenv("PG_DB_PORT"),
			TimeZone:               getEnvOrDefault("PG_DB_TIMEZONE", "Europe/Istanbul"),
			MaxOpenConns:           maxOpenConns,
			MaxIdleConns:           maxIdleConns,
			ConnMaxLifetimeMinutes: connMaxLifetimeMinutes,
			CloseTimeoutSeconds:    closeTimeoutSeconds,
		},
		Valkey: ValkeyConfig{
			Host: os.Getenv("VALKEY_HOST"),
			Port: os.Getenv("VALKEY_PORT"),
		},
		JWT: JWTConfig{
			Secret:            os.Getenv("JWT_SECRET"),
			AccessExpiration:  accessExp,
			RefreshExpiration: refreshExp,
		},
		Mail: MailConfig{
			Host:     os.Getenv("MAIL_HOST"),
			Port:     mailPort,
			Username: os.Getenv("MAIL_USERNAME"),
			Password: os.Getenv("MAIL_PASSWORD"),
		},
		Cache: CacheConfig{
			SizeMB: cacheSizeMB,
		},
	}, nil
}

func validate(c *Config) error {
	checks := []struct {
		value string
		name  string
	}{
		{c.App.Port, "PORT"},
		{c.App.V1Prefix, "V1_PREFIX"},
		{c.Postgres.Host, "PG_DB_HOST"},
		{c.Postgres.User, "PG_DB_USER"},
		{c.Postgres.Password, "PG_DB_PASSWORD"},
		{c.Postgres.Name, "PG_DB_NAME"},
		{c.Postgres.Port, "PG_DB_PORT"},
		{c.Valkey.Host, "VALKEY_HOST"},
		{c.Valkey.Port, "VALKEY_PORT"},
		{c.JWT.Secret, "JWT_SECRET"},
		{c.Mail.Host, "MAIL_HOST"},
		{c.Mail.Username, "MAIL_USERNAME"},
		{c.Mail.Password, "MAIL_PASSWORD"},
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

	if len(c.JWT.Secret) < 32 {
		return errors.ValidationError("JWT_SECRET must be at least 32 characters long", nil)
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
