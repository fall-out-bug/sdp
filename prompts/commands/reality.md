---
description: Reconstruct a single-repo reality baseline from local evidence and emit open reality artifacts
agent: builder
---

# /reality - Reality

## Overview

This command runs the OSS `@reality` baseline. It is local, evidence-first, and single-repo scoped.

## Runtime

```bash
sdp reality emit-oss [--quick|--deep|--bootstrap-sdp] [--focus=architecture|quality|testing|docs|security]
```

## Behavior

- Scans code, tests, configs, manifests, and in-repo docs.
- Emits `.sdp/reality/*.json` and `docs/reality/*.md`.
- Runs a heuristic cross-check pass inside the same repository.
- Does not spawn subagents or claim `reality-pro` behavior.

## Related

- Skill: `prompts/skills/reality/SKILL.md`
- Spec: `docs/specs/reality/OSS-SPEC.md`
