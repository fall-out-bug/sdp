# /bugfix — Quality Bug Fixes

When calling `/bugfix "description" --feature=F23 --issue-id=002`:

1. Load full prompt: `@prompts/commands/bugfix.md`
2. Create bugfix branch (from feature/* or bugfix/* from develop)
3. Implement fix with full TDD
4. Complete test suite
5. Quality gates (coverage, linters)
6. Merge to appropriate branch (not main!)
7. Close GitHub issue

## Quick Reference

**Input:** P1/P2 issue
**Output:** Quality fix with full tests

**Key Difference from Hotfix:**
| Aspect | Hotfix | Bugfix |
|--------|--------|--------|
| Severity | P0 | P1/P2 |
| Branch from | main | develop/feature |
| Testing | Fast | Full |
| Deploy | Production | Staging |
| Timeline | < 2h | < 24h |

**Next:** Merge to develop, later to main
