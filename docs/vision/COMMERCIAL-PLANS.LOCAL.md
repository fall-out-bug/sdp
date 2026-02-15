# SDP Commercial Plans (Local Branch Only)

> Scope: non-public planning notes.
> Branch: `codex/vision`.
> This file is intentionally separated from public OSS vision docs.

## 1) Packaging and Licensing Strategy

- `L0` Protocol assets (Claude plugin + prompts/guides): MIT, public forever.
- `L1` Safety layer (hooks, guard, traces, provenance): MIT, distributed via Homebrew.
- `L2` Core orchestration tooling: MIT, distributed via Homebrew.
- `L3` Enterprise extensions: private repositories, separate commercial terms, no hard dependency from `L0-L2`.

## 2) Commercial Positioning Hypotheses

- Wedge: lightweight CI entrypoint (GitHub Action) to seed adoption.
- Buyer persona for commercial track: compliance/governance stakeholders in regulated environments.
- Trust argument: audit-ready evidence chain + provenance + explainability.
- Maturity requirements before paid enterprise adoption: support model, legal clarity, compliance controls.

## 3) GTM Sequence (Internal)

1. Grow OSS adoption on `L0-L2` with low-friction onboarding.
2. Validate value with design partners in high-assurance domains.
3. Productize private enterprise capabilities in closed repos.
4. Keep public docs focused on OSS protocol/CLI/orchestrator narrative.

## 4) Hard Boundaries

- Public repo must not contain pricing, ARR, sales targets, or paid-tier promises.
- Public architecture and roadmap remain OSS-first (`L0-L2`).
- Commercial roadmap details stay in local/private planning artifacts only.

## 5) Notes

- This document aggregates the commercial planning context moved out of public vision text.
- If this branch is ever proposed for merge to public docs, this file should be reviewed and likely excluded.
