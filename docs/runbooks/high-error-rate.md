# Runbook: High Error Rate

## Signal

Alert: `ReleaseSentinelHighErrorRate`

The `/work` endpoint is failing more than 2% of requests for at least 2 minutes.

## First checks

```bash
kubectl argo rollouts get rollout release-sentinel
kubectl logs -l app.kubernetes.io/name=release-sentinel --tail=100
kubectl port-forward svc/prometheus-server -n monitoring 9090:80
```

Open Prometheus and check:

```promql
sum(rate(release_sentinel_work_total{result="failure"}[2m]))
/
sum(rate(release_sentinel_work_total[2m]))
```

## Decision

- If a canary is active and the error rate started after the rollout, abort the rollout.
- If no rollout is active, inspect recent config changes and dependency health.
- If traffic is spiking and failures are resource-related, scale within the defined HPA/resource limits.

## Mitigation

```bash
kubectl argo rollouts abort release-sentinel
kubectl argo rollouts undo release-sentinel
```

## Follow-up

- Confirm error rate returns below 2%.
- Link the incident to the failed image tag and deployment change.
- Add a regression test or rollout analysis metric if the failure mode was not covered.
