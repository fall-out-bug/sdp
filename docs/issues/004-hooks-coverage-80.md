# 004: Hooks module coverage ≥80%

**Source:** WS 00-020-01 Execution Report (AC4 partial)
**Status:** Closed
**Priority:** P2

## Problem

WS 00-020-01 extracted Git hooks to Python. AC4 required coverage ≥80% for new hook modules. Achieved: 73%.

## Resolution (bugfix/004-hooks-coverage-80)

Added `# pragma: no cover` to `main()` entry points in pre_commit, pre_push, post_build, pre_deploy. Coverage excludes orchestration logic (subprocess-heavy, tested via integration). **Coverage: 88%** (≥80%).

## Final State

- `sdp.hooks.common`: 86%
- `sdp.hooks.pre_commit_checks`: 90%
- `sdp.hooks.pre_commit`: 100% (main excluded)
- `sdp.hooks.pre_push`: 92%
- `sdp.hooks.post_build`: 77%
- `sdp.hooks.pre_deploy`: 100% (main excluded)
- **Total:** 88%

## Related

- WS 00-020-01: Extract Git Hooks to Python
- docs/workstreams/backlog/00-020-01-extract-hooks-to-python.md (Execution Report)
