# F061: Data Collection & AI Failure Benchmark

> Beads: sdp-6fgr | Priority: P1

---

## Problem

We claim SDP improves code quality. We have no data to prove it. Enterprise customers ask "what's the catch rate?" and we can't answer with numbers.

## Solution

Instrument SDP to collect structured metrics, then publish an "AI Code Quality Benchmark."

### What We Collect

| Metric | Source | Purpose |
|--------|--------|---------|
| Catch rate | Verification events | % of AI generations that fail verification |
| Iteration count | Build events | How many red→green cycles per workstream |
| Model performance | Generation events | Pass rate by model, language, domain |
| Acceptance test catch rate | Acceptance events | % caught by acceptance vs unit tests |
| MTTR | Incident + evidence | Time to root cause with vs without evidence |
| Defect escape rate | Post-deploy events | Bugs that make it to production despite verification |

### AI Failure Taxonomy

Classify failures by:
- **Model**: Claude vs GPT vs Gemini vs open-source
- **Language**: Go vs Python vs TypeScript vs Rust
- **Domain**: auth, payments, data processing, UI
- **Failure type**: wrong logic, missing edge case, hallucinated API, type error, test-passing-but-wrong
- **Severity**: caught by verification, caught by acceptance, escaped to production

### Publication

"AI Code Quality Benchmark" — quarterly:
- Aggregate catch rates across all instrumented repos
- Model comparison (anonymized)
- Failure taxonomy (most common AI mistakes)
- Trend over time (are models getting better?)

### Privacy

- Only aggregate metrics published — never raw evidence
- Opt-in: repos must enable telemetry sharing
- No code, no prompts, no company names
- Methodology published alongside data

## Users

- Enterprise evaluating SDP (hard data for ROI conversations)
- AI model providers (feedback on failure modes)
- Industry researchers (benchmark for AI coding quality)
- SDP developers (prioritize what to improve)

## Success Metrics

- First benchmark published with data from 10+ repos
- Catch rate > 5% (SDP finds real issues)
- At least 1 enterprise cites benchmark in evaluation

## Dependencies

- F054 (evidence log — data source)
- F056 (full instrumentation — complete data)

## Notes

- This is the "trust standard" play — data-driven, not marketing
- Start with SDP's own dogfooding data (repo sdp itself)
- Expand to enterprise repos with opt-in telemetry
- Benchmark positions SDP as the authority on AI code quality measurement
