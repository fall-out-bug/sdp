---
name: tdd
description: Enforce Test-Driven Development: Red → Green → Refactor (INTERNAL - used by @build)
---

# @tdd (INTERNAL)

TDD discipline. Called by @build, not users.

## Cycle

1. **RED** — Write failing test first. Run: `go test ./...` — must FAIL
2. **GREEN** — Minimal implementation. Run: `go test ./...` — must PASS
3. **REFACTOR** — Improve code. Run: `go test ./...` — still PASS
4. **COMMIT** — Save state

## Exit When

- All AC met
- `go test ./...` passes
- `go vet ./...` passes

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
