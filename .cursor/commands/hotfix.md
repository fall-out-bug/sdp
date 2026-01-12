# /hotfix — Emergency Production Fixes

When calling `/hotfix "description" --issue-id=001`:

1. Load full prompt: `@prompts/commands/hotfix.md`
2. Create hotfix branch from `main`
3. Implement minimal fix (no refactoring!)
4. Fast testing (smoke + critical path only)
5. Deploy to production
6. Monitor (5 min verification)
7. Merge to `main` + tag `hotfix-{ID}`
8. Backport to `develop` and all `feature/*` branches
9. Close GitHub issue
10. Send Telegram notification

## Quick Reference

**Input:** P0 CRITICAL issue
**Output:** Production fix < 2h

**Key Rules:**
- Minimal changes only
- No refactoring
- No new features
- Fast testing
- Backport mandatory

**Timeline:** < 2 hours from detection to deployment

**Next:** Monitor + postmortem
