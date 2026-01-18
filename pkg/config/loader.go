package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

func LoadConfig(service string, out any, envPrefix string) error {
	env := getEnvironment()

	v := viper.New()
	v.SetConfigType("yaml")
	if envPrefix != "" {
		v.SetEnvPrefix(envPrefix)
	}
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	// Search paths for config files
	searchPaths := []string{
		"config",      // local development
		"/app/config", // docker fallback
	}

	// Find base.yaml
	var baseFile string
	for _, p := range searchPaths {
		var path string
		if service != "" {
			path = filepath.Join(p, service, "base.yaml")
		} else {
			path = filepath.Join(p, "base.yaml")
		}
		if _, err := os.Stat(path); err == nil {
			baseFile = path
			break
		}
	}

	if baseFile == "" {
		return fmt.Errorf("cannot load base config: base.yaml not found in %v", searchPaths)
	}

	// Load base config
	v.SetConfigFile(baseFile)
	if err := v.MergeInConfig(); err != nil {
		return fmt.Errorf("cannot load base config: %w", err)
	}

	// Load environment-specific config (optional)
	envFile := ""
	for _, p := range searchPaths {
		var path string
		if service != "" {
			path = filepath.Join(p, service, fmt.Sprintf("%s.yaml", env))
		} else {
			path = filepath.Join(p, fmt.Sprintf("%s.yaml", env))
		}
		if _, err := os.Stat(path); err == nil {
			envFile = path
			break
		}
	}

	if envFile != "" {
		v.SetConfigFile(envFile)
		if err := v.MergeInConfig(); err != nil {
			return fmt.Errorf("cannot load %s config: %w", env, err)
		}
	}

	// Unmarshal to struct
	if err := v.Unmarshal(out); err != nil {
		return fmt.Errorf("cannot decode config: %w", err)
	}

	return nil
}

// getEnvironment returns the current environment based on APP_ENV
func getEnvironment() string {
	env := os.Getenv("APP_ENV")
	if env == "" {
		if _, err := os.Stat("/.dockerenv"); err == nil {
			return "development-docker"
		}
		return "development"
	}

	switch strings.ToUpper(env) {
	case "DEV", "DEVELOPMENT":
		return "development"
	case "PROD", "PRODUCTION":
		return "production"
	default:
		return env
	}
}
