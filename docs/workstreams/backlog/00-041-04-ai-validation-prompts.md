# 00-041-04: AI-Based Validation Prompts

> **Feature:** F041 - Claude Plugin Distribution
> **Status:** backlog
> **Size:** MEDIUM
> **Created:** 2026-02-02

## Goal

Create AI-based validation prompts to replace static analysis tools (pytest, mypy, ruff) for language-agnostic quality gates.

## Acceptance Criteria

- AC1: Coverage validator prompt created and tested
- AC2: Architecture validator prompt created and tested
- AC3: Error handling validator prompt created and tested
- AC4: Complexity validator prompt created and tested
- AC5: Accuracy ≥90% compared to tool-based validation

## Scope

### Input Files
- Existing tool-based validation logic (pytest coverage, grep patterns)
- SDP quality gate requirements (80% coverage, <200 LOC, etc.)

### Output Files
- `sdp-plugin/prompts/validators/coverage.md` (NEW)
- `sdp-plugin/prompts/validators/architecture.md` (NEW)
- `sdp-plugin/prompts/validators/errors.md` (NEW)
- `sdp-plugin/prompts/validators/complexity.md` (NEW)
- `sdp-plugin/prompts/validators/all.md` (orchestrator, NEW)

### Out of Scope
- Modifying skills (WS-00-041-03 integrates these validators)
- Go binary (WS-00-041-05)

## Implementation Steps

### Step 1: Coverage Validator

**File: sdp-plugin/prompts/validators/coverage.md**

```markdown
# Coverage Validator

You are a test coverage analyst. Analyze test coverage by:

## 1. Read Files

Read all test files and source files:
- Test files: `tests/`, `src/test/`, `*_test.go`
- Source files: `src/`, `src/main/`, `*.go`

## 2. Map Tests to Code

For each function/method in source files:
- Find corresponding test in test files
- Identify untested functions
- Identify untested branches (if/else, try/except, switch/case)

## 3. Calculate Coverage

```
coverage = (tested_functions / total_functions) * 100
```

## 4. Output Format

```markdown
## Coverage Report

**Total Coverage:** 85%

**Untested Functions:**
- `src/service.py:UserService.delete_user()` (lines 45-50)
- `src/models.py:User.validate_email()` (lines 78-82)

**Missing Branch Coverage:**
- `src/auth.py:login()` - error path not tested (line 45)

**Test Files Found:** tests/test_user.py, tests/test_auth.py
**Source Files Found:** src/service.py, src/models.py, src/auth.py

**Verdict:** ✅ PASS (≥80%) / ❌ FAIL (<80%)
```

## 5. Gate

If coverage <80%, FAIL with specific missing lines and functions.
```

### Step 2: Architecture Validator

**File: sdp-plugin/prompts/validators/architecture.md**

```markdown
# Architecture Validator

You are a Clean Architecture enforcer. Validate layer separation.

## 1. Read All Source Files

Parse all files in `src/` directory.

## 2. Build Dependency Graph

For each file:
- Parse import statements
- Map file path to layer:
  - `src/domain/` → Domain layer
  - `src/application/` → Application layer
  - `src/infrastructure/` → Infrastructure layer
  - `src/presentation/` → Presentation layer
- Track dependencies between files

## 3. Check Violations

**Forbidden Dependencies:**
- ❌ Domain imports Infrastructure
- ❌ Domain imports Application
- ❌ Domain imports Presentation
- ❌ Application imports Presentation

**Allowed Dependencies:**
- ✅ Infrastructure imports Domain
- ✅ Application imports Domain
- ✅ Presentation imports Application

## 4. Output Format

```markdown
## Architecture Report

**Violations Found:** 2

1. ❌ `src/domain/entities/user.py` (line 5)
   ```
   from src.infrastructure.persistence import Database
   ```
   Domain must not import infrastructure.

2. ❌ `src/application/auth.py` (line 12)
   ```
   from src.presentation.controllers import UserController
   ```
   Application must not import presentation.

**Allowed Dependencies (OK):**
- ✅ src/infrastructure/repositories.py → src/domain/entities.py
- ✅ src/application/services.py → src/domain/entities.py

**Verdict:** ✅ PASS (no violations) / ❌ FAIL (violations detected)
```

## 5. Gate

Zero violations required for PASS.
```

### Step 3: Error Handling Validator

**File: sdp-plugin/prompts/validators/errors.md**

```markdown
# Error Handling Validator

You are an error handling auditor. Find unsafe exception handling.

## 1. Search for Patterns

Search for unsafe error handling patterns:

**Python:**
- `except:` (bare except)
- `except Exception:` (too broad, unless logged+re-raised)
- `except: pass` (swallows errors)

**Java:**
- `catch(Exception e)` (too broad)
- `catch(Throwable e)` (catches everything)
- Empty catch blocks

**Go:**
- `recover()` without error check
- Ignored error returns (`func(), _`)

## 2. For Each Match

Check if error is:
- Logged with context
- Re-raised if critical
- Specific exception type (not catch-all)

## 3. Output Format

```markdown
## Error Handling Report

