---
name: strataudit
description: Evidence-backed strategy traceability audit over a document corpus; use when the user needs document-grounded alignment analysis across strategy, architecture, design, or implementation materials. Prefer an injected host-native runtime when available, otherwise use a configured OpenAI-compatible runtime; OpenRouter is the default network accelerator, not the only path.
version: 1.0.0
---

# @strataudit - Strategy Traceability Audit

Run a document-backed audit over a real corpus and produce traceability artifacts instead of free-form prose.

## Use When

- the user needs alignment analysis across strategic and delivery documents
- the answer must be backed by extracted entities, traces, and findings
- the corpus is large enough that a systematic audit beats a manual summary

## Do Not Use When

- the user wants only a short narrative summary
- there is no accessible corpus yet
- the problem is operational debugging rather than document traceability

## Runtime Order

1. use an injected host-native runtime from the harness when available
2. otherwise use a configured OpenAI-compatible runtime
3. use OpenRouter as the default network path when no better native option exists
4. use `sdp-strataudit run` only as CLI fallback

The CLI can resolve configured network runtimes. It cannot create a host-native runtime on its own.

## Workflow

1. locate or create `strataudit.yaml`
2. confirm source directories and level patterns
3. resolve runtime in the order above
4. run ingest → extract → link → analyze → report
5. return artifact paths plus the main trust caveats

## Output

- `.strataudit/report.json`
- `.strataudit/report.html`
- `.strataudit/similarity_distribution.json`
- `.strataudit/strataudit.db`

## Rules

- only rely on real document text
- do not fabricate quotes, traces, or initiatives
- preserve source language unless a later presentation layer explicitly derives display fields
- if runtime capabilities are missing, fail explicitly instead of silently degrading the audit

## References

- `docs/QUICKSTART.md`
- `docs/reference/skills.md`
- `sdp-strataudit run`
