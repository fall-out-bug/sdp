---
name: strataudit
description: Use when the user needs evidence-backed alignment analysis, traceability gaps, coverage, or a source-grounded evidence pack across strategy, architecture, design, or implementation documents in a real corpus.
version: 1.1.0
---

# @strataudit - Strategy Traceability Audit

Run a document-backed strategy audit over a real corpus or existing `.strataudit`
artifacts. This skill is for evidence-backed audit work, not free-form strategy prose.

## Use When

- the user needs alignment analysis across strategic and delivery documents
- the answer must be backed by extracted entities, traces, findings, or saved artifacts
- the user needs one of these modes: `corpus-audit`, `traceability-audit`, `coverage-audit`, `evidence-pack`, `report-redraft`

## Do Not Use When

- the user wants only a short narrative summary
- there is no accessible corpus and no existing artifacts
- the task is brainstorming or roadmap generation from scratch
- the problem is operational debugging rather than document traceability

## Safety Guards

- only rely on real document text or saved audit artifacts
- do not fabricate quotes, traces, or initiatives
- preserve source language in the evidence layer unless explicitly asked to derive display text
- similarity alone is never enough to call a trace verified
- if provenance is weak, downgrade the claim or refuse to make it

## Audit Modes

| Mode | Use when | Must emit |
|------|----------|-----------|
| `corpus-audit` | corpus quality and ingest readiness are unclear | corpus inventory, exclusions, level coverage, caveats |
| `traceability-audit` | the user wants cross-level alignment and missing links | traces, findings, caveats |
| `coverage-audit` | the user asks what is and is not covered | coverage summary with explicit denominators |
| `evidence-pack` | the user wants inspectable proof behind claims | source-backed findings, trace tables, caveats |
| `report-redraft` | the user wants a better report from existing artifacts | rewritten sections with unchanged trust boundaries |

Start with `corpus-audit` when corpus quality is unknown. Use `report-redraft`
only when an evidence pack or prior audit artifacts already exist.

## Runtime Order

1. use an injected host-native runtime from the harness when available
2. otherwise use a configured OpenAI-compatible runtime
3. use OpenRouter as the default network enhancer/fallback
4. use artifact-only mode when the question can be answered from existing outputs
5. use `sdp-strataudit run` only as CLI fallback

The CLI can resolve configured network runtimes. It cannot create a host-native runtime on its own.

## Workflow

1. choose the audit mode that matches the user's question
2. validate inputs and trust boundary
3. resolve runtime by policy or reuse existing artifacts
4. run ingest → extract → link → analyze → report, or inspect saved artifacts
5. return artifact paths plus trust caveats and what is not claimed

## Refuse When

- the user asks for verified alignment without inspectable provenance
- the corpus root is missing and there are no prior artifacts
- the requested mode needs runtime capabilities that are unavailable
- the requested summary is broader than the evidence pack supports

## Output

- `.strataudit/report.json`
- `.strataudit/report.html`
- `.strataudit/similarity_distribution.json`
- `.strataudit/strataudit.db`
- explicit runtime choice or artifact-only path
- key trust caveats and what is not claimed

## References

- `docs/QUICKSTART.md`
- `docs/reference/skills.md`
- `docs/reference/strataudit-evidence-policy.md`
- `docs/reference/strataudit-runtime-policy.md`
- `docs/reference/strataudit-output-modes.md`
- `sdp-strataudit run`

## Recovery

| Symptom | Fix |
|---------|-----|
| Skill produces no output | Check working directory is project root with `docs/workstreams/backlog/` |
| "checkpoint not found" | Run `sdp-orchestrate --feature <ID>` to create initial checkpoint |
| "workstream files missing" | Run `sdp-orchestrate --index` to verify, then `@feature` to regenerate |
| Skill hangs / no progress | Check `.sdp/log/events.jsonl` for last event; use `sdp reset --feature <ID>` if stuck |
| Review loop exceeds 3 rounds | Use `@review --override "reason"`, `@review --partial`, or `@review --escalate` |
