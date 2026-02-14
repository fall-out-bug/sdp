# F058: CI/CD Integration — GitHub Action

> Beads: sdp-4gvy | Priority: P1

---

## Problem

SDP verification runs locally during `@build`. No automated verification in CI/CD. AI-generated PRs merge without evidence chain validation.

## Solution

GitHub Action that runs SDP verification on every PR.

### GitHub Action

```yaml
# .github/workflows/sdp-verify.yml
name: SDP Verification
on: [pull_request]
jobs:
  verify:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: sdp-dev/verify-action@v1
        with:
          gates: [types, tests, coverage, evidence]
          evidence-required: true
          comment: true
```

### What It Does

1. **Verification gates**: runs configured quality checks (types, tests, semgrep, coverage)
2. **Evidence validation**: checks that evidence log exists and chain is intact
3. **Provenance gate**: blocks merge if AI-generated code lacks evidence records
4. **PR comment**: posts evidence chain summary as PR comment

### PR Comment Format

```markdown
## SDP Verification Report

| Gate | Status | Details |
|------|--------|---------|
| Types | ✅ Pass | 0 errors |
| Tests | ✅ Pass | 147 passed, 0 failed |
| Coverage | ✅ Pass | 87% (threshold: 80%) |
| Evidence | ✅ Pass | 12 events, chain intact |

### Evidence Summary
- 3 plan events (workstreams 00-054-01 to 00-054-03)
- 5 generation events (claude-sonnet-4)
- 3 verification events (all pass)
- 1 approval event (auto-approved)

Chain integrity: ✅ valid (12 events, no breaks)
```

## Constraints

- Action must work without SDP CLI pre-installed (self-contained or auto-install)
- Free tier: 50 runs/month per repo
- Evidence-required is optional (gradual adoption)
- GitLab CI variant as separate WS

## Users

- Teams using GitHub for AI-generated PRs
- Enterprise requiring automated audit trail validation
- Open source projects wanting AI code transparency

## Success Metrics

- Action runs in < 2 minutes on average PR
- Evidence chain validation catches at least 1 real issue in dogfooding
- 10+ repos using the action within 3 months of release

## Dependencies

- F054 (evidence log, chain validation)
- F057 (CLI commands — action calls `sdp verify` and `sdp log`)

## Notes

- Start with GitHub only — GitLab CI as separate workstream
- Action is open source — the verification logic is the open core
- Public OSS distribution encourages adoption; advanced signed-certificate workflows can stay outside this repository scope
