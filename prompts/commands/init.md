---
description: Initialize SDP in current project (interactive wizard)
agent: builder
---

# /init — Init

## Overview

This command implements the init skill from the SDP workflow.

See `/prompts/skills/init/SKILL.md` for complete documentation.

## Usage

```bash
/init [arguments]
```

## Implementation

The command delegates to the `init` skill, which provides:

- Systematic workflow
- Quality gates
- Proper error handling
- Documentation

## Related

- Skills: `prompts/skills/init/SKILL.md`
- Agents: `prompts/agents/builder.md`
