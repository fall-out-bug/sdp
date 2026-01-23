# Testing Anti-Patterns — Common Mistakes to Avoid

**Purpose:** Prevent common testing mistakes that lead to false confidence, flaky tests, and maintenance burden.

---

## Core Principle

**Tests should verify behavior, not implementation details. Tests should be reliable, maintainable, and provide real value.**

---

## Anti-Pattern 1: Mocking What You're Testing

**Problem:** Mocking the code under test defeats the purpose of testing. You're testing the mock, not the real code.

### ❌ Bad Example

```python
# ❌ WRONG: Mocking the function we're testing
from unittest.mock import patch

def test_calculate_total():
    with patch('module.calculate_total') as mock_calc:
        mock_calc.return_value = 100
        result = calculate_total([1, 2, 3])
        assert result == 100
    # This test doesn't actually test calculate_total!
```

### ✅ Good Example

```python
# ✅ CORRECT: Test the real function
def test_calculate_total():
    result = calculate_total([1, 2, 3])
    assert result == 6  # Real calculation

# ✅ CORRECT: Mock only external dependencies
def test_calculate_total_with_discount():
    with patch('module.fetch_discount_rate') as mock_fetch:
        mock_fetch.return_value = 0.1  # Mock external API
        result = calculate_total_with_discount([100, 200])
        assert result == 270  # 300 * 0.9
```

### Detection Rule

```python
# Ruff rule (optional)
# Check for mocks of functions in the same module
# Pattern: patch('module.function_under_test')
```

---

## Anti-Pattern 2: Test-Only Code Paths

**Problem:** Code that only exists for tests creates maintenance burden and doesn't reflect real behavior.

### ❌ Bad Example

```python
# ❌ WRONG: Test-only parameter
def process_data(data: list[str], test_mode: bool = False) -> dict:
    if test_mode:
        return {"test": "data"}  # Only used in tests
    # Real logic...
    return process_real_data(data)

# Test
def test_process_data():
    result = process_data([], test_mode=True)
    assert result == {"test": "data"}
    # This doesn't test real behavior!
```

### ✅ Good Example

```python
# ✅ CORRECT: No test-only paths
def process_data(data: list[str]) -> dict:
    if not data:
        return {}  # Real edge case handling
    return process_real_data(data)

# Test
def test_process_data_empty():
    result = process_data([])
    assert result == {}  # Tests real behavior

def test_process_data_normal():
    result = process_data(["item1", "item2"])
    assert "item1" in result
```

### Detection Rule

```python
# Check for parameters with "test" in name
# Pattern: def function(..., test_mode=False, test_flag=True)
```

---

## Anti-Pattern 3: Incomplete Mocks

**Problem:** Mocks that don't match real behavior lead to tests passing but production failing.

### ❌ Bad Example

```python
# ❌ WRONG: Incomplete mock
from unittest.mock import Mock

def test_fetch_user_data():
    mock_api = Mock()
    mock_api.get.return_value = {"name": "John"}  # Missing fields!
    
    result = fetch_user_data(mock_api, user_id=123)
    assert result.name == "John"
    # Production API returns {"name": "John", "email": "...", "id": 123}
    # Code might fail when accessing result.email
```

### ✅ Good Example

```python
# ✅ CORRECT: Complete mock matching real API
def test_fetch_user_data():
    mock_api = Mock()
    mock_api.get.return_value = {
        "name": "John",
        "email": "john@example.com",
        "id": 123
    }  # Matches real API response
    
    result = fetch_user_data(mock_api, user_id=123)
    assert result.name == "John"
    assert result.email == "john@example.com"
    assert result.id == 123

# ✅ BETTER: Use fixtures or factories
@pytest.fixture
def mock_user_response():
    return {
        "name": "John",
        "email": "john@example.com",
        "id": 123
    }

def test_fetch_user_data(mock_user_response):
    mock_api = Mock()
    mock_api.get.return_value = mock_user_response
    result = fetch_user_data(mock_api, user_id=123)
    assert result.name == "John"
```

### Detection Rule

```python
# Check for mocks with minimal return values
# Pattern: Mock().return_value = {...} with < 3 fields
# (Heuristic, not perfect)
```

---

## Anti-Pattern 4: Testing Implementation Details

**Problem:** Tests that break when refactoring (even when behavior is correct) create maintenance burden.

### ❌ Bad Example

```python
# ❌ WRONG: Testing internal state
def test_calculate_total():
    calculator = Calculator()
    calculator._cache = {}  # Accessing private attribute
    result = calculator.calculate([1, 2, 3])
    assert calculator._cache == {6: True}  # Testing internal cache
    # If we change caching strategy, test breaks even though behavior is correct
```

### ✅ Good Example

```python
# ✅ CORRECT: Testing public behavior
def test_calculate_total():
    calculator = Calculator()
    result = calculator.calculate([1, 2, 3])
    assert result == 6  # Test public interface

def test_calculate_total_cached():
    calculator = Calculator()
    result1 = calculator.calculate([1, 2, 3])
    result2 = calculator.calculate([1, 2, 3])
    assert result1 == result2 == 6
    # Test caching behavior, not implementation
```

### Detection Rule

```python
# Check for access to private attributes in tests
# Pattern: obj._attribute or obj.__attribute
```

---

## Anti-Pattern 5: Flaky Tests with Timeouts

**Problem:** Tests that depend on timing are unreliable and fail randomly.

### ❌ Bad Example

```python
# ❌ WRONG: Time-based test
import time

def test_async_operation():
    start = time.time()
    result = async_operation()
    elapsed = time.time() - start
    assert elapsed < 0.1  # Flaky! Depends on system load
    assert result == expected
```

