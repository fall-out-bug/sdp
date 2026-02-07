#!/usr/bin/env python3
"""Regenerate workstream files from roadmap document."""

import re
from pathlib import Path

# Read roadmap
roadmap_path = Path("docs/plans/2026-02-05-golang-migration-roadmap.md")
roadmap_content = roadmap_path.read_text()

# Extract workstream sections
ws_pattern = r'### (00-050-\d+): ([^\n]+)\n\*\*Size:\*\* ([^\n]+)\n\*\*Priority:\*\* ([^\n]+)\n\*\*Depends on:\*\* ([^\n]+)\n\n\*\*Goal:\*\* ([^\n]+)\n\n\*\*Scope Files:\*\*\n((?:- `[^\n]+`\n)+)\n\n\*\*Acceptance Criteria:\*\*\n((?:- AC\d+: [^\n]+\n)+)'

# For now, let's manually recreate each file based on the template format
# Using a simpler approach: read the existing template and fill in details

template = """---
ws_id: {ws_id}
parent: sdp-79u
feature: F050
status: backlog
size: {size}
project_id: 00
---

## WS-{ws_id}: {title}

### 🎯 Goal

{goal}

**What must WORK after completing this WS:**
{outcomes}

**Acceptance Criteria:**
{acceptance_criteria}

**⚠️ WS is NOT complete until Goal is achieved (all AC ✅).**

---

### 📋 Contracts

{contracts}

---

### 🔌 Interfaces

{interfaces}

---

### 📁 Scope Files

{scope_files}

---

### ✅ Testing

{testing}

---

### 🔍 Verification

{verification}

---

### 📝 Notes

**Dependencies:** {dependencies}

**Context:**
{context}

**Key Decisions:**
- TBD

**Risks:**
- TBD

**Mitigations:**
- TBD

"""

# Workstream data extracted from roadmap and system reminders
workstreams = {
    "00-050-02": {
        "title": "TDD Runner Implementation",
        "size": "MEDIUM",
        "goal": "Implement Red-Green-Refactor cycle with pytest integration. Enable automated test execution with phase tracking and error reporting.",
        "outcomes": """- Red phase runs pytest, expects failure
- Green phase runs pytest, expects success
- Refactor phase runs pytest, ensures tests still pass
- Each phase captures stdout/stderr for error messages
- Each phase reports duration (for telemetry)
- RunAllPhases executes Red→Green→Refactor sequentially
- Context support for cancellation (ctx.Cancel)""",
        "acceptance_criteria": """- [ ] AC1: Red phase runs pytest, expects failure (returncode != 0)
- [ ] AC2: Green phase runs pytest, expects success (returncode == 0)
- [ ] AC3: Refactor phase runs pytest, ensures tests still pass
- [ ] AC4: RunAllPhases executes all 3 phases sequentially
- [ ] AC5: Each phase reports duration (time.Since)
- [ ] AC6: Context cancellation stops execution immediately
- [ ] Coverage ≥ 80%""",
        "contracts": """**Technical Constraints:**
- Must use os/exec for pytest execution
- Must capture stdout/stderr separately
- Must support context.Context for cancellation
- Must return structured errors (phase + cause)""",
        "interfaces": """```go
package tdd

type Phase int

const (
    Red Phase = iota
    Green
    Refactor
)

type Runner struct {
    pytestPath string
}

func (r *Runner) RunPhase(ctx context.Context, phase Phase, wsPath string) (*PhaseResult, error)
func (r *Runner) RunAllPhases(ctx context.Context, wsPath string) ([]*PhaseResult, error)

type PhaseResult struct {
    Phase    Phase
    Success  bool
    Duration time.Duration
    Stdout   string
    Stderr   string
    Error    error
}
```""",
        "scope_files": """**Implementation:**
- `internal/tdd/runner.go`
- `internal/tdd/phase.go`
- `cmd/sdp/commands/build.go`

**Tests:**
- `internal/tdd/runner_test.go`""",
        "testing": """**Unit Tests:**
```bash
# Red phase execution
TestRedPhaseFails()

# Green phase execution
TestGreenPhaseSucceeds()

# Refactor phase execution
TestRefactorPhasePasses()

# Context cancellation
TestContextCancellation()
```

**Integration Tests:**
```bash
# Full TDD cycle
go test ./internal/tdd/... -run TestIntegrationTDDCycle
```""",
        "verification": """**Pre-build:**
- [ ] Review pytest command construction
- [ ] Verify context cancellation logic
- [ ] Check error propagation

**Post-build:**
```bash
# Run tests
go test ./internal/tdd/... -v -cover

# Manual integration test
cd /tmp/test-project
../../sdp build 00-050-02

# Expected: Executes Red → Green → Refactor
```""",
        "dependencies": "00-050-01",
        "context": """Simplified from Python (800 LOC → 400 LOC). TDD runner is core to SDP workflow. Context cancellation enables @oneshot to stop mid-execution."""
    },
}

# For now, let's just recreate 00-050-02.md as a test
ws = workstreams["00-050-02"]
content = template.format(
    ws_id="00-050-02",
    title=ws["title"],
    size=ws["size"],
    goal=ws["goal"],
    outcomes=ws["outcomes"],
    acceptance_criteria=ws["acceptance_criteria"],
    contracts=ws["contracts"],
    interfaces=ws["interfaces"],
    scope_files=ws["scope_files"],
    testing=ws["testing"],
    verification=ws["verification"],
    dependencies=ws["dependencies"],
    context=ws["context"],
)

Path("docs/workstreams/backlog/00-050-02.md").write_text(content)
print("Regenerated 00-050-02.md")
