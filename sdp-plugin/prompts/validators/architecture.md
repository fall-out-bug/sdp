# Architecture Validator

AI-based Clean Architecture enforcer that validates layer separation.

## Purpose

Check dependency graph for Clean Architecture violations by parsing imports and ensuring proper layer isolation.

## How to Use

```
Ask Claude: "Run architecture validator by:
1. Parsing imports from all source files
2. Mapping files to layers (domain/, application/, infrastructure/, presentation/)
3. Checking for violations:
   - ❌ domain/ imports infrastructure/
   - ❌ domain/ imports application/
   - ❌ application/ imports presentation/
4. Report violations with file paths and line numbers

Allowed dependencies:
- ✅ infrastructure/ imports domain/
- ✅ application/ imports domain/
- ✅ presentation/ imports application/
- ✅ Any layer imports external libraries

Report:
- Number of violations
- File paths and line numbers
- Verdict: ✅ PASS (no violations) or ❌ FAIL (violations found)"
```

## Layer Definitions

### Domain Layer
- **Purpose:** Business logic and entities
- **Location:** `src/domain/`, `domain/`, `entities/`
- **Dependencies:** NONE (only external libraries)
- **Forbidden:** Cannot import infrastructure, application, presentation

### Application Layer
- **Purpose:** Use cases and business workflows
- **Location:** `src/application/`, `application/`, `services/`
- **Dependencies:** Can import domain/
- **Forbidden:** Cannot import infrastructure/, presentation/

### Infrastructure Layer
- **Presentation:** UI, controllers, APIs
- **Location:** `src/infrastructure/`, `src/presentation/`, `api/`
- **Dependencies:** Can import domain/, application/
- **Forbidden:** None (can import any layer)

### Presentation Layer
- **Location:** `src/presentation/`, `controllers/`, `api/`
- **Dependencies:** Can import application/
- **Forbidden:** Cannot import domain/, infrastructure/

## Import Analysis

### Python

```python
# ✅ CORRECT: Infrastructure imports Domain
from src.domain.entities import User  # OK

# ❌ VIOLATION: Domain imports Infrastructure
from src.infrastructure.persistence import Database  # FAIL
```

### Java

```java
// ✅ CORRECT: Service imports Domain
import com.example.domain.entities.User;  // OK

// ❌ VIOLATION: Domain imports Infrastructure
import com.example.infrastructure.persistence.Database;  // FAIL
```

### Go

```go
// ✅ CORRECT: Infrastructure imports Domain
import "github.com/user/project/domain"  // OK

// ❌ VIOLATION: Domain imports Infrastructure
import "github.com/user/project/infrastructure"  // FAIL
```

## Output Format

### PASS Example

```markdown
## Architecture Report

**Violations Found:** 0

**Dependencies Analyzed:** 25 files

**Allowed Dependencies (OK):**
- ✅ src/infrastructure/repositories.py → src/domain/entities.py
- ✅ src/application/services.py → src/domain/entities.py
- ✅ src/presentation/controllers.py → src/application/services.py

**Verdict:** ✅ PASS (no violations)
```

### FAIL Example

```markdown
## Architecture Report

**Violations Found:** 3

1. ❌ `src/domain/entities/user.py` (line 5)
   ```python
   from src.infrastructure.persistence import Database
   ```
   **Rule:** Domain must not import infrastructure.
   **Fix:** Move persistence logic to application layer, use repository pattern

2. ❌ `src/domain/entities/user.py` (line 12)
   ```python
   from src.application.services import UserService
   ```
   **Rule:** Domain must not import application.
   **Fix:** Inject UserService as dependency, use inversion of control

3. ❌ `src/application/auth.py` (line 8)
   ```python
   from src.presentation.controllers import UserController
   ```
   **Rule:** Application must not import presentation.
   **Fix:** Use callbacks or events instead of direct controller access

**Allowed Dependencies (OK):**
- ✅ src/infrastructure/repositories.py → src/domain/entities.py
- ✅ src/application/services.py → src/domain/entities.py
- ✅ src/presentation/controllers.py → src/application/services.py

**Verdict:** ❌ FAIL (violations detected)

**Required Actions:**
1. Fix domain/infrastructure violation
2. Fix domain/application violation
3. Fix application/presentation violation
4. Re-run validator
```

## Detection Patterns

### Python Imports

```python
# Parse these patterns:
from src.{layer}.{module} import {Class}
from src.{layer} import {module}
import src.{layer}.{module}
```

### Java Imports

```java
// Parse these patterns:
import com.example.{layer}.{module}.{Class};
import com.example.{layer}.{module};
import static com.example.{layer}.{Class}.{method};
```

### Go Imports

```go
// Parse these patterns:
import "github.com/user/project/{layer}"
import "github.com/user/project/{layer}/{package}"
```

## Circular Dependencies

Check for circular dependencies across layers:

```bash
# Python
# domain → application → domain (FAIL)

# Java
# com.example.domain → com.example.application → com.example.domain (FAIL)

# Go
# github.com/user/domain → github.com/user/application → github.com/user/domain (FAIL)
```

## Quality Gate

- **PASS:** Zero architecture violations
- **FAIL:** Any violation found

## Refactoring Tips

### Fix Domain Layer Violations

**Problem:** Domain imports infrastructure

```python
# ❌ BEFORE (violation)
class User:
    def __init__(self):
        self.db = Database()  # Imports infrastructure
```

```python
# ✅ AFTER (fixed)
class User:
    def __init__(self, user_repository):
        self.user_repository = user_repository  # Injected dependency
```

### Fix Application Layer Violations

**Problem:** Application imports presentation

```python
# ❌ BEFORE (violation)
class AuthService:
    def __init__(self):
        self.controller = UserController()  # Imports presentation
```

```python
# ✅ AFTER (fixed)
class AuthService:
    def __init__(self, user_created_callback):
        self.user_created_callback = user_created_callback  # Callback/event
```

## Edge Cases

### External Libraries

```python
# ✅ OK - External libraries don't count
from sqlalchemy import create_engine  # OK
from flask import Flask  # OK
from numpy import np  # OK
```

### Test Files

```python
# ✅ OK - Tests can import anything for testing
from tests.fakes import FakeDatabase  # OK
```

### Factory Pattern

```python
# ✅ OK - Factories can create infrastructure
class UserFactory:
    @staticmethod
    def create_with_db():
        return User(db=Database())  # OK - factory method
```

## See Also

- [@build Skill](../skills/build.md) - Calls architecture validator
- [@review Skill](../skills/review.md) - Runs architecture validator
- [Quality Gates Reference](../../docs/quality-gates.md) - Architecture criteria
