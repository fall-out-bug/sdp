---
ws_id: "00-057-03"
feature_id: "F057"
title: "sdp log show & Filters"
status: completed
priority: "P1"
depends_on: ["00-054-07"]
blocks: ["00-057-04"]
project_id: "00"
completed_at: "2026-02-10"
---

# 00-057-03: sdp log show & Filters

## Goal

Extend `sdp log` with interactive browser, filters, and export.

## Acceptance Criteria

- [x] AC1–AC8: All implemented (log show, filters type/model/date/ws, export csv/json, log stats).

## Execution

Merged in PR #17 (F057). Deliverables: sdp-plugin/cmd/sdp/log*.go, internal/evidence/browser.go, export.go.
