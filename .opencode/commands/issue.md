---
description: Analyze bugs, classify severity (P0-P3), route to appropriate fix command (/hotfix, /bugfix, or backlog).
agent: planner
---

# /issue — Issue

## Overview

This command implements the issue skill from the SDP workflow.

See `/prompts/skills/issue/SKILL.md` for complete documentation.

## Usage

```bash
/issue [arguments]
```

## Implementation

The command delegates to the `issue` skill, which provides:

- Systematic workflow
- Quality gates
- Proper error handling
- Documentation

## Related

- Skills: `prompts/skills/issue/SKILL.md`
- Agents: `prompts/agents/planner.md`
