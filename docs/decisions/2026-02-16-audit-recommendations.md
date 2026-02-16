# Audit Recommendations — 2026-02-16

Consolidated from: repo audit (Cloud Agent), Codex review, owner feedback.

Status: **PLAN** — items to be scheduled into workstreams.

---

## Done (this session)

- [x] Remove `uv.lock` — was triggering Python dependabot PRs
- [x] Remove tracked `sdp-plugin/test_output.txt` (73KB) + add to .gitignore
- [x] Fix `go-release.yml` — hardcoded `go-version: '1.24'` → `go-version-file: '../.go-version'`

## Owner action (no agent permissions)

- [ ] Close stale dependabot PRs: #36, #37, #38 (Python packages, no longer relevant)
- [ ] GitHub Settings → Description: `Structured protocol for AI-assisted development. Discovery → Delivery → Evidence.`
- [ ] GitHub Settings → Topics: `ai-agents`, `developer-tools`, `code-quality`, `tdd`, `provenance`, `evidence`, `orchestration`, `claude-code`, `cursor`, `opencode`
- [ ] GitHub Settings → Social preview image (1280x640, protocol flow diagram)

---

## P0: Unblock v0.9.0 release

- [ ] Fix `RunGuided` cognitive complexity (gocognit 22 > 20) — split into 2 functions, unblocks CI lint
- [ ] Merge PR #46 "Release v0.9.0"
- [ ] Tag v0.9.0, trigger GoReleaser

## P1: L0 standalone distribution

- [ ] Create `scripts/install-prompts.sh` — curl + tar, copies `prompts/`, `CLAUDE.md`, creates symlinks. No Go, no Node
- [ ] Or: create GitHub template repo `fall-out-bug/sdp-starter` with just L0 (prompts + agents + CLAUDE.md + example workstream). "Use this template" button
- [ ] README section: "Install prompts only (no CLI)" with one-liner
- [ ] Terminal demo GIF in README (Charm VHS or asciinema) — 30-second `@feature → @oneshot → sdp trace` flow

## P2: Evidence layer dogfooding

- [ ] Enable evidence emission in own workflow — every `@build` emits `generation` event, every `@review` emits `verification` event
- [ ] Commit `.sdp/log/events.jsonl` to repo (or document why not). Currently 0 events — the differentiator feature has no data in its own repo
- [ ] `sdp evidence export --format=in-toto` — bridge to OpenClaw / SLSA ecosystem. One format, one command

## P3: Complexity / lint truth

- [ ] Resolve the three-number problem: `guard-rules.yml` says `max_cc: 10`, `.golangci.yml` says `gocyclo: 15` + `gocognit: 20`. Pick one source of truth
- [ ] Reduce `.golangci.yml` exclusions — currently 25 exclusion rules disabling gocognit/gocyclo for 15 packages. Either fix the code or raise the limit honestly
- [ ] LOC gate: change CI from warning to blocker, or raise the limit with documented exceptions. "For now, this is a warning" has been "for now" long enough

## P4: Code hygiene

- [ ] `src/sdp/` vs `sdp-plugin/internal/` — not duplicate (different layers), but not connected. Domain layer (`src/sdp/graph.Dispatcher`, `synthesis.Synthesizer`, `runtime.CircuitBreaker`) is unused by CLI. Either: (A) wire via `replace` directive in `sdp-plugin/go.mod`, or (B) document as future `sdp-orchestrator` repo and leave intentionally separate
- [ ] Qualify "4.96x parallel speedup" in PRODUCT_VISION.md and ROADMAP — add "(synthetic benchmark, mock executor)" or remove the number
- [ ] Adapter path drift: `.cursor/commands/` has 9 of 27 skills, `.opencode/commands/` has 23 of 27. Document which skills are available in which adapter, or sync them

## P5: Release pipeline hardening

- [ ] Create `fall-out-bug/homebrew-tap` repo (or remove brew instructions from README)
- [ ] Add release smoke-check step: `./dist/sdp doctor && ./dist/sdp skill list` after GoReleaser build
- [ ] Add fail-fast for missing GPG_PRIVATE_KEY / PASSPHRASE secrets in go-release.yml — human-readable error instead of silent failure
- [ ] `guard-rules.yml` references Python patterns (`except:`, `except Exception:`, `radon`, `mypy`, `mvn`) — clean up to Go reality

## P6: Agent system evolution

- [ ] Cross-model review: implement different model/provider for implement vs review agents. Current single-LLM "cross-review" is self-review with different prompts
- [ ] Contract tests for skills: Go tests that parse SKILL.md frontmatter and verify consistency (expected tools exist, referenced agent files exist, version format valid)
- [ ] Agent Capability Registry: machine-readable YAML in skill frontmatter — `allowed_tools`, `expected_artifacts`, `slo_seconds`. Makes skills testable and documentable
- [ ] Fix next-step engine: `sdp next` returns "sdp status" at 50% confidence regardless of state. Should read backlog, check dependencies, recommend concrete workstream

## P7: Promotion

- [ ] README badges (Go version, CI status, license, latest release)
- [ ] Habr article: "Как я построил протокол для AI-агентов за 50 дней" (story, not docs)
- [ ] dev.to article: English version of the same
- [ ] Reddit posts: r/ClaudeAI, r/cursor, r/ChatGPTCoding — experience report format
- [ ] Submit to awesome-lists (awesome-claude, awesome-ai-coding, awesome-developer-tools)
- [ ] Enable GitHub Discussions (pin: "What is SDP?", "Show your setup", "Roadmap feedback")
- [ ] Bi-weekly stream clips: record 30-min "Build X with SDP", cut 5-10 short clips for social
- [ ] Propose SDP as OpenCode community recipe/template

---

## Not now

- Homebrew tap — L0 distribution first
- npm package — overkill, requires Node
- MCP server — unclear market signal, revisit Q2
- Product Hunt / Hacker News — wait for v1.0 + working install + GIF + 5-min quickstart
- Discord server — not enough users yet, GitHub Discussions sufficient
