package service

import (
	"encoding/json"
	"log/slog"
	"math/rand"
	"net/http"
	"time"

	"github.com/yankawai/release-sentinel-devops-lab/apps/api/internal/config"
	"github.com/yankawai/release-sentinel-devops-lab/apps/api/internal/metrics"
)

type Workload struct {
	cfg      config.Config
	registry *metrics.Registry
	logger   *slog.Logger
	rand     *rand.Rand
}

func NewWorkload(cfg config.Config, registry *metrics.Registry, logger *slog.Logger) *Workload {
	return &Workload{
		cfg:      cfg,
		registry: registry,
		logger:   logger,
		rand:     rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func RegisterHandlers(mux *http.ServeMux, workload *Workload, cfg config.Config, registry *metrics.Registry) {
	mux.HandleFunc("GET /work", workload.HandleWork)
	mux.HandleFunc("GET /metrics", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/plain; version=0.0.4")
		_, _ = w.Write([]byte(registry.RenderPrometheus()))
	})
	mux.HandleFunc("GET /", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{
			"name":        "release-sentinel-api",
			"version":     cfg.Version,
			"environment": cfg.Environment,
		})
	})
}

func (w *Workload) HandleWork(rw http.ResponseWriter, _ *http.Request) {
	if w.cfg.Latency > 0 {
		time.Sleep(w.cfg.Latency)
	}

	if w.shouldFail() {
		w.registry.ObserveWork(false)
		w.logger.Warn("synthetic work failed", "configured_error_rate", w.cfg.ErrorRate)
		writeJSON(rw, http.StatusServiceUnavailable, map[string]string{
			"status": "degraded",
			"reason": "synthetic failure injection",
		})
		return
	}

	w.registry.ObserveWork(true)
	writeJSON(rw, http.StatusOK, map[string]string{
		"status": "processed",
	})
}

func (w *Workload) shouldFail() bool {
	if w.cfg.ErrorRate <= 0 {
		return false
	}
	if w.cfg.ErrorRate >= 1 {
		return true
	}
	return w.rand.Float64() < w.cfg.ErrorRate
}

func writeJSON(w http.ResponseWriter, statusCode int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(payload)
}
