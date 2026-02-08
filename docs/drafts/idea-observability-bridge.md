# F059: Observability Bridge Design

> Beads: sdp-pom6 | Priority: P1

---

## Problem

Evidence log answers "what happened during build." But when something breaks in production, you need to connect the evidence to runtime:
- Which deploy introduced this change?
- Which code paths are AI-generated?
- Which lines were written by which model?

Without this bridge, evidence stays isolated from production observability.

## Solution

A design document (not implementation — P2 builds it) specifying how SDP evidence connects to observability tools.

### 1. Deploy Markers

Tie evidence records to deploy events:
```
evidence record → commit SHA → deploy event → production issue
```

Spec: what metadata SDP adds to deploy events (commit range, evidence count, model distribution).

### 2. OTel Span Attributes

Mark AI-generated code paths in traces:
```
span.setAttribute("sdp.ai_generated", true)
span.setAttribute("sdp.model", "claude-sonnet-4")
span.setAttribute("sdp.evidence_id", "evt-abc123")
```

Spec: OTel semantic convention proposal for AI-generated code.

### 3. Diff-Level Provenance

Which lines are AI-generated vs human-edited:
```
git blame + evidence → per-line attribution
```

Spec: how to compute per-line AI/human attribution from evidence + git history.

### 4. Integration Spec

How SDP connects to:
- Honeycomb (structured events)
- Datadog (custom metrics + traces)
- Grafana (Loki for logs, Tempo for traces)
- Generic OTel collector

## Deliverables

1. `docs/design/OBSERVABILITY-BRIDGE.md` — the design document
2. OTel semantic convention draft (informal)
3. Integration architecture diagrams
4. Data flow: evidence → deploy → runtime → incident

## Users

- SRE teams correlating incidents with AI code changes
- Platform engineers integrating SDP into observability stack
- SDP developers implementing P2 observability features

## Success Metrics

- Design doc reviewed by 2+ SRE practitioners
- Clear enough to implement without design questions
- Covers all 4 integration points

## Dependencies

- F054 (evidence schema — must know what data is available)

## Notes

- This is a DESIGN feature, not implementation
- Implementation is P2 (F063 in roadmap)
- Value: ensures P2 implementation is well-designed from day 1
- Avoids building OTel integration that doesn't match real observability needs
