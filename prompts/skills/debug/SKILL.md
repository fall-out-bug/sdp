---
name: debug
description: Systematic debugging using scientific method for evidence-based root cause analysis
tools: Read, Write, Edit, Bash, Grep, Glob
---

# /debug - Systematic Debugging

Evidence-based debugging using scientific method. Not "try stuff and see" -- systematic investigation.

## When to Use

- Tests failing unexpectedly
- Production incidents
- Bug reports with unclear cause
- Performance degradation
- Integration failures

## The 4-Phase Method

### Phase 1: OBSERVE - Gather Facts

**Goal:** Collect evidence WITHOUT forming hypotheses

**Actions:**
- Read error messages/logs completely
- Check git diff for recent changes
- Verify environment (Python version, dependencies)
- Check configuration files
- Reproduce the bug consistently

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
- Design minimal experiment
- Run ONLY the experiment
- Record result objectively
- Move to next hypothesis if denied

**Output:** Experiment results with pass/fail

### Phase 4: CONFIRM - Verify Root Cause

**Goal:** Confirm root cause and verify fix

**Actions:**
- Reproduce bug with root cause isolated
- Implement minimal fix
- Verify fix resolves issue
- Add regression test

**Output:** Root cause report + fix

## Common Pitfalls

- **Skipping observation** -> jumping to conclusions
- **Testing multiple things at once** -> can't isolate cause
- **Confirmation bias** -> only looking for evidence that proves theory
- **Stopping at first fix** -> not understanding WHY it worked

## Exit When

- Root cause identified
- Fix implemented and verified
- Regression test added
