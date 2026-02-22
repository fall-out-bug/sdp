---
name: think
description: Deep structured thinking with parallel expert analysis before implementation (INTERNAL)
tools:
  - Read
  - Glob
  - Grep
  - Write
  - Shell
  - Task
  - WebSearch
  - WebFetch
---

# /think - Deep Structured Thinking

**INTERNAL SKILL** — Used by `@idea`, `@design`, and `@feature` for deep analysis.

Work in **three stages**: breakdown → parallel expert analysis → summary.

## Stage 1: Task Breakdown

Identify **aspects to think through** — parts of the task that need decisions.

Choose a **main expert** for the task as a whole.

**Output format:**

```
## Understanding the Task

[How you understood the task — 1-2 sentences]

---

### Expert Perspective

> "Analyzing as [Main Expert] because [reason]"
>
> **Principles from 3 experts:**
> 1. [Expert A]: "[principle]"
> 2. [Expert B]: "[principle]"
> 3. [Expert C]: "[principle]"

---

## Aspects to Think Through

| # | Aspect | Why Important | Expert |
|---|--------|---------------|--------|
| 1 | [Name] | [Why needs thinking] | [Who will analyze] |
| 2 | ... | ... | ... |
```

Usually 5-10 aspects. No more than 15.

### Expert Table

| Area | Expert | Principles |
|------|--------|------------|
| Go design | Rob Pike | simplicity, composition over inheritance, explicit errors |
| Distributed systems | Martin Kleppmann | eventual consistency, idempotency, partition tolerance |
| Architecture | Sam Newman | bounded context, single responsibility, loose coupling |
| K8s / DevOps | Kelsey Hightower | declarative config, immutable infrastructure, GitOps |
| API design | Theo Browne | type-safe contracts, fail fast, explicit errors |
| Database | Markus Winand | index-first thinking, avoid N+1, explain analyze |
| Refactoring | Martin Fowler | small steps, preserve behavior, extract till you drop |
| Testing | Kent C. Dodds | test behavior not implementation, colocation |
| Security | Troy Hunt | defense in depth, least privilege, validate all inputs |
| Reliability / SRE | Charity Majors | observability over monitoring, SLOs over SLAs, deploy small |
| Event-driven | Ben Stopford | event sourcing, CQRS, stream processing |
| Concurrency | Bryan Mills | share by communicating, goroutine lifecycle, cancellation |
| Vibecoding | Andrej Karpathy | prompt-first development, AI-native workflows, spec before code |
| Opencode / CLI agents | Thorsten Ball | composable CLI tools, unix philosophy for agents, stdin/stdout contracts |
| LLM orchestration | Harrison Chase | chains of thought, retrieval-augmented generation, structured output |
| Prompt engineering | Simon Willison | reproducible prompts, system prompt hygiene, tool-use patterns |
| LLMOps / Evals | Hamel Husain | eval-driven development, dataset curation, regression testing for LLMs |
| AgentOps | Shunyu Yao | ReAct loop, tool selection, agent memory and planning |
| Multi-agent systems | Andrew Ng | agentic design patterns, reflection, planning, multi-agent collaboration |
| AI safety / guardrails | Anthropic (team) | constitutional AI, RLHF, harmlessness-helpfulness tradeoff |

For other areas — find appropriate specialists yourself.

## Stage 2: Project Study + Parallel Expert Analysis

After breakdown, announce:

> "Identified N aspects. Now I'll study the project and launch experts for each."

Then launch **in parallel** expert agents — one per aspect (max 4 concurrent):

```
Task(subagent_type="expert"):
  "Aspect: [aspect name].
   Task context: [brief context].
   Study the project and propose solution options."
```

**IMPORTANT:** Launch all agents in ONE message in parallel.

### Expert Agent Workflow

Each expert agent:

1. **Studies the project** — uses Glob/Grep/Read to find relevant patterns, existing solutions, constraints
2. **Applies expert thinking** — chooses main expert + 3 additional expert principles
3. **Proposes 2-4 options** with pros/cons/when-suitable
4. **Makes a decision** for this specific project

Expert agent response format:

```
## Aspect: [aspect name]

### Project Context
[Relevant patterns, existing solutions, constraints found in codebase]

### Expert Analysis

> "Analyzing as [Main Expert] because [reason]"
>
> **Principles from 3 experts:**
> 1. [Expert A]: "[principle]"
> 2. [Expert B]: "[principle]"
> 3. [Expert C]: "[principle]"

### Solution Options

**A: [Name]**
- Essence: [description]
- Pros: [list]
- Cons: [list]
- When: [when suitable]

**B: [Name]**
...

### Decision from [Main Expert]

**Choice: [Option X]**

[Reasoning considering project context and expert principles]

**Risks:** [what to consider during implementation]
```

## Stage 3: Summary Document

When all experts return, create a **unified document**:

```markdown
# [Task Name]

> **Status:** Research complete
> **Date:** [date]
> **Goal:** [brief goal description]

---

## Overview

### Goals

1. **[Goal 1]** — description
2. **[Goal 2]** — description

### Key Decisions

| Aspect | Decision |
|--------|----------|
| [Aspect 1] | [Brief decision] |
| [Aspect 2] | [Brief decision] |

---

## 1. [Aspect Name]

> **Experts:** [Expert 1], [Expert 2], [Expert 3]

### Solution

[Detailed description of chosen option]

| Aspect | Details |
|--------|---------|
| ... | ... |

### Examples

```go
// Example code if applicable
```

---

## 2. [Next Aspect]
...

---

## Implementation Plan

### Phase 1: MVP

- [ ] Task 1
- [ ] Task 2

### Phase 2: Hardening
...

---

## Success Metrics

| Metric | Baseline | Target |
|--------|----------|--------|
| ... | — | ... |
```

**Save the document** to `docs/plans/YYYY-MM-DD-[topic]-design.md`

Then ask:

> "Summary saved to `docs/plans/...`. Which aspects to discuss further? Or ready to implement?"

## Single-Agent Mode (Simple Problems)

For problems with fewer than 3 aspects, skip parallel agents:

1. **Study** — Glob/Grep/Read relevant code
2. **Analyze** — apply expert thinking with named experts
3. **Propose** — 2-4 options with pros/cons
4. **Recommend** — clear decision with rationale

## Principles

- **Study first** — always read the codebase before analyzing
- **Named expertise** — reference real expert principles, not abstract advice
- **Specificity** — solutions for THIS project, not generic patterns
- **Honesty** — every option has cons, don't hide them
- **Parallel exploration** — multiple experts simultaneously
- **Clear recommendation** — don't leave the user hanging
- **Context** — consider what already exists in the project
