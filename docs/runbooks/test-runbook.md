# /test Runbook

**Purpose:** Generate or approve tests as contract for workstream implementation.

**Prerequisites:**
- `/design` completed (Interface section exists in WS)
- WS file in `tools/hw_checker/docs/workstreams/backlog/`
- T0 tier capability (architectural decisions only)

## When to Use

- After `/design` is complete
- Before `/build` (implementation phase)
- Need to create/validate test contract
- Want contract-driven development

## Workflow

### Phase 1: Verify Pre-Conditions

```bash
# 1. Check Interface section exists
grep "### Interface" tools/hw_checker/docs/workstreams/backlog/WS-060-01.md

# 2. Verify function signatures with types
grep -A 10 "### Interface" tools/hw_checker/docs/workstreams/backlog/WS-060-01.md

# 3. Check docstrings with Args/Returns/Raises
grep "Args:" tools/hw_checker/docs/workstreams/backlog/WS-060-01.md
grep "Returns:" tools/hw_checker/docs/workstreams/backlog/WS-060-01.md
grep "Raises:" tools/hw_checker/docs/workstreams/backlog/WS-060-01.md
```

### Phase 2: Execute /test Command

**In Claude Code:**
```
/test WS-060-01
```

**In Cursor:**
```
/test WS-060-01
```

**In OpenCode:**
```
/test WS-060-01
```

### Phase 3: What /test Does

The `/test` command (T0 tier only) will:

1. **Read context:**
   - Read WS file
   - Read PROJECT_MAP.md

2. **Verify /design complete:**
   - Check Interface section exists
   - Verify function signatures with types
   - Check docstrings with Args/Returns/Raises

3. **Create/Approve tests:**
   - If no tests exist → create full test set
   - If tests exist → verify completeness, supplement if needed
   - Ensure tests cover:
     - All public methods
     - Edge cases
     - Error conditions
     - Boundary conditions

4. **Update WS file:**
   - Add "Tests (DO NOT MODIFY)" section
   - Ensure tests are executable (fail with `NotImplementedError`)
   - Mark contract as read-only for T2/T3 models

5. **Verify completion criteria:**
   - Tests are executable
   - Tests fail before implementation (RED)
   - Tests pass after implementation (GREEN)
   - Coverage targets defined

## Contract Principle

**Tests = Single Source of Truth for Behavior**

### Contract Rules:

1. **Tests NOT changed in /build** — only implementation bodies
2. **Tests define behavior** — if test requires X, implementation must do X
3. **Tests are executable** — `pytest path/to/test.py` must run
4. **Tests fail before implementation** — `NotImplementedError` in functions → RED
5. **Tests green after implementation** — /build makes them GREEN

## Test Contract Example

### Before /test

```markdown
## WS-060-01: Create Submission

### Interface

```python
def create_submission(user_id: int, hw_id: int, files: list[dict]) -> Submission:
    """
    Create a new homework submission.

    Args:
        user_id: User ID
        hw_id: Homework ID
        files: List of file objects

    Returns:
        Submission object

    Raises:
        ValidationError: If files invalid
        DuplicateSubmissionError: If already submitted
    """
```

## Steps

1. Create Submission entity
2. Implement create_submission function
3. Add validation
4. Test with sample data
```

### After /test

```markdown
### Interface

```python
def create_submission(user_id: int, hw_id: int, files: list[dict]) -> Submission:
    """
    Create a new homework submission.

    Args:
        user_id: User ID
        hw_id: Homework ID
        files: List of file objects

    Returns:
        Submission object

    Raises:
        ValidationError: If files invalid
        DuplicateSubmissionError: If already submitted
    """
```

### Tests (DO NOT MODIFY)

```python
def test_create_submission_success():
    """Test successful submission creation."""
    user_id = 1
    hw_id = 1
    files = [{"name": "main.py", "content": "print('hello')"}]

    submission = create_submission(user_id, hw_id, files)

    assert submission.user_id == user_id
    assert submission.hw_id == hw_id
    assert len(submission.files) == len(files)

