---
description: Emergency P0 fixes. Fast-track production deployment with minimal changes. Branch from main, immediate deploy.
agent: builder
---

# /hotfix — Hotfix

## Overview

This command implements the hotfix skill from the SDP workflow.

See `/.claude/skills/hotfix/SKILL.md` for complete documentation.

## Usage

```bash
/hotfix [arguments]
```

## Implementation

The command delegates to the `hotfix` skill, which provides:

- Systematic workflow
- Quality gates
- Proper error handling
- Documentation

## Related

- Skills: `.claude/skills/hotfix/SKILL.md`
- Agents: `.claude/agents/builder.md`
