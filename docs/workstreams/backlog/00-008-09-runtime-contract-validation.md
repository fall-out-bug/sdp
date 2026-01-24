---
ws_id: 00-410-09
project_id: 00
feature: F008
status: backlog
size: SMALL
github_issue: null
assignee: null
started: null
completed: null
blocked_reason: null
ws_version: "2.0"
capability_tier: T2
---

## 00-410-09: Runtime Contract Validation

### 🎯 Цель (Goal)

**Что должно РАБОТАТЬ после завершения WS:**
- Автоматическая проверка что T2/T3 builds не изменяют Interface или Tests sections
- Runtime validation перед/после /build execution
- Build fails with clear error если contract нарушен

**Acceptance Criteria (критерии приёмки):**
- [ ] AC1: Snapshot Interface + Tests sections перед /build
- [ ] AC2: Compare snapshot с post-build state
- [ ] AC3: Raise `ContractViolationError` если T2/T3 модифицировал contract
- [ ] AC4: Unit-тесты для violation detection (changed signature, added test, removed test)

**⚠️ Правило:** WS НЕ завершён, пока Goal не достигнута (все AC ✅).

---

### Контекст

**Проблема:** Contract-Driven workflow (WS-410-04) определяет, что T2/T3 workstreams должны сохранять Interface и Tests read-only:

```
#### Interface (DO NOT MODIFY для T2/T3)
def function_name(arg: Type) -> ReturnType:
    """Docstring."""
    raise NotImplementedError

#### Tests (DO NOT MODIFY для T2/T3)
def test_function():
    assert function_name(...) == expected
```

Но enforcement только на уровне документации. **Нет runtime check.**

**Последствия:**
- T2/T3 agent может случайно изменить signature
- Тесты могут быть удалены или ослаблены
- Contract integrity не гарантирована

**Решение:** Runtime validation для T2/T3 builds.

**Приоритет:** LOW - safety net, но не критично если процесс работает правильно.

---

### Зависимость

- WS-410-04 (/test command) — обязательная (contract definition)
- WS-410-05 (Builder router) — обязательная (build execution)
- WS-410-02 (Validator) — желательная (tier check)

---

### Входные файлы

- `sdp/src/sdp/core/builder_router.py` — где происходит build
- `sdp/src/sdp/core/workstream.py` — Workstream dataclass
- WS markdown files — для snapshot Interface/Tests sections

---

### Contract (для T2 — Contract-Driven WS v2.0)

#### Input Files (read-only)
- `sdp/src/sdp/core/builder_router.py` — build execution point

#### Output Files (create/modify)
- `sdp/src/sdp/core/contract_validator.py` — новый модуль
- `sdp/tests/unit/core/test_contract_validator.py` — тесты

#### Interface (DO NOT MODIFY для T2)

```python
# sdp/src/sdp/core/contract_validator.py

from dataclasses import dataclass
from pathlib import Path
from typing import Optional

class ContractViolationError(Exception):
    """Raised when T2/T3 build violates contract.

    Args:
        ws_id: Workstream ID
        tier: Capability tier (T2 or T3)
        violation: Description of what was changed
    """

    def __init__(self, ws_id: str, tier: str, violation: str) -> None:
        """Initialize contract violation error."""
        self.ws_id = ws_id
        self.tier = tier
        self.violation = violation
        super().__init__(
            f"Contract violation in {ws_id} (tier {tier}): {violation}"
        )


@dataclass
class ContractSnapshot:
    """Snapshot of contract sections (Interface + Tests).

    Args:
        interface_content: Content of Interface section
        tests_content: Content of Tests section
    """
    interface_content: str
    tests_content: str

    def equals(self, other: "ContractSnapshot") -> bool:
        """Check if two snapshots are identical.

        Args:
            other: Another snapshot to compare

        Returns:
            True if interface and tests are identical
        """
        raise NotImplementedError


class ContractValidator:
    """Validate contract integrity for T2/T3 builds.

    Raises:
        ContractViolationError: If contract is modified during build
    """

    def snapshot_contract(self, ws_file: Path) -> Optional[ContractSnapshot]:
        """Extract Interface + Tests sections from WS file.

        Args:
            ws_file: Path to workstream markdown file

        Returns:
            ContractSnapshot if contract sections exist, None otherwise
        """
        raise NotImplementedError

    def validate_contract_integrity(
        self,
        before: ContractSnapshot,
        after: ContractSnapshot,
        ws_id: str,
        tier: str
    ) -> None:
        """Validate contract wasn't modified.

        Args:
            before: Snapshot before build
            after: Snapshot after build
            ws_id: Workstream ID
            tier: Capability tier

        Raises:
            ContractViolationError: If contract was modified
        """
        raise NotImplementedError
```

#### Tests (DO NOT MODIFY для T2)

