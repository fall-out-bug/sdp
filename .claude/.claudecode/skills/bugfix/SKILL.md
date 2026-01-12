---
name: bugfix
description: Quality bug fixes (P1/P2). Full TDD cycle, branch from feature/develop, no production deploy.
tools: Read, Write, Edit, Bash, Glob, Grep
---

# /bugfix - Quality Bug Fixes

Standard bug fixes with full quality cycle.

## When to Use

- P1 (HIGH) or P2 (MEDIUM) issues
- Feature broken but not production
- Reproducible errors
- Can wait for proper testing

## Invocation

```bash
/bugfix "description" --feature=F23 --issue-id=002
```

## Master Prompt

📄 **sdp/prompts/commands/bugfix.md** (530+ lines)

**Contains:**
- Full TDD workflow
- Branch strategy (feature/* or bugfix/* from develop)
- Complete test suite
- Quality gates
- Merge to develop (not main!)
- GitHub issue closure

## Workflow

1. Create bugfix branch (from feature or develop)
2. Implement fix with TDD
3. Full test suite
4. Quality checks
5. Merge to appropriate branch
6. Close issue

## Key Difference from Hotfix

| Aspect | Hotfix | Bugfix |
|--------|--------|--------|
| Severity | P0 | P1/P2 |
| Branch from | main | develop/feature |
| Testing | Fast | Full |
| Deploy | Production | Staging |
| SLA target | < 2h | < 24h |

## Output

Bug fixed in feature/develop branch + tests + issue closed

## Quick Reference

**Input:** P1/P2 issue  
**Output:** Quality fix with full tests  
**Next:** Merge to develop, later to main
