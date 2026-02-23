---
name: design
description: System design with progressive disclosure
version: 6.0.0
changes:
  - Converted to LLM-agnostic format
  - Removed tool-specific API references
  - Focus on WHAT, not HOW to invoke
---

# @design - System Design with Progressive Disclosure

Multi-agent system design (Arch + Security + SRE) with progressive discovery blocks.

---

## EXECUTE THIS NOW

When user invokes `@design <task_id>`:

### Step 1: Load Requirements

Load requirements from:
- `docs/intent/{task_id}.json` - Machine-readable intent from @idea
- `docs/drafts/idea-*.md` - Feature spec from @idea

Skip topics already covered by @idea.

### Step 2: Progressive Discovery (3-5 blocks)

**Block Structure:**
- Each block: 3 questions
- After each block: trigger point (Continue / Skip block / Done)
- User can skip blocks not relevant to feature

**Discovery Blocks:**

**Block 1: Data & Storage (3 questions)**
- Data models?
- Storage requirements?
- Persistence strategy?

**Block 2: API & Integration (3 questions)**
- API endpoints?
- External integrations?
- Authentication/authorization?

**Block 3: Architecture (3 questions)**
- Component structure?
- Layer boundaries?
- Error handling strategy?

**Block 4: Security (3 questions)**
- Input validation?
- Sensitive data handling?
- Rate limiting?

**Block 5: Operations (3 questions)**
- Monitoring?
- Deployment?
- Rollback strategy?

**After Each Block: Trigger Point**
- Continue (next discovery block)
- Skip block (skip remaining blocks)
- Done (generate workstreams with current info)

### Step 3: Workstream Generation

Generate workstreams based on:
- Architecture decisions (from discovery blocks)
- Quality gates (TDD, coverage, type hints)

**Output:** Workstream files in `docs/workstreams/backlog/00-FFF-SS.md`

### Step 4: Create Beads Tasks

```bash
bd create --title="WS-FFF-01: {title}" --type=task --priority=2
```

---

## When to Use

- After @idea requirements gathering
- Need architecture decisions
- Creating workstream breakdown

---

## Modes

| Mode | Blocks | Purpose |
|------|--------|---------|
| Default | 3-5 | Full discovery |
| `--quiet` | 2 | Minimal (Data + Architecture) |

---

## --quiet Mode

Minimal blocks (2 blocks, 6 questions):
1. Data & Storage
2. Core Architecture

---

## Output

**Primary:** Workstream files in `docs/workstreams/backlog/`

**Secondary:**
- `docs/drafts/<task_id>-design.md` - Design document

---

## Next Steps

```bash
@oneshot <feature>  # Execute all workstreams
@build <ws_id>      # Execute single workstream
```

---

---

## See Also

- `@idea` - Requirements gathering
- `@build` - Execute workstream
- `@oneshot` - Execute all workstreams
- `@feature` - Orchestrator that calls @idea + @design
