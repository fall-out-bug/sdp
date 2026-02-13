# Repository Analysis: SDP (Spec-Driven Protocol)

**Date:** 2026-02-13
**Scope:** Code quality, distribution, agent configuration, CI/CD, developer experience

---

## Executive Summary

SDP is an ambitious multi-agent orchestration framework for AI-assisted development. The codebase is substantial (~80K LOC Go, 465 files, 203 test files) with a well-thought-out architecture. However, several systemic issues reduce maintainability, hinder adoption, and create inconsistencies between documentation and reality.

---

## Strengths

### 1. Architecture and Design

- **Clean separation of concerns**: `sdp-plugin/internal/` has 37 packages, each with a clear domain (evidence, drift, guard, executor, etc.).
- **Comprehensive agent/skill system**: 26 skills and 25+ agent definitions in `prompts/` provide a rich vocabulary for AI-assisted workflows.
- **Symlink-based multi-tool support**: `.claude/`, `.cursor/`, `.opencode/` all point to canonical `prompts/` source — single source of truth for all AI tool integrations.
- **Evidence chain with hash-chaining**: The `events.jsonl` system provides auditability and traceability.
- **Parallel dispatcher with DAG**: Kahn's algorithm for workstream execution order is a solid foundation.

### 2. Testing

- **Good test coverage ratio**: 203 test files out of 465 total Go files (~44% of files are tests).
- **Diverse test types**: Unit tests, integration tests, smoke tests, performance tests, and E2E tests are all present.
- **CI runs race detection**: `-race` flag in test pipeline catches data races early.

### 3. CI/CD Pipeline

- **Cross-platform builds**: Linux (amd64/arm64), macOS (amd64/arm64), Windows (amd64).
- **GoReleaser integration**: Automated releases with GPG signing, Homebrew tap, Snapcraft, Docker images.
- **Dogfood workflow**: `sdp-verify-dogfood.yml` uses the project's own verify action on itself.
- **Contract validation**: Dedicated CI job for contract validation between components.
- **Custom GitHub Action**: `.github/actions/verify/` is a reusable composite action for quality gates.

### 4. Documentation

- **Extensive**: 100+ documentation files covering protocol, vision, ADRs, runbooks, compliance, threat model.
- **Architecture Decision Records (ADRs)**: 7+ ADRs documenting key decisions.
- **Incident reports**: Post-mortems for CI issues and branch confusion incidents.
- **Templates**: Workstream, idea, skill, migration, and release-note templates.

### 5. Developer Experience

- **Multiple installation methods**: Submodule, CLI build, standalone clone.
- **Git hooks**: Pre-commit, post-commit, pre-push, commit-msg hooks with safety checks.
- **Shell completion**: Bash/Zsh completion support via cobra.
- **Telemetry with consent**: First-run consent flow for telemetry — privacy-aware.

---

## Weaknesses

### 1. Go Version Inconsistency (Critical)

**Problem:** Multiple Go versions declared across the project:

| Location | Version |
|----------|---------|
| `go.mod` (root) | 1.25.6 |
| `sdp-plugin/go.mod` | 1.26 |
| All CI workflows | 1.26 |
| `Dockerfile` | 1.21-alpine |
| `README.md` badge | 1.24+ |

**Impact:** The Dockerfile uses Go 1.21 but the project requires 1.26 features. The root `go.mod` says 1.25.6. These are future Go versions that don't exist yet (as of Feb 2026, Go 1.24 is latest stable). This means:
- Docker builds will fail or use wrong version
- Contributors can't build without matching Go toolchain
- CI may break when GitHub Actions doesn't have Go 1.26

**Recommendation:** Align all Go version references to a single, currently-released version. Use `.go-version` file as single source of truth.

### 2. Massive LOC Violations (Critical)

**Problem:** The project enforces a 200 LOC per file rule, but 20+ source files violate it:

| File | LOC | Over by |
|------|-----|---------|
| `cmd/sdp/contract.go` | 572 | 372 |
| `agents/synthesis_agent.go` | 560 | 360 |
| `agents/code_analyzer.go` | 543 | 343 |
| `agents/contract_validator.go` | 485 | 285 |
| `doctor/doctor.go` | 468 | 268 |
| `cmd/sdp/decisions.go` | 415 | 215 |
| `orchestrator/feature_coordinator.go` | 403 | 203 |
| `ui/dashboard/app.go` | 394 | 194 |
| ... and 12+ more | | |

**Impact:** The project's own quality gates fail against its own code. This undermines trust in the enforcement system.

**Recommendation:** Prioritize splitting these files. The project's own tooling should pass its own gates. This is the single most important credibility issue.

