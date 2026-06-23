# Runbook: Latency Spike

## Signal

Alert: `ReleaseSentinelLatencyRegression`

Average `/work` latency is above 250 ms for at least 2 minutes.

## First checks

```bash
kubectl top pods
kubectl describe hpa release-sentinel
kubectl logs -l app.kubernetes.io/name=release-sentinel --tail=100
```

Prometheus query:

```promql
sum(rate(release_sentinel_http_request_duration_seconds_sum{path="/work"}[2m]))
/
sum(rate(release_sentinel_http_requests_total{path="/work"}[2m]))
```

## Decision

- If only canary pods are slow, abort the rollout.
- If all pods are slow and CPU is high, check HPA scaling and resource limits.
- If latency is injected intentionally through `LATENCY_MS`, revert the config.

## Mitigation

```bash
kubectl argo rollouts abort release-sentinel
kubectl rollout restart deployment/release-sentinel
```

Use the restart only when the latency is not tied to a canary and the process state is suspected.
