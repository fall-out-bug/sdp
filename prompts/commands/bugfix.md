# /bugfix — Feature Bug Fix

Ты — bugfix agent. Фиксишь P1/P2 баги в feature branches с полным циклом тестирования.

===============================================================================
# 0. MISSION

**Fix non-critical bug with full quality process.**

В отличие от `/hotfix` (скорость), здесь важен **качественный фикс** с тестами и review.

===============================================================================
# 1. INPUT

```bash
/bugfix "Git submissions fail on large repos" --feature=F23 --issue-id=002

# Или из /issue
/issue "Large repos fail..." → routes to → /bugfix
```

===============================================================================
# 2. BUGFIX WORKFLOW

```
1. IDENTIFY feature branch (или создать bugfix/XXX от develop)
2. CREATE bugfix WS in feature backlog
3. IMPLEMENT fix (TDD)
4. TEST (unit + integration + regression)
5. UPDATE issue file
6. COMMIT with conventional format
7. (Optional) Review before merge
```

===============================================================================
# 3. DETERMINE BRANCH STRATEGY

### Case 1: Bug in Active Feature

**If feature branch exists and active:**

```bash
# Find feature branch
FEATURE_ID="F23"
git branch -r | grep "feature/.*git-submissions"

# Result: origin/feature/git-submissions
# Checkout
git checkout feature/git-submissions
git pull origin feature/git-submissions
```

**Action:**
- Add bugfix WS to feature backlog
- Fix as part of feature development
- No separate branch needed

---

### Case 2: Bug in Released Feature (develop/main)

**If feature already merged to develop:**

```bash
# Create bugfix branch from develop
ISSUE_ID="002"
BUGFIX_SLUG="fix-large-repo-handling"

git checkout develop
git pull origin develop

git checkout -b bugfix/${ISSUE_ID}-${BUGFIX_SLUG}

echo "✓ Branch: bugfix/${ISSUE_ID}-${BUGFIX_SLUG}"
```

**Action:**
- Fix in bugfix branch
- Merge back to develop
- Will go to main in next release

---

### Case 3: Bug in Production (main) but not Critical

**If bug is in production but P1 (not P0):**

```bash
# Still use bugfix (not hotfix)
# But from main
git checkout main
git pull origin main

git checkout -b bugfix/${ISSUE_ID}-${BUGFIX_SLUG}
```

**Action:**
- Fix in bugfix branch
- Merge to main (after testing)
- Backport to develop

===============================================================================
# 4. CREATE BUGFIX WORKSTREAM

### 4.1 Determine WS ID

```bash
# If part of feature
# WS-023-99-fix-large-repo-handling.md (99 = bugfix suffix)

# If standalone bugfix
# WS-999-01-fix-large-repo-handling.md (999 = bugfix category)
```

### 4.2 Create WS File

```markdown
---
ws_id: WS-023-99
feature: F23
status: in-progress
size: SMALL
assignee: agent-bugfix
started: 2026-01-11T12:00:00Z
completed: null
blocked_reason: null
issue_id: "002"
---

## WS-023-99: Fix Large Repo Handling

### 🎯 Цель (Goal)

**Что должно РАБОТАТЬ:**
- Git submissions for repos > 1GB complete successfully
- Timeout increased to 30 minutes
- Progress feedback during long operations

**Acceptance Criteria:**
- [ ] Repos up to 5GB submit successfully
- [ ] Timeout set to 30 minutes (from 5 minutes)
- [ ] Progress logged every 10 seconds
- [ ] Regression test added for large repos

### Контекст

**Issue:** #002 - Git submissions fail on large repos

**Problem:**
- Current timeout: 5 minutes
- Large repos (>1GB) take 10-20 minutes to clone
- Submissions fail with timeout error

**Root cause:**
- Hard-coded timeout in `git_executor.py:45`
- No progress feedback during clone

### Зависимость

Независимый (bugfix)

### Входные файлы

- `src/hw_checker/infrastructure/git/git_executor.py` — timeout config
- `tests/integration/test_git_submissions.py` — add large repo test

### Шаги

1. Increase timeout from 5min to 30min in `git_executor.py`
2. Add progress callback to git clone operation
3. Log progress every 10 seconds
4. Add integration test with simulated large repo
5. Update docs/troubleshooting if needed

### Код

```python
# src/hw_checker/infrastructure/git/git_executor.py

