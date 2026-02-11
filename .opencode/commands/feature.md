---
description: Feature planning orchestrator (idea → design → workstreams)
agent: builder
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

- Systematic workflow
- Quality gates
- Proper error handling
- Documentation

## Related

- Skills: `prompts/skills/feature/SKILL.md`
- Agents: `prompts/agents/builder.md`
