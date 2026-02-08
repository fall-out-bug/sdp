# SDP Feature Map

> Everything SDP does, organized by mode and phase.

---

## Core Engine (Both Modes)

### Decomposition Engine

| Feature | Description | Priority |
|---------|-------------|----------|
| **NL → Units** | Natural language feature description → atomic workstream units | P0 |
| **Dependency Graph** | Topological sort, parallel-safe execution order | P0 (exists) |
| **Scope Inference** | Auto-detect which files a unit should touch | P0 |
| **Acceptance Criteria** | Auto-generate testable criteria per unit | P0 |
| **Adaptive Granularity** | Unit size scales with risk: 100 LOC for payments, 300 LOC for UI | P1 |
| **Decomposition Heuristics** | Learned from verified builds: "OAuth → 3 units, not 5" | P2 (moat) |

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

### Risk Engine

| Feature | Description | Priority |
|---------|-------------|----------|
| **Path-based risk** | `auth/`, `payments/`, `infra/` → high risk | P1 |
| **Content-based risk** | SQL queries, crypto, token handling → flag | P1 |
| **History-based risk** | Files with high bug rate → extra verification | P2 |
| **Custom risk profiles** | `.sdp.yml` overrides per project | P1 |

### Audit Trail

| Feature | Description | Priority |
|---------|-------------|----------|
| **Verification evidence** | Actual command output, not "should pass" | P0 (exists) |
| **Decision log** | What was decided and why (drive mode) | P0 |
| **Provenance record** | Which model, which spec, what cost, when | P1 |
| **Compliance export** | SOC2/HIPAA-ready audit format | P2 |
| **Verification certificates** | Signed per-PR proof of verification | P2 |

---

## Ship Mode Features

| Feature | Description | Priority |
|---------|-------------|----------|
| **`sdp ship "description"`** | One-command: decompose → plan → verify → PR | P0 |
| **Plan-by-default** | Show decomposition, ask Y/n/edit before execution | P0 |
| **Streaming progress** | Real-time progress bar per unit | P0 |
| **Per-unit rollback** | Failed unit 3? Retry just unit 3. | P0 |
| **`--auto-approve`** | Skip plan approval (trust the decomposition) | P0 |
| **`--cross-review`** | Add cross-model review for high-risk code | P1 |
| **`--dry-run`** | Show plan only, don't execute | P0 |
| **`--edit`** | Open plan in editor before execution | P1 |
| **`--retry N`** | Regenerate only failed unit N | P0 |
| **Cost estimate** | "~3 min, 3 units, ~$0.15" before execution | P1 |
| **Model routing** | Fast model for simple units, capable for complex | P2 |

---

## Drive Mode Features

| Feature | Description | Priority |
|---------|-------------|----------|
| **`sdp drive "description"`** | Full pipeline, human decides at every fork | P0 |
| **Decision points** | AI stops and asks: "Sessions or JWT?" | P0 |
| **Decision recording** | Every human choice is logged in artifacts | P0 |
| **Stage gates** | Human approval required: idea → design → build → review | P0 |
| **Expertise encoding** | Human decisions become spec constraints for AI | P0 |
| **Progressive disclosure** | Start with 3 questions, go deeper only if needed | P1 |
| **Domain teaching** | "In this project, we always use X for Y" → remembered | P2 |
| **Compliance sign-off** | Human approval at each gate for regulated environments | P1 |

---

## Platform & Integration Features

### CI/CD

| Feature | Description | Priority |
|---------|-------------|----------|
| **GitHub Action** | `sdp-dev/verify-action@v1` — one-YAML verification | P0 |
| **PR comments** | Verification report as PR comment | P0 |
| **GitLab CI** | Same as GitHub Action for GitLab | P2 |
| **Branch protection** | "All AI PRs must pass SDP verification" | P2 |

### IDE

| Feature | Description | Priority |
|---------|-------------|----------|
| **CLI** | First-class, always supported | P0 |
| **Cursor plugin** | Native integration (highest IDE priority) | P1 |
| **VS Code extension** | Ship/drive from command palette | P2 |
| **JetBrains plugin** | Ship/drive from IDE | P3 |

### SDK

| Feature | Description | Priority |
|---------|-------------|----------|
| **`sdp.Decompose()`** | Programmatic decomposition | P2 |
| **`sdp.Verify()`** | Programmatic verification | P2 |
| **`sdp.Audit()`** | Programmatic audit recording | P2 |
| **Provider adapters** | Claude, GPT, Gemini, local models | P2 |

---

## Enterprise Features

| Feature | Description | Priority |
|---------|-------------|----------|
| **Team dashboards** | "What % of AI code is SDP-verified?" | P3 |
| **Policy enforcement** | Org-wide: "all AI PRs require SDP" | P3 |
| **Audit export** | SOC2, ISO 27001, HIPAA formats | P2 |
| **SSO/SAML** | Enterprise auth | P3 |
| **On-premise** | Self-hosted for air-gapped environments | P3 |
| **Custom gate marketplace** | Third-party verification rules | P3 |

---

## Data & Moat Features

| Feature | Description | Priority |
|---------|-------------|----------|
| **Defect tracking** | What SDP caught that would have shipped | P1 |
| **Heuristic learning** | Decomposition patterns that reduce defects | P2 |
| **AI failure taxonomy** | Categorized dataset of "what AI gets wrong" | P2 |
| **Benchmark suite** | "Generate these 50 features, measure defect rates" | P2 |
| **Public dataset** | Open evidence that decomposition reduces defects | P1 |

---

## Priority Legend

| Priority | Meaning | Timeline |
|----------|---------|----------|
| **P0** | Must have for launch | Phase 1 (Weeks 1-6) |
| **P1** | Must have for enterprise | Phase 2 (Months 2-4) |
| **P2** | Platform expansion | Phase 3 (Months 4-8) |
| **P3** | Scale/enterprise | Phase 4 (Months 8-12) |

---

*SDP Features v1.0 — February 2026*
