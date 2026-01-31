# ADR-008: CI Split Strategy

## Status

Accepted

## Context

Current CI workflow uses `continue-on-error: true` for all checks.
This means:
- PRs can merge with failing tests
- Agents learn they can ignore CI failures
- Quality degrades over time

We need to balance:
- Strictness (blocking bad code)
- Developer experience (not blocking on minor issues)

## Decision

Split checks into two categories:

### Critical (Block PR)

| Check | Threshold | Reason |
|-------|-----------|--------|
| Tests | All pass | Core functionality |
| Coverage | ≥80% | Prevent untested code |
| mypy strict | No errors | Type safety |
| ruff errors | No errors | Code quality |

These run without `continue-on-error`. Failure blocks merge.

### Warning (Comment Only)

| Check | Threshold | Reason |
|-------|-----------|--------|
| File size | <200 LOC | AI readability |
| Complexity | CC <10 | Maintainability |
| ruff warnings | Report | Style suggestions |

These run with `continue-on-error: true`. Post comment but don't block.

## Implementation

Two separate workflows:
- `ci-critical.yml` — Required status check
- `ci-warnings.yml` — Informational

## Consequences

### Positive
- Bad code can't merge
- Clear distinction critical vs nice-to-have
- Agents learn which rules are hard requirements

### Negative
- More CI configuration
- Some PRs will be blocked (intended)
- Need branch protection setup
