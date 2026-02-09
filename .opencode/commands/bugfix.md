---
description: Quality bug fixes (P1/P2). Full TDD cycle, branch from feature/develop, no production deploy.
agent: builder
---

# /bugfix — Bugfix

## Overview

This command implements the bugfix skill from the SDP workflow.

See `/.claude/skills/bugfix/SKILL.md` for complete documentation.

## Usage

```bash
/bugfix [arguments]
```

## Implementation

The command delegates to the `bugfix` skill, which provides:

- Systematic workflow
- Quality gates
- Proper error handling
- Documentation

## Related

- Skills: `.claude/skills/bugfix/SKILL.md`
- Agents: `.claude/agents/builder.md`
