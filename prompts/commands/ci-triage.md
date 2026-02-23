---
description: Investigate failing GitHub Actions runs and produce root-cause plus Beads follow-up.
agent: builder
---

# /ci-triage — Ci-triage

## Overview

This command implements the ci-triage skill from the SDP workflow.

See `/prompts/skills/ci-triage/SKILL.md` for complete documentation.

## Usage

```bash
/ci-triage [arguments]
```

## Implementation

The command delegates to the `ci-triage` skill, which provides:

- Systematic workflow
- Quality gates
- Proper error handling
- Documentation

## Related

- Skills: `prompts/skills/ci-triage/SKILL.md`
- Agents: `prompts/agents/builder.md`
