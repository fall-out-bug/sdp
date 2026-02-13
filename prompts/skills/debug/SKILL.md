---
name: debug
description: Systematic debugging using scientific method for evidence-based root cause analysis
version: 2.0.0
changes:
  - Converted to LLM-agnostic format
  - Removed tool-specific API references
  - Focus on WHAT, not HOW to invoke
---

# @debug - Systematic Debugging

Evidence-based debugging using scientific method. Not "try stuff and see" -- systematic investigation.

---

## EXECUTE THIS NOW

When user invokes `@debug "<issue>"`:

### Phase 1: OBSERVE - Gather Facts

**Goal:** Collect evidence WITHOUT forming hypotheses

**Actions:**
1. Read error messages/logs completely
2. Check git diff for recent changes
3. Verify environment (Python version, dependencies)
4. Check configuration files
5. Reproduce the bug consistently

**Output:** Observation log with timestamps, error messages, environment state

### Phase 2: HYPOTHESIZE - Form Theories

**Goal:** Create testable theories about root cause

**Process:**
1. List ALL possible causes (brainstorm)
2. Rank by likelihood (use evidence)
3. Select TOP theory to test first
4. Define falsification test

**Output:** Hypothesis list with ranked theories

### Phase 3: EXPERIMENT - Test Theories

**Goal:** Run targeted tests to confirm/deny hypotheses

**Actions:**
1. Design minimal experiment
2. Run ONLY the experiment
3. Record result objectively
4. Move to next hypothesis if denied

**Output:** Experiment results with pass/fail

### Phase 4: CONFIRM - Verify Root Cause

**Goal:** Confirm root cause and verify fix

**Actions:**
1. Reproduce bug with root cause isolated
2. Implement minimal fix
3. Verify fix resolves issue
4. Add regression test

**Output:** Root cause report + fix

---

## When to Use

- Tests failing unexpectedly
- Production incidents
- Bug reports with unclear cause
- Performance degradation
- Integration failures

---

## Common Pitfalls

| Pitfall | Problem |
|---------|---------|
| Skipping observation | Jumping to conclusions |
| Testing multiple things | Can't isolate cause |
| Confirmation bias | Only looking for proving evidence |
| Stopping at first fix | Not understanding WHY it worked |

---

## Exit Criteria

- Root cause identified
- Fix implemented and verified
- Regression test added

---

## See Also

- `@bugfix` - Quality bug fixes (P1/P2)
- `@hotfix` - Emergency fixes (P0)
- `@issue` - Bug classification and routing
