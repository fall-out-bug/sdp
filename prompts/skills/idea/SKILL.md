---
name: idea
description: Interactive requirements gathering with progressive disclosure
version: 5.0.0
changes:
  - Converted to LLM-agnostic format
  - Removed tool-specific API references
  - Focus on WHAT, not HOW to invoke
---

# @idea - Requirements Gathering with Progressive Disclosure

Deep interviewing to capture comprehensive feature requirements using progressive disclosure (3-question cycles). Creates markdown spec, optionally creates Beads task.

---

## EXECUTE THIS NOW

When user invokes `@idea "feature description"`:

### Step 1: Read Context

Read existing project files to understand context:
- `PRODUCT_VISION.md` - Align with project goals
- `docs/specs/**/*` - Similar features

### Step 2: Progressive Interview (3-Question Cycles)

**Question Target:**
- Minimum: 12 questions (bounded exploration)
- Maximum: 27 questions (deep analysis)
- Average: 18-20 questions per feature

**3-Question Cycles:**

1. Ask 3 focused questions
2. Offer trigger point after each cycle
3. User chooses: continue / deep design / skip to @design

**Cycle 1 - Vision (3 questions):**
- What is the core mission of this feature?
- How does this align with PRODUCT_VISION.md?
- Who are the primary users?

**TRIGGER POINT (after each cycle):**
- Continue (more questions)
- Deep design (jump to @design with architectural exploration)
- Skip to @design (move to workstream decomposition)

**Cycle 2 - Problem & Users (3 questions):**
- What problem does this solve?
- What are the user pain points?
- What happens if we don't build this?

**Cycle 3 - Technical Approach (3 questions):**
- Storage/data requirements?
- Failure modes to handle?
- Integration points?

**Cycle 4 - UI/UX & Quality (3 questions):**
- UI/UX requirements?
- Performance targets?
- Security considerations?

**Cycle 5 - Testing & Edge Cases (3 questions):**
- Testing strategy?
- Edge cases to handle?
- Success metrics?

### Step 3: TMI Detection

If user provides extensive detail upfront:
- "detailed spec", "full implementation", "complete architecture"
- User writes >500 characters in initial prompt

Offer shortcuts:
- Continue with targeted questions (recommended)
- Skip to @design with detailed spec
- Use --quiet mode for minimal questions

### Step 4: Create Beads Task

```bash
bd create --title="{feature_title}" --type=feature --priority=2
```

Include in description:
- Context & Problem
- Goals & Non-Goals
- Technical Approach
- Concerns & Tradeoffs

### Step 5: Create Spec File

Create `docs/intent/{task_id}.json` with machine-readable intent.

---

## When to Use

- Starting new feature
- Unclear requirements
- Need comprehensive spec with tradeoffs explored

---

## Modes

| Mode | Questions | Purpose |
|------|-----------|---------|
| Default | 12-27 | Full progressive interview |
| `--quiet` | 3-5 | Minimal questions (core only) |
| `--spec path` | Varies | Use existing spec as base |

---

## --quiet Mode

Minimal questions (3-5 core only):
1. Mission?
2. Users?
3. Core requirement?

Skip deep-dive cycles, move directly to @design.

---

## Output

**Primary:** Beads task ID (e.g., `sdp-xxx`)

**Secondary:**
- `docs/intent/{task_id}.json` - Machine-readable intent
- Question count included in metadata

---

## Next Steps

```bash
@design sdp-xxx      # Decompose into workstreams
bd show sdp-xxx      # View task details
bd ready             # Check ready tasks
```

---

## Key Principles

1. **Progressive disclosure** - 3 questions at a time
2. **User-controlled depth** - trigger points after each cycle
3. **Respect brevity** - --quiet mode for experienced users
4. **No obvious questions** - explore tradeoffs, not yes/no
5. **TMI detection** - offer shortcuts when user over-explains

---

## Example Session

```
@idea "Add user authentication"

# Cycle 1: Vision (3 questions)
[Mission] What is the core mission?
[Alignment] How does this align with vision?
[Users] Who are the primary users?

# TRIGGER: Continue? (yes/deep design/skip)
User selects: Continue

# Cycle 2: Problem (3 questions)
...

# TRIGGER: Continue? (yes/deep design/skip)
User selects: Deep design

# Jump to @design with architectural exploration

Created Beads task: sdp-xxx
   Title: Add user authentication
   Questions asked: 6
   Priority: P2

# Next:
@design sdp-xxx
```

---

## Quick Reference

| Command | Purpose |
|---------|---------|
| `@idea "feature"` | Create task with progressive interview |
| `@idea "feature" --quiet` | Minimal questions (3-5 core only) |
| `bd show {id}` | View task details |
| `@design {id}` | Decompose into workstreams |

---

## See Also

- `@design` - Workstream decomposition
- `@build` - Execute workstream
- `@feature` - Orchestrator that calls @idea + @design
