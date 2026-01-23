# UAT Guide: F{XX} - {Feature Name}

**Created:** {YYYY-MM-DD}
**Feature:** F{XX}
**Workstreams:** WS-{XX}-01, WS-{XX}-02, ...

---

## Overview

{Что делает фича в 2-3 предложениях для человека}

---

## Prerequisites

Перед тестированием убедись:

- [ ] Docker запущен (`docker ps`)
- [ ] `poetry install` выполнен в `tools/hw_checker/`
- [ ] `.env` или `hw_checker.yaml` настроен
- [ ] База данных доступна (если нужно)
- [ ] Redis запущен (если нужно)

```bash
# Quick prerequisite check
cd tools/hw_checker
poetry run python -c "from hw_checker import __version__; print(f'Version: {__version__}')"
```

---

## Quick Verification (5 минут)

### 1. Smoke Test

```bash
cd tools/hw_checker

# Основная проверка
poetry run hwc {main_command}

# Ожидаемый результат:
# {описание что должно быть}
```

### 2. Visual Inspection

- [ ] Открой {что открыть: logs/UI/API}
- [ ] Проверь что {что должно отображаться}
- [ ] Убедись что {нет ошибок/warnings}

---

## Detailed Test Scenarios

### Scenario 1: Happy Path

**Описание:** {основной use case}

**Steps:**
1. {step 1}
2. {step 2}
3. {step 3}

**Expected:**
- {expectation 1}
- {expectation 2}

**Actual:** ____________________

**Status:** ⬜ Pass / ⬜ Fail

---

### Scenario 2: Error Handling

**Описание:** {как система обрабатывает ошибки}

**Steps:**
1. {trigger error condition}
2. {observe response}

**Expected:**
- Graceful error message (не stack trace)
- Логирование ошибки
- Система продолжает работать

**Actual:** ____________________

**Status:** ⬜ Pass / ⬜ Fail

---

### Scenario 3: Edge Cases

**Описание:** {граничные случаи}

**Steps:**
1. {edge case input}
2. {observe behavior}

**Expected:**
- {expected handling}

**Actual:** ____________________

**Status:** ⬜ Pass / ⬜ Fail

---

## Red Flags Checklist

**❌ Если видишь любой из этих признаков — агент накосячил:**

| # | Red Flag | What to Check | Severity |
|---|----------|---------------|----------|
| 1 | Stack trace в output | Logs, stderr | 🔴 HIGH |
| 2 | Пустой response | API response body | 🔴 HIGH |
| 3 | Timeout (>30s) | Network, DB connection | 🟡 MEDIUM |
| 4 | Warning в логах | Log files | 🟡 MEDIUM |
| 5 | Неожиданный формат данных | Response structure | 🟡 MEDIUM |
| 6 | Deprecated warnings | Console output | 🟢 LOW |

**Что делать если нашёл Red Flag:**
1. Скопируй error message / screenshot
2. Проверь соответствующий WS Execution Report
3. Создай issue или вернись к `/codereview`

---

## Code Sanity Checks

Быстрая проверка что код в порядке:

```bash
cd tools/hw_checker

# 1. Нет TODO/FIXME
grep -rn "TODO\|FIXME" src/hw_checker/{feature_module}/
# Ожидание: пусто

# 2. Размер файлов разумный
wc -l src/hw_checker/{feature_module}/*.py
# Ожидание: все < 200 строк

# 3. Clean Architecture соблюдена
grep -r "from hw_checker.infrastructure" src/hw_checker/domain/
# Ожидание: пусто

# 4. Тесты проходят
poetry run pytest tests/unit/test_{feature}*.py -v
# Ожидание: все passed

# 5. Coverage достаточный
poetry run pytest tests/unit/test_{feature}*.py --cov=hw_checker/{feature_module} --cov-report=term-missing
# Ожидание: >= 80%
```

---

## Performance Baseline (если применимо)

| Операция | Expected | Acceptable | Measured |
|----------|----------|------------|----------|
| {operation 1} | < 100ms | < 500ms | ___ms |
| {operation 2} | < 1s | < 5s | ___s |
| {operation 3} | < 5s | < 30s | ___s |

---

## Sign-off

### Pre-Sign-off Checklist

- [ ] Все scenarios пройдены
- [ ] Red flags отсутствуют
- [ ] Code sanity checks пройдены
- [ ] Performance в пределах baseline

### Approval

| Role | Name | Date | Signature |
|------|------|------|-----------|
| Developer (агент) | {agent} | {date} | ✅ |
| Reviewer | {reviewer} | {date} | ⬜ |
| **Human Tester** | ____________ | ____________ | ⬜ |

### Final Verdict

⬜ **APPROVED** — готово к deploy
⬜ **NEEDS WORK** — требуются исправления (см. комментарии ниже)

### Comments

```
{комментарии от человека-тестировщика}
```

---

## Related

- Feature Spec: `docs/specs/feature_{XX}/feature.md`
- Workstreams: `docs/workstreams/backlog/WS-{XX}-*.md`
- Review Results: см. каждый WS файл
