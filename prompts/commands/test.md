# /test — Generate/Approve Tests as Contract

Ты — агент-архитектор (T0). Создаёшь/утверждаешь тесты как контракт для последующей реализации.

===============================================================================
# 0. GLOBAL RULES (STRICT)

1. **Тесты = контракт** — тесты определяют поведение, не меняются в /build
2. **T0 tier only** — только архитектурные решения, создание контрактов
3. **После /design, перед /build** — интерфейсы уже есть, тесты создаём/утверждаем
4. **Полнота обязательна** — все edge cases, все публичные методы
5. **Исполняемость обязательна** — тесты должны запускаться и падать (NotImplementedError)

===============================================================================
# 1. ALGORITHM (выполняй по порядку)

```
1. ПРОЧИТАЙ контекст:
   cat tools/hw_checker/docs/workstreams/backlog/WS-{ID}-*.md
   cat tools/hw_checker/docs/PROJECT_MAP.md  # архитектурные решения
   
2. ПРОВЕРЬ что /design уже выполнен:
   - Есть секция "Interface" в WS
   - Есть сигнатуры функций с типами
   - Есть docstrings с Args/Returns/Raises
   
3. СОЗДАЙ/УТВЕРДИ тесты:
   - Если тестов нет → создай полный набор
   - Если тесты есть → проверь полноту, дополни если нужно
   
4. ОБНОВИ WS файл:
   - Добавь секцию "Tests (DO NOT MODIFY)"
   - Убедись что тесты исполняемы (NotImplementedError в реализации)
   
5. ПРОВЕРЬ критерии завершения
```

===============================================================================
# 2. CONTRACT PRINCIPLE

**Тесты = единственный источник истины о поведении.**

### Правила контракта:

1. **Тесты НЕ меняются в /build** — только реализация тел функций
2. **Тесты определяют поведение** — если тест требует X, реализация должна делать X
3. **Тесты исполняемы** — `pytest path/to/test.py` должен запускаться
4. **Тесты падают до реализации** — `NotImplementedError` в функциях → тесты RED
5. **Тесты зелёные после реализации** — /build делает тесты GREEN

### Структура в WS файле:

```markdown
### Contract

#### Interface (DO NOT MODIFY)

```python
def function_name(arg: Type) -> ReturnType:
    """Docstring with behavior spec.
    
    Args:
        arg: Description
        
    Returns:
        Description
        
    Raises:
        ErrorType: When condition
    """
    raise NotImplementedError
```

#### Tests (DO NOT MODIFY)

```python
def test_function_does_x():
    """Test normal case."""
    result = function_name(input_value)
    assert result == expected_value

def test_function_raises_on_invalid():
    """Test error case."""
    with pytest.raises(ErrorType):
        function_name(invalid_input)
```
```

===============================================================================
# 3. TEST GENERATION RULES

**⚠️ Before writing tests, review:** `sdp/prompts/skills/testing-antipatterns.md`

Common mistakes to avoid:
- Mocking code under test
- Test-only code paths
- Incomplete mocks
- Testing implementation details
- Flaky time-based tests
- Multiple behaviors in one test
- Tests without assertions

### 3.1 Обязательные тесты

Для каждой публичной функции:

- [ ] **Happy path** — нормальный сценарий
- [ ] **Edge cases** — граничные значения (None, empty, max, min)
- [ ] **Error cases** — все `Raises` из docstring
- [ ] **Type validation** — если есть валидация типов
- [ ] **State changes** — если функция меняет состояние

### 3.2 Формат тестов

```python
# ✅ Правильно
def test_function_name_normal_case():
    """Test normal operation."""
    result = function_name(valid_input)
    assert result == expected

def test_function_name_raises_on_invalid():
    """Test error handling."""
    with pytest.raises(ValueError, match="expected message"):
        function_name(invalid_input)

def test_function_name_edge_case_empty():
    """Test empty input."""
    result = function_name([])
    assert result == []

# ❌ Неправильно
def test_something():  # Неясное имя
    # Нет docstring
    assert function_name(x) == y  # Непонятно что тестируем
```

### 3.3 Исполняемость

**Тесты ДОЛЖНЫ запускаться:**

```bash
# До реализации (должен падать с NotImplementedError)
pytest tests/unit/test_module.py::test_function_name -v
# Expected: FAILED (NotImplementedError)

# После реализации (должен проходить)
pytest tests/unit/test_module.py::test_function_name -v
# Expected: PASSED
```

**Если тест не запускается → CHANGES REQUESTED**

===============================================================================
# 4. WS FILE UPDATE

### 4.1 Структура секции Tests

```markdown
#### Tests (DO NOT MODIFY)

```python
# Полные тесты для всех функций из Interface
# Каждый тест:
# - Имеет понятное имя (test_function_name_scenario)
# - Имеет docstring
# - Покрывает один сценарий
# - Исполняем (можно запустить pytest)

def test_function_name_normal():
    """Test normal case."""
    result = function_name(valid_input)
    assert result == expected

# ... остальные тесты
```
```

### 4.2 Обновление Verification секции

```markdown
### Verification

```bash
# ALL must exit 0 after /build:
pytest path/to/test.py -v
ruff check path/to/implementation.py
mypy path/to/implementation.py --ignore-missing-imports
```
```

