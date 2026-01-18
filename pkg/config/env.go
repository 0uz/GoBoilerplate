package config

import (
	"os"
	"strings"
)

const (
	EnvDevelopment = "development"
	EnvProduction  = "production"
)

// NormalizeEnv maps common variants (e.g., DEV/PROD) to canonical values.
func NormalizeEnv(env string) string {
	e := strings.TrimSpace(strings.ToLower(env))
	switch e {
	case "", "dev", "development":
		return EnvDevelopment
	case "prod", "production":
		return EnvProduction
	default:
		return EnvDevelopment
	}
}

func IsDev(env string) bool {
	return NormalizeEnv(env) == EnvDevelopment
}

func IsProd(env string) bool {
	return NormalizeEnv(env) == EnvProduction
}

// GetAppVersion returns the application version from APP_VERSION environment variable.
// If not set, returns "0.0.0-dev" as default.
func GetAppVersion() string {
	if v := os.Getenv("APP_VERSION"); v != "" {
		return v
	}
	return "0.0.0-dev"
}
