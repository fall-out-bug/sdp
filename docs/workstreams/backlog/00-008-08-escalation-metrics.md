---
ws_id: 00-410-08
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

## 00-410-08: Escalation Metrics & Monitoring

### 🎯 Цель (Goal)

**Что должно РАБОТАТЬ после завершения WS:**
- Трекинг и анализ escalation patterns для T2/T3 workstreams
- Dashboard с ключевыми метриками (escalation rate, avg attempts, cost)
- Alerts при аномально высоком escalation rate

**Acceptance Criteria (критерии приёмки):**
- [ ] AC1: Метрики собираются при каждом escalation event (WS ID, tier, attempts, diagnostics)
- [ ] AC2: CLI команда `sdp metrics escalations` выводит summary
- [ ] AC3: Реализован threshold-based alert (email/Slack при escalation rate > 20%)
- [ ] AC4: Unit-тесты для metrics collection и alert logic

**⚠️ Правило:** WS НЕ завершён, пока Goal не достигнута (все AC ✅).

---

### Контекст

**Проблема:** Политика D1 (WS-410-05) определяет retry limit для T2/T3:
```
T2/T3: 3 attempts → escalate to human with diagnostics
```

Но нет visibility:
- Сколько WS escalate в production?
- Какие тиры/features escalate чаще?
- Какова средняя стоимость escalation (human time)?
- Есть ли паттерны в failure diagnostics?

**Решение:** Metrics pipeline для escalation events с анализом и alerting.

**Приоритет:** LOW - nice-to-have для observability, но не критично.

---

### Зависимость

- WS-410-05 (Builder router) — обязательная (`HumanEscalationError`)
- WS-410-07 (Tier metrics) — желательная (можно переиспользовать storage)

---

### Входные файлы

- `sdp/src/sdp/core/builder_router.py` — `HumanEscalationError` raising
- `sdp/src/sdp/cli.py` — добавить команду `sdp metrics escalations`
- `.sdp/escalation_metrics.json` — новое хранилище метрик

---

### Contract (для T2 — Contract-Driven WS v2.0)

#### Input Files (read-only)
- `sdp/src/sdp/core/builder_router.py` — где происходит escalation

#### Output Files (create/modify)
- `sdp/src/sdp/core/escalation_metrics.py` — новый модуль
- `sdp/src/sdp/cli.py` — добавить команду `metrics escalations`
- `.sdp/escalation_metrics.json` — метрики storage

#### Interface (DO NOT MODIFY для T2)

```python
# sdp/src/sdp/core/escalation_metrics.py

from dataclasses import dataclass
from datetime import datetime
from typing import Optional

@dataclass
class EscalationEvent:
    """Record of human escalation event.

    Args:
        ws_id: Workstream ID that escalated
        tier: Capability tier (T2 or T3)
        attempts: Number of failed attempts before escalation
        timestamp: When escalation occurred
        diagnostics: Diagnostic info for human
        feature_id: Optional feature ID
    """
    ws_id: str
    tier: str
    attempts: int
    timestamp: datetime
    diagnostics: str
    feature_id: Optional[str] = None

    def __post_init__(self) -> None:
        """Validate escalation event fields."""
        if self.tier not in ("T2", "T3"):
            raise ValueError(f"Invalid tier for escalation: {self.tier}")
        if self.attempts <= 0:
            raise ValueError(f"Attempts must be positive: {self.attempts}")


class EscalationMetricsStore:
    """Store and analyze escalation metrics.

    Raises:
        FileNotFoundError: If storage path doesn't exist and create=False
    """

    def record_escalation(self, event: EscalationEvent) -> None:
        """Record escalation event.

        Args:
            event: Escalation event to record
        """
        raise NotImplementedError

    def get_escalation_rate(
        self,
        tier: Optional[str] = None,
        days: int = 7
    ) -> float:
        """Calculate escalation rate.

        Args:
            tier: Filter by tier (None = all tiers)
            days: Time window in days

        Returns:
            Escalation rate as fraction (0.0 - 1.0)
        """
        raise NotImplementedError

    def get_top_escalating_ws(self, limit: int = 10) -> list[tuple[str, int]]:
        """Get workstreams with most escalations.

        Args:
            limit: Max number of results

        Returns:
            List of (ws_id, escalation_count) tuples
        """
        raise NotImplementedError
```

#### Tests (DO NOT MODIFY для T2)

