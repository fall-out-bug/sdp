---
description: Evidence-backed strategy traceability audit over a document corpus; use when the user needs document-grounded alignment analysis across strategy, architecture, design, or implementation materials. Prefer an injected host-native runtime when available, otherwise use a configured OpenAI-compatible runtime; OpenRouter is the default network accelerator, not the only path.
agent: architect
---

# /strataudit — StratAudit

## Overview

This command implements the `strataudit` skill from the SDP workflow.

See `/prompts/skills/strataudit/SKILL.md` for complete documentation.

## Usage

```bash
/strataudit [arguments]
```

## Implementation

The command delegates to the `strataudit` skill, which provides:

- mode-based document-backed audit flow
- explicit runtime selection order
- structured artifact output
- trust-oriented failure and refusal behavior

## Related

- Skills: `prompts/skills/strataudit/SKILL.md`
- Reference: `docs/reference/strataudit-evidence-policy.md`
- Reference: `docs/reference/strataudit-runtime-policy.md`
- Reference: `docs/reference/strataudit-output-modes.md`
- Agents: `prompts/agents/architect.md`
