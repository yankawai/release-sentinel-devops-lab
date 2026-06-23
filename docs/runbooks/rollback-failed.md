# Runbook: Rollback Failed

## Signal

Rollback or rollout abort did not return the service to the expected stable state.

## First checks

```bash
kubectl argo rollouts get rollout release-sentinel
kubectl describe rollout release-sentinel
kubectl get rs,pods,svc -l app.kubernetes.io/name=release-sentinel
```

## Common causes

- stable service selector no longer matches stable ReplicaSet labels;
- bad image tag was promoted before metrics had enough samples;
- Prometheus query returns empty data and analysis cannot make a decision;
- pods are blocked by resource quotas or security policy.

## Mitigation

Pin the last known good image:

```bash
kubectl argo rollouts set image release-sentinel api=ghcr.io/yankawai/release-sentinel-api:0.1.0
kubectl argo rollouts promote release-sentinel --full
```

If Rollouts is unhealthy, temporarily switch traffic to the stable service selector and scale the last known good ReplicaSet.

## Follow-up

- Add a pre-promotion pause if Prometheus had insufficient samples.
- Tighten success conditions if the rollout promoted too early.
- Add a CI check for service selector and pod label consistency.
