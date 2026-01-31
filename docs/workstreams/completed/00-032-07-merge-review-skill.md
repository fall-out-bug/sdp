---
assignee: null
completed: '2026-01-30'
depends_on:
- 00-032-05
feature: F032
github_issue: null
project_id: 0
size: MEDIUM
status: completed
traceability:
- ac_description: '`SKILL.md` содержит ≤100 строк'
  ac_id: AC1
  confidence: 1.0
  status: mapped
  test_file: tests/integration/test_guard_flow.py
  test_name: test_happy_path_activate_edit_complete
- ac_description: Traceability check добавлен в workflow
  ac_id: AC2
  confidence: 1.0
  status: mapped
  test_file: tests/integration/test_guard_flow.py
  test_name: test_edit_blocked_without_active_ws
- ac_description: Полный checklist в `docs/reference/review-spec.md`
  ac_id: AC3
  confidence: 1.0
  status: mapped
  test_file: tests/integration/test_guard_flow.py
  test_name: test_edit_blocked_outside_scope
- ac_description: Validation passes
  ac_id: AC5
  confidence: 1.0
  status: mapped
  test_file: tests/unit/test_scope_manager.py
  test_name: test_is_in_scope_with_restricted_scope
- ac_description: '`sdp review trace` команда вызывается в workflow'
  ac_id: AC4
  confidence: 1.0
  status: mapped
  test_file: tests/integration/test_guard_flow.py
  test_name: test_concurrent_activation_blocked
ws_id: 00-032-07
---

## 00-032-07: Merge review.md into SKILL.md

### 🎯 Goal

**What must WORK after completing this WS:**
- `.claude/skills/review/SKILL.md` сокращён до ≤100 строк
- Добавлена проверка AC→Test traceability
- Детальный checklist в `docs/reference/review-spec.md`

**Acceptance Criteria:**
- [ ] AC1: `SKILL.md` содержит ≤100 строк
- [ ] AC2: Traceability check добавлен в workflow
- [ ] AC3: Полный checklist в `docs/reference/review-spec.md`
- [ ] AC4: `sdp review trace` команда вызывается в workflow
- [ ] AC5: Validation passes

**⚠️ WS is NOT complete until Goal is achieved (all AC ✅).**

---

### Context

**Problem**: Review проверяет код, но не требования. AC могут быть не покрыты тестами.

**Solution**: Добавить traceability check + сократить skill.

### Dependencies

- **00-032-05**: Skill Template Standard

### Steps

1. **Rewrite SKILL.md (~80 lines)**

   ```markdown
   # .claude/skills/review/SKILL.md
   ---
   name: review
   description: Quality review with traceability check
   tools: Read, Bash, Grep
   ---
   
   # @review - Quality Review
   
   Review feature by validating workstreams against quality gates and traceability.
   
   ## Quick Reference
   
   | Step | Action | Gate |
   |------|--------|------|
   | 1 | List WS | All WS found |
   | 2 | Traceability | All ACs have tests |
   | 3 | Quality gates | All checks pass |
   | 4 | Goal check | All ACs achieved |
   | 5 | Verdict | APPROVED or CHANGES_REQUESTED |
   
   ## Workflow
   
   ### Step 1: List Workstreams
   
   ```bash
   sdp ws list --feature {F-ID}
   ```
   
   ### Step 2: Check Traceability (NEW)
   
   ```bash
   sdp trace check {WS-ID}
   ```
   
   **Gate:** All ACs must have mapped tests.
   
   | AC | Test | Status |
   |----|------|--------|
   | AC1 | test_func_1 | ✅ |
   | AC2 | - | ❌ MISSING |
   
   If any AC missing → CHANGES_REQUESTED
   
   ### Step 3: Quality Gates
   
   ```bash
   sdp quality check --full
   ```
   
   Must pass:
   - Coverage ≥80%
   - mypy --strict
   - No except:pass
   - Files <200 LOC
   
   ### Step 4: Goal Achievement
   
   For each AC in WS:
   - [ ] Test exists and passes
   - [ ] Implementation matches description
   
   ### Step 5: Verdict
   
   **APPROVED** if:
   - All ACs traceable to tests
   - All tests pass
   - All quality gates pass
   
   **CHANGES_REQUESTED** if any fails.
   
   No middle ground. No "approved with notes."
   
   ## Quality Gates
   
   See [Quality Gates Reference](../../docs/reference/quality-gates.md)
   
   ## Errors
   
   | Error | Cause | Fix |
   |-------|-------|-----|
   | Missing trace | AC has no test | Add test for AC |
   | Coverage <80% | Insufficient tests | Add more tests |
   | Goal not met | AC not working | Fix implementation |
   
   ## See Also
   
   - [Full Review Spec](../../docs/reference/review-spec.md)
   - [Traceability Guide](../../docs/reference/traceability.md)
   ```

2. **Create full spec**

   Move all 17 checks from `prompts/commands/review.md` to `docs/reference/review-spec.md`.

3. **Update CLI**

   Ensure `sdp trace check` exists (from WS-21) or add placeholder.

### Output Files

- `.claude/skills/review/SKILL.md` (rewritten)
- `docs/reference/review-spec.md` (full checklist)
- `docs/reference/traceability.md`

### Completion Criteria

```bash
# Skill is short
wc -l .claude/skills/review/SKILL.md
# Expected: ≤100

# Validation
sdp skill validate .claude/skills/review/SKILL.md
```

---

## Execution Report

**Executed by:** ______  
**Date:** ______

### Goal Status
- [ ] AC1-AC5 — ✅

**Goal Achieved:** ______