### ✅ Good Example

```python
# ✅ CORRECT: Use proper async testing
import asyncio

@pytest.mark.asyncio
async def test_async_operation():
    result = await async_operation()
    assert result == expected  # Test behavior, not timing

# ✅ CORRECT: Mock time if needed
from unittest.mock import patch

def test_scheduled_task():
    with patch('time.time', return_value=1000):
        result = scheduled_task()
        assert result == expected
```

### Detection Rule

```python
# Check for time.sleep, time.time comparisons in tests
# Pattern: time.sleep(..., assert time.time() - start < ...
```

---

## Anti-Pattern 6: Testing Multiple Things in One Test

**Problem:** When a test fails, you don't know which part broke. Hard to debug and maintain.

### ❌ Bad Example

```python
# ❌ WRONG: Testing multiple behaviors
def test_user_operations():
    user = create_user("John")
    assert user.name == "John"
    user.update_email("john@example.com")
    assert user.email == "john@example.com"
    user.delete()
    assert user.is_deleted is True
    # If this fails, which operation broke?
```

### ✅ Good Example

```python
# ✅ CORRECT: One behavior per test
def test_create_user():
    user = create_user("John")
    assert user.name == "John"

def test_update_user_email():
    user = create_user("John")
    user.update_email("john@example.com")
    assert user.email == "john@example.com"

def test_delete_user():
    user = create_user("John")
    user.delete()
    assert user.is_deleted is True
```

### Detection Rule

```python
# Check for tests with > 3 assertions (heuristic)
# Pattern: Multiple assert statements in one test
```

---

## Anti-Pattern 7: Tests Without Assertions

**Problem:** Tests that don't verify anything provide false confidence.

### ❌ Bad Example

```python
# ❌ WRONG: No assertions
def test_process_data():
    data = [1, 2, 3]
    result = process_data(data)
    # No assertion! Test always passes
```

### ✅ Good Example

```python
# ✅ CORRECT: Always assert expected behavior
def test_process_data():
    data = [1, 2, 3]
    result = process_data(data)
    assert result == {"sum": 6, "count": 3}

# ✅ CORRECT: Even for void functions, verify side effects
def test_save_to_file():
    data = {"key": "value"}
    save_to_file(data, "test.json")
    assert os.path.exists("test.json")
    with open("test.json") as f:
        assert json.load(f) == data
```

### Detection Rule

```python
# Check for test functions without assert statements
# Pattern: def test_*(): ... (no assert)
```

---

## Detection Rules Summary

### Optional Lint Rules (Ruff/Pylint)

```python
# Rule 1: Detect mocking of code under test
# Pattern: patch('module.function_under_test')
ANTIPATTERN_MOCK_UNDER_TEST = "T001"

# Rule 2: Detect test-only parameters
# Pattern: def func(..., test_mode=False, test_flag=True)
ANTIPATTERN_TEST_ONLY_PARAM = "T002"

# Rule 3: Detect incomplete mocks (heuristic)
# Pattern: Mock().return_value = {...} with < 3 fields
ANTIPATTERN_INCOMPLETE_MOCK = "T003"

# Rule 4: Detect testing implementation details
# Pattern: obj._attribute or obj.__attribute in tests
ANTIPATTERN_IMPL_DETAILS = "T004"

# Rule 5: Detect time-based tests
# Pattern: time.sleep, time.time comparisons
ANTIPATTERN_TIME_BASED = "T005"

# Rule 6: Detect tests without assertions
# Pattern: def test_*(): ... (no assert)
ANTIPATTERN_NO_ASSERT = "T006"
```

### Pre-commit Hook Integration

Add to `sdp/hooks/pre-commit.sh`:

```bash
# Check 8: Testing anti-patterns
echo ""
echo "Check 8: Testing anti-patterns"
TEST_FILES=$(echo "$STAGED_FILES" | grep -E "tests/.*test_.*\.py$" || true)

if [ -n "$TEST_FILES" ]; then
    # Check for common anti-patterns
    # (Implementation depends on available tools)
    if command -v ruff &> /dev/null; then
        ruff check --select T001,T002,T003,T004,T005,T006 $TEST_FILES
    fi
fi
```

---

## Integration with TDD Skill

When using `/test` command, refer to this guide to avoid anti-patterns:

1. **Before writing tests:** Review anti-patterns 1-7
2. **During test review:** Check for anti-pattern violations
3. **In pre-commit:** Run detection rules automatically

See: `sdp/prompts/commands/test.md` for TDD workflow.

---

## Quick Reference Checklist

Before committing tests, verify:

- [ ] Not mocking code under test (Anti-Pattern 1)
- [ ] No test-only code paths (Anti-Pattern 2)
- [ ] Mocks match real behavior (Anti-Pattern 3)
- [ ] Testing behavior, not implementation (Anti-Pattern 4)
- [ ] No time-based assertions (Anti-Pattern 5)
- [ ] One behavior per test (Anti-Pattern 6)
- [ ] Every test has assertions (Anti-Pattern 7)

---

## Key Takeaways

1. **Test behavior, not implementation** — Tests should verify what the code does, not how it does it
2. **Mock external dependencies only** — Don't mock what you're testing
3. **Complete mocks** — Mocks should match real behavior
4. **Reliable tests** — Avoid timing, randomness, external state
5. **Clear failures** — One behavior per test makes debugging easier
6. **Always assert** — Every test must verify something

---

**Version:** 1.0.0  
**Last Updated:** 2024-01-15  
**Related:** `/test`, `/build`, `/codereview`
