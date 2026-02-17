---
name: think
description: Deep structured thinking with parallel agents (INTERNAL)
tools:
  - Read
  - Write
  - Bash
  - Task
---

# /think - Deep Structured Thinking

**INTERNAL SKILL** — Used by `@idea` and `@design` for deep analysis.

## When to Use

- Complex tradeoffs with no clear answer
- Architectural decisions with multiple valid approaches
- Unknown unknowns in requirements
- System-level implications

## Parallel Expert Agents Pattern

### Step 1: Define Expert Roles

| Expert | Focus | When to Use |
|--------|-------|-------------|
| **Architect** | System design, patterns | All architectural decisions |
| **Security** | Threats, auth, data | User data, APIs, external integration |
| **Performance** | Latency, scalability | High load, real-time |
| **UX** | User experience | User-facing features |
| **Ops** | Deployability, monitoring | Production systems |

### Step 2: Launch Parallel Analysis

```python
# Spawn 2-4 experts in parallel (single message)
Task(
    subagent_type="general-purpose",
    prompt="""You are the ARCHITECT expert.
    
PROBLEM: {problem}

Analyze from your perspective:
1. Key considerations?
2. Applicable patterns?
3. Risks?

Return 3-5 bullet points.""",
    description="Architect analysis"
)
# Launch other experts similarly...
```

### Step 3: Synthesize

After all experts complete:

```markdown
## Expert Analysis

**@architect:** Domain layer first, risk of tight coupling
**@security:** OAuth2 preferred, need rate limiting
**@performance:** Caching needed, ~500MB for 10K users

## Synthesis
Recommended approach combining all perspectives...

## Open Questions
What remains unknown...
```

## Single-Agent Mode (Simple Problems)

1. **Deconstruct** problem into dimensions
2. **Explore** 3+ angles (ideal/pragmatic/minimal)
3. **Synthesize** insights
4. **Present** findings with tradeoffs

## Output Format

```markdown
## Problem Analysis

### Context
{Brief problem statement}

### Expert Analysis
**@architect:** {analysis}
**@security:** {analysis}

### Synthesis
{Combined insights}

### Recommendation
{Clear recommendation with rationale}

### Open Questions
{What remains unknown}
```

## Principles

- **Parallel exploration** — Multiple experts simultaneously
- **Role-based expertise** — Each expert has defined perspective
- **Explicit tradeoffs** — State what you're optimizing for
- **Clear recommendation** — Don't leave user hanging
