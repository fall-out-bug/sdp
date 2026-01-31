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
- ac_description: Полная спецификация в `docs/reference/build-spec.md`
  ac_id: AC2
  confidence: 1.0
  status: mapped
  test_file: tests/integration/test_guard_flow.py
  test_name: test_edit_blocked_without_active_ws
- ac_description: Примеры в `docs/examples/build/`
  ac_id: AC3
  confidence: 1.0
  status: mapped
  test_file: tests/integration/test_guard_flow.py
  test_name: test_edit_blocked_outside_scope
- ac_description: Guard activation добавлена как Step 1
  ac_id: AC5
  confidence: 1.0
  status: mapped
  test_file: tests/unit/test_scope_manager.py
  test_name: test_is_in_scope_with_restricted_scope
- ac_description: TDD workflow ссылается на `@tdd` skill
  ac_id: AC4
  confidence: 1.0
  status: mapped
  test_file: tests/integration/test_guard_flow.py
  test_name: test_concurrent_activation_blocked
- ac_description: '`sdp skill validate .claude/skills/build/SKILL.md` passes'
  ac_id: AC6
  confidence: 1.0
  status: mapped
  test_file: tests/unit/test_skill_validator.py
  test_name: test_validate_skill_valid
ws_id: 00-032-06
---

## 00-032-06: Merge build.md into SKILL.md

### 🎯 Goal

**What must WORK after completing this WS:**
- `.claude/skills/build/SKILL.md` сокращён до ≤100 строк
- Детальная спецификация перенесена в `docs/reference/build-spec.md`
- Guard integration добавлена в workflow

**Acceptance Criteria:**
- [ ] AC1: `SKILL.md` содержит ≤100 строк
- [ ] AC2: Полная спецификация в `docs/reference/build-spec.md`
- [ ] AC3: Примеры в `docs/examples/build/`
- [ ] AC4: TDD workflow ссылается на `@tdd` skill
- [ ] AC5: Guard activation добавлена как Step 1
- [ ] AC6: `sdp skill validate .claude/skills/build/SKILL.md` passes

**⚠️ WS is NOT complete until Goal is achieved (all AC ✅).**

---

### Context

**Problem**: Текущий `build/SKILL.md` + `prompts/commands/build.md` = 500+ строк.

**Solution**: Сократить до essential workflow, детали в docs/.

### Dependencies

- **00-032-05**: Skill Template Standard (provides template)

### Steps

1. **Rewrite SKILL.md (target: ~80 lines)**

   ```markdown
   # .claude/skills/build/SKILL.md
   ---
   name: build
   description: Execute workstream with TDD and guard enforcement
   tools: Read, Write, Edit, Bash, Skill
   ---
   
   # @build - Execute Workstream
   
   Execute a single workstream following TDD discipline with automatic guard.
   
   ## Quick Reference
   
   | Step | Action | Gate |
   |------|--------|------|
   | 1 | Activate guard | `sdp guard activate` succeeds |
   | 2 | Read WS spec | AC present and clear |
   | 3 | TDD cycle | `@tdd` for each AC |
   | 4 | Quality check | `sdp quality check` passes |
   | 5 | Complete | Commit with WS-ID |
   
   ## Workflow
   
   ### Step 1: Activate Guard
   
   ```bash
   sdp guard activate {WS-ID}
   ```
   
   **Gate:** Must succeed. If fails, WS not ready.
   
   ### Step 2: Read Workstream
   
   ```bash
   Read("docs/workstreams/backlog/{WS-ID}-*.md")
   ```
   
   Extract:
   - Goal and Acceptance Criteria
   - Input/Output files
   - Steps to execute
   
   ### Step 3: TDD Cycle
   
   For each AC, call internal TDD skill:
   
   ```
   @tdd "AC1: {description}"
   ```
   
   Cycle: Red → Green → Refactor
   
   ### Step 4: Quality Check
   
   ```bash
   sdp quality check --module {module}
   ```
   
   Must pass:
   - Coverage ≥80%
   - mypy --strict
   - ruff (no errors)
   - Files <200 LOC
   
   ### Step 5: Complete
   
   ```bash
   sdp guard complete {WS-ID}
   git add .
   git commit -m "feat({scope}): {WS-ID} - {title}"
   ```
   
   ## Quality Gates
   
   See [Quality Gates Reference](../../docs/reference/quality-gates.md)
   
   ## Errors
   
   | Error | Cause | Fix |
   |-------|-------|-----|
   | No active WS | Guard not activated | `sdp guard activate` |
   | File not in scope | Editing wrong file | Check WS scope |
   | Coverage <80% | Missing tests | Add tests |
   
   ## See Also
   
   - [Full Build Spec](../../docs/reference/build-spec.md)
   - [TDD Skill](../tdd/SKILL.md)
   - [Guard Skill](../guard/SKILL.md)
   ```

2. **Create full spec in docs/**

   ```markdown
   # docs/reference/build-spec.md
   # Build Command Full Specification
   
   This document contains the complete specification for `@build`.
   For quick reference, see [SKILL.md](../../.claude/skills/build/SKILL.md).
   
   ## Overview
   
   The build command executes a single workstream following TDD.
   
   ## Prerequisites
   
   - Workstream file exists in `docs/workstreams/backlog/`
   - WS has Goal and Acceptance Criteria
   - Dependencies are satisfied
   - Scope is SMALL or MEDIUM
   
   ## Detailed Steps
   
   ### Pre-Build Validation
   
   Run automatically or manually:
   
   ```bash
   hooks/pre-build.sh {WS-ID}
   ```
   
   Checks:
   - WS file exists
   - Goal section present
   - AC defined
   - Dependencies complete
   
   [... rest of detailed content from original build.md ...]
   ```

3. **Create examples**

   ```markdown
   # docs/examples/build/tdd-cycle.md
   # TDD Cycle Example
   
   ## Red Phase
   
   ```python
   # Write test FIRST
   def test_user_can_login():
       result = auth_service.login("user@example.com", "password")
       assert result.success is True
   ```
   
   ```bash
   pytest tests/unit/test_auth.py::test_user_can_login -v
   # Expected: FAILED (no implementation yet)
   ```
   
   ## Green Phase
   
   ```python
   # Minimal implementation
   def login(self, email: str, password: str) -> LoginResult:
       return LoginResult(success=True)
   ```
   
   ```bash
   pytest tests/unit/test_auth.py::test_user_can_login -v
   # Expected: PASSED
   ```
   
   ## Refactor Phase
   
   ```python
   # Improve implementation
   def login(self, email: str, password: str) -> LoginResult:
       user = self.user_repo.find_by_email(email)
       if not user or not user.verify_password(password):
           return LoginResult(success=False, error="Invalid credentials")
       return LoginResult(success=True, user=user)
   ```
   ```

4. **Validate new skill**

   ```bash
   sdp skill validate .claude/skills/build/SKILL.md
   ```

### Output Files

- `.claude/skills/build/SKILL.md` (rewritten, ~80 lines)
- `docs/reference/build-spec.md` (full spec, ~400 lines)
- `docs/examples/build/tdd-cycle.md`
- `docs/examples/build/execution-report.md`

### Completion Criteria

```bash
# Skill is short enough
wc -l .claude/skills/build/SKILL.md
# Expected: ≤100

# Validation passes
sdp skill validate .claude/skills/build/SKILL.md

# References exist
test -f docs/reference/build-spec.md
test -f docs/examples/build/tdd-cycle.md
```

---

## Execution Report

**Executed by:** ______  
**Date:** ______

### Goal Status
- [ ] AC1-AC6 — ✅

**Goal Achieved:** ______
