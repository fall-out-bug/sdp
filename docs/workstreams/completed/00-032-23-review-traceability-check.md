---
assignee: null
completed: '2026-01-30'
depends_on:
- 00-032-22
feature: F032
github_issue: null
project_id: 0
size: SMALL
status: completed
traceability:
- ac_description: '@review skill вызывает `sdp trace check` в workflow'
  ac_id: AC1
  confidence: 1.0
  status: mapped
  test_file: tests/integration/test_guard_flow.py
  test_name: test_happy_path_activate_edit_complete
- ac_description: Verdict CHANGES_REQUESTED если любой AC unmapped
  ac_id: AC2
  confidence: 1.0
  status: mapped
  test_file: tests/integration/test_guard_flow.py
  test_name: test_edit_blocked_without_active_ws
- ac_description: Review output включает traceability table
  ac_id: AC3
  confidence: 1.0
  status: mapped
  test_file: tests/integration/test_guard_flow.py
  test_name: test_edit_blocked_outside_scope
- ac_description: ''
  ac_id: AC5
  confidence: 1.0
  status: mapped
  test_file: tests/unit/test_scope_manager.py
  test_name: test_is_in_scope_with_restricted_scope
- ac_description: CI workflow включает traceability step
  ac_id: AC4
  confidence: 1.0
  status: mapped
  test_file: tests/integration/test_guard_flow.py
  test_name: test_concurrent_activation_blocked
ws_id: 00-032-23
---

## 00-032-23: Review Traceability Check

### 🎯 Goal

**What must WORK after completing this WS:**
- @review автоматически запускает traceability check
- CHANGES_REQUESTED если coverage < 100%
- CI check для traceability

**Acceptance Criteria:**
- [ ] AC1: @review skill вызывает `sdp trace check` в workflow
- [ ] AC2: Verdict CHANGES_REQUESTED если любой AC unmapped
- [ ] AC3: Review output включает traceability table
- [ ] AC4: CI workflow включает traceability step

**⚠️ WS is NOT complete until Goal is achieved (all AC ✅).**

---

### Context

**Problem**: Review проверяет код, но не requirements traceability.

**Solution**: Интеграция traceability check в review workflow.

### Dependencies

- **00-032-22**: Auto-Detect Test Coverage

### Steps

1. **Update review skill**

   Add traceability check to `.claude/skills/review/SKILL.md`:

   ```markdown
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
   ```

2. **Update CI workflow**

   ```yaml
   # .github/workflows/ci-critical.yml (add step)
   - name: Check traceability
     run: |
       # Get WS ID from PR title or branch
       WS_ID=$(echo "${{ github.head_ref }}" | grep -oE '[0-9]{2}-[0-9]{3}-[0-9]{2}' || echo "")
       
       if [ -n "$WS_ID" ]; then
         poetry run sdp trace check "$WS_ID" || {
           echo "::error::Traceability check failed. Some ACs don't have tests."
           exit 1
         }
       else
         echo "No WS ID found in branch name, skipping traceability check"
       fi
   ```

3. **Create integration test**

   ```python
   # tests/integration/test_review_traceability.py
   import pytest
   from typer.testing import CliRunner
   from sdp.cli.main import app
   
   runner = CliRunner()
   
   class TestReviewWithTraceability:
       def test_review_fails_on_missing_traceability(
           self,
           mock_beads_incomplete_trace
       ):
           """Review should fail if traceability incomplete."""
           # This would be part of the full @review flow
           result = runner.invoke(app, ["trace", "check", "00-032-01"])
           
           assert result.exit_code == 1
           assert "INCOMPLETE" in result.output
       
       def test_review_passes_with_complete_traceability(
           self,
           mock_beads_complete_trace
       ):
           """Review should pass if traceability complete."""
           result = runner.invoke(app, ["trace", "check", "00-032-01"])
           
           assert result.exit_code == 0
           assert "COMPLETE" in result.output
   ```

4. **Update docs/reference/review-spec.md**

   Add traceability section to full review specification.

### Output Files

- `.claude/skills/review/SKILL.md` (updated)
- `.github/workflows/ci-critical.yml` (updated)
- `docs/reference/review-spec.md` (updated)
- `tests/integration/test_review_traceability.py`

### Completion Criteria

```bash
# Skill mentions traceability
grep -q "trace check" .claude/skills/review/SKILL.md

# CI has traceability step
grep -q "traceability" .github/workflows/ci-critical.yml

# Integration tests pass
pytest tests/integration/test_review_traceability.py -v
```

---

## Execution Report

**Executed by:** ______  
**Date:** ______

### Goal Status
- [ ] AC1-AC4 — ✅

**Goal Achieved:** ______

---

# Feature Complete

This is the final workstream of F032: SDP Protocol Enforcement Overhaul.

## Summary

After completing all 23 workstreams:

1. **Guard Skill** (WS-01 to WS-04) — Pre-edit enforcement
2. **Prompt Consolidation** (WS-05 to WS-10) — Short skills, deleted prompts/commands/
3. **CI Enforcement** (WS-11 to WS-14) — Critical blocks, warnings comment
4. **Real Beads** (WS-15 to WS-19) — Full Beads integration
5. **Traceability** (WS-20 to WS-23) — AC→Test tracking

## Expected Outcomes

- Agents cannot edit without active WS
- Agents follow short, focused prompts
- CI blocks bad code
- Beads tracks all state
- Review verifies requirements match implementation

## Next Steps

After F032:
- Monitor agent compliance
- Gather metrics on rule violations
- Iterate on enforcement mechanisms
