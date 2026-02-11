---
ws_id: "00-057-02"
feature_id: "F057"
title: "sdp apply Command"
status: completed
priority: "P1"
depends_on: ["00-054-05"]
blocks: ["00-057-04"]
project_id: "00"
completed_at: "2026-02-10"
---

# 00-057-02: sdp apply Command

## Goal

Implement `sdp apply` — execute workstreams from the terminal with streaming progress.

## Acceptance Criteria

- [x] AC1–AC8: All implemented (execute ready/specific, retry, dry-run, progress, JSON, dependency order, evidence chain).

## Execution

Merged in PR #17 (F057). Deliverables: sdp-plugin/cmd/sdp/apply.go, internal/executor/.
