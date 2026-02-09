---
description: System design with progressive disclosure
agent: planner
---

# /design — Design

## Overview

This command implements the design skill from the SDP workflow.

See `/.claude/skills/design/SKILL.md` for complete documentation.

## Usage

```bash
/design [arguments]
```

## Implementation

The command delegates to the `design` skill, which provides:

- Systematic workflow
- Quality gates
- Proper error handling
- Documentation

## Related

- Skills: `.claude/skills/design/SKILL.md`
- Agents: `.claude/agents/planner.md`
