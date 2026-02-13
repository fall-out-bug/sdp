# Quality Gates Pattern

> Mandatory checks before code is considered complete

## Pre-Commit Gates

| Gate | Command | Threshold |
|------|---------|-----------|
| Tests | `go test ./...` | 100% pass |
| Coverage | `go test -cover ./...` | >= 80% |
| Lint | `golangci-lint run` | 0 errors |
| Build | `go build ./...` | Success |

## Pre-PR Gates

| Gate | Command | Threshold |
|------|---------|-----------|
| All tests | `go test -race ./...` | 100% pass |
| Coverage report | `go test -coverprofile=c.out ./...` | >= 80% |
| Security scan | `gosec ./...` | 0 high/critical |
| Dependencies | `go mod tidy && go mod verify` | Clean |

## File Size Gate

```
Max LOC per file: 200
Warn at: 180
Block at: 200+
```

Check with:
```bash
find . -name "*.go" ! -name "*_test.go" -exec wc -l {} \; | \
  awk '$1 > 180 {print $2 ": " $1 " LOC"}'
```

## Architecture Gates

| Layer | Direction | Allowed |
|-------|-----------|---------|
| Domain | <- | Nothing (pure) |
| Application | <- | Domain |
| Infrastructure | <- | Domain, Application |
| Presentation | <- | All |

## Review Gates

Before `@review` can pass:

- [ ] All above gates pass
- [ ] No P0/P1 findings
- [ ] P2 findings tracked in beads
- [ ] Documentation updated
- [ ] Changelog updated (if user-facing)

## CI Gates

Automated in GitHub Actions:

1. Build (all platforms)
2. Test (with race detector)
3. Lint (golangci-lint)
4. Security (gosec, trivy)
5. Coverage (upload to codecov)

## Bypass Policy

Quality gates may ONLY be bypassed:

1. With explicit user approval
2. Documented reason in commit message
3. Follow-up issue created in beads
4. Time-boxed (must fix within 24h)

## See Also

- `@review` - Multi-agent quality review
- `.claude/patterns/tdd.md` - TDD pattern
- `docs/reference/PRINCIPLES.md` - Core principles
