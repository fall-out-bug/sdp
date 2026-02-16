---
description: System design with progressive disclosure
agent: planner
---

# /design — Design

## Overview

This command implements the design skill from the SDP workflow.

See `/prompts/skills/design/SKILL.md` for complete documentation.

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

- Skills: `prompts/skills/design/SKILL.md`
- Agents: `prompts/agents/planner.md`
