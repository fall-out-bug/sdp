# SDP Roadmap

> Protocol → Plugin → CLI → SDK → Enterprise

---

## Current State (February 2026)

### What Exists

- **Claude Code plugin**: 19 skills/agents (idea, design, build, review, deploy, vision, reality)
- **Go engine**: decomposition, dependency graph, parallel dispatcher, synthesis engine, circuit breaker
- **Verification**: TDD pipeline (red → green → refactor), coverage enforcement, contract validation
- **Enterprise traction**: 2 contracts (bank, airline), 1 evaluating (marketplace)

### What's Missing

- **Protocol specification**: no formal spec that other tools can implement
- **Model provenance**: no tracking of which model generated what
- **Evidence chain**: no cryptographic linking of spec → code → verification → approval
- **`sdp plan` / `sdp apply`**: the plan/apply UX doesn't exist
- **`sdp incident`**: no forensic trace tool
- **Data collection**: no measurement of verification effectiveness

---

## Phase 1: The Protocol (Weeks 1-6)

**Goal:** Formalize what exists into an open protocol. Ship the evidence chain.

### 1.1 Protocol Specification (Weeks 1-2)

Formal spec: what any SDP-compatible tool must do.

- [ ] **plan**: decompose feature → atomic units with dependencies
- [ ] **apply**: generate → verify → record per unit
- [ ] **evidence**: linked chain (spec → code → verification → approval)
- [ ] **incident**: trace from commit back through evidence chain
- [ ] Model provenance format (model, version, prompt hash, timestamp)
- [ ] Evidence format (JSON, cryptographic hash chain)
- [ ] Verification report format (actual command output, not assertions)

**Output:** `docs/protocol/SDP-SPEC-v1.md` — the open standard.

### 1.2 Model Provenance (Weeks 1-3)

Every piece of AI-generated code records:

- [ ] Model name and version
- [ ] Prompt hash (not the prompt itself — privacy)
- [ ] Temperature and parameters
- [ ] Timestamp
- [ ] Spec it was generated against
- [ ] Who initiated the generation

**This is P0.** Without provenance, the evidence chain is broken and compliance export is garbage.

### 1.3 Evidence Chain (Weeks 2-4)

Cryptographically linked record:

- [ ] spec → generation → verification → approval
- [ ] Tamper-evident hashing (each step references previous hash)
- [ ] Stored in `.sdp/evidence/` alongside the repo
- [ ] Human-readable + machine-parseable (JSON)
- [ ] `sdp incident <commit>` reads and presents the chain

### 1.4 Claude Code Plugin Update (Weeks 3-6)

Update existing skills to produce protocol-compliant artifacts:

- [ ] `@build` emits evidence chain
- [ ] `@review` emits verification report with actual output
- [ ] Model provenance tracked on every generation
- [ ] `sdp plan` / `sdp apply` skill wrappers

### 1.5 Data Collection (Week 1, ongoing)

- [ ] Instrument every run: decomposition quality, verification catch rate, iteration count
- [ ] Rich telemetry (not A/B test — observational dataset)
- [ ] Target: "AI Code Quality Benchmark" quarterly publication

**Success criteria:** Protocol spec published. Evidence chain working in Claude Code plugin. Provenance on every generation.

---

## Phase 2: The Tools (Months 2-4)

**Goal:** CLI + GitHub Action. Enterprise can now adopt.

### 2.1 CLI: plan / apply / incident (Months 2-3)

```bash
sdp plan "Add OAuth2 login"            # Show decomposition
sdp apply                               # Execute plan
sdp plan "Add auth" --auto-apply        # Ship mode
sdp plan "Add auth" --interactive       # Drive mode
sdp apply --retry 3                     # Retry failed unit
sdp incident abc1234                    # Forensic trace
```

- [ ] `sdp plan` — NL → decomposition → plan display
- [ ] `sdp apply` — verified execution with streaming progress
- [ ] `sdp incident` — forensic trace from commit to evidence chain
- [ ] Per-unit rollback
- [ ] JSON output (`--output=json`) for tool integration

### 2.2 GitHub Action (Month 2-3)

```yaml
- uses: sdp-dev/verify-action@v1
  with:
    gates: [types, semgrep, tests, provenance]
```

- [ ] Verification on every AI-generated PR
- [ ] Evidence chain in PR comment
- [ ] Model provenance check
- [ ] Free tier: 50 runs/month. Usage-based pricing above

### 2.3 Compliance Foundation (Month 3-4)

- [ ] Evidence export: SOC2-ready format
- [ ] Verification certificates (signed, timestamped)
- [ ] Vanta/Drata integration
- [ ] DORA compliance mapping for EU fintech
- [ ] Audit trail spec: model provenance, prompt hash, verification output, approval chain

