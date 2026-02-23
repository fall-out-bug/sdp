---
description: Beads task tracker integration for SDP workflows.
agent: builder
---

# /beads — Beads

## Overview

This command implements the beads skill from the SDP workflow.

See `/prompts/skills/beads/SKILL.md` for complete documentation.

## Usage

```bash
/beads [arguments]
```

## Implementation

The command delegates to the `beads` skill, which provides:

- Systematic workflow
- Quality gates
- Proper error handling
- Documentation

## Related

- Skills: `prompts/skills/beads/SKILL.md`
- Agents: `prompts/agents/builder.md`
