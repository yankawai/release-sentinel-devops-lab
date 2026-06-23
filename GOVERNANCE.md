# Governance

Release Sentinel is maintained as a small open-source project with a bias toward clear operational value.

## Maintainer responsibilities

- Keep `main` releasable.
- Review changes for reliability, security, and maintainability.
- Keep CI, runbooks, and documentation aligned with behavior.
- Avoid accepting features that make the project harder to reproduce locally.

## Decision criteria

Changes should improve at least one of these areas:

- release safety;
- observability;
- supply-chain security;
- incident response;
- local reproducibility;
- Kubernetes production readiness.

## Compatibility

The project is early-stage. Public interfaces are documented through the README, Helm values, endpoint behavior, and runbooks. Breaking changes should be called out in pull requests and release notes once tagged releases begin.

## Security

Security issues are handled through the process described in [SECURITY.md](SECURITY.md).
