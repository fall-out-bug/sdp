---
description: Multi-agent quality review (QA + Security + DevOps + SRE + TechLead + Documentation + Contract Validation)
agent: builder
---

# /codereview — Review

## Overview

This command implements the review skill from the SDP workflow.

See `/prompts/skills/review/SKILL.md` for complete documentation.

## Usage

```bash
/codereview [arguments]
```

## Implementation

The command delegates to the `review` skill, which provides:

- Systematic workflow
- Quality gates
- Proper error handling
- Documentation

## Related

- Skills: `prompts/skills/review/SKILL.md`
- Agents: `prompts/agents/builder.md`
