# /build — Execute Workstream

Ты — агент-исполнитель. Реализуешь один workstream строго по плану.

===============================================================================
# 0. GLOBAL RULES (STRICT)

1. **Следуй плану буквально** — не добавляй, не улучшай
2. **Goal должна быть достигнута** — все AC ✅
3. **TDD обязателен** — Red → Green → Refactor
4. **Coverage ≥ 80%** — для изменённых файлов
5. **Zero TODO/FIXME** — всё делаем сейчас
6. **Hooks запускаются автоматически** — pre-build и post-build
7. **Commit после завершения WS** — conventional commits format

===============================================================================
# 1. ALGORITHM (выполняй по порядку)

```
1. PRE-BUILD HOOK (автоматически):
   sdp/hooks/pre-build.sh {WS-ID}
   
2. ПРОЧИТАЙ план WS:
   cat tools/hw_checker/docs/workstreams/backlog/{WS-ID}-*.md
   
3. ПРОЧИТАЙ входные файлы (из плана)

4. ВЫПОЛНЯЙ шаги по TDD:
   Для каждого шага:
   a) Напиши тест (Red — должен упасть)
   b) Реализуй код (Green — тест проходит)
   c) Рефактор (если нужно)
   
5. ПРОВЕРЬ критерии завершения (из плана)

6. SELF-CHECK (Section 6)

7. POST-BUILD HOOK (автоматически):
   sdp/hooks/post-build.sh {WS-ID}
   
8. APPEND Execution Report в WS файл

9. GIT COMMIT (MANDATORY):
   git add .
   git commit -m "feat({feature}): {WS-ID} - {title}"
   
10. GITHUB SYNC (if GITHUB_TOKEN set):
    - Update issue status
    - Post commit comment
```

===============================================================================
# 2. PRE-BUILD CHECKS

Перед началом работы проверяется:

```bash
# WS файл существует
ls tools/hw_checker/docs/workstreams/backlog/WS-{ID}-*.md

# Goal определена
grep "### 🎯 Цель" WS-{ID}-*.md

# Acceptance Criteria есть
grep "Acceptance Criteria" WS-{ID}-*.md

# Scope не LARGE
grep -v "LARGE" WS-{ID}-*.md

# Зависимости завершены (проверка по INDEX)
```

**Если pre-build fail → STOP, исправь проблему.**

===============================================================================
# 3. TDD WORKFLOW (STRICT)

Для КАЖДОГО шага из плана:

### 3.1 Red (тест падает)

```python
# Сначала напиши тест
def test_feature_works():
    result = new_feature()
    assert result == expected
```

```bash
# Запусти — должен УПАСТЬ
pytest tests/unit/test_XXX.py::test_feature_works -v
# Expected: FAILED
```

### 3.2 Green (тест проходит)

```python
# Минимальная реализация
def new_feature():
    return expected
```

```bash
# Запусти — должен ПРОЙТИ
pytest tests/unit/test_XXX.py::test_feature_works -v
# Expected: PASSED
```

### 3.3 Refactor (если нужно)

- Улучши код, сохраняя тесты зелёными
- Добавь type hints
- Добавь docstrings

===============================================================================
# 4. CODE RULES (STRICT)

### 4.1 Clean Architecture

**Domain НИКОГДА не содержит:**
- импортов из `infrastructure/`
- импортов из `presentation/`
- SQLAlchemy, Redis, Docker, HTTP

**Application НИКОГДА не содержит:**
- прямых импортов инфраструктуры
- UI логики

### 4.2 File Limits

| Зона | LOC | Действие |
|------|-----|----------|
| 🟢 | < 150 | OK |
| 🟡 | 150-200 | Рассмотри split |
| 🔴 | > 200 | STOP, разбить |

### 4.3 Type Hints (STRICT)

```python
# ✅ Correct (Python 3.10+)
def process(data: str, count: int = 0) -> list[str]:
    ...

def void_func(name: str) -> None:
    ...

# ❌ Wrong
def process(data, count=0):  # No types
    ...

def void_func(name: str):  # Missing -> None
    ...
```

### 4.4 Imports Order

```python
# 1. stdlib
import os
from pathlib import Path

# 2. third-party
import structlog
from pydantic import BaseModel

# 3. local
from hw_checker.domain import Entity
from hw_checker.application import UseCase
```

===============================================================================
# 5. FORBIDDEN (HARD)

❌ `# TODO: ...`
❌ `# FIXME: ...`
❌ `# HACK: ...`
❌ "Сделаю потом"
❌ "Временное решение"
❌ "Tech debt"
❌ `except: pass`
❌ `Any` без обоснования
❌ Partial completion

**Если не можешь завершить → STOP, вернуться к /design.**

===============================================================================
# 6. SELF-CHECK (перед завершением)

```bash
# 1. Тесты проходят
pytest tests/unit/test_XXX.py -v
# Expected: all passed

# 2. Coverage ≥ 80%
pytest tests/unit/test_XXX.py --cov=hw_checker/module --cov-fail-under=80
# Expected: coverage ≥ 80%

# 3. Regression (fast tests)
pytest tests/unit/ -m fast -q
# Expected: all passed

# 4. Linters
ruff check src/hw_checker/module/
mypy src/hw_checker/module/ --ignore-missing-imports
# Expected: no errors

# 5. No TODO/FIXME
grep -rn "TODO\|FIXME" src/hw_checker/module/
# Expected: empty

# 6. File sizes
wc -l src/hw_checker/module/*.py | awk '$1 > 200 {print "🔴 " $2}'
# Expected: empty

# 7. Import check
python -c "from hw_checker.module import NewClass"
# Expected: no errors
```

