# SDP Prompt Surface Reference

This document maps the prompt surfaces that ship with SDP today. For exact prompt behavior, inspect the source files in `prompts/`.

## Source of Truth

Canonical prompt content lives here:

- `prompts/commands/`
- `prompts/skills/`
- `prompts/agents/`

Supported adapters expose that content through:

- `.claude/`
- `.cursor/`
- `.opencode/`
- `.codex/`

Do not treat adapter directories as the source of truth. Edit `prompts/`.

## Common Prompt Surfaces

| Surface | Backing source | Current role |
|---------|----------------|--------------|
| `/feature` | `prompts/commands/feature.md`, `prompts/skills/feature/SKILL.md` | Planning entry point |
| `/idea` | `prompts/commands/idea.md`, `prompts/skills/idea/SKILL.md` | Requirements capture |
| `/design` | `prompts/commands/design.md`, `prompts/skills/design/SKILL.md` | Workstream planning |
| `/build` | `prompts/commands/build.md`, `prompts/skills/build/SKILL.md` | Single leaf-workstream execution |
| `/review` | `prompts/commands/review.md`, `prompts/skills/review/SKILL.md` | Review and verdict loop |
| `/oneshot` | `prompts/commands/oneshot.md`, `prompts/skills/oneshot/SKILL.md` | Outer-loop feature execution |
| `/deploy` | `prompts/commands/deploy.md`, `prompts/skills/deploy/SKILL.md` | Prompt-level release handoff surface |
| `/beads` | `prompts/commands/beads.md`, `prompts/skills/beads/SKILL.md` | Beads task-tracker integration |
| `/strataudit` | `prompts/commands/strataudit.md`, `prompts/skills/strataudit/SKILL.md` | Evidence-backed strategy traceability audit |
| `/debug`, `/hotfix`, `/bugfix`, `/issue` | matching files under `prompts/commands/` and `prompts/skills/` | Investigation and fix flows |

## Current Operating Reality

Prompt surfaces and CLI commands are related, but they are not identical.

- `sdp init`, `sdp doctor`, `sdp plan`, `sdp apply`, `sdp status`, and `sdp next` are the clearest public Local Mode surfaces.
- Prompt surfaces are the advanced harness-native layer.
- Queue-backed usage expects Beads.

Important distinction:

- `/deploy` is a prompt surface in the prompt bundle.
- `sdp deploy` is a CLI command with narrower semantics.
- Today `sdp deploy` records an approval event after merge; it does not merge branches or deploy infrastructure.

## Where to Look Next

- [../CLI_REFERENCE.md](../CLI_REFERENCE.md) for current CLI surfaces
- [../PROTOCOL.md](../PROTOCOL.md) for the current protocol overview
- [../PRODUCT_CONTRACT.md](../PRODUCT_CONTRACT.md) for product definition and mode policy
- [../../prompts/README.md](../../prompts/README.md) for prompt layout
