# Security Policy

## Reporting Vulnerabilities

If you find a security vulnerability in SDP or its dependencies:

1. **DO NOT** open a public issue
2. Email security advisory to the maintainers (see repository contacts)
3. Include: CVE ID (if known), affected versions, reproduction steps

We will respond within 48 hours with:

- Confirmation of vulnerability
- Severity assessment (CRITICAL/HIGH/MEDIUM/LOW)
- Patch timeline (usually <7 days for CRITICAL/HIGH)

## Supported Versions

| Version | Supported |
|---------|-----------|
| v0.5.x  | ✅ Yes |
| v0.4.x  | ⚠️ Security fixes only |
| < v0.4  | ❌ No |

## Dependency Security

SDP uses **pip-audit** to scan for known vulnerabilities in dependencies:

- Runs automatically on every PR via GitHub Actions
- Blocks merge if vulnerabilities found
- Automated patching via Dependabot (weekly PRs)

### Vulnerability Exceptions

If a vulnerability must be temporarily ignored (e.g., no fix available):

1. Document the exception in this file with:
   - CVE ID and affected package
   - Reason for exception
   - Workaround or mitigation
   - Timeline for resolution
2. Use `poetry run pip-audit --ignore-vuln <VULN_ID>` in CI (temporary only)