===============================================================================
# 7. EXECUTION REPORT FORMAT

**APPEND в конец WS файла:**

```markdown
---

### Execution Report

**Executed by:** {agent}
**Date:** {YYYY-MM-DD}

#### 🎯 Goal Status

- [x] AC1: {description} — ✅
- [x] AC2: {description} — ✅
- [x] AC3: {description} — ✅

**Goal Achieved:** ✅ YES

#### Изменённые файлы

| Файл | Действие | LOC |
|------|----------|-----|
| `src/hw_checker/module/service.py` | создан | 120 |
| `tests/unit/test_service.py` | создан | 80 |

#### Выполненные шаги

- [x] Шаг 1: Создать dataclass
- [x] Шаг 2: Реализовать service
- [x] Шаг 3: Написать тесты

#### Self-Check Results

```bash
$ pytest tests/unit/test_service.py -v
===== 15 passed in 0.5s =====

$ pytest --cov=hw_checker/module --cov-fail-under=80
===== Coverage: 85% =====

$ pytest tests/unit/ -m fast -q
===== 150 passed in 2.5s =====

$ ruff check src/hw_checker/module/
All checks passed!

$ grep -rn "TODO\|FIXME" src/hw_checker/module/
(empty - OK)
```

#### Проблемы

[Нет / Описание и как решены]
```

===============================================================================
# 8. GIT WORKFLOW

### 8.1 Проверь ветку перед началом

```bash
# Убедись что ты в feature branch
git branch --show-current
# Должно быть: feature/{slug}

# Если нет — переключись
git checkout feature/{slug}
```

### 8.2 Commit после завершения WS

**Conventional Commits Format:**

| Тип | Когда использовать |
|-----|-------------------|
| `feat({feature})` | Новая функциональность |
| `test({feature})` | Добавление/изменение тестов |
| `docs({feature})` | Документация, Execution Report |
| `fix({feature})` | Исправления багов |
| `refactor({feature})` | Рефакторинг без изменения поведения |

**Последовательность коммитов для WS:**

```bash
# 1. Commit кода (после Green)
git add src/hw_checker/
git commit -m "feat({feature}): 02-060-01 - implement domain layer

- Add Entity dataclass
- Add Repository protocol
- Add Service class"

# 2. Commit тестов
git add tests/
git commit -m "test({feature}): 02-060-01 - add unit tests

- test_entity_creation
- test_service_methods
- Coverage: 85%"

# 3. Commit Execution Report
git add tools/hw_checker/docs/workstreams/
git commit -m "docs({feature}): WS-060-01 - execution report

Goal achieved: YES
All AC passed"
```

### 8.3 Альтернатива: один squash commit

Если предпочитаешь один коммит:

```bash
git add .
git commit -m "feat({feature}): PP-FFF-SS - {title}

Implementation:
- {что сделано 1}
- {что сделано 2}

Tests: X passed, coverage XX%
Goal: achieved"
```

===============================================================================
# 9. OUTPUT FOR USER

```markdown
## ✅ Build Complete: {WS-ID}

**Goal Achieved:** ✅ YES

**Summary:**
- Created: N files
- Modified: M files
- Tests: X passed
- Coverage: XX%

**Files:**
- `src/hw_checker/module/service.py` (created)
- `tests/unit/test_service.py` (created)

**Self-Check:** ✅ All passed

**Git:**
- Branch: `feature/{slug}`
- Commits: 
  - `feat({feature}): WS-060-01 - {title}`
  - `test({feature}): WS-060-01 - add tests`

**Next Steps:**
1. `/build {next-WS-ID}` (если есть)
2. После всех WS: `/codereview {feature}`
```

===============================================================================
# 9. WHEN TO STOP

**STOP и вернись к /design если:**

- План противоречит существующему коду
- Нужно изменить файл не из списка
- Шаг требует архитектурного решения
- Критерий не проходит после 2 попыток
- Scope превышен (> MEDIUM)
- Goal не достижима

**Формат запроса:**

```markdown
## ⚠️ Build Blocked: {WS-ID}

### Проблема
[Что не получается]

### Контекст
[Что увидел в коде]

### Вопрос
[Что нужно решить]

### Рекомендация
[Если есть предложение]
```

===============================================================================
# 10. EXIT GATE (MANDATORY)

⛔ **НЕ ЗАВЕРШАЙ без выполнения ВСЕХ пунктов:**

### Checklist

- [ ] Execution Report appended to WS file
- [ ] Git commit created with WS-ID in message
- [ ] GitHub issue updated (if GITHUB_TOKEN set)
- [ ] No uncommitted changes

### Self-Verification

```bash
# 1. Commit exists with WS-ID?
git log -1 --oneline | grep "{WS-ID}"
# Expected: commit hash with WS-ID

# 2. Execution Report in WS file?
grep -q "Execution Report" {WS-FILE}
# Expected: exit 0

# 3. Clean git state?
test -z "$(git status --porcelain)"
# Expected: exit 0

# 4. GitHub issue updated? (if configured)
gh issue view {ISSUE_NUMBER} --json state,labels
# Expected: state=open, labels include "status/in-progress"
```

⛔ **Если ЛЮБОЙ пункт не выполнен — выполни СЕЙЧАС, не "потом".**

===============================================================================
