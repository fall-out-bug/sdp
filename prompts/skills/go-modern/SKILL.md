---
name: go-modern
description: Apply modern Go idioms that match the repo's Go version before writing or reviewing Go code.
---

# @go-modern

Use this skill when editing or reviewing Go in `sdp`.

The repo targets modern Go. Prefer standard-library idioms that are already available in the checked-in Go version.

## Core Rule

Before changing Go code:

1. Detect the module Go version from `go.mod` or `sdp-plugin/go.mod`.
2. Use only language and stdlib features supported by that version.
3. Prefer standard-library modernizations that reduce custom code without changing behavior.

## Prefer These Patterns

- `slices.SortFunc` over `sort.Slice`
- `slices.Contains` over manual membership loops
- `maps.Copy` or `maps.Clone` over manual copy loops
- `strings.Cut` over `strings.SplitN(..., 2)` or `strings.Index` plus slicing
- `strings.CutPrefix` and `strings.CutSuffix` over `HasPrefix` or `HasSuffix` plus trim
- `min(...)` and `max(...)` over trivial compare-and-assign blocks
- `any` over `interface{}` for parameters, locals, and variadics when behavior does not change

## First-Wave Safe Changes

Use these freely in polish PRs when tests stay green:

- mechanical helper swaps (`Cut`, `CutPrefix`, `SortFunc`, `Contains`, `Copy`, `Clone`)
- removing tiny custom helpers now replaced by stdlib
- local `interface{}` to `any` updates
- comments or docs that replace outdated Go guidance such as `golint`

## Avoid In The First Pass

Do not batch risky modernizations into a style PR:

- context plumbing solely to replace `exec.Command` with `exec.CommandContext`
- wide JSON tag changes such as `omitempty` to `omitzero`
- concurrency rewrites (`wg.Go`, atomics, cancellation semantics) unless the feature already needs them
- broad repo-wide `interface{}` churn across public structs or serialized payloads

## Review Checklist

When reviewing Go changes, ask:

1. Is there a standard-library helper that makes this code shorter and clearer?
2. Is the rewrite behavior-preserving and easy to verify?
3. Is the new idiom supported by the module's Go version?
4. Does the change remove custom parsing, sorting, or copying logic instead of adding abstraction?

## Verification

After modernization changes:

- run package tests for touched areas
- run `go build ./...`, `go test ./...`, and `go vet ./...`
- run `golangci-lint run ./...` where configured
- use `golangci-lint run --enable-only modernize --issues-exit-code 0 ./...` for audit snapshots

## Recovery

| Symptom | Fix |
|---------|-----|
| Skill produces no output | Check working directory is project root with `docs/workstreams/backlog/` |
| "checkpoint not found" | Run `sdp-orchestrate --feature <ID>` to create initial checkpoint |
| "workstream files missing" | Run `sdp-orchestrate --index` to verify, then `@feature` to regenerate |
| Skill hangs / no progress | Check `.sdp/log/events.jsonl` for last event; use `sdp reset --feature <ID>` if stuck |
| Review loop exceeds 3 rounds | Use `@review --override "reason"`, `@review --partial`, or `@review --escalate` |

## See Also

- `@build` - implementation workflow
- `@review` - multi-agent review
- `CONTRIBUTING.md` - repository code style
- `DEVELOPMENT.md` - local verification flow
