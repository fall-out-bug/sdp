---
name: verify-workstream
description: Validate workstream documentation against codebase reality.
---

# Verify Workstream

Before executing any workstream, validate that the documentation matches the actual codebase reality.

## When to Use

Use this skill **before** starting any workstream execution (`@build` or `@oneshot`).

## Quick Reference

| Step | Action | Gate |
|------|--------|------|
| 1 | Read workstream description | Frontmatter parsed |
| 2 | Find scope files | All files located |
| 3 | Read actual implementation | Code understood |
| 4 | Compare docs vs reality | Discrepancies listed |
| 5 | Recommend action | Clear next steps |

## Workflow

### Step 1: Read Workstream Description

Parse the workstream frontmatter to understand:

```yaml
# From WS frontmatter
scope_files:
  - src/sdp/quality/validators.py
  - src/sdp/quality/models.py

goal: Create generic validation layer
acceptance_criteria:
  - AC1: Generic validators implemented
```

Extract:
- **Goal:** What should this workstream achieve?
- **Scope Files:** Which files will be modified?
- **Acceptance Criteria:** What defines success?
- **Contracts:** What interfaces/modules are documented?

### Step 2: Find Scope Files

Use Glob to locate all files in scope:

```bash
# Find all Python files in scope
Glob("**/*.py", path="src/sdp/quality")

# Find specific files mentioned in docs
Glob("src/sdp/quality/validators.py")
Glob("src/sdp/quality/models.py")
```

**Gate:** All files must exist. If any file is missing, alert the user.

### Step 3: Read Actual Implementation

For each file in scope, read and analyze:

```python
# Read file structure
import ast
import inspect

# For each file:
1. Parse module structure (classes, functions)
2. Identify actual implementation patterns
3. Check for business logic vs generic logic
4. Note dependencies and imports
```

**What to Look For:**
- **File Structure:** Does the file contain what's documented?
- **Function/Class Names:** Are documented functions actually present?
- **Logic Type:** Is it generic validation or business logic?
- **Dependencies:** What does this file actually depend on?

### Step 4: Compare Docs vs Reality

Create a comparison table:

```markdown
## Documentation vs Reality Analysis

### File: src/sdp/quality/validators.py

| Aspect | Documentation | Reality | Status |
|--------|---------------|---------|--------|
| **Purpose** | Generic validation layer | Contains business logic (UserValidator, PaymentValidator) | ❌ Mismatch |
| **Functions** | validate_contract() | validate_user(), validate_payment() | ❌ Missing |
| **Logic Type** | Generic/reusable | Domain-specific | ❌ Mismatch |

### File: src/sdp/quality/models.py

| Aspect | Documentation | Reality | Status |
|--------|---------------|---------|--------|
| **Purpose** | Validation models | Dataclasses (User, Payment) | ❌ Mismatch |
| **Content** | Validation schemas | Domain entities | ❌ Wrong layer |

## Summary

- **Total Files Checked:** 2
- **Mismatches Found:** 3
- **Severity:** HIGH (workstream description does not match reality)
```

### Step 5: Recommend Action

Based on severity, recommend next steps:

**If MISMATCH Found (High Severity):**
```markdown
## ⚠️ Recommendation: PAUSE and Clarify

The workstream description assumes a generic validation layer, but the actual implementation contains business logic.

**Options:**
1. **Update Workstream** to reflect business logic location
2. **Create new file** for generic validators (if migration is intended)
3. **Split workstream** into: (a) extract generic validation, (b) implement business validators

**User Choice Required:** Which approach should we take?
```

**If MINOR Mismatch (Low Severity):**
```markdown
## ✅ Recommendation: Proceed with Adaptations

Minor discrepancies found:
- Missing function: validate_contract() (can be added)
- Additional imports: yaml, toml (not documented)

**Action:** Document deviations in workstream frontmatter and proceed.
```

**If NO Mismatch:**
```markdown
## ✅ Recommendation: Proceed

All scope files match documentation. Ready to execute workstream.
```

## Output Format

Always end with a clear recommendation:

```markdown
## Verification Complete

**Workstream:** 00-040-04
**Files Checked:** 2
**Mismatches:** 3

**Severity:** HIGH
**Recommendation:** [PAUSE / PROCEED / PROCEED WITH ADAPTATIONS]

**Next Steps:** [Clear action items]

---

[Comparison table]
[Discrepancy details]
```

## Example Usage

```markdown
User: /verify-workstream 00-040-04

Claude:
## Verifying Workstream 00-040-04...

### Step 1: Reading Workstream Description
- Goal: Remove redundant validation code
- Scope: src/sdp/quality/validators.py, models.py
- AC: Extract generic validation to separate module

### Step 2: Locating Scope Files
✓ Found: src/sdp/quality/validators.py (245 LOC)
✓ Found: src/sdp/quality/models.py (180 LOC)

### Step 3: Reading Implementation
[Analyzing code structure...]

### Step 4: Comparing Docs vs Reality
[Comparison table...]

### Step 5: Recommendation
⚠️ PAUSE - Workstream assumes generic validators but file contains business logic
```

## Integration with @build

This skill is automatically invoked before `@build` execution:

```markdown
User: @build 00-040-04

Claude:
→ /verify-workstream 00-040-04
✓ Verification complete
→ Proceeding with @build execution...
```

## Success Metrics

- **Detection Rate:** ≥90% of documentation mismatches caught before execution
- **Time Saved:** Prevents wasted work on misaligned implementations
- **User Satisfaction:** Reduction in "wrong_approach" friction events

## Error Handling

**If File Not Found:**
```markdown
❌ ERROR: Scope file not found
Expected: src/sdp/quality/validation.py
Actual: File does not exist

**Action:** Update workstream scope to reflect actual codebase structure
```

**If Parse Error:**
```markdown
❌ ERROR: Unable to parse workstream frontmatter
**Action:** Ensure workstream file has valid YAML frontmatter
```

## Related Skills

- `/reality-check` - Quick validation for single files
- `/build` - Automatically runs verification before execution
- `/review` - Checks for drift during quality reviews
