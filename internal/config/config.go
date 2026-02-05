package config

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	AppName         string
	Environment     string
	Port            int
	LogLevelStr     string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	IdleTimeout     time.Duration
	DatabaseURL     string
	MigrationsPath  string
	WorkerPoolSize  int
	WorkerQueueSize int
}

func Load() (*Config, error) {
	port, err := strconv.Atoi(getEnv("PORT", "8080"))
	if err != nil {
		return nil, fmt.Errorf("invalid PORT: %w", err)
	}

	cfg := &Config{
		AppName:         getEnv("APP_NAME", "payments-ledger"),
		Environment:     getEnv("ENVIRONMENT", "development"),
		Port:            port,
		LogLevelStr:     getEnv("LOG_LEVEL", "info"),
		ReadTimeout:     parseDuration("READ_TIMEOUT", "5s"),
		WriteTimeout:    parseDuration("WRITE_TIMEOUT", "10s"),
		IdleTimeout:     parseDuration("IDLE_TIMEOUT", "120s"),
		DatabaseURL:     getEnv("DATABASE_URL", ""),
		MigrationsPath:  getEnv("MIGRATIONS_PATH", "migrations"),
		WorkerPoolSize:  parseInt("WORKER_POOL_SIZE", 10),
		WorkerQueueSize: parseInt("WORKER_QUEUE_SIZE", 100),
	}

	return cfg, nil
}

func (c *Config) LogLevel() slog.Level {
	switch strings.ToLower(c.LogLevelStr) {
	case "debug":
		return slog.LevelDebug
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func parseInt(envKey string, fallback int) int {
	raw := getEnv(envKey, "")
	if raw == "" {
		return fallback
	}
	v, err := strconv.Atoi(raw)
	if err != nil {
		return fallback
	}
	return v
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}

func parseDuration(envKey, fallback string) time.Duration {
	raw := getEnv(envKey, fallback)
	d, err := time.ParseDuration(raw)
	if err != nil {
		d, _ = time.ParseDuration(fallback)
	}
	return d
}
