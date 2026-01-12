---
name: design
description: Analyze requirements and plan workstreams. Decomposes features into executable WS with clear goals and dependencies.
tools: Read, Write, Edit, Bash, Glob, Grep
---

# /design - Analyze + Plan

Transform idea drafts or feature specs into detailed workstream specifications.

## When to Use

- After `/idea` creates a draft
- When a feature spec needs to be broken down
- Before starting implementation

## Invocation

```bash
/design idea-{slug}
# or
/design feature-{XX}
```

## Workflow

**IMPORTANT:** This skill delegates to the master prompt.

### Step 1: Load Master Prompt

```bash
cat sdp/prompts/commands/design.md
```

**This file contains:**
- Complete algorithm
- Pre-flight checks (PROJECT_MAP, INDEX, duplicates)
- WS decomposition rules
- Scope limits (SMALL/MEDIUM/LARGE)
- Git workflow (GitFlow)
- GitHub integration
- Output format

### Step 2: Execute Instructions

Follow **all steps** in `sdp/prompts/commands/design.md`:

1. Read mandatory files (PROJECT_MAP.md, INDEX.md)
2. Analyze draft/spec
3. Determine WS count and dependencies
4. Create all WS files
5. Update INDEX.md
6. Git workflow (branch, commit, push)
7. GitHub integration (if `gh` available)
8. Output summary

### Step 3: Verification

Before finishing:
- [ ] All WS files created (no references to non-existent files)
- [ ] All WS scope ≤ MEDIUM
- [ ] INDEX.md updated
- [ ] Git commit created
- [ ] GitHub issues created (if gh available)

## Key Rules

From master prompt:

1. **Read PROJECT_MAP.md FIRST** - all architecture decisions
2. **Check INDEX.md** - prevent duplicates
3. **Create ALL WS files** - no dangling references
4. **Scope ≤ MEDIUM** - split if larger
5. **No time estimates** - only LOC/tokens
6. **GitFlow** - feature branch from develop
7. **GitHub sync** - create issues for WS

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
