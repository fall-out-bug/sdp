# Traceability Guide

## Overview

Traceability ensures that every Acceptance Criterion (AC) has a corresponding test that validates it. This creates a direct link from requirements to verification.

## Why Traceability Matters

1. **Requirements Coverage:** Ensures all requirements are tested
2. **Change Impact:** Easy to find what tests break when requirements change
3. **Regression Safety:** Prevents requirements from being accidentally broken
4. **Audit Trail:** Clear evidence that requirements are met

## Traceability Format

### In Workstream Files

```markdown
## Acceptance Criteria

- [ ] AC1: User can login with valid credentials
- [ ] AC2: Invalid credentials are rejected
- [ ] AC3: Session is created after successful login
```

### In Test Files

```python
def test_ac1_user_can_login():
    """AC1: User can login with valid credentials."""
    # Test implementation

def test_ac2_invalid_credentials_rejected():
    """AC2: Invalid credentials are rejected."""
    # Test implementation

def test_ac3_session_created_after_login():
    """AC3: Session is created after successful login."""
    # Test implementation
```

## Naming Convention

**Pattern:** `test_ac{N}_{description}`

**Examples:**
- `test_ac1_user_can_login`
- `test_ac2_invalid_credentials_rejected`
- `test_ac3_session_created`

## Traceability Matrix

Create a matrix to track AC-to-test mappings:

| AC | Description | Test | Status |
|----|-------------|------|--------|
| AC1 | User can login | test_ac1_user_can_login | ✅ PASS |
| AC2 | Invalid creds rejected | test_ac2_invalid_credentials_rejected | ✅ PASS |
| AC3 | Session created | test_ac3_session_created | ❌ FAIL |

## Automated Traceability Check

The `sdp trace check` command automates traceability verification:

```bash
# Check single workstream
sdp trace check 00-001-01

# Check entire feature
sdp trace check --feature F01
```

**Output:**
```
Checking traceability for 00-001-01...

AC1: User can login with valid credentials
  ✅ test_ac1_user_can_login (PASS)

AC2: Invalid credentials are rejected
  ✅ test_ac2_invalid_credentials_rejected (PASS)

AC3: Session is created after successful login
  ❌ No test found

Summary: 2/3 ACs traced (66%)
Status: CHANGES_REQUIRED
```

## Manual Traceability Check

If automated check not available, use grep:

```bash
# Extract ACs from WS file
grep -A 1 "AC[0-9]" docs/workstreams/completed/00-001-01.md

# Find tests
grep -r "def test_ac" tests/

# Check test status
pytest tests/unit/test_auth.py -v
```

## Multiple Tests Per AC

Some ACs require multiple tests:

**AC:** User can login with valid credentials

**Tests:**
- `test_ac1_user_can_login_with_email`
- `test_ac1_user_can_login_with_username`
- `test_ac1_user_can_login_case_insensitive`

All tests should reference the same AC in docstring.

## Edge Cases and AC Coverage

Beyond basic AC tests, add edge case tests:

```python
def test_ac1_user_can_login():
    """AC1: User can login with valid credentials."""
    # Basic happy path

def test_ac1_login_with_trailing_spaces():
    """AC1 edge case: Email with trailing spaces."""
    # Edge case

def test_ac1_login_case_insensitive():
    """AC1 edge case: Case insensitive email."""
    # Edge case
```

## Orphaned Tests

Tests without AC references are "orphaned":

```python
def test_password_hashing():
    """Test password is hashed before storage."""
    # No AC reference - orphaned!
```

**Fix:** Link to relevant AC or create new AC if needed.

## Updating Traceability

When AC changes:

1. Update AC in WS file
2. Update test docstring
3. Update test implementation if needed
4. Re-run traceability check

## Tools

### Extract ACs

```bash
# From WS file
grep -E "AC[0-9]+:" docs/workstreams/completed/00-001-01.md
```

### Find Tests

```bash
# All AC tests
grep -r "def test_ac" tests/

# Specific AC
grep -r "AC1" tests/
```

### Run AC Tests

```bash
# Run all AC tests
pytest -k "test_ac" -v

# Run tests for specific AC
pytest -k "test_ac1" -v
```

## Best Practices

1. **One AC, At Least One Test:** Every AC needs at least one test
2. **Explicit References:** Use AC number in test name and docstring
3. **Test Independence:** Each test should be runnable independently
4. **Clear Assertions:** Assert on AC outcome, not implementation details
5. **Update Together:** When AC changes, update tests immediately

## Example: Full Traceability

**Workstream File:**
```markdown
## Acceptance Criteria

- [x] AC1: User can login with valid credentials
- [x] AC2: Invalid credentials are rejected
- [x] AC3: Session is created after successful login
- [x] AC4: Session has 30-minute timeout
```

**Test File:**
```python
def test_ac1_user_can_login():
    """AC1: User can login with valid credentials."""
    result = auth_service.login("user@example.com", "password123")
    assert result.success is True

def test_ac2_invalid_credentials_rejected():
    """AC2: Invalid credentials are rejected."""
    result = auth_service.login("user@example.com", "wrongpass")
    assert result.success is False

def test_ac3_session_created_after_login():
    """AC3: Session is created after successful login."""
    result = auth_service.login("user@example.com", "password123")
    assert result.session is not None

def test_ac4_session_timeout():
    """AC4: Session has 30-minute timeout."""
    session = create_session()
    assert session.timeout_minutes == 30
```

**Traceability Matrix:**

| AC | Test | Status |
|----|------|--------|
| AC1 | test_ac1_user_can_login | ✅ PASS |
| AC2 | test_ac2_invalid_credentials_rejected | ✅ PASS |
| AC3 | test_ac3_session_created_after_login | ✅ PASS |
| AC4 | test_ac4_session_timeout | ✅ PASS |

**Result:** 100% traceability ✅
