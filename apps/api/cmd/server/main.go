package main

import (
	"context"
	"errors"
	"flag"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/yankawai/release-sentinel/apps/api/internal/config"
	"github.com/yankawai/release-sentinel/apps/api/internal/health"
	"github.com/yankawai/release-sentinel/apps/api/internal/metrics"
	"github.com/yankawai/release-sentinel/apps/api/internal/service"
)

func main() {
	healthcheck := flag.Bool("healthcheck", false, "check local health endpoint")
	flag.Parse()
	if *healthcheck {
		runHealthcheck()
		return
	}

	cfg, err := config.Load()
	if err != nil {
		slog.Error("configuration error", "error", err)
		os.Exit(1)
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: cfg.LogLevel}))
	registry := metrics.NewRegistry()
	workload := service.NewWorkload(cfg, registry, logger)

	mux := http.NewServeMux()
	health.RegisterHandlers(mux, cfg)
	service.RegisterHandlers(mux, workload, cfg, registry)

	server := &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           loggingMiddleware(logger, metricsMiddleware(registry, mux)),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	errCh := make(chan error, 1)
	go func() {
		logger.Info("api server started", "addr", cfg.HTTPAddr, "version", cfg.Version, "environment", cfg.Environment)
		errCh <- server.ListenAndServe()
	}()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			logger.Error("graceful shutdown failed", "error", err)
			os.Exit(1)
		}
		logger.Info("api server stopped")
	case err := <-errCh:
		if !errors.Is(err, http.ErrServerClosed) {
			logger.Error("api server failed", "error", err)
			os.Exit(1)
		}
	}
}

func runHealthcheck() {
	client := http.Client{Timeout: 2 * time.Second}
	response, err := client.Get("http://127.0.0.1:8080/healthz")
	if err != nil {
		os.Exit(1)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		os.Exit(1)
	}
}

func loggingMiddleware(logger *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startedAt := time.Now()
		recorder := &statusRecorder{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(recorder, r)

		logger.Info("http request",
			"method", r.Method,
			"path", r.URL.Path,
			"status", recorder.statusCode,
			"duration_ms", time.Since(startedAt).Milliseconds(),
			"remote_addr", r.RemoteAddr,
		)
	})
}

func metricsMiddleware(registry *metrics.Registry, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startedAt := time.Now()
		recorder := &statusRecorder{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(recorder, r)

		registry.ObserveHTTP(r.URL.Path, r.Method, recorder.statusCode, time.Since(startedAt))
	})
}

type statusRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (r *statusRecorder) WriteHeader(statusCode int) {
	r.statusCode = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}
