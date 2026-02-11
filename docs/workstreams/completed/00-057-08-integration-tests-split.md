---
ws_id: "00-057-08"
feature_id: "F057"
title: "Split integration tests to <=200 LOC per file"
status: completed
priority: "P1"
depends_on: ["00-057-04"]
blocks: ["00-057-09"]
project_id: "00"
completed_at: "2026-02-10"
---

# 00-057-08: Split integration tests to <=200 LOC per file

## Goal

Split integration test file so every test file is <=200 LOC.

## Acceptance Criteria

- [x] AC1–AC4: All implemented (tests split into 13 files).

## Execution

Merged in PR #17 (F057). Deliverables: cmd/sdp/*_integration_test.go (13 files).
