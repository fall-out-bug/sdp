---
description: Deployment orchestration. Generates artifacts and EXECUTES GitFlow merge.
agent: deployer
---

# /deploy — Deploy

## Overview

This command implements the deploy skill from the SDP workflow.

See `/prompts/skills/deploy/SKILL.md` for complete documentation.

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

- Skills: `prompts/skills/deploy/SKILL.md`
- Agents: `prompts/agents/deployer.md`
