---
name: ci-triage
description: Investigate failing GitHub Actions runs and produce root-cause plus Beads follow-up.
disable-model-invocation: true
---

# CI Triage

Use this skill when a branch or PR has failing CI and you need fast root-cause with actionable follow-up.

## Inputs

- PR number or URL (preferred)
- Or branch name (e.g. `feature/F063`)

## Workflow

### 1) Collect CI signal

```bash
gh pr checks <pr-number>
gh run list --branch <branch> --limit 20
```

### 2) Identify failing runs and jobs

```bash
gh run view <run-id>
gh run view <run-id> --log-failed
```

If logs are missing, query run metadata and check-suite details:

```bash
gh api repos/<owner>/<repo>/actions/runs/<run-id>
gh api repos/<owner>/<repo>/actions/runs/<run-id>/jobs
```

### 3) Classify root cause

Pick one primary class:

- Workflow config issue
- Test/regression failure
- Environment/dependency drift
- Permissions/secrets issue
- Flaky/infra transient

### 4) Propose smallest safe fix

- Point to exact workflow/job/step
- Provide command-level repro if possible
- Suggest one patch path and one fallback

### 5) Track in Beads

Create a follow-up item when CI is not green:

```bash
bd create --title="CI: <short root cause>" --type=bug --priority=1
bd dep add <new-id> <parent-feature-id>
bd sync
```

## Output Template

```markdown
## CI Triage Result

- Branch/PR: ...
- Failing workflow(s): ...
- Root cause: ...
- Evidence: run/job/step links
- Proposed fix: ...
- Beads follow-up: <id>
```
