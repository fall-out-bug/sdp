# Oneshot Skill Failure Analysis

> **Status:** Research complete
> **Date:** 2026-02-13
> **Goal:** Понять, почему @oneshot F067 не выполнился автоматически

---

## Overview

### Goals

1. **Identify root cause** — почему я не запустил Task agent после загрузки skill
2. **Fix skill structure** — сделать инструкции однозначными
3. **Document pattern** — чтобы другие skills не имели той же проблемы

### Key Decisions

| Aspect | Decision |
|--------|----------|
| Skill Invocation | Skill tool только показывает контент, AI должен сам выполнить инструкции |
| SKILL.md Structure | Python код в skill — это ПРИМЕР/ДОКУМЕНТАЦИЯ, не автоматическое выполнение |
| Orchestrator Reference | `**READ FIRST:**` buried в code block, не prominent |
| Skill vs Agent | Skill = orchestration docs, Agent = subagent prompt |
| Ambiguity | Три цели смешаны: education, action, reference |

---

## 1. Skill Invocation Flow

> **Expert:** Theo Browne (API Design)

### Problem

Когда Skill tool возвращает контент SKILL.md, он НЕ автоматически выполняет Task(). AI видит Python code block и интерпретирует его как документацию/пример.

### Root Cause

```markdown
## Orchestrator Agent Prompt
```python
Task(
    subagent_type="general-purpose",
    prompt=...
)
```

Этот блок выглядит как "вот как работает orchestrator" (описание), а не "ЗАПУСТИ ЭТО" (инструкция).

### Solution

Добавить явную инструкцию в конец SKILL.md:

```markdown
---

**Execution Instructions:**

After displaying this skill content, you MUST immediately invoke the orchestrator:

```python
Task(
    subagent_type="general-purpose",
    prompt="""[orchestrator prompt]"""
)
```

DO NOT stop after displaying this skill. The Task() call above is not an example - it is the instruction you must execute.
```

---

## 2. SKILL.md Structure

> **Expert:** Dan Abramov (React/Architecture)

### Problem

Python code blocks в SKILL.md смешивают две цели:
- **Документация**: "вот как система работает"
- **Инструкция**: "сделай это"

### Current Pattern

```
tools: Task, Read, Bash    ← declaration (permissions)
## Orchestrator Agent Prompt
```python                   ← is this example or instruction?
Task(...)
```

### Solution

Добавить explicit labels:

```markdown
## Orchestrator Agent Prompt (REFERENCE ONLY)

> **NOTE:** The Python code below shows the orchestrator structure.
> It is NOT executed automatically. You must invoke Task() yourself.
```

---

## 3. Orchestrator Reference

> **Expert:** Don Norman (UX Design)

### Problem

Инструкция `**READ FIRST:** Read(".claude/agents/orchestrator.md")` buried внутри Python code block (line 62). Люди её не видят.

### Additional Problem

Symlink `.claude/agents -> prompts/agents` может не существовать.

### Solution

1. **Fix symlink**: `ln -s prompts/agents .claude/agents`
2. **Prominent reference**: Вынести orchestrator reference в начало SKILL.md:

```markdown
# @oneshot - Autonomous Feature Execution

> **Prerequisites:** Read `.claude/agents/orchestrator.md` first for full agent behavior.

## Quick Start
...
```

---

## 4. Skill vs Agent Boundary

> **Expert:** Dan Abramov (Separation of Concerns)

### Mental Model

```
User: @oneshot F050
       ↓
Me reads: .claude/skills/oneshot/SKILL.md (orchestration docs)
       ↓
Me decides: Need to spawn orchestrator subagent
       ↓
Me calls: Task(subagent_type="general-purpose",
               prompt="Read prompts/agents/orchestrator.md...")
       ↓
Subagent: Executes workstreams autonomously
```

### Key Insight

`tools: Task, Read, Bash` — это **permissions**, не **requirements**.
Skill MAY use these tools, не MUST.

### The Contract

| Component | Role | File |
|-----------|------|------|
| Skill | Orchestration documentation | `.claude/skills/*/SKILL.md` |
| Agent | Subagent behavior prompt | `prompts/agents/*.md` |
| Me (AI) | Orchestrator — reads skill, spawns agents | — |

---

## 5. Ambiguity in Instructions

> **Expert:** Dianna Mullin (Technical Writing)