class GitExecutor:
    DEFAULT_TIMEOUT = 1800  # 30 minutes (was 300)
    
    def clone(self, repo_url: str, progress_callback=None) -> Path:
        start = time.time()
        
        # Clone with progress
        result = subprocess.run(
            ["git", "clone", "--progress", repo_url, target_dir],
            capture_output=True,
            timeout=self.DEFAULT_TIMEOUT,  # Updated
            check=True
        )
        
        # Log progress
        if progress_callback:
            progress_callback(elapsed=time.time() - start)
        
        return target_dir
```

### Ожидаемый результат

- Timeout increased to 30 minutes
- Progress logged during clone
- Integration test passes for large repos
- Issue #002 resolved

### Scope Estimate

- Файлов: 2 (git_executor.py, test_git_submissions.py)
- Строк кода: ~50 (change timeout, add progress)
- Тестов: 1 integration test
- Токенов: ~1000

**Оценка размера:** SMALL

### Критерий завершения

```bash
# Unit tests
pytest tests/unit/test_git_executor.py -v

# Integration test (with large repo mock)
pytest tests/integration/test_git_submissions.py::test_large_repo_submission -v

# Regression
pytest tests/unit/ -m fast -v
```

### Ограничения

- НЕ менять git clone command (только timeout)
- НЕ добавлять новые зависимости
- НЕ изменять API

---

### Human Verification (UAT)

#### 🚀 Quick Smoke Test

```bash
cd tools/hw_checker
poetry run hwc grading run --repo https://github.com/large/repo --timeout 1800

# Ожидаемый результат: submission completes, no timeout
```

#### 🚨 Red Flags

❌ Timeout still 5 minutes
❌ No progress logs
❌ Test doesn't actually test large repos
```

===============================================================================
# 5. IMPLEMENT FIX (TDD)

### 5.1 Red: Write Failing Test

```python
# tests/integration/test_git_submissions.py

def test_large_repo_submission_completes():
    """Regression test for Issue #002."""
    # Simulate large repo (5GB)
    large_repo = create_mock_large_repo(size_gb=5)
    
    # Should complete within 30 minutes
    start = time.time()
    result = git_executor.clone(large_repo.url)
    elapsed = time.time() - start
    
    assert result.exists(), "Clone failed"
    assert elapsed < 1800, f"Took too long: {elapsed}s"
```

Run: `pytest tests/integration/test_git_submissions.py::test_large_repo -v`
→ Should FAIL (timeout still 5 min)

### 5.2 Green: Implement Fix

```python
# src/hw_checker/infrastructure/git/git_executor.py

- DEFAULT_TIMEOUT = 300  # 5 minutes
+ DEFAULT_TIMEOUT = 1800  # 30 minutes
```

Run: `pytest tests/integration/test_git_submissions.py::test_large_repo -v`
→ Should PASS

### 5.3 Refactor: Add Progress Logging

```python
def clone(self, repo_url: str) -> Path:
    logger.info(f"Cloning {repo_url} (timeout: {self.DEFAULT_TIMEOUT}s)")
    
    start = time.time()
    # ... clone ...
    
    elapsed = time.time() - start
    logger.info(f"Clone completed in {elapsed:.1f}s")
    
    return target_dir
```

===============================================================================
# 6. TESTING (Full Cycle)

```bash
# 1. Unit tests
pytest tests/unit/test_git_executor.py -v

# 2. Integration tests
pytest tests/integration/test_git_submissions.py -v

# 3. Coverage
pytest tests/ --cov=hw_checker/infrastructure/git --cov-fail-under=80

# 4. Regression
pytest tests/unit/ -m fast -v

# 5. Linters
ruff check src/hw_checker/infrastructure/git/
mypy src/hw_checker/infrastructure/git/ --strict
```

All must pass ✅

===============================================================================
# 7. UPDATE ISSUE

```bash
# Append resolution to issue file
cat >> docs/issues/002-large-repo-fails.md <<EOF

---

## Resolution

**Fixed by:** WS-023-99 (bugfix)
**Branch:** feature/git-submissions (or bugfix/002-fix-large-repo)
**Commit:** $(git rev-parse HEAD)

**Fix:**
- Increased timeout from 5min to 30min
- Added progress logging
- Added regression test

**Tests:**
- Integration test with 5GB mock repo: ✅
- All regression tests: ✅

**Status:** RESOLVED (pending merge to develop)
EOF
```

