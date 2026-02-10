---
ws_id: 00-054-02
feature: F054
status: pending
size: MEDIUM
project_id: "00"
parent: 00-054-01
---

## Goal
Implement `sdp apply` command with workstream execution.

## Acceptance Criteria
- [x] AC1: Execute all ready workstreams (no blockers)
- [x] AC2: Execute specific workstream with --ws flag
- [x] AC3: Retry failed workstreams with --retry flag
- [x] AC4: Dry-run mode shows execution plan
- [x] AC5: Streaming progress with visual progress bar
- [x] AC6: JSON output format for CI/CD
- [x] AC7: Respect dependency order
- [x] AC8: Emit full evidence chain

## Dependencies
- 00-054-01 (base setup must be complete)

## Scope Files
**Implementation:**
- cmd/sdp/apply.go (NEW)
- internal/executor/executor.go (ENHANCE)
- internal/executor/progress.go (ENHANCE)

**Tests:**
- internal/executor/executor_test.go (ENHANCE)