### Problem

SKILL.md смешивает три цели:
1. **User education** (how oneshot works)
2. **User action** (what AI should execute)
3. **Implementation reference** (for maintainers)

### Solution: Split into Three Sections

```markdown
# @oneshot - Autonomous Feature Execution

## Quick Start
[2-3 sentence summary + one example]

## What You Do (REQUIRED)
[Step-by-step actions with exact Task() calls]

## How It Works (CONCEPTUAL)
[Architecture diagram and flow]

## Implementation Details (REFERENCE)
[For skill maintainers]
```

### Pattern for All Skills

| Section | Purpose | Audience |
|---------|---------|----------|
| Quick Start | Get running fast | Users |
| What You Do | Exact commands | AI executors |
| How It Works | Understanding | Everyone |
| Implementation Details | Maintenance | Developers |

---

## Implementation Plan

### Phase 1: Fix oneshot SKILL.md (Immediate) ✅ DONE

- [x] Created new LLM-agnostic format with "EXECUTE THIS NOW" section
- [x] Removed Python `Task()` code blocks
- [x] Added explicit CLI reference (`sdp orchestrate`)
- [x] Created `commands.json` mapping file
- [x] Created `.claude/patterns/` directory with reusable patterns:
  - `tdd.md` - TDD pattern
  - `git-safety.md` - Git safety rules
  - `quality-gates.md` - Quality gates
  - `session-complete.md` - Session completion checklist

### Phase 2: Fix Other Orchestrator Skills ✅ DONE

- [x] `@build/SKILL.md` — converted to LLM-agnostic format
- [x] `@review/SKILL.md` — converted to LLM-agnostic format
- [ ] `@vision/SKILL.md` — needs conversion
- [ ] `@reality/SKILL.md` — needs conversion
- [ ] Other 22 skills — needs conversion

### Phase 3: Create Template

- [x] Documented the new pattern in this file
- [ ] Create formal skill template
- [ ] Add validation to CI
- [ ] Update CLAUDE.md to reference new structure

---

## Success Metrics

| Metric | Baseline | Target |
|--------|----------|--------|
| Skill execution success | 0% (didn't execute) | 100% (executes on first try) |
| Ambiguity reports | 1 (this session) | 0 |
| Time to understand skill | ~10 min investigation | <30 sec |

---

## Root Cause Summary

```
@oneshot F067 invoked
       ↓
Skill tool shows SKILL.md content
       ↓
I see: Python code block with Task()
       ↓
I interpret: "This is how oneshot works" (documentation)
       ↓
I stop: "Content displayed, waiting for next instruction"
       ↓
PROBLEM: No explicit "NOW EXECUTE THIS" signal
```

**The skill system worked correctly. The ambiguity was in the documentation structure.**

---

## Next Steps

1. **Fix oneshot SKILL.md** with explicit execution instructions ✅ DONE
2. **Test** by running `@oneshot F067` again
3. **Apply pattern** to other orchestrator skills
4. **Document** the three-section pattern in skill template

---

## Implemented Solution (2026-02-13)

### New Architecture

Created LLM-agnostic architecture that works with any LLM (Opus, GLM, Codex) in any tool (Claude Code, Cursor, Windsurf):

```
.claude/
├── commands.json          # Mapping: @command → skill file
├── skills/
│   └── oneshot.md         # Unified skill+agent
└── patterns/              # Reusable knowledge blocks
    ├── tdd.md
    ├── git-safety.md
    ├── quality-gates.md
    └── session-complete.md
```

### Key Changes

1. **No Python `Task()` blocks** - Removed tool-specific syntax
2. **CLI-first approach** - Skills reference actual Go CLI commands
3. **"EXECUTE THIS NOW" section** - Explicit instruction at top of each skill
4. **Progressive disclosure** - Quick Start → What Happens → Details
5. **Reusable patterns** - Common patterns extracted to `.claude/patterns/`

### Converted Skills

| Skill | Status | CLI Command |
|-------|--------|-------------|
| oneshot | ✅ v6.0.0 | `sdp orchestrate` |
| build | ✅ v6.0.0 | `sdp apply --ws` |
| review | ✅ v11.0.0 | `sdp quality review` |

### Remaining Work

- Convert remaining 23 skills to LLM-agnostic format
- Add CI validation for skill format
- Update agent files to match new pattern
