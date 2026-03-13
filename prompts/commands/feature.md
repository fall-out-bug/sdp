---
description: Feature planning orchestrator (idea → design → workstreams)
agent: planner
---

# /feature — Feature

## Overview

This command implements the feature skill from the SDP workflow.

See `/prompts/skills/feature/SKILL.md` for complete documentation.

## Usage

```bash
/feature [arguments]
```

## Implementation

The command delegates to the `feature` skill, which provides:

- Full-description handoff with no truncation or pre-summary
- Multi-service and multi-language topology discovery before decomposition
- Mandatory clarification for admin/user-facing behavior when requirements are underspecified
- Systematic workflow
- Quality gates
- Proper error handling
- Documentation

## Related

- Skills: `prompts/skills/feature/SKILL.md`
- Agents: `prompts/agents/builder.md`
