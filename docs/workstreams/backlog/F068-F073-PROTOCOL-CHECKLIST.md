# F068-F073 Protocol Checklist (UX-First Track)

> Execution checklist for roadmap implementation after F067.
> Scope: F068 (first-run UX) -> F073 (trust/explainability).

## 1) Entry Criteria (Before F068 starts)

- [ ] F067 closure report is available and linked
- [ ] Current CLI command map is documented (`help/status/plan/apply/log`)
- [ ] Baseline metrics source is agreed (events log, manual runbook, or both)
- [ ] Owners assigned for UX, CLI, protocol, and docs
- [ ] Feature dependencies and sequencing approved

## 2) Protocol Decomposition Quality Gate (Per Feature)

- [ ] Feature has summary file with dependency graph and registry
- [ ] Every workstream has explicit `depends_on` and `blocks`
- [ ] Every workstream has measurable acceptance criteria (minimum 5)
- [ ] `scope_files` includes canonical targets only (no ambiguous paths)
- [ ] Definition of Done includes artifact + metric validation

## 3) UX Learning Curve Checklist (F068/F069)

- [ ] New user reaches first successful run in <= 15 minutes (target)
- [ ] `sdp init --guided` path covers prerequisites and next command
- [ ] `sdp --help` and docs present intent-first command journeys
- [ ] `sdp status` always surfaces deterministic next action
- [ ] Recommendation quality metrics are captured and reviewed

## 4) Failure and Recovery Checklist (F070)

- [ ] Failure taxonomy and error codes are documented and test-covered
- [ ] Every critical failure class has inline recovery playbook
- [ ] Resume/checkpoint flow exposes clear options and safe defaults
- [ ] Post-failure diagnostics report is generated automatically
- [ ] MTTR target is defined and measured during rollout

## 5) Team UX Checklist (F071)

- [ ] Team operating model (roles, ownership, handoff) is published
- [ ] Handoff package format is standardized and validated
- [ ] Scope collision signals are visible before merge-time conflicts
- [ ] Team notifications avoid sensitive data leakage
- [ ] Pilot team adoption report is produced with follow-up backlog

## 6) Interop and Migration Checklist (F072)

- [ ] Interop matrix defines supported tools, versions, and fidelity tiers
- [ ] Import path supports dry-run + validated write mode
- [ ] Export path marks lossy/unsupported fields explicitly
- [ ] Migration wizard produces auditable migration report
- [ ] Compatibility/deprecation policy is versioned and linked in docs

## 7) Trust and Explainability Checklist (F073)

- [ ] Explainability levels (`brief/standard/audit`) are contract-defined
- [ ] Recommendations include rationale and confidence
- [ ] Trace output follows canonical narrative order
- [ ] Gate failures include explicit policy reason + remediation
- [ ] Public trust pack and GA criteria are published

## 8) Cross-Cutting Release Gate

- [ ] UX KPIs and technical KPIs are both required for release
- [ ] Every feature has completion report with metric delta vs baseline
- [ ] Follow-up issues created for unmet targets
- [ ] README and PRODUCT_VISION reflect current strategic direction
- [ ] Roadmap remains consistent with protocol decomposition files

## 9) Multi-Repo Readiness Checklist (Vision Alignment)

- [ ] Boundaries are explicit between `protocol`, `cli`, and `orchestrator`
- [ ] Workstream scope references the intended layer ownership
- [ ] Contracts between layers are documented before implementation
- [ ] No feature introduces hidden coupling that blocks future repo split
- [ ] Backlog items include migration notes for eventual multi-repo extraction

## 10) Competitive Adoption Traceability

- [ ] Competitive adoption map exists and is versioned:
  `/Users/fall_out_bug/projects/vibe_coding/sdp/docs/reference/2026-02-14-competitive-analysis-sdp-vs-oss.md`
- [ ] Every adopted external pattern is mapped to a concrete WS ID
- [ ] Every non-adopted pattern has an explicit rationale (scope, complexity, or governance)
- [ ] F072 interop coverage includes OpenSpec, Task Master, and CCPM mapping baselines
- [ ] F070/F071 include memory continuity and handoff learnings without replacing SDP evidence chain
