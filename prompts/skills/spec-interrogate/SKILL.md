---
name: spec-interrogate
description: Use when a spec or plan needs clean-context Socratic criticism before planning or implementation.
version: 1.1.0
tags:
  - discovery
  - spec
  - quality-gate
  - socratic
compatibility:
  - claude-code
  - opencode
  - cursor
  - codex
  - pi
changes:
  - "1.1.0: Add pi-backed Socratic critic/judge protocol, rubrics, provider rotation, and evidence contract"
  - "1.0.0: Initial public skill for spec hardening with auditable report and evidence contract"
---

# spec-interrogate

## Purpose

Pressure-test a text artifact before the next delivery step.
The interrogator receives only the artifact, selected rubrics, and invocation parameters. No chat history, beads context, author explanation, tools, or repository context.

This is not a gate. It produces evidence for the author and for downstream tooling. Policy decisions about whether work may pass a stage belong to `sdp-gate`.

## Use When

- Before turning a non-trivial spec into a plan or implementation task.
- Before using a SpecKit artifact as a contract for `sdp-trace`, `sdp-gate`, or another downstream component.
- When the author is too close to the document and may be relying on unstated context.

Do not use this for code review, implementation, or general research.

## Inputs

```bash
@spec-interrogate <artifact-path> \
  [--mode socratic|cold-read|adversarial|impl-test] \
  [--rubrics RUBRIC_LIST_OR_FILE] \
  [--questions N] \
  [--rounds M] \
  [--feature-id F] \
  [--evidence-path PATH] \
  [--report-path PATH]
```

Defaults:

- `--mode socratic`
- `--questions 12`
- `--rounds 3`
- `--rubrics default`
- `--evidence-path .sdp/evidence/spec-interrogate.json`
- `--report-path .sdp/reports/spec-interrogate.md`

## Rubrics

Use the full selected rubric set in one critic pass. Do not split one round into one subagent per rubric unless the artifact is too large for a single model context.

Default rubrics:

- problem and goal
- system boundary and non-goals
- roles and actors
- primary scenarios
- assumptions and dependencies
- edge cases and failure behavior
- security and access
- observability and metrics
- testability and acceptance
- rollout, migration, and backward compatibility
- open questions and risks

Tailor the rubric set only by removing irrelevant categories or adding domain-critical ones. Record the final rubric set in evidence.

## Roles

**Author** owns the artifact and edits it.

**Critic** is a fresh model invocation. It asks Socratic questions grouped by rubric. The critic MUST NOT propose solutions, write patches, choose product behavior, or rewrite the artifact.

**Judge** is a different fresh model invocation. It compares the original artifact, revised artifact, critic questions, and author resolution notes. The judge does not rewrite the spec and does not suggest new solutions.

## Provider Policy

Use `pi` with clean context:

```bash
pi --provider <provider> --model <model> --no-tools --no-context-files --no-session -p "<prompt>"
```

provider rotation is required:

- critic providers rotate between rounds: `zai/glm-5.1`, `kimi-coding/k2p6`, `minimax/MiniMax-M2.7`
- the next critic round must not reuse the previous critic provider
- the judge provider must differ from the current critic provider
- if a provider fails, record the failure and use the next provider; do not hide degraded coverage

## Protocol (mode `socratic`)

1. Save the original artifact snapshot or hash.
2. Select rubrics and critic provider.
3. Invoke the critic with only the artifact, rubrics, and output schema.
4. Require critic output as questions only.
5. Author edits the artifact. Chat answers do not count unless reflected in the artifact.
6. Author writes resolution notes per question: `resolved`, `rejected`, or `deferred`, with rationale.
7. Invoke the judge with original artifact, revised artifact, critic questions, and resolution notes.
8. If blocking contradictions remain, repeat with a different critic provider.
9. Stop when exit criteria are met or `--rounds` is exhausted.

## Critic Output

The critic output must be a JSON object:

```json
{
  "role": "critic",
  "critic_provider": "zai/glm-5.1",
  "rubrics": ["problem and goal"],
  "questions": [
    {
      "id": "Q1",
      "rubric": "problem and goal",
      "severity": "blocking",
      "artifact_ref": "section heading or line reference",
      "question": "What observable CTO question does this spec answer?",
      "why_it_matters": "Without this, acceptance cannot distinguish telemetry from governance.",
      "cannot_verify_until_answered": "Whether the component answers degradation over time."
    }
  ]
}
```

Allowed severities: `blocking`, `major`, `minor`.

Reject critic output that contains fixes, patches, rewritten text, implementation plans, policy verdicts, or "you should..." recommendations. Re-run with a stricter prompt if needed.

## Judge Output

The judge output must be a JSON object:

```json
{
  "role": "judge",
  "judge_provider": "kimi-coding/k2p6",
  "critic_provider": "zai/glm-5.1",
  "items": [
    {
      "question_id": "Q1",
      "status": "resolved",
      "evidence_ref": "revised artifact section heading or line reference",
      "assessment": "The revised artifact states the CTO-facing degradation question and acceptance evidence."
    }
  ],
  "new_contradictions": [],
  "scope_creep": [],
  "verdict": "PASS"
}
```

Allowed item statuses: `resolved`, `partially_resolved`, `unresolved`, `new_contradiction`, `scope_creep`.

Allowed verdicts:

- `PASS`: no blocking unresolved items and no new contradictions
- `REWORK`: blocking or major unresolved items remain
- `ABORT`: author stops the process or model coverage is too degraded to trust

## Evidence Contract

Every run writes a report and evidence JSON. Stdout is not the system of record.

```json
{
  "interrogate_verdict": "REWORK",
  "artifact_path": "specs/001-sdp-trace-time-series-evidence-substrate/spec.md",
  "feature_id": "F163",
  "mode": "socratic",
  "rounds_completed": 1,
  "max_rounds": 3,
  "rubrics": ["problem and goal", "system boundary and non-goals"],
  "critic_provider": "zai/glm-5.1",
  "judge_provider": "kimi-coding/k2p6",
  "provider_failures": [],
  "open_questions_count": 2,
  "open_questions": [],
  "new_contradictions": [],
  "scope_creep": [],
  "report_path": ".sdp/reports/spec-interrogate.md",
  "report_summary": "Two major questions remain unresolved."
}
```

The report must include artifact path, rubrics, providers, critic questions, author resolution notes, judge conclusion, explicit verdict, and next action.

## Exit Criteria

Stop with `PASS` only when:

- no `blocking` question is unresolved
- no `major` question is unresolved unless it is explicitly deferred with owner and rationale
- the judge reports no new contradictions
- scope creep is absent or intentionally accepted by the author

Stop with `REWORK` when the round cap is hit with unresolved blocking or major issues.

## Other Modes

`cold-read`, `adversarial`, and `impl-test` are single-pass variants. They still use clean-context `pi`, rubrics when relevant, evidence JSON, and the same not a gate boundary.

## Anti-Patterns

- Passing chat history or repository context to critic or judge.
- Letting the critic propose fixes.
- Treating `PASS` as permission to merge or deploy.
- Hiding provider failures.
- Counting a chat answer as resolved when the artifact did not change.

## Acceptance Boundary

This skill works only with text artifacts: specs, plans, design docs, schemas, and SpecKit files. If the target is executable code, use review tooling instead.
