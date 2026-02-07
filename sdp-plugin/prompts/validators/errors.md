# Error Handling Validator

AI-based error handling auditor that finds unsafe exception handling patterns.

## Purpose

Detect error handling anti-patterns that swallow errors or fail to handle them properly.

## How to Use

```
Ask Claude: "Run error handling validator by:
1. Searching for unsafe error handling patterns
2. Checking if errors are logged or re-raised
3. Verifying specific exception types (not catch-all)

Patterns to find:
- Python: bare except, except: pass, except Exception (without logging)
- Java: empty catch blocks, catch(Exception e) without logging
- Go: ignored errors (func(), _), recover() without error check

For each match:
- Is error logged with context?
- Is error re-raised if critical?
- Is specific exception type (not catch-all)?

Report:
- Number of violations
- File paths and line numbers
- Verdict: ✅ PASS (no violations) or ❌ FAIL (violations found)"
```

## Patterns to Detect

### Python

**❌ Bare except:**
```python
try:
    risky()
except:
    pass  # Swallows ALL errors
```

**❌ except Exception without logging:**
```python
try:
    risky()
except Exception as e:
    return  # Error lost
```

**✅ Correct:**
```python
try:
    risky()
except ValueError as e:
    logger.error(f"Invalid value: {e}")
    raise
except Exception as e:
    logger.error(f"Unexpected error: {e}")
    raise
```

### Java

**❌ Empty catch block:**
```java
try {
    risky();
} catch (Exception e) {
    // Empty - error lost
}
```

**❌ Catch-all without logging:**
```java
try {
    risky();
} catch (Exception e) {
    return;  // Error not logged
}
```

**✅ Correct:**
```java
try {
    risky();
} catch (ValueException e) {
    logger.error("Invalid value", e);
    throw e;
} catch (Exception e) {
    logger.error("Unexpected error", e);
    throw new RuntimeException(e);
}
```

### Go

**❌ Ignored error:**
```go
data, _ := fetchData()  // Error lost
```

**❌ recover without check:**
```go
defer func() {
    if r := recover(); r != nil {
        // r is error, but not checked
    }
}()
```

**✅ Correct:**
```go
data, err := fetchData()
if err != nil {
    log.Printf("Fetch failed: %v", err)
    return err
}
```

## Output Format

### PASS Example

```markdown
## Error Handling Report

**Violations:** 0

**Files Analyzed:** 15

**Safe Error Handling (OK):**
- ✅ src/service.py:52 (catches ValueError, logs, re-raises)
- ✅ src/api.py:78 (catches specific exception, handles)
- ✅ src/auth.py:123 (propagates errors correctly)

**Verdict:** ✅ PASS (no violations)
```

### FAIL Example

```markdown
## Error Handling Report

**Violations:** 5

1. ❌ `src/service.py:process_data()` (lines 45-47)
   ```python
   try:
       risky_operation()
   except:
       pass  # Swallows all errors!
   ```
   **Severity:** HIGH
   **Fix:** Catch specific exception, log error, re-raise if critical

2. ❌ `src/api.py:handle_request()` (line 78)
   ```python
   try:
       process_request()
   except Exception as e:
       return  # Error not logged
   ```
   **Severity:** MEDIUM
   **Fix:** Log error before returning

3. ❌ `src/auth.py:login()` (line 123)
   ```python
   try:
       authenticate_user()
   except Exception:
       pass  # Security risk: auth errors hidden
   ```
   **Severity:** CRITICAL
   **Fix:** Log security events, re-raise auth failures

4. ❌ `service/data.go:LoadData()` (line 56)
   ```go
   data, _ := fetchData()  # Error ignored
   ```
   **Severity:** HIGH
   **Fix:** Check error and handle

5. ❌ `service/auth.go:Recover()` (line 78)
   ```go
   defer func() {
       if r := recover(); r != nil {
       // Error not checked
       log.Printf("Recovered: %v", r)
       // What is r? Is it the actual error?
       }
   }()
   ```
   **Severity:** MEDIUM
   **Fix:** Check error type, handle appropriately

**Verdict:** ❌ FAIL (violations detected)

**Required Actions:**
1. Fix bare except in service.py:process_data()
2. Log exception in api.py:handle_request()
3. Fix auth.py:login() - log security events
4. Check error in data.go:LoadData()
5. Improve error handling in auth.go:Recover()

**Priority:** CRITICAL issues must be fixed immediately
```

