# SDP Feature Map

> Organized by protocol layer, not by release phase.

---

## Protocol Layer (Open Standard)

### Core Protocol

| Feature | Description | Priority |
|---------|-------------|----------|
| **plan** | Decompose feature → atomic units with dependencies | P0 |
| **apply** | Generate → verify → record per unit | P0 |
| **evidence** | Cryptographic chain: spec → code → verification → approval | P0 |
| **incident** | Trace from commit back through evidence chain | P0 |

### Model Provenance

| Feature | Description | Priority |
|---------|-------------|----------|
| **Model identity** | Model name, version, provider | P0 |
| **Generation params** | Temperature, prompt hash (not prompt), timestamp | P0 |
| **Spec reference** | Which spec this code was generated against | P0 |
| **Approval record** | Who approved, when, what they saw | P0 |

### Evidence Format

| Feature | Description | Priority |
|---------|-------------|----------|
| **Hash chain** | Tamper-evident linking of evidence records | P0 |
| **JSON format** | Human-readable + machine-parseable | P0 |
| **Verification output** | Actual command output, not "tests passed" | P0 |
| **Decision log** | What was decided and why (drive/interactive mode) | P0 |

---

## Verification Engine (Open Source)

### Verification Stack

| Layer | What It Catches | Cost | Priority |
|-------|----------------|------|----------|
| **Type checking** | Type errors (#1 AI bug class) | Free | P0 |
| **Static analysis** | Security patterns, anti-patterns (semgrep) | Free | P0 |
| **Test execution** | Functional correctness | Free | P0 |
| **Coverage gates** | Untested code paths | Free | P0 |
| **Property-based testing** | Edge cases, invariant violations | 1.2x tokens | P1 |
| **Cross-model review** | Correlated blind spots | 2.5x tokens | P2 |
| **Snapshot testing** | Unintended behavioral changes | Free | P2 |

### Decomposition Engine

| Feature | Description | Priority |
|---------|-------------|----------|
| **NL → Units** | Natural language → atomic workstream units | P0 |
| **Dependency Graph** | Topological sort, parallel-safe execution order | P0 (exists) |
| **Scope Inference** | Auto-detect which files a unit should touch | P0 |
| **Acceptance Criteria** | Auto-generate testable criteria per unit | P0 |
| **Adaptive Granularity** | Unit size scales with risk (100 LOC payments, 300 LOC UI) | P1 |
| **Decomposition Heuristics** | Learned patterns: "OAuth → 3 units, not 5" | P2 (moat) |

### Risk Engine

| Feature | Description | Priority |
|---------|-------------|----------|
| **Path-based risk** | `auth/`, `payments/`, `infra/` → high risk | P1 |
| **Content-based risk** | SQL queries, crypto, token handling → flag | P1 |
| **History-based risk** | Files with high bug rate → extra verification | P2 |
| **Custom risk profiles** | `.sdp.yml` overrides per project | P1 |
| **Model selection policy** | Route high-risk to capable models, low-risk to fast models | P1 |

---

## Orchestration (Proprietary)

### plan/apply UX

| Feature | Description | Priority |
|---------|-------------|----------|
| **`sdp plan`** | Show decomposition before execution | P0 |
| **`sdp apply`** | Execute plan with streaming progress | P0 |
| **`sdp incident`** | Forensic trace from commit to evidence chain | P0 |
| **`--auto-apply`** | Plan + apply immediately (ship mode) | P0 |
| **`--interactive`** | Stop at every fork (drive mode) | P1 |
| **`--retry N`** | Regenerate only failed unit N | P0 |
| **`--output=json`** | JSON output for tool integration | P0 |
| **`--dry-run`** | Show plan only, equivalent to `plan` | P0 |
| **Cost estimate** | "~3 min, 3 units, ~$0.15" before execution | P1 |

### Drive Mode (interactive)

| Feature | Description | Priority |
|---------|-------------|----------|
| **Decision points** | AI stops and asks: "Sessions or JWT?" | P1 |
| **Decision recording** | Every human choice logged in evidence chain | P1 |
| **Stage gates** | Human approval at each stage | P1 |
| **Expertise encoding** | Human decisions become spec constraints | P1 |
| **Progressive disclosure** | Start with 3 questions, deepen only if needed | P1 |

### Failure UX

| Feature | Description | Priority |
|---------|-------------|----------|
| **Graceful degradation** | Clear errors when decomposition fails | P1 |
| **"Can't decompose" UX** | Explain why, suggest alternatives | P1 |
| **Partial success** | "3/5 units passed, 2 need attention" | P0 |
| **Hallucination detection** | Flag when AI invents dependencies/files | P1 |

---

## Tools (User Surfaces)

### Claude Code Plugin

| Feature | Description | Priority |
|---------|-------------|----------|
| **Protocol compliance** | Existing skills emit evidence chain | P0 |
| **Model provenance** | Track provenance on every `@build` | P0 |
| **`sdp plan` skill** | Plan/apply wrapper for existing workflow | P0 |
| **`sdp incident` skill** | Forensic trace in Claude Code | P1 |

### CLI

| Feature | Description | Priority |
|---------|-------------|----------|
| **`sdp plan`** | Standalone plan command | P1 |
| **`sdp apply`** | Standalone apply command | P1 |
| **`sdp incident`** | Forensic trace command | P1 |
| **Streaming progress** | Real-time progress per unit | P1 |

### CI/CD

| Feature | Description | Priority |
|---------|-------------|----------|
| **GitHub Action** | `sdp-dev/verify-action@v1` — one-YAML verification | P1 |
| **PR evidence comment** | Evidence chain summary in PR | P1 |
| **Provenance check** | Verify all AI code has provenance | P1 |
| **GitLab CI** | Same for GitLab | P2 |

### IDE

| Feature | Description | Priority |
|---------|-------------|----------|
| **Cursor plugin** | Plan/apply from IDE (highest IDE priority) | P2 |
| **VS Code extension** | Plan/apply from command palette | P2 |
| **JetBrains plugin** | Plan/apply from IDE | P3 |

### SDK

| Feature | Description | Priority |
|---------|-------------|----------|
| **`sdp.Verify()`** | Verification engine as library | P2 |
| **`sdp.Evidence()`** | Evidence bundle generation | P2 |
| **`sdp.Audit()`** | Audit trail recording | P2 |
| **Provider adapters** | Claude, GPT, Gemini, local models | P2 |
| **JSON-in/JSON-out** | API for external tool integration | P1 |

---

## Enterprise Features

| Feature | Description | Priority |
|---------|-------------|----------|
| **Compliance export** | SOC2/HIPAA/DORA-ready evidence format | P1 |
| **Verification certificates** | Signed per-PR proof of verification | P1 |
| **Vanta/Drata integration** | Export to compliance platforms | P2 |
| **Team templates** | Shared decomposition patterns | P1 |
| **Conflict detection** | Two devs, same codebase awareness | P2 |
| **Team policies** | "All AI PRs require evidence chain" | P2 |
| **Billing/metering** | Usage tracking, invoicing | P1 |
| **Team dashboards** | "What % of AI code is SDP-verified?" | P3 |
| **Policy enforcement** | Org-wide verification requirements | P3 |
| **SSO/SAML** | Enterprise auth | P3 |
| **On-premise** | Self-hosted for air-gapped environments | P3 |

---

## Data & Moat

| Feature | Description | Priority |
|---------|-------------|----------|
| **Verification telemetry** | Instrument every run: catch rate, iterations | P0 |
| **AI failure taxonomy** | Categorized: what AI gets wrong by model/language/domain | P1 |
| **Decomposition heuristics** | Learned patterns from verified builds | P2 |
| **Benchmark suite** | "Generate 50 features, measure defect rates" | P2 |
| **Quarterly benchmark** | "AI Code Quality Benchmark" publication | P1 |

---

## Priority Legend

| Priority | Meaning | Phase |
|----------|---------|-------|
| **P0** | Protocol + plugin launch | Phase 1 (Weeks 1-6) |
| **P1** | CLI + enterprise foundation | Phase 2 (Months 2-4) |
| **P2** | SDK + platform expansion | Phase 3 (Months 4-8) |
| **P3** | Scale + enterprise platform | Phase 4 (Months 8-12) |

---

*SDP Features v2.0 — February 2026*
