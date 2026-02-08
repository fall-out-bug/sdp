# SDP Roadmap

> From philosophy to product. From product to standard.

---

## Current State (February 2026)

### What Exists

- **Go-based core**: decomposition engine, dependency graph, parallel dispatcher, synthesis engine
- **19 skills/agents**: idea, design, build, review, deploy, vision, reality, etc.
- **Contract validation**: OpenAPI spec comparison between frontend/backend
- **Working TDD pipeline**: red → green → refactor with coverage enforcement
- **Enterprise contracts**: 2 active (bank, airline), 1 evaluating (marketplace)

### What's Missing

- **`sdp ship` command**: the one-command UX doesn't exist yet
- **`sdp drive` command**: the human-in-the-loop pipeline is manual (separate skill calls)
- **GitHub Action**: no CI/CD integration
- **Evidence dataset**: no data proving SDP reduces defects
- **Risk-proportional gates**: everything gets the same treatment regardless of risk

---

## Phase 1: The Wedge (Weeks 1-6)

**Goal:** Get SDP into repos. Prove the thesis with data.

### 1.1 Ship `sdp ship` (Weeks 1-3)

The one-command experience:

```bash
sdp ship "Add OAuth2 login"
# Planning... → Proceed? [Y/n/edit] → Building... → PR created
```

**Deliverables:**
- [ ] Decomposition from natural language → workstream units
- [ ] Plan-by-default UX (terraform plan style)
- [ ] Sequential verified execution with streaming progress
- [ ] Per-unit rollback on failure
- [ ] PR creation with verification summary

**Quality gates (default):** types + tests + static analysis

### 1.2 Ship `sdp drive` (Weeks 2-4)

The human-in-the-loop experience:

```bash
sdp drive "Add OAuth2 login"
# Stops at every fork: "Sessions or JWT?" → you decide
```

**Deliverables:**
- [ ] Interactive decomposition with decision points
- [ ] Human approval gates at each stage (idea → design → build → review)
- [ ] Decision recording in artifacts (why sessions, not JWT)
- [ ] Same verification pipeline as ship mode

### 1.3 Ship the GitHub Action (Weeks 3-5)

```yaml
- uses: sdp-dev/verify-action@v1
  with:
    gates: [types, semgrep, tests]
```

**Deliverables:**
- [ ] GitHub Action: `sdp-dev/verify-action@v1`
- [ ] Default gates: type checking + static analysis + test execution
- [ ] PR comment with verification report
- [ ] Free for open source, paid for private repos

### 1.4 Start Collecting Data (Week 1, ongoing)

- [ ] 10 features built with `sdp ship` vs 10 without
- [ ] Metrics: defects found in review, time to merge, rework rate, token spend
- [ ] Internal dogfooding on SDP repo itself

**Success criteria:** 100 repos using GitHub Action by end of Phase 1.

---

## Phase 2: The Evidence (Months 2-4)

**Goal:** Prove SDP works with data. Publish results.

### 2.1 Publish the Dataset

- [ ] "SDP-verified code: X% fewer defects, Y% faster to merge, Z% less rework"
- [ ] Open dataset, reproducible methodology
- [ ] Blog posts with real numbers, not claims
- [ ] Case study from enterprise contracts (anonymized)

### 2.2 Risk-Proportional Verification

- [ ] Automatic risk inference from file path / module type
- [ ] `payments/`, `auth/`, `infra/` → full verification + audit trail
- [ ] `components/`, `pages/`, `utils/` → types + tests only
- [ ] User-configurable overrides via `.sdp.yml`

### 2.3 Property-Based Testing

- [ ] AI-generated Hypothesis/fast-check strategies per unit
- [ ] Second-highest ROI verification layer
- [ ] Auto-detect: "this function handles money → generate property tests for commutativity and rounding"

### 2.4 Compliance Reporting

- [ ] SOC2-ready audit trails per feature
- [ ] Per-PR verification certificates (signed, timestamped)
- [ ] Export to compliance tools (Vanta, Drata, etc.)
- [ ] This is where enterprise money lives

**Success criteria:** Published dataset with statistically significant defect reduction. 3+ enterprise customers on paid plan.

---

## Phase 3: The Platform (Months 4-8)

**Goal:** Let others build on SDP. Expand from CLI to SDK.

### 3.1 Agent SDK

```go
result := sdp.Ship("Add OAuth2 login", sdp.Options{
    Gates:    []string{"types", "tests", "semgrep"},
    AutoApprove: false,
})
```

- [ ] `sdp.Decompose()` — break feature into units
- [ ] `sdp.Verify()` — run verification gates on code
- [ ] `sdp.Audit()` — record verification evidence
- [ ] Embeddable in other tools (IDE plugins, CI systems, platforms)

### 3.2 Cross-Model Review (Premium)

- [ ] High-risk code only (auto-detected or user-flagged)
- [ ] Model A generates, Model B reviews (decorrelated errors)
- [ ] Justified by evidence dataset, not by thesis
- [ ] Premium pricing tier

### 3.3 Multi-Provider Support

- [ ] Claude, GPT, Gemini as generation backends
- [ ] Provider-agnostic decomposition (same units, any model)
- [ ] Cost-optimal routing: "use fast model for utils, capable model for auth"

### 3.4 IDE Integrations

- [ ] Cursor plugin (native, highest priority)
- [ ] VS Code extension
- [ ] JetBrains plugin
- [ ] CLI remains first-class citizen

**Success criteria:** SDK adopted by 2+ external tools. Cross-model review shows measurable improvement over single-model.

---

## Phase 4: The Standard (Months 8-12+)

**Goal:** If we win, extract the protocol.

### 4.1 Extract the Protocol

- [ ] Thin interface: Decompose → Verify → Audit
- [ ] Like Rack (Ruby) or WSGI (Python) — the minimum contract for composability
- [ ] Only if 1000+ repos use SDP (protocol from success, not from planning)

### 4.2 Ecosystem

- [ ] Third-party verification gates (custom semgrep rulesets, domain-specific tests)
- [ ] Marketplace for decomposition heuristics
- [ ] Community-contributed risk profiles

### 4.3 Enterprise Platform

- [ ] Multi-team dashboards: "what % of AI code is SDP-verified?"
- [ ] Policy enforcement: "all AI PRs must pass SDP verification"
- [ ] Audit export for SOC2/ISO 27001/HIPAA

**Success criteria:** Industry recognition. Cursor/GitHub integration or acquisition interest.

---

## Success Metrics

| Metric | Phase 1 | Phase 2 | Phase 3 | Phase 4 |
|--------|---------|---------|---------|---------|
| Repos using SDP | 100 | 500 | 2,000 | 10,000 |
| Enterprise customers | 3 | 10 | 30 | 100 |
| Defect reduction (proven) | — | X% (published) | X% (reproduced) | Industry benchmark |
| Revenue | $0 | $50K ARR | $500K ARR | $2M+ ARR |
| Integrations | CLI + GH Action | + Compliance | + SDK + IDE | + Protocol |

---

## Key Risks

| Risk | Mitigation |
|------|-----------|
| Models improve so fast decomposition becomes unnecessary | Monitor. If 5000-line generations verify clean consistently, pivot to verification-only |
| Cursor/GitHub ships native verification | Be the engine they embed, or the standard they adopt. Speed matters. |
| Enterprise sales cycle too slow | GitHub Action as self-serve wedge. Bottom-up adoption → top-down purchase |
| No measurable defect reduction | Kill the product. The thesis is testable — if data says no, we're wrong |

---

*SDP Roadmap v1.0 — February 2026*
