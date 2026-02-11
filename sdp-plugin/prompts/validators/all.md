# All Validators Orchestrator

Run all validators in sequence and aggregate results into a unified quality gates report.

## Purpose

Execute complete quality validation by running all 4 validators and producing a consolidated verdict.

## How to Use

```
Ask Claude: "Run all quality validators by:
1. Running coverage validator (analyze test coverage)
2. Running architecture validator (check layer violations)
3. Running error handling validator (find unsafe patterns)
4. Running complexity validator (check file/function size)

Aggregate Results:
- Collect verdicts from all 4 validators
- Calculate overall pass rate
- Produce unified report with action items

Output:
- Summary table with all validator results
- Overall verdict: ✅ PASS (all pass) or ❌ FAIL (any fail)
- Consolidated list of required actions
- Priority ordering (CRITICAL > HIGH > MEDIUM > LOW)"
```

## Execution Order

Validators run in this sequence:

1. **Coverage Validator** - Test coverage analysis
2. **Architecture Validator** - Clean Architecture compliance
3. **Error Handling Validator** - Error pattern safety
4. **Complexity Validator** - Code complexity metrics

## Output Format

### PASS Example (All Validators Pass)

```markdown
## Quality Gates Summary

**Execution Date:** 2025-01-15 14:30:00
**Project:** sdp-plugin
**Validators Run:** 4/4

### Validator Results

| Validator | Verdict | Details |
|-----------|---------|---------|
| **Coverage** | ✅ PASS | 87% coverage (≥80% required) |
| **Architecture** | ✅ PASS | No layer violations |
| **Error Handling** | ✅ PASS | No unsafe patterns found |
| **Complexity** | ✅ PASS | All files <200 LOC, max complexity 8 |

### Overall Verdict: ✅ PASS

**Pass Rate:** 4/4 (100%)

**Summary:**
- All quality gates passed
- Test coverage exceeds threshold
- Clean architecture maintained
- Error handling is safe
- Code complexity is acceptable

**Next Steps:**
- Ready for review
- Proceed with deployment workflow
```

### FAIL Example (Multiple Violations)

```markdown
## Quality Gates Summary

**Execution Date:** 2025-01-15 14:35:00
**Project:** sdp-plugin
**Validators Run:** 4/4

### Validator Results

| Validator | Verdict | Details |
|-----------|---------|---------|
| **Coverage** | ❌ FAIL | 72% coverage (≥80% required) |
| **Architecture** | ✅ PASS | No layer violations |
| **Error Handling** | ❌ FAIL | 3 bare except clauses found |
| **Complexity** | ❌ FAIL | 2 files >200 LOC |

### Overall Verdict: ❌ FAIL

**Pass Rate:** 1/4 (25%)

### Critical Issues (Fix Immediately)

1. **[CRITICAL]** `src/auth.py:login()` (line 45-50)
   ```python
   try:
       return authenticate(username, password)
   except:
       return None  # Security risk: hides auth errors
   ```
   **Fix:** Log security events, re-raise auth failures
   **Validator:** Error Handling

### High Priority Issues

2. **[HIGH]** Coverage insufficient: 72% (target ≥80%)
   **Untested Functions:**
   - `src/service.py:UserService.delete_user()` (lines 45-50)
   - `src/api.py:DataController.export_data()` (lines 120-135)
   **Fix:** Add tests for untested functions
   **Validator:** Coverage

3. **[HIGH]** `src/service.py:PaymentProcessor` (lines 1-250, 250 LOC)
   **Violation:** File exceeds 200 LOC
   **Recommendation:** Split into multiple classes:
     - `PaymentValidator` - Validation logic
     - `PaymentExecutor` - Execution logic
     - `PaymentResultHandler` - Result processing
   **Validator:** Complexity

4. **[HIGH]** `src/models.py:UserFactory` (lines 100-210, 110 LOC)
   **Violation:** File exceeds 200 LOC
   **Recommendation:** Extract factory methods to separate module
   **Validator:** Complexity

### Medium Priority Issues

5. **[MEDIUM]** `src/service.py:process_data()` (line 78)
   ```python
   try:
       risky_operation()
   except Exception:
       pass  # Error not logged
   ```
   **Fix:** Log error before swallowing
   **Validator:** Error Handling

### Required Actions (Priority Order)

1. **[CRITICAL]** Fix `src/auth.py:login()` bare except (security risk)
2. **[HIGH]** Add tests for `UserService.delete_user()` and `DataController.export_data()`
3. **[HIGH]** Split `PaymentProcessor` into 3 classes (250 LOC → 3 × ~80 LOC)
4. **[HIGH]** Split `UserFactory` into 2 modules (110 LOC → 2 × ~55 LOC)
5. **[MEDIUM]** Fix `src/service.py:process_data()` error handling

### Verification Steps

After fixing all issues:
1. Re-run coverage validator → Verify ≥80%
2. Re-run error validator → Verify 0 violations
3. Re-run complexity validator → Verify all files <200 LOC
4. Re-run this orchestrator → Confirm overall PASS

**Estimated Fix Time:** 2-3 hours

**Estimated Test Time:** 15 minutes
```

