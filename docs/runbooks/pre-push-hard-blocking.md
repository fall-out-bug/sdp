# Pre-Push Hook Hard Blocking Implementation

**Date:** 2026-01-30
**Issue:** P1-05
**Status:** ✅ Complete

## Summary

The pre-push hook has been enhanced to support hard blocking mode instead of warning-only behavior. This prevents bad code from reaching the remote repository.

## Changes Made

### 1. Updated `hooks/pre-push.sh`

**Before:**
- Regression test failures: warning only, push allowed
- Coverage < 80%: warning only, push allowed
- No remediation steps
- No configuration options

**After:**
- Regression test failures: blocks push when `SDP_HARD_PUSH=1`
- Coverage < 80%: blocks push when `SDP_HARD_PUSH=1`
- Clear remediation steps for each failure type
- Two modes: Warning (default) and Hard Blocking

### 2. Enhanced Error Messages

**Regression Test Failure:**
```
❌ Regression tests failed

To fix this issue:
  1. Run: cd tools/hw_checker && poetry run pytest tests/unit/ -m fast -v
  2. Fix failing tests
  3. Commit the fixes
  4. Push again

🚫 PUSH BLOCKED (SDP_HARD_PUSH=1)
To bypass: git push --no-verify
```

**Coverage Failure:**
```
❌ Coverage is below 80% (currently: 65%)

To fix this issue:
  1. Run: cd tools/hw_checker && poetry run pytest --cov=. --cov-report=term-missing
  2. Add tests for uncovered lines
  3. Commit the tests
  4. Push again

🚫 PUSH BLOCKED (SDP_HARD_PUSH=1)
To bypass: git push --no-verify
```

### 3. Environment Variable Control

**SDP_HARD_PUSH** - Controls blocking behavior
- `0` or unset: Warning mode (default, backward compatible)
- `1`: Hard blocking mode (blocks on failures)

### 4. Fixed Installation Script

Updated `hooks/install-hooks.sh`:
- Fixed incorrect path (`sdp/hooks` → `hooks`)
- Updated hook descriptions
- Added documentation about `SDP_HARD_PUSH`

### 5. Updated Documentation

**Files updated:**
- `docs/runbooks/git-hooks-installation.md`
  - Added behavior mode documentation
  - Added remediation steps
  - Updated troubleshooting section

- `docs/workstreams/completed/00-007-02-git-hooks.md`
  - Updated execution report
  - Noted hard blocking capability

### 6. Created Test Scripts

**`tests/test_pre_push_hook.sh`** - Automated verification
- Checks hook installation
- Verifies environment variable support
- Validates hard blocking logic
- Confirms remediation messages
- Tests both regression and coverage checks

**`tests/test_pre_push_manual.sh`** - Manual behavior demonstration
- Shows expected behavior in both modes
- Documents error messages
- Provides usage examples

## Behavior Modes

### Warning Mode (Default)

**When:** `SDP_HARD_PUSH` is unset or `0`

**Behavior:**
- Runs regression tests
- Checks coverage
- Shows warnings for failures
- Allows push to proceed
- Exit code: 0 (always)

**Use case:** Development phase, gradual adoption

**Usage:**
```bash
git push
```

### Hard Blocking Mode

**When:** `SDP_HARD_PUSH=1`

**Behavior:**
- Runs regression tests
- Checks coverage
- Blocks push on failures
- Shows remediation steps
- Exit code: 1 (on failures)

**Use case:** Production enforcement, quality gates

**Usage:**
```bash
# Enable for current session
export SDP_HARD_PUSH=1
git push

# One-time use
SDP_HARD_PUSH=1 git push

# Enable permanently (add to ~/.bashrc or ~/.zshrc)
echo 'export SDP_HARD_PUSH=1' >> ~/.bashrc
source ~/.bashrc
```

## Bypass Options

### Emergency Bypass (Not Recommended)

```bash
git push --no-verify
```

**Use when:**
- Critical hotfix that can't wait for test fixes
- Infrastructure failure preventing test execution
- Explicit team decision to bypass

