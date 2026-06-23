# Roadmap

Release Sentinel is focused on practical release safety for Kubernetes services. The roadmap favors features that improve deployment confidence, observability, and operational response.

## Current

- Go service with controlled failure and latency injection.
- Helm deployment with probes, resources, security context, HPA, PDB, and NetworkPolicy.
- Argo Rollouts canary analysis backed by Prometheus.
- Grafana dashboard and Prometheus alert rules.
- CI validation with tests, Helm rendering, Docker build, SBOM, and Trivy scanning.
- Operational runbooks for error rate, latency regression, and rollback failure.

## Next

- Cosign image signing and verification examples.
- SLSA provenance for release artifacts.
- Kyverno or OPA policies for deployment guardrails.
- Argo CD GitOps bootstrap manifests.
- Chaos scenarios for pod disruption and dependency latency.
- Example multi-service rollout with dependency-aware analysis.

## Later

- Terraform bootstrap for a disposable Kubernetes environment.
- Multi-cluster promotion flow.
- Release scorecard CLI for summarizing rollout health.
- More Grafana panels for saturation, burn rate, and deployment events.

## Non-goals

- Owning cloud-provider-specific production infrastructure.
- Replacing Argo Rollouts, Prometheus, or Grafana.
- Hiding operational complexity behind magic defaults.