**Violations:** 3

1. ❌ `src/service.py` (lines 45-47)
   ```python
   try:
       risky_operation()
   except:
       pass  # Swallows all errors!
   ```
   Fix: Catch specific exception, log error, re-raise.

2. ❌ `src/auth.py` (line 78)
   ```python
   except Exception as e:
       return  # Error not logged
   ```
   Fix: Log error before returning.

3. ⚠️ `src/api.py` (line 123)
   ```go
   data, _ := fetchData()  # Error ignored
   ```
   Fix: Check error and handle.

**Safe Error Handling (OK):**
- ✅ src/user.py:52 (catches ValueError, logs, re-raises)

**Verdict:** ✅ PASS (no violations) / ❌ FAIL (violations detected)
```

## 4. Gate

Zero unsafe error handling violations for PASS.
```

### Step 4: Complexity Validator

**File: sdp-plugin/prompts/validators/complexity.md**

```markdown
# Complexity Validator

You are a complexity analyst. Identify overly complex code.

## 1. For Each Function

Calculate metrics:
- **Cyclomatic Complexity:** Number of branches + loops + conditions
- **Lines of Code:** Total lines in function
- **Nesting Depth:** Maximum indentation level

## 2. Thresholds

- Complexity <10: ✅ OK
- Complexity 10-20: ⚠️ Warning (consider refactoring)
- Complexity >20: ❌ FAIL (too complex)
- Lines >200: ❌ FAIL
- Nesting >4: ❌ FAIL

## 3. Output Format

```markdown
## Complexity Report

**Violations:** 2

1. ❌ `src/service.py:process_data()` (lines 50-150, 101 LOC)
   **Cyclomatic Complexity:** 15 (too many branches)
   **Nesting Depth:** 5 (too deep)

   Recommendation: Extract sub-functions:
   - `validate_input()`
   - `transform_data()`
   - `save_result()`

2. ❌ `src/auth.py:authenticate()` (lines 200-250, 51 LOC)
   **Cyclomatic Complexity:** 22 (exceeds threshold)

   Recommendation: Break into smaller functions with single responsibility.

**Complex Functions (Warning):**
- ⚠️ `src/models.py:User.serialize()` - Complexity 18

**Simple Functions (OK):**
- ✅ src/utils.py:format_date() - Complexity 2, 5 LOC
- ✅ src/api.py:handle_request() - Complexity 5, 12 LOC

**Verdict:** ✅ PASS (no violations) / ❌ FAIL (violations detected)
```

## 4. Gate

No functions exceed thresholds for PASS.
```

### Step 5: Orchestrator

**File: sdp-plugin/prompts/validators/all.md**

```markdown
# All Validators

Run all validators in sequence and aggregate results.

## Execution Order

1. `/coverage-validator`
2. `/architecture-validator`
3. `/error-validator`
4. `/complexity-validator`

## Aggregate Verdict

- **PASS:** All 4 validators pass
- **FAIL:** Any validator fails

## Output Format

```markdown
## Quality Gates Summary

✅ Coverage: 85% (≥80% required)
✅ Architecture: No violations
❌ Error Handling: 3 bare except clauses
✅ Complexity: All functions <200 LOC

**Overall Verdict:** ❌ FAIL

**Required Actions:**
1. Fix bare except clauses in src/service.py:45, src/auth.py:78
2. Re-run validators after fixes
```
```

## Verification (90%+ Accuracy Target)

```bash
# Benchmark on real codebases

for project in sdp hw_checker spring-petclinic gin; do
  echo "Testing $project..."

  # Tool-based (ground truth)
  cd tests/$project
  if [ -f "pyproject.toml" ]; then
    pytest --cov=src/ > tool-coverage.txt 2>&1
  elif [ -f "pom.xml" ]; then
    mvn verify > tool-coverage.txt 2>&1
  elif [ -f "go.mod" ]; then
    go test -cover > tool-coverage.txt 2>&1
  fi

  # AI-based
  claude "/coverage-validator" > ai-coverage.txt 2>&1

  # Compare
  echo "Tool: $(grep 'coverage' tool-coverage.txt)"
  echo "AI:   $(grep 'coverage' ai-coverage.txt)"
done

# Expected: AI matches tools in ≥90% of cases
```

## Quality Gates

- All 4 validators created
- Output format consistent across validators
- Benchmark shows ≥90% accuracy
- False positive rate ≤5%
- False negative rate ≤5%

## Dependencies

- 00-041-03 (Remove Python Dependencies from Skills) - skills integrate these validators

## Blocks

- 00-041-06 (Cross-Language Validation) - needs validators for testing
