# INFRA-004: Python Quality Gates Breaking Go Plugin CI/CD

> **Severity:** P0 (CRITICAL)
> **Status:** OPEN
> **Type:** Configuration/CI-CD
> **Created:** 2026-02-06
> **Estimated Fix:** 15 minutes

## Problem

Python quality gate workflows are running on Go plugin repository, causing **all CI/CD checks to fail**.

### Error Message

```
Unable to locate executable file: poetry
```

### Root Cause Analysis

**Workflow Mismatch:**
- Repository: `fall-out-bug/sdp` (main Python SDP)
- Subdirectory: `sdp-plugin/` (Go-only plugin)
- Workflows running: Python-specific quality gates
- Expected: Go-only workflows

**Failing Workflows:**
1. `.github/workflows/ci-critical.yml` - "Critical Quality Gates"
2. `.github/workflows/ci-warnings.yml` - "Quality Warnings"
3. `.github/workflows/sdp-quality-gate.yml` - "SDP Quality Gate"

**What These Workflows Expect:**
- `pyproject.toml` (Python project config)
- `poetry.lock` (Python dependency lock file)
- Poetry package manager installed
- Python source files (.py)

**What Actually Exists:**
- `go.mod` (Go modules)
- `go.sum` (Go dependencies)
- Go source files (.go)
- **NO Python files**

### Impact

- **ALL PR checks fail** (blocking merges)
- **False negative failures** (code is fine, checks are wrong)
- **Developer confusion** (why is Poetry being checked?)
- **Wasted CI/CD minutes**
- **Reputation damage** (repo looks broken)

### Affected Workflows

**Failing (Python):**
- ❌ Critical Quality Gates
- ❌ Quality Warnings
- ❌ SDP Quality Gate

**Working (Go):**
- ✅ Go CI (go-ci.yml)
- ✅ Go Release (go-release.yml)

## Investigation Details

### Workflow Trigger Analysis

**ci-critical.yml:**
```yaml
on:
  pull_request:
    branches: [ main, dev ]
  push:
    branches: [ main, dev ]
```
**Problem:** Runs on ALL PRs to main/dev, regardless of language

**Expected:** Should only run when Python files change

### Why This Happened

**Historical Context:**
1. Main repo (`fall-out-bug/sdp`) was Python-first
2. Python workflows created at root level
3. Go plugin (`sdp-plugin/`) added as subdirectory
4. Workflows not updated to exclude Go subdirectory
5. Result: Python checks run on Go code

## Solution

### Option 1: Delete Python Workflows (RECOMMENDED)

**Rationale:** Go plugin should be in separate repository or have its own workflows.

**Action:**
```bash
# Remove Python-specific workflows
rm .github/workflows/ci-critical.yml
rm .github/workflows/ci-warnings.yml
rm .github/workflows/sdp-quality-gate.yml
```

**Benefits:**
- ✅ Fixes CI/CD immediately
- ✅ No confusion (Go project = Go workflows)
- ✅ Faster CI/CD (fewer workflows)
- ✅ Cleaner repository

**Downsides:**
- ❌ Loses Python quality gates (if needed elsewhere)
- ❌ Breaking change for Python SDP users

### Option 2: Add Path Filters (COMPROMISE)

**Update workflows to only run on Python files:**

```yaml
on:
  pull_request:
    branches: [ main, dev ]
    paths:
      - '**.py'
      - 'pyproject.toml'
      - 'poetry.lock'
      - '.python-version'
      - 'requirements*.txt'
```

**Benefits:**
- ✅ Keeps Python workflows
- ✅ Doesn't run on Go plugin changes
- ✅ No breaking changes

**Downsides:**
- ❌ Still confusing (mixed Python/Go in same repo)
- ❌ More complex workflow configuration

### Option 3: Split Repositories (LONG-TERM)

**Move Go plugin to separate repository:**
- `fall-out-bug/sdp-python` - Python SDP
- `fall-out-bug/sdp-go-plugin` - Go plugin

**Benefits:**
- ✅ Clean separation of concerns
- ✅ Independent CI/CD
- ✅ No confusion
- ✅ Can use different workflows per repo

**Downsides:**
- ❌ Requires migration
- ❌ Breaking change for users
- ❌ More repos to manage

## Recommended Action Plan

### Immediate (P0 - DO NOW)

1. **Delete Python workflows** (Option 1)
   ```bash
   rm .github/workflows/ci-critical.yml
   rm .github/workflows/ci-warnings.yml
   rm .github/workflows/sdp-quality-gate.yml
   git add .github/workflows/
   git commit -m "fix(infra): remove Python workflows from Go plugin"
   git push
   ```

2. **Verify CI/CD works**
   - Create test PR
   - Confirm only Go workflows run
   - Confirm all checks pass

### Short-term (P1 - THIS WEEK)

3. **Document Go plugin structure**
   - Update README.md
   - Explain workflow structure
   - Add troubleshooting section

### Long-term (P2 - NEXT RELEASE)

4. **Consider repository split** (Option 3)
   - Discuss with team
   - Plan migration if needed
   - Update documentation

## Verification

- [ ] Python workflows deleted
- [ ] Test PR created
- [ ] Only Go workflows run (go-ci.yml)
- [ ] All checks pass
- [ ] No Poetry errors
- [ ] CI/CD time reduced

## Timeline

- **2026-02-06 18:10:** Issue detected (CI failures)
- **2026-02-06 20:XX:** Root cause identified (Python workflows on Go code)
- **Pending:** Workflow deletion
- **Pending:** CI/CD verification

## Related Issues

- INFRA-001: Symbolic Link Loop (FIXED)
- INFRA-002: GitHub Actions Permissions (FIXED)
- INFRA-003: Mixed Python/Go Workflows (DOCUMENTED)

## References

- [GitHub Actions: Path Filters](https://docs.github.com/en/actions/using-workflows/workflow-syntax-for-github-actions#onpushpull_requestpaths)
- [GitHub Actions: Using Path Filters](https://docs.github.com/en/actions/using-workflows triggers-workflow-syntax-for-github-actions#onpushpull_requestpaths-ignore)

---

**Status:** 🔴 OPEN - Awaiting workflow deletion