### 3. Linter Exclusion Sprawl (High)

**Problem:** `.golangci.yml` has 30+ exclusion rules, effectively disabling `errcheck` and `gocognit` for most of the codebase:

```yaml
# These packages have errcheck disabled:
cmd/sdp/, internal/quality/, internal/guard/, internal/memory/,
internal/coordination/, internal/task/, internal/git/,
internal/watcher/, internal/decision/, internal/orchestrator/,
internal/ui/
```

**Impact:** Error handling verification is essentially disabled for the majority of production code. The linter configuration exists but provides false confidence.

**Recommendation:** Gradually remove exclusions. Start with `errcheck` — unchecked errors are the #1 source of silent failures in Go.

### 4. Dual Source Tree (Medium)

**Problem:** Go source code exists in two separate trees:
- `sdp-plugin/internal/` — 324 files (main CLI)
- `src/sdp/` — ~50 files (graph, synthesis, monitoring, agents)

The `sdp-plugin/` has its own `go.mod`, while `src/sdp/` is covered by the root `go.mod`. This creates confusion about which package is canonical.

**Impact:** Contributors don't know where to add code. Tests in `tests/` cover `src/sdp/` but not `sdp-plugin/internal/`. Build commands differ between trees.

**Recommendation:** Consolidate into a single Go module. Either merge `src/sdp/` into `sdp-plugin/internal/` or use Go workspaces (`go.work`).

### 5. Self-Referential Submodule (Medium)

**Problem:** `.gitmodules` contains:
```
[submodule ".sdp/.sdp"]
    path = .sdp/.sdp
    url = https://github.com/fall-out-bug/sdp.git
```

The repo includes itself as a submodule inside `.sdp/.sdp`.

**Impact:** Recursive cloning creates confusion. The submodule is likely there for dogfooding (using SDP within SDP), but this pattern is fragile and can cause Git issues.

**Recommendation:** Remove the self-referential submodule. For dogfooding, use symlinks or direct references to the repo root.

### 6. Missing Releases (Medium)

**Problem:** CHANGELOG documents versions up to 0.10.0, but GitHub releases only go up to v0.7.1 (released 2026-02-02). Tags exist up to v1.0.0-go but no corresponding releases.

**Impact:** Users can't install the latest version. The Homebrew tap, Docker images, and GoReleaser are configured but releases aren't being cut.

