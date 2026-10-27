package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config encapsulates server runtime settings.
type Config struct {
	Host        string
	Port        int
	BasePath    string
	Environment string
}

// Load extracts configuration from environment variables with sensible production defaults.
func Load() (*Config, error) {
	cfg := &Config{
		Host:        getEnv("SERVER_HOST", "0.0.0.0"),
		Port:        8080,
		BasePath:    getEnv("SERVER_BASE_PATH", "/v1"),
		Environment: getEnv("APP_ENV", "development"),
	}

	if portStr := os.Getenv("SERVER_PORT"); portStr != "" {
		p, err := strconv.Atoi(portStr)
		if err != nil {
			return nil, fmt.Errorf("invalid SERVER_PORT value %q: %w", portStr, err)
		}
		cfg.Port = p
	}

	return cfg, nil
}

func getEnv(key, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok && val != "" {
		return val
	}
	return defaultVal
}
