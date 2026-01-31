# 004: Hooks module coverage ≥80%

**Source:** WS 00-020-01 Execution Report (AC4 partial)
**Status:** Open
**Priority:** P2

## Problem

WS 00-020-01 extracted Git hooks to Python. AC4 required coverage ≥80% for new hook modules. Achieved: 73%.

## Current State

- `sdp.hooks.common`: 86%
- `sdp.hooks.pre_commit_checks`: 90%
- `sdp.hooks.pre_commit`: 63%
- `sdp.hooks.pre_push`: 73%
- `sdp.hooks.post_build`: 64%
- `sdp.hooks.pre_deploy`: 64%
- **Total:** 73%

## Remediation

1. Add integration tests that run hooks as subprocess (coverage subprocess mode)
2. Or add unit tests with comprehensive subprocess mocks for `main()` flows
3. Or add `# pragma: no cover` to `main()` entry points with justification in Execution Report

## Related

- WS 00-020-01: Extract Git Hooks to Python
- docs/workstreams/backlog/00-020-01-extract-hooks-to-python.md (Execution Report)
