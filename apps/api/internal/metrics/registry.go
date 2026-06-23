package metrics

import (
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"
)

type Registry struct {
	mu       sync.RWMutex
	requests map[httpMetricKey]*httpMetric
	work     workMetric
}

type httpMetricKey struct {
	Path   string
	Method string
	Status int
}

type httpMetric struct {
	Count         uint64
	TotalDuration time.Duration
}

type workMetric struct {
	Success uint64
	Failure uint64
}

func NewRegistry() *Registry {
	return &Registry{requests: make(map[httpMetricKey]*httpMetric)}
}

func (r *Registry) ObserveHTTP(path, method string, statusCode int, duration time.Duration) {
	r.mu.Lock()
	defer r.mu.Unlock()

	key := httpMetricKey{Path: normalizePath(path), Method: method, Status: statusCode}
	metric, ok := r.requests[key]
	if !ok {
		metric = &httpMetric{}
		r.requests[key] = metric
	}
	metric.Count++
	metric.TotalDuration += duration
}

func (r *Registry) ObserveWork(success bool) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if success {
		r.work.Success++
		return
	}
	r.work.Failure++
}

func (r *Registry) RenderPrometheus() string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var builder strings.Builder
	builder.WriteString("# HELP release_sentinel_http_requests_total Total HTTP requests.\n")
	builder.WriteString("# TYPE release_sentinel_http_requests_total counter\n")

	keys := make([]httpMetricKey, 0, len(r.requests))
	for key := range r.requests {
		keys = append(keys, key)
	}
	sort.Slice(keys, func(i, j int) bool {
		return fmt.Sprintf("%s-%s-%d", keys[i].Path, keys[i].Method, keys[i].Status) < fmt.Sprintf("%s-%s-%d", keys[j].Path, keys[j].Method, keys[j].Status)
	})

	for _, key := range keys {
		metric := r.requests[key]
		builder.WriteString(fmt.Sprintf("release_sentinel_http_requests_total{path=%q,method=%q,status=%q} %d\n", key.Path, key.Method, fmt.Sprint(key.Status), metric.Count))
	}

	builder.WriteString("# HELP release_sentinel_http_request_duration_seconds_sum Total request duration by route.\n")
	builder.WriteString("# TYPE release_sentinel_http_request_duration_seconds_sum counter\n")
	for _, key := range keys {
		metric := r.requests[key]
		builder.WriteString(fmt.Sprintf("release_sentinel_http_request_duration_seconds_sum{path=%q,method=%q,status=%q} %.6f\n", key.Path, key.Method, fmt.Sprint(key.Status), metric.TotalDuration.Seconds()))
	}

	builder.WriteString("# HELP release_sentinel_work_total Synthetic work outcomes.\n")
	builder.WriteString("# TYPE release_sentinel_work_total counter\n")
	builder.WriteString(fmt.Sprintf("release_sentinel_work_total{result=\"success\"} %d\n", r.work.Success))
	builder.WriteString(fmt.Sprintf("release_sentinel_work_total{result=\"failure\"} %d\n", r.work.Failure))

	return builder.String()
}

func normalizePath(path string) string {
	if path == "" {
		return "/"
	}
	return path
}
