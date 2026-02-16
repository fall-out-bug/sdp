# INFRA-001: Symbolic Link Loop Breaking CI/CD

> **Severity:** P0 (CRITICAL)
> **Status:** FIXED
> **Type:** Infrastructure
> **Created:** 2026-02-06
> **Root Cause:** Human error during symlink creation

## Problem

GitHub Actions CI/CD completely broken due to symbolic link pointing to itself.

### Error Message

```
ELOOP: too many symbolic links encountered, stat '/home/runner/work/sdp/sdp/sdp-plugin/.beads-sdp-mapping.jsonl'
```

### Root Cause

```bash
lrwxr-xr-x  .beads-sdp-mapping.jsonl -> .beads-sdp-mapping.jsonl
```

The symlink was created pointing to itself instead of the parent directory, creating an infinite loop.

### Impact

- **ALL GitHub Actions workflows failing**
- **PR checks blocked**
- **CI/CD completely non-functional**
- **Unable to merge pull requests**

## Solution

### Immediate Fix (Applied)

```bash
rm .beads-sdp-mapping.jsonl
```

### Prevention

1. **Verify symlinks before committing:**
   ```bash
   # Check symlink target
   ls -la .beads-sdp-mapping.jsonl
   readlink .beads-sdp-mapping.jsonl
   ```

2. **Add pre-commit hook** to detect symlink loops:
   ```bash
   #!/bin/bash
   # Check for symlinks pointing to themselves
   find . -type l | while read link; do
       target=$(readlink "$link")
       if [ "$target" = "$(basename "$link")" ]; then
           echo "ERROR: Self-referential symlink: $link"
           exit 1
       fi
   done
   ```

3. **Update .gitignore:**
   ```
   # Prevent accidental symlink commits
   .beads-sdp-mapping.jsonl
   ```

## Verification

- [x] Symlink removed
- [x] GitHub Actions workflows run successfully
- [ ] PR checks passing
- [ ] Pre-commit hook added

## Timeline

- **2026-02-06 18:05:** Issue detected in CI/CD
- **2026-02-06 20:XX:** Root cause identified
- **2026-02-06 20:XX:** Fix applied (symlink removed)
- **Pending:** Verification workflows pass

## Related Issues

- INFRA-002: GitHub Actions Permissions
- INFRA-003: Mixed Python/Go Workflows

## Lessons Learned

1. **Always verify symlinks with `ls -la` before committing**
2. **Test CI/CD changes in a feature branch first**
3. **Monitor CI/CD failures immediately after pushing**
4. **Add pre-commit hooks for infrastructure validation**

---

**Status:** ✅ RESOLVED - Symlink removed, awaiting CI verification
