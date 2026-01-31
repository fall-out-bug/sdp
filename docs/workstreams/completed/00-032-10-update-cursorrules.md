---
assignee: null
completed: '2026-01-30'
depends_on:
- 00-032-09
feature: F032
github_issue: null
project_id: 0
size: SMALL
status: completed
traceability:
- ac_description: '`.cursorrules` содержит ≤50 строк'
  ac_id: AC1
  confidence: 1.0
  status: mapped
  test_file: tests/integration/test_guard_flow.py
  test_name: test_happy_path_activate_edit_complete
- ac_description: Ссылки на PROTOCOL.md и skills
  ac_id: AC2
  confidence: 1.0
  status: mapped
  test_file: tests/integration/test_guard_flow.py
  test_name: test_edit_blocked_without_active_ws
- ac_description: Только критичные правила inline (forbidden patterns)
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
- ac_description: Guard enforcement упомянут
  ac_id: AC4
  confidence: 1.0
  status: mapped
  test_file: tests/integration/test_guard_flow.py
  test_name: test_concurrent_activation_blocked
ws_id: 00-032-10
---

## 00-032-10: Update .cursorrules

### 🎯 Goal

**What must WORK after completing this WS:**
- `.cursorrules` сокращён до ≤50 строк
- Только критичные правила inline
- Остальное — ссылки на PROTOCOL.md и skills

**Acceptance Criteria:**
- [ ] AC1: `.cursorrules` содержит ≤50 строк
- [ ] AC2: Ссылки на PROTOCOL.md и skills
- [ ] AC3: Только критичные правила inline (forbidden patterns)
- [ ] AC4: Guard enforcement упомянут

**⚠️ WS is NOT complete until Goal is achieved (all AC ✅).**

---

### Context

**Problem**: Текущий `.cursorrules` 233 строки. Дублирует PROTOCOL.md.

**Solution**: Сократить до pointer с критичными правилами.

### Dependencies

- **00-032-09**: Delete prompts/commands/

### Steps

1. **Rewrite .cursorrules (~45 lines)**

   ```markdown
   # SDP Project Rules
   
   This project uses **Spec-Driven Protocol (SDP)** v0.6.0.
   
   ## Commands
   
   Use skills for all work:
   - `@idea` — Gather requirements
   - `@design` — Plan workstreams
   - `@build` — Execute workstream (guard enforced)
   - `@review` — Quality review
   - `@deploy` — Production deployment
   
   ## Guard Enforcement
   
   All edits require active workstream:
   
   ```bash
   sdp guard activate {WS-ID}  # Before editing
   sdp guard check {file}      # Verify allowed
   ```
   
   ## Critical Rules
   
   **Forbidden:**
   - ❌ `except: pass`
   - ❌ Files > 200 LOC
   - ❌ TODO without WS
   - ❌ Edit without active WS
   
   **Required:**
   - ✅ TDD (Red → Green → Refactor)
   - ✅ Coverage ≥80%
   - ✅ Type hints (mypy --strict)
   - ✅ Conventional commits
   
   ## Documentation
   
   - [PROTOCOL.md](PROTOCOL.md) — Full specification
   - [Skills](.claude/skills/) — Command details
   - [Quality Gates](docs/reference/quality-gates.md)
   
   **Version:** 0.6.0
   ```

### Output Files

- `.cursorrules` (rewritten)

### Completion Criteria

```bash
# Check line count
wc -l .cursorrules
# Expected: ≤50

# References exist
grep "PROTOCOL.md" .cursorrules
grep "skills" .cursorrules
```

---

## Execution Report

**Executed by:** ______  
**Date:** ______

### Goal Status
- [ ] AC1-AC4 — ✅

**Goal Achieved:** ______
