# F063: Guardian-style Hooks and Guard Rails for SDP — Workstream Summary

> Feature: planned via `idea-guardian-hooks-guardrails` | 4 workstreams | ~8h total

## Mandatory Requirement

- Every feature/workstream description in F063 must include `Source of inspiration` with an explicit link to the reference artifact.

## Source of inspiration

- [guardian-cli README](https://github.com/AlexGladkov/guardian-cli/blob/main/README.md)
- [guardian-cli hooks manager](https://github.com/AlexGladkov/guardian-cli/blob/main/internal/git/hooks.go)
- [guardian-cli enforcement command](https://github.com/AlexGladkov/guardian-cli/blob/main/internal/cli/check.go)
- [guardian-cli exceptions and meta-check](https://github.com/AlexGladkov/guardian-cli/blob/main/internal/engine/engine.go)

## Dependency Graph

```
00-063-01 Hook Ownership & Installer Unification
    ├─► 00-063-02 Staged Guard Checks + CI Diff Range
    ├─► 00-063-03 Config Alignment (rules/config/docs/CI)
    └─► 00-063-04 Exceptions TTL + Governance Meta-check
```

## Workstream Registry

| WS ID | Title | Est. | Depends On |
|-------|-------|------|------------|
| 00-063-01 | Hook Ownership and Installer Unification | 2h | - |
| 00-063-02 | Staged Guard Checks and CI Diff Range | 2.5h | 00-063-01 |
| 00-063-03 | Guard Configuration and Gate Alignment | 1.5h | 00-063-01 |
| 00-063-04 | Exceptions TTL and Governance Meta-check | 2h | 00-063-02, 00-063-03 |
