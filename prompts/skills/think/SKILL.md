---
name: think
description: Deep structured thinking with parallel expert analysis before implementation (INTERNAL)
---

# @think (INTERNAL)

Used by @idea, @design, @feature. Three stages: breakdown → parallel expert analysis → summary.

## Stage 1: Task Breakdown

Identify aspects to think through. Choose main expert. Output: Understanding + Aspects table (5-10 rows).

## Stage 2: Parallel Expert Analysis

Launch expert agents (max 4 concurrent) — one per aspect:

```
Task(subagent_type="expert"):
  "Aspect: [name]. Task context: [brief]. Study project and propose options."
```

Each expert: study project → apply expert principles → propose 2-4 options → make decision.

## Stage 3: Summary

When experts return, create unified document. Save to `docs/plans/YYYY-MM-DD-[topic]-design.md`. Ask: "Which aspects to discuss? Or ready to implement?"

## Single-Agent Mode (<3 aspects)

Skip parallel agents: Study → Analyze (named experts) → Propose options → Recommend.

## Principles

- Study first — read codebase before analyzing
- Named expertise — reference real expert principles
- Specificity — solutions for THIS project
- Honesty — every option has cons
- Clear recommendation — don't leave user hanging