### 2.4 Risk-Proportional Verification (Month 3-4)

- [ ] Path-based risk: `auth/`, `payments/` → full trail. `components/` → light
- [ ] Content-based risk: SQL, crypto, tokens → flag
- [ ] Custom risk profiles via `.sdp.yml`

**Success criteria:** CLI shipped. GitHub Action in 100+ repos. 3 enterprise customers on paid tier. First "AI Code Quality Benchmark" published.

---

## Phase 3: The Platform (Months 4-8)

**Goal:** SDK for embedding. Enterprise features. Evidence standard.

### 3.1 Verification SDK

```go
result := sdp.Verify(files, sdp.Gates{"types", "tests", "semgrep"})
evidence := sdp.Evidence(result)  // Signed evidence bundle
```

- [ ] `sdp.Verify()` — verification engine as library
- [ ] `sdp.Evidence()` — evidence bundle generation
- [ ] `sdp.Audit()` — audit trail recording
- [ ] Provider adapters: Claude, GPT, Gemini
- [ ] JSON-in/JSON-out API for tool integration

**Note:** No `sdp.Decompose()` in SDK. Decomposition stays in CLI/plugin where we control the experience. SDK = verification + evidence only.

### 3.2 Cross-Model Review (Premium)

- [ ] Model A generates, Model B reviews (decorrelated errors)
- [ ] Auto-triggered for high-risk code (auth, payments, data deletion)
- [ ] Model selection policy engine
- [ ] Justified by evidence dataset, not thesis

### 3.3 Team Features

- [ ] Shared decomposition templates ("at our company, auth always = 5 units")
- [ ] Conflict detection (two devs, same codebase)
- [ ] Team-wide verification policies
- [ ] Billing/metering infrastructure

### 3.4 IDE Integrations

- [ ] Cursor plugin (highest priority — plan/apply from IDE)
- [ ] VS Code extension
- [ ] JetBrains plugin

**Success criteria:** SDK embedded in 2+ external tools. Team features used by 10+ enterprise teams.

---

## Phase 4: The Standard (Months 8-12+)

**Goal:** SDP becomes how enterprises prove AI code was verified.

### 4.1 Evidence Standard

- [ ] Published spec: "SDP Evidence Format v1.0"
- [ ] Adopted by 2+ tools beyond SDP itself
- [ ] Industry working group (if traction warrants)

### 4.2 Enterprise Platform

- [ ] Multi-team dashboards: "what % of AI code is SDP-verified?"
- [ ] Policy enforcement: "all AI PRs must have evidence chain"
- [ ] Audit export: SOC2 / ISO 27001 / HIPAA / DORA
- [ ] SSO/SAML
- [ ] On-premise for air-gapped environments

### 4.3 Ecosystem

- [ ] Third-party verification gates
- [ ] Custom risk profiles
- [ ] Community-contributed decomposition heuristics
- [ ] AI failure taxonomy (public dataset)

**Success criteria:** SDP Evidence Format cited in enterprise security policies. Acquisition or deep integration interest.

---

## Success Metrics

| Metric | Phase 1 | Phase 2 | Phase 3 | Phase 4 |
|--------|---------|---------|---------|---------|
| Protocol implementations | 1 (plugin) | 2 (+ CLI) | 4 (+ SDK, GH Action) | 10+ |
| Repos with evidence chains | 10 (dogfood) | 200 | 1,000 | 5,000+ |
| Enterprise customers | 2 (existing) | 5 | 20 | 50+ |
| Revenue | $0 | $50K ARR | $500K ARR | $2M+ ARR |
| Evidence records generated | 100 | 5,000 | 50,000 | 500,000 |

---

## Kill Criteria

> After 500 SDP runs: if verification catch rate < 5% AND post-merge defect rate is not measurably different from baseline → kill the product.

Specific. Testable. Honest.

---

## Key Risks

| Risk | Mitigation |
|------|-----------|
| Models improve, decomposition unnecessary | Monitor. Pivot to verification-only if evidence shows no decomposition benefit |
| Cursor/GitHub ships native verification | Protocol is the moat. If they implement SDP protocol, we win |
| Enterprise sales cycle too slow | GitHub Action as self-serve wedge. Bottom-up → top-down |
| No measurable defect reduction | Kill the product (see kill criteria) |
| Open protocol gets forked | Move fast. Build the evidence dataset. Data moat > code moat |

---

*SDP Roadmap v2.0 — February 2026*
*Protocol → Plugin → CLI → SDK → Enterprise*
