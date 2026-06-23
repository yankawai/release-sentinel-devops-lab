package config

import (
	"testing"
	"time"
)

func TestValidateRejectsInvalidErrorRate(t *testing.T) {
	cfg := Config{
		HTTPAddr:        ":8080",
		ErrorRate:       1.1,
		ShutdownTimeout: time.Second,
	}

	if err := cfg.Validate(); err == nil {
		t.Fatal("expected invalid error rate to be rejected")
	}
}

func TestLoadReadsFailureInjectionSettings(t *testing.T) {
	t.Setenv("ERROR_RATE", "0.25")
	t.Setenv("LATENCY_MS", "150")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("load config: %v", err)
	}

	if cfg.ErrorRate != 0.25 {
		t.Fatalf("unexpected error rate: %f", cfg.ErrorRate)
	}
	if cfg.Latency != 150*time.Millisecond {
		t.Fatalf("unexpected latency: %s", cfg.Latency)
	}
}
