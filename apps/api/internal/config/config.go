package config

import (
	"errors"
	"fmt"
	"log/slog"
	"math"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	HTTPAddr        string
	Environment     string
	Version         string
	CommitSHA       string
	ErrorRate       float64
	Latency         time.Duration
	ShutdownTimeout time.Duration
	LogLevel        slog.Level
}

func Load() (Config, error) {
	cfg := Config{
		HTTPAddr:        readString("HTTP_ADDR", ":8080"),
		Environment:     readString("ENVIRONMENT", "local"),
		Version:         readString("APP_VERSION", "dev"),
		CommitSHA:       readString("COMMIT_SHA", "unknown"),
		ShutdownTimeout: readDuration("SHUTDOWN_TIMEOUT", 10*time.Second),
		Latency:         readDuration("LATENCY_MS", 0),
		LogLevel:        readLogLevel("LOG_LEVEL", slog.LevelInfo),
	}

	errorRate, err := readFloat("ERROR_RATE", 0)
	if err != nil {
		return Config{}, err
	}
	cfg.ErrorRate = errorRate

	if err := cfg.Validate(); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func (c Config) Validate() error {
	if strings.TrimSpace(c.HTTPAddr) == "" {
		return errors.New("HTTP_ADDR must not be empty")
	}
	if c.ErrorRate < 0 || c.ErrorRate > 1 {
		return fmt.Errorf("ERROR_RATE must be between 0 and 1, got %.3f", c.ErrorRate)
	}
	if c.Latency < 0 {
		return fmt.Errorf("LATENCY_MS must not be negative, got %s", c.Latency)
	}
	if c.ShutdownTimeout <= 0 {
		return fmt.Errorf("SHUTDOWN_TIMEOUT must be positive, got %s", c.ShutdownTimeout)
	}
	return nil
}

func readString(key, fallback string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	return value
}

func readFloat(key string, fallback float64) (float64, error) {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback, nil
	}

	parsed, err := strconv.ParseFloat(value, 64)
	if err != nil || math.IsNaN(parsed) || math.IsInf(parsed, 0) {
		return 0, fmt.Errorf("%s must be a finite float", key)
	}

	return parsed, nil
}

func readDuration(key string, fallback time.Duration) time.Duration {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}

	if strings.HasSuffix(value, "ms") || strings.HasSuffix(value, "s") || strings.HasSuffix(value, "m") {
		parsed, err := time.ParseDuration(value)
		if err == nil {
			return parsed
		}
	}

	milliseconds, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}

	return time.Duration(milliseconds) * time.Millisecond
}

func readLogLevel(key string, fallback slog.Level) slog.Level {
	switch strings.ToLower(strings.TrimSpace(os.Getenv(key))) {
	case "debug":
		return slog.LevelDebug
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	case "info", "":
		return slog.LevelInfo
	default:
		return fallback
	}
}
