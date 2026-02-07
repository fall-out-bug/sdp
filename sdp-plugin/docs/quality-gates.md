# SDP Quality Gates

**Binary pass/fail criteria** for workstream completion and quality assurance.

**Language-Agnostic:** Python, Java, Go supported

---

## Overview

Quality gates are **MUST PASS** criteria. Each gate has:
- ✅ **PASS** - Criteria met, workstream approved
- ❌ **FAIL** - Criteria not met, must fix before completion
- 🔧 **Measurement** - How to check the criterion (language-specific)
- 📋 **Example** - What pass/fail looks like

---

## 1. AI-Readiness Gate

### Criteria

| # | Rule | Limit | Python | Java | Go |
|---|------|-------|--------|------|-----|
| 1.1 | Max lines per file | < 200 LOC | `find src/ -name "*.py" -exec wc -l {} + \| awk '$1 > 200'` | `find src/ -name "*.java" -exec wc -l {} + \| awk '$1 > 200'` | `find . -name "*.go" -exec wc -l {} + \| awk '$1 > 200'` |
| 1.2 | Max cyclomatic complexity | < 10 | `radon cc src/ -a -s` | Manual review | `gocyclo -over 10 .` |
| 1.3 | Type hints coverage | 100% | `mypy src/ --strict` | `javac -Xlint:all` | `go vet ./...` |
| 1.4 | Function length | < 50 lines | `radon cc src/ -a -s` | Manual review | `gocyclo -over 15 .` |

### PASS Criteria (Python)

```bash
# All checks must return empty (no violations)
find src/ -name "*.py" -exec wc -l {} + | awk '$1 > 200'  # No output
mypy src/ --strict  # Success: 0 errors
```

### PASS Criteria (Java)

```bash
# All checks must return empty (no violations)
find src/main -name "*.java" -exec wc -l {} + | awk '$1 > 200'  # No output
javac -Xlint:all  # Success: 0 errors
```

### PASS Criteria (Go)

```bash
# All checks must return empty (no violations)
find . -name "*.go" -exec wc -l {} + | awk '$1 > 200'  # No output
go vet ./...  # Success: 0 warnings
```

### FAIL Examples

❌ **File too long (any language):**
```
service.java (250 lines) - FAIL
service.go (250 lines) - FAIL
service.py (250 lines) - FAIL
```

✅ **Correct:**
```
service.java (180 lines) - PASS (split into multiple functions)
service.go (180 lines) - PASS (split into multiple functions)
service.py (180 lines) - PASS (split into multiple functions)
```

---

## 2. Test Coverage Gate

### Criteria

| Language | Threshold | Command |
|----------|-----------|---------|
| **Python** | ≥80% | `pytest tests/unit/ --cov=src/ --cov-fail-under=80` |
| **Java** | ≥80% | `mvn verify` (JaCoCo) or `gradle test jacocoTestReport` |
| **Go** | ≥80% | `go test -coverprofile=coverage.out && go tool cover -func=coverage.out \| grep total` |

### PASS Examples

**Python:**
```bash
$ pytest tests/ --cov=src/ --cov-report=term-missing
---------- coverage: platform linux, python 3.10 ----------
Name                          Stmts   Miss  Cover   Missing
---------------------------------------------------------
src/service.py                   50      5    90%    23-27
src/models.py                    30      2    93%    45, 67
---------------------------------------------------------
TOTAL                            80      7    91%    ✅ PASS (≥80%)
```

**Java:**
```bash
$ mvn verify
[INFO] --- jacoco-maven-plugin:0.8.11:check (default) @ sdp ---
[INFO] Loading execution data file /target/jacoco.exec
[INFO] Analyzed bundle 'sdp' with 80 classes
[INFO] All coverage checks have been met.
[INFO]   Rule 0: CoveredRatio = 0.85 (required minimum = 0.80) ✅ PASS
```

**Go:**
```bash
$ go test -coverprofile=coverage.out ./...
ok      github.com/user/sdp/service    0.123s   coverage: 85.4% of statements
ok      github.com/user/sdp/models     0.045s   coverage: 92.1% of statements

$ go tool cover -func=coverage.out | grep total
github.com/user/sdp/    87.3%  ✅ PASS (≥80%)
```

### FAIL Examples

❌ **Coverage below 80%:**
```
Python: 75% - FAIL
Java: 72% - FAIL
Go: 68% - FAIL
```

---

## 3. Clean Architecture Gate

### Criteria

| Layer | Python Check | Java Check | Go Check |
|-------|-------------|------------|----------|
| **Domain** | `grep -r "from.*infrastructure" src/domain/` (empty) | `grep -r "import.*infrastructure" src/main/java/domain/` (empty) | `grep -r "infrastructure" domain/` (empty) |
| **Application** | `grep -r "from.*presentation" src/application/` (empty) | `grep -r "import.*presentation" src/main/java/application/` (empty) | `grep -r "presentation" application/` (empty) |

