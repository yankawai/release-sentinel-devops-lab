package metrics

import (
	"strings"
	"testing"
	"time"
)

func TestRenderPrometheusIncludesHTTPAndWorkMetrics(t *testing.T) {
	registry := NewRegistry()
	registry.ObserveHTTP("/work", "GET", 200, 50*time.Millisecond)
	registry.ObserveWork(true)
	registry.ObserveWork(false)

	output := registry.RenderPrometheus()

	for _, expected := range []string{
		`release_sentinel_http_requests_total{path="/work",method="GET",status="200"} 1`,
		`release_sentinel_work_total{result="success"} 1`,
		`release_sentinel_work_total{result="failure"} 1`,
	} {
		if !strings.Contains(output, expected) {
			t.Fatalf("missing metric %q in output:\n%s", expected, output)
		}
	}
}