```python
# sdp/tests/unit/core/test_escalation_metrics.py

def test_escalation_event_validates_tier():
    """Must reject invalid tiers."""
    with pytest.raises(ValueError, match="Invalid tier"):
        EscalationEvent(
            ws_id="WS-410-01",
            tier="T0",  # Invalid for escalation
            attempts=3,
            timestamp=datetime.now(),
            diagnostics="Test"
        )


def test_escalation_event_validates_attempts():
    """Must reject non-positive attempts."""
    with pytest.raises(ValueError, match="must be positive"):
        EscalationEvent(
            ws_id="WS-410-01",
            tier="T2",
            attempts=0,  # Invalid
            timestamp=datetime.now(),
            diagnostics="Test"
        )


def test_record_escalation():
    """Must record escalation event."""
    store = EscalationMetricsStore()
    event = EscalationEvent(
        ws_id="WS-410-01",
        tier="T2",
        attempts=3,
        timestamp=datetime.now(),
        diagnostics="Build failed: syntax error"
    )

    store.record_escalation(event)

    # Verify persisted
    events = store._load_events()
    assert len(events) == 1
    assert events[0].ws_id == "WS-410-01"


def test_get_escalation_rate():
    """Must calculate escalation rate correctly."""
    store = EscalationMetricsStore()

    # Record 3 T2 escalations, 1 T3 escalation
    for i in range(3):
        store.record_escalation(EscalationEvent(
            ws_id=f"WS-{i}",
            tier="T2",
            attempts=3,
            timestamp=datetime.now(),
            diagnostics="Test"
        ))

    store.record_escalation(EscalationEvent(
        ws_id="WS-3",
        tier="T3",
        attempts=3,
        timestamp=datetime.now(),
        diagnostics="Test"
    ))

    # Assume 20 total builds in period
    rate = store.get_escalation_rate(days=7)
    assert rate == 4 / 20  # 4 escalations out of 20 builds


def test_get_top_escalating_ws():
    """Must return most escalating workstreams."""
    store = EscalationMetricsStore()

    # WS-A: 3 escalations, WS-B: 1 escalation
    for _ in range(3):
        store.record_escalation(EscalationEvent(
            ws_id="WS-A",
            tier="T2",
            attempts=3,
            timestamp=datetime.now(),
            diagnostics="Test"
        ))

    store.record_escalation(EscalationEvent(
        ws_id="WS-B",
        tier="T2",
        attempts=3,
        timestamp=datetime.now(),
        diagnostics="Test"
    ))

    top = store.get_top_escalating_ws(limit=2)
    assert top == [("WS-A", 3), ("WS-B", 1)]
```

**⚠️ Правило для T2:** Секции Interface и Tests являются **READ-ONLY**. Запрещено изменять сигнатуры функций, docstrings, тесты. Только реализация тел функций.

---

### Шаги

1. **Реализовать `EscalationMetricsStore`** согласно интерфейсу
2. **Интегрировать с `BuilderRouter`**:
   ```python
   # При raise HumanEscalationError
   escalation_metrics.record_escalation(EscalationEvent(...))
   ```
3. **Добавить CLI команды**:
   ```bash
   # Show escalation summary
   sdp metrics escalations

   # Filter by tier
   sdp metrics escalations --tier T2

   # Show top escalating WS
   sdp metrics escalations --top 10
   ```
4. **Реализовать alerting** (optional):
   ```python
   if escalation_rate > ALERT_THRESHOLD:
       send_alert(f"High escalation rate: {escalation_rate:.1%}")
   ```

---

### Ожидаемый результат

- Escalation events tracked в `.sdp/escalation_metrics.json`
- CLI dashboard для viewing metrics
- Alerts при high escalation rate (optional)

---

### Scope Estimate

- Файлов: ~2 создано + ~1 изменено
- Строк: ~200-350 (SMALL)
- Токенов: ~1200-2000

---

### Критерий завершения

```bash
# Unit tests pass
pytest sdp/tests/unit/core/test_escalation_metrics.py -v

# CLI works
sdp metrics escalations
sdp metrics escalations --tier T2

# Type checks
mypy sdp/src/sdp/core/escalation_metrics.py

# Lint
ruff check sdp/src/sdp/core/escalation_metrics.py
```

---

### Ограничения

- НЕ изменять `HumanEscalationError` signature
- НЕ удалять существующие metrics (append-only)
- Storage должен быть thread-safe

---

### Related Issues

- Referenced in: `docs/reviews/F410-cross-ws-review.md` (Low Priority recommendation)
- Depends on: WS-410-05 (Builder router)
- Related to: WS-410-07 (Tier metrics — similar storage pattern)
