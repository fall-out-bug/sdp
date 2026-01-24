---
ws_id: PP-FFF-WW
feature: FFF
status: backlog|active|completed|blocked
size: SMALL|MEDIUM|LARGE
project_id: PP
github_issue: null
assignee: null
depends_on:
  - PP-FFF-WW  # Optional: list of dependent WS IDs
---

## WS-PP-FFF-WW: Title

### 🎯 Цель (Goal)

**Что должно РАБОТАТЬ после завершения WS:**
- [First specific outcome]
- [Second specific outcome]

**Acceptance Criteria:**
- [ ] AC1: [First criterion - specific, measurable]
- [ ] AC2: [Second criterion - specific, measurable]
- [ ] AC3: [Third criterion - specific, measurable]
- [ ] Coverage ≥ 80%
- [ ] mypy --strict passes

**⚠️ WS НЕ завершён, пока Goal не достигнута (все AC ✅).**

---

### Контекст

[Background information about the workstream]

### Зависимость

[List dependencies or write "Независимый" for no dependencies]

### Входные файлы

[List input files or sections to read]

### Шаги

1. **[Step 1 title]**

   [Detailed instructions for step 1]

2. **[Step 2 title]**

   [Detailed instructions for step 2]

### Код

```python
# Готовый код для copy-paste
# Полные type hints
```

### Ожидаемый результат

[Description of expected outcome]

### Scope Estimate

- Файлов: ~[number]
- Строк: ~[number] ([SMALL|MEDIUM|LARGE])
- Токенов: ~[number]

### Критерий завершения

```bash
# Verification commands
pytest tests/unit/test_module.py -v
pytest --cov=module --cov-fail-under=80
mypy module/ --strict
```

### Ограничения

- НЕ [constraint 1]
- НЕ [constraint 2]

---

## Execution Report

**Executed by:** ______
**Date:** ______

### Goal Status
- [ ] AC1-AC3 — ✅

**Goal Achieved:** ______

### Files Changed
| File | Action | LOC |
|------|--------|-----|
| | | |

### Commit
______
