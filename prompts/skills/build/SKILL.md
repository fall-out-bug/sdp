---
name: build
description: Execute workstream with TDD and guard enforcement
cli: sdp apply --ws
version: 6.0.0
---

# build

> **CLI:** `sdp apply --ws <workstream-id>`

Execute a single workstream following TDD discipline with guard enforcement.

---

## EXECUTE THIS NOW

When user invokes `@build 00-067-01`, run:

```bash
sdp apply --ws 00-067-01
```

This is not an example. Execute this command.

---

## Quick Start

```bash
sdp apply --ws 00-067-01         # Execute workstream
sdp apply --ws 00-067-01 --retry 3  # Allow 3 retries
sdp apply                         # Execute all ready workstreams
```

---

## What Happens

```
Step 1: Resolve Task ID
    └─ sdp resolve {ID} (detects ws_id, beads_id, issue_id)

Step 2: Activate Guard
    └─ sdp guard activate {WS-ID}

Step 3: Read Workstream Spec
    └─ docs/workstreams/backlog/{WS-ID}.md

Step 4: TDD Cycle
    └─ RED: Write failing test
    └─ GREEN: Implement minimum code
    └─ REFACTOR: Clean up

Step 5: Quality Gates
    └─ Test coverage >= 80%
    └─ LOC <= 200 per file
    └─ Lint passes

Step 6: Commit
    └─ git commit with workstream reference
```

---

## Identifier Formats

Accepts any format:

```bash
@build 00-067-01      # Workstream ID (PP-FFF-SS)
@build 99-F064-01     # Fix workstream (99-{FEATURE}-{SEQ})
@build sdp-xxx        # Beads task ID (resolved)
@build ISSUE-0001     # Issue ID (resolved)
```

---

## Verbosity

```bash
@build 00-067-01 --quiet    # Exit status only: ✅
@build 00-067-01            # Summary
@build 00-067-01 --verbose  # Step-by-step
@build 00-067-01 --debug    # Internal state
```

---

## Quality Gates

| Gate | Threshold | Check Command |
|------|-----------|---------------|
| Tests | 100% pass | `go test ./...` |
| Coverage | >= 80% | `go test -cover ./...` |
| Lint | 0 errors | `golangci-lint run` |
| File Size | <= 200 LOC | `wc -l *.go` |

**LOC check (MANDATORY):**
```bash
for file in *.go; do
  loc=$(wc -l < "$file")
  if [ "$loc" -gt 200 ]; then
    echo "ERROR: $file is $loc LOC (max: 200)"
    exit 1
  fi
done
```

---

## Beads Integration

When Beads is enabled (`bd --version` works and `.beads/` exists):

1. **Before work:** `bd update {beads_id} --status in_progress`
2. **On success:** `bd close {beads_id} --reason "WS completed"`
3. **On failure:** `bd update {beads_id} --status blocked`

---

## Git Safety

**CRITICAL:** Features MUST be in feature branches. Never commit to dev or main.

Before any git command:
```bash
pwd
git branch --show-current
sdp guard context check
```

---

## Errors

| Error | Fix |
|-------|-----|
| No active WS | `sdp guard activate {WS-ID}` |
| File not in scope | Check WS scope_files |
| Coverage <80% | Add tests |
| LOC >200 | Split file |

---

## See Also

- `.claude/patterns/tdd.md` - TDD pattern
- `.claude/patterns/quality-gates.md` - Quality gates
- `@oneshot` - Execute all workstreams

**Implementation:** `sdp-plugin/cmd/sdp/apply.go`
