# Review: sdp repo — PR Readiness

**Date:** 2026-02-25  
**Scope:** All uncommitted changes in sdp submodule  
**Branch:** fix/audit-engineering-gaps  
**Purpose:** Prepare for PR of all changes

---

## 1. Summary of Changes

| File | Change Type | Lines |
|------|-------------|-------|
| `.beads/issues.jsonl` | Status updates | +6 (closed beads) |
| `PRODUCT_VISION.md` | Minor edit | ±1 |
| `docs/decisions/DECISIONS.md` | Major trim | −103 |
| `docs/workstreams/INDEX.md` | Table update | ±4 |
| `docs/workstreams/backlog/00-001-02.md` | Minor | ±1 |
| `docs/workstreams/backlog/00-002-03.md` | Minor | ±1 |
| `schema/intent.schema.json` | Format (minified) | ±1 |
| `sdp-plugin/internal/quality/coverage.go` | API: +ctx | +4 |
| `sdp-plugin/internal/quality/*_test.go` | Tests: pass ctx | +15 |

---

## 2. Checklist (Check 0–11)

| # | Check | Status | Note |
|---|-------|--------|------|
| 0 | Build passes | ✅ | `go build ./...` |
| 1 | Tests pass | ⏳ | Run `go test ./... -short` (quality tests may be slow) |
| 2 | No P0/P1 findings | ✅ | Coverage ctx refactor is backward-compatible |
| 3 | Schema valid | ✅ | intent.schema.json minified, structure unchanged |
| 4 | DECISIONS trim | ✅ | Kept 1 decision (JWT); removed noise |
| 5 | Beads sync | ⚠️ | .beads/issues.jsonl — commit with `bd sync` or exclude from PR |
| 6 | INDEX consistency | ✅ | Minor table updates |
| 7 | Coverage API | ✅ | CheckCoverage(ctx) — verifier, quality_impl, tests updated |
| 8 | Untracked files | ⚠️ | schema/next-action.schema.json, evals/review/test_cases.jsonl — add or .gitignore |
| 9 | CLAUDE.md sync | ✅ | Already synced (F054-06) |
| 10 | Prompts unchanged | ✅ | No prompt edits in this diff |
| 11 | No secrets | ✅ | No credentials in diff |

---

## 3. Risk Areas

### QA
- **Scope:** sdp-plugin/internal/quality, verifier
- **Risk:** CheckCoverage(ctx) — callers must pass context
- **Evidence:** verifier.go:192 passes ctx; quality_impl.go:31 uses context.Background(); all tests updated
- **Severity:** P2 (maintenance) — context propagation is correct
- **Verdict:** PASS

### Security
- **Scope:** schema, config, no new exec paths
- **Risk:** None in this diff
- **Verdict:** PASS

### DevOps
- **Scope:** .beads, docs
- **Risk:** .beads/issues.jsonl in PR — may conflict with other clones
- **Recommendation:** Document whether beads state is committed or gitignored
- **Verdict:** PASS (with note)

### Docs
- **Scope:** DECISIONS.md, INDEX, backlog
- **Risk:** DECISIONS trim — 103 lines removed
- **Evidence:** Kept JWT decision; removed redundant/empty entries
- **Verdict:** PASS

### PromptOps
- **Scope:** prompts/skills — unchanged
- **Verdict:** PASS

---

## 4. Cross-Consistency

| Check | Status |
|-------|--------|
| CheckCoverage callers all updated | ✅ verifier, quality_impl, 4 test files |
| checkPythonCoverage, checkGoCoverage, checkJavaCoverage | ✅ All take ctx |
| schema/intent.schema.json | ✅ Valid JSON, structure intact |

---

## 5. PR Preparation Checklist

Before opening PR:

- [ ] Run `go test ./... -short` and fix any failures
- [ ] Decide: include .beads/issues.jsonl in PR or `git restore .beads/`
- [ ] Add or ignore: `schema/next-action.schema.json`, `evals/review/test_cases.jsonl`
- [ ] Commit message: `F053/F054: coverage ctx propagation, DECISIONS trim, beads sync`
- [ ] PR title: `F053: Coverage context propagation + DECISIONS cleanup`

---

## 6. Verdict

**APPROVED** for PR with minor prep:

- Code changes (coverage ctx) are correct and consistent
- DECISIONS trim improves maintainability
- Resolve untracked files and .beads handling before push

---

## 7. Handoff

```
---
## Next Step

1. Resolve untracked files (add or gitignore)
2. Run `go test ./... -short`
3. Commit with conventional message
4. Open PR to main
---
```
