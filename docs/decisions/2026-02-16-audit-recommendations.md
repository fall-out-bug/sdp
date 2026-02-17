# Audit Recommendations — 2026-02-16

Consolidated from: repo audit (Cloud Agent), Codex review, owner feedback, roadmap analysis.

Status: **PLAN** — items to be scheduled into workstreams.

---

## Done (this session)

- [x] Remove `uv.lock` — was triggering Python dependabot PRs
- [x] Remove tracked `sdp-plugin/test_output.txt` (73KB) + add to .gitignore
- [x] Fix `go-release.yml` — hardcoded `go-version: '1.24'` → `go-version-file: '../.go-version'`

## Done (owner, post-session)

- [x] Fix CI lint (`RunGuided` cognitive complexity)
- [x] Merge PR #46, release v0.9.1
- [x] One-liner installer (`install.sh`)
- [x] Sync adapter skills (OpenCode/Cursor alignment)
- [x] GitHub description and topics
- [x] Repo cleanup per audit

## Owner action (remaining)

- [ ] Close stale dependabot PRs: #36, #37, #38 (Python packages)
- [ ] GitHub Settings → Social preview image (1280x640, protocol flow diagram)

---

# Part 1: Code & Infrastructure

## P0: Evidence layer dogfooding

- [ ] Enable evidence emission in own workflow — every `@build` emits `generation` event, every `@review` emits `verification` event
- [ ] Commit `.sdp/log/events.jsonl` to repo (or document why it's excluded). Currently 0 events — the differentiator feature has no data in its own repo
- [ ] `sdp evidence export --format=in-toto` — bridge to OpenClaw / SLSA ecosystem

## P1: Complexity / lint truth

- [ ] Resolve the three-number problem: `guard-rules.yml` says `max_cc: 10`, `.golangci.yml` says `gocyclo: 15` + `gocognit: 20`. Pick one source of truth
- [ ] Reduce `.golangci.yml` exclusions — currently 25 rules disabling gocognit/gocyclo for 15 packages. Either fix the code or raise the limit honestly
- [ ] LOC gate: change CI from warning to blocker, or raise the limit with documented exceptions

## P2: Code hygiene

- [ ] `src/sdp/` vs `sdp-plugin/internal/` — not duplicate (different layers: domain vs application), but not connected. Domain layer (`src/sdp/graph.Dispatcher`, `synthesis.Synthesizer`, `runtime.CircuitBreaker`) is unused by CLI. Either: (A) wire via `replace` directive in `sdp-plugin/go.mod`, or (B) document as future `sdp-orchestrator` repo
- [ ] Qualify "4.96x parallel speedup" — add "(synthetic benchmark, mock executor)" or remove
- [ ] `guard-rules.yml` references Python patterns (`except:`, `except Exception:`, `radon`, `mypy`, `mvn`) — clean up to Go reality
- [ ] "Already Built" section in ROADMAP references Python tools (mypy, Semgrep, `except: pass`) — update to Go equivalents (golangci-lint, gocognit, go vet)

## P3: Release pipeline hardening

- [ ] Create `fall-out-bug/homebrew-tap` repo (or remove brew instructions from README/GoReleaser)
- [ ] Add release smoke-check: `./dist/sdp doctor && ./dist/sdp skill list` after GoReleaser build
- [ ] Add fail-fast for missing GPG_PRIVATE_KEY / PASSPHRASE secrets — human-readable error instead of silent failure

## P4: Agent system evolution

- [ ] Cross-model review: different model/provider for implement vs review agents. Single-LLM "cross-review" is self-review with different prompts
- [ ] Contract tests for skills: Go tests that parse SKILL.md frontmatter and verify consistency (tools exist, agent files exist, version format valid)
- [ ] Agent Capability Registry: machine-readable YAML in skill frontmatter — `allowed_tools`, `expected_artifacts`, `slo_seconds`
- [ ] Fix next-step engine: `sdp next` returns "sdp status" at 50% confidence regardless of state. Should read backlog, check dependencies, recommend concrete workstream

---

# Part 2: Roadmap Recommendations

## Roadmap hygiene (fix now)

- [ ] Mark M1 as "Status: **COMPLETE** (Feb 2026)" in ROADMAP.md — currently reads as future work
- [ ] Update Feature → Milestone Map (lines 601-608): all M1 features should be "Done", not "Backlog"
- [ ] Fix dead links: `SPECTRUM.md` and `LAYERED-ADOPTION.md` are referenced but don't exist
- [ ] Remove Python references from "Already Built" section: `mypy --strict`, `Semgrep security patterns`, `no except: pass` → replace with Go equivalents
- [ ] "159 workstreams" total — separate "planned with WS" from "idea without estimate (TBD)"

## M2 "Smart Casual" — trim to realistic scope

Current M2 has 22+ planned WS + 7 features with no WS estimates. That's 3-4 months, not "4-6 weeks".

### Keep (core M2):

| Feature | WS | Why |
|---|---|---|
| Schema Consolidation | 3 | Foundation for interop. Without this, nothing else works |
| Scope Collision Detection | 2 | `internal/collision/` partly exists. Real value for parallel features |
| F073: Trust & Explainability | 3 | Evidence phase 2, direct path to OpenClaw |
| Discovery → Delivery Contract (merge Agentic Discovery + Bridge) | 3 | One pipeline, not two features. Owner intake → handoff → delivery |
| Explore/Commit Path Policy | 1 | Spec + routing config, not a code feature |
| Drift detection upgrade | 1 | Expand existing `sdp drift detect`, not a new feature |
| F078: Wire existing runtime | 2 | Connect `src/sdp/runtime/` (circuit breaker, retry, degraded) to CLI via `replace` |
| **Total** | **~15 WS** | **4-6 weeks realistic** |

### Move out of M2:

| Feature | Where | Why |
|---|---|---|
| F060: Shared Contracts | → M3 | Enterprise feature, no teams using parallel features yet |
| F071: Team UX & Collaboration | → M3 | Team UX without team = airport for a village |
| F077: Runtime Hooks Platform (full) | → M3 | 5 WS overkill. Keep 1 hook point (`post-build`) in M2, full platform later |
| Architecture Awareness | → Skill update | Not a feature. Wire @reality → @design in existing skills |
| Review as Protocol Self-Check | → 1 WS in M2 | Expand `sdp drift detect`, not a standalone feature |

## M3 "Blazer" — MCP decision point

### Keep:
| Feature | Why |
|---|---|
| F074: Layered OSS Packaging (`sdp init` L0/L1/L2) | Distribution story, critical for adoption |
| F058: CI/CD GitHub Action | Adoption accelerator — one YAML line gives evidence check in PR |
| F079: Model Routing & Economics | Enables cross-model review, answers "single-LLM self-review" critique |
| F069: Next-Step Engine (fix) | Already started, currently broken (50% confidence fallback). Finish it |

### Decide:
| Feature | Question |
|---|---|
| MCP Server | You said MCP unclear. Option: build "Runtime-Agnostic Agent Adapter" as interface, MCP as one implementation. Don't bet entire M3 on MCP |
| F057: CLI plan/apply/log | 45 CLI commands already exist. What's actually new here? Clarify scope or merge into existing |
| F072: Interop & Migration | Import/export from which tools? Without concrete target, this is abstract. Defer until user request |

## M4 "Tie" — direction correct, trim scope

### Keep:
| Feature | Why |
|---|---|
| F055: Compliance Design Doc | EU AI Act mapping. Can write earlier than M4 — it's a document, not code |
| F056: Full Skills Instrumentation | All 19 skills emit evidence. Logical completion |
| Honest hash chain labeling | "Corruption detection, NOT tamper-proof" — mature positioning |

### Move to Horizon:
| Feature | Why |
|---|---|
| F059: Observability Bridge | OTel span attributes for teams with Datadog/Honeycomb. Zero current user demand |
| F061: Data Collection & Benchmark | MTTR tracking, AI failure taxonomy. Research, not product |

## Kill criteria — add tracking

Current kill criteria are well-written but not tracked:

- [ ] Add counter: "builds with evidence" (target: 100 before evaluating)
- [ ] Add counter: "times sdp trace used during real incident" (kill if 0 after 100 builds)
- [ ] Add counter: "external users" (kill if 0 after 6 months = Aug 2026)
- [ ] Add counter: "acceptance test pass rate vs baseline" (kill if no measurable improvement after 50 builds)

## New skill: @audit

Missing skill that this audit session exposed. Different from `@reality` (analyzes code) and `@review` (checks feature quality).

`@audit` analyzes **project as product**:
1. Build health (compile, test, coverage)
2. Distribution (releases, install paths, README accuracy)
3. Docs freshness (dead links, stale technology references, Python/Go mismatch)
4. Dead code (unused packages, tracked artifacts, orphaned files)
5. Lint truth (are rules enforced or excluded everywhere?)
6. Evidence health (log populated? hash chain intact? events emitted?)
7. Dependency health (outdated deps, security advisories, stale dependabot)
8. Community health (stars, forks, open issues, PR response time)

---

# Part 3: Promotion

## Quick wins (1 evening):
- [ ] README badges (Go version, CI status, license, latest release)
- [ ] Social preview image for link sharing
- [ ] Enable GitHub Discussions (pin: "What is SDP?", "Show your setup", "Roadmap feedback")

## Content (1 per week):
- [ ] Habr article: "Как я построил протокол для AI-агентов за 50 дней" (story format)
- [ ] dev.to article: English version
- [ ] Reddit posts: r/ClaudeAI, r/cursor, r/ChatGPTCoding — experience report, not promo
- [ ] Submit to awesome-lists (awesome-claude, awesome-ai-coding, awesome-developer-tools)

## Ongoing:
- [ ] Terminal demo GIF in README (Charm VHS) — 30-second flow
- [ ] GitHub template repo `fall-out-bug/sdp-starter` — L0 via "Use this template"
- [ ] Bi-weekly stream clips: 30-min build session → 5-10 short clips
- [ ] Propose SDP as OpenCode community recipe/template

---

# Not now (with reasons)

| Item | Reason |
|---|---|
| Homebrew tap | L0 distribution via installer is sufficient |
| npm package | Overkill, requires Node |
| MCP server | Market signal unclear, revisit after M2 |
| Product Hunt / Hacker News | Wait for v1.0 + GIF + working 5-min quickstart |
| Discord server | Not enough users, GitHub Discussions sufficient |
| IDE plugins (Cursor/VS Code native) | MCP or adapter interface first |
| Full observability bridge (OTel) | No user demand |

---

*Generated: 2026-02-16. Sources: Cloud Agent audit, Codex review, roadmap analysis v8.1.*
