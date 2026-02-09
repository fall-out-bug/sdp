---
description: Contract test generation and validation workflow.
agent: builder
---

# /test — Test

## Overview

This command implements the test skill from the SDP workflow.

See `/.claude/skills/test/SKILL.md` for complete documentation.

## Usage

```bash
/test [arguments]
```

## Implementation

The command delegates to the `test` skill, which provides:

- Systematic workflow
- Quality gates
- Proper error handling
- Documentation

## Related

- Skills: `.claude/skills/test/SKILL.md`
- Agents: `.claude/agents/builder.md`
