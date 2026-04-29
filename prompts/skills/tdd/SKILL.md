---
name: tdd
description: Enforce Test-Driven Development: Red → Green → Refactor (INTERNAL - used by @build)
---

# @tdd (INTERNAL)

TDD discipline. Called by @build, not users.

## Cycle

1. **RED** — Write failing test first. Run test suite (see Quality Gates in AGENTS.md) — must FAIL
2. **GREEN** — Minimal implementation. Run test suite (see Quality Gates in AGENTS.md) — must PASS
3. **REFACTOR** — Improve code. Run test suite (see Quality Gates in AGENTS.md) — still PASS
4. **COMMIT** — Save state

## Exit When

- All AC met
- Test suite passes (see Quality Gates in AGENTS.md)
- Static analysis passes (see Quality Gates in AGENTS.md)

## Example (Go)

```go
// RED: test first
func TestEmailValid(t *testing.T) {
    v := NewValidator()
    if !v.IsValid("a@b.com") { t.Error("expected valid") }
    if v.IsValid("x") { t.Error("expected invalid") }
}
// Run: FAIL (undefined NewValidator)

// GREEN: minimal impl
func NewValidator() *V { return &V{} }
func (v *V) IsValid(s string) bool { return strings.Contains(s, "@") }
// Run: PASS

// REFACTOR: improve, tests still pass
```

For Go refactors, prefer modern stdlib idioms from `@go-modern` when they preserve behavior.

## Write Plan (F101)

Before writing test files or implementation files, emit a write plan:

1. **Enumerate** — List every file the skill will CREATE / MODIFY / DELETE with a one-line reason (test files, fix files, event log).
2. **Flags:**
   - `--dry-run` — Emit write plan only. Do NOT create, modify, or delete any file.
   - `--yes` — Skip confirmation prompt. Execute immediately. Intended for CI/non-interactive.
3. **Confirm** — Present the plan to the user and wait for explicit approval (unless `--yes`).
4. **Log** — Append write plan event to `.sdp/log/events.jsonl` (**sanitize file paths** before logging: strip newlines, ensure valid JSON escaping):
   ```json
   {"spec_version":"v1.0","event_id":"<uuid>","timestamp":"<ISO-8601>","source":{"system":"sdp-lab","component":"tdd"},"event_type":"decision.made","payload":{"decision_type":"write_plan","plan":[{"path":"...","action":"CREATE|MODIFY|DELETE","reason":"..."}]},"context":{"feature_id":"<F-id if known>","workstream_id":"<ws-id if applicable>"}}
   ```
   Include context fields only when the ID is known at plan time. Omit unavailable fields rather than inventing placeholders.
   > **Note:** Phase 1 uses prompt-level write boundaries (CLI out of scope). Aligns with `schema/contracts/orchestration-event.schema.json` via `event_type: "decision.made"`. Phase 2 CLI will emit natively.

**Output format:**
```
WRITE PLAN for @tdd <target>:
  CREATE: path/to/*_test.go — Failing test (RED phase)
  CREATE: path/to/impl.go — Minimal implementation (GREEN phase)
  MODIFY: .sdp/log/events.jsonl — Write plan event log

Proceed? [y/n]
```

**Modes:**
- No flag: Show plan → Confirm → Execute
- `--dry-run`: Show plan → STOP
- `--yes`: Show plan → Execute immediately (no prompt)

---

| Symptom | Fix |
|---------|-----|
| Skill produces no output | Check working directory is project root with `docs/workstreams/backlog/` |
| "checkpoint not found" | Run `sdp-orchestrate --feature <ID>` to create initial checkpoint |
| "workstream files missing" | Run `sdp-orchestrate --index` to verify, then `@feature` to regenerate |
| Skill hangs / no progress | Check `.sdp/log/events.jsonl` for last event; use `sdp reset --feature <ID>` if stuck |
| Review loop exceeds 3 rounds | Use `@review --override "reason"`, `@review --partial`, or `@review --escalate` |
