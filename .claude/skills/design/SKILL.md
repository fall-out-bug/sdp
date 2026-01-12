---
name: design
description: Analyze requirements and plan workstreams. Decomposes features into executable WS with clear goals and dependencies.
tools: Read, Write, Edit, Bash, Glob, Grep, EnterPlanMode, AskUserQuestion
---

# /design - Analyze + Plan

Transform idea drafts or feature specs into detailed workstream specifications using interactive planning mode.

## When to Use

- After `/idea` creates a draft
- When a feature spec needs to be broken down
- Before starting implementation
- When architectural decisions need user input

## Invocation

```bash
/design idea-{slug}
# or
/design feature-{XX}
```

## Workflow

**IMPORTANT:** This skill uses EnterPlanMode for interactive codebase exploration and planning.

### Step 0: Enter Plan Mode

**Always start** with EnterPlanMode for codebase analysis:

```markdown
EnterPlanMode()
```

This transitions to plan mode where you can:
- Explore codebase thoroughly (Glob, Grep, Read)
- Understand existing patterns and architecture
- Use AskUserQuestion for architectural decisions
- Design implementation approach without writing code

### Step 1: Exploration Phase (In Plan Mode)

**Read mandatory context:**
```bash
# Core documentation
@PROJECT_MAP.md          # Architecture decisions
@docs/workstreams/INDEX.md  # Existing workstreams
@PROTOCOL.md             # WS rules and sizing
@CODE_PATTERNS.md        # Code patterns

# Feature specification
@docs/drafts/idea-{slug}.md
# or
@docs/specs/feature_{XX}/feature.md
```

**Explore relevant code:**
- Use Glob to find similar features
- Use Grep to understand patterns
- Read existing implementations
- Identify integration points

### Step 2: Interactive Planning (In Plan Mode)

**Use AskUserQuestion** for architectural decisions:

```markdown
AskUserQuestion({
  "questions": [{
    "question": "Which authentication approach should we use for this feature?",
    "header": "Auth Method",
    "options": [
      {
        "label": "JWT tokens (Recommended)",
        "description": "Stateless, scalable. Requires token refresh logic."
      },
      {
        "label": "Session-based",
        "description": "Simpler, but requires session storage. Less scalable."
      },
      {
        "label": "OAuth2 integration",
        "description": "Delegates to external provider. More setup complexity."
      }
    ],
    "multiSelect": false
  }]
})
```

**Ask about:**
- Architecture patterns (which layer structure?)
- External dependencies (which libraries?)
- Data storage (database schema decisions)
- API design (REST vs GraphQL, endpoint structure)
- Testing strategy (unit vs integration scope)
- Deployment approach (where to deploy, what infrastructure?)

### Step 3: Design Workstreams (In Plan Mode)

**Based on exploration and user answers:**

1. Decompose feature into WS (5-30 workstreams)
2. Ensure each WS is SMALL or MEDIUM scope
3. Define dependencies between WS
4. Write clear Goals and Acceptance Criteria
5. Consider integration points

**WS Decomposition Pattern:**
```
WS-XXX-01: Domain layer (entities, value objects)
WS-XXX-02: Application layer (use cases, ports)
WS-XXX-03: Infrastructure adapters
WS-XXX-04: Presentation layer (API/UI)
WS-XXX-05: Integration tests
```

### Step 4: Exit Plan Mode

**When design is complete:**

```markdown
ExitPlanMode()
```

This will:
- Present plan to user for approval
- Wait for user confirmation
- Exit to execution mode

### Step 5: Execute Plan (After Approval)

**Load master prompt:**
```bash
cat prompts/commands/design.md
```

**Follow all steps:**
1. Create all WS files in `docs/workstreams/backlog/`
2. Update INDEX.md
3. Git workflow (branch, commit, push)
4. GitHub integration (create issues if `gh` available)
5. Output summary

### Step 6: Verification

Before finishing:
- [ ] All WS files created (no references to non-existent files)
- [ ] All WS scope ≤ MEDIUM
- [ ] INDEX.md updated
- [ ] Git commit created
- [ ] GitHub issues created (if gh available)

## Key Rules

**Planning Mode (EnterPlanMode):**

1. **ALWAYS start with EnterPlanMode** - before any WS creation
2. **Explore thoroughly** - use Glob, Grep, Read to understand codebase
3. **Ask questions** - use AskUserQuestion for architectural decisions
4. **Design first, code later** - complete planning before ExitPlanMode
5. **Get approval** - ExitPlanMode requests user confirmation

**Execution Mode (After ExitPlanMode):**

6. **Read PROJECT_MAP.md FIRST** - all architecture decisions
7. **Check INDEX.md** - prevent duplicates
8. **Create ALL WS files** - no dangling references
9. **Scope ≤ MEDIUM** - split if larger
10. **No time estimates** - only LOC/tokens
11. **GitFlow** - feature branch from develop
12. **GitHub sync** - create issues for WS

## Output

Delegate output format to master prompt. Should include:
- Created WS list
- Dependency graph
- Git info (branch, commit)
- GitHub info (issues created)
- Next steps

## Master Prompt Location

📄 **sdp/prompts/commands/design.md** (496 lines)

**Why reference?**
- Single source of truth
- No synchronization issues
- Easy updates
- Consistent behavior across Cursor and Claude Code

## Quick Reference

**Input:** `docs/drafts/idea-{slug}.md` or `docs/specs/feature_XX/feature.md`  
**Output:** `docs/workstreams/backlog/WS-XXX-*.md` + INDEX.md  
**Next:** `/build WS-XXX-01`