## Severity Classification

| Severity | Description | Example |
|----------|-------------|---------|
| **CRITICAL** | Security or data loss risk | Auth errors hidden, data corruption |
| **HIGH** | Error lost, silent failures | Bare except, ignored errors |
| **MEDIUM** | Error not logged | Empty catch without context |
| **LOW** | Poor error messages | Generic exception type |

## Common Mistakes

### 1. Bare Except (Python)

**Problem:** Catches all exceptions including system interrupts

```python
try:
    work()
except:  # ❌ Catches KeyboardInterrupt too!
    pass
```

**Fix:**
```python
try:
    work()
except Exception:  # ✅ Excludes KeyboardInterrupt
    pass
```

### 2. Except: Pass (Python)

**Problem:** Silently fails, no record of error

```python
try:
    critical_operation()
except:
    pass  # ❌ No logging, no re-raise
```

**Fix:**
```python
try:
    critical_operation()
except Exception as e:
    logger.error(f"Critical op failed: {e}")
    raise  # ✅ Re-raise if critical
```

### 3. Ignored Error Returns (Go)

**Problem:** Error return value ignored

```go
data, err := fetchData()
// ❌ What if err != nil?
process(data)
```

**Fix:**
```go
data, err := fetchData()
if err != nil {
    log.Printf("Fetch failed: %v", err)
    return err
}
process(data)  // ✅ Only runs if data valid
```

### 4. Empty Catch Blocks (Java)

**Problem:** Exception caught but ignored

```java
try {
    operation()
} catch (Exception e) {
    // TODO: handle this
}
```

**Fix:**
```java
try {
    operation()
} catch (Exception e) {
    logger.error("Operation failed", e);
    throw new RuntimeException(e);
}
```

## Security Considerations

### Authentication Errors

**❌ CRITICAL:**
```python
def login(username, password):
    try:
        return authenticate(username, password)
    except:
        return None  # ❌ Hacks can test passwords silently!
```

**✅ CORRECT:**
```python
def login(username, password):
    try:
        return authenticate(username, password)
    except AuthenticationError as e:
        logger.warning(f"Auth failed for {username}: {e}")
        raise  # ✅ Log and re-raise
```

### Database Errors

**❌ HIGH:**
```python
def save_user(user):
    try:
        db.save(user)
    except:
        return "Error"  # ❌ User doesn't know if save succeeded
```

**✅ CORRECT:**
```python
def save_user(user):
    try:
        db.save(user)
        return "Success"
    except DatabaseError as e:
        logger.error(f"Failed to save user: {e}")
        raise  # ✅ Log and re-raise
```

## Testing Error Handling

### Python Test Example

```python
import pytest
import logging

def test_error_logging(caplog):
    """Test that errors are logged."""
    with caplog.at_level(logging.ERROR):
        result = risky_operation()
        assert result == "error"
        assert "Operation failed" in caplog.text

def test_specific_exception():
    """Test that specific exceptions are used."""
    with pytest.raises(ValueError):
        raise ValueError("test")
```

### Java Test Example

```java
@Test
public void testErrorLogging() {
    // Test that exceptions are logged
    Logger logger = LoggerFactory.getLogger(MyClass.class);
    appender = ListAppender.createAppender();
    logger.addAppender(appender);

    try {
        riskyOperation();
    } catch (Exception e) {
        String logged = appender.getEvents().get(0).getFormattedMessage();
    }
}
```

## Quality Gate

- **PASS:** No unsafe error handling patterns
- **FAIL:** Any violation found (especially CRITICAL severity)

## See Also

- [@build Skill](../skills/build.md) - Calls error validator
- [@review Skill](../skills/review.md) - Runs error validator
- [Quality Gates Reference](../../docs/quality-gates.md) - Error handling criteria
