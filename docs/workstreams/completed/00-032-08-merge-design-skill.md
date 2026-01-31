---
assignee: null
completed: '2026-01-30'
depends_on:
- 00-032-05
feature: F032
github_issue: null
project_id: 0
size: SMALL
status: completed
traceability:
- ac_description: '`SKILL.md` содержит ≤100 строк'
  ac_id: AC1
  confidence: 1.0
  status: mapped
  test_file: tests/integration/test_guard_flow.py
  test_name: test_happy_path_activate_edit_complete
- ac_description: Scope files автоматически определяются из плана
  ac_id: AC2
  confidence: 1.0
  status: mapped
  test_file: tests/integration/test_guard_flow.py
  test_name: test_edit_blocked_without_active_ws
- ac_description: Full spec в `docs/reference/design-spec.md`
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
- ac_description: Validation passes
  ac_id: AC4
  confidence: 1.0
  status: mapped
  test_file: tests/integration/test_guard_flow.py
  test_name: test_concurrent_activation_blocked
ws_id: 00-032-08
---

## 00-032-08: Merge design.md into SKILL.md

### 🎯 Goal

**What must WORK after completing this WS:**
- `.claude/skills/design/SKILL.md` сокращён до ≤100 строк
- Автоматическая генерация scope_files для WS
- Детали в `docs/reference/design-spec.md`

**Acceptance Criteria:**
- [ ] AC1: `SKILL.md` содержит ≤100 строк
- [ ] AC2: Scope files автоматически определяются из плана
- [ ] AC3: Full spec в `docs/reference/design-spec.md`
- [ ] AC4: Validation passes

**⚠️ WS is NOT complete until Goal is achieved (all AC ✅).**

---

### Context

**Problem**: Design skill 590 строк. Дублирует Beads workflow.

**Solution**: Сократить до core workflow, убрать дублирование.

### Dependencies

- **00-032-05**: Skill Template Standard

### Steps

1. **Rewrite SKILL.md (~80 lines)**

   ```markdown
   # .claude/skills/design/SKILL.md
   ---
   name: design
   description: Decompose feature into workstreams with scope
   tools: Read, Write, Bash, AskUserQuestion
   ---
   
   # @design - Feature Decomposition
   
   Analyze requirements and create workstreams with dependencies and scope.
   
   ## Quick Reference
   
   | Step | Action | Gate |
   |------|--------|------|
   | 1 | Read feature | Requirements clear |
   | 2 | Explore codebase | Context gathered |
   | 3 | Ask architecture | Decisions made |
   | 4 | Create workstreams | All WS have AC + scope |
   | 5 | Verify deps | No cycles |
   
   ## Workflow
   
   ### Step 1: Read Feature
   
   ```bash
   bd show {feature-id}
   # or
   Read("docs/drafts/{feature}.md")
   ```
   
   ### Step 2: Explore Codebase
   
   ```bash
   Glob("src/**/*.py")
   Grep("relevant patterns")
   ```
   
   ### Step 3: Architecture Questions
   
   Use AskUserQuestion for:
   - Complexity level (simple/medium/large)
   - Layers needed (domain/repo/service/api)
   - Database changes
   - External integrations
   
   ### Step 4: Create Workstreams
   
   For each WS:
   
   ```yaml
   ws_id: 00-032-01
   title: Domain entities
   size: MEDIUM
   depends_on: []
   scope_files:
     - src/domain/entities.py
     - src/domain/value_objects.py
     - tests/unit/test_entities.py
   acceptance_criteria:
     - AC1: Entity created with required fields
     - AC2: Value objects immutable
   ```
   
   **Key:** Include `scope_files` for guard enforcement.
   
   ### Step 5: Verify Dependencies
   
   ```bash
   sdp ws graph {feature-id}
   ```
   
   Check for cycles. Ensure topological order possible.
   
   ## Quality Gates
   
   See [Quality Gates Reference](../../docs/reference/quality-gates.md)
   
   ## Errors
   
   | Error | Cause | Fix |
   |-------|-------|-----|
   | Cycle detected | Circular deps | Break cycle |
   | Missing scope | No files listed | Add scope_files |
   | Too large | WS >500 LOC | Split WS |
   
   ## See Also
   
   - [Full Design Spec](../../docs/reference/design-spec.md)
   - [Sizing Guide](../../docs/reference/ws-sizing.md)
   ```

2. **Move details to docs/**

   `docs/reference/design-spec.md` contains:
   - Workstream sizing guidelines
   - Dependency patterns
   - Scope file detection heuristics

### Output Files

- `.claude/skills/design/SKILL.md` (rewritten)
- `docs/reference/design-spec.md`
- `docs/reference/ws-sizing.md`

### Completion Criteria

```bash
wc -l .claude/skills/design/SKILL.md
# ≤100

sdp skill validate .claude/skills/design/SKILL.md
```

---

## Execution Report

**Executed by:** ______  
**Date:** ______

### Goal Status
- [ ] AC1-AC4 — ✅

**Goal Achieved:** ______
