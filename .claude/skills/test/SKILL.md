---
name: test
description: Contract test generation and validation workflow.
---

# @test - Contract Test Generation

Generate and validate contract tests that define the interface between components.

## Purpose

Contract tests are **immutable specifications** of component interfaces. Once created, they cannot be modified during implementation (`/build` phase).

## Workflow

### Step 1: Analyze Requirements

Read the feature specification or workstream document:
```bash
# User provides feature ID or spec
Read("docs/specs/{feature_id}.md")
# OR
Read("docs/workstreams/backlog/{WS-ID}.md")
```

Extract:
- **Public interfaces** - Functions, classes, APIs that external code depends on
- **Data structures** - Input/output formats, schemas
- **Error conditions** - Expected failures, edge cases
- **Invariants** - Rules that must always hold true

### Step 2: Design Test Contracts

For each interface, define:

**1. Function Signature Test**
```python
def test_function_signature():
    """Contract: Function name and parameters MUST NOT change."""
    # Arrange
    component = Component()

    # Act & Assert
    assert hasattr(component, "method_name")
    # Check parameter count
    import inspect
    sig = inspect.signature(component.method_name)
    assert len(sig.parameters) == 2  # param1, param2
```

**2. Input/Output Contract**
```python
def test_input_output_contract():
    """Contract: Input → Output mapping MUST NOT change."""
    # Given valid input
    result = component.process({"key": "value"})

    # Contract: returns dict with specific fields
    assert isinstance(result, dict)
    assert "result" in result
    assert "status" in result
```

**3. Error Conditions Contract**
```python
def test_error_conditions():
    """Contract: Error behavior MUST NOT change."""
    # Given invalid input
    with pytest.raises(ValueError) as exc_info:
        component.process(None)

    # Contract: specific error message
    assert "cannot be None" in str(exc_info.value)
```

**4. Invariants Contract**
```python
def test_invariants():
    """Contract: Business rules MUST always hold."""
    result = component.calculate(x=10, y=5)

    # Invariant: result must be non-negative
    assert result >= 0
    # Invariant: result must be divisible by x
    assert result % x == 0
```

### Step 3: Create Test File

Generate test file in appropriate location:

**Python:**
```bash
# Location: tests/contract/test_{component}.py
```

**Go:**
```bash
# Location: {package}_test.go with TestContract prefix
```

**Structure:**
```python
"""Contract tests for {Component}.

⚠️ CONTRACT TESTS - DO NOT MODIFY once approved

These tests define the public interface contract.
Changes require explicit review and approval.
"""

class Test{Component}Contract:
    """Contract tests for {Component}."""

    def test_signature(self):
        """Contract: Method signature is stable."""
        # Implementation...

    def test_input_output(self):
        """Contract: Input/output behavior is stable."""
        # Implementation...

    # ... more contract tests
```

### Step 4: Review with Stakeholder

Before implementation, get approval:

```markdown
## Contract Review: {Component}

### Interfaces Covered
- `Component.method(param1, param2)` - 2 tests
- `Component.calculate(x, y)` - 3 tests

### Invariants Defined
- Result must be non-negative
- Result must be divisible by x

### Error Conditions
- Raises ValueError for None input
- Raises TypeError for invalid types

**Approval Required:**
- [ ] Product owner approves interface design
- [ ] Tech lead approves test completeness
- [ ] Security approves error handling

Once approved, contracts are **LOCKED** for /build phase.
```

### Step 5: Mark as Approved

Once approved, add marker:

```python
# CONTRACT APPROVED: 2026-02-06 by @techlead
# Changes require explicit approval
```

## Rules for /build Phase

When `/build` executes:

**✅ ALLOWED:**
- Implement functions to pass contract tests
- Add private helper methods
- Refactor implementation (as long as contracts pass)
- Add new tests for implementation details

**❌ FORBIDDEN:**
- Modify contract test files
- Change function signatures
- Change input/output formats
- Remove error conditions
- Relax invariants

**If interface change is needed:**
1. Stop /build
2. Create new workstream: "Update contract for {Component}"
3. Get explicit approval for contract change
4. Return to /build

## Examples

### Example 1: API Contract

```python
# tests/contract/test_api_client.py

class TestAPIClientContract:
    """Contract tests for APIClient."""

    def test_get_endpoint_signature(self):
        """Contract: get() has stable signature."""
        client = APIClient()
        import inspect
        sig = inspect.signature(client.get)
        params = list(sig.parameters.keys())
        assert params == ["url", "params", "headers"]

    def test_get_returns_response(self):
        """Contract: get() returns Response object."""
        client = APIClient()
        response = client.get("https://api.example.com/data")

        assert hasattr(response, "status_code")
        assert hasattr(response, "json")
        assert hasattr(response, "headers")
```

### Example 2: Data Pipeline Contract

```python
# tests/contract/test_pipeline.py

class TestPipelineContract:
    """Contract tests for DataPipeline."""

    def test_transform_input_output(self):
        """Contract: transform() accepts dict, returns dict."""
        pipeline = DataPipeline()
        input_data = {"records": [{"id": 1}]}

        output = pipeline.transform(input_data)

        # Contract: output is always dict
        assert isinstance(output, dict)
        # Contract: output always has 'records' key
        assert "records" in output
        # Contract: records is always a list
        assert isinstance(output["records"], list)

    def test_transform_handles_empty_input(self):
        """Contract: transform() handles empty records."""
        pipeline = DataPipeline()
        output = pipeline.transform({"records": []})

        assert output["records"] == []
        # Contract: does not raise for empty input
```

## Output

After `/test` completes:

```
✅ Contract tests generated: tests/contract/test_{component}.py
   - 5 interface tests
   - 3 invariant tests
   - 2 error condition tests

📋 Ready for review by: @techlead
🔒 Once approved, contracts are LOCKED for /build

Next step: /build {WS-ID} (implementation without changing contracts)
```

## Integration with Workflow

**Full workflow:**
1. `/design {feature-id}` - Plan architecture
2. `/test {WS-ID}` - Generate and approve contract tests
3. `/build {WS-ID}` - Implement (contracts are immutable)
4. `/review {feature-id}` - Quality check

**Contract change workflow:**
```
Current: /build → Oops, interface needs change
↓
Stop /build
Create new WS: "Update contract for X"
Run /test with new contracts
Get approval
Resume /build
```

## Common Pitfalls

**❌ Don't test implementation details:**
```python
# Bad - tests internal implementation
def test_uses_cache():
    assert component._cache_enabled  # Private field

# Good - tests observable behavior
def test_caching_effect():
    result1 = component.compute(key)
    result2 = component.compute(key)
    assert result2 == result1  # Same result = cache works
```

**❌ Don't make contracts too strict:**
```python
# Bad - overly specific
def test_exact_error_message():
    with pytest.raises(ValueError) as exc:
        component.process(None)
    assert str(exc.value) == "Value cannot be None"  # Too rigid

# Good - flexible but clear
def test_error_message_content():
    with pytest.raises(ValueError) as exc:
        component.process(None)
    assert "cannot be None" in str(exc.value)  # Contains key info
```

**✅ Do focus on stability:**
```python
# Good - stable contract
def test_function_exists():
    assert hasattr(module, "public_function")

def test_acceptable_parameters():
    import inspect
    sig = inspect.signature(module.public_function)
    # Can accept 2-3 parameters (flexible)
    assert 2 <= len(sig.parameters) <= 3
```

## Version

**1.0.0** - Initial /test command for contract test generation
