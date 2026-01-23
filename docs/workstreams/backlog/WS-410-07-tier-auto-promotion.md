---
ws_id: WS-410-07
feature: F410
status: backlog
size: MEDIUM
github_issue: null
assignee: null
started: null
completed: null
blocked_reason: null
ws_version: "2.0"
capability_tier: T1
---

## WS-410-07: Tier Auto-Promotion System

### 🎯 Цель (Goal)

**Что должно РАБОТАТЬ после завершения WS:**
- Система автоматического повышения/понижения capability tier на основе метрик успешности
- Трекинг success rate для каждого workstream + tier combination
- Автоматическое обновление `capability_tier` в WS frontmatter при достижении thresholds

**Acceptance Criteria (критерии приёмки):**
- [ ] AC1: Создана база метрик для хранения execution history (WS ID → attempts, successes, tier)
- [ ] AC2: Реализована логика promotion/demotion с configurable thresholds
- [ ] AC3: Автоматическое обновление `capability_tier` в WS файлах при promotion
- [ ] AC4: Добавлена CLI команда `sdp tier promote-check` для manual trigger
- [ ] AC5: Unit-тесты для promotion logic и threshold checks

**⚠️ Правило:** WS НЕ завершён, пока Goal не достигнута (все AC ✅).

---

### Контекст

**Проблема:** Capability tiers (T0-T3) в текущей реализации F410 статичны и устанавливаются вручную при создании WS. Это приводит к:
- Недоиспользованию автоматизации (workstreams остаются T3 даже после многих успехов)
- Ручному overhead для обновления tiers
- Отсутствию обратной связи о качестве workstream decomposition

**Решение:** Автоматическая promotion/demotion на основе success metrics:

```
T3 (Junior) ──10 successes──→ T2 (Mid-level) ──20 successes──→ T1 (Senior)
     ↑                              ↑                               ↓
     └──────────────── 3+ consecutive failures ───────────────────┘
```

**Приоритет:** MEDIUM - улучшает UX, но не критично для MVP.

---

### Зависимость

- WS-410-01 (Template) — обязательная (capability_tier field)
- WS-410-02 (Validator) — обязательная (tier validation)
- WS-410-05 (Builder router) — обязательная (tier usage)

---

### Входные файлы

- `sdp/src/sdp/core/workstream.py` — Workstream dataclass
- `sdp/src/sdp/core/capability_tier_validator.py` — валидация tiers
- `docs/workstreams/backlog/*.md` — WS файлы для обновления
- `sdp/src/sdp/cli.py` — добавить команду `tier promote-check`

---

### Шаги

1. **Создать метрики хранилище** (`tier_metrics.py`):
   ```python
   @dataclass
   class TierMetrics:
       ws_id: str
       current_tier: str
       total_attempts: int
       successful_attempts: int
       consecutive_failures: int
       last_updated: datetime

   class TierMetricsStore:
       def record_attempt(ws_id: str, success: bool) -> None
       def get_metrics(ws_id: str) -> TierMetrics
       def check_promotion_eligible(ws_id: str) -> Optional[str]  # Returns new tier or None
   ```

2. **Реализовать promotion logic** (`tier_promoter.py`):
   ```python
   PROMOTION_RULES = {
       "T3->T2": {"min_successes": 10, "min_success_rate": 0.8},
       "T2->T1": {"min_successes": 20, "min_success_rate": 0.85},
   }

   DEMOTION_RULES = {
       "consecutive_failures": 3,
   }

   def check_promotion(metrics: TierMetrics) -> Optional[str]:
       # Check if eligible for tier upgrade
       pass

   def check_demotion(metrics: TierMetrics) -> Optional[str]:
       # Check if should be demoted
       pass
   ```

3. **Интегрировать с BuilderRouter**:
   ```python
   # В builder_router.py после каждого build execution
   def execute_build(ws: Workstream, ...) -> BuildResult:
       result = ...  # execute build

       # Record metrics
       tier_metrics.record_attempt(ws.ws_id, success=result.success)

       # Check promotion
       new_tier = tier_metrics.check_promotion_eligible(ws.ws_id)
       if new_tier:
           update_workstream_tier(ws.file_path, new_tier)
           logger.info(f"{ws.ws_id}: promoted {ws.capability_tier} → {new_tier}")

       return result
   ```

4. **Добавить CLI команду**:
   ```bash
   # Manually trigger promotion check for all WS
   sdp tier promote-check

   # Check specific WS
   sdp tier promote-check WS-410-01

   # Show metrics
   sdp tier metrics WS-410-01
   ```

5. **Добавить unit-тесты**:
   - Promotion after 10 T3 successes
   - Demotion after 3 consecutive failures
   - Edge cases (exactly threshold, rate below threshold)
   - File update correctness

---

### Код

