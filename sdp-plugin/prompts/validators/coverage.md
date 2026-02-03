# Coverage Validator

AI-based test coverage analyzer that works with any programming language.

## Purpose

Analyze test coverage by reading test and source files, mapping tests to code, and calculating coverage percentage.

## How to Use

```
Ask Claude: "Run coverage validator by:
1. Reading all test files (tests/, test_*, *_test.go, etc.)
2. Reading all source files (src/, main/, lib/, etc.)
3. Mapping each function/class to its corresponding test
4. Calculating coverage: (tested_functions / total_functions) * 100

Report:
- Total coverage percentage
- List of untested functions
- Missing branch coverage
- Verdict: ✅ PASS (≥80%) or ❌ FAIL (<80%)"
```

## Coverage Calculation

### Step 1: Identify Test Files

**Python:** `tests/`, `test_*.py`

**Java:** `src/test/java/`, `*Test.java`

**Go:** `*_test.go`, `*_test.go`

### Step 2: Identify Source Files

**Python:** `src/`, `lib/`, `*.py` (excluding test files)

**Java:** `src/main/java/`, `*.java` (excluding test files)

**Go:** `*.go` (excluding `*_test.go`)

### Step 3: Map Tests to Functions

For each function/method in source files:
- Find corresponding test in test files
- Check if test exists and covers the function
- Mark as tested or untested

### Step 4: Calculate Coverage

```
coverage = (tested_functions / total_functions) * 100
```

## Output Format

```markdown
## Coverage Report

**Total Coverage:** 85%

**Test Files Found:** 3
- tests/test_user.py
- tests/test_service.py
- tests/test_api.py

**Source Files Found:** 5
- src/service.py (10 functions, 8 tested)
- src/models.py (5 functions, 5 tested)
- src/api.py (8 functions, 5 tested)

**Untested Functions:**
- `src/service.py:UserService.delete_user()` (lines 45-50)
- `src/api.py:DataController.export_data()` (lines 120-135)
- `src/models.py:User.validate_email()` (lines 78-82)

**Missing Branch Coverage:**
- `src/auth.py:login()` - error path not tested (line 45)
- `src/service.py:process_data()` - edge case not tested (line 78)

**Verdict:** ✅ PASS (≥80%)
```

## Failure Example

```markdown
## Coverage Report

**Total Coverage:** 65%

**Untested Functions:**
- `src/service.py:UserService.*()` - 5 functions untested
- `src/models.py:User.*()` - 2 functions untested
- `src/api.py:DataController.*()` - 3 functions untested

**Verdict:** ❌ FAIL (<80%)

**Required Actions:**
1. Add tests for UserService methods
2. Add tests for User model methods
3. Add tests for DataController methods
4. Re-run validator
```

## Language-Specific Patterns

### Python

```python
# Source: src/service.py
class UserService:
    def create_user(self, name: str) -> User:
        return User(name=name)

    def delete_user(self, user_id: int) -> bool:
        # ...
```

```python
# Test: tests/test_service.py
def test_create_user():
    user = UserService().create_user("Alice")
    assert user.name == "Alice"  # ✅ Tests create_user()

    # ❌ No test for delete_user() - uncovered
```

### Java

```java
// Source: src/main/java/service/UserService.java
public class UserService {
    public User createUser(String name) {
        return new User(name);
    }

    public boolean deleteUser(int userId) {
        // ...
    }
}
```

```java
// Test: src/test/java/service/UserServiceTest.java
@Test
public void testCreateUser() {
    User user = userService.createUser("Alice");
    assertEquals("Alice", user.getName());  // ✅ Tests createUser()

    // ❌ No test for deleteUser() - uncovered
}
```

### Go

```go
// Source: service/user.go
type UserService struct{}

func (s *UserService) CreateUser(name string) *User {
    return &User{Name: name}
}

func (s *UserService) DeleteUser(userID int) bool {
    // ...
}
```

```go
// Test: service/user_test.go
func TestUserService_CreateUser(t *testing.T) {
    user := UserService{}.CreateUser("Alice")
    if user.Name != "Alice" {
        t.Errorf("expected Alice, got %s", user.Name)
    }  // ✅ Tests CreateUser()

    // ❌ No test for DeleteUser() - uncovered
}
```

## Edge Cases

### Partial Coverage

```python
# Function tested but not all branches
def login(username: str, password: str) -> bool:
    if not authenticate(username, password):
        return False  # ❌ Tested only this path
    return create_session(username)  # ❌ Not tested
```

**Solution:** Add test for successful login path.

### Indirect Coverage

```python
# Class is imported but methods not called
class Helper:
    def utility_method(self):
        pass
```

**Coverage:** 0% if methods never called in tests.

### Test Files Not Executed

```bash
# Test file exists but skipped
@pytest.mark.skip("TODO: fix later")
def test_something():
    pass
```

**Coverage:** Test skipped = not covered.

## Quality Gate

- **PASS:** Coverage ≥80%
- **FAIL:** Coverage <80%

## Tips for Increasing Coverage

1. **Add tests for untested functions**
   - Create test file if missing
   - Add test methods for uncovered functions

2. **Add branch coverage**
   - Test both success and failure paths
   - Test edge cases (null, empty, invalid input)

3. **Test integration points**
   - Test functions that call other functions
   - Test with realistic data

4. **Use coverage tools**
   - Python: `pytest --cov=src/ --cov-report=term-missing`
   - Java: `mvn jacoco:report` (shows missing lines)
   - Go: `go tool cover -html=coverage.out`

## See Also

- [@build Skill](../skills/build.md) - Calls coverage validator
- [@review Skill](../skills/review.md) - Runs coverage validator
- [Quality Gates Reference](../../docs/quality-gates.md) - Coverage criteria
