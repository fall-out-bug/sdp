# INFRA-003: Mixed Python/Go Workflows

> **Severity:** P2 (MEDIUM)
> **Status:** OPEN
> **Type:** Cleanup/Optimization
> **Created:** 2026-02-06
> **Estimated Fix:** 30 minutes

## Problem

GitHub Actions workflows contain Python-specific setup steps for a Go-only project, wasting CI time and causing confusion.

### Observation

Workflows run `actions/setup-python@v5` with Python 3.10 and Poetry cache, but this is a pure Go project.

### Example Wasted Time

```yaml
- name: Setup Python
  uses: actions/setup-python@v5
  with:
    python-version: 3.10
    cache: poetry  # ❌ Not used in Go project
```

### Root Cause

CI workflows copied from Python SDP repository without adaptation for Go plugin.

### Impact

- **Slower CI/CD** (unnecessary Python setup adds ~30-60 seconds)
- **Confusing logs** (developers wonder why Python is being set up)
- **Maintenance burden** (Python steps must be kept up-to-date for no reason)
- **Misleading documentation** (workflow files don't match actual tech stack)

## Solution

### Remove Python-Specific Steps

**Remove from workflows:**
1. `actions/setup-python@v5` steps
2. Poetry cache configuration
3. Python quality checks (mypy, ruff, pytest)

**Keep/Update:**
1. Go setup steps (already present)
2. Go test commands
3. `golangci-lint` for linting
4. `govulncheck` for security

### Clean Workflow Structure

**Ideal Go CI workflow:**
```yaml
name: CI
on:
  push:
    branches: [ main, dev ]
  pull_request:
    branches: [ main, dev ]

permissions:
  contents: read
  issues: write
  pull-commands: write

jobs:
  build:
    name: Build
    runs-on: ubuntu-latest
    strategy:
      matrix:
        go: ['1.21']
    steps:
      - uses: actions/checkout@v4

      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: ${{ matrix.go }}
          cache: true

      - name: Download dependencies
        run: go mod download

      - name: Run go vet
        run: go vet ./...

      - name: Run tests
        run: go test -v -race -coverprofile=coverage.out ./...

      - name: Check coverage
        run: |
          coverage=$(go tool cover -func=coverage.out | grep total | awk '{print $3}' | sed 's/%//')
          echo "Coverage: $coverage%"
          if (( $(echo "$coverage < 75" | bc -l) )); then
            echo "Coverage $coverage% is below 75% threshold"
            exit 1
          fi

  lint:
    name: Lint
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.21'
          cache: true

      - name: Run golangci-lint
        uses: golangci/golangci-lint-action@v6
        with:
          version: latest
          args: --timeout=5m

  security:
    name: Security
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.21'
          cache: true

      - name: Run govulncheck
        run: |
          go install golang.org/x/vuln/cmd/govulncheck@latest
          govulncheck ./...
```

## Affected Files

- `.github/workflows/ci.yml`
- `.github/workflows/release.yml`
- Any other `.yml` files with Python setup

## Implementation Plan

1. Audit all workflow files for Python-specific steps
2. Remove unnecessary Python setup
3. Verify Go-specific steps are complete
4. Test with new commit
5. Verify CI runs faster

## Benefits

- **Faster CI/CD** (30-60 seconds saved per run)
- **Cleaner logs** (only relevant output)
- **Easier maintenance** (no Python updates to track)
- **Clear documentation** (workflows match actual tech stack)

## Verification

- [ ] Python steps removed from all workflows
- [ ] Go-only workflows tested
- [ ] CI time reduced
- [ ] Logs are cleaner
- [ ] All checks still pass

## Timeline

- **2026-02-06 20:XX:** Issue identified
- **Pending:** Workflow cleanup
- **Pending:** Verification

## Related Issues

- INFRA-001: Symbolic Link Loop (FIXED)
- INFRA-002: GitHub Actions Permissions (OPEN)

## References

- [Go Setup Action](https://github.com/actions/setup-go)
- [golangci-lint Action](https://github.com/golangci/golangci-lint-action)

---

**Status:** 🔵 OPEN - Can be deferred until P0/P1 issues resolved
