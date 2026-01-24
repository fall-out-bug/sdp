---
ws_id: 00-410-06
project_id: 00
feature: F008
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

## 00-410-06: Model Selection Optimization

### 🎯 Цель (Goal)

**Что должно РАБОТАТЬ после завершения WS:**
- Builder router выбирает модель на основе cost/availability/context-size, а не просто первую из списка
- Реализована взвешенная логика выбора с учётом метрик модели
- Добавлены тесты для разных стратегий выбора

**Acceptance Criteria (критерии приёмки):**
- [ ] AC1: Добавлены поля `cost_per_1m_tokens`, `availability_pct`, `context_window` в `ModelProvider`
- [ ] AC2: Реализована функция weighted selection с параметрами веса (cost, availability, context)
- [ ] AC3: `select_model_for_tier()` использует weighted selection вместо `models[0]`
- [ ] AC4: Добавлены unit-тесты для разных весовых конфигураций
- [ ] AC5: `model-mapping.md` обновлён с метриками для всех моделей

**⚠️ Правило:** WS НЕ завершён, пока Goal не достигнута (все AC ✅).

---

### Контекст

**Проблема:** Текущая реализация `select_model_for_tier()` в `builder_router.py:199-201` просто возвращает первую модель из списка для данного tier:

```python
# Current implementation (naive)
def select_model_for_tier(tier: str, registry: ModelRegistry) -> ModelProvider:
    models = registry.get_models_for_tier(tier)
    if not models:
        raise ValueError(f"No models available for tier {tier}")
    return models[0]  # ⚠️ Always picks first - no optimization
```

**Нужно:** Добавить интеллектуальный выбор на основе:
1. **Cost** - минимизировать $/1M tokens (важно для T2/T3 с высокой частотой)
2. **Availability** - предпочитать модели с высоким uptime (99%+)
3. **Context** - учитывать требования к context window (для больших workstreams)

**Приоритет:** HIGH - это наиболее заметный gap в текущей реализации F410.

---

### Зависимость

- WS-410-03 (Model mapping registry) — обязательная
- WS-410-05 (Builder router) — обязательная

---

### Входные файлы

- `sdp/src/sdp/core/builder_router.py` — текущий naive selection
- `sdp/src/sdp/core/model_mapping.py` — ModelProvider dataclass
- `sdp/docs/model-mapping.md` — модели по tier (нужно обновить с метриками)
- `sdp/tests/unit/core/test_builder_router.py` — добавить тесты для weighted selection

---

### Шаги

1. **Обновить `ModelProvider` dataclass** с метриками:
   ```python
   @dataclass
   class ModelProvider:
       provider: str
       model: str
       context: str
       tool_use: bool
       # New fields:
       cost_per_1m_tokens: float  # USD per 1M tokens (input)
       availability_pct: float    # Uptime percentage (0.0-1.0)
       context_window: int        # Max tokens
   ```

2. **Реализовать weighted selection**:
   ```python
   def select_model_weighted(
       models: list[ModelProvider],
       weights: dict[str, float] = {"cost": 0.4, "availability": 0.3, "context": 0.3},
       required_context: int = 0
   ) -> ModelProvider:
       """Select model based on weighted scoring."""
       # Filter by context requirements
       # Normalize metrics
       # Compute weighted score
       # Return highest score
   ```

3. **Обновить `select_model_for_tier()`**:
   ```python
   def select_model_for_tier(
       tier: str,
       registry: ModelRegistry,
       required_context: int = 0,
       weights: Optional[dict[str, float]] = None
   ) -> ModelProvider:
       models = registry.get_models_for_tier(tier)
       if not models:
           raise ValueError(f"No models available for tier {tier}")

       # Use weighted selection
       return select_model_weighted(models, weights or DEFAULT_WEIGHTS, required_context)
   ```

4. **Обновить `model-mapping.md`** с метриками:
   ```markdown
   | Provider | Model | Cost ($/1M) | Availability | Context |
   |----------|-------|-------------|--------------|---------|
   | Anthropic | Claude Haiku 4 | 0.25 | 99.9% | 200K |
   | OpenAI | GPT-4o-mini | 0.15 | 99.5% | 128K |
   ```

