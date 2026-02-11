# Debug Report: CI Dogfood — Duplicate Workflow Root Cause

**Date:** 2026-02-11  
**Method:** @debug (4-phase systematic debugging)  
**Result:** Root cause identified and fixed  

---

## Phase 1: OBSERVE

### Symptoms

- SDP Verify (Dogfood) CI stays red on PR #22
- `sdp-path` input appeared empty in action logs
- Download from GitHub Releases fails (no release exists)

### Evidence

- GitHub API: run 21899664208 used path **`.github/workflows/test-sdp-verify.yml`**
- Run 21899664208 job steps: Checkout, Set up Go, **Run SDP Verify** — no Build step
- Two workflow files with identical `name`, triggers (`on: pull_request/push` branches main, dev):
  - `sdp-verify-dogfood.yml` — has Build + sdp-path (our changes)
  - `test-sdp-verify.yml` — no Build, no sdp-path

### Diff

```diff
# test-sdp-verify.yml vs sdp-verify-dogfood.yml
- No Build SDP CLI step
- No sdp-path passed to action
- Action tries to download → fails (no release)
```

---

## Phase 2: HYPOTHESIZE

| Theory | Likelihood | Falsification |
|--------|------------|---------------|
| **steps.build.outputs not passed to composite action** | Low | Debug step would show value in workflow |
| **Duplicate workflow: test-sdp-verify runs without Build** | High | API shows run path = test-sdp-verify.yml |
| GITHUB_OUTPUT format wrong | Low | Build step in sdp-verify-dogfood works |

---

## Phase 3: EXPERIMENT

- Confirmed: run 21899664208 `path` = `test-sdp-verify.yml`
- Confirmed: that workflow has no Build step
- Both workflows trigger on same PR → two runs; the one from `test-sdp-verify.yml` fails

---

## Phase 4: CONFIRM

### Root Cause

**Duplicate workflow files.** `test-sdp-verify.yml` is a legacy/duplicate of `sdp-verify-dogfood.yml` with the same triggers. It lacks the Build step and sdp-path, so it tries to download SDP from Releases and fails.

### Fix

1. **Delete** `.github/workflows/test-sdp-verify.yml` (redundant)
2. **Remove** debug step from `sdp-verify-dogfood.yml` (no longer needed)

### Canonical Workflow

- `sdp-verify-dogfood.yml` — per F058 docs, the single dogfood workflow
- Builds SDP locally, passes `sdp-path` to action, skips download

---

## Lessons

- When CI is "still red" after fixes, check **which workflow file** actually ran (API `path` field)
- Duplicate workflows with same triggers cause confusion — consolidate or differentiate
