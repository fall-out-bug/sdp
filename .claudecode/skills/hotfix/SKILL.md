---
name: hotfix
description: Emergency P0 fixes. Fast-track production deployment with minimal changes. Branch from main, deploy < 2h.
tools: Read, Write, Edit, Bash, Glob, Grep
---

# /hotfix - Emergency Production Fixes

Fast-track critical bug fixes for production.

## When to Use

- P0 CRITICAL issues only
- Production down or severely degraded
- All/most users affected
- Data loss/corruption risk

## Invocation

```bash
/hotfix "description" --issue-id=001
```

## Master Prompt

📄 **sdp/prompts/commands/hotfix.md** (420+ lines)

**Contains:**
- Emergency workflow (SLA target: < 2h)
- Branch from `main` (not develop!)
- Minimal changes only
- Fast testing (smoke + critical path)
- Production deployment
- Backport to develop + feature branches
- GitHub issue closure
- Telegram notification

## Workflow

1. Create hotfix branch from `main`
2. Implement minimal fix
3. Fast testing (no full test suite)
4. Deploy to production
5. Monitor (5 min verification)
6. Merge to `main` + tag
7. Backport to `develop` and all `feature/*`
8. Close issue

## Key Rules

- **Minimal changes** - no refactoring!
- **No new features** - fix bug only
- **Fast testing** - smoke + critical path
- **SLA target: < 2h** - emergency window (not a scope estimate)
- **Backport mandatory** - to all branches

## Output

Hotfix deployed to production + backported + issue closed

## Quick Reference

**Input:** P0 issue  
**Output:** Production fix (SLA target: < 2h)  
**Next:** Monitor + postmortem