5. **Добавить unit-тесты**:
   - Cost-optimized selection (weight cost=1.0)
   - Availability-optimized (weight availability=1.0)
   - Context-filtered selection
   - Equal weighting (default)

---

### Код

```python
# sdp/src/sdp/core/model_mapping.py

@dataclass
class ModelProvider:
    """Model provider with performance metrics."""
    provider: str
    model: str
    context: str
    tool_use: bool
    cost_per_1m_tokens: float  # NEW
    availability_pct: float    # NEW
    context_window: int        # NEW

# sdp/src/sdp/core/builder_router.py

DEFAULT_WEIGHTS = {
    "cost": 0.4,
    "availability": 0.3,
    "context": 0.3,
}

def select_model_weighted(
    models: list[ModelProvider],
    weights: dict[str, float],
    required_context: int = 0,
) -> ModelProvider:
    """Select model using weighted scoring.

    Args:
        models: Available models for tier
        weights: Weight dict (cost, availability, context)
        required_context: Minimum context window required

    Returns:
        Best model according to weighted score

    Raises:
        ValueError: If no models meet requirements
    """
    # Filter by context
    candidates = [m for m in models if m.context_window >= required_context]
    if not candidates:
        raise ValueError(f"No models with context >= {required_context}")

    # Normalize metrics (0-1 scale, higher is better)
    max_cost = max(m.cost_per_1m_tokens for m in candidates)
    scores = []
    for model in candidates:
        cost_score = 1 - (model.cost_per_1m_tokens / max_cost)  # Lower cost = higher score
        avail_score = model.availability_pct  # Already 0-1
        context_score = model.context_window / 200_000  # Normalize to 200K max

        weighted_score = (
            weights["cost"] * cost_score +
            weights["availability"] * avail_score +
            weights["context"] * context_score
        )
        scores.append((weighted_score, model))

    # Return highest score
    scores.sort(key=lambda x: x[0], reverse=True)
    return scores[0][1]


def select_model_for_tier(
    tier: str,
    registry: ModelRegistry,
    required_context: int = 0,
    weights: Optional[dict[str, float]] = None,
) -> ModelProvider:
    """Select model provider for given capability tier.

    Uses weighted selection based on cost, availability, context.

    Args:
        tier: Capability tier (T0, T1, T2, T3)
        registry: Model registry
        required_context: Minimum context tokens needed
        weights: Custom weights (defaults to DEFAULT_WEIGHTS)

    Returns:
        Selected ModelProvider instance

    Raises:
        ValueError: If tier is invalid or no models available
    """
    models = registry.get_models_for_tier(tier)

    if not models:
        raise ValueError(f"No models available for tier {tier}")

    return select_model_weighted(
        models,
        weights or DEFAULT_WEIGHTS,
        required_context
    )
```

---

### Ожидаемый результат

- Weighted selection реализован и протестирован
- `model-mapping.md` содержит полные метрики
- `select_model_for_tier()` использует интеллектуальный выбор
- Тесты покрывают разные стратегии выбора

---

### Scope Estimate

- Файлов: ~3 изменено + ~1 тест файл
- Строк: ~200-400 (MEDIUM)
- Токенов: ~1500-2500

---

### Критерий завершения

```bash
# Unit tests
pytest sdp/tests/unit/core/test_model_selection.py -v

# All router tests pass
pytest sdp/tests/unit/core/test_builder_router.py -v

# Type checks
mypy sdp/src/sdp/core/builder_router.py
mypy sdp/src/sdp/core/model_mapping.py

# Lint
ruff check sdp/src/sdp/core/
```

---

### Ограничения

- НЕ менять публичные API `BuilderRouter` и `ModelRegistry`
- НЕ удалять существующие тесты
- Backward compatibility: если weights=None, использовать DEFAULT_WEIGHTS

---

### Related Issues

- Referenced in: `docs/reviews/F410-cross-ws-review.md` (High Priority recommendation)
