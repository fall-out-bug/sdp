# Complexity Validator

AI-based code complexity analyzer that identifies overly complex code.

## Purpose

Identify functions/methods that exceed complexity thresholds:
- Cyclomatic complexity >10
- Function length >200 lines
- Nesting depth >4

## How to Use

```
Ask Claude: "Run complexity validator by:
1. Counting lines per file
2. Calculating cyclomatic complexity per function
3. Checking nesting depth per function

Thresholds:
- Complexity <10: ✅ OK
- Complexity 10-20: ⚠️ Warning (consider refactoring)
- Complexity >20: ❌ FAIL (too complex)
- Lines >200: ❌ FAIL
- Nesting >4: ❌ FAIL

Report:
- Files exceeding thresholds
- Complexity scores for each function
- Nesting depth violations
- Refactoring recommendations
- Verdict: ✅ PASS (no violations) or ❌ FAIL (violations found)"
```

## Complexity Metrics

### Cyclomatic Complexity

Number of independent paths through code:

```python
# Complexity = 1 (single path)
def simple_function(x):  # ✅ Complexity: 1
    return x + 1

# Complexity = 4 (2 binary decisions + 1 path)
def medium_function(x, y, z):  # ✅ Complexity: 4
    if x > 0:
        if y > 0:
            return 1
        else:
            return 2
    else:
        return 3

# Complexity = 15 (many branches)  # ❌ FAIL
def complex_function(a, b, c, d, e, f, g):
    if a > 0:
        if b > 0:
            if c > 0:
                return 1
            elif d > 0:
                if e > 0:
                    if f > 0:
                        return 2
                    else:
                        return 3
                else:
                    return 4
            else:
                return 5
        else:
            return 6
    else:
        return 7
```

### Function Length

Count non-empty, non-comment lines:

```python
# 50 lines - OK
def medium_function():
    """Do something."""
    # ... 40 lines of code ...
    return result
```

```python
# 250 lines - ❌ FAIL
def very_long_function():
    """Too much logic."""
    # ... 247 lines of code ...
    return result
```

**Threshold:** <200 LOC per function

### Nesting Depth

Maximum level of nested control structures:

```python
# Depth 2 - OK
def ok_function():
    if condition:
        if nested:
            return  # ✅ Depth 2
```

```python
# Depth 5 - ❌ FAIL
def deep_function():
    if level1:
        if level2:
            if level3:
                if level4:
                    if level5:  # ❌ Too deep!
                        return
```

**Threshold:** ≤4 nesting depth

## Output Format

### PASS Example

```markdown
## Complexity Report

**Violations:** 0

**Files Analyzed:** 10

**Complex Functions (OK):**
- ✅ src/service.py:UserService.create_user() - Complexity: 3, 15 LOC (OK)
- ✅ src/service.py:UserService.delete_user() - Complexity: 4, 20 LOC (OK)
- ✅ src/models.py:User.validate_email() - Complexity: 2, 8 LOC (OK)

**File Sizes:**
- ✅ src/service.py - 180 LOC (OK)
- ✅ src/models.py - 90 LOC (OK)
- ✅ src/api.py - 150 LOC (OK)

**Nesting Depth:**
- ✅ Max depth: 3 (within threshold)

**Verdict:** ✅ PASS (no violations)
```

### FAIL Example

```markdown
## Complexity Report

**Violations:** 4

**Complex Functions (>20 complexity):**
1. ❌ `src/service.py:PaymentProcessor.process_transaction()` (lines 50-150, 101 LOC)
   **Cyclomatic Complexity:** 25
   **Nesting Depth:** 6
   **Recommendation:** Extract sub-functions:
   - `validate_transaction()`
   - `calculate_fees()`
   - `execute_payment()`
   - `handle_result()`

2. ❌ `src/api.py:DataController.export_data()` (lines 200-250, 51 LOC)
   **Cyclomatic Complexity:** 22
   **Nesting Depth:** 5
   **Recommendation:** Split into multiple controllers

**Files >200 LOC:**
3. ❌ `src/service.py:LegacyPaymentService` (lines 1-350, 350 LOC)
   **Recommendation:** Split into multiple classes

**Nesting Depth >4:**
4. ❌ `src/models.py:User.from_dict()` (lines 100-150)
   **Nesting Depth:** 5
   **Recommendation:** Flatten nesting using helper functions

**Verdict:** ❌ FAIL (violations detected)

**Required Actions:**
1. Refactor PaymentProcessor.process_transaction() (split into 4 functions)
2. Refactor DataController.export_data() (extract methods)
3. Split LegacyPaymentService into 3 classes
4. Flatten User.from_dict() nesting

**Complex Functions Requiring Refactoring:**
- src/service.py:PaymentProcessor.process_transaction()
- src/api.py:DataController.export_data()
- src/models.py:User.from_dict()
```

