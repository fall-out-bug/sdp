---
name: idea
description: Interactive requirements gathering. Structured dialogue to create feature draft with goals, scope, and open questions.
tools: Read, Write, Edit, Bash, Glob, Grep
---

# /idea - Requirements Gathering

Interactive dialogue to capture feature requirements.

## When to Use

- Starting new feature
- Unclear requirements
- Need structured requirements document

## Invocation

```bash
/idea "feature description"
```

## Master Prompt

📄 **sdp/prompts/commands/idea.md** (220+ lines)

**Contains:**
- Interactive dialogue questions
- Structured draft template
- Goals vs Non-Goals
- User stories
- Open questions
- Output location

## Workflow

1. Ask clarifying questions:
   - Problem being solved?
   - Target users?
   - Success criteria?
   - Out of scope?
2. Create draft file
3. Ask follow-up questions

## Output

Draft: `docs/drafts/idea-{slug}.md`

Sections:
- Context & Problem
- Goals & Non-Goals
- User Stories
- Scope
- Assumptions
- Open Questions

## Quick Reference

**Input:** Feature idea  
**Output:** `docs/drafts/idea-{slug}.md`  
**Next:** `/design idea-{slug}`
