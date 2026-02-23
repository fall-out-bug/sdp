---
description: Quick documentation vs code reality validation.
agent: builder
---

# /reality-check — Reality-check

## Overview

This command implements the reality-check skill from the SDP workflow.

See `/prompts/skills/reality-check/SKILL.md` for complete documentation.

## Usage

```bash
/reality-check [arguments]
```

## Implementation

The command delegates to the `reality-check` skill, which provides:

- Systematic workflow
- Quality gates
- Proper error handling
- Documentation

## Related

- Skills: `prompts/skills/reality-check/SKILL.md`
- Agents: `prompts/agents/builder.md`
