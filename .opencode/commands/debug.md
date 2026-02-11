---
description: Systematic debugging using scientific method for evidence-based root cause analysis
agent: planner
---

# /debug — Debug

## Overview

This command implements the debug skill from the SDP workflow.

See `/prompts/skills/debug/SKILL.md` for complete documentation.

## Usage

```bash
/debug [arguments]
```

## Implementation

The command delegates to the `debug` skill, which provides:

- Systematic workflow
- Quality gates
- Proper error handling
- Documentation

## Related

- Skills: `prompts/skills/debug/SKILL.md`
- Agents: `prompts/agents/planner.md`
