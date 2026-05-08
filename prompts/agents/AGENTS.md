# prompts/agents — Agent Contract

## Scope

This subtree owns role prompts for SDP agent personas and reviewers.

## Contract

Every top-level agent must state:

- what intent or workflow it supports
- what SDP stage or entity it updates
- what artifact, verdict, or handoff it must emit

## Dependencies

Agent prompts may reference root `AGENTS.md`, `docs/reference/agent-catalog.md`,
and relevant skill files.

Do not copy full skill workflows into agent prompts. Agents choose and execute
skills; skills own the procedural steps.

## Runtime Assumptions

Harnesses may load these prompts through symlinks, generated adapters, or explicit
dispatch prompts. Keep role behavior harness-neutral unless the file is explicitly
for one harness.

## Local Rules

- Prefer fewer durable roles over many overlapping personas.
- If a role does not own a unique transition, collapse it into an existing agent
  or a review dimension.
- Review agents must emit findings with severity and evidence.
- Execution agents must emit concrete changed files, gates run, and blockers.