### Temporary Disable

```bash
# Disable hard blocking temporarily
export SDP_HARD_PUSH=0
git push
export SDP_HARD_PUSH=1  # Re-enable
```

## Acceptance Criteria Status

- ✅ Coverage < 80% blocks push with error (when `SDP_HARD_PUSH=1`)
- ✅ Regression failures block push with error (when `SDP_HARD_PUSH=1`)
- ✅ Clear error messages with remediation steps
- ✅ Optional `SDP_HARD_PUSH=1` flag for warm-up period
- ✅ All verification tests pass
- ✅ Documentation updated

## Testing

### Automated Tests

```bash
# Run verification tests
bash tests/test_pre_push_hook.sh
```

**Expected output:**
```
All tests passed!
✓ Installation check
✓ Environment variable support (SDP_HARD_PUSH)
✓ Hard blocking logic (exit 1 on failures)
✓ Warning mode (default, SDP_HARD_PUSH=0)
✓ Remediation steps for failures
✓ Regression test validation
✓ Coverage validation (≥80%)
✓ Bypass instruction (--no-verify)
```

### Manual Testing

```bash
# Run manual behavior test
bash tests/test_pre_push_manual.sh

# Test warning mode (default)
git push

# Test hard blocking mode
export SDP_HARD_PUSH=1
git push

# Test bypass
git push --no-verify
```

## Migration Path

### Phase 1: Warning Mode (Current)
- Default behavior
- Team gets used to seeing warnings
- Fix critical issues as they appear

### Phase 2: Opt-in Hard Blocking
- Individual developers opt-in: `export SDP_HARD_PUSH=1`
- Team builds confidence in blocking mode

### Phase 3: Team-wide Hard Blocking
- Add to team documentation: `export SDP_HARD_PUSH=1` in setup instructions
- All developers use hard blocking mode

### Phase 4: CI/CD Enforcement
- CI/CD pipeline runs with `SDP_HARD_PUSH=1`
- PR checks enforce quality gates
- Bad code never reaches main branch

## Rollback Plan

If issues arise:

1. **Disable hard blocking:**
   ```bash
   unset SDP_HARD_PUSH
   ```

2. **Disable hook entirely:**
   ```bash
   mv .git/hooks/pre-push .git/hooks/pre-push.disabled
   ```

3. **Revert changes:**
   ```bash
   git checkout HEAD~1 hooks/pre-push.sh
   ```

## Files Changed

### Modified
- `hooks/pre-push.sh` - Main hook implementation
- `hooks/install-hooks.sh` - Fixed path, updated descriptions
- `docs/runbooks/git-hooks-installation.md` - Added behavior documentation
- `docs/workstreams/completed/00-007-02-git-hooks.md` - Updated execution report

### Created
- `tests/test_pre_push_hook.sh` - Automated verification tests
- `tests/test_pre_push_manual.sh` - Manual behavior demonstration
- `docs/runbooks/pre-push-hard-blocking.md` - This document

## Next Steps

1. **Team communication:** Announce new hard blocking capability
2. **Individual opt-in:** Developers start using `SDP_HARD_PUSH=1`
3. **Monitor issues:** Track any problems with false positives
4. **Refine thresholds:** Adjust coverage threshold if needed (currently 80%)
5. **Team-wide adoption:** Add to onboarding documentation
6. **CI/CD integration:** Add to automated PR checks

## References

- Pre-push hook: `/Users/fall_out_bug/projects/vibe_coding/sdp/hooks/pre-push.sh`
- Installation script: `/Users/fall_out_bug/projects/vibe_coding/sdp/hooks/install-hooks.sh`
- Documentation: `/Users/fall_out_bug/projects/vibe_coding/sdp/docs/runbooks/git-hooks-installation.md`
- Test script: `/Users/fall_out_bug/projects/vibe_coding/sdp/tests/test_pre_push_hook.sh`

---

**Result:** Pre-push hook now supports hard blocking mode with clear remediation steps and backward-compatible warning mode.
