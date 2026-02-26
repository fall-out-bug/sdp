# SDP Release Checklist

This checklist must be completed before any SDP release.

## Pre-Release Checks

### Engineering Quality

- [ ] All tests passing (`go test ./...`)
- [ ] Test coverage >= 80%
- [ ] No linter errors (`golangci-lint run`)
- [ ] No security vulnerabilities (`go vet`)
- [ ] All files < 200 LOC

### Documentation

- [ ] README.md updated
- [ ] CHANGELOG.md updated with version
- [ ] All command help text reviewed
- [ ] Breaking changes documented

### Build Verification

- [ ] Binary builds successfully (`go build ./...`)
- [ ] Cross-platform builds work (darwin, linux, windows)
- [ ] Version embedded correctly

## UX Quality Gates

### Onboarding Metrics

| Metric | Threshold | Check |
|--------|-----------|-------|
| TTFV | < 15 min | [ ] |
| Setup Completion | > 95% | [ ] |
| First Apply Success | 100% | [ ] |
| Discoverability Score | > 80% | [ ] |

### First-Run Verification

- [ ] `sdp demo` completes successfully on clean machine
- [ ] `sdp init --guided` passes all steps
- [ ] `sdp doctor` reports no errors
- [ ] `sdp status --text` shows meaningful output

### Help Quality

- [ ] `sdp --help` shows commands by intent
- [ ] Each command has examples
- [ ] Common journeys documented

## Post-Release Verification

- [ ] Binary download works
- [ ] Installation instructions correct
- [ ] Demo walkthrough verified on:
  - [ ] macOS
  - [ ] Linux
  - [ ] Windows (if supported)

## Rollback Criteria

If any of these occur after release, consider rollback:

1. Critical bug in guided setup
2. TTFV exceeds 20 minutes
3. Setup completion rate < 90%
4. > 3 critical support tickets in first week

## Sign-off

| Role | Name | Date |
|------|------|------|
| Engineering | | |
| QA | | |
| Product | | |

---

**Version:** 1.0
**Last Updated:** 2026-02-16