**Recommendation:** Create GitHub releases for v0.8.0, v0.9.0, v0.10.0. Automate release creation on tag push (the workflow exists but hasn't been triggered).

### 7. Coverage Threshold Gap (Medium)

**Problem:** Documentation states 80% coverage requirement, but CI only enforces 60%:
```yaml
# go-ci.yml
go tool cover -func=coverage.out | grep total | awk '{if ($3+0 < 60.0) ...'
```

**Impact:** The quality gate claims 80% but enforces 60%. This creates a false sense of security.

**Recommendation:** Align CI threshold to documented 80%, or document the actual 60% threshold.

---

## Agent Configuration Improvements

### 1. OpenCode Configuration is Outdated

The `.opencode/opencode.json` references `prompts/commands/build.md` but the actual path is `prompts/skills/build/SKILL.md`. The agents defined (builder, reviewer, planner, deployer, orchestrator) map to only 5 of the 25+ available agent definitions.

**Recommendation:** Update paths and add missing agent configurations for OpenCode.

### 2. Skill File Sizes

Some skill files are very large (reality: 703 lines, review: 697 lines). These are prompt templates, so size matters for context window consumption.

**Recommendation:** Split large skills into composable sub-skills. Keep each SKILL.md under 300 lines.

### 3. Duplicate Agent Definitions

`sdp-plugin/prompts/agents/` contains a subset (10 files) of the 25 agents in `prompts/agents/`. These appear to be older copies, not symlinks.

**Recommendation:** Remove `sdp-plugin/prompts/agents/` and use symlinks to `prompts/agents/`, or consolidate.

### 4. Agent Versioning

Some agents have version headers (`version: 2.2.0`), others don't. No mechanism to check agent version compatibility.

**Recommendation:** Add version to all agent definitions. Create a version compatibility matrix.

---

## Deployment Improvements

### 1. Dockerfile Uses Ancient Go Version

`Dockerfile` specifies `golang:1.21-alpine` while the project requires 1.26. Docker builds will produce a binary compiled with the wrong Go version (or fail outright).

**Recommendation:** Update Dockerfile to match `go.mod` version. Use build args for version pinning:
```dockerfile
ARG GO_VERSION=1.24
FROM golang:${GO_VERSION}-alpine AS builder
```

### 2. CGO_ENABLED Conflict

The Dockerfile sets `CGO_ENABLED=0` but `sdp-plugin` depends on `modernc.org/sqlite` (which is a pure-Go SQLite, so this actually works). However, this dependency is heavy and may cause issues.

**Recommendation:** Document the SQLite dependency rationale. Consider if the memory/drift features could use a lighter storage backend for CLI-only users.

### 3. No Install Script

README offers `go build` and `brew install` but no `curl | bash` install script. The GoReleaser config mentions `brew install fall-out-bug/tap/sdp` but the tap likely doesn't exist yet (no releases since v0.7.1).

**Recommendation:** Create a standalone install script (similar to Beads' `install.sh`). Publish to Homebrew tap with each release.

### 4. Missing `go install` Support

`PROTOCOL.md` says `go install github.com/fall-out-bug/sdp@latest` but this will install from the root module, not `sdp-plugin/`. The CLI entry point is at `sdp-plugin/cmd/sdp/main.go`.

**Recommendation:** Either restructure so `cmd/sdp/` is at root, or update docs to: `go install github.com/fall-out-bug/sdp/sdp-plugin/cmd/sdp@latest`

---

## Development Workflow Improvements

### 1. No `go.work` File

With two Go modules (`go.mod` at root and `sdp-plugin/go.mod`), there's no Go workspace file to manage them together.

**Recommendation:** Add `go.work`:
```
go 1.24

use (
    .
    ./sdp-plugin
)
```

### 2. Missing Development Setup Documentation

No `DEVELOPMENT.md` or setup guide for contributors. `CONTRIBUTING.md` is high-level but doesn't cover:
- Required Go version
- How to run tests
- How to build the CLI
- How to run linters locally
- IDE setup recommendations

**Recommendation:** Create a `DEVELOPMENT.md` with step-by-step local setup instructions.

### 3. No Dependabot for Go

Dependabot is configured for pip (`dependabot/pip/rich-14.3.2` branch exists) but not for Go modules.

**Recommendation:** Add `.github/dependabot.yml` for Go module updates.

### 4. Quality Gate Config Duplication

Quality rules are defined in three places:
- `.sdp/guard-rules.yml` (12 rules)
- `quality-gate.toml` (10 sections)
- `ci-gates.toml` (3 sections)

These are not synchronized and use different formats (YAML vs TOML).

**Recommendation:** Consolidate into a single configuration file. Use `.sdp/guard-rules.yml` as the canonical source and derive CI config from it.

### 5. Python Artifacts Remain

Files like `quality-gate.toml` reference Python tools (pytest, mypy, ruff, PEP8), and there's a `poetry.lock` at root. The project has pivoted to Go, but Python configuration remains.

**Recommendation:** Remove or clearly separate Python-era artifacts. If Python support is still needed, document it explicitly.

---

## Priority Recommendations

| Priority | Issue | Effort | Impact |
|----------|-------|--------|--------|
| P0 | Fix LOC violations in own codebase | High | Critical for credibility |
| P0 | Align Go versions across all configs | Low | Build reliability |
| P1 | Cut missing releases (v0.8.0-v0.10.0) | Low | User adoption |
| P1 | Fix Dockerfile Go version | Low | Docker build reliability |
| P1 | Align CI coverage threshold to 80% | Low | Quality gate integrity |
| P1 | Remove/reduce linter exclusions | Medium | Code quality |
| P2 | Consolidate dual source tree | High | Developer experience |
| P2 | Add go.work file | Low | Developer experience |
| P2 | Update OpenCode configuration | Low | Multi-tool support |
| P2 | Create DEVELOPMENT.md | Low | Contributor onboarding |
| P2 | Remove Python artifacts | Low | Codebase clarity |
| P3 | Consolidate quality gate configs | Medium | Maintainability |
| P3 | Add Dependabot for Go | Low | Dependency freshness |
| P3 | Remove self-referential submodule | Low | Git hygiene |
| P3 | Split large skill files | Medium | Context window efficiency |

---

## Conclusion

SDP has a strong architectural foundation and an impressive breadth of features. The main risks are:

1. **Credibility gap** — The project's own code violates its quality gates (LOC, coverage threshold mismatch, linter exclusions). This is the single most important issue to address.

2. **Distribution gap** — Users can't easily install the latest version. No releases since v0.7.1 despite significant development.

3. **Complexity creep** — 37 internal packages, 91 command files, and 324 internal source files for a CLI tool suggests potential over-engineering. Consider which features are core vs. optional.

The project would benefit most from a "practice what you preach" sprint: split oversized files, align CI thresholds, cut proper releases, and clean up version inconsistencies.
