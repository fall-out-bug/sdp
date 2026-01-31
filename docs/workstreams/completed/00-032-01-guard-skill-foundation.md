---
assignee: null
completed: '2026-01-30'
depends_on: []
feature: F032
github_issue: null
project_id: 0
size: MEDIUM
status: completed
traceability:
- ac_description: Skill файл `.claude/skills/guard/SKILL.md` создан и валиден
  ac_id: AC1
  confidence: 1.0
  status: mapped
  test_file: tests/integration/test_guard_flow.py
  test_name: test_happy_path_activate_edit_complete
- ac_description: '`GuardSkill.check_edit(file)` возвращает `GuardResult` с `allowed:
    bool`'
  ac_id: AC2
  confidence: 1.0
  status: mapped
  test_file: tests/integration/test_guard_flow.py
  test_name: test_edit_blocked_without_active_ws
- ac_description: Если нет активного WS, `allowed=False` с сообщением "No active WS"
  ac_id: AC3
  confidence: 1.0
  status: mapped
  test_file: tests/integration/test_guard_flow.py
  test_name: test_edit_blocked_outside_scope
- ac_description: Unit тесты покрывают все edge cases (≥80%)
  ac_id: AC5
  confidence: 1.0
  status: mapped
  test_file: tests/unit/test_scope_manager.py
  test_name: test_is_in_scope_with_restricted_scope
- ac_description: Если файл не в scope, `allowed=False` с списком scope_files
  ac_id: AC4
  confidence: 1.0
  status: mapped
  test_file: tests/integration/test_guard_flow.py
  test_name: test_concurrent_activation_blocked
ws_id: 00-032-01
---

## 00-032-01: Guard Skill Foundation

### 🎯 Goal

**What must WORK after completing this WS:**
- Internal skill `@guard` блокирует edit без активного WS
- Guard проверяет что файл находится в scope текущего WS
- Guard интегрирован в workflow `@build`

**Acceptance Criteria:**
- [ ] AC1: Skill файл `.claude/skills/guard/SKILL.md` создан и валиден
- [ ] AC2: `GuardSkill.check_edit(file)` возвращает `GuardResult` с `allowed: bool`
- [ ] AC3: Если нет активного WS, `allowed=False` с сообщением "No active WS"
- [ ] AC4: Если файл не в scope, `allowed=False` с списком scope_files
- [ ] AC5: Unit тесты покрывают все edge cases (≥80%)

**⚠️ WS is NOT complete until Goal is achieved (all AC ✅).**

---

### Context

**Problem**: Агенты могут редактировать любые файлы без контроля. Нет проверки что работа идёт в рамках WS.

**Solution**: Pre-edit guard skill который:
1. Проверяет наличие активного WS
2. Проверяет что файл в scope этого WS
3. Блокирует edit если проверки не прошли

### Dependencies

None (independent, foundational WS)

### Input Files

- `.claude/skills/build/SKILL.md` (reference workflow)
- `src/sdp/beads/client.py` (Beads integration)
- `src/sdp/beads/models.py` (BeadsTask model)

### Steps

1. **Create GuardResult model**

   ```python
   # src/sdp/guard/models.py
   from dataclasses import dataclass
   
   @dataclass
   class GuardResult:
       """Result of pre-edit guard check."""
       allowed: bool
       ws_id: str | None
       reason: str
       scope_files: list[str]
   ```

2. **Create GuardSkill class**

   ```python
   # src/sdp/guard/skill.py
   from sdp.beads import BeadsClient
   from sdp.guard.models import GuardResult
   
   class GuardSkill:
       """Pre-edit guard for WS scope enforcement."""
       
       def __init__(self, beads_client: BeadsClient):
           self._client = beads_client
           self._active_ws: str | None = None
       
       def activate(self, ws_id: str) -> None:
           """Set active workstream."""
           ws = self._client.get_task(ws_id)
           if not ws:
               raise ValueError(f"WS not found: {ws_id}")
           self._active_ws = ws_id
       
       def check_edit(self, file_path: str) -> GuardResult:
           """Check if file edit is allowed."""
           if not self._active_ws:
               return GuardResult(
                   allowed=False,
                   ws_id=None,
                   reason="No active WS. Run @build <ws_id> first.",
                   scope_files=[]
               )
           
           ws = self._client.get_task(self._active_ws)
           scope = ws.sdp_metadata.get("scope_files", [])
           
           if not scope:
               # No scope defined = all files allowed
               return GuardResult(
                   allowed=True,
                   ws_id=self._active_ws,
                   reason="No scope restrictions",
                   scope_files=[]
               )
           
           if file_path not in scope:
               return GuardResult(
                   allowed=False,
                   ws_id=self._active_ws,
                   reason=f"File {file_path} not in WS scope",
                   scope_files=scope
               )
           
           return GuardResult(
               allowed=True,
               ws_id=self._active_ws,
               reason="File in scope",
               scope_files=scope
           )
   ```

3. **Create SKILL.md**

   ```markdown
   # .claude/skills/guard/SKILL.md
   ---
   name: guard
   description: Pre-edit gate enforcing WS scope (INTERNAL)
   tools: Read, Bash
   ---
   
   # @guard - Pre-Edit Gate (INTERNAL)
   
   **INTERNAL SKILL** — Called automatically before file edits.
   
   ## Purpose
   
   Enforce that all edits happen within active WS scope.
   
   ## Check Flow
   
   1. Is there an active WS? → No → BLOCK
   2. Is file in WS scope? → No → BLOCK  
   3. Allow edit
   
   ## CLI Integration
   
   ```bash
   # Activate WS (called by @build)
   sdp guard activate 00-032-01
   
   # Check file (called before edit)
   sdp guard check src/sdp/guard/skill.py
   ```
   ```

4. **Write unit tests**

   ```python
   # tests/unit/test_guard_skill.py
   import pytest
   from unittest.mock import Mock
   from sdp.guard.skill import GuardSkill
   from sdp.guard.models import GuardResult
   
   class TestGuardSkill:
       def test_no_active_ws_blocks_edit(self):
           """AC3: No active WS returns allowed=False."""
           client = Mock()
           guard = GuardSkill(client)
           
           result = guard.check_edit("any/file.py")
           
           assert result.allowed is False
           assert "No active WS" in result.reason
       
       def test_file_not_in_scope_blocks_edit(self):
           """AC4: File outside scope returns allowed=False."""
           client = Mock()
           client.get_task.return_value = Mock(
               sdp_metadata={"scope_files": ["src/allowed.py"]}
           )
           guard = GuardSkill(client)
           guard._active_ws = "00-032-01"
           
           result = guard.check_edit("src/forbidden.py")
           
           assert result.allowed is False
           assert "not in WS scope" in result.reason
           assert "src/allowed.py" in result.scope_files
       
       def test_file_in_scope_allows_edit(self):
           """File in scope returns allowed=True."""
           client = Mock()
           client.get_task.return_value = Mock(
               sdp_metadata={"scope_files": ["src/allowed.py"]}
           )
           guard = GuardSkill(client)
           guard._active_ws = "00-032-01"
           
           result = guard.check_edit("src/allowed.py")
           
           assert result.allowed is True
       
       def test_no_scope_allows_all(self):
           """No scope_files defined allows all edits."""
           client = Mock()
           client.get_task.return_value = Mock(sdp_metadata={})
           guard = GuardSkill(client)
           guard._active_ws = "00-032-01"
           
           result = guard.check_edit("any/file.py")
           
           assert result.allowed is True
   ```

### Output Files

- `.claude/skills/guard/SKILL.md`
- `src/sdp/guard/__init__.py`
- `src/sdp/guard/models.py`
- `src/sdp/guard/skill.py`
- `tests/unit/test_guard_skill.py`

### Completion Criteria

```bash
# Module imports
python -c "from sdp.guard import GuardSkill, GuardResult"

# Tests pass
pytest tests/unit/test_guard_skill.py -v

# Coverage ≥80%
pytest tests/unit/test_guard_skill.py --cov=src/sdp/guard --cov-fail-under=80

# Type check
mypy src/sdp/guard/ --strict
```

### Constraints

- DO NOT add external dependencies
- DO NOT persist state (WS-03 will add persistence)
- DO NOT integrate with CLI yet (WS-02 will add CLI)

---

## Execution Report

**Executed by:** ______  
**Date:** ______  
**Duration:** ______ minutes

### Goal Status
- [ ] AC1-AC5 — ✅

**Goal Achieved:** ______

### Files Changed
| File | Action | LOC |
|------|--------|-----|
| .claude/skills/guard/SKILL.md | Create | 50 |
| src/sdp/guard/__init__.py | Create | 10 |
| src/sdp/guard/models.py | Create | 20 |
| src/sdp/guard/skill.py | Create | 80 |
| tests/unit/test_guard_skill.py | Create | 100 |

### Commit
______
