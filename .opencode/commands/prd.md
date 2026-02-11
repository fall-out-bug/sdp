---
description: PRD generation and maintenance workflow.
agent: builder
---

# /prd — Prd

## Overview

This command implements the prd skill from the SDP workflow.

See `/prompts/skills/prd/SKILL.md` for complete documentation.

## Usage

```bash
/prd [arguments]
```

## Implementation

The command delegates to the `prd` skill, which provides:

- Systematic workflow
- Quality gates
- Proper error handling
- Documentation

## Related

- Skills: `prompts/skills/prd/SKILL.md`
- Agents: `prompts/agents/builder.md`
