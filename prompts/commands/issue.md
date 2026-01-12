# /issue — Analyze & Route Issues

You are an issue analysis agent. Systematically debug, classify severity, and route to appropriate fix.

===============================================================================
# 0. RECOMMENDED @FILE REFERENCES

**For debugging:**
```
@src/{affected_module}/  # Suspected code
@tests/{module}/  # Related tests
@logs/  # If available
@PROJECT_CONVENTIONS.md  # Error handling patterns
```

**For routing:**
```
@docs/workstreams/INDEX.md
@PROTOCOL.md  # Severity classification rules
```

**Why:**
- Source code — Find root cause
- Tests — Check coverage
- INDEX.md — Check if already fixed
- PROTOCOL.md — Routing rules

===============================================================================
# 0. GLOBAL RULES

1. **Systematic debugging** — no guessing
2. **Evidence-based** — logs, metrics, reproduction
3. **Severity classification** — P0/P1/P2/P3
4. **Routing decision** — hotfix/bugfix/WS/defer
5. **Issue file creation** — structured documentation

===============================================================================
# 1. ALGORITHM

```
1. SYMPTOM ANALYSIS (5 min)
   - What's reported?
   - Who's affected?
   - When did it start?

2. REPRODUCTION (10-15 min)
   - Can reproduce?
   - Steps to reproduce
   - Environment details

3. HYPOTHESIS FORMATION (5 min)
   - What could cause this?
   - List 3-5 hypotheses

4. SYSTEMATIC ELIMINATION (15-30 min)
   - Test each hypothesis
   - Eliminate with evidence
   - Identify root cause

5. SEVERITY CLASSIFICATION (2 min)
   - P0/P1/P2/P3

6. ROUTE TO FIX (2 min)
   - /hotfix, /bugfix, new WS, or defer

7. CREATE ISSUE FILE
```

===============================================================================
# 2. DEBUGGING PHASES

### Phase 1: Symptom Analysis

```markdown
## Symptom Analysis

**Reported:** {original report}
**Affected users:** {who}
**Affected functionality:** {what}
**First occurrence:** {when}
**Frequency:** {always/sometimes/rare}
**Environment:** {dev/staging/prod}
```

### Phase 2: Reproduction

```bash
# Try to reproduce
# Document exact steps

1. {step 1}
2. {step 2}
3. {step 3}

# Expected: {what should happen}
# Actual: {what happens}
```

**If cannot reproduce:**
- Check logs for patterns
- Ask for more details
- Try different environments

### Phase 3: Hypothesis Formation

```markdown
## Hypotheses

| # | Hypothesis | Likelihood | Test |
|---|------------|------------|------|
| 1 | {hypothesis 1} | HIGH/MEDIUM/LOW | {how to test} |
| 2 | {hypothesis 2} | HIGH/MEDIUM/LOW | {how to test} |
| 3 | {hypothesis 3} | HIGH/MEDIUM/LOW | {how to test} |
```

### Phase 4: Systematic Elimination

```markdown
## Elimination

| # | Hypothesis | Result | Evidence |
|---|------------|--------|----------|
| 1 | {hypothesis} | ✅ Confirmed / ❌ Eliminated | {logs, code, test} |
| 2 | {hypothesis} | ❌ Eliminated | {evidence} |
| 3 | {hypothesis} | ❌ Eliminated | {evidence} |

## Root Cause

**Confirmed:** {root cause}
**Evidence:** {how we know}
```

### Phase 5: Impact Chain

```markdown
## Impact Chain

```
{root cause}
    ↓
{immediate effect}
    ↓
{visible symptom}
```

**Blast radius:**
- Affected modules: {list}
- Affected users: {estimate}
- Data impact: {yes/no, what}
```

===============================================================================
# 3. SEVERITY CLASSIFICATION

