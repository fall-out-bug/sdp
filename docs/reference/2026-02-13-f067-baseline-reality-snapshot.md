# F067 Baseline Reality Snapshot

> **Date:** 2026-02-13
> **Purpose:** Anchor all F067 improvements with deterministic baseline metrics

---

## 1. Toolchain Version Matrix

| Location | Current Version | Target | Gap |
|----------|-----------------|--------|-----|
| `go.mod` (root) | 1.25.6 | 1.24 (stable) | Invalid version |
| `sdp-plugin/go.mod` | 1.26 | 1.24 (stable) | Invalid version |
| CI workflows (5 files) | 1.26 | 1.24 (stable) | Invalid version |
| `sdp-plugin/Dockerfile` | 1.21-alpine | 1.24-alpine | 3 versions behind |
| `README.md` badge | 1.24+ | 1.24 | Correct |

**Target:** WS-04 will align all to Go 1.24

---

## 2. CI Stability Baseline

### Workflows (5 total)
- `.github/workflows/go-ci.yml` - Main CI
- `.github/workflows/go-release.yml` - Release builds
- `.github/workflows/release-action.yml` - Action release
- `.github/workflows/sdp-verify-dogfood.yml` - Dogfood verification
- `.github/workflows/test-verify-action.yml` - Action tests

### Known Issues
- Go versions 1.25.6 and 1.26 don't exist yet
- Docker image uses outdated 1.21-alpine

---

## 3. Quality Gate Drift

| Source | Coverage Threshold | Status |
|--------|-------------------|--------|
| `docs/PROTOCOL.md` | 80% | Documented |
| `.sdp/config.yml` | 80% | Config |
| `.github/workflows/go-ci.yml` | **60%** | MISMATCH |
| `.sdp/guard-rules.yml` | - | Other rules |

**Gap:** CI enforces 60% while docs say 80%

**Additional Drift:**
- `quality-gate.toml` contains Python-specific references (pytest, mypy, ruff)
- `ci-gates.toml` is separate from `.sdp/guard-rules.yml`
- Guard checks run with `|| true` in CI (non-blocking)

**Target:** WS-05 will align thresholds and consolidate configs

---

## 4. Prompt/Agent Duplication

| Location | File Count | Status |
|----------|------------|--------|
| `prompts/agents/` | 25 | Canonical |
| `sdp-plugin/prompts/agents/` | 11 | **Duplicate tree** |
| `.claude/agents` | Symlink | Points to `../prompts/agents` |
| `.claude/skills` | Symlink | Points to `../prompts/skills` |

**Issue:** 11 duplicate agent files in `sdp-plugin/prompts/agents/` may diverge from canonical source.

**Target:** WS-02 will remove duplicates and add drift detection

---

## 5. Repository Hygiene

### Tracked Artifacts (should not be in git)

| File | Issue | Action |
|------|-------|--------|
| `coverage_quality.out` | Build artifact | Remove from tracking |
| `.sdp/log/events.jsonl` | Runtime evidence | Keep (audit trail) |
| `sdp-plugin/.sdp/log/events.jsonl` | Runtime evidence | Keep |
| `sdp-plugin/cmd/sdp/.sdp/log/events.jsonl` | Runtime evidence | Keep |

### Binary Files (scripts, acceptable)
- `.claude/hooks/*.sh` - Hook scripts (OK)
- `.github/actions/verify/scripts/*.sh` - CI scripts (OK)

### Python Artifacts
- `poetry.lock` exists at root (Go project)
- `quality-gate.toml` references Python tools

**Target:** WS-07 will clean up artifacts and add gitignore rules

---

## 6. LOC Violation Baseline

### Files Exceeding 200 LOC (19 files)

| File | LOC | Over By |
|------|-----|---------|
| `sdp-plugin/cmd/sdp/contract.go` | 580 | 380 |
| `src/sdp/agents/synthesis_agent.go` | 560 | 360 |
| `src/sdp/agents/code_analyzer.go` | 543 | 343 |
| `src/sdp/agents/contract_validator.go` | 485 | 285 |
| `sdp-plugin/internal/doctor/doctor.go` | 468 | 268 |
| `sdp-plugin/cmd/sdp/decisions.go` | 415 | 215 |
| `sdp-plugin/internal/orchestrator/feature_coordinator.go` | 403 | 203 |
| `sdp-plugin/internal/ui/dashboard/app.go` | 394 | 194 |
| `sdp-plugin/cmd/sdp/metrics.go` | 386 | 186 |
| `sdp-plugin/internal/telemetry/analyzer.go` | 377 | 177 |
| `src/sdp/agents/contract_generator.go` | 368 | 168 |
| `sdp-plugin/internal/orchestrator/orchestrator.go` | 367 | 167 |
| `sdp-plugin/cmd/sdp/telemetry.go` | 345 | 145 |
| `src/sdp/recovery/lock_recovery.go` | 334 | 134 |
| `sdp-plugin/internal/ui/completion.go` | 332 | 132 |
| `sdp-plugin/internal/hooks/hooks.go` | 306 | 106 |
| `src/sdp/monitoring/metrics.go` | 304 | 104 |
| `src/sdp/cli/contract_commands.go` | 291 | 91 |
| `sdp-plugin/internal/telemetry/tracker.go` | 270 | 70 |

**Total:** 19 files violating 200 LOC rule

**Target:** WS-09 will split top 10 offenders (50% reduction target)

---

## 7. Linter Exclusion Baseline

### `.golangci.yml` Exclusions

| Metric | Count | Status |
|--------|-------|--------|
| Total `linters:` entries | 28 | High |
| Packages with `errcheck` disabled | Multiple | Concern |
| Packages with `gocognit` disabled | Multiple | Concern |

### Key Exclusions (from file)
- Test files: `gocyclo`, `gocognit`, `errcheck` disabled
- Multiple production packages excluded from `errcheck`

**Target:** WS-14 will reduce exclusions by 50%

---

## 8. Target Metrics by Workstream

| WS | Current State | Target |
|----|---------------|--------|
| **WS-02** | 11 duplicate prompt files | 0 duplicates |
| **WS-03** | Broken adapter paths | All paths valid |
| **WS-04** | 5 Go versions, 2 invalid | 1 version (1.24) |
| **WS-05** | CI 60% vs docs 80% | Aligned at 80% |
| **WS-06** | Release gap (v0.7.1 vs v0.10.0) | Releases current |
| **WS-07** | `coverage_quality.out` tracked | Removed |
| **WS-09** | 19 files > 200 LOC | <= 10 files |
| **WS-10** | 2 modules, no go.work | go.work + docs |
| **WS-11** | Self-submodule, poetry.lock | Cleaned |
| **WS-12** | No DEVELOPMENT.md | Created |
| **WS-13** | Skills > 500 LOC | <= 300 LOC |
| **WS-14** | 28 linter exclusions | <= 14 exclusions |

---

## 9. Summary

### Critical Issues (Fix First)
1. **Go versions 1.25.6/1.26 don't exist** - WS-04
2. **Coverage threshold mismatch** (60% vs 80%) - WS-05
3. **19 files violate 200 LOC rule** - WS-09

### High Priority
4. **11 duplicate prompt files** - WS-02
5. **Docker image 3 versions behind** - WS-04
6. **28 linter exclusions** - WS-14

### Medium Priority
7. **No go.work for dual modules** - WS-10
8. **Python artifacts remain** - WS-07, WS-11
9. **No DEVELOPMENT.md** - WS-12

---

**Generated:** 2026-02-13
**Feature:** F067 Repository Hardening
**Baseline ID:** F067-BASELINE-001
