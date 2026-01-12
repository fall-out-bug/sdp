---
name: review
description: Code review with metrics-based quality checks. Reviews entire features, enforces quality gates, generates UAT guide.
tools: Read, Write, Edit, Bash, Glob, Grep
---

# /review - Code Review

Comprehensive code review for features or workstreams with strict quality gates.

## When to Use

- After all WS in a feature are completed
- Before human UAT
- Part of `/oneshot` flow
- To verify quality standards

## Invocation

```bash
/review F60         # Review entire feature
/review WS-060      # Review all WS-060-XX
```

## Workflow

**IMPORTANT:** This skill delegates to the master prompt.

### Load Master Prompt

```bash
cat sdp/prompts/commands/review.md
```

**This file contains:**
- 17-point quality checklist
- Metrics-based validation (coverage, complexity, LOC)
- Goal achievement verification
- Cross-WS consistency checks
- UAT guide generation
- Delivery notification template
- Verdict rules (APPROVED / CHANGES_REQUESTED only)

### Execute Instructions

Follow `sdp/prompts/commands/review.md`:

1. Find all WS in scope
2. For each WS:
   - Check 0: Goal achieved? (BLOCKING)
   - Checks 1-17: tests, coverage, complexity, etc.
   - Collect metrics
3. Cross-WS checks
4. Generate UAT guide
5. Send notification (if blockers)
6. Output verdict

## Key Checks

From master prompt:

- **Check 0:** Goal Achievement (100% AC) 🔴 BLOCKING
- **Check 2:** Test Coverage (≥80%) 
- **Check 4:** AI-Readiness (file <200 LOC, CC <10)
- **Check 5:** Clean Architecture (no violations)
- **Check 6:** Type Hints (100% coverage)
- **Check 9:** TODO/FIXME (must be 0)
- **Check 12:** Git History (commits present)
- **Check 14:** Human Verification (UAT) guide

## Metrics Summary Table

Collected for each WS:

| Check | Target | Actual | Status |
|-------|--------|--------|--------|
| Goal Achievement | 100% | {X/Y}% | ✅/🔴 |
| Test Coverage | ≥80% | {N}% | ✅/⚠️/🔴 |
| Cyclomatic Complexity | <10 | avg {N} | ✅/⚠️/🔴 |
| File Size | <200 LOC | max {N} | ✅/⚠️/🔴 |
| Type Hints | 100% | {N}% | ✅/🔴 |
| TODO/FIXME | 0 | {count} | ✅/🔴 |

## Verdict Rules

From master prompt:

- **APPROVED:** All checks ✅, all WS Goals achieved
- **CHANGES REQUESTED:** Any 🔴 BLOCKING issue

**NO "APPROVED WITH NOTES"** - that's a loophole for tech debt!

## Output

### Per-WS Summary

```markdown
| WS | Verdict | Goal | Coverage |
|----|---------|------|----------|
| WS-060-01 | ✅ APPROVED | ✅ | 87% |
| WS-060-02 | ❌ CHANGES REQUESTED | ❌ AC2 | 75% |
```

### Delivery Notification

```markdown
## ✅ Review Complete: F60

**Feature:** LMS Integration
**Status:** APPROVED
**Elapsed (telemetry):** 2h 15m

### Metrics
| Metric | Target | Actual |
|--------|--------|--------|
| Coverage | ≥80% | 86% |
| Complexity | <10 | avg 4.8 |

### Impact
{Business impact statement}

### Next Steps
1. Human UAT (5-10 min)
2. `/deploy F60` if UAT passes
```

## UAT Guide

Generated at: `docs/uat/F{XX}-uat-guide.md`

Sections:
- Quick Smoke Test (30 sec)
- Detailed Scenarios (5-10 min)
- Red Flags checklist
- Sign-off

## Master Prompt Location

📄 **sdp/prompts/commands/review.md** (460+ lines)

## Quick Reference

**Input:** Feature ID or WS prefix  
**Output:** Verdict + UAT Guide + Metrics  
**Next:** Human UAT → `/deploy F{XX}`