def test_create_submission_duplicate():
    """Test duplicate submission raises error."""
    user_id = 1
    hw_id = 1
    files = [{"name": "main.py", "content": "print('hello')"}]

    create_submission(user_id, hw_id, files)

    with pytest.raises(DuplicateSubmissionError):
        create_submission(user_id, hw_id, files)

def test_create_submission_invalid_files():
    """Test invalid files raise ValidationError."""
    user_id = 1
    hw_id = 1
    files = []  # Empty files

    with pytest.raises(ValidationError):
        create_submission(user_id, hw_id, files)

def test_create_submission_file_size_limit():
    """Test files exceed size limit."""
    user_id = 1
    hw_id = 1
    files = [{"name": "large.py", "content": "x" * 100_000_000}]

    with pytest.raises(ValidationError, match="file size exceeds limit"):
        create_submission(user_id, hw_id, files)
```

## Steps

1. Create Submission entity
2. Implement create_submission function
3. Add validation
4. Test with sample data
```

## Capability Tiers

| Tier | Capabilities | When to Use | Contract Access |
|-------|-------------|-------------|-----------------|
| **T0** | Architectural decisions, contract creation | `/test` command (always T0) | ✅ Read/Write |
| T1 | Basic implementation | Strong models | ✅ Read only |
| T2 | Refactoring with constraints | Medium models | ✅ Read only |
| T3 | Fills in implementation | Weak models | ✅ Read only |

**For T2/T3:**
- Contract (Tests section) is READ-ONLY
- Cannot modify Interface or Tests
- Only implement function bodies to satisfy contract

## Verification After /test

```bash
# 1. Check Tests section exists
grep "### Tests (DO NOT MODIFY)" tools/hw_checker/docs/workstreams/backlog/WS-060-01.md

# 2. Verify tests are executable
# Create test file
cat tools/hw_checker/docs/workstreams/backlog/WS-060-01.md | \
  sed -n '/```python/,/```/p' > /tmp/test_temp.py

# Run tests (should fail with NotImplementedError)
cd tools/hw_checker
pytest /tmp/test_temp.py -v

# Expected: FAILED - NotImplementedError
```

## Integration with Other Commands

```
/design idea-slug
    ↓
/test WS-060-01    # Generate test contract (T0)
    ↓
/build WS-060-01   # Implement (tests turn GREEN)
    ↓
/test WS-060-02
    ↓
/build WS-060-02
```

## Troubleshooting

### /test not found

**Problem:** Command not available

**Solution:**
- Claude Code: Check `.claude/skills/test/SKILL.md` exists
- Cursor: Check `.cursor/commands/test.md` exists
- OpenCode: Check `.opencode/commands/test.md` exists

### Interface section missing

**Problem:** /test complains "Interface section not found"

**Solution:**
```bash
# Run /design first
/design idea-lms-integration
```

### Tests not executable

**Problem:** Tests fail to run with syntax errors

**Solution:**
1. Check test code format
2. Verify imports are correct
3. Check pytest is installed
4. Run tests manually: `pytest tests/unit/test_*.py`

### /test tries to modify existing tests

**Problem:** /test changes tests that should be read-only

**Solution:**
1. Check if WS has T2/T3 tier specified
2. Verify contract read-only flag is set
3. Manually mark Tests section as read-only

## Best Practices

1. **Run after /design:** Don't skip /design, it creates Interface section
2. **Review tests:** Always review generated tests before /build
3. **Test edge cases:** Ensure tests cover error conditions and boundaries
4. **Mark read-only:** Always mark Tests section as DO NOT MODIFY
5. **T0 tier only:** /test is always T0, architectural decisions only

## References

- Master prompt: `sdp/prompts/commands/test.md`
- Contract-Driven WS spec: `tools/hw_checker/docs/workstreams/completed/2026-01/WS-410-01-contract-driven-ws-spec.md`
- Capability-tier validator: `sdp/src/sdp/validators/capability_tier.py`
