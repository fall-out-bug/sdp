---
description: Deployment orchestration. Generates artifacts and EXECUTES GitFlow merge.
agent: deployer
---

# /deploy — Deploy

## Overview

This command implements the deploy skill from the SDP workflow.

See `/.claude/skills/deploy/SKILL.md` for complete documentation.

## Usage

```bash
/deploy [arguments]
```

## Implementation

The command delegates to the `deploy` skill, which provides:

- Systematic workflow
- Quality gates
- Proper error handling
- Documentation

## Related

- Skills: `.claude/skills/deploy/SKILL.md`
- Agents: `.claude/agents/deployer.md`
