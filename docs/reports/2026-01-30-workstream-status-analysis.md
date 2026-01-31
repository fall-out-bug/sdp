# Workstream Status Analysis

**Date:** 2026-01-30  
**Scope:** docs/workstreams/backlog  
**Purpose:** Correct statuses, move completed, identify outdated

---

## Analysis Summary

### 1. F032 (00-032-*) — COMPLETED (move to completed/)

**Evidence:** F032 execution report confirms all 28 WS completed. Codebase verification:

| WS | Evidence | Action |
|----|----------|--------|
| 00-032-01 | src/sdp/guard/skill.py | → completed |
| 00-032-02 | src/sdp/cli/guard.py | → completed |
| 00-032-03 | src/sdp/guard/tracker.py | → completed |
| 00-032-04 | tests/integration/test_guard_flow.py | → completed |
| 00-032-05 | templates/skill-template.md | → completed |
| 00-032-06 | .claude/skills/build/SKILL.md (88 lines) | → completed |
| 00-032-07 | .claude/skills/review/SKILL.md | → completed |
| 00-032-08 | .claude/skills/design/SKILL.md | → completed |
| 00-032-09 | prompts/commands/ only README.md | → completed |
| 00-032-10 | .cursorrules (45 lines) | → completed |
| 00-032-11 | Already in completed/ | Delete from backlog |
| 00-032-12 | Already in completed/ | Delete from backlog |
| 00-032-13 | Already in completed/ | Delete from backlog |
| 00-032-14 | Already in completed/ | Delete from backlog |
| 00-032-15 | src/sdp/beads/ (install check) | → completed |
| 00-032-16 | src/sdp/beads/ (CLI wrapper) | → completed |
| 00-032-17 | src/sdp/beads/scope_manager.py | → completed |
| 00-032-18 | src/sdp/beads/sync_service.py | → completed |
| 00-032-19 | Beads factory auto-detect | → completed |
| 00-032-20 | src/sdp/traceability/models.py | → completed |
| 00-032-21 | src/sdp/cli/trace.py | → completed |
| 00-032-22 | src/sdp/traceability/detector.py | → completed |
| 00-032-23 | Review skill + CI trace check | → completed |
| 00-032-24 | src/sdp/validators/time_estimate_checker.py | → completed |
| 00-032-25 | templates/workstream-v2.md | → completed |
| 00-032-26 | src/sdp/validators/ws_completion.py | → completed |
| 00-032-27 | hooks/post-ws-complete.sh | → completed |
| 00-032-28 | src/sdp/validators/supersede_checker.py | → completed |

**00-032-00** (feature overview): status → completed (feature done)

---

### 2. F013 (00-013-*) — COMPLETED (move to completed/)

**Evidence:** All skills exist in .claude/skills/

| WS | Skill/Code | Action |
|----|------------|--------|
| 00-013-01 | src/sdp/schema/validator.py, docs/schema/intent.schema.json | → completed |
| 00-013-02 | .claude/skills/tdd/, src/sdp/tdd/runner.py | → completed |
| 00-013-03 | .claude/skills/debug/ | → completed |
| 00-013-04 | .claude/skills/feature/ | → completed |
| 00-013-05 | .claude/skills/idea/ | → completed |
| 00-013-06 | .claude/skills/design/ | → completed |
| 00-013-07 | .claude/skills/build/ | → completed |
| 00-013-08 | .claude/skills/oneshot/ | → completed |
| 00-013-09 | .claude/skills/think/ | → completed |

---

### 3. F012 (00-012-*) — SUPERSEDED (stay in backlog)

**Status:** Already correct (status: superseded, superseded_by: 00-032-01)

**Action:** Add supersede_reason for clarity. No file move — superseded items stay in backlog per SDP convention.

```
supersede_reason: "F012 daemon/agent framework superseded by F032 Guard + Beads integration"
```

---

### 4. F020 (00-020-*) — BACKLOG (keep)

**Evidence:** src/sdp/hooks/ exists but contains only ws_complete.py (F032). Full hook extraction (pre_commit.py, pre_push.py etc.) NOT done.

| WS | Status | Reason |
|----|--------|--------|
| 00-020-01 | backlog | Hooks still bash, not Python |
| 00-020-02 | backlog | Depends on 00-020-01 |

---

### 5. F025, F030, F031 — BACKLOG (keep)

| WS | Status | Reason |
|----|--------|--------|
| 00-025-01 | backlog | pip-audit security not added |
| 00-030-01 | backlog | GitHub integration tests |
| 00-030-02 | backlog | Adapter tests |
| 00-030-03 | backlog | Core functionality tests |
| 00-031-01 | backlog | SDPError migration |

---

### 6. BEADS-001 — BACKLOG (keep, update)

**Analysis:** F032 implemented Beads CLI, scope, sync. BEADS-001 asks for skill-level integration (@idea creates Beads task, etc.). Partially done — skills use Beads. Mark as backlog with note.

**Action:** Add status: backlog, add note "Partially implemented by F032 (infrastructure). Remaining: full skill-Beads workflow verification."

---

## Execution Plan

1. **Delete from backlog** (duplicates): 00-032-11, 00-032-12, 00-032-13, 00-032-14
2. **Move to completed** (24 files): 00-032-01..10, 00-032-15..28
3. **Move to completed** (9 files): 00-013-01..09
4. **Update in backlog**: 00-032-00 (status: completed), F012 (add supersede_reason), BEADS-001 (add note)

**Total moves:** 33 files  
**Total deletes:** 4 files (duplicates)  
**Stay in backlog:** 6 files (F020×2, F025×1, F030×3, F031×1, BEADS-001×1, F032-00×1)

---

## Execution Complete (2026-01-30)

**Actions performed:**
1. ✅ Deleted F032 duplicates from backlog (00-032-11..14)
2. ✅ Moved 24 F032 workstreams to completed (01-10, 15-28)
3. ✅ Moved 9 F013 workstreams to completed (01-09)
4. ✅ Updated 00-032-00 status: planning → completed
5. ✅ Added supersede_reason to all 14 F012 workstreams
6. ✅ Added note to BEADS-001 (partial F032 implementation)
7. ✅ Moved 14 F012 workstreams to `docs/workstreams/canceled/`

**Backlog now contains:** 10 files
- 2 F020 (backlog)
- 1 F025 (backlog)
- 3 F030 (backlog)
- 1 F031 (backlog)
- 1 F032-00 (feature overview, completed)
- 1 BEADS-001 (backlog)

**Canceled folder:** 14 F012 workstreams (superseded by F032)