### PASS Criteria

```bash
# Domain layer must have NO external dependencies
# Python:
grep -r "from.*infrastructure" src/domain/  # No output
grep -r "from.*application" src/domain/  # No output

# Java:
grep -r "import.*infrastructure" src/main/java/domain/  # No output
grep -r "import.*presentation" src/main/java/domain/  # No output

# Go:
grep -r "infrastructure" domain/  # No output
grep -r "presentation" domain/  # No output
```

### FAIL Examples

❌ **Domain layer violation (Python):**
```python
# src/domain/entities/user.py
from src.infrastructure.persistence import Database  # FAIL - domain imports infrastructure
```

✅ **Correct:**
```python
# src/domain/entities/user.py (no external imports) - PASS
class User:
    def __init__(self, name: str):
        self.name = name
```

❌ **Domain layer violation (Java):**
```java
// com.example.domain.entity.User.java
import com.example.infrastructure.persistence.Database;  // FAIL
```

✅ **Correct:**
```java
// com.example.domain.entity.User.java (no external imports) - PASS
package com.example.domain.entity;

public class User {
    private String name;
}
```

---

## 4. Error Handling Gate

### Criteria

| Pattern | Python | Java | Go |
|---------|--------|------|-----|
| **Bare except** | `except:` (forbidden) | `catch(Exception e)` (forbidden unless logged) | `recover()` (must check error) |
| **Swallowed errors** | `except: pass` (forbidden) | Empty catch block (forbidden) | Ignored returns `func(), _` (forbidden) |

### PASS Examples

**Python:**
```python
# ✅ PASS - Specific exception
try:
    risky_operation()
except ValueError as e:
    logger.error(f"Invalid value: {e}")
    raise
```

**Java:**
```java
// ✅ PASS - Specific exception with logging
try {
    riskyOperation();
} catch (ValueException e) {
    logger.error("Invalid value", e);
    throw e;
}
```

**Go:**
```go
// ✅ PASS - Error checked
data, err := fetchData()
if err != nil {
    log.Printf("Fetch failed: %v", err)
    return err
}
```

### FAIL Examples

❌ **Python:**
```python
# FAIL - Bare except
try:
    risky_operation()
except:
    pass  # Swallows all errors!
```

❌ **Java:**
```java
// FAIL - Catch-all exception
try {
    riskyOperation();
} catch (Exception e) {
    // Error not logged or re-thrown
}
```

❌ **Go:**
```go
// FAIL - Error ignored
data, _ := fetchData()  // Error lost
```

---

## 5. Type Safety Gate

### Criteria

| Language | Requirement | Command |
|----------|-------------|---------|
| **Python** | 100% type hints | `mypy src/ --strict` (0 errors) |
| **Java** | No raw types | `javac -Xlint:unchecked` (0 warnings) |
| **Go** | No unsafe conversions | `go vet ./...` (0 warnings) |

### PASS Examples

**Python:**
```python
# ✅ PASS - Full type hints
def process(data: str, count: int) -> str:
    return data * count
```

**Java:**
```java
// ✅ PASS - Generic types
List<String> names = new ArrayList<>();  // Type-safe
```

**Go:**
```go
// ✅ PASS - Explicit types
func process(data string, count int) string {
    return strings.Repeat(data, count)
}
```

### FAIL Examples

❌ **Python:**
```python
# FAIL - Missing type hints
def process(data, count):  # No types
    return data * count
```

❌ **Java:**
```java
// FAIL - Raw type
List names = new ArrayList();  // Raw type
```

---

## Enforcement

### Automatic Enforcement

Quality gates are enforced via:
1. **AI Validators** - Analyze code and report violations
2. **Pre-commit Hooks** - Optional, run language-specific tools
3. **Review Skill** - `@review` runs all validators

### Manual Enforcement

If AI validation is insufficient:
- Run language-specific tools manually
- Review violations in output
- Fix code and re-validate

---

## Quick Reference

### Quality Gate Checklist

- [ ] Coverage ≥80% (pytest/mvn test/go test)
- [ ] Type hints complete (mypy/javac/go vet)
- [ ] No bare exceptions (code review)
- [ ] File size <200 LOC
- [ ] Architecture violations = 0
- [ ] No TODO/FIXME without WS

### Language-Specific Commands

**Python:**
```bash
pytest tests/ --cov=src/ --cov-fail-under=80
mypy src/ --strict
ruff check src/
```

**Java:**
```bash
mvn verify  # Runs tests + JaCoCo
mvn checkstyle:check
```

**Go:**
```bash
go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out | grep total
go vet ./...
golint ./...
```
