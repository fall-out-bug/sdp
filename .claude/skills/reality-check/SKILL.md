# Reality Check

Quick validation that documentation matches actual code before making changes.

## When to Use

Use this skill when:
- About to modify a file based on documentation
- Unsure if docs reflect current implementation
- Starting work on an unfamiliar part of codebase
- **Quick check** (faster than full `/verify-workstream`)

## Quick Reference

| Step | Action | Time |
|------|--------|------|
| 1 | Read actual code | 30s |
| 2 | Compare with assumption | 30s |
| 3 | Report mismatch | 30s |

**Total Time:** ~90 seconds (vs. 5-10 minutes for full verify-workstream)

## Workflow

### Step 1: Read Actual Code First

**Before** reading documentation, read the actual implementation:

```bash
# Example: User says "add validation to models.py"
# DON'T read docs first!

# READ THE FILE FIRST:
Read("src/sdp/quality/models.py")

# Analyze:
- What does this file ACTUALLY contain?
- What patterns are used?
- What is the actual architecture?
```

**Key Question:** "What does the code actually do, not what it should do?"

### Step 2: Compare With Assumptions

After reading code, compare with documentation/assumptions:

```markdown
## Reality Check: src/sdp/quality/models.py

### Actual Code Structure:
```
# File contains:
- dataclass User (id, name, email)
- dataclass Payment (id, amount, status)
- No validation logic present
- Pure domain entities
```

### Documentation/Assumption:
- "Validation models with type checking"
- "Contains validate() methods"
- "Enforces business rules"

### Comparison:
| Aspect | Assumption | Reality | Match? |
|--------|------------|---------|--------|
| **Content** | Validation logic | Dataclasses | ❌ No |
| **Purpose** | Type enforcement | Domain modeling | ❌ No |
| **Methods** | validate() | None | ❌ No |
```

### Step 3: Report and Recommend

Based on mismatch severity:

**If MISMATCH Found:**
```markdown
## ⚠️ Reality Check Failed

**Assumption:** File contains validation logic
**Reality:** File contains domain dataclasses

**Impact:** HIGH - Adding validation here would violate architecture

**Recommendation:**
1. Create separate `validators.py` for validation logic
2. Keep `models.py` as pure domain entities
3. Don't mix concerns

**User Choice:** How should we proceed?
```

**If MATCH Found:**
```markdown
## ✅ Reality Check Passed

Documentation matches actual implementation.

**File:** src/sdp/quality/models.py
**Content:** Domain dataclasses (as documented)
**Status:** Safe to proceed with modifications
```

## Example Usage

```markdown
User: Add email validation to User model in models.py

Claude:
Let me reality-check first...

[Reads src/sdp/quality/models.py]

## Reality Check Results

**File Contains:**
- dataclass User (id: int, name: str, email: str)
- No validation logic present
- Pure domain entity pattern

**Assumption Check:**
❌ User model is a domain entity, NOT a validation model

**Recommendation:**
Don't add validation to models.py (violates clean architecture).

**Better Approach:**
Create UserValidator in validators.py:
```python
class UserValidator:
    def validate_email(self, email: str) -> bool:
        return "@" in email and "." in email
```

Shall I proceed with creating the validator?
```

## Output Format

Always include:

```markdown
## Reality Check: <filename>

### What Code Actually Does:
[Code structure summary]

### What We Assumed:
[Documentation/expectation]

### Comparison:
[Quick table]

### Recommendation:
[✅ Proceed / ⚠️ Stop / 🔄 Adapt]
```

## Integration with Other Skills

This is a **lightweight version** of `/verify-workstream`:

| Aspect | /reality-check | /verify-workstream |
|--------|----------------|-------------------|
| **Scope** | Single file | Full workstream |
| **Time** | 90 seconds | 5-10 minutes |
| **Detail** | Quick check | Comprehensive analysis |
| **Use Case** | Quick validation | Pre-build verification |

**Usage Pattern:**
```
/reality-check → Quick validation during conversation
/verify-workstream → Before @build execution
```

## Common Patterns

### Pattern 1: Documentation Claims X, Code Has Y

```markdown
## Reality Check Failed

**Docs Say:** "Generic validation functions"
**Code Has:** Business logic (UserValidator, PaymentValidator)

**Pattern:** Documentation drift → code evolved, docs didn't

**Fix:** Update docs OR extract generic logic (user choice)
```

### Pattern 2: Assumption Based on Filename

```markdown
## Reality Check Failed

**Assumption:** "models.py contains data models"
**Reality:** "models.py contains validation logic"

**Pattern:** Filename doesn't match content (architectural drift)

**Fix:** Rename file OR restructure code
```

### Pattern 3: Missing Implementation

```markdown
## Reality Check Failed

**Docs Say:** "File contains validate_contract()"
**Reality:** Function not found

**Pattern:** Documentation ahead of implementation

**Fix:** Implement function OR remove from docs
```

## Anti-Patterns to Avoid

❌ **Don't read docs first** - Start with code
❌ **Don't assume docs are correct** - They may be outdated
❌ **Don't skip this step** - 90 seconds saves hours of rework

## Success Metrics

- **Mismatch Detection:** Catches documentation drift before implementation
- **Time Saved:** Prevents "wrong_approach" friction
- **Architecture Preservation:** Maintains clean architecture boundaries

## Related Skills

- `/verify-workstream` - Full workstream validation
- `/build` - Uses reality-check during execution
- `/review` - Checks for drift in completed work

## Quick Reference Card

```
┌─────────────────────────────────────┐
│  REALITY CHECK: 3 Steps, 90s        │
├─────────────────────────────────────┤
│  1. Read actual code (FIRST!)       │
│  2. Compare with docs/assumptions   │
│  3. Report: ✅ Proceed or ⚠️ Stop   │
└─────────────────────────────────────┘
```

**Remember:** Code never lies. Documentation always lags.
