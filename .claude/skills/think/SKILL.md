---
name: think
description: Deep structured thinking for complex problems (INTERNAL - used by @idea and @design)
tools: Read, Write, Edit, Bash, Glob, Grep, AskUserQuestion
---

# /think - Deep Structured Thinking

**INTERNAL SKILL** — Used by `@idea` and `@design` for deep analysis of requirements and architecture.

## Purpose

When a problem needs deeper exploration than surface-level questions:
- Complex tradeoffs with no clear answer
- Architectural decisions with multiple valid approaches
- Unknown unknowns in requirements
- System-level implications

## When @idea or @design Should Call This

**@idea calls /think when:**
- Requirements have significant ambiguity
- Multiple user types with conflicting needs
- Technical approach unclear
- Success metrics debatable

**@design calls /think when:**
- Architecture has multiple valid approaches
- Integration points complex
- Failure modes unclear
- Performance/security tradeoffs significant

## The Thinking Process

### Phase 1: Deconstruct the Problem

Break down into dimensions:
- **Users** — Who is affected? How?
- **Scope** — What's in/out?
- **Constraints** — Technical, business, time
- **Risks** — What could go wrong?
- **Dependencies** — What else does this touch?

### Phase 2: Explore Multiple Angles

For each significant dimension, explore:

**Angle 1: The Ideal Solution**
- What would perfect look like?
- What constraints prevent it?

**Angle 2: The Pragmatic Solution**
- What's realistic given constraints?
- What are acceptable tradeoffs?

**Angle 3: The Minimal Solution**
- What's the smallest thing that works?
- What can be deferred?

### Phase 3: Synthesize Insights

Look for:
- Patterns across angles
- Hidden assumptions
- Risks not yet considered
- Dependencies not yet mapped

### Phase 4: Present Findings

Format:
```markdown
## Problem Analysis

### Dimensions
- Users: ...
- Scope: ...
- Constraints: ...

### Angles Explored
1. **Ideal**: ... → Blocked by ...
2. **Pragmatic**: ... → Tradeoffs: ...
3. **Minimal**: ... → Deferred: ...

### Key Insights
- ...
- ...

### Recommendation
With the above analysis, recommended approach: ...

### Open Questions
- ...
```

## Principles

- **Explore before concluding** — Don't jump to answers
- **Multiple angles** — At least 3 perspectives per major decision
- **Explicit tradeoffs** — State what you're optimizing for
- **Surface assumptions** — Make implicit things explicit
- **Identify unknowns** — Distinguish known from unknown

## Exit When

- Analysis covers all major dimensions
- At least 2-3 angles explored for each decision
- Tradeoffs are explicit
- Recommendation is clear with rationale
