---
name: spec-interrogate
description: Use when a spec or plan may hide ambiguities and needs a context-stripped challenge before planning or implementation.
version: 1.0.0
changes:
  - "1.0.0: Initial public skill for spec hardening with auditable report and evidence contract"
---

# spec-interrogate

> Challenge a text artifact with a fresh interrogator before you commit to planning or implementation.

The interrogator receives only the artifact and invocation parameters. No chat history. No implied context. No author explanation. The goal is simple: surface what an implementer still cannot infer from the document itself.

---

## Use When

- Before `sdp phase plan` for non-trivial Discovery output
- Before committing to a risky architecture, API, or rollout plan
- When the author is too close to the document and may be missing undefined terms, missing error paths, or scope leaks

Do not use this for code review or implementation. Use the relevant review/build skill instead.

---

## Inputs

```bash
@spec-interrogate <artifact-path> \
  [--mode socratic|cold-read|adversarial|impl-test] \
  [--questions N] \
  [--rounds M] \
  [--feature-id F] \
  [--evidence-path PATH] \
  [--report-path PATH]
```

Defaults:

- `--mode socratic`
- `--questions 5`
- `--rounds 5`
- `--evidence-path .sdp/evidence/spec-interrogate.json`
- `--report-path .sdp/reports/spec-interrogate.md`

---

## Shared Output Contract

Every run must create:

1. A human-readable report at `report-path`
2. A machine-readable evidence file at `evidence-path`

The report is mandatory even on `PASS`. It must contain:

- artifact path
- selected mode
- short summary of what was tested
- ordered unresolved questions or gaps
- explicit verdict
- next action

The evidence file must include:

```json
{
  "interrogate_verdict": "PASS | REWORK | ABORT",
  "artifact_path": "docs/discovery/my-feature/validation.md",
  "feature_id": "F042",
  "mode": "socratic",
  "rounds_completed": 2,
  "max_rounds": 5,
  "open_questions_count": 0,
  "open_questions": [],
  "report_path": ".sdp/reports/spec-interrogate.md",
  "report_summary": "No unresolved implementation-blocking questions remain."
}
```

Each unresolved question must be structured:

```json
{
  "id": "Q1",
  "type": "scope-ambiguity",
  "impact": "plan-blocking",
  "question": "What is the fallback behavior when the upstream model call times out?"
}
```

Stdout is not the system of record. The report and evidence files are.

---

## Modes

### `socratic`

Iterative dialogue via document edits.

1. Interrogator asks up to `N` high-impact questions.
2. Author edits the artifact instead of replying in chat.
3. Interrogator re-reads and checks whether the questions are resolved.
4. Repeat until convergence or `--rounds` is exhausted.

Verdict:

- `PASS` when no unresolved questions remain
- `REWORK` when the round cap is hit and blocking questions remain
- `ABORT` when the author explicitly stops

### `cold-read`

Cheap first pass.

1. Interrogator reads once.
2. It writes:
   - what it believes the artifact says
   - what it still cannot infer
   - what it refuses to assume
3. Unresolved inferences become `open_questions[]`.

Verdict:

- `PASS` when `open_questions_count = 0`
- `REWORK` otherwise
- `ABORT` when the author explicitly stops

Accounting: `rounds_completed = 1`, `max_rounds = 1`

### `adversarial`

Artifact-level attack review.

1. Interrogator reads once.
2. It lists trust-boundary gaps, abuse paths, failure modes, and mitigation holes.
3. Blocking gaps become `open_questions[]`.

Verdict:

- `PASS` when no blocking gaps remain
- `REWORK` otherwise
- `ABORT` when the author explicitly stops

Accounting: `rounds_completed = 1`, `max_rounds = 1`

### `impl-test`

Checks whether another agent could implement the artifact without hallucinating.

1. Interrogator tries to outline a minimal implementation plan from the artifact alone.
2. Any step that requires invented assumptions becomes `open_questions[]`.
3. The report must separate grounded steps from assumption-dependent steps.

Verdict:

- `PASS` when the outline requires zero invented assumptions
- `REWORK` otherwise
- `ABORT` when the author explicitly stops

Accounting: `rounds_completed = 1`, `max_rounds = 1`

---

## Question Taxonomy

Prioritize only questions that matter for planning or implementation:

1. `why`
2. `undefined-term`
3. `missing-error-path`
4. `scope-ambiguity`
5. `unstated-assumption`

Do not waste rounds on style or formatting unless they obscure meaning.

---

## SDP Integration

This is agent discipline before the Plan gate, not CLI enforcement.

```bash
@spec-interrogate docs/discovery/<slug>/validation.md --feature-id <F>

# only after PASS:
sdp phase plan --feature-id <F> --strict --evidence-path .sdp/evidence/plan.json
```

If the verdict is `REWORK`, do not call `sdp phase plan`. Resume work using the unresolved questions in the generated report.

---

## Skip Rules

Skip only when:

- the task is trivial
- there is no Discovery artifact to interrogate
- `--skip-interrogate` is explicitly documented in beads with a reason

"Probably fine" is not a valid skip reason.

---

## Examples

```bash
# iterative hardening
@spec-interrogate docs/discovery/my-feature/validation.md --feature-id F042

# cheap sanity check
@spec-interrogate docs/plans/arch-decision.md --mode cold-read

# risky architecture review
@spec-interrogate docs/plans/auth-redesign.md --mode adversarial
```

---

## Acceptance Boundary

This skill is for text artifacts only: specs, plans, design docs, and schema-oriented documents.
If the target is code, use review tooling instead.