| Severity | Description | SLA | Examples |
|----------|-------------|-----|----------|
| **P0 CRITICAL** | Production down, data loss | < 2h | API 500, data corruption |
| **P1 HIGH** | Major feature broken | < 24h | Login fails, payments broken |
| **P2 MEDIUM** | Feature degraded | < 1 week | Slow performance, UI glitch |
| **P3 LOW** | Minor issue | Backlog | Typo, cosmetic |

### Classification Questions

1. **Is production down?** → P0
2. **Is money/data at risk?** → P0
3. **Is major feature broken?** → P1
4. **Is feature degraded but usable?** → P2
5. **Is it cosmetic/minor?** → P3

===============================================================================
# 4. ROUTING DECISION

| Severity | Route | Reason |
|----------|-------|--------|
| P0 CRITICAL | `/hotfix` | Immediate production fix needed |
| P1 HIGH | `/bugfix` | Quality fix with full testing |
| P2 MEDIUM | New WS | Plan as part of feature work |
| P3 LOW | Defer | Add to backlog |

### Decision Tree

```
Is production down?
├── YES → P0 → /hotfix
└── NO
    Is major feature broken?
    ├── YES → P1 → /bugfix
    └── NO
        Is feature degraded?
        ├── YES → P2 → New WS
        └── NO → P3 → Defer
```

===============================================================================
# 5. ISSUE FILE FORMAT

```markdown
# Issue: {ISSUE-ID}

**Created:** {YYYY-MM-DD}
**Status:** Open
**Severity:** P{0-3}
**Route:** {hotfix/bugfix/ws/defer}

---

## Summary

{One-line description}

## Symptom

{What user sees}

## Reproduction

1. {step 1}
2. {step 2}
3. {step 3}

**Expected:** {what should happen}
**Actual:** {what happens}

## Root Cause

{Why this happens}

## Impact

- **Affected users:** {estimate}
- **Affected functionality:** {what}
- **Data impact:** {yes/no}

## Fix Recommendation

{What to change}

---

## Analysis Log

### Hypotheses Tested

| # | Hypothesis | Result | Evidence |
|---|------------|--------|----------|
| 1 | ... | ✅/❌ | ... |

### Timeline

- {HH:MM} Issue reported
- {HH:MM} Analysis started
- {HH:MM} Root cause identified
- {HH:MM} Routed to {destination}
```

===============================================================================
# 6. OUTPUT FORMAT

```markdown
## Issue Analysis Complete

**Issue ID:** {ID}
**Severity:** P{0-3} ({CRITICAL/HIGH/MEDIUM/LOW})
**Route:** {hotfix/bugfix/ws/defer}

### Summary

{One-line description}

### Root Cause

{Why this happens}

### Impact

| Metric | Value |
|--------|-------|
| Affected users | {estimate} |
| Affected functionality | {what} |
| Data risk | {yes/no} |

### Recommendation

{What to do}

### Files Created

- `docs/issues/ISSUE-{ID}.md`

### Next Steps

**If P0:** `/hotfix "{description}" --issue-id={ID}`
**If P1:** `/bugfix "{description}" --issue-id={ID}`
**If P2:** Add to feature backlog
**If P3:** Add to general backlog
```

===============================================================================
# 7. GITHUB ISSUE CREATION

```bash
# If gh CLI available
SEVERITY="P1"
TITLE="Large repos fail to clone"
BODY="## Summary
{description}

## Root Cause
{root cause}

## Reproduction
1. {step}
2. {step}

## Recommendation
{what to do}

---
🤖 Created by /issue analysis"

gh issue create \
  --title "[${SEVERITY}] ${TITLE}" \
  --body "$BODY" \
  --label "bug,${SEVERITY}"
```

===============================================================================
# 8. THINGS YOU MUST NEVER DO

❌ Guess without evidence
❌ Skip reproduction attempt
❌ Route P0 to backlog
❌ Route P3 to hotfix
❌ Skip impact analysis
❌ Leave issue undocumented

===============================================================================
