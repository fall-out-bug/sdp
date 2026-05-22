---
name: reality
description: Single-repo reality scan that reconstructs code, docs, drift, integrations, and SDP readiness from local evidence
version: 2.1.0
changes:
  - Aligned public contract with OSS runtime
  - Added explicit mode surface for quick, deep, focus, and bootstrap-sdp
  - Removed fake multi-agent Task fan-out claims from OSS
---

# @reality - Single-Repo Reality Baseline

Use `@reality` to recover what is actually present in one repository and emit the open reality artifact set.

OSS `@reality` is evidence-first and local-only. It does not pretend to be `reality-pro`.

## What OSS Does

1. Scans one repository working tree.
2. Builds a baseline from code, tests, configs, manifests, and in-repo docs.
3. Emits machine-readable artifacts in `.sdp/reality/`.
4. Emits human-readable summaries in `docs/reality/`.
5. Runs a heuristic cross-check pass across code, tests, configs, manifests, and docs inside the same run.

## What OSS Does Not Do

- No Task-tool fan-out.
- No "8 expert agents" promise.
- No multi-repo orchestration.
- No persistent repository memory.
- No consulting-grade adversarial synthesis.

Those belong in `reality-pro`, not here.

## Runtime Surface

Primary command:

```bash
sdp reality emit-oss [--quick|--deep|--bootstrap-sdp] [--focus=architecture|quality|testing|docs|security]
sdp reality validate
```

Mode rules:

- `--deep` is the default OSS baseline.
- `--quick` emits the same artifact families with reduced evidence detail.
- `--bootstrap-sdp` keeps the same baseline scan but emphasizes first-workstream recommendations and agent-readiness notes.
- `--focus=...` adds reporting emphasis for one domain without changing the open artifact contract.

## Output Contract

Human-readable outputs:

- `docs/reality/summary.md`
- `docs/reality/architecture.md`
- `docs/reality/quality.md`
- `docs/reality/bootstrap.md`

Machine-readable outputs:

- `.sdp/reality/reality-summary.json`
- `.sdp/reality/feature-inventory.json`
- `.sdp/reality/architecture-map.json`
- `.sdp/reality/integration-map.json`
- `.sdp/reality/quality-report.json`
- `.sdp/reality/drift-report.json`
- `.sdp/reality/readiness-report.json`

## OSS Review Semantics

OSS still avoids single-pass overconfidence, but it does so with a local cross-check strategy:

1. Primary source-first pass over code and executable config.
2. Secondary heuristic review against tests, manifests, and docs.
3. Synthesis that downgrades weak documentation claims into drift or unresolved questions.

This is limited review, not multi-agent arbitration.

## Mode Guidance

| Mode | Use When | Result |
|---|---|---|
| `--quick` | Need a fast baseline before planning | Same artifact set, less detail in findings |
| `--deep` | Need the normal OSS baseline | Full local single-repo scan |
| `--focus=architecture` | Need boundaries and entrypoints | Architecture report gets extra emphasis |
| `--focus=quality` | Need hotspots and maintainability | Quality findings get extra emphasis |
| `--focus=testing` | Need verification posture | Readiness emphasizes missing tests |
| `--focus=docs` | Need doc drift review | Drift and trust signals get extra emphasis |
| `--focus=security` | Need boundary review | Integration surfaces and failure notes are highlighted |
| `--bootstrap-sdp` | Need first SDP-safe slices | Bootstrap report emits starter workstreams |

## Examples

```bash
sdp reality emit-oss
sdp reality emit-oss --quick
sdp reality emit-oss --focus=docs
sdp reality emit-oss --bootstrap-sdp --focus=architecture
```

## Related

- `docs/specs/reality/OSS-SPEC.md`
- `docs/specs/reality/ARTIFACT-CONTRACT.md`
- `docs/reference/reality-oss.md`
- `prompts/commands/reality.md`
