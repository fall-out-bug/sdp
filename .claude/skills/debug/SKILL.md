---
name: debug
description: Systematic 4-phase debugging process using scientific method for evidence-based root cause analysis
tools: Read, Write, Edit, Bash, Glob, Grep
---

# /debug - Systematic Debugging

Systematic 4-phase root cause analysis using scientific method.

## When to Use

- You have a specific bug to fix
- You need evidence-based debugging (not trial-and-error)
- You want to follow systematic process
- You need to prevent infinite fix loops

## Invocation

```bash
/debug "Description of the issue"
# Example: /debug "API returns 500 on /submissions endpoint"
```

## Workflow

**IMPORTANT:** This skill delegates to the master prompt.

### Load Master Prompt

First, read the systematic debugging protocol:

```bash
cat prompts/skills/systematic-debugging.md
```

**This file contains:**
- 4-phase debugging process (Evidence → Pattern → Hypothesis → Implementation)
- Evidence collection checklist
- Pattern analysis techniques
- Hypothesis testing protocol
- Root-cause tracing method
- Failsafe rule (3 strikes → stop, question architecture)

### Execute 4 Phases

Follow the systematic debugging process from the master prompt:

#### Phase 1: Evidence Collection
- Collect error messages
- Document reproduction steps
- Check recent changes
- Capture environment state

#### Phase 2: Pattern Analysis
- Find working examples
- Compare working vs. broken
- Identify patterns

#### Phase 3: Hypothesis Testing
- Form ONE hypothesis
- Design minimal test
- Execute test
- Record result (PASS/FAIL)

#### Phase 4: Implementation
- Write failing test first
- Implement minimal fix
- Verify fix (unit + regression + integration)
- Document root cause

### Failsafe Rule

**After 3 failed fix attempts → STOP, escalate to architecture review**

Do NOT continue debugging. Create architecture WS instead.

## Output Format

```markdown
# Debug Session: [Issue Description]

## Phase 1: Evidence Collection

**Error Messages:**
```
[Error logs]
```

**Reproduction Steps:**
1. [Step 1]
2. [Step 2]

**Recent Changes:**
- [File 1]: [Change]

**Environment:**
- Python: [version]
- OS: [version]

## Phase 2: Pattern Analysis

**Working Examples:**
- [Example 1]

**Comparison:**
| Aspect | Working | Broken | Difference |
|--------|---------|--------|------------|
| [Aspect] | [value] | [value] | [diff] |

## Phase 3: Hypothesis Testing

**Hypothesis #1:** [Clear statement]

**Test:**
```python
[Minimal test code]
```

**Result:** PASS / FAIL

## Phase 4: Implementation

**Failing Test:**
```python
def test_fix():
    # Reproduce bug
    assert broken_function() == expected  # Fails initially
```

**Fix:**
```python
def broken_function():
    # Minimal fix
    pass
```

**Verification:**
- Unit test: ✅ PASS
- Regression: ✅ PASS
- Integration: ✅ PASS

**Root Cause:** [Clear explanation]
```

## Related Commands

- `/issue` - For full issue analysis (severity, routing, GitHub)
- `/hotfix` - For P0 production fixes
- `/bugfix` - For P1/P2 feature bugs
