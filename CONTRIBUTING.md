# Contributing

This repository is structured as a DevOps/SRE project, so changes should preserve the release-safety story: application behavior, deployment automation, observability, and operational documentation should stay aligned.

## Development checks

Run the full local validation before opening a pull request:

```bash
make validate
```

For runtime changes, also run:

```bash
make run
make smoke
```

For release-behavior changes, render the Helm chart:

```bash
make helm-template
```

## Change guidelines

- Keep the Go service small and focused on producing realistic release signals.
- Keep Kubernetes manifests explicit about security context, probes, and resources.
- Update runbooks when alert behavior or rollback behavior changes.
- Do not add credentials, tokens, kubeconfigs, or cloud account details.
- Prefer reproducible local commands over environment-specific manual steps.

## Pull request checklist

- Tests and validation pass.
- README or docs are updated when behavior changes.
- Security posture is not weakened without a clear reason.
- New operational failure modes have a runbook or demo note.
