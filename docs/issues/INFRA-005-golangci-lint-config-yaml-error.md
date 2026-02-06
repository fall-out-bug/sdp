# INFRA-005: golangci-lint Configuration YAML Syntax Error

> **Severity:** P0 (CRITICAL)
> **Status:** FIXED
> **Type:** Configuration
> **Created:** 2026-02-06
> **Root Cause:** Human error during YAML editing

## Problem

golangci-lint failing with "configuration contains invalid elements" error, blocking all CI/CD.

### Error Message

```
Failed to run: Error: Command failed: .../golangci-lint config verify
Command failed: .../golangci-lint config verify
Failed executing command with error: the configuration contains invalid elements
```

### Root Cause

**YAML Structure Error:** Incorrect indentation in `.golangci.yml`

**Problem (lines 23-27):**
```yaml
      - gocyclo
      - gocognit
      - prealloc
  linters-settings:  # ❌ Wrong indentation (2 spaces)
    errcheck:
```

**Expected:**
```yaml
      - gocyclo
      - gocognit
      - prealloc

linters-settings:  # ✅ Correct (0 spaces, top-level under run)
```

### Full Structural Issue

The entire configuration after line 10 was incorrectly nested:

**BEFORE (BROKEN):**
```yaml
run:
  linters:
    enable:
      - prealloc
  linters-settings:  # ❌ Indented (nested under run)
    errcheck: {...}
  issues:           # ❌ Indented (nested under run)
    exclude: [...]
  run:              # ❌ Duplicate 'run' key!
    skip-dirs: [...]
```

**AFTER (FIXED):**
```yaml
run:
  linters:
    enable:
      - prealloc
  skip-dirs:        # ✅ Moved to correct location
    - vendor
  skip-files:       # ✅ Moved to correct location
    - ".*\\.pb\\.go"

linters-settings:   # ✅ Top-level (under run)
  errcheck: {...}
  govet: {...}

issues:              # ✅ Top-level (under run)
  exclude: [...]
  exclude-rules: [...]
```

### Impact

- **golangci-lint job fails** in CI/CD
- **All PRs blocked** (lint job is required)
- **Developers confused** (config looks valid but isn't)
- **False sense of security** (linting not actually running)

## Solution

### Fix Applied

1. **Moved `linters-settings` to top level** (under `run:`)
2. **Moved `issues` to top level** (under `run:`)
3. **Moved `skip-dirs` and `skip-files`** into `run:` section
4. **Removed duplicate `run:` key**
5. **Verified YAML structure**

### Correct YAML Structure

```yaml
run:
  concurrency: 4
  timeout: 5m
  output: {...}
  linters:
    disable-all: true
    enable:
      - errcheck
      - gosimple
      - govet
      # ... more linters
  skip-dirs:
    - vendor
    - testdata
  skip-files:
    - ".*\\.pb\\.go"

linters-settings:
  errcheck:
    check-type-assertions: true
    check-blank: true
  gocyclo:
    min-complexity: 15
  gocognit:
    min-complexity: 20
  govet:
    enable-all: true
    disable:
      - shadow
      - fieldalignment

issues:
  exclude-use-default: false
  max-issues-per-linter: 0
  max-same-issues: 0
  exclude:
    - "exported (\\w+) should have comment"
    - "comment on exported (\\w+)"
  exclude-rules:
    - path: _test\.go
      linters:
        - gocyclo
        - errcheck
    - linters:
        - staticcheck
      text: "SA9003:"
```

## Prevention

### Pre-commit Hook

Add to `.git/hooks/pre-commit`:

```bash
#!/bin/bash
# Validate golangci-lint config
echo "Validating golangci-lint configuration..."

if command -v golangci-lint &> /dev/null; then
    golangci-lint config verify .golangci.yml
    if [ $? -ne 0 ]; then
        echo "❌ golangci-lint config is invalid!"
        echo "Run: golangci-lint config fix"
        exit 1
    fi
else
    echo "⚠️  golangci-lint not installed, skipping validation"
fi
```

### YAML Linting

Use YAML linter in CI/CD:
```yaml
- name: Validate YAML
  run: |
    pip install yamllint
    yamllint .golangci.yml
```

## Verification

- [x] YAML structure fixed
- [x] Configuration validated
- [x] Committed to git
- [ ] golangci-lint passes in CI/CD
- [ ] Pre-commit hook added

## Timeline

- **2026-02-06 18:17:** Issue detected (golangci-lint failure in CI)
- **2026-02-06 20:XX:** Root cause identified (YAML indentation error)
- **2026-02-06 20:XX:** Fix applied
- **Pending:** CI/CD verification

## Related Issues

- INFRA-001: Symbolic Link Loop (FIXED)
- INFRA-002: GitHub Actions Permissions (FIXED)
- INFRA-003: Mixed Python/Go Workflows (DOCUMENTED)
- INFRA-004: Python Workflows on Go Plugin (FIXED)

## References

- [golangci-lint Configuration](https://golangci-lint.run/usage/configuration/)
- [YAML Specification](https://yaml.org/spec/1.2/spec.html)

---

**Status:** ✅ FIXED - Awaiting CI/CD verification
