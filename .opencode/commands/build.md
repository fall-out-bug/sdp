---
description: Execute workstream with TDD and guard enforcement
agent: builder
---

# /build — Build

## Overview

This command implements the build skill from the SDP workflow.

See `/.claude/skills/build/SKILL.md` for complete documentation.

## Usage

```bash
/build [arguments]
```

## Implementation

The command delegates to the `build` skill, which provides:

- Systematic workflow
- Quality gates
- Proper error handling
- Documentation

## Related

- Skills: `.claude/skills/build/SKILL.md`
- Agents: `.claude/agents/builder.md`
