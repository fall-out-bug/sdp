---
name: think
description: Deep structured thinking with parallel agents (INTERNAL - used by @idea and @design)
tools: Read, Write, Edit, Bash, Glob, Grep, AskUserQuestion, Task
---

# /think - Deep Structured Thinking

**INTERNAL SKILL** — Used by `@idea` and `@design` for deep analysis with parallel expert agents.

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

## Parallel Expert Agents Pattern

### Step 1: Define Expert Roles

For complex problems, spawn 2-4 parallel expert agents:

| Expert | Focus | When to Use |
|--------|-------|-------------|
| **Architect** | System design, patterns, modularity | All architectural decisions |
| **Security** | Threats, auth, data protection | User data, APIs, external integration |
| **Performance** | Latency, throughput, scalability | High load, real-time requirements |
| **UX** | User experience, discoverability | User-facing features |
| **Ops** | Deployability, monitoring, failure modes | Production systems |

### Step 2: Launch Parallel Analysis

```python
# Spawn experts in parallel (single message)
Task(
    subagent_type="general-purpose",
    prompt="""You are the ARCHITECT expert for this problem.

PROBLEM: {problem_description}

Your expertise: System design, patterns, modularity, clean architecture

Explore the problem from your perspective:
1. What are the key architectural considerations?
2. What patterns apply here?
3. What are the risks?

Return your analysis in 3-5 bullet points.""",
    description="Architect analysis"
)

# Launch other experts similarly...
```

**All agents run in parallel** — user sees all thoughts simultaneously.

### Step 3: Aggregate and Synthesize

After all experts complete, synthesize:

```markdown
## Expert Analysis Summary

### @architect
- Clean architecture suggests domain layer first
- Risk: tight coupling to existing services

### @security
- Need threat modeling for user data
- OAuth2 preferred over custom auth

### @performance
- Caching strategy needed for read-heavy workload
- Database indexing critical for query performance

### Synthesis
Combining all perspectives, recommended approach: ...
```

## Single-Agent Mode (Simple Problems)

For straightforward problems, skip parallel agents:

1. **Deconstruct** the problem into dimensions
2. **Explore** 3+ angles (ideal/pragmatic/minimal)
3. **Synthesize** insights
4. **Present** findings

## Output Format

```markdown
## Problem Analysis

### Context
{Brief problem statement}

### Expert Analysis

**@architect:** {analysis}
**@security:** {analysis}
**@performance:** {analysis}

### Synthesis
{Combined insights}

### Recommendation
{Clear recommendation with rationale}

### Open Questions
{What remains unknown}
```

## Principles

- **Parallel exploration** — Multiple experts run simultaneously
- **Real-time thoughts** — User sees all expert thoughts as they complete
- **Role-based expertise** — Each expert has defined perspective
- **Synthesis** — Combine insights into coherent recommendation
- **Explicit tradeoffs** — State what you're optimizing for

## Exit When

- All experts have completed analysis (parallel mode) OR
- All angles explored (single mode)
- Tradeoffs are explicit
- Recommendation is clear with rationale