```python
# sdp/src/sdp/core/tier_metrics.py

from dataclasses import dataclass, field
from datetime import datetime
from pathlib import Path
from typing import Optional
import json

@dataclass
class TierMetrics:
    """Metrics for tier promotion/demotion."""
    ws_id: str
    current_tier: str
    total_attempts: int = 0
    successful_attempts: int = 0
    consecutive_failures: int = 0
    last_updated: datetime = field(default_factory=datetime.now)

    @property
    def success_rate(self) -> float:
        """Calculate success rate (0.0 - 1.0)."""
        if self.total_attempts == 0:
            return 0.0
        return self.successful_attempts / self.total_attempts


class TierMetricsStore:
    """Store and retrieve tier metrics."""

    def __init__(self, storage_path: Path = Path(".sdp/tier_metrics.json")):
        """Initialize metrics store.

        Args:
            storage_path: Path to JSON file for persistence
        """
        self.storage_path = storage_path
        self.storage_path.parent.mkdir(parents=True, exist_ok=True)
        self._load()

    def _load(self) -> None:
        """Load metrics from storage."""
        if not self.storage_path.exists():
            self._metrics = {}
            return

        with open(self.storage_path) as f:
            data = json.load(f)
            self._metrics = {
                ws_id: TierMetrics(**m) for ws_id, m in data.items()
            }

    def _save(self) -> None:
        """Save metrics to storage."""
        data = {
            ws_id: {
                "ws_id": m.ws_id,
                "current_tier": m.current_tier,
                "total_attempts": m.total_attempts,
                "successful_attempts": m.successful_attempts,
                "consecutive_failures": m.consecutive_failures,
                "last_updated": m.last_updated.isoformat(),
            }
            for ws_id, m in self._metrics.items()
        }
        with open(self.storage_path, "w") as f:
            json.dump(data, f, indent=2)

    def record_attempt(self, ws_id: str, tier: str, success: bool) -> None:
        """Record a build attempt.

        Args:
            ws_id: Workstream ID
            tier: Current capability tier
            success: Whether attempt succeeded
        """
        if ws_id not in self._metrics:
            self._metrics[ws_id] = TierMetrics(ws_id=ws_id, current_tier=tier)

        metrics = self._metrics[ws_id]
        metrics.total_attempts += 1
        if success:
            metrics.successful_attempts += 1
            metrics.consecutive_failures = 0
        else:
            metrics.consecutive_failures += 1
        metrics.last_updated = datetime.now()

        self._save()

    def get_metrics(self, ws_id: str) -> Optional[TierMetrics]:
        """Get metrics for workstream.

        Args:
            ws_id: Workstream ID

        Returns:
            TierMetrics if exists, None otherwise
        """
        return self._metrics.get(ws_id)


# sdp/src/sdp/core/tier_promoter.py

from typing import Optional
from sdp.core.tier_metrics import TierMetrics

PROMOTION_RULES = {
    "T3": {"min_successes": 10, "min_success_rate": 0.80, "promotes_to": "T2"},
    "T2": {"min_successes": 20, "min_success_rate": 0.85, "promotes_to": "T1"},
}

DEMOTION_THRESHOLD = 3  # consecutive failures


def check_promotion(metrics: TierMetrics) -> Optional[str]:
    """Check if workstream is eligible for tier promotion.

    Args:
        metrics: Current tier metrics

    Returns:
        New tier if promotion eligible, None otherwise
    """
    tier = metrics.current_tier
    if tier not in PROMOTION_RULES:
        return None  # T0 and T1 cannot be promoted

    rules = PROMOTION_RULES[tier]
    if (
        metrics.successful_attempts >= rules["min_successes"]
        and metrics.success_rate >= rules["min_success_rate"]
    ):
        return rules["promotes_to"]

    return None


def check_demotion(metrics: TierMetrics) -> Optional[str]:
    """Check if workstream should be demoted.

    Args:
        metrics: Current tier metrics

    Returns:
        New tier if demotion needed, None otherwise
    """
    if metrics.consecutive_failures >= DEMOTION_THRESHOLD:
        # Demote one tier down
        tier_order = ["T0", "T1", "T2", "T3"]
        current_idx = tier_order.index(metrics.current_tier)
        if current_idx < len(tier_order) - 1:
            return tier_order[current_idx + 1]

    return None


def check_tier_change(metrics: TierMetrics) -> Optional[str]:
    """Check if tier should change (promotion or demotion).

    Args:
        metrics: Current tier metrics

    Returns:
        New tier if change needed, None otherwise
    """
    # Check demotion first (higher priority)
    new_tier = check_demotion(metrics)
    if new_tier:
        return new_tier

    # Check promotion
    return check_promotion(metrics)
```

---

### Ожидаемый результат

- Автоматическая promotion T3→T2→T1 при достижении thresholds
- Автоматическая demotion при repeated failures
- CLI команды для manual check и metrics view
- Метрики персистентны в `.sdp/tier_metrics.json`

---

### Scope Estimate

- Файлов: ~4 создано + ~2 изменено
- Строк: ~400-600 (MEDIUM)
- Токенов: ~2500-3500

---

### Критерий завершения

```bash
# Unit tests
pytest sdp/tests/unit/core/test_tier_promoter.py -v
pytest sdp/tests/unit/core/test_tier_metrics.py -v

# CLI command works
sdp tier promote-check
sdp tier metrics WS-410-01

# Type checks
mypy sdp/src/sdp/core/tier_promoter.py
mypy sdp/src/sdp/core/tier_metrics.py

# Lint
ruff check sdp/src/sdp/core/
```

---

### Ограничения

- НЕ изменять существующие capability tier rules (T0-T3)
- НЕ автоматически продвигать T0 (Architect tier — ручной только)
- Metrics хранилище должно быть thread-safe (если concurrent builds)

---

### Related Issues

- Referenced in: `docs/reviews/F410-cross-ws-review.md` (Medium Priority recommendation)
- Depends on: WS-410-01, WS-410-02, WS-410-05
