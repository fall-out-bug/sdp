---
description: Codebase analysis and architecture validation - what's actually there vs documented
agent: builder
---

# /reality — Reality

## Overview

This command implements the reality skill from the SDP workflow.

See `/prompts/skills/reality/SKILL.md` for complete documentation.

## Usage

```bash
/reality [arguments]
```

## Implementation

The command delegates to the `reality` skill, which provides:

- Systematic workflow
- Quality gates
- Proper error handling
- Documentation

## Related

- Skills: `prompts/skills/reality/SKILL.md`
- Agents: `prompts/agents/builder.md`
