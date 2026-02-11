---
description: Autonomous multi-agent execution with checkpoints, resume, and PR-less modes
agent: orchestrator
---

# /oneshot — Oneshot

## Overview

This command implements the oneshot skill from the SDP workflow.

See `/prompts/skills/oneshot/SKILL.md` for complete documentation.

## Usage

```bash
/oneshot [arguments]
```

## Implementation

The command delegates to the `oneshot` skill, which provides:

- Systematic workflow
- Quality gates
- Proper error handling
- Documentation

## Related

- Skills: `prompts/skills/oneshot/SKILL.md`
- Agents: `prompts/agents/orchestrator.md`