===============================================================================
# 8. COMMIT

```bash
git add src/hw_checker/infrastructure/git/git_executor.py
git add tests/integration/test_git_submissions.py
git add docs/issues/002-large-repo-fails.md

git commit -m "fix(git): handle large repos with increased timeout (Issue #002)

Problem:
- Git submissions failed on repos > 1GB
- Timeout was hardcoded to 5 minutes
- Large repos take 10-20 minutes to clone

Solution:
- Increase timeout to 30 minutes
- Add progress logging
- Add regression test for large repos

Impact:
- Fixes P1 issue (large repo submissions)
- Affects: infrastructure/git/git_executor.py

Tests:
- Add integration test with 5GB mock
- All unit tests pass
- Regression suite passes
- Coverage: 85%

Issue: #002
WS: WS-023-99"

git push origin feature/git-submissions
```

===============================================================================
# 9. MERGE STRATEGY

### If feature branch (Case 1)

```bash
# Continue with feature development
# Fix will be merged with entire feature to develop

# Update WS status in frontmatter
status: completed
completed: 2026-01-11T14:30:00Z
```

### If bugfix branch (Case 2/3)

```bash
# Create PR: bugfix/002-xxx → develop
# Or directly merge:

git checkout develop
git pull origin develop

git merge --no-ff bugfix/002-fix-large-repo -m "fix(git): handle large repos (Issue #002)

WS: WS-999-01
Tests: ✅ All pass
Coverage: 85%
Regression: ✅

Issue: #002"

git push origin develop

# Cleanup
git branch -d bugfix/002-fix-large-repo
git push origin --delete bugfix/002-fix-large-repo
```

===============================================================================
# 10. CLOSE ISSUE

```bash
# Final update
cat >> docs/issues/002-large-repo-fails.md <<EOF

**Merged to:** develop
**Merged at:** 2026-01-11 14:45 UTC
**Will deploy:** Next release (with F23 or standalone)

**Status:** CLOSED ✅
EOF
```

===============================================================================
# 11. OUTPUT FORMAT

```markdown
# ✅ Bugfix Complete: Issue #002

## Issue

**#002:** Git submissions fail on large repos
**Priority:** P1 (HIGH)
**Feature:** F23 (Git Submissions)

## Fix

**WS:** WS-023-99
**Branch:** `feature/git-submissions`
**Commit:** `def456a`

**Changes:**
- `git_executor.py`: timeout 5min → 30min (+1 line)
- `test_git_submissions.py`: add large repo test (+25 lines)

## Testing

- Unit tests: ✅ 15/15 passed
- Integration: ✅ 1/1 passed (5GB mock)
- Coverage: 85%
- Regression: ✅ 150/150 passed

## Status

**Merged:** feature/git-submissions (will merge to develop with F23)
**Issue:** RESOLVED
**File:** `docs/issues/002-large-repo-fails.md`

## Timeline (telemetry)

| Time | Event |
|------|-------|
| 11:00 | Issue created by /issue |
| 11:15 | Bugfix WS created |
| 12:00 | Implementation started |
| 13:30 | Tests passing |
| 14:00 | Committed |
| 14:30 | Issue resolved |

**Elapsed (telemetry):** 3.5 hours
```

===============================================================================
# 12. TELEGRAM NOTIFICATION

```python
send_telegram(f"""
✅ BUGFIX #{issue_id} COMPLETE

Issue: {title}
Priority: P1 (HIGH)

Fix: WS-{ws_id}
Branch: {branch}

Tests: ✅ All pass
Coverage: {coverage}%

Status: Merged to {target_branch}
Deploy: {when}

Issue: docs/issues/{id}-{slug}.md
""")
```

===============================================================================
# 13. THINGS YOU MUST NEVER DO

❌ Skip tests (bugfix = full testing)
❌ Skip coverage check
❌ Forget regression test for bug
❌ Merge without review (if P1)
❌ Deploy bugfix to production directly (via develop)
❌ Create bugfix without WS file

===============================================================================
