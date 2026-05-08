# prompts/skills — Agent Contract

## Scope

This subtree owns structured `SKILL.md` workflow prompts published as SDP protocol
artifacts.

## Contract

Each skill defines an executable workflow: triggers, preconditions, steps, outputs,
stop conditions, and recovery behavior.

Skills must keep YAML frontmatter valid and aligned with `sdp.manifest.yaml` while
the manifest remains the inventory gate.

## Dependencies

Skills may reference root `AGENTS.md`, `docs/reference/skills.md`,
`docs/reference/skill-authoring.md`, and task-specific reference docs.

Do not encode package API contracts here. Put subtree-specific facts in the
nearest module-local `AGENTS.md`.

## Runtime Assumptions

Harnesses may load these files through Claude plugin format, symlinks, generated
adapters, or explicit prompt paths. Keep prose harness-neutral unless a section is
explicitly marked as harness-specific.

## Local Rules

- Keep workflows executable and bounded.
- State refusal, blocker, and stop conditions explicitly.
- State completion criteria explicitly: verification run, scoped staging, commit,
  push or exact blocker.
- Never instruct agents to use `git add .`; enumerate scoped files from the write
  plan or command output.
- Do not duplicate Beads, branch, or publish policy unless the skill changes how
  that policy is applied.
- Deprecated skills must route to the canonical intent or replacement skill.
