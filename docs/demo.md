# Demo

## Local service

```bash
make test
make run
```

In another terminal:

```bash
make smoke
```

## Local degraded release

```bash
ERROR_RATE=0.35 LATENCY_MS=250 make run
```

Then run:

```bash
BASE_URL=http://127.0.0.1:8080 k6 run tests/load/release-validation.js
```

The k6 thresholds should fail because the service is behaving like a bad release.

## Compose observability

```bash
make compose-up
```

Open:

- Prometheus: `http://127.0.0.1:9090`
- Grafana: `http://127.0.0.1:3000`

Grafana anonymous viewer access is enabled for the local demo. The admin password is intentionally local-only and should not be reused anywhere else.

## Kubernetes release control

Render the chart:

```bash
make helm-template
```

Deploy the chart:

```bash
helm upgrade --install release-sentinel deploy/helm/release-sentinel
```

Deploy the Rollout resources after Argo Rollouts and Prometheus are installed:

```bash
kubectl apply -f deploy/rollouts/
kubectl argo rollouts get rollout release-sentinel --watch
```

To simulate a bad canary, update the rollout image or env to use a version with elevated `ERROR_RATE`. The analysis template should abort the rollout when success rate falls below 98%.
