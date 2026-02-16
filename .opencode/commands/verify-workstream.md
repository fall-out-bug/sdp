---
description: Validate workstream documentation against codebase reality.
agent: builder
---

# /verify-workstream — Verify-workstream

## Overview

This command implements the verify-workstream skill from the SDP workflow.

See `/prompts/skills/verify-workstream/SKILL.md` for complete documentation.

## Usage

```bash
/verify-workstream [arguments]
```

## Implementation

The command delegates to the `verify-workstream` skill, which provides:

- Systematic workflow
- Quality gates
- Proper error handling
- Documentation

## Related

- Skills: `prompts/skills/verify-workstream/SKILL.md`
- Agents: `prompts/agents/builder.md`
