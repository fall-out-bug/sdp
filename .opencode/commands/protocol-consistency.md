---
description: Audit consistency across workstream docs, CLI capabilities, and CI workflows.
agent: builder
---

# /protocol-consistency — Protocol-consistency

## Overview

This command implements the protocol-consistency skill from the SDP workflow.

See `/prompts/skills/protocol-consistency/SKILL.md` for complete documentation.

## Usage

```bash
/protocol-consistency [arguments]
```

## Implementation

The command delegates to the `protocol-consistency` skill, which provides:

- Systematic workflow
- Quality gates
- Proper error handling
- Documentation

## Related

- Skills: `prompts/skills/protocol-consistency/SKILL.md`
- Agents: `prompts/agents/builder.md`
