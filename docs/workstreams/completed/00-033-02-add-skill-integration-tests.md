---
ws_id: 00-033-02
feature: F032
status: completed
size: MEDIUM
project_id: 00
github_issue: null
assignee: null
depends_on: []
---

## WS-00-033-02: Add BEADS-001 Skill Integration Tests

### 🎯 Goal

**What must WORK after completing this WS:**
- Skill-Beads integration covered by tests (mock + real)
- BEADS-001 Definition of Done satisfied

**Acceptance Criteria:**
- [x] AC1: Tests for @build + Beads (bd update/close, status lifecycle)
- [x] AC2: Tests for @review + Beads (bd list --parent, ws_id resolution)
- [x] AC3: Tests for @idea + Beads (create_task)
- [x] AC4: Tests for @design + Beads (migrate)
- [x] AC5: Mock tests work without bd installed
- [x] AC6: `pytest tests/ -k beads` — all pass

**⚠️ WS is NOT complete until Goal is achieved (all AC ✅).**

---

### Context

**Source:** [BEADS-001 Review 2026-01-30](../../reports/2026-01-30-BEADS-001-review.md) — CHANGES_REQUESTED

**BEADS-001 AC:** "Tests for all skill integrations (mock + real Beads)"

**Reference:** [BEADS-001 Phase 2.3–2.5](BEADS-001-skills-integration.md)

---

### Scope

- Files: tests/unit/beads/, tests/integration/beads/
- LOC: ~150–200

---

### Execution Report

**Date:** 2026-01-30  
**Result:** ✅ COMPLETED

**Summary:**
- AC1–AC5: Already covered by `tests/integration/beads/test_skills_integration.py` (6 tests)
- AC6: Fixed 2 pre-existing beads test failures:
  1. `test_beads_check_installed` — BeadsCLICheck message now includes "installed" for consistency
  2. `test_handles_beads_unavailable` — Added `shutil.which` + `BEADS_USE_MOCK=false` patches so factory attempts CLIBeadsClient; subprocess.run raises FileNotFoundError → BeadsClientError

**Changes:**
- `src/sdp/health_checks/checks.py`: Message `"Beads CLI v{version} at {path}"` → `"Beads CLI installed (v{version}) at {path}"`
- `tests/unified/test_e2e/test_beads_client.py`: Patched `shutil.which`, `@patch.dict('os.environ', {'BEADS_USE_MOCK': 'false'})` for `test_handles_beads_unavailable`

**Verification:** `pytest tests/ -k beads` — 108 passed, 3 skipped