## Refactoring Patterns

### Extract Method

**Before (too complex):**
```python
def process_transaction(transaction):
    if transaction.type == "payment":
        if transaction.amount > 100:
            if transaction.currency == "USD":
                if transaction.valid:
                    # 50 lines of logic
                    return "processed"
                else:
                    return "invalid"
            else:
                return "unsupported"
        else:
            return "small"
    else:
        return "unknown"
```

**After (simplified):**
```python
def process_transaction(transaction):
    if transaction.type == "payment":
        return process_payment(transaction)
    else:
        return "unknown"

def process_payment(transaction):
    if transaction.amount <= 100:
        return "small"
    if transaction.currency != "USD":
        return "unsupported"
    if not transaction.valid:
        return "invalid"
    # Payment logic extracted
    return process_large_payment(transaction)

def process_large_payment(transaction):
    # Extracted 50 lines of logic here
    return "processed"
```

### Extract Class

**Before (single file >200 LOC):**
```python
# service.py: 350 lines
class LegacyPaymentService:
    def method1(self): ...  # 50 lines
    def method2(self): ...  # 50 lines
    # ... 5 more methods
    def method7(self): ...  # 50 lines
```

**After (split into 3 classes):**
```python
# service/payment_processor.py: 100 lines
class PaymentProcessor:
    def method1(self): ...
    def method2(self): ...
    def method3(self): ...

# service/payment_validator.py: 100 lines
class PaymentValidator:
    def validate(self): ...
    def check_rules(self): ...

# service/payment_executor.py: 100 lines
class PaymentExecutor:
    def execute(self): ...
    def finalize(self): ...
```

### Flatten Nesting

**Before (deep nesting):**
```python
def create_user(data):
    if "name" in data:
        if "email" in data:
            if "email" in data["email"]:
                if validate_email(data["email"]):
                    if is_unique(data["email"]):
                        return create(data["name"], data["email"])
```

**After (flattened):**
```python
def create_user(data):
    if "name" not in data:
        raise ValueError("Missing name")
    if "email" not in data:
        raise ValueError("Missing email")

    if not validate_email(data["email"]):
        raise ValueError("Invalid email")

    if not is_unique(data["email"]):
        raise ValueError("Email already exists")

    return create(data["name"], data["email"])
```

## Language-Specific Patterns

### Python

Count complexity as:
- Number of `if`, `elif`, `for`, `while`, `except`, `and`, `or` plus 1

### Java

Count complexity as:
- Number of `if`, `else if`, `for`, `while`, `case`, `catch` plus 1

### Go

Count complexity as:
- Number of `if`, `else`, `for`, `switch`, `case` plus 1

## Quick Reference

| Metric | Python | Java | Go | Threshold |
|--------|--------|------|-----|----------|
| **Cyclomatic Complexity** | +1 per branch | +1 per branch | +1 per branch | <10 OK, >20 FAIL |
| **Function Length** | LOC | LOC | LOC | <200 OK |
| **Nesting Depth** | Indent level | Brace level | Indent level | ≤4 OK |

## Tools (Optional Validation)

If language tools available, use for additional checks:

**Python:**
```bash
radon cc src/ -a -s  # Cyclomatic complexity
radon cc src/ -a -s | grep "C"  # List complex functions
```

**Java:**
```bash
# Check complexity manually or use IDE warnings
# Most IDEs show complexity metrics
```

**Go:**
```bash
gocyclo -over 10 .  # List complex functions
gocyclo -over 15 .  # List very complex functions
```

## Quality Gate

- **PASS:** No functions exceed thresholds
- **FAIL:** Any violation found

## See Also

- [@build Skill](../skills/build.md) - Calls complexity validator
- [@review Skill](../skills/review.md) - Runs complexity validator
- [Quality Gates Reference](../../docs/quality-gates.md) - Complexity criteria