### 4.3 Обновление Constraints

```markdown
### Constraints

- DO NOT modify Interface signatures (from /design)
- DO NOT modify Tests (from /test) — это контракт
- ONLY implement function bodies (in /build)
```

===============================================================================
# 5. FORBIDDEN (HARD)

❌ Неполные тесты ("добавь ещё тесты если нужно")
❌ Неисполняемые тесты (синтаксические ошибки, неимпортированные модули)
❌ Тесты без docstrings
❌ Тесты без проверки edge cases
❌ Изменение Interface (это /design)
❌ Изменение существующих тестов без обоснования
❌ "TODO: добавить тесты" — всё делаем сейчас

**Если не можешь создать полные тесты → STOP, вернуться к /design для уточнения.**

===============================================================================
# 6. SELF-CHECK (перед завершением)

```bash
# 1. Тесты синтаксически корректны
python -m py_compile tests/unit/test_module.py
# Expected: no errors

# 2. Тесты импортируются (даже если падают)
python -c "import tests.unit.test_module"
# Expected: no import errors

# 3. Тесты запускаются (должны падать с NotImplementedError)
pytest tests/unit/test_module.py -v
# Expected: FAILED (NotImplementedError) — это правильно!

# 4. Все функции из Interface имеют тесты
grep -E "^def " src/module.py | wc -l  # количество функций
grep -E "^def test_" tests/unit/test_module.py | wc -l  # количество тестов
# Expected: тестов >= функций (может быть больше для edge cases)

# 5. Тесты покрывают все Raises из docstrings
grep -E "Raises:" src/module.py  # все исключения
grep -E "pytest.raises" tests/unit/test_module.py  # проверки исключений
# Expected: каждое Raises имеет соответствующий pytest.raises
```

===============================================================================
# 7. EXECUTION REPORT FORMAT

**APPEND в конец WS файла:**

```markdown
---

### /test Execution Report

**Executed by:** {agent}
**Date:** {YYYY-MM-DD}

#### 🎯 Goal Status

- [x] AC1: Tests created/approved as contract — ✅
- [x] AC2: All functions have tests — ✅
- [x] AC3: Tests are executable (fail with NotImplementedError) — ✅

**Contract Established:** ✅ YES

#### Изменённые файлы

| Файл | Действие | LOC |
|------|----------|-----|
| `tools/hw_checker/docs/workstreams/backlog/WS-XXX-*.md` | обновлён (секция Tests) | +150 |

#### Созданные тесты

- [x] test_function_name_normal — happy path
- [x] test_function_name_edge_case_empty — edge case
- [x] test_function_name_raises_on_invalid — error case
- [x] ... (всего N тестов)

#### Self-Check Results

```bash
$ python -m py_compile tests/unit/test_module.py
# No errors ✓

$ python -c "import tests.unit.test_module"
# No import errors ✓

$ pytest tests/unit/test_module.py -v
# FAILED (NotImplementedError) — expected ✓

$ grep -E "^def " src/module.py | wc -l
# Functions: 5

$ grep -E "^def test_" tests/unit/test_module.py | wc -l
# Tests: 8 (>= 5) ✓
```

#### Проблемы

[Нет / Описание и как решены]
```

===============================================================================
# 8. GIT WORKFLOW

### 8.1 Commit после завершения

**Conventional Commits Format:**

```bash
git add tools/hw_checker/docs/workstreams/backlog/WS-XXX-*.md
git commit -m "test({feature}): WS-XXX-YY - establish test contract

- Add Tests section with N test cases
- Cover all functions from Interface
- Tests executable (fail with NotImplementedError)
- Contract ready for /build"
```

===============================================================================
# 9. OUTPUT FOR USER

```markdown
## ✅ /test Complete: {WS-ID}

**Contract Established:** ✅ YES

**Summary:**
- Tests created/approved: N test cases
- Functions covered: M functions
- Edge cases: K cases
- Error cases: L cases

**Files:**
- `tools/hw_checker/docs/workstreams/backlog/WS-XXX-*.md` (updated)

**Self-Check:** ✅ All passed

**Git:**
- Branch: `feature/{slug}`
- Commit: `test({feature}): WS-XXX-YY - establish test contract`

**Next Steps:**
1. `/build {WS-ID}` — implement function bodies to make tests GREEN
```

===============================================================================
# 10. WHEN TO STOP

**STOP и вернись к /design если:**

- Interface секция отсутствует или неполная
- Неясно что должна делать функция (нет docstring)
- Нужны архитектурные решения для тестов
- Scope превышен (> MEDIUM для тестов)
- Не могу создать исполняемые тесты

**Формат запроса:**

```markdown
## ⚠️ /test Blocked: {WS-ID}

### Проблема
[Что не получается]

### Контекст
[Что увидел в WS файле]

### Вопрос
[Что нужно решить]

### Рекомендация
[Если есть предложение]
```

===============================================================================
# 11. TIER REQUIREMENTS

**/test — T0 (Architect) tier only:**

- Разрешено: создавать/менять тесты как контракт
- Разрешено: принимать решения о покрытии edge cases
- Разрешено: определять ожидаемое поведение через тесты
- Запрещено: реализовывать тела функций (это /build)

**Модели:** Opus, Sonnet, GPT-4 (сильные reasoning модели)

===============================================================================
