# INFRA-002: GitHub Actions Token Permissions Missing

> **Severity:** P1 (HIGH)
> **Status:** OPEN
> **Type:** Configuration
> **Created:** 2026-02-06
> **Estimated Fix:** 15 minutes

## Problem

GitHub Actions workflows fail with 403 "Resource not accessible by integration" when trying to post comments on issues/PRs.

### Error Message

```
RequestError [HttpError]: Resource not accessible by integration
  status: 403
  url: 'https://api.github.com/repos/fall-out-bug/sdp/issues/13/comments'
```

### Root Cause

Workflow files missing `permissions` declaration. Default token permissions are read-only, but workflows try to:
- Create issue comments (`issues: write` required)
- Post PR comments (`pull-requests: write` required)

### Impact

- **Workflows fail at "Comment on failure" step**
- **Failed checks don't get comments on PRs**
- **Developers don't get automated feedback**
- **CI appears to fail even when tests pass** (comment step fails after success)

## Affected Workflows

All `.github/workflows/*.yml` files that use `actions/github-script@v7` to post comments.

## Solution

### Add Permissions to Workflow Files

Add to each workflow file (ci.yml, release.yml, etc.):

```yaml
permissions:
  contents: read
  issues: write
  pull-requests: write
```

### Example Fix

**Before:**
```yaml
name: CI
on:
  push:
    branches: [ main, dev ]
jobs:
  build:
    runs-on: ubuntu-latest
```

**After:**
```yaml
name: CI
on:
  push:
    branches: [ main, dev ]
permissions:
  contents: read
  issues: write
  pull-requests: write
jobs:
  build:
    runs-on: ubuntu-latest
```

## Implementation Plan

1. Update `.github/workflows/ci.yml`
2. Update `.github/workflows/release.yml`
3. Test with a new commit
4. Verify comments appear on PRs

## Verification

- [ ] Workflow files updated
- [ ] Test commit pushed
- [ ] PR comments appear successfully
- [ ] No more 403 errors in logs

## Timeline

- **2026-02-06 18:05:** Issue detected in CI/CD logs
- **2026-02-06 20:XX:** Root cause identified
- **Pending:** Fix applied and verified

## Related Issues

- INFRA-001: Symbolic Link Loop (CRITICAL)
- INFRA-003: Mixed Python/Go Workflows

## References

- [GitHub Actions: Managing workflow permissions](https://docs.github.com/en/actions/security-guides/automatic-token-authentication#managing-permissions-for-the-github_token)
- [GitHub Actions: Permissions](https://docs.github.com/en/actions/reference/workflow-syntax-for-github-actions#permissions)

---

**Status:** 🔴 OPEN - Awaiting fix
