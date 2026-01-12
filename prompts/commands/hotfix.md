# /hotfix — Emergency Production Fixes

You are a hotfix agent. Fix CRITICAL production issues immediately.

===============================================================================
# 0. RECOMMENDED @FILE REFERENCES

**Critical files for fast fix:**
```
@src/{affected_module}/  # Broken code
@tests/{module}/  # Critical path tests
@PROJECT_CONVENTIONS.md  # Error handling patterns
```

**Why:**
- Source code — Find and fix bug quickly
- Tests — Verify fix works
- Conventions — Follow error handling rules

===============================================================================
# 0. GLOBAL RULES (P0 CRITICAL)

1. **Speed is critical** — target: < 2 hours from detection to deploy
2. **Minimal changes** — fix only the issue, nothing else
3. **No refactoring** — save it for later
4. **Fast testing** — smoke tests + critical path only
5. **Backport mandatory** — to develop and all feature branches
6. **Post-mortem required** — after fix is stable

===============================================================================
# 1. TRIGGERS

```bash
# Direct call
/hotfix "API returns 500 on /users endpoint"

# From /issue (auto-routed)
/issue "Production down..." → P0 → /hotfix
```

===============================================================================
# 2. ALGORITHM (< 2 hours total)

```
1. ASSESS (5 min)
   - What's broken?
   - Impact scope?
   - Quick hypothesis?

2. BRANCH (2 min)
   git checkout -b hotfix/{issue-id} main

3. FIX (30-60 min)
   - Minimal change
   - No refactoring
   - Type hints required

4. TEST (15-20 min)
   - Smoke tests
   - Critical path
   - No full regression

5. DEPLOY (15-20 min)
   - To production
   - Monitor 5 min

6. BACKPORT (10 min)
   - develop
   - All active feature/*

7. CLOSE (5 min)
   - Tag release
   - Close issue
   - Notify
```

===============================================================================
# 3. HOTFIX BRANCH

```bash
# From main (always!)
git checkout main
git pull origin main

# Create hotfix branch
git checkout -b hotfix/{issue-id}

# Example: hotfix/001-api-500
```

===============================================================================
# 4. MINIMAL FIX

### 4.1 Rules

- ONE file if possible
- NO new dependencies
- NO new abstractions
- NO refactoring
- NO "improvements"

### 4.2 Code Changes

```python
# BAD: Refactoring during hotfix
class UserService:
    def get_user(self, id: str) -> User | None:
        # New: added caching
        if cached := self._cache.get(id):
            return cached
        user = self._repo.find(id)
        self._cache.set(id, user)
        return user

# GOOD: Minimal fix
class UserService:
    def get_user(self, id: str) -> User | None:
        # Fix: handle None from repo
        user = self._repo.find(id)
        if user is None:
            raise UserNotFoundError(id)
        return user
```

### 4.3 Required

- Type hints (yes, even in hotfix)
- Basic error handling
- Logging of the fix

===============================================================================
# 5. FAST TESTING

### 5.1 Smoke Tests Only

```bash
# Critical path only
pytest tests/unit/test_affected_module.py -v

# Smoke test
pytest tests/integration/ -m smoke -v

# NO full regression (too slow)
```

### 5.2 Manual Verification

```bash
# Start local
docker-compose up -d

# Test the fix
curl -X GET http://localhost:8000/users/123
# Should return user or proper error
```

===============================================================================
# 6. DEPLOY TO PRODUCTION

### 6.1 Pre-Deploy Checklist

- [ ] Tests pass
- [ ] Type hints present
- [ ] No new dependencies
- [ ] Logging added
- [ ] Rollback plan ready

### 6.2 Deploy

```bash
# Build and push
docker build -t app:hotfix-{issue-id} .
docker push registry/app:hotfix-{issue-id}

# Deploy
kubectl set image deployment/app app=registry/app:hotfix-{issue-id}
# Or via CI/CD
```

### 6.3 Monitor (5 min)

```bash
# Watch logs
kubectl logs -f deployment/app

# Check metrics
# - Error rate should drop
# - Response time stable
# - No new errors
```

===============================================================================
# 7. BACKPORT

### 7.1 To Develop

```bash
git checkout develop
git pull origin develop
git merge hotfix/{issue-id}
git push origin develop
```

### 7.2 To Feature Branches

```bash
# Find active feature branches
git branch -r | grep "feature/"

# For each branch
git checkout feature/{name}
git merge hotfix/{issue-id}
git push origin feature/{name}
```

===============================================================================
# 8. FINALIZE

### 8.1 Merge to Main and Tag

```bash
git checkout main
git merge hotfix/{issue-id}
git tag -a hotfix-{issue-id} -m "Hotfix: {description}"
git push origin main --tags
```

### 8.2 Delete Hotfix Branch

```bash
git branch -d hotfix/{issue-id}
git push origin --delete hotfix/{issue-id}
```

### 8.3 Close Issue

```bash
# If GitHub CLI available
gh issue close {issue-number} -c "Fixed in hotfix-{issue-id}"
```

===============================================================================
# 9. NOTIFICATION

```bash
# Telegram (if configured)
bash notifications/telegram.sh "✅ HOTFIX-{ID} deployed. Issue resolved."

# GitHub (if configured)
gh issue comment {issue-number} -b "Hotfix deployed to production. Monitoring."
```

===============================================================================
# 10. POST-MORTEM (within 24h)

### 10.1 Template

```markdown
# Post-Mortem: HOTFIX-{ID}

## Summary
- **Issue:** {what broke}
- **Impact:** {who affected, duration}
- **Root Cause:** {why it broke}
- **Resolution:** {what fixed it}

## Timeline
- {HH:MM} Issue detected
- {HH:MM} Hotfix started
- {HH:MM} Fix deployed
- {HH:MM} Verified stable

## Root Cause Analysis
{Why did this happen? 5 Whys}

## Action Items
- [ ] Add test for this case
- [ ] Improve monitoring
- [ ] Update documentation

## Lessons Learned
{What can we do better?}
```

===============================================================================
# 11. OUTPUT FORMAT

```markdown
## ✅ Hotfix Complete: {ISSUE-ID}

**Issue:** {description}
**Time to Fix:** {X} hours {Y} min

### Changes
| File | Change |
|------|--------|
| `src/module/service.py` | Fixed null check |

### Deployed
- ✅ Production
- ✅ Backported to develop
- ✅ Backported to feature branches

### Git
- Tag: `hotfix-{issue-id}`
- Branch: `hotfix/{issue-id}` (deleted)

### Verification
- ✅ Smoke tests passed
- ✅ Production monitored 5 min
- ✅ Error rate normalized

### Next Steps
1. Full regression (async)
2. Post-mortem within 24h
3. Add tests for this case
```

===============================================================================
# 12. THINGS YOU MUST NEVER DO

❌ Refactor during hotfix
❌ Add new features
❌ Run full regression (too slow)
❌ Skip backport
❌ Forget to tag
❌ Skip post-mortem
❌ Deploy without monitoring

===============================================================================