## Severity Classification

| Severity | Description | Example | Action Timeline |
|----------|-------------|---------|-----------------|
| **CRITICAL** | Security or data loss risk | Auth errors hidden, silent crashes | Fix immediately |
| **HIGH** | Quality gate failure | Low coverage, large files, bare except | Fix before commit |
| **MEDIUM** | Code quality issue | Error not logged, generic exception | Fix in next iteration |
| **LOW** | Style or optimization | Long lines, deep nesting | Fix when convenient |

## Quality Gate Rules

### Overall Verdict Logic

```python
if all_validators_pass():
    return "✅ PASS"
else:
    if any_critical_issues():
        return "❌ FAIL (CRITICAL issues)"
    elif any_high_priority_issues():
        return "❌ FAIL (HIGH priority issues)"
    else:
        return "⚠️ WARNING (MEDIUM/LOW issues)"
```

### Individual Validator Thresholds

| Validator | PASS Condition | FAIL Condition |
|-----------|----------------|----------------|
| **Coverage** | ≥80% | <80% |
| **Architecture** | 0 violations | ≥1 violation |
| **Error Handling** | 0 HIGH/CRITICAL | ≥1 HIGH/CRITICAL |
| **Complexity** | All files <200 LOC, CC <10 | Any file >200 LOC or CC >20 |

## Integration Points

### Used By

- **@review skill** - Runs all validators as quality check
- **@build skill** - Runs all validators in post-build validation
- **Pre-commit hooks** - Optional gate before git commit

### Validator Dependencies

Each validator is independent and can run standalone:
- `/coverage-validator` - Coverage only
- `/architecture-validator` - Architecture only
- `/error-validator` - Error handling only
- `/complexity-validator` - Complexity only

This orchestrator runs all 4 in sequence.

## Common Workflows

### Workflow 1: Pre-Commit Check

```bash
# Before committing workstream
git add .
claude "Run all quality validators on staged files"
# Expected: Overall PASS before commit
git commit -m "feat: implement workstream"
```

### Workflow 2: Feature Review

```bash
# After completing feature
claude "@review F01"
# Internally runs all validators
# Expected: Overall PASS before deploy
```

### Workflow 3: Continuous Quality

```bash
# Run validators on each workstream completion
claude "@build 00-001-01"
# Post-build step runs all validators automatically
# Expected: PASS before marking WS complete
```

## Troubleshooting

### Issue: Validator Fails to Run

**Symptom:** "Cannot read validator file"

**Fix:** Ensure all 4 validator files exist in `sdp-plugin/prompts/validators/`:
- coverage.md
- architecture.md
- errors.md
- complexity.md

### Issue: False Positive

**Symptom:** Validator reports violation but code is correct

**Fix:**
1. Review validator output for specific line numbers
2. Check if pattern matches correctly
3. Update validator prompt if pattern is too strict
4. Re-run validator after fix

### Issue: Slow Execution

**Symptom:** Validators take >5 minutes on large codebase

**Fix:**
1. Run validators on specific directories only: "Run validators on src/service/"
2. Run validators independently: "/coverage-validator" (faster than all)
3. Increase timeout: "Run all validators with 10 minute timeout"

## Performance Benchmarks

| Project Size | Lines of Code | Execution Time |
|--------------|---------------|----------------|
| Small | <1,000 LOC | ~30 seconds |
| Medium | 1,000-10,000 LOC | ~2 minutes |
| Large | >10,000 LOC | ~5 minutes |

**Bottleneck:** Coverage validator (reads all source + test files)

**Optimization:** Cache file list between validators

## Quality Gate

- **PASS:** All 4 validators pass (100% pass rate)
- **FAIL:** Any validator fails with HIGH or CRITICAL issues
- **WARNING:** Only MEDIUM/LOW issues (allowed to proceed with caution)

## See Also

- [Coverage Validator](./coverage.md) - Test coverage analysis
- [Architecture Validator](./architecture.md) - Clean Architecture enforcement
- [Error Handling Validator](./errors.md) - Error pattern safety
- [Complexity Validator](./complexity.md) - Code complexity metrics
- [@review Skill](../skills/review.md) - Runs all validators automatically
- [Quality Gates Reference](../../docs/quality-gates.md) - Complete quality criteria
