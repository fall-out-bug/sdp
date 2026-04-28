# Pi Review Gate Specification

Status: draft
Feature: F017
Owner surface: `sdp pi-review`

## Purpose

`sdp pi-review` is the SDP review gate for working-tree-first autonomous delivery. It gives Codex and other harnesses an external model opinion through `pi`, records a compact SDP verdict, and routes actionable findings into beads.

The gate is review-only. It may write review artifacts and beads. It must not patch code.

## User Surface

Primary command:

```bash
sdp pi-review --scope auto --base main --feature F017 --create-beads --write-verdict
```

Shell alias:

```bash
sdp-pi-review
```

Skill alias:

```text
sdp:pi-review
```

## Scope Model

`working-tree` is the canonical target. PR review is optional metadata, not the base abstraction.

| Scope | Behavior |
|---|---|
| `working-tree` | Review staged, unstaged, and untracked reviewable files. |
| `branch` | Review `git diff <base>...HEAD` plus touched files. |
| `auto` | Use PR metadata when available; otherwise fall back to working-tree review. |

The gate must print and persist the exact reviewed file list. Hidden dirty files are worse than noisy review.

## Context Packet

The context packet is the only input sent to model reviewers. It includes:

- git state: branch, base, head SHA, dirty status
- diff for the selected scope
- full content of touched reviewable files
- `AGENTS.md` / harness rules discovered for the repository
- linked feature, workstream, and beads summaries when available
- deterministic test evidence from the configured review command

The packet must include hashes for diff, touched files, rules, and test evidence. Raw prompt text is not stored in `.sdp/review_verdict.json`.

## Test Evidence

Models do not run tests independently in MVP. `sdp pi-review` runs one deterministic project review command before model review:

1. Read `.sdp/review.yaml` if present.
2. Fall back to `.sdp/config.yml` acceptance command if present.
3. Fall back to conservative ecosystem detection only for local drafts.

For delivery-loop use, missing or failing test evidence prevents `APPROVED`.

## Model Policy

MVP reviewers:

- GLM through Z.AI/OpenAI-compatible transport
- Kimi through Moonshot/OpenAI-compatible transport
- OpenRouter fallback when a primary model is unavailable

Model ids must be config aliases, not hard-coded literals. Provider offerings change; the gate stores resolved provider/model ids in telemetry for each run.

## Synthesis Rules

The synthesizer decides the final SDP verdict from model outputs and deterministic evidence.

- `APPROVED`: tests passed and no blocking findings remain.
- `CHANGES_REQUESTED`: at least one blocking finding remains.
- `ESCALATED`: the gate cannot decide safely because of model failure, repeated contradictory findings, missing required context, or repeated non-convergence.

SDP severity is `P0`-`P3`.

- `P0`/`P1` always block.
- `P2`/`P3` are tracked by default.
- A `P2`/`P3` may block only when the synthesizer gives a concrete blocking reason.

If GLM and Kimi conflict, the synthesizer may mark a finding false-positive only with evidence from the packet. A plausible `P0`/`P1` from any model blocks unless explicitly disproven.

## Artifacts

Compact verdict:

```text
.sdp/review_verdict.json
```

Run artifacts:

```text
.sdp/runs/pi-review/<run_id>/packet.json
.sdp/runs/pi-review/<run_id>/models/<model>.json
.sdp/runs/pi-review/<run_id>/run.json
```

`.sdp/review_verdict.json` is the canonical gate output. Raw model outputs are telemetry artifacts and should be gitignored by default.

## Beads Integration

With `--create-beads`, every accepted finding gets a beads issue.

Required fields in the issue body:

- `source = pi-review`
- feature id
- workstream id when known
- run id
- severity
- `blocking = true|false`
- file and line when available
- evidence excerpt or packet reference
- recommendation
- source model ids

Existing open findings should be updated or deduplicated by stable hash:

```text
rule + symbol_path + normalized_snippet + recommendation
```

## Delivery Loop Integration

The delivery loop treats `sdp pi-review` as its review gate:

```text
build/fix -> sdp pi-review -> beads findings -> fix -> repeat until clean
```

The loop may run until clean, but must have a runaway detector:

- wall-clock budget
- repeated same-finding count
- provider failure budget
- human escalation after repeated non-convergence

No automation may silently mark `pi-review` skipped.

## Non-Goals For MVP

- Full PR review UI.
- Inline GitHub review comments.
- Automatic code fixes by the review command.
- Full multi-provider council beyond GLM and Kimi.
- Raw transcript storage in `.sdp/review_verdict.json`.

