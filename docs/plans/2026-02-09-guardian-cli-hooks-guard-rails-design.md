# Guardian-cli Hooks and Guard Rails Adoption

> **Status:** Research complete
> **Date:** 2026-02-09
> **Goal:** Identify and adapt guardian-cli mechanisms to strengthen SDP hooks and guard rails.

---

## Table of Contents

1. [Overview](#overview)
2. [Guardian-cli Hook System](#1-guardian-cli-hook-system)
3. [Guardian-cli Guard-Rail Engine](#2-guardian-cli-guard-rail-engine)
4. [SDP Hooks and Guard Rails Baseline](#3-sdp-hooks-and-guard-rails-baseline)
5. [Integration Strategy](#4-integration-strategy)
6. [Implementation Plan](#implementation-plan)

---

## Overview

### Source of inspiration

- [guardian-cli README](https://github.com/AlexGladkov/guardian-cli/blob/main/README.md)
- [guardian-cli hooks manager](https://github.com/AlexGladkov/guardian-cli/blob/main/internal/git/hooks.go)
- [guardian-cli check command](https://github.com/AlexGladkov/guardian-cli/blob/main/internal/cli/check.go)
- [guardian-cli rule engine](https://github.com/AlexGladkov/guardian-cli/blob/main/internal/engine/checker.go)
- [guardian-cli exceptions and meta-check](https://github.com/AlexGladkov/guardian-cli/blob/main/internal/engine/meta_check.go)

### Goals

1. **Reduce drift** between hook scripts, CLI behavior, and documentation.
2. **Improve guard rail UX** with clear outputs, safe hook ownership, and low friction.
3. **Enable staged enforcement** (local hooks + CI) using consistent rules and exit codes.

### Key Decisions

| Aspect | Decision |
| --- | --- |
| Hook system | Port guardian-style hook ownership (marker + idempotent) and keep non-blocking notifications; add optional enforcement hooks. |
| Guard-rail engine | Adopt rule registry + diff-based staged scan + JSON/human output; add CI diff-range detection. |
| SDP baseline | Consolidate on Go CLI enforcement and align docs/configs to the Go-first path. |
| Integration | Phase in pre-commit staged guard scan first; unify hook installer template next; governance later. |

---

## 1. Guardian-cli Hook System

> **Experts:** Kelsey Hightower, Sam Newman, Troy Hunt

### Adoptable patterns

Guardian-cli installs only the hooks it owns, marks them, and keeps them non-blocking. This minimizes friction and avoids clobbering existing hook managers. The same pattern can be used in SDP to safely install hook templates without overwriting user hooks.

| Aspect | Details |
| --- | --- |
| Ownership | Marker comment to detect guardian-managed hooks, enabling safe uninstall. |
| Idempotency | Install is repeatable without duplicating content. |
| UX | Background execution with quiet output to avoid blocking workflows. |
| State | "Since last check" state reduces notification noise. |

### Example pattern (rule of thumb)

```sh
#!/bin/sh
# SDP-MANAGED-HOOK
nohup sdp notify --since-last-check --quiet >/dev/null 2>&1 &
```

---

## 2. Guardian-cli Guard-Rail Engine

> **Experts:** Sam Newman, Theo Browne, Troy Hunt

### Adoptable patterns

Guardian-cli uses a rule registry with YAML rules, diff-based checks on added lines, and clear exit codes. Exceptions support expiry and path globs. CI diff-range detection reduces false positives.

| Aspect | Details |
| --- | --- |
| Rule registry | Rule type dispatch by `type` enables extensibility. |
| Diff-based checks | Evaluates added lines for fast local and CI runs. |
| Exceptions | `rule_id` + path globs + expiry for temporary waivers. |
| Governance | Meta-check blocks unauthorized changes to policy files. |
| Outputs | Human and JSON output with stable exit codes (0/1/2). |

### Example rule shape

```yaml
rules:
  - id: domain_no_infra
    description: Domain must not depend on infra
    type: imports_forbidden
    config:
      from_globs: ["domain/**"]
      forbid_globs: ["infra/**"]
    severity: error
```

---

## 3. SDP Hooks and Guard Rails Baseline

> **Experts:** Kelsey Hightower, Martin Fowler, Theo Browne

### Current state and gaps

| Area | Current | Gap |
| --- | --- | --- |
| Hook entry points | Repo `hooks/*.sh` + `sdp hooks install` template | Two systems drift and are not aligned. |
| Pre-commit/pre-push | Scripts expect Python modules or `sdp pre-commit` | Missing Python modules and CLI subcommands. |
| Guard scope | `sdp guard` exists | Scope files are not populated; enforcement is weak. |
| Config alignment | `quality-gate.toml` + `.sdp/config.yml` | CLI ignores these; docs/CI differ. |
| Evidence | CLI emits evidence | Hooks do not emit evidence events. |

---

## 4. Integration Strategy

> **Experts:** Sam Newman, Kelsey Hightower, Troy Hunt

### Recommended path

1. **Pre-commit staged guard scan (MVP).** Add a Go subcommand (e.g., `sdp guard check --staged`) that scans staged files and fails on violations. Wire it into `hooks/pre-commit.sh`.
2. **Unify hook installer.** Update `sdp hooks install` to install the repo hook scripts or a single canonical template using a marker for safe ownership.
3. **Config and output parity.** Map `quality-gate.toml` and `.sdp/config.yml` into guard/quality outputs; add JSON output mode for CI.
4. **Governance and exceptions.** Add expiring exceptions and a meta-check for rule file changes once enforcement is stable.

---

## Implementation Plan

### Phase 1: MVP

- [ ] Add `sdp guard check --staged` with clear exit codes and human output.
- [ ] Update `hooks/pre-commit.sh` to call the staged guard scan.
- [ ] Emit evidence events for guard results (local and CI).
- [ ] Document the canonical hook path and the required CLI commands.

### Phase 2: Hardening

- [ ] Update `sdp hooks install` to use marker-based ownership and install the canonical template.
- [ ] Add CI diff-range detection and JSON output for guard checks.
- [ ] Align `quality-gate.toml` thresholds with CLI outputs and CI gates.

### Phase 3: Governance

- [ ] Introduce rule registry config (YAML) and exception expiry.
- [ ] Add a meta-check for policy file changes requiring approval.
- [ ] Provide migration notes from legacy hooks to the unified path.

---

## Success Metrics

| Metric | Baseline | Target |
| --- | --- | --- |
| Developers using canonical hooks | Unknown | 90% in 30 days |
| Guard checks run locally | Inconsistent | 80% of commits |
| Hook-related failures in CI | High drift | Reduced by 50% |
| Config/doc drift incidents | Frequent | Near-zero |
