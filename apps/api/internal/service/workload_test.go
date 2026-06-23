package service

import (
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/yankawai/release-sentinel/apps/api/internal/config"
	"github.com/yankawai/release-sentinel/apps/api/internal/metrics"
)

func TestHandleWorkReturnsSuccessWhenErrorRateIsZero(t *testing.T) {
	registry := metrics.NewRegistry()
	workload := NewWorkload(config.Config{ErrorRate: 0}, registry, slog.New(slog.NewTextHandler(io.Discard, nil)))
	request := httptest.NewRequest(http.MethodGet, "/work", nil)
	response := httptest.NewRecorder()

	workload.HandleWork(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("unexpected status: %d", response.Code)
	}
	if !strings.Contains(registry.RenderPrometheus(), `release_sentinel_work_total{result="success"} 1`) {
		t.Fatalf("success metric was not recorded")
	}
}

func TestHandleWorkReturnsFailureWhenErrorRateIsOne(t *testing.T) {
	registry := metrics.NewRegistry()
	workload := NewWorkload(config.Config{ErrorRate: 1}, registry, slog.New(slog.NewTextHandler(io.Discard, nil)))
	request := httptest.NewRequest(http.MethodGet, "/work", nil)
	response := httptest.NewRecorder()

	workload.HandleWork(response, request)

	if response.Code != http.StatusServiceUnavailable {
		t.Fatalf("unexpected status: %d", response.Code)
	}
	if !strings.Contains(registry.RenderPrometheus(), `release_sentinel_work_total{result="failure"} 1`) {
		t.Fatalf("failure metric was not recorded")
	}
}

func TestHandleWorkAppliesConfiguredLatency(t *testing.T) {
	registry := metrics.NewRegistry()
	workload := NewWorkload(config.Config{Latency: 20 * time.Millisecond}, registry, slog.New(slog.NewTextHandler(io.Discard, nil)))
	request := httptest.NewRequest(http.MethodGet, "/work", nil)
	response := httptest.NewRecorder()

	startedAt := time.Now()
	workload.HandleWork(response, request)

	if time.Since(startedAt) < 20*time.Millisecond {
		t.Fatal("configured latency was not applied")
	}
}
