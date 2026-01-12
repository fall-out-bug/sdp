# /bugfix — Quality Bug Fixes

You are a bugfix agent. Fix P1/P2 bugs in feature branches with full testing cycle.

===============================================================================
# 0. RECOMMENDED @FILE REFERENCES

**Always start with:**
```
@docs/issues/{issue-id}.md  # Issue analysis
@src/{affected_module}/  # Buggy code
@tests/{module}/  # Existing tests
@PROJECT_CONVENTIONS.md  # Testing conventions
```

**Why:**
- Issue file — Understand problem
- Source code — Locate bug
- Tests — Check coverage
- Conventions — Follow project rules

===============================================================================
# 0. GLOBAL RULES

1. **Quality over speed** — unlike hotfix, full TDD here
2. **Full test coverage** — for the fix
3. **Branch from develop** — not main
4. **Complete review cycle** — same as features
5. **No direct prod deploy** — goes through staging

**Key Difference from Hotfix:**

| Aspect | Hotfix | Bugfix |
|--------|--------|--------|
| Severity | P0 CRITICAL | P1/P2 |
| Branch from | main | develop/feature |
| Testing | Fast (smoke only) | Full TDD |
| Deploy | Production directly | Staging first |
| Timeline | < 2h | < 24h |
| Refactoring | Forbidden | Allowed (minimal) |

===============================================================================
# 1. TRIGGERS

```bash
# Direct call
/bugfix "Large repos fail to clone" --feature=F23 --issue-id=002

# From /issue (auto-routed)
/issue "Large repos fail..." → P1 → /bugfix
```

===============================================================================
# 2. ALGORITHM

```
1. ANALYZE (15-30 min)
   - Reproduce issue
   - Identify root cause
   - Define fix scope

2. BRANCH (2 min)
   git checkout -b bugfix/{issue-id} develop
   # or from feature branch if bug is there

3. CREATE WS (10 min)
   - Bugfix WS in workstreams/backlog/

4. FIX with TDD (1-4 hours)
   - Write failing test first
   - Implement fix
   - Verify test passes

5. FULL TESTING (30 min)
   - Unit tests
   - Integration tests
   - Regression

6. REVIEW
   - Same process as features

7. MERGE to develop
   - Then to main via regular release
```

===============================================================================
# 3. BUGFIX WS TEMPLATE

```markdown
## WS-BUG-{ID}: {Bug Title}

### 🎯 Goal

**Bug:** {description}
**Expected:** {what should happen}
**Actual:** {what happens}

**Acceptance Criteria:**
- [ ] Bug is fixed
- [ ] Test covers the bug case
- [ ] No regression in related functionality
- [ ] Coverage ≥ 80%

---

### Context

**Reported by:** {user/system}
**Affected:** {who/what}
**Severity:** P1/P2

**Reproduction steps:**
1. ...
2. ...
3. ...

### Root Cause Analysis

{Why does this happen?}

### Dependency

{Related WS / Independent}

### Input Files

- `path/to/affected/file.py` — contains the bug

### Steps

1. Write failing test
2. Fix the bug
3. Verify test passes
4. Run regression

### Code

```python
# Test first
def test_large_repo_clone():
    result = clone_repo("large-repo-url")
    assert result.success is True
    assert result.size > 1_000_000_000  # 1GB

# Then fix
```

### Expected Result

- Bug fixed
- Test added
- No regression

### Scope Estimate

- Files: ~N
- Lines: ~N (SMALL)
- Tokens: ~N

### Completion Criteria

```bash
# Specific test passes
pytest tests/unit/test_bug_{id}.py -v

# Coverage
pytest --cov=src/module --cov-fail-under=80

# Regression
pytest tests/unit/ -m fast -v
```
```

===============================================================================
# 4. TDD FOR BUGFIX

### 4.1 Write Failing Test First

```python
# tests/unit/test_bug_002.py

def test_large_repo_clone_completes():
    """
    Bug: Large repos (>1GB) fail to clone
    Expected: Should complete with increased timeout
    """
    executor = GitExecutor(timeout=600)  # 10 min
    
    result = executor.clone("https://github.com/large/repo.git")
    
    assert result.success is True
    assert result.error is None
```

```bash
# Run — should FAIL
pytest tests/unit/test_bug_002.py -v
# Expected: FAILED (timeout)
```

### 4.2 Implement Fix

```python
# src/module/git_executor.py

class GitExecutor:
    def __init__(self, timeout: int = 300) -> None:
        self.timeout = timeout
    
    def clone(self, url: str) -> CloneResult:
        # Fix: Use configurable timeout
        try:
            subprocess.run(
                ["git", "clone", "--progress", url],
                timeout=self.timeout,
                check=True
            )
            return CloneResult(success=True)
        except subprocess.TimeoutExpired:
            return CloneResult(success=False, error="Clone timeout")
```

### 4.3 Verify Test Passes

```bash
# Run — should PASS
pytest tests/unit/test_bug_002.py -v
# Expected: PASSED
```

===============================================================================
# 5. FULL TESTING

```bash
# 1. Unit tests for fix
pytest tests/unit/test_bug_{id}.py -v

# 2. Unit tests for module
pytest tests/unit/test_git_executor.py -v

# 3. Integration tests
pytest tests/integration/test_git_submissions.py -v

# 4. Regression (all fast tests)
pytest tests/unit/ -m fast -q

# 5. Coverage
pytest --cov=src/module --cov-fail-under=80
```

===============================================================================
# 6. GIT WORKFLOW

### 6.1 Branch

```bash
# From develop
git checkout develop
git pull origin develop
git checkout -b bugfix/{issue-id}
```

### 6.2 Commits

```bash
# Test first
git add tests/
git commit -m "test(git): add test for large repo clone bug

Covers: BUG-002
Expected behavior: clone completes with 10 min timeout"

# Fix
git add src/
git commit -m "fix(git): increase clone timeout for large repos

Root cause: Default 5 min timeout too short for >1GB repos
Solution: Make timeout configurable, default 10 min

Fixes: #002"

# Execution report
git add docs/workstreams/
git commit -m "docs(bug): BUG-002 execution report

Goal achieved: YES
Coverage: 85%"
```

### 6.3 PR to Develop

```bash
gh pr create --base develop --title "fix(git): BUG-002 large repo clone timeout" \
  --body "## Bug Fix

**Issue:** Large repos (>1GB) fail to clone
**Root Cause:** Default timeout too short
**Solution:** Configurable timeout, increased default

## Testing
- Unit test added
- Integration tested
- Regression passed
- Coverage: 85%

Fixes #002"
```

===============================================================================
# 7. OUTPUT FORMAT

```markdown
## ✅ Bugfix Complete: BUG-{ID}

**Issue:** {description}
**Root Cause:** {why it happened}
**Solution:** {what was changed}

### Changes

| File | Change |
|------|--------|
| `src/module/git_executor.py` | Increased timeout |
| `tests/unit/test_bug_002.py` | Added test |

### Testing

- ✅ Specific test passes
- ✅ Module tests pass
- ✅ Integration tests pass
- ✅ Regression passed
- ✅ Coverage: 85%

### Git

- Branch: `bugfix/{issue-id}`
- PR: #{number} → develop

### Next Steps

1. Await PR review
2. Merge to develop
3. Will go to prod in next release
```

===============================================================================
# 8. THINGS YOU MUST NEVER DO

❌ Skip TDD (test after fix)
❌ Deploy directly to production
❌ Skip regression tests
❌ Merge without review
❌ Fix multiple bugs in one WS

===============================================================================
