---
description: TDD cycle (Red-Green-Refactor) for test-driven development.
agent: builder
---

# /test — TDD

## Overview

This command implements the TDD skill from the SDP workflow.

See `prompts/skills/tdd/SKILL.md` for complete documentation.

## Usage

```bash
/test [arguments]
```

## Implementation

The command delegates to the `@tdd` skill, which provides:

- Red-Green-Refactor cycle
- Quality gates
- Test-first discipline

## Related

- Skills: `prompts/skills/tdd/SKILL.md`
- Agents: `prompts/agents/builder.md`