```python
# sdp/tests/unit/core/test_contract_validator.py

def test_snapshot_extracts_interface_and_tests():
    """Must extract both sections from WS file."""
    ws_content = """
    #### Interface (DO NOT MODIFY для T2/T3)
    def foo() -> int:
        pass

    #### Tests (DO NOT MODIFY для T2/T3)
    def test_foo():
        assert foo() == 42
    """

    validator = ContractValidator()
    snapshot = validator.snapshot_contract(ws_content)

    assert "def foo() -> int:" in snapshot.interface_content
    assert "def test_foo():" in snapshot.tests_content


def test_validate_detects_interface_change():
    """Must detect interface signature change."""
    before = ContractSnapshot(
        interface_content="def foo(x: int) -> int:",
        tests_content="def test_foo(): pass"
    )

    after = ContractSnapshot(
        interface_content="def foo(x: str) -> int:",  # Changed type
        tests_content="def test_foo(): pass"
    )

    validator = ContractValidator()

    with pytest.raises(ContractViolationError, match="Interface modified"):
        validator.validate_contract_integrity(before, after, "WS-410-01", "T2")


def test_validate_detects_test_removal():
    """Must detect removed test."""
    before = ContractSnapshot(
        interface_content="def foo() -> int:",
        tests_content="def test_foo(): pass\ndef test_bar(): pass"
    )

    after = ContractSnapshot(
        interface_content="def foo() -> int:",
        tests_content="def test_foo(): pass"  # test_bar removed
    )

    validator = ContractValidator()

    with pytest.raises(ContractViolationError, match="Tests modified"):
        validator.validate_contract_integrity(before, after, "WS-410-01", "T2")


def test_validate_passes_for_unchanged_contract():
    """Must pass when contract unchanged."""
    snapshot = ContractSnapshot(
        interface_content="def foo() -> int:",
        tests_content="def test_foo(): pass"
    )

    validator = ContractValidator()

    # Should not raise
    validator.validate_contract_integrity(snapshot, snapshot, "WS-410-01", "T2")


def test_snapshot_returns_none_for_no_contract():
    """Must return None if no contract sections."""
    ws_content = """
    ### Some other section
    No contract here
    """

    validator = ContractValidator()
    snapshot = validator.snapshot_contract(ws_content)

    assert snapshot is None
```

**⚠️ Правило для T2:** Секции Interface и Tests являются **READ-ONLY**. Запрещено изменять сигнатуры функций, docstrings, тесты. Только реализация тел функций.

---

### Шаги

1. **Реализовать `ContractValidator`** согласно интерфейсу:
   - Парсинг WS markdown для извлечения Interface/Tests
   - Diff comparison для обнаружения изменений
   - Raise `ContractViolationError` при нарушениях

2. **Интегрировать с `BuilderRouter`**:
   ```python
   def execute_build(ws: Workstream) -> BuildResult:
       # Only for T2/T3
       if ws.capability_tier in ("T2", "T3"):
           before = contract_validator.snapshot_contract(ws.file_path)

       result = ... # Execute build

       if ws.capability_tier in ("T2", "T3") and before:
           after = contract_validator.snapshot_contract(ws.file_path)
           contract_validator.validate_contract_integrity(
               before, after, ws.ws_id, ws.capability_tier
           )

       return result
   ```

3. **Добавить unit-тесты** согласно контракту

4. **Error handling**:
   - Clear error message о том, что именно изменилось
   - Rollback механизм (optional) для отката изменений

---

### Ожидаемый результат

- T2/T3 builds автоматически проверяются на contract integrity
- Build fails если contract нарушен
- Clear error message указывает на нарушение

---

### Scope Estimate

- Файлов: ~2 создано
- Строк: ~200-350 (SMALL)
- Токенов: ~1200-2000

---

### Критерий завершения

```bash
# Unit tests pass
pytest sdp/tests/unit/core/test_contract_validator.py -v

# Integration test: T2 build that violates contract fails
# (manual test or add integration test)

# Type checks
mypy sdp/src/sdp/core/contract_validator.py

# Lint
ruff check sdp/src/sdp/core/contract_validator.py
```

---

### Ограничения

- НЕ изменять WS markdown format (должен оставаться human-readable)
- НЕ блокировать T0/T1 builds (validation только для T2/T3)
- Diff должен игнорировать whitespace/formatting changes

---

### Related Issues

- Referenced in: `docs/reviews/F410-cross-ws-review.md` (Low Priority recommendation)
- Depends on: WS-410-04 (/test), WS-410-05 (Builder router)
- Related to: WS-410-02 (Capability tier validator)

---

### Notes

**Alternative approach:** Instead of runtime validation, enforce at commit time via pre-commit hook. Trade-offs:
- **Pro:** Catches violations earlier
- **Con:** Requires git integration, harder to test
- **Decision:** Runtime validation is simpler and easier to test.
