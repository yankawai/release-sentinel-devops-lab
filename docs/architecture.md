# Architecture

Release Sentinel is intentionally small at the application layer and opinionated at the delivery layer. The service exists to produce realistic release signals: success rate, latency, health, readiness, and structured logs. The platform layer decides whether a release is healthy enough to continue.

## Runtime contract

The API exposes five stable endpoints:

- `/healthz` reports process liveness.
- `/readyz` reports traffic readiness.
- `/version` reports deployed version metadata.
- `/work` simulates business work and emits success/failure metrics.
- `/metrics` exposes Prometheus metrics.

Failure injection is controlled through environment variables:

- `ERROR_RATE` accepts a float between `0` and `1`.
- `LATENCY_MS` accepts milliseconds or a Go duration string.

This keeps the bad-release scenario deterministic without adding external dependencies.

## Release control

The Helm chart deploys the stable service with ordinary Kubernetes primitives: Deployment, Service, HPA, PDB, NetworkPolicy, and optional ServiceMonitor.

The Argo Rollouts manifest represents the promotion path. It shifts traffic in steps and runs a Prometheus-backed analysis between steps. If the measured success rate drops below `98%`, the rollout is aborted before full promotion.

## Observability

Prometheus scrapes `/metrics` and evaluates release SLO alerts. Grafana provisions a release dashboard with error rate, throughput, and average `/work` latency. Loki and OpenTelemetry configs are included so logs and traces have a defined integration path when the lab is connected to a fuller observability stack.

## Security posture

The container runs as non-root on a distroless base image. Kubernetes manifests set read-only root filesystem, drop Linux capabilities, disable privilege escalation, and define resource requests/limits. CI builds the image, generates an SBOM, and blocks high/critical Trivy findings.
