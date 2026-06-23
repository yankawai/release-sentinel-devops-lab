# Security Policy

## Supported versions

The `main` branch is the supported version of this project.

## Reporting a vulnerability

Open a private vulnerability report through GitHub Security Advisories when possible. If that is unavailable, open an issue with enough detail to reproduce the finding but do not include live credentials, tokens, or private infrastructure details.

## Security controls in this repository

- distroless non-root runtime image;
- read-only root filesystem in Kubernetes;
- dropped Linux capabilities;
- disabled privilege escalation;
- resource requests and limits;
- SBOM generation in CI;
- Trivy image scanning for high and critical vulnerabilities;
- GitHub secret scanning and push protection.

## Out of scope

This project intentionally avoids cloud provider credentials and production cluster access. Cloud-specific hardening, IAM, private networking, and runtime admission control should be added in a real deployment.
